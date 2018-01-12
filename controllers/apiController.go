package controllers

import (
	"github.com/astaxie/beego"
	"onesteam/library"
	"onesteam/services/request/steam"
	"onesteam/models/users"
)


type GameInfoController struct {
	beego.Controller
}

type SteamUserInfoController struct {
	beego.Controller
}

func (c *GameInfoController) Get() {
	steamApi := steam.NewSteam()
	getMatch, _ := steamApi.GetGameByGameId("3655921259")
	output, _ := library.ReturnJsonWithError(0, "", getMatch)
	c.Ctx.WriteString(output)
}

func (c *SteamUserInfoController) Get() {
	//getPlayer, _ := users.BindSteamId("124543174")
	steamApi := steam.NewSteam()
	steamIdArr := []string{"154710346", "4294967295", "313272154", "124543174"}
	steamIdArr = steamApi.ConvertSteamIdTo64Bit(steamIdArr)
	getPlayer, _ := steamApi.GetUserInfoBysteamIds(steamIdArr)
	output, _ := library.ReturnJsonWithError(0, "", getPlayer)
	c.Ctx.WriteString(output)
}




