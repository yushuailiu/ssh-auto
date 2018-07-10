package check

import (
	"os"
	"fmt"
	"os/user"
)

var (
	userHome = ""
	configDirectoryPath = "/.ssh-auto"
	configFileName = "/servers.ini"
)

func getConfigDirectoryPath() string {
	return userHome + configDirectoryPath
}

func GetConfigFilePath() string {
	return userHome + configDirectoryPath  + configFileName
}

func init()  {
	userInfo,err := user.Current()
	if err != nil {
		panic(err)
	}
	userHome = userInfo.HomeDir
}

func HasInit() bool {
	_, err := os.Stat(GetConfigFilePath())
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func InitConfig()  {
	if _, err := os.Stat(getConfigDirectoryPath()); os.IsNotExist(err) {
		fmt.Println(getConfigDirectoryPath())
		err = os.Mkdir(getConfigDirectoryPath(), os.ModePerm)
	}
	if _, err := os.Stat(GetConfigFilePath()); os.IsNotExist(err) {
		os.Create(GetConfigFilePath())
	}
	fmt.Println("init config")
}