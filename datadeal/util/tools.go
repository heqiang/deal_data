package util

import (
	"bufio"
	"deal_data/global"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Tool struct {
	header       string
	dateTimeStr  string
	time         string
	filePath     string
	imgPath      string
	jsonPathName string
	proxy        string
}

func NewTool(proxy string, filePath string) *Tool {
	currTime := time.Now()
	tool := &Tool{
		header:      global.Header,
		dateTimeStr: currTime.Format("20060102150405"),
		time:        currTime.Format("2006-01-02"),
		proxy:       proxy,
	}
	jsonName := fmt.Sprintf("%s_newsty.json", tool.dateTimeStr)
	tool.jsonPathName = filepath.Join(filePath, jsonName)
	MkFile(tool.jsonPathName)
	return tool
}

func (tool *Tool) ProfileImgDownload(imgBasePath, fileName, url string, uuid string) {
	if !strings.HasPrefix(url, "http") {
		msg := fmt.Sprintf("%s:no schema", url)
		zap.L().Error(msg)
	}

	ProfilePathSuffix := fmt.Sprintf("%s_profile_image", tool.dateTimeStr)
	completeProfilePath := filepath.Join(imgBasePath, ProfilePathSuffix)

	Mkdir(completeProfilePath)

	fileDownload(completeProfilePath, fileName, url, tool.header, tool.proxy, uuid)
}

func (tool *Tool) ImgOrFileDownload(filePath, fileName, url, fileType string, newsId int) {
	filePathSuffix := fmt.Sprintf("%s_%s", tool.dateTimeStr, fileType)
	completeFilePath := filepath.Join(filePath, filePathSuffix)

	Mkdir(completeFilePath)

	fileDownload(completeFilePath, fileName, url, tool.header, tool.proxy, newsId)
}

func (tool *Tool) WriteNewsToJson(news map[string]interface{}, newsId int) {

	file, err := os.OpenFile(tool.jsonPathName, os.O_APPEND|os.O_CREATE, 7777)
	if err != nil {
		zap.L().Error(fmt.Sprintf("新闻写入json文件失败,err:%s,新闻id：%d", err, newsId))
		return
	}
	strNews, err := json.Marshal(news)
	if err != nil {
		zap.L().Error(fmt.Sprintf("新闻序列化失败,err:%s,新闻id：%d", err, newsId))
		return
	}
	writer := bufio.NewWriter(file)
	_, _ = writer.Write(strNews)
	_ = writer.Flush()

}
