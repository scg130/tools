package tools

import (
	"fmt"
	"testing"
)

func TestCompareHashAndPasswd(t *testing.T) {
	type args struct {
		passwd string
		target string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		struct {
			name string
			args args
		}{
			name: "compare",
			args: args{
				passwd: "scg013012",
				target: "$2a$10$ugCOaLNdl.FIjcVg5FW4WO5fjubFJMgRglGnp9bIcnjDuXHERAG2G",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if CompareHashAndPasswd(tt.args.passwd, tt.args.target) {
				t.Errorf("success")
			} else {
				t.Errorf("failure")
			}
		})
	}
}

func TestGeneratePasswd(t *testing.T) {
	type args struct {
		passwd string
	}
	tests := []struct {
		name string
		args args
	}{
		struct {
			name string
			args args
		}{
			name: "test-generate-passwd",
			args: args{
				passwd: "scg013012",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePasswd(tt.args.passwd)
			if err != nil {
				t.Errorf("GeneratePasswd() error = %v", err)
				return
			}
			fmt.Println(got)
		})
	}
}
