package server

import (
	"github.com/yushuailiu/ssh-auto/check"
	json "github.com/bitly/go-simplejson"
	"io/ioutil"
	"fmt"
	encodejson "encoding/json"
)

func Init() error {
	configFile := check.GetConfigFilePath()
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	servers = make(map[string]*Server)
	if len(b) == 0 {
		return nil
	}

	jsonList,err := json.NewJson(b)
	if err != nil {
		return err
	}
	arrayList,err := jsonList.Map()
	if err != nil {
		return err
	}
	for _, item := range arrayList {
		mapItem,ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("config file format error")
		}
		port,err := mapItem["port"].(encodejson.Number).Int64()
		if err != nil {
			return err
		}

		id,err := mapItem["id"].(encodejson.Number).Int64()
		if err != nil {
			return err
		}

		server := &Server{
			Id: int(id),
			Name: mapItem["name"].(string),
			Ip: mapItem["ip"].(string),
			Port:int(port),
			User:mapItem["user"].(string),
			Password:mapItem["password"].(string),
			LoginType:mapItem["login_type"].(string),
			SshFile:mapItem["pem_file"].(string),
		}
		servers[server.Name] = server
	}
	return nil
}