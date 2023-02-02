package util

import (
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Req(reqRul string, proxyIp string) ([]byte, error) {
	defer func() {
		//匿名函数来捕获异常
		if err := recover(); err != nil {
			fmt.Println("异常", err)
		}
	}()

	fmt.Println(fmt.Sprintf("正在下载图片:%s", reqRul))
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxyIp)
	}
	client := &http.Client{Transport: &http.Transport{Proxy: proxy}}

	req, err := http.NewRequest("GET", reqRul, nil)
	if err != nil {
		zap.L().Warn("request请求构建失败")
	}

	res, err := client.Do(req)
	if err != nil {
		zap.L().Warn(fmt.Sprintf("请求失败:%s", reqRul))
	}

	defer res.Body.Close()
	respByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		zap.L().Warn(fmt.Sprintf("图片响应读取失败:%s", reqRul))
	}
	fmt.Println(fmt.Sprintf("图片下载完毕:%s", reqRul))
	return respByte, nil
}
