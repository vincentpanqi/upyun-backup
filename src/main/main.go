package main

import (
	"fmt"
	"os"
	"strings"
	"upyun"
)

func help() {
	var helpDoc = `Upyun Backup

Usage:
	Backup the data from upyun bucket

Commands:
	upyun_backup snapshot snapshotFile - Make a snapshot of all the folders and files in the bucket.
	upyun_backup backup snapshotFile - Start to backup files by the specified snapshot file.

Build Version:
  upyun_backup v1.0.20141030
`
	fmt.Println(helpDoc)
}

func main() {
	cmdArgs := os.Args
	if len(cmdArgs) != 3 {
		help()
		return
	}

	//load config
	var conf upyun.Conf
	var err error
	confFile := "upyun_backup.conf"
	conf, err = upyun.LoadConfig(confFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	//init logging
	logConfig := `{"filename" : "upyun_backup.log"}`
	upyun.InitLogs(logConfig, conf.Debug)

	//check domain
	domainConf := conf.Domain
	domain := upyun.DOMAIN_AUTO
	switch domainConf {
	case 0:
		domain = upyun.DOMAIN_AUTO
	case 1:
		domain = upyun.DOMAIN_DIANXIN
	case 2:
		domain = upyun.DOMAIN_LIANTONG
	case 3:
		domain = upyun.DOMAIN_YIDONG
	default:
		upyun.L.Warning("Invalid domain configuration, will use default")
	}

	//info message
	upyun.L.Informational("User: `%s'", conf.User)
	upyun.L.Informational("Password: `%s'", conf.Password)
	upyun.L.Informational("Bucket: `%s'", conf.Bucket)
	upyun.L.Informational("LocalDir: `%s'", conf.LocalDir)
	upyun.L.Informational("Domain: `%s'", domain)
	upyun.L.Informational("Routine: `%d'", conf.Routine)
	upyun.L.Informational("Debug: `%v'", conf.Debug)

	//execute the command
	cmdName := cmdArgs[1]
	snapFile := cmdArgs[2]

	backuper := upyun.UpyunBackup{
		Domain: domain,
	}
	switch strings.ToLower(cmdName) {
	case "snapshot":
		backuper.SnapshotFiles(conf, snapFile)
	case "backup":
		backuper.BackupFiles(conf, snapFile)
	default:
		upyun.L.Error("Unsupported command `%s'", cmdName)
	}

	upyun.L.Close()
}
