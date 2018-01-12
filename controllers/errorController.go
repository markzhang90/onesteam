package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error404() {
	c.Data["content"] = "page not found"
	c.TplName = "index.tpl"
	beego.Warning("this is 404 here")
}

func (c *ErrorController) Error500() {
	c.Data["content"] = "internal server error"
	c.TplName = "index.tpl"
}

func (c *ErrorController) ErrorDb() {
	c.Data["content"] = "database is now down"
	c.TplName = "index.tpl"
}