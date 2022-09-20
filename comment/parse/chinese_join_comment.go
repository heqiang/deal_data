package parse

import (
	"deal_data/datadeal/util"
	"fmt"
	"github.com/lestrrat-go/libxml2/types"
	"github.com/lestrrat-go/libxml2/xpath"
	"go.uber.org/zap"
	"path/filepath"
	"strconv"
)

type ChineseJoinComment struct {
	articleUrl      string
	newsUuid        string
	commentPath     string
	proxy           string
	tool            *util.Tool
	commentJsonPath string
}

func NewChineseJoinComment(articleUrl, newsUuid, commentPath, proxy string) *ChineseJoinComment {
	comment := &ChineseJoinComment{
		articleUrl:  articleUrl,
		newsUuid:    newsUuid,
		commentPath: commentPath,
		proxy:       proxy,
		tool:        util.NewTool(proxy, commentPath),
	}
	comment.commentJsonPath = filepath.Join(commentPath, fmt.Sprintf("%s_comments.json", comment.tool.DateTimeStr))
	util.MkFile(comment.commentJsonPath)

	return comment
}

func (comment *ChineseJoinComment) Run() {
	fmt.Println(fmt.Sprintf("正在抓取%s的评论", comment.articleUrl))

	resp, err := util.Req(comment.articleUrl, comment.proxy, "comment")
	if err != nil {
		zap.L().Error(fmt.Sprintf("文章:%s\n评论请求失败", comment.articleUrl))
		return
	}
	doc := resp.(types.Document)
	defer doc.Free()
	nodes := xpath.NodeList(doc.Find(`//section[@class="reply-list"]/section/article`))
	if len(nodes) != 0 {
		for _, node := range nodes {
			var comments commentStruct
			comments.commentId = util.GenUuid()
			comments.newsId = comment.articleUrl
			comments.createTime = comment.getText(node, `//div[@class="comments-user"]/small/text()`)
			comments.content = comment.getText(node, `//section[@class="comments-content"]//text()`)
			likeCount, err := strconv.Atoi(comment.getText(node, `//button[@id="vote-up"]/span/text()`))
			if err != nil {
				zap.L().Warn(fmt.Sprintf("%s获取like_count 失败,err:%s", comment.articleUrl, err))
				return
			}
			comments.likeCount = likeCount
			comments.replayCount = 0
			comments.isHost = false

		}
	}
}

func (comment *ChineseJoinComment) getText(node types.Node, rule string) string {
	return xpath.String(node.Find(rule))
}
