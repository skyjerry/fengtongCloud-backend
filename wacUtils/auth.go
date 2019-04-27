package wacUtils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

var (
	key = []byte("-jwt-ldwangzhen@gmail.com")
)

type UserPayload struct {
	jwt.StandardClaims
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

func GetToken(userId int, userName string) string {
	var userPayload UserPayload
	userPayload.StandardClaims.NotBefore = int64(time.Now().Unix())
	userPayload.StandardClaims.ExpiresAt = int64(time.Now().Unix() + beego.AppConfig.DefaultInt64("tokenExpireTime", 10))
	userPayload.UserId = userId
	userPayload.UserName = userName

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userPayload)
	ss, err := token.SignedString(key)
	if err != nil {
		return ""
	}
	return ss
}

func CheckToken(token string) int {
	if len(token) > 200 {
		return 1
	}
	t, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})

	//过期
	if t.Claims.Valid() != nil {
		logs.Warning("token过期")
		return 1
	}

	if err == nil {
		return 0
	}

	logs.Warning("未知错误")
	return -1
}

func CheckAuthorization(authorization string) int {
	if !strings.Contains(authorization, "Bearer ") {
		logs.Warning("缺少头部")
		return 1
	}
	strs := strings.Split(authorization, " ")
	return CheckToken(strs[1])
}
