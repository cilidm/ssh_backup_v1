package config

import (
	"github.com/go-ini/ini"
	"log"
	"os"
	"ssh_backup/util"
	"sync"
)

var (
	Cfg         *ini.File
	once        sync.Once
	FileSetting = &FileInfo{}
	Server      = &ServerInfo{}
)

type ServerInfo struct {
	ServerType string
	SourceHost string
	SourceUser string
	SourcePwd  string
	SourcePort int
	TargetHost string
	TargetUser string
	TargetPwd  string
	TargetPort int
}

type FileInfo struct {
	Source        string
	Target        string
	ExceptDir     string
	MaxChannelNum int
	SaveHtmlLog   bool
}

func init() {
	iniPath := "conf.ini"
	once.Do(func() {
		var err error
		Cfg, err = ini.Load(iniPath)
		if err != nil {
			log.Fatal("未发现配置文件")
			os.Exit(500)
		}
		err = Cfg.Section("server").MapTo(Server)
		util.CheckErr(err, "映射配置文件出错，请检查server配置")
		err = Cfg.Section("file_info").MapTo(FileSetting)
		util.CheckErr(err, "映射配置文件出错，请检查file_info配置")
	})
}
