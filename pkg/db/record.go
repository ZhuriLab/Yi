package db

import (
	"gorm.io/gorm"
	"time"
)

/**
  @author: yhy
  @since: 2023/1/13
  @desc: //TODO
**/

type Record struct {
	gorm.Model
	Id          int    `gorm:"primary_key" json:"id"`
	Project     string `json:"project"`
	Url         string `json:"url"`
	Color       string `json:"color"`
	Title       string `json:"title"`
	CurrentTime string `json:"current_time"`
	Msg         string `json:"msg"`
}

func AddRecord(record Record) {
	record.CurrentTime = time.Now().Format("2006-01-02 15:04:05")
	// 	record.CurrentTime = fmt.Sprintf("%s/%d %d:%d", t.Format("01"), t.Day(), t.Hour(), t.Minute())
	GlobalDB.Create(&record)
}

func GetRecord() (records []Record) {
	globalDBTmp := GlobalDB.Model(&Record{})
	globalDBTmp.Order("id desc").Find(&records)
	return
}
