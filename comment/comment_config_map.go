package comment

import (
	"deal_data/comment/parse"
)

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

// ParseMap TODO 评论独立出去使用python解析
// 然后使用grpc实现调用
var ParseMap = map[string]map[string]interface{}{
	"chinese.joins.com": nil,
}

func InitMap(articleUrl, newsUuid, commentPath, proxy, lang, dataTimeStr string) {
	ParseMap["chinese.joins.com"]["parse"] = parse.NewChineseJoinComment(articleUrl, newsUuid, commentPath, proxy, lang, dataTimeStr)
	ParseMap["chinese.joins.com"]["type"] = parse.ChineseComment
}
