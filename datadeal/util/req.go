package util

import (
	"deal_data/global"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
)

func Req(reqRul string, proxyIp string) (io.Reader, error) {
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
		return nil, nil
	}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Warn(fmt.Sprintf("url req fail,%s", reqRul))
		return nil, err
	}
	defer resp.Body.Close()

	return resp.Body, nil

}
