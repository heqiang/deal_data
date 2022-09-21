package util

import (
	"deal_data/global"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/avast/retry-go"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Req(reqRul string, proxyIp string) ([]byte, error) {
	var respByte []byte
	_ = retry.Do(
		func() error {
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
				return err
			}
			resp, err := client.Do(req)
			if err != nil {
				zap.L().Warn(fmt.Sprintf("url req fail,%s", reqRul))
			}
			defer resp.Body.Close()
			respByte, _ = ioutil.ReadAll(resp.Body)
			return nil
		}, retry.Attempts(3),
	)

	return respByte, nil
}

func GetDocument(articleUrl, proxy string) (*goquery.Document, error) {
	resBytes, err := Req(articleUrl, proxy)
	if err != nil {
		return nil, err
	}

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(resBytes)))
	if err != nil {
		return nil, err
	}

	return dom, nil
}
