package server

import (
	"github.com/urfave/cli"
	"github.com/yushuailiu/ssh-auto/print"
	"fmt"
	"net"
	"os"
	"strconv"
	"github.com/theckman/go-flock"
	"github.com/yushuailiu/ssh-auto/check"
	"context"
	"time"
	"golang.org/x/crypto/ssh/terminal"
)

func Add(c *cli.Context) error {

	fileLock := flock.NewFlock(check.GetConfigFilePath())
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)

	defer cancel()

	locked, err := fileLock.TryLockContext(ctx, 2 * time.Second)
	if !locked || err != nil {
		fmt.Println(err)
		return fmt.Errorf("multi process is in editing status")
	}

	server := &Server{}
	for len(server.Name) == 0 {
		print.Info("Enter a name for the server: ")
		fmt.Scanln(&server.Name)
		ok := checkServerName(server.Name)
		if !ok {
			print.Error(fmt.Sprintf("the name of %s has exist", server.Name))
			server.Name = ""
		}
	}

	for net.ParseIP(server.Ip) == nil {
		print.Info("enter Ip: ")
		fmt.Scanln(&server.Ip)
	}

	for server.Port == 0 {
		print.Info("enter Port: ")
		input := ""
		fmt.Scanln(&input)
		port,err := strconv.Atoi(input)
		if err != nil {
			continue
		}
		server.Port = port
	}

	for server.User == "" {
		print.Info("enter the User: ")
		fmt.Scanln(&server.User)
	}

	for server.LoginType != "password" && server.LoginType != "ssh" {
		print.Info("login type(password or ssh): ")
		fmt.Scanln(&server.LoginType)
	}

	if server.IsLoginByPassword() {
		for server.LoginType == "password" && server.Password == "" {
			print.Info("enter password: ")
			psd, _ := terminal.ReadPassword(0)
			server.Password = string(psd)
		}
	} else {
		print.Info("enter password(ssh login type can be empty): ")
		fmt.Scanln(&server.Password)

		for true {
			print.Info("enter the private key file path: ")
			fmt.Scanln(&server.SshFile)
			info, err := os.Stat(server.SshFile)
			if err == nil && !info.IsDir() {
				break
			}
			print.Errorln("the file not exist or is not a file, please enter the right private key file path.")
		}
	}

	err = server.save()
	if err != nil {
		return err
	}
	if fileLock.Locked() {
		// do work
		fileLock.Unlock()
	}

	if err != nil {
		return err
	}

	return nil
}

func (server *Server) save() error {
	if _, ok := servers[server.Name]; ok {
		return fmt.Errorf("server name has register")
	}
	server.Id = generateNewServerId()
	servers[server.Name] = server
	return saveServers()
}

func checkServerName(name string) bool {
	_,ok := servers[name]
	return !ok
}