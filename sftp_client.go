package EasySftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	FEasySftp = "EasySftp"
)

type Sftp interface {
	Connect(user, password, host string, port int) (error)
	UploadFile(localFilePath string, remotePath string)
	DownloadFile(localPath string, remoteFilePath string)
	Mkdir(remoteFilePath string) (error)
	Close()
}

type EasySftp struct {
	*sftp.Client
}

// 简单工厂创建对应的sftp结构
func FactorySftp(factoryName string) (sftp Sftp) {
	switch factoryName {
	case FEasySftp:
		sftp = &EasySftp{}
	}
	return sftp
}

// 连接
func (p *EasySftp) Connect(user, password, host string, port int) (error) {
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	clientConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	sshClient, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return err
	}
	p.Client, err = sftp.NewClient(sshClient)
	return err
}

const depthMax = 6 //大于设定的深度
// 多文件上传
// localPath 本地需要上传文件夹或文件
// remotePath 远程需要上传到的文件夹
// depth 文件夹层数，一般写0，用于计算文件夹层数
func MultipleUploadFile(sftp Sftp, localPath string, remotePath string, depth int) {
	if depth > depthMax {
		return
	}
	files, err := ioutil.ReadDir(localPath)
	if err != nil {
		return
	}
	for _, file := range files {
		realLocalPath := path.Join(localPath, file.Name())
		if file.IsDir() {
			sftp.Mkdir(path.Join(remotePath, file.Name()))
			MultipleUploadFile(sftp, realLocalPath, remotePath, depth+1)
			continue
		} else {
			sftp.UploadFile(realLocalPath, remotePath)
		}
	}
}

// 单文件文件上传
// localPath 本地需要上传的文件
// remotePath 远程需要上传到的文件夹
func (p *EasySftp) UploadFile(localPath string, remotePath string) {
	srcFile, err := os.Open(localPath)
	if err != nil {
		fmt.Println("os.Open error : ", localPath)
		log.Fatal(err)
	}
	defer srcFile.Close()
	remoteFileName := filepath.Base(localPath)
	dstFile, err := p.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		log.Fatal(err)
	}
	defer dstFile.Close()
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localPath)
		log.Fatal(err)
	}
	dstFile.Write(ff)
	fmt.Println(localPath + "  copy file to remote server finished!")
}

// 单文件文件下载
// localPath 本地需要下载到的文件夹
// remotePath 远程需要下载文件
func (p *EasySftp) DownloadFile(localPath string, remoteFilePath string) {
	srcFile, err := p.Open(remoteFilePath)
	if err != nil {
		fmt.Println("sftpClient.open error : ", remoteFilePath)
		log.Fatal(err)
	}
	defer srcFile.Close()
	var localFileName = path.Base(remoteFilePath)
	dstFile, err := os.Create(path.Join(localPath, localFileName))
	if err != nil {
		fmt.Println("os.Open error : ", localFileName)
		log.Fatal(err)
	}
	defer dstFile.Close()
	if _, err = srcFile.WriteTo(dstFile); err != nil {
		log.Fatal(err)
		fmt.Println("download error : ", remoteFilePath)
	}
	fmt.Println(remoteFilePath, "download file to remote server finished!")
}

// 远程创建文件夹
func (p *EasySftp) Mkdir(remotePath string) error {
	return p.Client.Mkdir(remotePath)
}

func (p *EasySftp) Close() {
	p.Client.Close()
}
