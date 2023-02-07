package util

import (
	"deal_data/config"
	"testing"
)

func TestMkdir(t *testing.T) {
	MkFile("E:\\goproject\\deal_data\\data\\test.txt")
}

func TestReq(t *testing.T) {
	filPath := "./"
	fileName := "test.jpg"
	url := "https://www.stjornarradid.is/library/Heimsljos/52192293662_3fdcb3422c_k.jpg?proc=singleNewsItem"
	proxy := config.Proxy
	id := 1
	fileDownload(filPath, fileName, url, proxy, id)
}
