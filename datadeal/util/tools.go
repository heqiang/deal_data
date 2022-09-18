package util

import (
	"deal_data/global"
	"fmt"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
	"time"
)

type Tool struct {
	header      string
	dateTimeStr string
	time        string
	filePath    string
	imgPath     string
	jsonFile    string
	proxy       string
}

func NewTool(proxy string) *Tool {
	currTime := time.Now()
	return &Tool{
		header:      global.Header,
		dateTimeStr: currTime.Format("20060102150405"),
		time:        currTime.Format("2006-01-02"),
		proxy:       proxy,
	}
}

func (tool *Tool) ProfileImgDownload(imgBasePath, fileName, url string) {
	if !strings.HasPrefix(url, "http") {
		msg := fmt.Sprintf("%s:no schema", url)
		zap.L().Error(msg)
	}

	ProfilePathSuffix := fmt.Sprintf("%s_profile_image", tool.dateTimeStr)
	completeProfilePath := filepath.Join(imgBasePath, ProfilePathSuffix)

	Mkdir(completeProfilePath)

	resBytes, err := Req(url, tool.header, tool.proxy)
	if err != nil {
		msg := fmt.Sprintf("头像图片请求失败:%s", url)
		zap.L().Error(msg)
		return
	}
	fileDownload(completeProfilePath, fileName, resBytes)
}

func (tool *Tool) ImgOrFileDownload(filePath, fileName, url, fileType string, newsId int) {
	filePathSuffix := fmt.Sprintf("%s_%s", tool.dateTimeStr, fileType)
	completeFilePath := filepath.Join(filePath, filePathSuffix)

	Mkdir(completeFilePath)

	resBytes, err := Req(url, tool.header, tool.proxy)
	if err != nil {
		msg := fmt.Sprintf("文件请求失败fileUrl:%s,err:%s,新闻id:%d", url, err, newsId)
		zap.L().Error(msg)
		return
	}
	fileDownload(completeFilePath, fileName, resBytes)
}
