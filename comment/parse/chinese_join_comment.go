package parse

import (
	"deal_data/comment/parse/commentstruct"
	"deal_data/datadeal/util"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"path/filepath"
	"strconv"
	"time"
)

var ChineseComment *ChineseJoinComment

type ChineseJoinComment struct {
	commentstruct.BaseParseStruct
	articleUrl      string
	newsUuid        string
	commentPath     string
	proxy           string
	tool            *util.Tool
	commentJsonPath string
	lang            string
}

func NewChineseJoinComment(articleUrl, newsUuid, commentPath, proxy, lang, dataTimeStr string) *ChineseJoinComment {
	comment := &ChineseJoinComment{
		articleUrl:  articleUrl,
		newsUuid:    newsUuid,
		commentPath: commentPath,
		proxy:       proxy,
		tool:        util.NewTool(proxy, commentPath, dataTimeStr),
		lang:        lang,
	}
	util.Mkdir(comment.commentPath)
	comment.commentJsonPath = filepath.Join(commentPath, fmt.Sprintf("%s_comments.json", comment.tool.DateTimeStr))

	return comment
}

func (comment *ChineseJoinComment) Run() {
	fmt.Println(fmt.Sprintf("正在抓取:%s的评论", comment.articleUrl))
	dom, err := util.GetDocument(comment.articleUrl, comment.proxy)
	if err != nil {
		zap.L().Error(fmt.Sprintf("文章:%s\n评论请求失败,err:%s", comment.articleUrl, err))
		return
	}

	selection := dom.Find(".reply-list section article")
	if len(selection.Nodes) > 0 {
		selection.Each(func(i int, selection *goquery.Selection) {
			var comments commentstruct.CommentStruct
			comments.CommentId = util.GenUuid()
			comments.NewsId = comment.newsUuid
			comments.CreateTime = selection.Find(".list-comments .comments-user small").First().Text()
			comments.Content = selection.Find(".list-comments .comments-content").Text()
			comments.LikeCount, _ = strconv.Atoi(selection.Find(".list-comments .text-danger").Text())
			comments.ReplayCount = 0
			comments.IsHost = false
			comments.ReplayCommentId = ""
			comments.WhichPlatform = "中央日报"
			comments.UserId = util.GenUuid()
			comments.UserName = selection.Find(".list-comments .comments-user strong").Text()
			comments.CrawlTime = time.Now().Format("2006-01-02 15:04:05")
			comments.InsertTime = ""
			comments.Lang = comment.lang
			comments.UserPicLink = ""
			fmt.Println(comments.UserName)
			comment.tool.WriteToJson(comments, comment.commentJsonPath, comment.articleUrl)
		})
	}

}
