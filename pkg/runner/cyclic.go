package runner

import (
	"Yi/pkg/db"
	"fmt"
	"os"
	"sync"
	"time"
)

/**
  @author: yhy
  @since: 2022/12/7
  @desc: 循环执行
**/

func Cyclic() {
	for {
		// todo 不够优雅，万一监控的项目过多，导致一天还没执行完呢
		// 等待24小时后再循环执行
		time.Sleep(24 * 60 * time.Minute)
		// 更新规则库
		UpdateRule()

		count := 0
		today := time.Now().Format("2006-01-02") + "/"
		DirNames = DirName{
			ZipDir:    Pwd + "/db/zip/" + today,
			ResDir:    Pwd + "/db/results/" + today,
			DbDir:     Pwd + "/db/database/" + today,
			GithubDir: Pwd + "/github/" + today,
		}
		os.MkdirAll(DirNames.ZipDir, 0755)
		os.MkdirAll(DirNames.ResDir, 0755)
		os.MkdirAll(DirNames.DbDir, 0755)
		os.MkdirAll(DirNames.GithubDir, 0755)

		var projects []db.Project
		globalDBTmp := db.GlobalDB.Model(&db.Project{})
		globalDBTmp.Order("id asc").Find(&projects)

		var wg sync.WaitGroup
		limit := make(chan bool, Option.Thread)

		for _, p := range projects {
			if p.DBPath == "" {
				continue
			}
			wg.Add(1)
			limit <- true
			go func(project db.Project) {
				defer func() {
					<-limit
					wg.Done()
				}()

				// 更新了才会去生成数据库
				update, dbPath, pushedAt := CheckUpdate(project.Url, project.PushedAt, project.Project)

				if !update {
					return
				}

				count++
				project.DBPath = dbPath
				project.PushedAt = pushedAt

				db.UpdateProject(project.Id, project)

				Exec(project, nil)
			}(p)

		}

		wg.Wait()
		close(limit)

		record := db.Record{
			Color: "primary",
			Title: "新一轮扫描",
			Msg:   fmt.Sprintf("新一轮扫描结束, 总共扫描了 %d 个项目", count),
		}
		db.AddRecord(record)
	}

}
