package config

import (
	"sync"
)

var (
	Conf        *ServerConfig
	Header      = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"
	AbsDataPath string
	Proxy       = "http://127.0.0.1:9910"
	Mutex       sync.Mutex
	Cond        sync.Cond
)
