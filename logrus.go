package tools

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

//cfg 配置文件
type Cfg struct {
	LogLvl     string   // 日志级别
	EsAddrs    []string //ES addr
	EsUser     string   //ES user
	EsPassword string   //ES password
}

//setupLogrus 初始化logrus 同时把logrus的logger var 引用到这个common.Logger
func SetupLogrus(cc Cfg) error {
	//logFileName := fmt.Sprintf("%s_%s.log", os.Args[1], time.Now().Format("06_01_02T15_04_05"))
	//
	//f, err := os.Create(logFileName)
	//if err != nil {
	//	return err
	//}

	logLvl, err := logrus.ParseLevel(cc.LogLvl)
	if err != nil {
		return err
	}
	logrus.SetLevel(logLvl)
	//logrus.SetReportCaller(true)
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	//使用console默认输出

	//logrus.SetOutput(f)

	logrus.SetReportCaller(true)
	//开启 logrus ES hook
	esh := NewEsHook(cc)
	logrus.AddHook(esh)

	return nil
}

//esHook 自定义的ES hook
type EsHook struct {
	cmd    string // 记录启动的命令
	client *elastic.Client
}

//newEsHook 初始化
func NewEsHook(cc Cfg) *EsHook {
	es, err := elastic.NewClient(
		elastic.SetURL(cc.EsAddrs...),
		elastic.SetBasicAuth(cc.EsUser, cc.EsPassword),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(15*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "ES:", log.LstdFlags)),
		//elastic.SetInfoLog(log.New(os.Stdout, "ES:", log.LstdFlags)),
	)

	if err != nil {
		log.Fatal("failed to create Elastic V6 Client: ", err)
	}
	return &EsHook{client: es, cmd: strings.Join(os.Args, " ")}
}

//Fire logrus hook interface 方法
func (hook *EsHook) Fire(entry *logrus.Entry) error {
	doc := newEsLog(entry)
	doc["cmd"] = hook.cmd
	go hook.sendEs(doc)
	return nil
}

//Levels logrus hook interface 方法
func (hook *EsHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.WarnLevel,
	}
}

//sendEs 异步发送日志到es
func (hook *EsHook) sendEs(doc appLogDocModel) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("send entry to es failed: ", r)
		}
	}()
	_, err := hook.client.Index().Index(doc.indexName()).Type("_doc").BodyJson(doc).Do(context.Background())
	if err != nil {
		log.Println(err)
	}

}

//appLogDocModel es model
type appLogDocModel map[string]interface{}

func newEsLog(e *logrus.Entry) appLogDocModel {
	ins := map[string]interface{}{}
	for kk, vv := range e.Data {
		ins[kk] = vv
	}
	ins["time"] = time.Now().Local()
	ins["lvl"] = e.Level
	ins["message"] = e.Message
	ins["caller"] = fmt.Sprintf("%s:%d  %#v", e.Caller.File, e.Caller.Line, e.Caller.Func)
	return ins
}

// indexName es index name 时间分割
func (m *appLogDocModel) indexName() string {
	return "logrus-" + time.Now().Local().Format("2006-01-02")
}
