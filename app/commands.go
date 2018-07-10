package app

import (
	"github.com/urfave/cli"
	"github.com/yushuailiu/ssh-auto/server"
)

func getCommands() []cli.Command {

	return []cli.Command{
		{
			Name:	"add",
			Aliases:[]string{"a"},
			Usage: "add a server info",
			Action: server.Add,
		},
		{
			Name:	"delete",
			Aliases:[]string{"d"},
			Usage:	"delete a serve by id or name",
			Action: server.Delete,
		},
		{
			Name:	"edit",
			Aliases:[]string{"e"},
			Usage:	"edit a serve by id or name",
			Action: server.Edit,
		},
		{
			Name:	"list",
			Aliases:[]string{"l"},
			Usage:	"list all servers",
			Action: server.List,
			Flags:	[]cli.Flag{
				cli.BoolFlag{
					Name:	"detail, d",
					Usage:	"show detail of server",
				},
				cli.StringFlag{
					Name:	"filter, f",
					Usage:	"filter servers",
				},
			},
		},
		{
			Name:	"login",
			Usage:	"login a server by id",
			Action: server.Login,
		},
	}
}