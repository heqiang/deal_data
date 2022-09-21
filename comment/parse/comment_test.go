package parse

import "testing"

func TestChineseJoinComment_Run(t *testing.T) {
	url := "https://chinese.joins.com/news/articleView.html?idxno=107289"
	uuid := "test"
	comment := NewChineseJoinComment(url, uuid, ".", "http://127.0.0.1:9910", "en", "20222")
	comment.Run()
}
