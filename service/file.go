package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"ssh_backup/client"
	conf "ssh_backup/config"
	"ssh_backup/model"
	"ssh_backup/util"
	"strings"
	"time"

	"github.com/cilidm/toolbox/levelDB"
	"github.com/cilidm/toolbox/logging"
	"github.com/pkg/sftp"
)

// =============== save ===============
var (
	saveEnd     = make(chan bool)
	searchEnd   = make(chan bool)
	pathHandle  = make(chan string, conf.FileSetting.MaxChannelNum)
	maxReadChan = make(chan bool, conf.FileSetting.MaxChannelNum)
)

//func SaveEndMonitor() {
//	select {
//	case <-saveEnd:
//		for {
//			if len(maxReadChan) == 0 {
//				searchEnd <- true
//			}
//		}
//	}
//}

func GetFileToLDB() {
	begin := time.Now()
	cli := client.Instance()
	if cli == nil {
		fmt.Println("无法连接远程服务器")
		os.Exit(1)
	}
	dirPath := conf.FileSetting
	DirWalk(dirPath.Source, cli)
	select {
	case <-searchEnd:
		fmt.Println("任务结束，耗时", time.Since(begin))
	}
}

// 遍历远端文件夹
func DirWalk(dirPath string, client *sftp.Client) {
	files, err := client.Glob(filepath.Join(dirPath, "*"))
	if err != nil {
		logging.Error(err)
		return
	}
	except := strings.Split(conf.FileSetting.ExceptDir, ",")
	for _, v := range files {
		if util.IsContain(except, v) {
			continue
		}
		stat, err := client.Stat(v)
		if err != nil {
			logging.Error(err)
			return
		}
		if stat.IsDir() {
			DirWalk(v, client)
		} else {
			pathHandle <- v
		}
	}
	close(pathHandle)
}

func PathMonitor() {
	for {
		val, ok := <-pathHandle
		if !ok {
			if len(maxReadChan) == 0 {
				searchEnd <- true
				break
			}
		} else {
			maxReadChan <- true
			go SyncSaveFileInfo(val)
		}
	}
}

// 将文件下载并存入ldb
func SyncSaveFileInfo(v string) {
	if checkFileStatus(v) {
		<-maxReadChan
		return
	}
	s, _ := client.Instance().Stat(v)
	var file model.FileInfo
	file.FileName = s.Name()
	file.FileSize = s.Size()
	file.Status = model.Prossing
	file.FileSource = v
	file.FileTarget = util.GetNewPath(conf.FileSetting.Source, conf.FileSetting.Target, v)
	file.CreatedAt = time.Now().Format(util.FormatTime)
	SaveFileFromLDBHandler(file)
	<-maxReadChan
}

// 校验文件是否已存在
func checkFileStatus(v string) bool {
	has, _ := levelDB.GetServer().FindByKey(util.GetLdbKey(conf.FileSetting.Source, v))
	if string(has) != "" {
		var oldFile model.FileInfo
		json.Unmarshal(has, &oldFile)
		if oldFile.Status == model.Processed {
			logging.Warn(oldFile.FileName, "已存在")
			return true
		}
	}
	return false
}

// =============== read ===============
func SaveFileFromLDBHandler(v model.FileInfo) {
	fmt.Println("开始传输文件", v.FileName)
	dir, fileName := path.Split(v.FileTarget)
	if strings.ToLower(conf.Server.ServerType) == "ssh" {
		if err := UploadFile(client.Instance(), client.Instance(), dir, fileName, v.FileSource, v.FileSize); err != nil {
			v.Status = model.ProssErr
			levelDB.GetServer().Insert(util.GetLdbKey(conf.FileSetting.Source, v.FileSource), &v)
			logging.Error("upload file err :", err)
			return
		} else {
			v.Status = model.Processed
			levelDB.GetServer().Insert(util.GetLdbKey(conf.FileSetting.Source, v.FileSource), &v)
			return
		}
	} else {
		has, err := util.CheckFile(v.FileTarget)
		if has != nil && err == nil {
			if has.Size() == v.FileSize {
				v.Status = model.Processed
				levelDB.GetServer().Insert(util.GetLdbKey(conf.FileSetting.Source, v.FileSource), &v)
				return
			}
		}
		util.PathExists(dir)
		if err := GetFile(client.Instance(), v.FileSource, v.FileTarget); err != nil {
			v.Status = model.ProssErr
			levelDB.GetServer().Insert(util.GetLdbKey(conf.FileSetting.Source, v.FileSource), &v)
			logging.Error(err)
			return
		}
	}
	v.Status = model.Processed
	levelDB.GetServer().Insert(util.GetLdbKey(conf.FileSetting.Source, v.FileSource), &v)
}
