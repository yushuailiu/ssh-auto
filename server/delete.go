package server

import (
	"github.com/urfave/cli"
	"strconv"
	"fmt"
	"github.com/theckman/go-flock"
	"github.com/yushuailiu/ssh-auto/check"
	"context"
	"time"
)

func Delete(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return fmt.Errorf("command error, command format: ssh-auto login [serverId]")
	}
	id, err := strconv.Atoi(c.Args().Get(0))

	if err != nil {
		return err
	}

	fileLock := flock.NewFlock(check.GetConfigFilePath())
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)

	defer cancel()

	locked, err := fileLock.TryLockContext(ctx, 2 * time.Second)
	if !locked || err != nil {
		fmt.Println(err)
		return fmt.Errorf("multi process is in editing status")
	}

	deleteServerById(id)

	saveServers()

	if fileLock.Locked() {
		// do work
		fileLock.Unlock()
	}

	return nil
}