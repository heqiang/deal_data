package util

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func Req(reqRul, proxyIp string, dstFile *os.File) (io.Reader, error) {
	fmt.Println(fmt.Sprintf("正在下载图片:%s", reqRul))

	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxyIp)
	}
	client := &http.Client{Transport: &http.Transport{Proxy: proxy}}

	req, err := http.NewRequest("GET", reqRul, nil)
	if err != nil {
		zap.L().Warn("request请求构建失败")
	}

	resp, err := client.Do(req)
	if err != nil {
		zap.L().Warn(fmt.Sprintf("url req fail,%s", reqRul))
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println(fmt.Sprintf("图片：%s下载完成", reqRul))

	imgSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	newReader := bufio.NewReaderSize(resp.Body, imgSize)
	_, err = io.Copy(dstFile, resp.Body)
	if err != nil {
		fmt.Println(err)
		zap.L().Warn(fmt.Sprintf("文件:%s写入失败,err:%v", reqRul, err))
		return nil, nil
	}
	return newReader, nil

}
