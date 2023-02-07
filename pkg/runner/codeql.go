package runner

import (
	"Yi/pkg/db"
	"Yi/pkg/logging"
	"Yi/pkg/utils"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

/**
  @author: yhy
  @since: 2022/10/13
  @desc: //TODO
**/

func Analyze(database string, name string, language string, qls []string) map[string]string {
	if qls == nil {
		if strings.ToLower(language) == "go" {
			qls = QLFiles.GoQL
		} else if strings.ToLower(language) == "java" {
			qls = QLFiles.JavaQL
		}
	}

	res := make(map[string]string)
	filePath := DirNames.ResDir + name
	os.MkdirAll(filePath, 0755)

	logging.Logger.Infof("[[%s:%s]] analyze start ...", name, database)
	fileName := fmt.Sprintf("%s/%d.json", filePath, time.Now().Unix())
	for i, ql := range qls {
		cmd := exec.Command("codeql", "database", "analyze", "--rerun", database, Option.Path+ql, "--format=sarif-latest", "-o", fileName)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout // 标准输出
		cmd.Stderr = &stderr // 标准错误
		err := cmd.Run()
		_, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		if err != nil {
			logging.Logger.Errorf("Analyze cmd.Run() failed with %s --  %s, %s %s", err, errStr, database, name)
			continue
		}

		lines := utils.LoadFile(fileName)

		var result string

		for _, i := range lines {
			result += i
		}
		res[fileName] = result

		ProgressBar[name] = float32(i) / float32(len(qls)) * 100
	}

	logging.Logger.Infof("[[%s:%s]] analysis completed.", name, database)
	record := db.Record{
		Project: name,
		Url:     name,
		Color:   "success",
		Title:   name,
		Msg:     fmt.Sprintf("%s 分析完毕", name),
	}
	db.AddRecord(record)
	return res
}

// CreateDb 拉取仓库，本地创建数据库
func CreateDb(gurl string, res *githubRes, name string) string {
	if !utils.StringInSlice(res.Language, Languages) {
		return ""
	}
	dbName := utils.GetName(gurl)
	err := GitClone(gurl, dbName)

	if err != nil {
		logging.Logger.Errorln("create db err:", err)
		return ""
	}

	// todo 批量跑就抽风，导致有的项目无法生成数据库 "There's no CodeQL extractor named 'Go' installed."
	cmd := exec.Command("codeql", "database", "create", DirNames.DbDir+dbName, "-s", DirNames.GithubDir+dbName, "--language="+strings.ToLower(res.Language), "--overwrite")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err = cmd.Run()
	out, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		logging.Logger.Errorf("CreateDb cmd.Run() failed with %s\n %s --  %s\n", err, out, errStr)
		return ""
	}

	// 很奇怪，有的生成数据库不是在项目目录下，而是在第二级目录下
	dbPath := filepath.Dir(path.Join(utils.CodeqlDb(DirNames.DbDir+dbName), "*"))
	return dbPath
}

// UpdateRule 每天拉取一下官方仓库，更新规则
func UpdateRule() {
	if Option.Path != "" {
		utils.RunGitCommand(Option.Path, "git", "pull")
	}
}
