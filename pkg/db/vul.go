package db

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

/**
  @author: yhy
  @since: 2022/10/18
  @desc: //TODO
**/

type Vul struct {
	gorm.Model
	Id            int    `gorm:"primary_key" json:"id"`
	Project       string `json:"project"`
	RuleId        string `json:"rule_id"`
	Url           string `json:"url"`
	DefaultBranch string `json:"default_branch"`
	PushedAt      string `json:"pushed_at"`
	Location      datatypes.JSON
	//CodeFlows     string `gorm:"type:text"`
	ResDir string `json:"res_dir"`
}

func AddVul(vul Vul) {
	if ExistBlacklist(vul.Location.String()) {
		return
	}
	GlobalDB.Create(&vul)
}

// GetVuls 查看漏洞信息
func GetVuls(pageNum int, pageSize int, maps interface{}) (count int64, vuls []Vul) {
	globalDBTmp := GlobalDB.Model(&Vul{})
	query := maps.(map[string]interface{})

	if query["project"] != nil {
		globalDBTmp = globalDBTmp.Where("project LIKE ?", "%"+query["project"].(string)+"%")
	}

	if query["rule_id"] != nil {
		globalDBTmp = globalDBTmp.Where("rule_id LIKE ?", "%"+query["rule_id"].(string)+"%")
	}

	globalDBTmp.Count(&count)
	globalDBTmp.Offset(pageNum).Limit(pageSize).Order("id asc").Find(&vuls)

	return
}

func DeleteVul(id string) {
	globalDBTmp := GlobalDB.Model(&Vul{})
	globalDBTmp.Where("id = ?", id).Unscoped().Delete(&Vul{})
}

// ExistVul  判断数据库中ip、端口是否存在
func ExistVul(id string) (bool, Vul) {
	var vul Vul
	globalDBTmp := GlobalDB.Model(&Vul{})
	globalDBTmp.Where("id = ? ", id).Limit(1).First(&vul)

	if vul.Id > 0 {
		return true, vul
	}

	return false, vul
}
