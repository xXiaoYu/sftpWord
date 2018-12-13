package main

import (
	"fmt"
	"mods/EasySftp"
	"os"
	"path"
	"strings"
)

func main() {
	sftp := EasySftp.FactorySftp(EasySftp.FEasySftp)
	if err := sftp.Connect("root", "", "", 22); err != nil {
		fmt.Println(err)
		return
	}
	defer sftp.Close()
	var (
		remotePath       = "/root/test/" // 服务器项目路径
		localPath        = "C:/gowork/myApplication/ftpWork/test/"
		localProgramPath = "C:/gowork/myApplication/ftpWork/" // 本地项目路径
	)
	programPath := strings.Replace(path.Dir(localPath), localProgramPath, "", 1)
	remotePath = path.Join(remotePath, programPath)
	if f, _ := os.Stat(localPath); f.IsDir() {
		EasySftp.MultipleUploadFile(sftp,localPath,remotePath,0)
	} else {
		sftp.UploadFile(localPath, remotePath)
	}
	//programPath := strings.Replace(path.Dir(localPath),localProgramPath,"",1)
	//remotePath = remotePath+programPath
	//sftp.DownloadFile(localProgramPath, "/root/test/demo.txt")

}
