package util

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
)

func Req(reqRul string, header string, proxyIp string) ([]byte, error) {
	client := http.Client{}
	if proxyIp != "" {
		uri, err := url.Parse(proxyIp)
		if err != nil {
			zap.L().Warn(fmt.Sprintf("parse proxyIp error: %s", err))
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(uri)}
	}
	req, err := http.NewRequest("GET", reqRul, nil)
	req.Header.Add("User-Agent", header)
	if err != nil {
		zap.L().Warn(fmt.Sprintf("NewRequest Failed"))
	}
	response, _ := client.Do(req)
	defer response.Body.Close()

	return io.ReadAll(response.Body)

}
