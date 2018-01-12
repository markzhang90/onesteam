package request

import (
	"errors"
	"onestory/library"
	"github.com/astaxie/beego/logs"
)

func GetWeatherInfo(location string) (interface{}, error) {
	if len(location) < 1 {
		return "", errors.New("location missing")
	}
	var requestVars = make(map[string]string)
	requestVars["key"] = "77514aacee204dc697a27743f714d434";
	requestVars["cityname"] = location;
	stringRes, err := HttpGet("http://api.avatardata.cn/Weather/Query", requestVars)
	if err != nil{
		return nil, err
	}
	mapRes, err := library.Json2Map(stringRes)
	if err != nil{
		return nil, err
	}
	logs.Warn(mapRes["error_code"])
	if mapRes["error_code"].(float64) != 0 {
		return nil, errors.New(mapRes["reason"].(string))
	}
	return mapRes["result"], nil
}


