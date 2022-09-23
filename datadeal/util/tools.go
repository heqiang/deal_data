package util

import (
	"bufio"
	"deal_data/comment/parse/comm_struct"
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
		header:      global.Header,
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
	global.Cond.L.Lock()
	defer global.Cond.L.Unlock()
	jsonName := fmt.Sprintf("%s_newsty.json", tool.DateTimeStr)
	tool.jsonPathName = filepath.Join(tool.fileDir, jsonName)
	MkFile(tool.jsonPathName)

	file, err := os.OpenFile(tool.jsonPathName, os.O_APPEND|os.O_CREATE, 7777)
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
	_, _ = writer.WriteString(string(byteNews) + "\n")
	_ = writer.Flush()
}

func (tool *Tool) WriteCommentToJson(comment comm_struct.Comment, jsonPath, url string) {
	defer global.Mutex.Unlock()
	global.Mutex.Lock()
	file, err := os.OpenFile(jsonPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println(err)
		zap.L().Error(fmt.Sprintf("json文件打开失败,err:%s,新闻url：%s", err, url))
		return
	}
	defer file.Close()

	byteNews, err := json.Marshal(comment)
	if err != nil {
		zap.L().Error(fmt.Sprintf("评论序列化失败,err:%s,新闻url：%s", err, url))
		return
	}
	_, err = file.Write(byteNews)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, _ = writer.WriteString(string(byteNews) + "\n")
	_ = writer.Flush()
	fmt.Println("评论写入json完毕")
}
