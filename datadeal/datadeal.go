package datadeal

import (
	"deal_data/datadeal/util"
	"deal_data/global"
	"deal_data/mysqlservice"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
	"time"
)

type DataDeal struct {
	tool     *util.Tool
	country  string
	currTime string
	filePath string
}

func NewDataDeal(proxy, country string) *DataDeal {
	dataDeal := &DataDeal{
		currTime: time.Now().Format("2006-01-02"),
	}
	dataDeal.country = getDirection(country)
	dataDeal.filePath = filepath.Join(global.AbsDataPath, dataDeal.country, dataDeal.currTime)
	dataDeal.tool = util.NewTool(proxy, dataDeal.filePath)

	util.Mkdir(dataDeal.filePath)
	return dataDeal
}

func getDirection(country string) (newCountry string) {
	countryList := []string{"china_hk", "china_tw", "china_macao", "china_xz"}
	for _, cou := range countryList {
		if strings.Contains(country, cou) {
			newCountry = strings.Replace(country, "china_hk", "hk", -1)
			return strings.ToUpper(newCountry)
		}
	}
	return
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
