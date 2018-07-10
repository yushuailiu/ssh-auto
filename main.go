package main

import (
	"github.com/yushuailiu/ssh-auto/app"
	"github.com/yushuailiu/ssh-auto/check"
	"github.com/yushuailiu/ssh-auto/server"
	"github.com/yushuailiu/ssh-auto/print"
)

func main() {


	if !check.HasInit() {
		check.InitConfig()
	}
	
	err := server.Init()
	if err != nil {
		panic(err)
	}

	cliApp := &app.App{}

	err = cliApp.Run()

	if err != nil {
		print.Errorln(err.Error())
	}
}
