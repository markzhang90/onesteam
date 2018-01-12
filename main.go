package main

import (
	_ "onesteam/routers"
	"github.com/astaxie/beego"
	"time"
	"github.com/astaxie/beego/logs"
	"onesteam/controllers"
)

func main() {
	beego.ErrorController(&controllers.ErrorController{})

	//err log config
	projectName := "./logs/" + beego.AppConfig.String("appname") + "." + time.Now().Format("2006-01-02-15")
	beego.SetLogger(logs.AdapterFile, `{"filename":"`+projectName+`","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":3}`)

	beego.Run()
}

