package main

import (
	"github.com/astaxie/beego/plugins/cors"
	"time"
	_ "wac/routers"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		//AllowMethods:    []string{"GET", "POST"},
		//AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		//AllowHeaders: []string{"*"},
		//ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		//ExposeHeaders:    []string{"*"},
		//AllowCredentials: true,
		MaxAge: 999999 * time.Hour,
	}))

	beego.Run()
}
