package service

import (
	"encoding/json"
	"github.com/cilidm/toolbox/gconv"
	"github.com/cilidm/toolbox/levelDB"
	"github.com/cilidm/toolbox/logging"
	"os"
	"path"
	conf "ssh_backup/config"
	"ssh_backup/model"
	"ssh_backup/template"
	"ssh_backup/util"
)

func WriteHtmlLog() {
	if conf.FileSetting.SaveHtmlLog == false {
		return
	}
	ldb := levelDB.GetServer()
	val, err := ldb.FindByPrefix(util.GetLdbPreKey(conf.FileSetting.Source))
	if err != nil {
		logging.Error(err)
		return
	}
	var files []model.FileInfoJson
	for _, v := range val {
		var (
			f        model.FileInfo
			fileJson model.FileInfoJson
		)
		if err := json.Unmarshal([]byte(v), &f); err != nil {
			logging.Error(err)
			return
		}
		util.CopyFields(&fileJson, f)
		fileJson.FileSize = util.SizeFormat(gconv.Float64(f.FileSize))
		files = append(files, fileJson)
	}
	jsonStr, _ := json.Marshal(files)
	WriteHtmlLogHandler(string(jsonStr))
}

func WriteHtmlLogHandler(jsonStr string) {
	_, dirName := path.Split(conf.FileSetting.Source)
	fp := "./" + dirName + ".html"
	str := template.Top
	str += template.Script
	str += jsonStr
	str += template.Bottom
	fd, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	fd.WriteString(str)
	fd.Close()
}
