package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestMain(t *testing.M) {
	setup()
	exitCode := t.Run()
	os.Exit(exitCode)
}

func setup() {
	path := filepath.Join("TestData", "dir", "SubDir")
	_, err := os.Stat(path)
	if err == nil {
		return
	}

	err = os.Mkdir(path, 0777)
	if err != nil {
		fmt.Println("failed to make dir:", err)
		return
	}
}

func Test_HashList_GetDirectoryInfo(t *testing.T) {
	tests := []struct {
		name            string
		source          HashList
		wantDirectories int
		wantFiles       int
	}{
		{name: "normal-01",
			source: HashList{
				List: []HashData{
					HashData{RelativePath: "test.bmp", HashValue: "aaaa"},
					HashData{RelativePath: "data", HashValue: "-"},
					HashData{RelativePath: "data2", HashValue: "-"}}},
			wantDirectories: 2,
			wantFiles:       1,
		},
		{name: "normal-02",
			source: HashList{
				List: []HashData{
					HashData{RelativePath: "data", HashValue: "-"},
					HashData{RelativePath: "data2", HashValue: "-"}}},
			wantDirectories: 2,
			wantFiles:       0,
		},
		{name: "normal-03",
			source: HashList{
				List: []HashData{
					HashData{RelativePath: "test.bmp", HashValue: "aaaa"},
					HashData{RelativePath: "hoge.txt", HashValue: "hjsk"}}},
			wantDirectories: 0,
			wantFiles:       2,
		},
		{name: "nil",
			source:          HashList{},
			wantDirectories: 0,
			wantFiles:       0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDir, gotF := tt.source.GetDirectoryInfo(); gotDir != tt.wantDirectories || gotF != tt.wantFiles {
				t.Errorf("HashList.GetDirectoryInfo() = (%v, %v), want (%v, %v)", gotDir, gotF, tt.wantDirectories, tt.wantFiles)
			}
		})
	}
}

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

		bytes, err := os.ReadFile(test.file)
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

func TestVerifyHashList(t *testing.T) {
	master := HashList{List: []HashData{
		{RelativePath: "test.bmp", HashValue: "aaaa"},
		{RelativePath: "data", HashValue: "-"}}}
	same := HashList{List: []HashData{
		{RelativePath: "test.bmp", HashValue: "aaaa"},
		{RelativePath: "data", HashValue: "-"}}}
	otherHash := HashList{List: []HashData{
		{RelativePath: "test.bmp", HashValue: "bbbb"},
		{RelativePath: "data", HashValue: "-"}}}
	otherPath := HashList{List: []HashData{
		{RelativePath: "sample.txt", HashValue: "aaaa"},
		{RelativePath: "data", HashValue: "-"}}}
	otherDir := HashList{List: []HashData{
		{RelativePath: "test.bmp", HashValue: "aaaa"},
		{RelativePath: "hoge", HashValue: "-"}}}
	morePath := HashList{List: []HashData{
		{RelativePath: "test.bmp", HashValue: "aaaa"},
		{RelativePath: "data", HashValue: "-"},
		{RelativePath: "sample.txt", HashValue: "aaaa"}}}

	type args struct {
		source             HashList
		target             HashList
		doHashCheck        bool
		doUnnecessaryCheck bool
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
	}{
		{name: "Test-Same1", args: args{source: master, target: same, doHashCheck: true, doUnnecessaryCheck: true}, wantResult: true},
		{name: "Test-Same2", args: args{source: master, target: same, doHashCheck: false, doUnnecessaryCheck: true}, wantResult: true},
		{name: "Test-OtherHash1", args: args{source: master, target: otherHash, doHashCheck: true, doUnnecessaryCheck: true}, wantResult: false},
		{name: "Test-OtherHash2", args: args{source: master, target: otherHash, doHashCheck: false, doUnnecessaryCheck: true}, wantResult: true},
		{name: "Test-OtherPath1", args: args{source: master, target: otherPath, doHashCheck: true, doUnnecessaryCheck: true}, wantResult: false},
		{name: "Test-OtherPath2", args: args{source: master, target: otherPath, doHashCheck: false, doUnnecessaryCheck: true}, wantResult: false},
		{name: "Test-OtherDir1", args: args{source: master, target: otherDir, doHashCheck: true, doUnnecessaryCheck: true}, wantResult: false},
		{name: "Test-OtherDir2", args: args{source: master, target: otherDir, doHashCheck: false, doUnnecessaryCheck: true}, wantResult: false},
		{name: "Test-More1", args: args{source: master, target: morePath, doHashCheck: true, doUnnecessaryCheck: true}, wantResult: false},
		{name: "Test-More2", args: args{source: master, target: morePath, doHashCheck: false, doUnnecessaryCheck: true}, wantResult: false},
		{name: "Test-More3", args: args{source: master, target: morePath, doHashCheck: true, doUnnecessaryCheck: false}, wantResult: true},
		{name: "Test-More4", args: args{source: master, target: morePath, doHashCheck: false, doUnnecessaryCheck: false}, wantResult: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyHashList(tt.args.source, tt.args.target, tt.args.doHashCheck, tt.args.doUnnecessaryCheck); got.VerifyResult != tt.wantResult {
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
				HashData{RelativePath: "SubDir", HashValue: "-"},
				HashData{RelativePath: "SubDir2", HashValue: "-"},
				HashData{RelativePath: filepath.Join("SubDir2", "Test2.txt"), HashValue: "f9f2385a7d7cd1e6e9a801ab9bbbf7ae9998706af0c2f8b608a39b42ab94d88f"},
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
				HashData{RelativePath: "SubDir", HashValue: "-"},
				HashData{RelativePath: "SubDir2", HashValue: "-"},
				HashData{RelativePath: filepath.Join("SubDir2", "Test2.txt"), HashValue: "f9f2385a7d7cd1e6e9a801ab9bbbf7ae9998706af0c2f8b608a39b42ab94d88f"},
				HashData{RelativePath: "Test.txt", HashValue: "532eaabd9574880dbf76b9b8cc00832c20a6ec113d682299550d7a6e0f345e25"},
			}},
			wantErr: false,
			sortErr: false,
		},
		{name: "check-sort-01", args: args{dir: filepath.Join("TestData", "sample")},
			want: HashList{List: []HashData{
				HashData{RelativePath: "dir.json", HashValue: "6924024bc2bd0646a186723b9bd272dafdd1910eeabd8fe89462502e3f48b704"},
				HashData{RelativePath: "invalid.json", HashValue: "f9f2385a7d7cd1e6e9a801ab9bbbf7ae9998706af0c2f8b608a39b42ab94d88f"},
			}},
			wantErr: false,
			sortErr: false,
		},
		{name: "check-sort-02", args: args{dir: filepath.Join("TestData", "sample")},
			want: HashList{List: []HashData{
				HashData{RelativePath: "invalid.json", HashValue: "f9f2385a7d7cd1e6e9a801ab9bbbf7ae9998706af0c2f8b608a39b42ab94d88f"},
				HashData{RelativePath: "dir.json", HashValue: "6924024bc2bd0646a186723b9bd272dafdd1910eeabd8fe89462502e3f48b704"},
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
