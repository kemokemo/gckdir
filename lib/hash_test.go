package lib

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestHashData(t *testing.T) {
	var tests = []struct {
		file string
		want HashData
	}{
		{filepath.Join("TestData", "HashData", "result_empty.json"), HashData{RelativePath: "aaa.txt", HashValue: "aaaa"}},
		{filepath.Join("TestData", "HashData", "result_true.json"), HashData{RelativePath: "bbb.txt", HashValue: "1234", VerifyResult: true}},
	}

	for _, test := range tests {
		data, err := json.MarshalIndent(test.want, "", "    ")
		if err != nil {
			t.Errorf(`MarshalIndent failed. %q`, err)
			return
		}

		bytes, err := ioutil.ReadFile(test.file)
		if err != nil {
			t.Errorf(`ReadFile failed. %q`, err)
			return
		}

		master := strings.TrimRight(string(bytes), "\n")
		testData := string(data)
		comp := strings.Compare(testData, master)
		if comp != 0 {
			t.Errorf(`Compare failed. strings.Compare returns %d`, comp)
		}
	}

}

func TestGenerateHashList(t *testing.T) {
	list, err := generateHashList(filepath.Join("TestData", "dir"))
	if err != nil {
		t.Errorf(`GenerateHashList(filepath.Join("TestData", "dir")) = %v, %q`, list, err)
	}

	count := len(list.List)
	if count == 0 {
		t.Errorf(`HashList.List has no item.`)
		return
	}
	if count > 1 {
		t.Errorf(`HashList.List has %d item.`, count)
		for _, item := range list.List {
			t.Errorf(`HashData is %v`, item)
		}
	}

	data := list.List[0]
	if data.RelativePath != "Test.txt" {
		t.Errorf(`HashList.List[0].FileName = %s`, data.RelativePath)
	}
	if data.HashValue != "532eaabd9574880dbf76b9b8cc00832c20a6ec113d682299550d7a6e0f345e25" {
		t.Errorf(`HashList.List[0].HashValue = %s`, data.HashValue)
	}
}

func TestVerifyHashList(t *testing.T) {
	master := HashList{List: []HashData{HashData{RelativePath: "test.bmp", HashValue: "aaaa"}}}
	same := HashList{List: []HashData{HashData{RelativePath: "test.bmp", HashValue: "aaaa"}}}
	otherHash := HashList{List: []HashData{HashData{RelativePath: "test.bmp", HashValue: "bbbb"}}}
	otherPath := HashList{List: []HashData{HashData{RelativePath: "sample.txt", HashValue: "aaaa"}}}
	morePath := HashList{List: []HashData{
		HashData{RelativePath: "test.bmp", HashValue: "aaaa"},
		HashData{RelativePath: "sample.txt", HashValue: "aaaa"}}}

	type args struct {
		source HashList
		target HashList
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
	}{
		{name: "Test-SameOne", args: args{source: master, target: same}, wantResult: true},
		{name: "Test-OtherOne1", args: args{source: master, target: otherHash}, wantResult: false},
		{name: "Test-OtherOne2", args: args{source: master, target: otherPath}, wantResult: false},
		{name: "Test-More", args: args{source: master, target: morePath}, wantResult: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyHashList(tt.args.source, tt.args.target); got.VerifyResult != tt.wantResult {
				t.Errorf("VerifyHashList() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}
