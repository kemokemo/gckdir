package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/ahmetb/go-linq"
)

// PathList is a list of paths
type PathList struct {
	SourcePath string
	TargetPath string
}

// HashList is a list of HashData struct
type HashList struct {
	List         []HashData `json:"hash_list"`
	VerifyResult bool       `json:"verify_result,omitempty"`
}

// GetDirectoryInfo returns the number of directories and files.
// If HashList.List is nil or empty, returns 0 and 0.
func (h *HashList) GetDirectoryInfo() (int, int) {
	directories := linq.From(h.List).CountWith(func(arg1 interface{}) bool {
		return arg1.(HashData).HashValue == "-"
	})
	files := len(h.List) - directories
	return directories, files
}

// HashData is hash value and file name struct
type HashData struct {
	RelativePath string `json:"relative_path"`
	HashValue    string `json:"hash_value"`
	VerifyResult bool   `json:"verify_result,omitempty"`
	Reason       string `json:"reason_of_failed,omitempty"`
}

// GetHashList returns the hash list of source.
// For source, a json file or a directory is supported.
func GetHashList(source string) (HashList, error) {
	source = filepath.Clean(source)
	list := HashList{}
	info, err := os.Stat(source)
	if err != nil {
		return list, err
	}
	if info.IsDir() {
		return generateHashList(source)
	}
	return readHashList(source)
}

// GenerateHashList generates hash information of dir path.
func generateHashList(dir string) (HashList, error) {
	type result struct {
		Data  HashData
		Error error
	}

	dir = filepath.Clean(dir)
	list := HashList{}
	c := make(chan result)
	var wg sync.WaitGroup
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// root path is not used.
		if strings.Compare(dir, path) == 0 {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			data := HashData{}
			rel, err := filepath.Rel(dir, path)
			if err != nil {
				log.Println("Get relative path error.", err)
				c <- result{Data: data, Error: err}
				return
			}
			data.RelativePath = rel

			if info.IsDir() {
				// directories do not have hash value
				data.HashValue = "-"
			} else {
				bytes, err := ioutil.ReadFile(path)
				if err != nil {
					log.Println("ReadFile error.", err)
					c <- result{Data: data, Error: err}
					return
				}
				hash := sha256.Sum256(bytes)
				data.HashValue = hex.EncodeToString(hash[:])
			}
			c <- result{Data: data, Error: nil}
		}()
		return nil
	})
	if err != nil {
		return list, err
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for r := range c {
		if r.Error != nil {
			log.Println(err)
		}
		list.List = append(list.List, r.Data)
	}
	sort.SliceStable(list.List, func(i int, j int) bool {
		return list.List[i].RelativePath < list.List[j].RelativePath
	})

	return list, nil
}

// ReadHashList reads source json file and return a hash list.
func readHashList(source string) (HashList, error) {
	list := HashList{}
	source = filepath.Clean(source)
	data, err := ioutil.ReadFile(source)
	if err != nil {
		return list, err
	}

	err = json.Unmarshal(data, &list)
	if err != nil {
		return list, err
	}
	return list, nil
}

// VerifyHashList verifies hash list of source and target.
// Then, returns HashList that has VerifyResult.
func VerifyHashList(source, target HashList, doHashCheck, doUnnecessaryCheck bool) HashList {
	result := verifyWithSource(source, target, doHashCheck)
	result = verifyWithTarget(result, target, doUnnecessaryCheck)
	sort.SliceStable(result.List, func(i int, j int) bool {
		return result.List[i].RelativePath < result.List[j].RelativePath
	})
	result.VerifyResult = linq.From(result.List).All(func(arg1 interface{}) bool {
		return arg1.(HashData).VerifyResult == true
	})
	return result
}

func verifyWithSource(source, target HashList, doHashCheck bool) HashList {
	result := HashList{VerifyResult: false}

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
				RelativePath: item.RelativePath,
				HashValue:    item.HashValue,
				VerifyResult: false,
				Reason:       message,
			}
			result.List = append(result.List, fail)
			log.Printf(`Required item does not exist. "%s"`, item.RelativePath)
			continue
		}

		data := selected.(HashData)
		if data.HashValue == item.HashValue {
			data.VerifyResult = true
		} else {
			if doHashCheck {
				data.VerifyResult = false
				data.Reason = "Hash value does not match"
				log.Printf(`Hash value does not match. "%s"`, item.RelativePath)
			} else {
				data.VerifyResult = true
				data.Reason = "Ignore that the hash values do not match"
			}
		}
		result.List = append(result.List, data)
	}

	return result
}

// verifyWithTarget verifies the result of verification with the hash list
// of the target directory to check wheather the unnecessary items exist.
func verifyWithTarget(result, target HashList, doUnnecessaryCheck bool) HashList {
	for _, item := range target.List {
		more := linq.From(result.List).
			SingleWith(func(c interface{}) bool {
				return c.(HashData).RelativePath == item.RelativePath
			})

		if more == nil && doUnnecessaryCheck {
			item.VerifyResult = false
			item.Reason = "Unnecessary item exists."
			log.Printf(`Unnecessary item exists. "%s"`, item.RelativePath)
			result.List = append(result.List, item)
		} else if more == nil && !doUnnecessaryCheck {
			item.VerifyResult = true
			item.Reason = "Ignore that an unnecessary item exists."
			result.List = append(result.List, item)
		}
	}
	return result
}
