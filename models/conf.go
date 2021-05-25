package models

type MysqlConf struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	DB       string `json:"db"`
}

type AesConf struct {
	Key string `json:"key"`
}

type LogConf struct {
	Path string `json:"path"`
	File string `json:"file"`
}