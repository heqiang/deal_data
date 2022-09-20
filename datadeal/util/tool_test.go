package util

import (
	"fmt"
	"os"
	"testing"
)

func TestNewTool(t *testing.T) {
	tool := NewTool("", "./")
	tool.ProfileImgDownload("./", "test.jpg",
		"https://www.keaidian.com/uploads/allimg/190424/24110307_8.jpg", "sss")
	_, err := os.OpenFile("test.jpg", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
}

func TestParse(t *testing.T) {
	articleUrl := "https://www.baidu.com"
	host := Parse(articleUrl)
	fmt.Println(host)
}
