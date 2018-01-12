package users

import (
	"onesteam/services/request/steam"
	"fmt"
)

type(
	SteamUser struct{
		Avatar string
		Lastlogoff string
		Steamid string
		Personaname string
	}
)


func GetSteamInfoById(steamId string, ch chan interface{}) (SteamUser, error){
	//get steam info
	var player SteamUser

	steamApi := steam.NewSteam()
	steamIdArr := []string{steamId}
	steamIdArr = steamApi.ConvertSteamIdTo64Bit(steamIdArr)
	players, errPlayers := steamApi.GetUserInfoBysteamIds(steamIdArr)

	if errPlayers != nil {
		ch <- nil
		return player, errPlayers
	}
	if len(players) != 1 {
		ch <- nil
		return player, fmt.Errorf("get steam user fail")
	}

	getPlay := players[0].(map[string]interface{})
	//myGetPlay := getPlay.(map[string]interface{})

	if _, ok := getPlay["steamid"] ; !ok{
		ch <- nil
		return player, fmt.Errorf("get steam id fail")
	}

	player.Avatar = getPlay["avatar"].(string)
	player.Lastlogoff = fmt.Sprintf("%.0f", getPlay["lastlogoff"].(float64))
	player.Steamid = getPlay["steamid"].(string)
	player.Personaname = getPlay["personaname"].(string)
	ch <- player
	return player, nil
}

func getUserAllSync(openid string)  {
	
}