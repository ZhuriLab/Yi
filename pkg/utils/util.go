package utils

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

/**
  @author: yhy
  @since: 2022/10/13
  @desc: //TODO
**/

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

func GetName(path string) string {
	if strings.HasSuffix(path, "/") {
		path = strings.TrimRight(path, "/")
	}

	ss := strings.Split(path, "/")

	return ss[len(ss)-1]
}

// WriteCounter 下载进度条
type WriteCounter struct {
	Total    uint64
	FileName string
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading [%s] ... %s complete", wc.FileName, humanize.Bytes(wc.Total))
}

func RandStr() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ%*.!@"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < 10; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

// CodeqlDb 获取文件夹下 codeql-database.yml 的上级目录路径
func CodeqlDb(dir string) string {
	var filePath string
	var walkFunc = func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.Contains(path, "codeql-database.yml") {
				filePath = path
				return nil
			}
		}
		return nil
	}
	filepath.Walk(dir, walkFunc)
	return filepath.Dir(filePath)
}

func StringInSlice(s string, slice []string) bool {
	if slice == nil {
		return false
	}
	sort.Strings(slice)
	index := sort.SearchStrings(slice, s)
	if index < len(slice) && strings.ToLower(slice[index]) == strings.ToLower(s) {
		return true
	}
	return false
}

// isSupportedProtocol checks given protocols are supported
func isSupportedProtocol(value string) bool {
	return value == "http" || value == "https" || value == "socks5"
}

func Difference(slice1, slice2 []string) []string {
	intersect := func(slice1, slice2 []string) []string {
		m := make(map[string]int)
		nn := make([]string, 0)
		for _, v := range slice1 {
			m[v]++
		}

		for _, v := range slice2 {
			times, _ := m[v]
			if times == 1 {
				nn = append(nn, v)
			}
		}
		return nn
	}

	m := make(map[string]int)
	nn := make([]string, 0)
	inter := intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 && value != "" {
			nn = append(nn, value)
		}
	}
	return nn
}
