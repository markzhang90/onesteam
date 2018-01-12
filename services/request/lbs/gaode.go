package lbs

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/config"
	"qiniupkg.com/x/errors.v7"
	"onestory/library"
	"onestory/services/request"
)

const weatherApi = "http://restapi.amap.com/v3/weather/weatherInfo"
const locationApi = "http://restapi.amap.com/v3/ip"

/**
获取天气信息
 */
func GetWeatherByLocation(location string) (interface{}, error) {

	secreteKey, err := getSecreteKey()

	if err != nil {
		return "", err
	}
	var requestVars = make(map[string]string)

	requestVars["key"] = secreteKey;
	requestVars["city"] = location;
	stringRes, err := request.HttpGet(weatherApi, requestVars)
	if err != nil {
		return nil, err
	}

	mapRes, err := library.Json2Map(stringRes)
	if err != nil {
		return nil, err
	}

	if _, ok := mapRes["status"].(string); !ok {

		return nil, errors.New("get weather fail")
	}

	if "1" != mapRes["status"] {
		return nil, errors.New("get weather fail")
	}
	var mapLives = make([]map[string]interface{}, 0)
	for _, value := range mapRes["lives"].([]interface {}) {
		mapValue , _ := value.(map[string]interface {})
		mapValue["currentCity"] = mapValue["city"]
		mapValue["weatherDesc"] = mapValue["weather"]
		mapLives = append(mapLives, mapValue)
	}

	return mapLives, nil
}


func GetLocationByIp(ip string) (string, error) {
	secreteKey, err := getSecreteKey()

	if err != nil {
		return "", err
	}
	var requestVars = make(map[string]string)

	requestVars["key"] = secreteKey;
	requestVars["ip"] = ip;
	stringRes, err := request.HttpGet(locationApi, requestVars)
	if err != nil {
		return "", err
	}

	mapRes, err := library.Json2Map(stringRes)
	if err != nil {
		return "", err
	}

	cityRes, ok := mapRes["city"].(string)
	if ok{
		return cityRes, err
	}else{
		return "", errors.New("get city fail")
	}

}
func getSecreteKey() (string, error) {

	thirdConf, err := config.NewConfig("ini", "conf/third.conf")

	if err != nil {
		logs.Warning(err)
		return "", err
	}

	screctKeyMap, err := thirdConf.GetSection("gaode")

	screctKey := screctKeyMap["apikey"]

	if err != nil || len(screctKey) <= 1 {
		return "", errors.New("key missing")
	}
	return screctKey, nil
}

