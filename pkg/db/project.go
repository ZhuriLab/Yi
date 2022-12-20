package db

import (
	"gorm.io/gorm"
)

/**
  @author: yhy
  @since: 2022/12/7
  @desc: //TODO
**/

type Project struct {
	gorm.Model
	Id            int    `gorm:"primary_key" json:"id"`
	Project       string `json:"project"`
	Url           string `json:"url"`
	Language      string `json:"language"`
	DBPath        string `json:"db_path"`
	Count         int    `json:"count"`
	PushedAt      string `json:"pushed_at"`
	Vul           int    `json:"vul"`
	DefaultBranch string `json:"default_branch"`
	LastScanTime  string
}

func AddProject(project Project) (int, int) {
	GlobalDB.Create(&project)
	return project.Id, project.Count
}

// GetProjects 查看项目信息
func GetProjects(pageNum int, pageSize int, maps interface{}) (count int64, projects []Project) {
	globalDBTmp := GlobalDB.Model(&Project{})
	query := maps.(map[string]interface{})

	if query["project"] != nil {
		globalDBTmp = globalDBTmp.Where("project LIKE ?", "%"+query["project"].(string)+"%")
	}

	if query["language"] != nil {
		globalDBTmp = globalDBTmp.Where("language LIKE ?", "%"+query["language"].(string)+"%")
	}

	globalDBTmp.Count(&count)
	if pageNum == 0 && pageSize == 0 {
		globalDBTmp.Find(&projects)

	} else {
		globalDBTmp.Offset(pageNum).Limit(pageSize).Order("vul desc,count desc").Find(&projects)
	}
	return
}

// UpdateProjectArg 更新字段
func UpdateProjectArg(id int, arg string, count int) bool {
	globalDBTmp := GlobalDB.Model(&Project{})
	globalDBTmp.Where("id = ?", id).Update(arg, count)
	return true
}

func DeleteProject(id string) {
	globalDBTmp := GlobalDB.Model(&Project{})
	globalDBTmp.Where("id = ?", id).Unscoped().Delete(&Project{})
}

// Exist  判断数据库中ip、端口是否存在
func Exist(url string) (bool, Project) {
	var project Project
	globalDBTmp := GlobalDB.Model(&Project{})
	globalDBTmp.Where("url = ? ", url).Limit(1).First(&project)

	if project.Id > 0 {
		return true, project
	}

	return false, project
}

func UpdateProject(id int, project Project) {
	globalDBTmp := GlobalDB.Model(&Project{})
	globalDBTmp.Where("id = ?", id).Updates(project)
}
