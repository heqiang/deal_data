package mysql

import "time"

type News struct {
	Id              int    `orm:"id" json:"id"`
	Uuid            string `orm:"uuid" json:"uuid"`
	SiteDomain      string `orm:"site_domain" json:"site_domain"`
	SourceName      string `orm:"source_name" json:"source_name"`
	Url             string `orm:"url" json:"url"`
	Title           string `orm:"title" json:"title"`
	Author          string `orm:"author" json:"author"`
	Content         string `orm:"content" json:"content"`
	CommentCount    int    `orm:"comment_count" json:"comment_count"`
	ReadCount       int    `orm:"read_count" json:"read_count"`
	LikeCount       int    `orm:"like_count" json:"like_count"`
	ForwardCount    int    `orm:"forward_count" json:"forward_count"`
	Type            string `orm:"type" json:"type"`
	Lang            string `orm:"lang" json:"lang"`
	Direction       string `orm:"direction" json:"direction"`
	BoardTheme      string `orm:"board_theme" json:"board_theme"`
	OriginalTags    string `orm:"original_tags" json:"original_tags"`
	SiteBoardName   string `orm:"site_board_name" json:"site_board_name"`
	RepostSource    string `orm:"repost_source" json:"repost_source"`
	IfRepost        int    `orm:"if_repost" json:"if_repost"`
	IfFrontPosition int    `orm:"if_front_position" json:"if_front_position"`
	PublishTime     string `orm:"publish_time" json:"publish_time"`
	InsertTime      string `orm:"insert_time" json:"insert_time"`
	DealCode        int    `orm:"deal_code" json:"deal_code"`
}

func (*News) TableName() string {
	return "news"
}

func GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
