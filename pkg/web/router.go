package web

/**
  @author: yhy
  @since: 2022/12/6
  @desc: //TODO
**/

import (
	"Yi/pkg/db"
	"Yi/pkg/runner"
	"Yi/pkg/utils"
	"embed"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Vul struct {
	Id       int
	Project  string
	RuleId   string
	Url      string
	Location map[string]string
	PushedAt string
	ResDir   string
}

//go:embed static
var static embed.FS

//go:embed templates
var templates embed.FS

func Init() {
	gin.SetMode("release")
	router := gin.Default()

	// 静态资源加载
	//router.Static("/static", Path("static"))

	router.StaticFS("/static", http.FS(static))

	router.Static("/db/results/", "./db/results/")

	// 模板加载
	//templatesPath := Path("templates/*")
	//router.LoadHTMLGlob(templatesPath)

	// 设置模板资源
	router.SetHTMLTemplate(template.Must(template.New("").ParseFS(templates, "templates/*")))

	// basic 认证
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		runner.Option.UserName: runner.Option.Pwd,
	}))

	authorized.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/index")
	})

	authorized.GET("/index", func(c *gin.Context) {
		project := c.Query("project")
		language := c.Query("language")

		maps := make(map[string]interface{})

		if project != "" {
			maps["project"] = project
		}
		if language != "" {
			maps["language"] = language
		}

		pageSize, _ := strconv.Atoi(c.Query("pageSize"))

		if pageSize == 0 {
			pageSize = 20
		}

		result := 0
		page, _ := strconv.Atoi(c.Query("current"))
		if page == 0 {
			page = 1
		} else if page > 0 {
			result = (page - 1) * pageSize
		}

		total, data := db.GetProjects(result, pageSize, maps)

		for i, pro := range data {
			t, _ := time.Parse(time.RFC3339, pro.PushedAt)
			data[i].PushedAt = t.Format("2006-01-02 15:04:05")
			data[i].LastScanTime = data[i].UpdatedAt.Format("2006-01-02 15:04:05")
		}

		p := utils.NewPaginator(c.Request, pageSize, total)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"total":     total,
			"projects":  data,
			"paginator": p,
			"year":      time.Now().Year(),
			"msg":       db.Msg,
		})
	})

	authorized.GET("/addProject", func(c *gin.Context) {
		url := c.Query("url")
		tag := c.Query("tag")
		if url != "" {
			url = strings.TrimRight(url, "/")
			record := db.Record{
				Project: url,
				Url:     url,
				Color:   "success",
				Title:   url,
				Msg:     fmt.Sprintf("%s 添加成功, 正在生成数据库...", url),
			}
			db.AddRecord(record)

			go runner.ApiAdd(url, tag)
			c.Redirect(302, "/index")
		}

	})

	authorized.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.tmpl", gin.H{
			"year": time.Now().Year(),
			"msg":  db.Msg,
		})
	})

	authorized.GET("/record", func(c *gin.Context) {
		records := db.GetRecord()

		db.Msg = 0
		c.HTML(http.StatusOK, "record.tmpl", gin.H{
			"records": records,
			"year":    time.Now().Year(),
			"msg":     db.Msg,
		})
	})

	authorized.GET("/vul", func(c *gin.Context) {

		project := c.Query("project")
		rule_id := c.Query("rule_id")

		maps := make(map[string]interface{})

		if project != "" {
			maps["project"] = project
		}
		if rule_id != "" {
			maps["rule_id"] = rule_id
		}

		pageSize, _ := strconv.Atoi(c.Query("pageSize"))

		if pageSize == 0 {
			pageSize = 10
		}

		result := 0
		page, _ := strconv.Atoi(c.Query("current"))
		if page == 0 {
			page = 1
		} else if page > 0 {
			result = (page - 1) * pageSize
		}

		total, data := db.GetVuls(result, pageSize, maps)

		var vuls []Vul
		for _, vul := range data {
			location := make(map[string]string)
			for _, k := range jsoniter.Get(vul.Location).Keys() {
				location[k] = fmt.Sprintf("%s/blob/%s/%s", vul.Url, vul.DefaultBranch, k)
			}

			vuls = append(vuls, Vul{
				Id:       vul.Id,
				Project:  vul.Project,
				RuleId:   vul.RuleId,
				Url:      vul.Url,
				Location: location,
				PushedAt: vul.PushedAt,
				ResDir:   vul.ResDir,
			})
		}

		p := utils.NewPaginator(c.Request, pageSize, total)

		c.HTML(http.StatusOK, "vulnerability.tmpl", gin.H{
			"total":     total,
			"vuls":      vuls,
			"paginator": p,
			"year":      time.Now().Year(),
			"msg":       db.Msg,
		})
	})

	authorized.GET("/blacklist", func(c *gin.Context) {
		id := c.Query("id")

		exist, vul := db.ExistVul(id)
		if exist {
			db.DeleteVul(id)
			db.AddBlacklist(db.Blacklist{Blacklist: vul.Location.String()})
		}

		c.Redirect(302, "/vul")
	})

	authorized.GET("/del", func(c *gin.Context) {
		id := c.Query("id")
		db.DeleteVul(id)
		c.Redirect(302, "/vul")
	})

	authorized.GET("/download", func(c *gin.Context) {
		fileDir := c.Query("fileDir")

		f := strings.Split(fileDir, "/")

		fileName := f[len(f)-1]
		//打开文件
		_, errByOpenFile := os.Open(fileDir)
		//非空处理
		if errByOpenFile != nil {
			c.Redirect(http.StatusFound, "/404")
			return
		}
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Transfer-Encoding", "binary")
		c.File(fileDir)
		return
	})

	pprof.Register(router)
	router.Run(":" + runner.Option.Port)
}
