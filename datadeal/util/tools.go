package util

import (
	"bufio"
	"deal_data/config"
	"deal_data/service/mysql"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

type Tool struct {
	header       string
	DateTimeStr  string
	time         string
	fileDir      string
	imgPath      string
	jsonPathName string
	proxy        string
}

func NewTool(proxy string, fileDir, dataTimeStr string) *Tool {
	currTime := time.Now()
	tool := &Tool{
		header:      config.Header,
		DateTimeStr: dataTimeStr,
		time:        currTime.Format("2006-01-02"),
		proxy:       proxy,
		fileDir:     fileDir,
	}

	return tool
}

func (tool *Tool) ProfileImgDownload(imgBasePath, fileName, url string, uuid string) {
	if !strings.HasPrefix(url, "http") {
		msg := fmt.Sprintf("%s:no schema", url)
		zap.L().Error(msg)
	}

	ProfilePathSuffix := fmt.Sprintf("%s_profile_image", tool.DateTimeStr)
	completeProfilePath := filepath.Join(imgBasePath, ProfilePathSuffix)

	Mkdir(completeProfilePath)

	fileDownload(completeProfilePath, fileName, url, tool.proxy, uuid)
}

func (tool *Tool) ImgOrFileDownload(filePath, fileName, url, fileType string, newsId int) {
	filePathSuffix := fmt.Sprintf("%s_%s", tool.DateTimeStr, fileType)
	completeFilePath := filepath.Join(filePath, filePathSuffix)

	Mkdir(completeFilePath)

	fileDownload(completeFilePath, fileName, url, tool.proxy, newsId)
}

func (tool *Tool) WriteNewsToJson(news map[string]interface{}, newsId int) {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()
	jsonName := fmt.Sprintf("%s_newsty.json", tool.DateTimeStr)
	tool.jsonPathName = filepath.Join(tool.fileDir, jsonName)
	MkFile(tool.jsonPathName)

	file, err := os.OpenFile(tool.jsonPathName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		zap.L().Error(fmt.Sprintf("新闻写入json文件失败,err:%s,新闻id：%d", err, newsId))
		return
	}
	defer file.Close()

	byteNews, err := json.Marshal(news)
	if err != nil {
		zap.L().Error(fmt.Sprintf("新闻序列化失败,err:%s,新闻id：%d", err, newsId))
		return
	}

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(string(byteNews) + "\n")
	if err != nil {
		zap.L().Error(fmt.Sprintf("新闻写入json文件失败,新闻id:%d,err:%v,", newsId, err))
		return
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("新闻写入json文件失败")
		zap.L().Error(fmt.Sprintf("新闻写入json文件失败,新闻id:%d,err:%v,", newsId, err))
		return
	}

	err = mysql.Db.UpdateNew(newsId, 2)
	if err != nil {
		zap.L().Error(fmt.Sprintf("新闻处理状态更新2失败,err:%s,debugStack:%s,新闻id:%d", err, debug.Stack(), newsId))
		return
	}
	fmt.Println(fmt.Sprintf("%s ->%d数据文本处理结束,更新状态为2", mysql.GetNowTime(), newsId))
}
