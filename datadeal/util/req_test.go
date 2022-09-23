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
	articlUrl := `https://static.ctwant.com/images/cover/54/208154/md-b42c12db8e9ceb89d9ddc2bef831c905.jpg`
	res, err := Req(articlUrl, global.Proxy)
	if err != nil {
		fmt.Println(err)
	}
	file, err := os.OpenFile("./test.jpg", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		msg := fmt.Sprintf("%s 文件打开失败", "")
		zap.L().Error(msg)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, _ = writer.Write(res)
	_ = writer.Flush()
	fmt.Println("ok")
}
