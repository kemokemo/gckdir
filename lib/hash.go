package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ahmetb/go-linq"
)

// HashList is list of HashData struct
type HashList struct {
	List          []HashData `json:"hash_list"`
	CompareResult bool       `json:"comp_result,omitempty"`
}

// HashData is hash value and file name struct
type HashData struct {
	RelativePath  string `json:"relative_path"`
	HashValue     string `json:"hash_value"`
	CompareResult bool   `json:"comp_result,omitempty"`
	Reason        string `json:"false_reason,omitempty"`
}

// GenerateHashList generates hash information of dir path.
func GenerateHashList(dir string) (HashList, error) {
	list := HashList{}
	walk := func(path string, info os.FileInfo, err error) error {
		// root path is not used.
		if strings.Compare(dir, path) == 0 {
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			log.Println("Get relative path error.", err)
			return err
		}
		data := HashData{}
		data.RelativePath = rel

		if info.IsDir() {
			// directories do not have hash value
			data.HashValue = "-"
		} else {
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println("ReadFile error.", err)
				return err
			}
			hash := sha256.Sum256(bytes)
			data.HashValue = hex.EncodeToString(hash[:])
		}

		list.List = append(list.List, data)
		return nil
	}

	err := filepath.Walk(dir, walk)
	if err != nil {
		return list, err
	}

	return list, nil
}

// CompareHashList compares hash list of source and target.
// Then, returns HashList that has CompareResult.
func CompareHashList(source, target HashList) HashList {
	result := compareWithSource(source, target)
	result = compareWithTarget(result, target)
	sort.Slice(result.List, func(i int, j int) bool {
		return result.List[i].RelativePath < result.List[j].RelativePath
	})
	result.CompareResult = linq.From(result.List).All(func(arg1 interface{}) bool {
		return arg1.(HashData).CompareResult == true
	})
	return result
}

func compareWithSource(source, target HashList) HashList {
	result := HashList{CompareResult: false}

	for _, item := range source.List {
		selected := linq.From(target.List).
			SingleWith(func(c interface{}) bool {
				return c.(HashData).RelativePath == item.RelativePath
			})

		if selected == nil {
			var message string
			if item.HashValue == "-" {
				message = "This directory does not exist"
			} else {
				message = "This file does not exist"
			}

			fail := HashData{
				RelativePath:  item.RelativePath,
				HashValue:     item.HashValue,
				CompareResult: false,
				Reason:        message,
			}
			result.List = append(result.List, fail)
			log.Printf(`"%s" does not exist.`, item.RelativePath)
			continue
		}

		data := selected.(HashData)
		if data.HashValue == item.HashValue {
			data.CompareResult = true
		} else {
			data.CompareResult = false
			data.Reason = "Hash value does not match"
			log.Printf(`Hash value of "%s" does not match.`, item.RelativePath)
		}
		result.List = append(result.List, data)
	}

	return result
}

// compareWithTarget compares the result of comparison with the hash list
// of the target directory to check wheather the unnecessary items exist.
func compareWithTarget(result, target HashList) HashList {
	for _, item := range target.List {
		more := linq.From(result.List).
			SingleWith(func(c interface{}) bool {
				return c.(HashData).RelativePath == item.RelativePath
			})

		if more == nil {
			item.CompareResult = false
			item.Reason = "Unnecessary item exists."
			log.Printf(`Unnecessary item "%s" exists.`, item.RelativePath)
			result.List = append(result.List, item)
		}
	}
	return result
}
