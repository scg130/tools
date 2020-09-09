package tools

import (
	"fmt"
	"testing"
)

func TestValidToken(t *testing.T) {
	type args struct {
		tokenStr string
		key      string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		struct {
			name    string
			args    args
			want    interface{}
			wantErr bool
		}{
			name: "valid",
			args: args{
				tokenStr: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEYXRhIjoiZGF0YSIsImV4cCI6MTU5OTYzODg2MH0.xvzW2t1h1Hfa7BIvHcOVWXWBuZIZQnmjex3hbJi7iFo",
				key:      "key",
			},
			want:    "want",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidToken(tt.args.tokenStr, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	type args struct {
		key  string
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		struct {
			name    string
			args    args
			want    string
			wantErr bool
		}{
			name:    "generate",
			args:    args{key: "key", data: "data"},
			want:    "want",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateToken(tt.args.key, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}
