package global

import (
	"deal_data/config"
	"deal_data/mysqlservice"
)

var (
	Db           *mysqlservice.MysqlConn
	ServerConfig config.ServerConfig
	Header       map[string]string
)

func init() {
	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"

	Header = map[string]string{"user-agent": userAgent}
}
