package datadeal

import (
	"deal_data/datadeal/util"
	"deal_data/global"
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

type contents []struct {
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
		cons := contents{}
		err := json.Unmarshal([]byte(content), &cons)
		if err != nil {
			zap.L().Error(fmt.Sprintf("content 反序列化失败,err:%s,新闻id:%d", err, newsId))
			return
		}
		for _, con := range cons {
			if con.Type == "image" {

			}
		}
	}
}
