package comment

import (
	"deal_data/comment/parse"
)

// ParseMap TODO 评论独立出去使用python解析
// 然后使用grpc实现调用
var ParseMap = map[string]map[string]interface{}{
	"chinese.joins.com": nil,
}

func InitMap(articleUrl, newsUuid, commentPath, proxy, lang, dataTimeStr string) {
	ParseMap["chinese.joins.com"]["parse"] = parse.NewChineseJoinComment(articleUrl, newsUuid, commentPath, proxy, lang, dataTimeStr)
	ParseMap["chinese.joins.com"]["type"] = parse.ChineseComment
}
