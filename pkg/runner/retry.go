package runner

import (
	"Yi/pkg/db"
	"Yi/pkg/logging"
	"time"
)

/**
  @author: yhy
  @since: 2023/2/16
  @desc: 数据库生成错误的进行重试
**/

var IsRetry bool

// Retry  todo 不优雅
func Retry() {
	for {
		if IsRetry { // 运行完再进行重试, 就不用协程了，一个个跑
			for _, p := range RetryProject {
				logging.Logger.Printf("项目(%s)重试", p.Url)

				_, project := db.Exist(p.Url)

				if p.Code == 1 {
					// 从 github 获取
					_, dbPath, res := GetRepos(p.Url)
					project.DBPath = dbPath
					project.Language = res.Language
					project.PushedAt = res.PushedAt
					project.DefaultBranch = res.DefaultBranch
				} else if p.Code == 2 { // 手动生成
					project.DBPath = CreateDb(p.Url, project.Language)
				}

				db.UpdateProject(project.Id, project)
				if project.DBPath != "" {
					Exec(project, nil)
				}
			}
			RetryProject = nil
			RetryProject = make(map[string]ProError)
		}
		time.Sleep(time.Second)
	}
}
