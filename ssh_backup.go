package main

import (
	"github.com/cilidm/toolbox/levelDB"
	"github.com/cilidm/toolbox/logging"
	"ssh_backup/client"
	conf "ssh_backup/config"
	"ssh_backup/service"
)

func checkTargetDirs() {
	if conf.Server.ServerType == "ssh" {
		_, err := service.ClientPathExists(conf.FileSetting.Target, client.Instance())
		if err != nil {
			logging.Fatal("目标文件夹校验失败，请确定是否有权限写入", err.Error())
		}
	}
}

func initLevelDB() {
	err := levelDB.InitServer("runtime")
	if err != nil {
		logging.Fatal(err.Error())
	}
}

func init() {
	checkTargetDirs()
	initLevelDB()
	go service.PathMonitor()
}

func main() {
	service.GetFileToLDB()
	service.WriteHtmlLog()
	defer levelDB.GetServer().Close()
}
