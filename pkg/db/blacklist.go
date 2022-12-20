package db

import (
	"gorm.io/gorm"
)

/**
  @author: yhy
  @since: 2022/12/13
  @desc: //TODO
**/

type Blacklist struct {
	gorm.Model
	Id        int    `gorm:"primary_key" json:"id"`
	Blacklist string `json:"blacklist"`
}

func AddBlacklist(blacklist Blacklist) {
	GlobalDB.Create(&blacklist)
}

func ExistBlacklist(str string) bool {
	var blacklist Blacklist
	globalDBTmp := GlobalDB.Model(&Blacklist{})
	globalDBTmp.Where("Blacklist = ? ", str).First(&blacklist)

	if blacklist.Id > 0 {
		return true
	}

	return false
}
