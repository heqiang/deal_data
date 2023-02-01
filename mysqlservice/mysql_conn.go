package mysqlservice

import (
	"deal_data/config"
	"deal_data/datadeal/util"
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

// Select TODO 记录最大值实现每次查询不会查询到重复数据
func (conn *MysqlConn) Select() (newsList []News, err error) {
	selectResult := conn.Db.Limit(50).Where("id>?", MaxId).Where("deal_code =?", 0).Find(&newsList)
	if selectResult.Error != nil {
		return []News{}, selectResult.Error
	}
	return
}

// UpdateNew 数据状态更新
// 0 未处理
// 1 处理中
// 2 处理结束
func (conn *MysqlConn) UpdateNew(id int, statusNum int) (err error) {
	updateRes := conn.Db.Model(&News{}).Where("id = ?", id).Update("deal_code", statusNum).Update("deal_time", util.GetNowTime())
	if updateRes.Error != nil {
		return updateRes.Error
	}

	return nil
}
