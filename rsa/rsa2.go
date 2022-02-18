package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"hash"
	"crypto/x509"
)

func Sign(data []byte,SignType,pemPriKey string) (sign string, err error) {
	var h hash.Hash
	var hType crypto.Hash
	switch SignType {
	case "RSA":
			h = sha1.New()
			hType = crypto.SHA1
	case "RSA2":
			h = sha256.New()
			hType = crypto.SHA256
	}
	h.Write(data)
	d := h.Sum(nil)
	pk,err := ParsePrivateKey(pemPriKey)
	if err != nil {
		return "", err
	}
	bs,err := rsa.SignPKCS1v15(rand.Reader,pk,hType,d)
	if err != nil {
		return "", err
	}
	sign = base64.StdEncoding.EncodeToString(bs)
	return sign,nil
}

func ParsePrivateKey(privateKey string) (pk *rsa.PrivateKey,err error) {
	block,_ := pem.Decode([]byte(privateKey))
	if block == nil {
		err = errors.New("privateKey 格式错误")
		return
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err == nil {
			pk = rsaPrivateKey
		}
	default:
		err = errors.New("privateKey 格式错误")	
	}
	return
}