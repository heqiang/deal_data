package util

import (
	"deal_data/config"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
)

func Mkdir(dirPath string) {
	_, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(dirPath, 0777)
			if err != nil {
				msg := fmt.Sprintf("%s 目录创建失败,err:%s", dirPath, err)
				zap.L().Warn(msg)
				return
			}
			os.Chmod(dirPath, 0777)
		}
	}
}

func MkFile(filePath string) {
	//确保文件目录存在
	fileDirList := strings.Split(filePath, string(filepath.Separator))
	fileDir := strings.Join(fileDirList[0:len(fileDirList)-1], string(filepath.Separator))
	Mkdir(fileDir)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fp, err := os.Create(filePath)
			if err != nil {
				fileSuffix := strings.Split(filePath, ".")[1]
				msg := fmt.Sprintf("%s 文件创建失败,err:%s", fileSuffix, err)
				zap.L().Warn(msg)
				return
			}
			defer fp.Close()

		}
	}
}

func fileDownload(filePath, fileName string, url, proxy string, id interface{}) {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()

	imgFilePath := filepath.Join(filePath, fileName)
	imageFile, err := os.OpenFile(imgFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 777)
	if err != nil {
		fmt.Println(err)
		msg := fmt.Sprintf("%s 文件创建失败", imgFilePath)
		zap.L().Error(msg)
		return
	}
	defer imageFile.Close()

	_, err = Req(url, proxy, imageFile)
	if err != nil {
		fmt.Println(err)
		msg := fmt.Sprintf("文件请求失败fileUrl:%s,err:%s,新闻id或者uuid:%d", url, err, id)
		zap.L().Error(msg)
		return
	}

}
