package parse

import (
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
	articleUrl      string
	newsUuid        string
	commentPath     string
	proxy           string
	tool            *util.Tool
	commentJsonPath string
	lang            string
}

type Comment struct {
	CommentId       string `json:"comment_id"`
	NewsId          string `json:"news_id"`
	SourceUrl       string `json:"source_url"`
	CreateTime      string `json:"create_time,omitempty"`
	Content         string `json:"content,omitempty"`
	LikeCount       int    `json:"like_count,omitempty"`
	ReplayCount     int    `json:"replay_count,omitempty"`
	IsHost          bool   `json:"is_host,omitempty"`
	ReplayCommentId string `json:"replay_comment_id,omitempty"`
	WhichPlatform   string `json:"which_platform,omitempty"`
	UserId          string `json:"user_id,omitempty"`
	UserName        string `json:"user_name,omitempty"`
	CrawlTime       string `json:"crawl_time,omitempty"`
	InsertTime      string `json:"insert_time,omitempty"`
	Lang            string `json:"lang_recg,omitempty"`
	UserPicLink     string `json:"user_pic_link,omitempty"`
}

func NewChineseJoinComment(articleUrl, newsUuid, commentPath, proxy, lang, dataTimeStr string) *ChineseJoinComment {
	chineseComment := &ChineseJoinComment{
		articleUrl:  articleUrl,
		newsUuid:    newsUuid,
		commentPath: commentPath,
		proxy:       proxy,
		tool:        util.NewTool(proxy, commentPath, dataTimeStr),
		lang:        lang,
	}
	util.Mkdir(chineseComment.commentPath)
	chineseComment.commentJsonPath = filepath.Join(commentPath, fmt.Sprintf("%s_comments.json", chineseComment.tool.DateTimeStr))

	return chineseComment
}

func (comm *ChineseJoinComment) Run() {
	fmt.Println(fmt.Sprintf("正在抓取:%s的评论", comm.articleUrl))
	dom, err := util.GetDocument(comm.articleUrl, comm.proxy)
	if err != nil {
		zap.L().Error(fmt.Sprintf("文章:%s\n评论请求失败,err:%s", comm.articleUrl, err))
		return
	}

	selection := dom.Find(".reply-list section article")
	if len(selection.Nodes) > 0 {
		selection.Each(func(i int, selection *goquery.Selection) {
			var comments Comment
			comments.CommentId = util.GenUuid()
			comments.NewsId = comm.newsUuid
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
			comments.Lang = comm.lang
			comments.UserPicLink = ""
			fmt.Println(comments.UserName)
			comm.tool.WriteToJson(comments, comm.commentJsonPath, comm.articleUrl)
		})
	}

}
