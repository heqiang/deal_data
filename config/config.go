package config

type ServerConfig struct {
	Name         string `mapstructure:"name" json:"name,omitempty"`
	*MysqlConfig `mapstructure:"mysql" json:"*_mysql_config,omitempty"`
	*LogConfig   `mapstructure:"log" json:"*_log_config,omitempty"`
}
type MysqlConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     int    `json:"port" mapstructure:"port"`
	Name     string `json:"name" mapstructure:"name"`
	User     string `json:"user" mapstructure:"user"`
	Password string `json:"password" mapstructure:"password"`
	Db       string `json:"db" mapstructure:"db"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
