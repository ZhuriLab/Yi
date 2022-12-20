package db

import (
	"Yi/pkg/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"
)

/**
  @author: yhy
  @since: 2022/9/30
  @desc: //TODO
**/

var GlobalDB *gorm.DB

func init() {
	var err error
	GlobalDB, err = gorm.Open(sqlite.Open("CodeQLVul.db"), &gorm.Config{})
	if err != nil {
		logging.Logger.Errorln("failed to connect database")
	}

	// Migrate the schema
	err = GlobalDB.AutoMigrate(&Vul{}, &Project{}, &Blacklist{})
	if err != nil {
		logging.Logger.Errorf("db.Setup err: %v", err)
	}

	GlobalDB.Logger = dblogger.Default.LogMode(dblogger.Silent)
}
