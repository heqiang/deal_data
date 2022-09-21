package util

import (
	"bufio"
	"deal_data/global"
	"fmt"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestReq(t *testing.T) {
	articlUrl := `https://cdnn21.img.ria.ru/images/07e6/09/14/1818134137_243:0:3072:1591_1920x0_80_0_0_e32f488b92484e7f6d7721451e078f92.jpg`
	res, err := Req(articlUrl, global.Proxy)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(res))
	file, err := os.OpenFile("./img/test.jpg", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		msg := fmt.Sprintf("%s 文件打开失败", "")
		zap.L().Error(msg)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, _ = writer.Write(res)
	_ = writer.Flush()
}
