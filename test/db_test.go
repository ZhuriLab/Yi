package test

import (
	"Yi/pkg/db"
	"Yi/pkg/runner"
	"Yi/pkg/utils"
	"Yi/pkg/web"
	"fmt"
	"testing"
)

/**
  @author: yhy
  @since: 2022/12/13
  @desc: //TODO
**/

func TestDb(t *testing.T) {
	maps := make(map[string]interface{})

	maps["language"] = "Go"

	aa, _ := db.GetProjects(0, 0, maps)

	fmt.Println(aa)
}

func TestWeb(t *testing.T) {
	runner.Option.Session = utils.NewSession("")
	runner.Option.UserName = "yhy"
	runner.Option.Pwd = "123"
	runner.Option.Port = "8888"
	web.Init()
}
