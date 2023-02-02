package mysqlservice

import (
	"deal_data/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var MaxId int

type MysqlConn struct {
	Db *gorm.DB
}

func NewMysqlConn(config *config.MysqlConfig) *MysqlConn {
	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn = fmt.Sprintf(dsn, config.User, config.Password, config.Host, config.Port, config.Db)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		panic("mysql Conn failed")
	}
	sqlDb, err := db.DB()
	sqlDb.SetMaxIdleConns(100) //空闲连接数
	sqlDb.SetMaxOpenConns(200) //最大连接数
	sqlDb.SetConnMaxLifetime(time.Minute)
	return &MysqlConn{Db: db}
}

func (conn *MysqlConn) SelectZero() (newsList []News, err error) {
	selectResult := conn.Db.Limit(100).Where("deal_code =?", 0).Find(&newsList)
	if selectResult.Error != nil {
		return []News{}, selectResult.Error
	}
	newsId := make([]int, len(newsList))
	for i, v := range newsList {
		newsId[i] = v.Id
	}
	conn.Db.Table("news_test").Where("id in ？", newsId).Updates(map[string]interface{}{"name": "hello", "age": 18})
	return
}

// UpdateNew 数据状态更新
// 0 未处理
// 1 处理中
// 2 处理结束
func (conn *MysqlConn) UpdateNew(id int, statusNum int) (err error) {
	updateRes := conn.Db.Model(&News{}).Where("id = ?", id).Update("deal_code", statusNum).Update("deal_time", GetNowTime())
	if updateRes.Error != nil {
		return updateRes.Error
	}

	return nil
}
