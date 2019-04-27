package controllers

import "github.com/astaxie/beego"

type BaseController struct {
	beego.Controller
}

func (b *BaseController) ApiResponse(status int, msg string, data map[string]interface{}) {
	returnData := map[string]interface{}{
		"code": status,
		"msg":  msg,
		"data": data,
	}
	b.Data["json"] = returnData
	b.ServeJSON()
}
