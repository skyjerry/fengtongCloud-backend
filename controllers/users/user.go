package users

import (
	"github.com/astaxie/beego"
	"time"
	//"github.com/astaxie/beego"
	"wac/controllers"
	"wac/models"
	"wac/wacUtils"
)

type UserController struct {
	controllers.BaseController
}

func (c *UserController) Login() {
	username := c.GetString("username")
	password := c.GetString("password")
	user, errStr := models.Login(username, password)
	if errStr != "" {
		c.ApiResponse(403, errStr, map[string]interface{}{})
	}

	token := wacUtils.GetToken(user.Id, user.Username)
	if token == "" {
		c.ApiResponse(500, "服务器错误", map[string]interface{}{})
	}

	c.ApiResponse(200, "登录成功", map[string]interface{}{"access_token": token, "expired_at": time.Now().Unix() + beego.AppConfig.DefaultInt64("tokenExpireTime", 10)})
}
