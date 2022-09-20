package parse

type commentStruct struct {
	commentId       string `json:"comment_id"`
	newsId          string `json:"news_id"`
	sourceUrl       string `json:"source_url"`
	createTime      string `json:"create_time,omitempty"`
	content         string `json:"content,omitempty"`
	likeCount       int    `json:"like_count,omitempty"`
	replayCount     int    `json:"replay_count,omitempty"`
	isHost          bool   `json:"is_host,omitempty"`
	ReplayCommentId string `json:"replay_comment_id,omitempty"`
	whichPlatform   string `json:"which_platform,omitempty"`
	userId          string `json:"user_id,omitempty"`
	userName        string `json:"user_name,omitempty"`
	crawlTim        string `json:"crawl_tim,omitempty"`
	insertTime      string `json:"insert_time,omitempty"`
	lang            string `json:"lang_recg,omitempty"`
	userPicLink     string `json:"user_pic_link,omitempty"`
}
