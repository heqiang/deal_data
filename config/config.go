package config

type ServerConfig struct {
	Name         string `mapstructure:"name"`
	*MysqlConfig `mapstructure:"mysql"`
	*LogConfig   `mapstructure:"log"`
	*MqConfig    `mapstructure:"mq"`
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

type MqConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	*QueueConfig    `mapstructure:"queue"`
	*ExchangeConfig `mapstructure:"exchange"`
}

type QueueConfig struct {
	QueueName  string `mapstructure:"queueName,omitempty"`
	DurAble    bool   `mapstructure:"durable,omitempty"`
	AutoDelete bool   `mapstructure:"autoDelete,omitempty"`
	Exclusive  bool   `mapstructure:"exclusive,omitempty"`
	NoWait     bool   `mapstructure:"noWait,omitempty"`
}

type ExchangeConfig struct {
	ExchangeName string `mapstructure:"exchangeName,omitempty"`
	Kind         string `mapstructure:"kind,omitempty"`
	DurAble      bool   `mapstructure:"durable,omitempty"`
	AutoDelete   bool   `mapstructure:"autoDelete,omitempty"`
	Internal     bool   `mapstructure:"internal,omitempty"`
	NoWait       bool   `mapstructure:"noWait,omitempty"`
}
