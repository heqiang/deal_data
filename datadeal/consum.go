package datadeal

import (
	"context"
	"deal_data/comment"
	"deal_data/datadeal/util"
	"deal_data/global"
	"deal_data/mysqlservice"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"path/filepath"
	"runtime/debug"
	"time"
)

type oldConStruct struct {
	Type        string      `json:"type"`
	Name        interface{} `json:"name,omitempty"`
	Md5Src      string      `json:"md5src,omitempty"`
	Description string      `json:"description,omitempty"`
	Src         string      `json:"src,omitempty"`
	Data        string      `json:"data,omitempty"`
}

type DataDeal struct {
	tool        *util.Tool
	country     string
	currTime    string
	filePath    string
	DateTimeStr string
}

func (p *Pipeline) Consume(in <-chan mysqlservice.News, ctx context.Context) {
	for data := range in {
		err := global.Db.UpdateNew(data.Id, 1)
		fmt.Println(fmt.Sprintf("%s:正在处理数据:%d", time.Now().Format("2006-01-02 15:04:05"), data.Id))
		if err != nil {
			fmt.Println(err)
		}
		//TODO 协程超时退出,数量控制
		p.w.Run(func() {
			news := data
			defer func() {
				if err := recover(); err != nil {
					zap.L().Error(fmt.Sprintf("数据处理异常:%s,err:%s,新闻id:%d,", err, debug.Stack(), news.Id))
					return
				} else {
					err = global.Db.UpdateNew(news.Id, 2)
					if err != nil {
						zap.L().Error(fmt.Sprintf("新闻处理状态更新2失败,err:%s,debugStack:%s,新闻id:%d", err, debug.Stack(), news.Id))
						return
					}
					fmt.Println(fmt.Sprintf("%s:%d数据处理结束", time.Now().Format("2006-01-02 15:04:05"), data.Id))
				}
			}()
			//数据处理的主要逻辑
			deal := NewDataDeal(global.Proxy, news.Direction)
			deal.TransNewsToJson(news)
			deal.download(news.Content, news.Id)
		})
	}
}

func NewDataDeal(proxy, country string) *DataDeal {
	dataDeal := &DataDeal{
		currTime:    time.Now().Format("2006-01-02"),
		DateTimeStr: time.Now().Format("20060102150405"),
	}
	dataDeal.country = util.GetDirection(country)
	dataDeal.filePath = filepath.Join(global.AbsDataPath, dataDeal.country, dataDeal.currTime)
	dataDeal.tool = util.NewTool(proxy, dataDeal.filePath, dataDeal.DateTimeStr)

	util.Mkdir(dataDeal.filePath)
	return dataDeal
}

func (d *DataDeal) parseComment(articleUrl, uuid string) {
	siteDomain := util.ParseHost(articleUrl)
	if _, ok := comment.CommentMap[siteDomain]; ok {
		msgInfo := map[string]string{
			"articleUrl":      articleUrl,
			"commentPath":     d.filePath,
			"commentJsonName": fmt.Sprintf("%s_newsty.json", d.DateTimeStr),
			"newUUID":         uuid,
		}
		msgByte, err := json.Marshal(msgInfo)
		if err != nil {
			zap.L().Error(fmt.Sprintf("%s序列化错误,err:%s", articleUrl, err))
			return
		}
		global.RabbitMq.SentMessage(string(msgByte))
	}
}

func (d *DataDeal) download(content string, newsId int) {
	err := global.Db.UpdateNew(newsId, 1)
	if err != nil {
		zap.L().Error(fmt.Sprintf("新闻处理状态更新1失败,err:%s,新闻id:%d", err, newsId))
		return
	}
	if content != "" {
		var cons []oldConStruct
		err := json.Unmarshal([]byte(content), &cons)
		if err != nil {
			zap.L().Error(fmt.Sprintf("content 反序列化失败,err:%s,新闻id:%d", err, newsId))
			return
		}

		for _, con := range cons {
			if con.Type == "image" || con.Type == "file" {
				d.tool.ImgOrFileDownload(d.filePath, con.Md5Src, con.Src, con.Type, newsId)
			}
		}
	}
}

// TransNewsToJson 新闻数据写入json
func (d *DataDeal) TransNewsToJson(news mysqlservice.News) {
	newsMap := map[string]interface{}{
		"news_id":           news.Uuid,
		"source_name":       news.SourceName,
		"site_domain":       news.SiteDomain,
		"author":            d.transStrToList(news.Author),
		"url":               news.Url,
		"time":              news.PublishTime,
		"site_board_name":   news.SiteBoardName,
		"board_theme":       news.BoardTheme,
		"if_front_position": news.IfFrontPosition,
		"type":              "",
		"crawl_time":        news.InsertTime,
		"lang":              news.Lang,
		"direction":         news.Direction,
		"comment_count":     news.CommentCount,
		"forward_count":     news.ForwardCount,
		"like_count":        news.LikeCount,
		"read_count":        news.ReadCount,
		"original_tags":     d.transStrToList(news.OriginalTags),
		"if_repost":         news.IfRepost,
		"repost_source":     news.RepostSource,
		"insert_time":       news.InsertTime,
		"title":             news.Title,
		"content":           d.transContent(news.Content),
	}
	d.tool.WriteNewsToJson(newsMap, news.Id)

}

func (d *DataDeal) transStrToList(newsItem string) []interface{} {
	var resList []interface{}

	err := json.Unmarshal([]byte(newsItem), &resList)
	if err != nil {
		zap.L().Warn("json 序列化失败")
		return []interface{}{}
	}
	return resList
}

func (d *DataDeal) transContent(contents string) []interface{} {
	var newContent []interface{}
	contentsList := d.transStrToList(contents)
	if len(contentsList) >= 1 {
		for _, con := range contentsList {
			conStruct := con.(map[string]interface{})
			// content中的结构转换
			if conStruct["type"] == "image" || conStruct["type"] == "file" {
				newCon := d.transImageOrFileCon(conStruct, conStruct["type"].(string))
				newContent = append(newContent, newCon)
			}
			if conStruct["type"] == "text" {
				newContent = append(newContent, conStruct)
			}
		}
	}
	return newContent
}

func (d *DataDeal) transImageOrFileCon(con map[string]interface{}, fileType string) (newCons map[string]interface{}) {
	newCons = make(map[string]interface{})
	newCons["type"] = fileType
	newCons["data"] = map[string]interface{}{
		"description": con["description"],
		"md5src":      con["md5src"],
		"src":         con["src"],
		"name":        con["name"],
	}
	return
}
