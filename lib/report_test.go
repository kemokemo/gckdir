package lib

import (
	"os"
	"testing"
)

func TestCreateReport(t *testing.T) {
	list := HashList{
		CompareResult: true,
		List: []HashData{
			HashData{RelativePath: "test.bmp", HashValue: "bbbb",
				CompareResult: true, Reason: ""},
		}}

	err := CreateReport(os.Stdout, list)
	if err != nil {
		t.Errorf(`CreateReport failed. %v`, err)
	}
}
