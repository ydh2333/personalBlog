package util

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/yaml.v3"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type ConfigJWT struct {
	JWT struct {
		Secret     string `yaml:"secret"`
		Expiration int64  `yaml:"expiration"`
	} `yaml:"jwt"`
}

func GenerateToken(name string) (string, error) {
	var cfg ConfigJWT
	// 读配置
	if data, err := os.ReadFile("conf/config.yaml"); err != nil {
		panic(fmt.Sprintf("读配置失败：%v", err))
	} else if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(fmt.Sprintf("解析YAML失败：%v", err))
	}

	var mySigningKey = []byte(cfg.JWT.Secret)
	cmyClaims := MyClaims{
		Username: name,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + cfg.JWT.Expiration,
			Issuer:    "yangduheng",
		},
	}

	// StandardClaims
	// MapClaims
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, cmyClaims)
	// 加密
	ss, err := t.SignedString(mySigningKey)
	return ss, err
}

func ParseToken(ss string) (*MyClaims, error) {
	var cfg ConfigJWT
	// 读配置
	if data, err := os.ReadFile("conf/config.yaml"); err != nil {
		panic(fmt.Sprintf("读配置失败：%v", err))
	} else if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(fmt.Sprintf("解析YAML失败：%v", err))
	}

	var mySigningKey = []byte(cfg.JWT.Secret)
	// 解析
	token, err := jwt.ParseWithClaims(ss, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*MyClaims), nil
}
