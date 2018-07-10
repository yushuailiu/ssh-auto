package app

import (
	"github.com/urfave/cli"
	"os"
)

type App struct {
	
}

func (app *App) Run() (err error) {
	innerApp := cli.NewApp()
	innerApp.Commands = getCommands()
	err = innerApp.Run(os.Args)
	return
}