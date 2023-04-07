package runner

import (
	"Yi/pkg/db"
	"Yi/pkg/logging"
	"Yi/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/corpix/uarand"
	jsoniter "github.com/json-iterator/go"
	"github.com/thoas/go-funk"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

/**
  @author: yhy
  @since: 2022/12/7
  @desc: //TODO
**/

type GithubRes struct {
	Language      string `json:"language"`
	PushedAt      string `json:"pushed_at"`
	DefaultBranch string `json:"default_branch"`
}

// ProError 项目数据库获取错误的,进行重试
type ProError struct {
	Url  string
	Code int
}

var RetryProject = make(map[string]ProError)

// GetRepos 从 github 下载构建好的数据库
func GetRepos(url_tmp string) (error, string, GithubRes) {
	// https://github.com/prometheus/prometheus  -> https://api.github.com/repos/prometheus/prometheus
	guri := strings.ReplaceAll(url_tmp, "github.com", "api.github.com/repos")

	res := GetTimeBran(guri, url_tmp)
	// https://api.github.com/repos/grafana/grafana 这里只会显示项目中使用最多的语言，但并不一定是项目的主语言，比如这个显示TypeScript，但其实用了 Go 写的

	// repos 中的语言只是占比最多的语言种类，有可能 go 写的 typescript 占比最多
	res.Language = GetLanguage(guri, url_tmp) // todo 现在只是适配 Go,java 语言，后期尽量适配主流语言是，目前主力只看 Go 项目

	if res.Language != "" {
		guri = fmt.Sprintf("%s/code-scanning/codeql/databases/%s", guri, res.Language)
	} else {
		return errors.New("no Language"), "", res
	}

	err, dbPath, code := GetDb(guri, url_tmp, res.Language)
	if code != 0 { // 没有生成对应的数据库
		RetryProject[url_tmp] = ProError{
			Url:  url_tmp,
			Code: code,
		}
	}
	return err, dbPath, res
}

// GetTimeBran 获取项目的更新时间和主分支 https://api.github.com/repos/prometheus/prometheus
func GetTimeBran(guri, url_tmp string) GithubRes {
	req, _ := http.NewRequest("GET", guri, nil)
	req.Header.Set("Accept", "application/vnd.github.v3.text-match+json")
	if Option.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Option.Token))
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	res := GithubRes{}

	req.Close = true
	Option.Session.RateLimiter.Take()
	resp, err := Option.Session.Client.Do(req)

	if err != nil {
		logging.Logger.Errorln("GetRepos client.Do(req) err:", err)
		// 网络错误导致的，需要重试
		RetryProject[url_tmp] = ProError{
			Url:  url_tmp,
			Code: 1,
		}
		res.Language = ""
		return res
	}
	defer resp.Body.Close()

	if resp.Body != nil {
		result, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(result, &res)
	}
	return res
}

// GetLanguage 获取项目的代码语言  https://api.github.com/repos/prometheus/prometheus/languages
func GetLanguage(guri, url_tmp string) string {
	req, _ := http.NewRequest("GET", guri+"/languages", nil)
	req.Header.Set("Accept", "application/vnd.github.v3.text-match+json")
	if Option.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Option.Token))
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Close = true
	Option.Session.RateLimiter.Take()
	resp, err := Option.Session.Client.Do(req)

	if err != nil {
		logging.Logger.Errorln("GetLanguage client.Do(req) err:", err)
		// 网络错误导致的，需要重试
		RetryProject[url_tmp] = ProError{
			Url:  url_tmp,
			Code: 1,
		}
		return ""
	}

	defer resp.Body.Close()

	if resp.Body != nil {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return ""
		}
		var language string
		var m float64 = 1

		// 去除 HTML,TypeScript,JavaScript,CSS,SCSS 这些玩意，之后，该项目的编写语言就是使用最多的语言
		for k, v := range jsoniter.Get(body).GetInterface().(map[string]interface{}) {
			if funk.Contains(k, "HTML") || funk.Contains(k, "TypeScript") || funk.Contains(k, "JavaScript") || funk.Contains(k, "CSS") || funk.Contains(k, "SCSS") {
				continue
			}

			if v.(float64) > m {
				language = k
				m = v.(float64)
			}
		}
		if funk.Contains(Languages, language) {
			return language
		}
	}

	return ""
}

// GetDb 下载/生成 数据库  https://api.github.com/repos/prometheus/prometheus/code-scanning/codeql/databases/{languages}
/*
	0 : 成功
	1 : 网络或文件创建错误
	2 : github上没有生成对应的数据库
*/
func GetDb(guri, url, languages string) (error, string, int) {
	req, _ := http.NewRequest("GET", guri, nil)
	req.Header.Set("Accept", "application/zip")
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Content-Type", "application/octet-stream")
	if Option.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Option.Token))
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Close = true

	Option.Session.RateLimiter.Take()
	resp, err := Option.Session.Client.Do(req)

	if err != nil {
		logging.Logger.Errorln(guri, "HttpRequest Do err: ", err)
		return err, "", 1
	}
	defer resp.Body.Close()

	name := utils.GetName(url)
	filePath := DirNames.ZipDir + name + ".zip"

	var dbPath string
	if resp != nil && resp.StatusCode == 200 {
		out, err := os.Create(filePath)
		defer out.Close()
		if err != nil {
			logging.Logger.Errorln("os.Create(filePath) err:", err)
			return err, "", 1
		}

		if _, err = io.Copy(out, resp.Body); err != nil {
			logging.Logger.Errorln(url, " HttpRequest io.Copy err: ", err)
			return err, "", 1
		}
	} else { // 说明 该项目没有在 github配置 codeql 扫描(404)，或者项目所有者配置了访问需要权限(403)
		dbPath = CreateDb(url, languages)
		if dbPath == "" {
			return err, "", 2
		} else {
			return nil, dbPath, 0
		}
	}

	err = utils.DeCompress(filePath, DirNames.DbDir+name+"/")
	if err != nil {
		logging.Logger.Errorln("DeCompress err:", err)
		return err, "", 1
	}

	dbPath = filepath.Dir(path.Join(utils.CodeqlDb(DirNames.DbDir+name), "*"))
	logging.Logger.Debugln(url, " downloadDb success.")
	return nil, dbPath, 0
}

// CheckUpdate 检查项目是否更新
func CheckUpdate(project db.Project) (bool, string, string) {
	guri := strings.ReplaceAll(project.Url, "github.com", "api.github.com/repos")

	res := GetTimeBran(guri, project.Url)

	var (
		dbPath string
		code   int
	)

	if project.PushedAt < res.PushedAt { // 说明更新了

		guri = fmt.Sprintf("%s/code-scanning/codeql/databases/%s", guri, project.Language)

		_, dbPath, code = GetDb(guri, project.Url, project.Language)

		if code != 0 { // 没有生成对应的数据库
			RetryProject[project.Url] = ProError{
				Url:  project.Url,
				Code: code,
			}
		} else {
			// 是否在重试列表中
			delete(RetryProject, project.Url)
		}
	}

	if dbPath != "" {
		logging.Logger.Debugln(project.Url, " update, start a new scan.", dbPath)
		record := db.Record{
			Project: project.Project,
			Url:     project.Url,
			Color:   "warning",
			Title:   project.Project + " 更新",
			Msg:     fmt.Sprintf("%s 项目更新, 重新生成Codeql数据库", project.Url),
		}
		db.AddRecord(record)

		return true, dbPath, res.PushedAt
	}

	return false, "", ""
}
