package routers

import (
	"onesteam/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/steam/gamedetail", &controllers.GameInfoController{})
	beego.Router("/steam/user", &controllers.SteamUserInfoController{})
}
