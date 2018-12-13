package EasySftp

import (
	"fmt"
	"testing"
)

func TestEasySftp_Connect(t *testing.T) {
	sftp := connect()
	if sftp == nil {
		t.Fail()
	}
	defer sftp.Close()
}

func connect() Sftp {
	sftp := FactorySftp(FEasySftp)
	var (
		user     = "root"
		password = "xx"
		host     = "192.168.0.200"
		port     = 22
	)
	if err := sftp.Connect(user, password, host, port); err != nil {
		fmt.Println(err)
		return nil
	}
	return sftp
}

func TestMultipleUploadFile(t *testing.T) {
	sftp := connect()
	defer sftp.Close()
	sftp.UploadFile("","")
}

func TestMultipleUploadFile2(t *testing.T) {
	sftp := connect()
	defer sftp.Close()
	MultipleUploadFile(sftp,"","",0)
}

func TestEasySftp_Mkdir(t *testing.T) {
	sftp := connect()
	defer sftp.Close()
	if err := sftp.Mkdir("");err!=nil{
		t.Fail()
	}
}

func TestEasySftp_DownloadFile(t *testing.T) {
	sftp := connect()
	defer sftp.Close()
	sftp.DownloadFile("","")
}


