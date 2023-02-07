package datadeal

import (
	"context"
	"deal_data/config"
	"deal_data/datadeal/util"
	"deal_data/datadeal/worker"
	"deal_data/service/mysql"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"go.uber.org/zap"
	"path/filepath"
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

type Pipeline struct {
	w *worker.Worker
}

func NewPipeline() *Pipeline {
	return &Pipeline{w: worker.New(100)}
}

func (p *Pipeline) Consume(in <-chan mysql.News) {
	for data := range in {
		p.w.Run(func(news mysql.News) {
			deal := NewDataDeal(config.Proxy, news.Direction)
			deal.TransNewsToJson(news)
			ch := make(chan struct{}, 0)

			retryNum := 0
			_ = retry.Do(func() error {
				ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
				defer cancel()

				go func(deal *DataDeal, news mysql.News, ctx context.Context) {
					deal.download(news.Content, news.Id)
					ch <- struct{}{}
				}(deal, news, ctx)

				select {
				case <-ch:
				case <-ctx.Done():
					retryNum++
				}
				return nil
			}, retry.Attempts(5))

			if retryNum == 5 {
				zap.L().Warn(fmt.Sprintf("图片下载五次失败,新闻url:%s,新闻%d,新闻:%v", news.Url, news.Id, news))
			}
		}, data)
	}
}

func NewDataDeal(proxy, country string) *DataDeal {
	dataDeal := &DataDeal{
		currTime:    time.Now().Format("2006-01-02"),
		DateTimeStr: time.Now().Format("20060102"),
	}
	dataDeal.country = util.GetDirection(country)
	dataDeal.filePath = filepath.Join(config.AbsDataPath, dataDeal.country, dataDeal.currTime)
	dataDeal.tool = util.NewTool(proxy, dataDeal.filePath, dataDeal.DateTimeStr)

	util.Mkdir(dataDeal.filePath)

	return dataDeal
}

func (d *DataDeal) download(content string, newsId int) {
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
func (d *DataDeal) TransNewsToJson(news mysql.News) {
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
			conStruct, ok := con.(map[string]interface{})
			if ok {
				if conStruct["type"] == "image" || conStruct["type"] == "file" {
					newCon := d.transImageOrFileCon(conStruct, conStruct["type"].(string))
					newContent = append(newContent, newCon)
				}
				// content中的结构转换
				if conStruct["type"] == "text" {
					newContent = append(newContent, conStruct)
				}
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
