package third

import (
	"errors"
	"onestory/library"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/config"
	"onestory/services/request"
)

type (
	WehcatSmallApp struct {
		serviceName string
		conf map[string]string
	}

)

func NewWechatSmallApp(confName string)  *WehcatSmallApp{
	thirdConf, err := config.NewConfig("ini", "conf/third.conf")

	if err != nil {
		logs.Warn(err)
		panic(err.Error())
	}

	sectionMap, err := thirdConf.GetSection(confName)
	if err != nil{
		logs.Warn(err)
		return nil
	}
	return &WehcatSmallApp{serviceName:confName, conf:sectionMap}
}


func (wechat *WehcatSmallApp)GetLoginOpenIdFronCode(code string) (interface{}, error) {
	if len(code) < 1 {
		return nil, errors.New("code missing")
	}
	var requestVars = make(map[string]string)

	if _, ok := wechat.conf["appid"]; ok {
		requestVars["appid"] = wechat.conf["appid"]
	}else{
		return nil, errors.New("appid missing")
	}

	if _, ok := wechat.conf["grant_type"]; ok {
		requestVars["grant_type"] = wechat.conf["grant_type"]
	}else{
		return nil, errors.New("grant_type missing")
	}

	if _, ok := wechat.conf["secret"]; ok {
		requestVars["secret"] = wechat.conf["secret"]
	}else{
		return nil, errors.New("secret missing")
	}

	requestVars["js_code"] = code

	stringRes, err := request.HttpGet("https://api.weixin.qq.com/sns/jscode2session", requestVars)
	if err != nil{
		return nil, err
	}
	mapRes, err := library.Json2Map(stringRes)

	if err != nil{
		return nil, err
	}
	if _, ok := mapRes["openid"]; ok{
		return mapRes, nil
	}

	errMsg := "获取openid失败"
	if _,ok := mapRes["errmsg"]; ok{
		errMsg = mapRes["errmsg"].(string)
	}

	return mapRes, errors.New(errMsg)
}
