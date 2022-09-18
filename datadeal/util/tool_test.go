package util

import (
	"os"
	"testing"
)

func TestNewTool(t *testing.T) {
	tool := NewTool("")
	tool.ProfileImgDownload("./", "test.jpg",
		"https://www.keaidian.com/uploads/allimg/190424/24110307_8.jpg")
	_, err := os.OpenFile("test.jpg", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
}
