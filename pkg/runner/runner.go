package runner

import (
	"Yi/pkg/db"
	"Yi/pkg/logging"
	"Yi/pkg/utils"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"os"
	"sync"
)

/**
  @author: yhy
  @since: 2022/10/13
  @desc: //TODO
**/

func Run() {
	logging.Logger.Infoln("Yi Starting ... ")

	projects := make(chan db.Project, Option.Thread)

	go func() {
		if Option.Target != "" {
			exist, project := db.Exist(Option.Target)
			if exist {
				projects <- project
				return
			}
			name := utils.GetName(Option.Target)
			err, dbPath, res := DownloadDb(Option.Target, "")

			project = db.Project{
				Project:       name,
				DBPath:        dbPath,
				Url:           Option.Target,
				Language:      res.Language,
				PushedAt:      res.PushedAt,
				Count:         0,
				DefaultBranch: res.DefaultBranch,
			}

			if err != nil || dbPath == "" {
				// 下载失败，有可能是对方没有使用 codeql 自动化扫描，所以没有创建数据库，这里直接拉取到本地，手动创建
				dbPath = CreateDb(Option.Target, res, name)
			}

			if dbPath == "" { // 即使不能扫描，也加入数据库
				db.AddProject(project)
				return
			}

			projects <- project

		} else if Option.Targets != "" {
			targets := utils.LoadFile(Option.Targets)
			limit := make(chan bool, Option.Thread)
			var wg sync.WaitGroup

			for _, target := range targets {
				limit <- true
				wg.Add(1)
				go func(target string) {
					defer func() {
						wg.Done()
						<-limit
					}()

					exist, project := db.Exist(target)
					if exist {
						projects <- project
						return
					}
					name := utils.GetName(target)
					err, dbPath, res := DownloadDb(target, "")

					if err != nil || dbPath == "" {
						// 下载失败，有可能是对方没有使用 codeql 自动化扫描，所以没有创建数据库，这里直接拉取到本地，手动创建
						dbPath = CreateDb(target, res, name)
					}

					if res == nil {
						res.Language = ""
						res.PushedAt = ""
					}

					project = db.Project{
						Project:       name,
						DBPath:        dbPath,
						Url:           target,
						Language:      res.Language,
						PushedAt:      res.PushedAt,
						DefaultBranch: res.DefaultBranch,
						Count:         0,
					}

					if dbPath != "" {
						projects <- project
					} else {
						db.AddProject(project)
					}

				}(target)
			}
			wg.Wait()
			close(limit)
		}

		close(projects)
	}()

	limit := make(chan bool, Option.Thread)
	var wg sync.WaitGroup

	for project := range projects {
		if project.DBPath == "" {
			continue
		}
		wg.Add(1)
		limit <- true
		go WgExec(project, &wg, limit)
	}

	wg.Wait()
	close(limit)
}

func WgExec(project db.Project, wg *sync.WaitGroup, limit chan bool) {
	exist, p := db.Exist(project.Url)

	if exist {
		project.Id = p.Id
		project.Count = p.Count
	} else {
		id, count := db.AddProject(project)
		project.Id = id
		project.Count = count
	}

	Exec(project, nil)

	<-limit
	wg.Done()
}

var LocationMaps = make(map[string]bool)

func Exec(project db.Project, qls []string) {
	if !utils.StringInSlice(project.Language, Languages) {
		return
	}
	for fileName, res := range Analyze(project.DBPath, project.Project, project.Language, qls) {
		results := jsoniter.Get([]byte(res), "runs", 0, "results")

		if results.Size() > 0 {
			for i := 0; i < results.Size(); i++ {
				location := "{"

				msg := "ruleId: " + results.Get(i).Get("ruleId").ToString() + "\t "
				if results.Get(i).Get("locations").Size() > 0 {
					for j := 0; j < results.Get(i).Get("locations").Size(); j++ {
						msg += "locations: " + results.Get(i).Get("locations").Get(j).Get("physicalLocation", "artifactLocation", "uri").ToString() + "\t startLine: " + results.Get(i).Get("locations").Get(j).Get("physicalLocation", "region", "startLine").ToString() + "\t | "

						line := results.Get(i).Get("locations").Get(j).Get("physicalLocation", "region", "startLine").ToString()
						location += fmt.Sprintf("\"%s#L%s\":\"%s\",", results.Get(i).Get("locations").Get(j).Get("physicalLocation", "artifactLocation", "uri").ToString(), line, line)
					}
				}

				if results.Get(i).Get("relatedLocations").Size() > 0 {
					for j := 0; j < results.Get(i).Get("relatedLocations").Size(); j++ {
						msg += "relatedLocations: " + results.Get(i).Get("relatedLocations").Get(j).Get("physicalLocation", "artifactLocation", "uri").ToString() + "\t startLine: " + results.Get(i).Get("relatedLocations").Get(j).Get("physicalLocation", "region", "startLine").ToString() + "\t | "

						line := results.Get(i).Get("relatedLocations").Get(j).Get("physicalLocation", "region", "startLine").ToString()
						location += fmt.Sprintf("\"%s#L%s\":\"%s\",", results.Get(i).Get("relatedLocations").Get(j).Get("physicalLocation", "artifactLocation", "uri").ToString(), line, line)
					}
				}

				location += "}"

				if _, ok := LocationMaps[location]; ok {
					continue
				}

				vul := db.Vul{
					Project:  project.Project,
					RuleId:   results.Get(i).Get("ruleId").ToString(),
					Location: []byte(location),
					//CodeFlows:     results.Get(i).Get("codeFlows").ToString(),
					Url:           project.Url,
					ResDir:        fileName,
					PushedAt:      project.PushedAt,
					DefaultBranch: project.DefaultBranch,
				}

				db.AddVul(vul)

				db.UpdateProjectArg(project.Id, "vul", 1)
				logging.Logger.Infof("%s(%s) Found: %s", project.Project, fileName, msg)
			}
		} else {
			err := os.Remove(fileName) //删除文件

			if err != nil {
				logging.Logger.Infof("file remove Error! %s\n", err)
			}
		}
	}
	db.UpdateProjectArg(project.Id, "count", project.Count+1)
}

func ApiAdd(target string) {
	var exist bool
	var project db.Project
	exist, project = db.Exist(target)
	if !exist {
		name := utils.GetName(target)
		err, dbPath, res := DownloadDb(target, "")
		if err != nil {
			return
		}
		project = db.Project{
			Project:       name,
			DBPath:        dbPath,
			Url:           target,
			Language:      res.Language,
			PushedAt:      res.PushedAt,
			DefaultBranch: res.DefaultBranch,
			Count:         0,
		}
		id, _ := db.AddProject(project)
		project.Id = id
	}

	Exec(project, nil)
}
