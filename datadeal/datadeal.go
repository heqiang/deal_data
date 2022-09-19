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

type contentsStruct struct {
	Type        string      `json:"type"`
	Name        interface{} `json:"name,omitempty"`
	Md5Src      string      `json:"md5src,omitempty"`
	Description string      `json:"description,omitempty"`
	Src         string      `json:"src,omitempty"`
	Data        string      `json:"data,omitempty"`
}

func NewDataDeal(proxy, country string) *DataDeal {
	dataDeal := &DataDeal{
		tool:     util.NewTool(proxy),
		currTime: time.Now().Format("2006-01-02"),
	}
	dataDeal.country = getDirection(country)
	dataDeal.filePath = filepath.Join(global.AbsDataPath, dataDeal.country, dataDeal.currTime)
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

func (d *DataDeal) fileDownload(content string, newsId int) {
	if content != "" {
		var cons []contentsStruct
		err := json.Unmarshal([]byte(content), &cons)
		if err != nil {
			zap.L().Error(fmt.Sprintf("content 反序列化失败,err:%s,新闻id:%d", err, newsId))
			return
		}

		for _, con := range cons {
			if con.Type == "image" {
				d.tool.ImgOrFileDownload(d.filePath, con.Md5Src, con.Src, "image", newsId)
			} else if con.Type == "file" {
				d.tool.ImgOrFileDownload(d.filePath, con.Md5Src, con.Src, "file", newsId)

			}

		}
	}
}

func (d *DataDeal) TransNewsToJson(news mysqlservice.News) map[string]interface{} {
	newsMap := map[string]interface{}{
		"news_id":           news.Id,
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
	return newsMap

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

func (d *DataDeal) transContent(contents string) []contentsStruct {
	//var newContent []contentsStruct
	contentsList := d.transStrToList(contents)
	if len(contentsList) != 1 {
		for _, con := range contentsList {
			conStruct := con.(contentsStruct)
			if conStruct.Type == "image" {
				continue
			}

		}
	}

	return []contentsStruct{}
}
func (d *DataDeal) transImageOrFileCon(con contentsStruct, fileType string) map[string]interface{} {
	newCon := make(map[string]interface{})
	newCon["type"] = fileType
	newCon["data"] = con
	return newCon
}

func (d *DataDeal) transImageOrFileDict() {

}
