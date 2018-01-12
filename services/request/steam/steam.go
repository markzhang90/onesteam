package steam

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"onesteam/services/request"
	"onesteam/library"
	"fmt"
	"strings"
	"strconv"
)

var (
	dota2ApiHost = "https://api.steampowered.com/IDOTA2Match_570/"
	steamApiHost = "http://api.steampowered.com/"
)

type (
	steamApi struct {
		apiKey string
	}
)

func NewSteam() *steamApi {
	thirdConf, err := config.NewConfig("ini", "conf/third.conf")

	if err != nil {
		logs.Warn(err)
		panic(err.Error())
	}

	sectionMap, err := thirdConf.GetSection("steam")
	if err != nil {
		logs.Warn(err)
		return nil
	}

	if getApiKey, ok := sectionMap["api_key"]; ok {
		return &steamApi{getApiKey}
	}
	return nil
}

func (steamApi *steamApi) GetGameByGameId(gameId string) (map[string]interface{}, error) {
	dota2ApiGetMatchDetails := dota2ApiHost + "GetMatchDetails/V001/"
	var queryMap = make(map[string]string)
	queryMap["match_id"] = gameId
	queryMap["key"] = steamApi.apiKey
	getRes, errRes := request.HttpGet(dota2ApiGetMatchDetails, queryMap)
	if errRes != nil {
		logs.Warn(errRes)
		return nil, errRes
	}

	decodeJson, errDecode := library.Json2Map(getRes)
	if errDecode != nil {
		logs.Warn(errDecode)
		return nil, errDecode
	}
	if _, ok := decodeJson["result"]; !ok {
		return nil, fmt.Errorf("no result")
	}

	result, ok := decodeJson["result"].(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("no result")
	}

	if _, ok := result["error"]; ok {
		return nil, fmt.Errorf(result["error"].(string))
	}

	return result, nil
}

func (steamApi *steamApi) ConvertSteamIdTo64Bit(steamIdArr []string) []string {
	var steamArr = make([]string, len(steamIdArr))
	for _, value := range steamIdArr {
		intVal, _ := strconv.Atoi(value)
		valueNew := 1<<56 | 1<<52 | 1<<32 | intVal
		steamArr = append(steamArr, strconv.Itoa(valueNew))
	}
	return steamArr
}

func (steamApi *steamApi) GetUserInfoBysteamIds(steamArr []string) ([]interface{}, error) {
	steamApiGetPlayerSummaries := steamApiHost + "ISteamUser/GetPlayerSummaries/v0002/"
	var queryMap = make(map[string]string)
	steamIds := strings.Join(steamArr, ",")
	queryMap["steamids"] = steamIds
	queryMap["key"] = steamApi.apiKey
	getRes, errRes := request.HttpGet(steamApiGetPlayerSummaries, queryMap)

	if errRes != nil {
		return nil, errRes
	}

	decodeJson, errDecode := library.Json2Map(getRes)
	if errDecode != nil {
		logs.Warn(errDecode)
		return nil, errDecode
	}
	if _, ok := decodeJson["response"]; !ok {
		return nil, fmt.Errorf("no result")
	}
	response, ok := decodeJson["response"].(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("no result")
	}

	if _, ok := response["players"]; !ok {
		return nil, fmt.Errorf("no plays")
	}

	playList := response["players"].([]interface{})
	return playList, nil
}
