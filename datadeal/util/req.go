package util

import (
	"deal_data/global"
	"fmt"
	"github.com/lestrrat-go/libxml2"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
)

// Req 这个函数主要是文件或者评论链接的请求
func Req(reqRul string, proxyIp string, reqType string) (interface{}, error) {
	client := http.Client{}
	if proxyIp != "" {
		uri, err := url.Parse(proxyIp)
		if err != nil {
			zap.L().Warn(fmt.Sprintf("parse proxyIp error: %s", err))
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(uri)}
	}
	req, err := http.NewRequest("GET", reqRul, nil)
	req.Header.Add("User-Agent", global.Header)
	if err != nil {
		zap.L().Warn(fmt.Sprintf("NewRequest Failed"))
	}
	response, _ := client.Do(req)
	defer response.Body.Close()
	if reqType == "file" {
		return io.ReadAll(response.Body)
	}
	return libxml2.ParseHTMLReader(response.Body)

}
