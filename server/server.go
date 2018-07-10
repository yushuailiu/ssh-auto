package server

import (
	"os"
	"encoding/json"
	"github.com/yushuailiu/ssh-auto/check"
	"github.com/yushuailiu/ssh-auto/print"
	"fmt"
)

var (
	servers map[string]*Server
)

type Server struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Ip string `json:"ip"`
	Port int `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	LoginType string `json:"login_type"`
	SshFile string `json:"pem_file"`
}

func (server *Server) IsLoginByPassword() bool {
	return server.LoginType == "password"
}

func saveServers() error {
	info,err := json.Marshal(servers)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(check.GetConfigFilePath(), os.O_WRONLY|os.O_TRUNC, 06044)
	if err != nil {
		return err
	}
	defer f.Close()

	_,err = f.Write(info)

	return err
}

func maxServerId() int {
	maxId := 0
	for _, server := range servers {
		if server.Id > maxId {
			maxId = server.Id
		}
	}
	return maxId
}

func generateNewServerId() int {
	maxId := maxServerId()
	return maxId + 1
}

func getServerById(id int) *Server {
	for _, server := range servers {
		if server.Id == id {
			return server
		}
	}
	return nil
}

func deleteServerById(id int)  {
	for name, server := range servers {
		if server.Id == id {
			print.Infoln(fmt.Sprintf("delete %s success", name))
			delete(servers, name)
		}
	}
}