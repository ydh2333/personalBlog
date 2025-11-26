package util

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var mySigningKey = []byte("MyKey")

func GenerateToken(name string) string {

	c := MyClaims{
		Username: name,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 60*60*2,
			Issuer:    "yangduheng",
		},
	}

	// StandardClaims
	// MapClaims
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 加密
	ss, err := t.SignedString(mySigningKey)
	if err != nil {
		fmt.Println(err)
	}
	return ss
}

func ParseToken(ss string) (*MyClaims, error) {
	// 解析
	token, err := jwt.ParseWithClaims(ss, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	//fmt.Println(token)
	//fmt.Println("token.claims:", token.Claims.(*MyClaims).Username)
	return token.Claims.(*MyClaims), nil
}
