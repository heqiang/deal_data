package util

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func Mkdir(dirPath string) {
	_, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(dirPath, 7777)
			if err != nil {
				msg := fmt.Sprintf("%s 目录创建失败,err:%s", dirPath, err)
				zap.L().Warn(msg)
				return
			}
		}
	}
}

func MkFile(filePath string) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fp, err := os.Create(filePath)
			if err != nil {
				msg := fmt.Sprintf("%s 文件创建失败,err:%s", filePath, err)
				zap.L().Warn(msg)
				return
			}
			defer fp.Close()
		}
	}
}

func fileDownload(filePath, fileName string, url, header, proxy string, id interface{}) {
	resBytes, err := Req(url, header, proxy)
	if err != nil {
		msg := fmt.Sprintf("文件请求失败fileUrl:%s,err:%s,新闻id或者uuid:%d", url, err, id)
		zap.L().Error(msg)
		return
	}

	imgFilePath := filepath.Join(filePath, fileName)
	file, err := os.OpenFile(imgFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		msg := fmt.Sprintf("%s 文件打开失败", imgFilePath)
		zap.L().Error(msg)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, _ = writer.Write(resBytes)
	_ = writer.Flush()
}
