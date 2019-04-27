// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strings"
	"wac/controllers"
	"wac/controllers/users"
	myutils "wac/wacUtils"
)

func init() {
	//登录
	beego.Router("/login", &users.UserController{}, "post:Login")
	beego.Router("/ping", &controllers.MainController{}, "get:Ping")
	//v1后台接口list
	v1ApplicationNs := beego.NewNamespace("/v1",
		//node
		beego.NSRouter("/nodes", &controllers.NodeController{}, "get:GetNodes"),
		beego.NSRouter("/node/:nodeName", &controllers.NodeController{}, "get:GetNode"),
		beego.NSRouter("/node/:nodeName/start", &controllers.NodeController{}, "post:StartNode"),
		beego.NSRouter("/node/:nodeName/stop", &controllers.NodeController{}, "post:StopNode"),

		//deployment
		beego.NSRouter("/deployments", &controllers.DeployController{}, "get:GetDeployments"),
		beego.NSRouter("/deployment/update", &controllers.DeployController{}, "post:UpdateDeployment"),
		beego.NSRouter("/deployment/:deploymentName/scale", &controllers.DeployController{}, "post:ScaleDeployment"),

		//pods
		beego.NSRouter("/pods", &controllers.PodController{}, "get:GetPods"),

		beego.NSRouter("/dashboard", &controllers.DashboardController{}, "get:GetDashboardInfo"),
	)

	beego.AddNamespace(v1ApplicationNs)

	//过滤越权请求
	var FilterUser = func(ctx *context.Context) {
		uri := ctx.Request.RequestURI

		if !strings.Contains(uri, "/login") && !strings.Contains(uri, "/ping") {
			authorization := ctx.Request.Header.Get("Authorization")
			if myutils.CheckAuthorization(authorization) != 0 {
				returnData := map[string]interface{}{
					"code": 401,
					"msg":  "未登录",
					"data": map[string]interface{}{},
				}
				ctx.Output.JSON(returnData, false, false)
			}

		}

	}
	beego.InsertFilter("/v1/*", beego.BeforeExec, FilterUser)

	//beego.Router("/", &controllers.MainController{})
	//beego.Router("/getPods", &controllers.MainController{}, "*:GetPods")
	//beego.Router("/getServices", &controllers.MainController{}, "*:GetServices")
	//beego.Router("/getDeployments", &controllers.MainController{}, "*:GetDeployments")
	//beego.Router("/getNodes", &controllers.MainController{}, "*:GetNodes")
	//beego.Router("/getNode", &controllers.MainController{}, "*:GetNode")
	//
	//beego.Router("/scaleService", &controllers.MainController{}, "*:ScaleService")
	//beego.Router("/stopNode", &controllers.MainController{}, "*:StopNode")

}
