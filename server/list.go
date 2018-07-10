package server

import (
	"github.com/urfave/cli"
	"strings"
	"github.com/yushuailiu/ssh-auto/print"
	"fmt"
)

func List(c *cli.Context) error {

	filter := c.String("filter")

	list := map[string]*Server{}

	detail := c.Bool("detail")

	if len(filter) != 0 {
		for name, server := range servers {
			if strings.Contains(name, filter) {
				list[name] = server
			}
		}
	} else {
		list = servers
	}

	showList(list, detail)


	return nil
}

func showList(list map[string]*Server, detail bool)  {
	for _, server := range list {
		if detail {
			print.Info(fmt.Sprintf("[%d] %s %s %d %s %s %s %s\n", server.Id, server.Name, server.Ip,
				server.Port, server.User, server.LoginType, server.Password, server.SshFile))
		} else {
			print.Info(fmt.Sprintf("[%d] %s %s\n", server.Id, server.Name, server.Ip))
		}
	}
}