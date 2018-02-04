package lib

import (
	"bytes"
	"html/template"
	"reflect"
	"testing"
)

func TestCreateReport(t *testing.T) {
	validResult := HashList{
		VerifyResult: true,
		List: []HashData{
			HashData{RelativePath: "test.bmp", HashValue: "bbbb", VerifyResult: true, Reason: ""},
		}}

	type args struct {
		paths  PathList
		result HashList
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "ValidReport", args: args{result: validResult, paths: PathList{"", ""}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := CreateReport(w, tt.args.paths, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("CreateReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_passFail(t *testing.T) {
	type args struct {
		result bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "PassCase", args: args{result: true}, want: "Pass"},
		{name: "FailCase", args: args{result: false}, want: "Fail"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := passFail(tt.args.result); got != tt.want {
				t.Errorf("passFail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rowAttr(t *testing.T) {
	type args struct {
		result bool
	}
	tests := []struct {
		name string
		args args
		want template.HTMLAttr
	}{
		{name: "SuccessCase", args: args{result: true}, want: template.HTMLAttr("success")},
		{name: "DangerCase", args: args{result: false}, want: template.HTMLAttr("danger")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rowAttr(tt.args.result); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rowAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}
