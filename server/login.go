package server

import (
	"strconv"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"golang.org/x/crypto/ssh/terminal"
	"fmt"
	"net"
	"github.com/urfave/cli"
)

func Login(c *cli.Context) error {
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
	auths, err := server.parseAuth()
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: server.User,
		Auth: auths,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client,err := ssh.Dial("tcp", server.getServerHost(), config)

	if err != nil {
		return err
	}
	defer client.Close()

	session,err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	defer terminal.Restore(fd, oldState)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}

	err = session.Shell()
	if err != nil {
		return err
	}

	err = session.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) getServerHost() string {
	return server.Ip + ":" + strconv.Itoa(server.Port)
}

func (server *Server)parseAuth() ([]ssh.AuthMethod, error) {
	auths := []ssh.AuthMethod{}
	switch server.LoginType {
	case "password":
		auths = append(auths, ssh.Password(server.Password))
		break
	case "ssh":
		authMethod, err := server.getSshAuth()
		if err != nil {
			return nil, err
		}
		auths = append(auths, authMethod)
		break
	default:
		return nil, fmt.Errorf("error password type")
	}
	return auths, nil
}

func (server *Server) getSshAuth() (ssh.AuthMethod, error) {
	bytesKey,err := ioutil.ReadFile(server.SshFile)
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	if server.Password == "" {
		signer, err = ssh.ParsePrivateKey(bytesKey)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(bytesKey, []byte(server.Password))
	}
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}