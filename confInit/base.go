package confInit

import (
	"encoding/base64"
	"fmt"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
	log "github.com/sirupsen/logrus"
	"sync"
	"testAPI/models"
	"testAPI/utils"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const defaultRoot = "test"

var (
	m sync.RWMutex
	mc models.MysqlConf
	lc models.LogConf

	ac models.AesConf

	DB *gorm.DB
)

func Init()  {
	m.Lock()
	defer m.Unlock()

	err := config.Load(file.NewSource(
		file.WithPath("./config/application.yml"),
	))

	if err != nil {
		log.Error("加载配置文件失败: ", err)
		return
	}

	if err := config.Get(defaultRoot, "log").Scan(&lc); err != nil {
		log.Error("Log 配置读取失败: ", err)
		return
	}
	log.Info("Log 配置读取成功!")
	utils.LoggerToFile(lc.Path, lc.File)

	if err := config.Get(defaultRoot, "aes").Scan(&ac); err != nil {
		log.Error("Aes 配置读取失败: ", err)
		return
	}
	log.Info("Aes 配置读取成功!")

	key := []byte(ac.Key)

	if err := config.Get(defaultRoot, "mysql").Scan(&mc); err != nil {
		log.Error("Mysql 配置文件读取失败: ", err)
		return
	}
	log.Info("读取 Mysql 配置成功!")

	str, err := base64.StdEncoding.DecodeString(mc.Password)
	if err != nil {
		log.Error("Base64 decode failed: ", err)
		return
	}

	pwd, err := utils.AesDecrypt(str, key)
	if err != nil {
		log.Error("ASE decrypt failed: ", err)
		return
	}
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=true",
		mc.User, string(pwd), mc.Host, mc.DB)
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Error("Open mysql failed: ", err)
		return
	}

	DB.DB().SetMaxIdleConns(100)
	DB.DB().SetMaxOpenConns(200)
	DB.DB().SetConnMaxLifetime(10 * time.Second)

	if err := DB.DB().Ping(); err != nil {
		log.Error("连接数据库失败: ", err)
		return
	}
	log.Info("数据库连接成功!")
}