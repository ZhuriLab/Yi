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
	"time"
)

/**
  @author: yhy
  @since: 2022/12/7
  @desc: //TODO
**/

type githubRes struct {
	Language      string `json:"language"`
	PushedAt      string `json:"pushed_at"`
	DefaultBranch string `json:"default_branch"`
}

// DownloadDb 从 github 下载构建好的数据库
func DownloadDb(url_tmp string, language string) (error, string, *githubRes) {
	// https://github.com/prometheus/prometheus  -> https://api.github.com/repos/prometheus/prometheus/code-scanning/codeql/databases/go
	guri := strings.ReplaceAll(url_tmp, "github.com", "api.github.com/repos")

	var res *githubRes

	if language == "" {
		res = GetRepos(guri)
		if res.Language != "" && utils.StringInSlice(res.Language, Languages) {
			guri = fmt.Sprintf("%s/code-scanning/codeql/databases/%s", guri, res.Language)
		} else {
			return errors.New("no Language"), "", res
		}
	}

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
		return err, "", res
	}
	defer resp.Body.Close()

	name := utils.GetName(url_tmp)
	filePath := DirNames.ZipDir + name + ".zip"

	if resp != nil && resp.StatusCode == 200 {
		out, err := os.Create(filePath)
		defer out.Close()
		if err != nil {
			logging.Logger.Errorln("os.Create(filePath) err:", err)
			return err, "", res
		}

		if _, err = io.Copy(out, resp.Body); err != nil {
			logging.Logger.Errorln(url_tmp, " HttpRequest io.Copy err: ", err)
			return err, "", res
		}
	} else { // 说明 该项目没有在 github配置 codeql 扫描(404)，或者项目所有者配置了访问需要权限(403)
		return err, "", res
	}

	err = utils.DeCompress(filePath, DirNames.DbDir+name+"/")
	if err != nil {
		logging.Logger.Errorln("DeCompress err:", err)
		return err, "", res
	}

	dbPath := filepath.Dir(path.Join(utils.CodeqlDb(DirNames.DbDir+name), "*"))
	logging.Logger.Debugln(url_tmp, " downloadDb success.")
	return nil, dbPath, res
}

// GetRepos 获取项目的代码语言和更新时间
func GetRepos(guri string) *githubRes {
	req, _ := http.NewRequest("GET", guri, nil)
	req.Header.Set("Accept", "application/vnd.github.v3.text-match+json")
	if Option.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Option.Token))
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	res := &githubRes{}

	req.Close = true
	Option.Session.RateLimiter.Take()
	resp, err := Option.Session.Client.Do(req)

	if err != nil {
		logging.Logger.Errorln("GetRepos client.Do(req) err:", err)
		res.Language = ""
		return res
	}
	defer resp.Body.Close()

	if resp.Body != nil {
		result, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(result, &res)
	}
	// https://api.github.com/repos/grafana/grafana/ 这里只会显示项目中使用最多的语言，但并不一定是项目的主语言，比如这个显示TypeScript，但其实用了 Go 写的

	// repos 中的语言只是占比最多的语言种类，有可能 go 写的 typescript 占比最多
	language := GetLanguage(guri) // todo 现在只是适配 Go,java 语言，后期尽量适配主流语言是，目前主力只看 Go 项目

	if utils.StringInSlice(language, Languages) {
		res.Language = language
	}

	return res
}

func GetLanguage(guri string) string {
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
		return ""
	}

	defer resp.Body.Close()

	if resp.Body != nil {

		body, _ := ioutil.ReadAll(resp.Body)
		results := jsoniter.Get(body).Keys()

		if funk.Contains(results, "Go") {
			return "Go"
		} else if funk.Contains(results, "Java") {
			return "Java"
		}
	}

	return ""
}

// CheckUpdate 检查项目是否更新
func CheckUpdate(url string, lastTime string, name string) (bool, string, string) {
	guri := strings.ReplaceAll(url, "github.com", "api.github.com/repos")
	res := GetRepos(guri)
	if res.Language != "" {
		guri = fmt.Sprintf("%s/code-scanning/codeql/databases/%s", guri, res.Language)
	} else {
		return false, "", ""
	}

	t := time.Date(2022, 12, 06, 16, 42, 28, 0, time.UTC)

	var dbPath string
	var err error
	if t.Format(lastTime) < t.Format(res.PushedAt) { // 说明更新了
		err, dbPath, _ = DownloadDb(url, res.Language)

		if err != nil || dbPath == "" {
			// 下载失败，有可能是对方没有使用 codeql 自动化扫描，所以没有创建数据库，这里直接拉取到本地，手动创建
			dbPath = CreateDb(url, res, name)
		}
	}

	if dbPath != "" {
		logging.Logger.Debugln(url, " update, start a new scan.", dbPath)

		record := db.Record{
			Project: name,
			Url:     url,
			Color:   "warning",
			Title:   name + " 更新",
			Msg:     fmt.Sprintf("%s 项目更新, 重新生成Codeql数据库", url),
		}
		db.AddRecord(record)

		return true, dbPath, res.PushedAt
	}

	return false, "", ""
}
