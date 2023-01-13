package db

import (
	"fmt"
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
	record := Record{
		Color: "dark",
		Title: "黑名单新增",
		Msg:   fmt.Sprintf("%s", blacklist.Blacklist),
	}

	AddRecord(record)

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
