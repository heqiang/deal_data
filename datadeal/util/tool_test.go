package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

//func TestNewTool(t *testing.T) {
//	tool := NewTool("", "./")
//	tool.ProfileImgDownload("./", "test.jpg",
//		"https://www.keaidian.com/uploads/allimg/190424/24110307_8.jpg", "sss")
//	_, err := os.OpenFile("test.jpg", os.O_WRONLY|os.O_CREATE, 0666)
//	if err != nil {
//		panic(err)
//	}
//}
//
//func TestParse(t *testing.T) {
//	articleUrl := "https://www.baidu.com"
//	host := Parse(articleUrl)
//	fmt.Println(host)
//}

func TestWrite(t *testing.T) {
	file, err := os.OpenFile("test.txt", os.O_WRONLY|os.O_CREATE, 777)
	if err != nil {
		fmt.Println(err)
	}
	write := bufio.NewWriter(file)
	_, err = write.WriteString("sss" + "\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	write.WriteString("ssddds" + "\n")
	write.WriteString("ggg" + "\n")
	write.Flush()
}

func TestW(t *testing.T) {
	file, err := os.OpenFile(filepath.Join("../../comment/parse/commentJson", "jj.json"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	for i := 0; i < 5; i++ {
		write.WriteString("http://c.biancheng.net/golang/ \n")
	}
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}
