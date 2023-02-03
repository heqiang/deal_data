package util

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
)

func Req(reqRul string, proxyIp string) (io.Reader, error) {
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

	return resp.Body, nil

}
