package runner

import (
	"Yi/pkg/db"
	"Yi/pkg/utils"
	"sync"
)

/**
  @author: yhy
  @since: 2022/12/13
  @desc: //TODO
**/

// NewRules 规则更新时，从数据库中获取项目，用新规则跑一遍
func NewRules(oldQls *QLFile, newQls *QLFile) {
	goQLs := utils.Difference(oldQls.GoQL, newQls.GoQL)
	javaQls := utils.Difference(oldQls.JavaQL, newQls.JavaQL)

	globalDBTmp := db.GlobalDB.Model(&db.Project{})

	if len(goQLs) != 0 {
		var projects []db.Project
		globalDBTmp.Where("language = Go").Order("id asc").Find(&projects)
		scan(projects)
	}

	if len(javaQls) != 0 {
		var projects []db.Project
		globalDBTmp.Where("language = Java").Order("id asc").Find(&projects)
		scan(projects)
	}

}

func scan(projects []db.Project) {
	var wg sync.WaitGroup
	limit := make(chan bool, Option.Thread)

	for _, project := range projects {
		if project.DBPath == "" {
			continue
		}
		wg.Add(1)
		limit <- true
		Exec(project, nil)
		<-limit
		wg.Done()
	}

	wg.Wait()
	close(limit)
}
