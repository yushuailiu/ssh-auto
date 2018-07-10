package server

import (
	"github.com/urfave/cli"
	"strconv"
	"fmt"
	"os"
	"github.com/yushuailiu/ssh-auto/print"
	"net"
	"github.com/theckman/go-flock"
	"github.com/yushuailiu/ssh-auto/check"
	"context"
	"time"
	"golang.org/x/crypto/ssh/terminal"
)

func Edit(c *cli.Context) error {
	fileLock := flock.NewFlock(check.GetConfigFilePath())
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)

	defer cancel()

	locked, err := fileLock.TryLockContext(ctx, 2 * time.Second)
	if !locked || err != nil {
		fmt.Println(err)
		return fmt.Errorf("multi process is in editing status")
	}


	if len(c.Args()) != 1 {
		return fmt.Errorf("command error, command format: ssh-auto login [serverId]")
	}
	id, err := strconv.Atoi(c.Args().Get(0))

	if err != nil {
		return err
	}
	
	server := getServerById(id)
	
	if server == nil {
		return fmt.Errorf("has no server id is %d", id)
	}

	time := 0
	for time == 0 || len(server.Name) == 0 {
		print.Info(fmt.Sprintf("Enter a name for the server(default:%s): ", server.Name))
		var name string
		fmt.Scanln(&name)
		if len(name) == 0 || name == server.Name {
			break
		}

		ok := checkServerName(name)
		if !ok {
			print.Errorln(fmt.Sprintf("the name of %s has exist", server.Name))
			continue
		}
		server.Name = name
		time ++
	}

	time = 0
	for time == 0 || net.ParseIP(server.Ip) == nil {
		print.Info(fmt.Sprintf("enter Ip(default:%s): ", server.Ip))
		var ip string
		fmt.Scanln(&ip)

		if len(ip) == 0 {
			break
		}

		if net.ParseIP(server.Ip) == nil {
			print.Errorln(fmt.Sprintf("%s is not a valid ip", ip))
			continue
		}

		server.Ip = ip
		time ++
	}

	time = 0
	for time == 0 || server.Port <= 0 {
		print.Info(fmt.Sprintf("enter Port(default:%d): ", server.Port))
		input := ""
		fmt.Scanln(&input)
		if len(input) == 0 {
			break
		}

		port,err := strconv.Atoi(input)
		if err != nil {
			continue
		}
		server.Port = port
		time ++
	}

	time = 0
	for time == 0 || server.User == "" {
		var user string
		print.Info(fmt.Sprintf("enter the User(default:%s): ", server.User))
		fmt.Scanln(&user)
		if len(user) == 0 {
			break
		}

		server.User = user
		time ++
	}

	time = 0
	for time == 0 || (server.LoginType != "password" && server.LoginType != "ssh") {
		print.Info(fmt.Sprintf("login type(password or ssh)(default:%s): ", server.LoginType))
		var loginType string
		fmt.Scanln(&loginType)

		if len(loginType) == 0 {
			break
		}

		if loginType != "password" && server.LoginType != "ssh" {
			print.Errorln("enter login type must be password or ssh")
			continue
		}
		server.LoginType = loginType
		time ++
	}

	if server.IsLoginByPassword() {
		server.SshFile = ""
		time = 0
		for time == 0 || (server.LoginType == "password" && server.Password == "") {
			print.Info(fmt.Sprintf("enter password(default:%s): ", server.Password))
			var password string
			psd, _ := terminal.ReadPassword(0)
			password = string(psd)
			if len(password) == 0 {
				break
			}
			server.Password = password
			time ++
		}
	} else {
		print.Info("enter password(ssh login type can be empty): ")
		fmt.Scanln(&server.Password)

		for true {
			print.Info(fmt.Sprintf("enter the private key file path(default:%s): ", server.SshFile))
			var sshFile string
			fmt.Scanln(&sshFile)
			if len(sshFile) == 0 {
				break
			}

			info, err := os.Stat(sshFile)
			if err == nil && !info.IsDir() {
				server.SshFile = sshFile
				break
			}
			print.Errorln("the file not exist or is not a file, please enter the right private key file path.")
		}
	}

	saveServers()

	if fileLock.Locked() {
		// do work
		fileLock.Unlock()
	}

	return nil
}