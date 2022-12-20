package utils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

/**
  @author: yhy
  @since: 2022/10/9
  @desc: //TODO
**/

func WriteFile(fileName string, fileData string) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(fileData)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

// LoadFile content to slice
func LoadFile(filename string) (lines []string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Println("LoadFile err, ", err)
		return
	}
	defer f.Close() //nolint
	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Text() != "" {
			lines = append(lines, s.Text())
		}
	}
	return
}

func SaveFile(path string, data []byte) (err error) {
	// Remove file if exist
	_, err = os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			log.Println("旧文件删除失败", err.Error())
		}
	}

	// save file
	return ioutil.WriteFile(path, data, 0644)
}

// DeCompress 解压
func DeCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

// RemoveDir 删除 ./github 下所有的项目
func RemoveDir() {
	Pwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir, _ := ioutil.ReadDir(Pwd + "/github")
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"github", d.Name()}...))
	}
}

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
