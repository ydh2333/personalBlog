package util

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ConfigDB struct {
	DB struct {
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		User    string `yaml:"username"`
		Pwd     string `yaml:"password"`
		Dbname  string `yaml:"db_name"`
		Charset string `yaml:"charset"`
	} `yaml:"db"`
}

func SqlConnect() (*gorm.DB, error) {
	// 读配置
	var cfg ConfigDB
	if data, err := os.ReadFile("conf/config.yaml"); err != nil {
		panic(fmt.Sprintf("读配置失败：%v", err))
	} else if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(fmt.Sprintf("解析YAML失败：%v", err))
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true", cfg.DB.User, cfg.DB.Pwd, cfg.DB.Host, cfg.DB.Port, cfg.DB.Dbname, cfg.DB.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, ErrDBConnect
	}
	return db, nil
}
