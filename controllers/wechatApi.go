package controllers

import (
	"github.com/astaxie/beego"
	"onesteam/library"
	"onesteam/services/request/third"
	"github.com/astaxie/beego/logs"
	"time"
	"strconv"
	"onestory/models"
	"onesteam/models/users"
)

type (
	LoginWehchatController struct {
		beego.Controller
	}
	InitWehchatController struct {
		beego.Controller
	}
)

func (c *LoginWehchatController) Get()  {
	code := c.GetString("code", "")
	if len(code) < 1 {
		output, _ := library.ReturnJsonWithError(1, "获取code失败", "")
		c.Ctx.WriteString(output)
		return
	}

	wechat := third.NewWechatSmallApp("wechat-smallapp")
	if wechat == nil {
		output, _ := library.ReturnJsonWithError(1, "获取配置失败", "")
		c.Ctx.WriteString(output)
		return
	}
	//call wehchat api
	callRes, err := wechat.GetLoginOpenIdFronCode(code);

	if err != nil{
		output, _ := library.ReturnJsonWithError(library.CodeErrApi, "微信登录失败", callRes)
		c.Ctx.WriteString(output)
		return
	}

	var openid string
	if weChatBack , ok := callRes.(map[string]interface{}); ok{
		openid = weChatBack["openid"].(string)
		logs.Warn(openid)
	}else{
		output, _ := library.ReturnJsonWithError(library.CodeErrApi, "微信登录失败", callRes)
		c.Ctx.WriteString(output)
		return
	}

	//get userInfo by openId

	userDb := users.NewUser()
	userprofile, errGetDb := userDb.GetUserByOpenIdOrCreate(openid)

	if err != nil{
		output, _ := library.ReturnJsonWithError(library.CodeErrApi, "微信登录失败", errGetDb.Error())
		c.Ctx.WriteString(output)
		return
	}

	clearRes := userDb.ClearProfileOut(userprofile)

	output, _ := library.ReturnJsonWithError(0, "", clearRes)
	c.Ctx.WriteString(output)
	return
}


/**
we chat init info
 */
func (c *InitWehchatController) Get() {

	cookiekey := beego.AppConfig.String("passid")
	//get from cache
	passId, _ := c.GetSecureCookie(cookiekey, "passid")

	if len(passId) <= 0 {
		passId = c.GetString("passid", "")
		if len(passId) < 1{
			output, _ := library.ReturnJsonWithError(library.GetUserFail, "ref", nil)
			c.Ctx.WriteString(output)
			return
		}
	}

	cahchedUser, err := models.GetUserFromCache(passId, false)
	if err != nil {
		output, _ := library.ReturnJsonWithError(library.GetUserFail, "ref", nil)
		c.Ctx.WriteString(output)
		return
	}

	var clearRes = make(map[string]interface{})
	var userInfo = make(map[string]string)
	userInfo["Nick_name"] = cahchedUser.Nick_name
	userInfo["Avatar"] = cahchedUser.Avatar
	clearRes["User_info"] = userInfo
	var uid = cahchedUser.Id
	userPost := models.NewPost()
	countAll, err := userPost.QueryCountUserPost(uid);

	if err != nil {
		countAll = -1
	}

	clearRes["Post_count"] = countAll

	var today = time.Now().Format("20060102");

	todayInt, _ := strconv.Atoi(today)

	todayArr := []int{todayInt}

	result, errGet := userPost.QueryUserPostByDate(cahchedUser.Id, todayArr, true, 1);

	clearRes["Today"] = false;
	clearRes["Id"] = -1;

	if errGet == nil {
		if len(result) > 0 {
			clearRes["Today"] = true;
			clearRes["Id"] = result[0].Id;
		}
	}

	output, _ := library.ReturnJsonWithError(0, "", clearRes)
	c.Ctx.WriteString(output)
	return
}