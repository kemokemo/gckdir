package lib

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
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
	master := HashList{List: []HashData{
		HashData{RelativePath: "test.bmp", HashValue: "aaaa"},
		HashData{RelativePath: "data", HashValue: "-"}}}
	same := HashList{List: []HashData{
		HashData{RelativePath: "test.bmp", HashValue: "aaaa"},
		HashData{RelativePath: "data", HashValue: "-"}}}
	otherHash := HashList{List: []HashData{
		HashData{RelativePath: "test.bmp", HashValue: "bbbb"},
		HashData{RelativePath: "data", HashValue: "-"}}}
	otherPath := HashList{List: []HashData{
		HashData{RelativePath: "sample.txt", HashValue: "aaaa"},
		HashData{RelativePath: "data", HashValue: "-"}}}
	otherDir := HashList{List: []HashData{
		HashData{RelativePath: "test.bmp", HashValue: "aaaa"},
		HashData{RelativePath: "hoge", HashValue: "-"}}}
	morePath := HashList{List: []HashData{
		HashData{RelativePath: "test.bmp", HashValue: "aaaa"},
		HashData{RelativePath: "data", HashValue: "-"},
		HashData{RelativePath: "sample.txt", HashValue: "aaaa"}}}

	type args struct {
		source      HashList
		target      HashList
		doHashCheck bool
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
	}{
		{name: "Test-Same1", args: args{source: master, target: same, doHashCheck: true}, wantResult: true},
		{name: "Test-Same2", args: args{source: master, target: same, doHashCheck: false}, wantResult: true},
		{name: "Test-OtherHash1", args: args{source: master, target: otherHash, doHashCheck: true}, wantResult: false},
		{name: "Test-OtherHash2", args: args{source: master, target: otherHash, doHashCheck: false}, wantResult: true},
		{name: "Test-OtherPath1", args: args{source: master, target: otherPath, doHashCheck: true}, wantResult: false},
		{name: "Test-OtherPath2", args: args{source: master, target: otherPath, doHashCheck: false}, wantResult: false},
		{name: "Test-OtherDir1", args: args{source: master, target: otherDir, doHashCheck: true}, wantResult: false},
		{name: "Test-OtherDir2", args: args{source: master, target: otherDir, doHashCheck: false}, wantResult: false},
		{name: "Test-More1", args: args{source: master, target: morePath, doHashCheck: true}, wantResult: false},
		{name: "Test-More2", args: args{source: master, target: morePath, doHashCheck: false}, wantResult: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyHashList(tt.args.source, tt.args.target, tt.args.doHashCheck); got.VerifyResult != tt.wantResult {
				t.Errorf("VerifyHashList() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}

func TestGetHashList(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name    string
		args    args
		want    HashList
		wantErr bool
	}{
		{name: "sample-data", args: args{source: filepath.Join("TestData", "sample", "dir.json")},
			want: HashList{List: []HashData{
				HashData{RelativePath: "Test.txt", HashValue: "532eaabd9574880dbf76b9b8cc00832c20a6ec113d682299550d7a6e0f345e25"},
			}},
			wantErr: false},
		{name: "invalid-file", args: args{source: filepath.Join("TestData", "sample", "dir2.json")},
			want:    HashList{},
			wantErr: true},
		{name: "sample-dir", args: args{source: filepath.Join("TestData", "dir")},
			want: HashList{List: []HashData{
				HashData{RelativePath: "Test.txt", HashValue: "532eaabd9574880dbf76b9b8cc00832c20a6ec113d682299550d7a6e0f345e25"},
			}},
			wantErr: false},
		{name: "invalid-dir", args: args{source: filepath.Join("TestData", "dir2")},
			want:    HashList{},
			wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHashList(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHashList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHashList() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_readHashList(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name    string
		args    args
		want    HashList
		wantErr bool
	}{
		{name: "sample-data", args: args{source: filepath.Join("TestData", "sample", "dir.json")},
			want: HashList{List: []HashData{
				HashData{RelativePath: "Test.txt", HashValue: "532eaabd9574880dbf76b9b8cc00832c20a6ec113d682299550d7a6e0f345e25"},
			}},
			wantErr: false},
		{name: "invalid-file", args: args{source: filepath.Join("TestData", "sample", "dir2.json")},
			want:    HashList{},
			wantErr: true},
		{name: "invalid-data", args: args{source: filepath.Join("TestData", "sample", "invalid.json")},
			want:    HashList{},
			wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readHashList(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("readHashList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readHashList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateHashList(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    HashList
		wantErr bool
		sortErr bool
	}{
		{name: "sample-dir", args: args{dir: filepath.Join("TestData", "dir")},
			want: HashList{List: []HashData{
				HashData{RelativePath: "Test.txt", HashValue: "532eaabd9574880dbf76b9b8cc00832c20a6ec113d682299550d7a6e0f345e25"},
			}},
			wantErr: false,
			sortErr: false,
		},
		{name: "check-sort-01", args: args{dir: filepath.Join("TestData", "sample")},
			want: HashList{List: []HashData{
				HashData{RelativePath: "dir.json", HashValue: "27ef37cea442f72b3209e769cd88537f967c14fe9c744d654d2232cb6483eeb8"},
				HashData{RelativePath: "invalid.json", HashValue: "f9f2385a7d7cd1e6e9a801ab9bbbf7ae9998706af0c2f8b608a39b42ab94d88f"},
			}},
			wantErr: false,
			sortErr: false,
		},
		{name: "check-sort-02", args: args{dir: filepath.Join("TestData", "sample")},
			want: HashList{List: []HashData{
				HashData{RelativePath: "invalid.json", HashValue: "f9f2385a7d7cd1e6e9a801ab9bbbf7ae9998706af0c2f8b608a39b42ab94d88f"},
				HashData{RelativePath: "dir.json", HashValue: "27ef37cea442f72b3209e769cd88537f967c14fe9c744d654d2232cb6483eeb8"},
			}},
			wantErr: false,
			sortErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateHashList(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateHashList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			g, err := json.MarshalIndent(got, "    ", "")
			if err != nil {
				t.Errorf("failed to marshal json %v", err)
			}
			w, err := json.MarshalIndent(tt.want, "    ", "")
			if err != nil {
				t.Errorf("failed to marshal json %v", err)
			}
			if !reflect.DeepEqual(string(g), string(w)) != tt.sortErr {
				t.Errorf("generateHashList() = %v, w %v", string(g), string(w))
			}
		})
	}
}
