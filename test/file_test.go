package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/**
  @author: yhy
  @since: 2022/12/13
  @desc: //TODO
**/

var listfile []string //获取文件列表

func Listfunc(path string, f os.FileInfo, err error) error {
	var strRet string
	strRet, _ = os.Getwd()
	ostype := os.Getenv("GOOS") // windows, linux

	if ostype == "windows" {
		strRet += "\\"
	} else {
		strRet += "/"
	}

	if f == nil {
		return err
	}
	if f.IsDir() {
		return nil
	}

	strRet += path //+ "\r\n"

	ok := strings.HasSuffix(strRet, ".ql")
	if ok {
		listfile = append(listfile, strRet) //将目录push到listfile []string中
	}

	return nil
}

func getFileList(path string) string {
	err := filepath.Walk(path, Listfunc) //
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}

	return " "
}

func ListFileFunc(p []string) {
	for _, value := range p {
		fmt.Println(value)
	}
}

func TestFile(t *testing.T) {
	getFileList("/Users/yhy/CodeQL/codeql/java/ql/src/Security/CWE/CWE-074")
	ListFileFunc(listfile)
}
