package main

import (
	"Yi/pkg/runner"
	"Yi/pkg/web"
)

/**
  @author: yhy
  @since: 2022/10/13
  @desc: //TODO
**/

func main() {
	runner.ParseArguments()
	go runner.Cyclic()
	go runner.Run()

	web.Init()
}
