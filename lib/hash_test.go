package lib

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
)

func TestHashData(t *testing.T) {
	var tests = []struct {
		file string
		want HashData
	}{
		{"TestData\\HashData\\result_empty.json", HashData{RelativePath: "aaa.txt", HashValue: "aaaa"}},
		{"TestData\\HashData\\result_true.json", HashData{RelativePath: "bbb.txt", HashValue: "1234", CompareResult: true}},
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
	list, err := GenerateHashList("TestData\\dir")
	if err != nil {
		t.Errorf(`GenerateHashList("TestData\\dir") = %v, %q`, list, err)
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

func TestCompareHashList(t *testing.T) {
	master := HashList{
		List: []HashData{
			HashData{RelativePath: "test.bmp", HashValue: "aaaa"},
		}}
	same := HashList{
		List: []HashData{
			HashData{RelativePath: "test.bmp", HashValue: "aaaa"},
		}}
	otherHash := HashList{
		List: []HashData{
			HashData{RelativePath: "test.bmp", HashValue: "bbbb"},
		}}
	otherPath := HashList{
		List: []HashData{
			HashData{RelativePath: "sample.txt", HashValue: "aaaa"},
		}}
	morePath := HashList{
		List: []HashData{
			HashData{RelativePath: "test.bmp", HashValue: "aaaa"},
			HashData{RelativePath: "sample.txt", HashValue: "aaaa"},
		}}

	// same one
	result := CompareHashList(same, master)
	if result.CompareResult == false {
		t.Errorf(`CompareHashList(same, master) = false`)
	}

	// other one
	result = CompareHashList(otherHash, master)
	if result.CompareResult == true {
		t.Errorf(`CompareHashList(otherHash, master) = true`)
	}

	result = CompareHashList(otherPath, master)
	if result.CompareResult == true {
		t.Errorf(`CompareHashList(otherPath, master) = true`)
	}

	// more
	result = CompareHashList(morePath, master)
	if result.CompareResult == true {
		t.Errorf(`CompareHashList(morePath, master) = true`)
	}
}
