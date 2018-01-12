package users

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"onesteam/services"
	"crypto/md5"
	"encoding/hex"
	"onesteam/services/rediscli"
	"encoding/json"
	"time"
	"errors"
)

type (
	UserProfileDb struct {
		tableName string
		DbConnect *services.DbService
	}

	UserProfile struct {
		Id          int `orm:"auto"`
		Passid      string
		Openid      string
		Thirdid      string
		Email       string
		Avatar      string
		Phone       int64
		Password    string
		Update_time int64
		Nick_name   string
		Ext         string
		Active         int64
	}

	UserCache struct {
		UserProfile
		LastLogin int64
	}

)



func NewUser() (*UserProfileDb) {

	dbService, err := services.NewService("onesteam")
	if err != nil{
		logs.Warn(err)
	}
	return &UserProfileDb{"user_profile", dbService}
}

func (userDb *UserProfileDb) LoginUserByEmail(email string, password string) (UserProfile, error)  {

	o := userDb.DbConnect.Orm
	o.Using(userDb.DbConnect.DbName)
	targetUser := UserProfile{Email: email, Password:encriptPass(password)}
	errdb := o.Read(&targetUser, "email","password")

	if errdb != nil {
		logs.Warning("LoginUser fail " + errdb.Error())
		return targetUser, errdb
	}

	return targetUser, nil
}

func (userDb *UserProfileDb) LoginUserByPhone(phone int64, password string) (UserProfile, error)  {

	o := userDb.DbConnect.Orm
	o.Using(userDb.DbConnect.DbName)
	targetUser := UserProfile{Phone: phone, Password:encriptPass(password)}
	errdb := o.Read(&targetUser, "phone","password")

	if errdb != nil {
		logs.Warning("LoginUser fail " + errdb.Error())
		return targetUser, errdb
	}

	return targetUser, nil
}


func (userDb *UserProfileDb) GetUserProfileByPhone(phone int64) (targetUser UserProfile, err error) {
	o := userDb.DbConnect.Orm
	o.Using(userDb.DbConnect.DbName)
	targetUser = UserProfile{Phone: phone}
	err = o.Read(&targetUser, "phone")

	if err != nil {
		logs.Warning("get user fail " + err.Error())
		return targetUser, err
	}
	return targetUser, nil
}


func (userDb *UserProfileDb) GetUserProfileByEmail(email string) (targetUser UserProfile, err error) {
	o := userDb.DbConnect.Orm
	o.Using(userDb.DbConnect.DbName)
	targetUser = UserProfile{Email: email}
	err = o.Read(&targetUser, "email")
	if err != nil {
		logs.Warning("get user fail " + err.Error())
		return targetUser, err
	}
	return targetUser, nil
}

func (userDb *UserProfileDb) GetUserProfileByPassId(passId string) (targetUser UserProfile, err error) {
	o := userDb.DbConnect.Orm
	o.Using(userDb.DbConnect.DbName)
	targetUser = UserProfile{Passid: passId}
	err = o.Read(&targetUser, "passid")
	if err != nil {
		logs.Warning("get user fail " + err.Error())
		return targetUser, err
	}
	return targetUser, nil
}

func (userDb *UserProfileDb) GetUserProfileById(uid int) (targetUser UserProfile, err error) {
	o := userDb.DbConnect.Orm
	o.Using(userDb.DbConnect.DbName)
	targetUser = UserProfile{Id: uid}
	err = o.Read(&targetUser)
	if err != nil {
		logs.Warning("get user fail " + err.Error())
		return targetUser, err
	}
	return targetUser, nil
}

func (userDb *UserProfileDb) AddNewUserProfile(userprofileData UserProfile)(int64, error){
	//check email
	_, errCheck := userDb.GetUserProfileByEmail(userprofileData.Email)
	if !(errCheck != nil && errCheck == orm.ErrNoRows) {
		return 0, errors.New("user exist")
	}
	o := userDb.DbConnect.Orm
	o.Using(userDb.DbConnect.DbName)
	profile := new(UserProfile)
	profile.Passid = userprofileData.Passid
	profile.Openid = userprofileData.Openid
	profile.Thirdid = userprofileData.Thirdid
	profile.Email = userprofileData.Email
	profile.Phone = userprofileData.Phone
	profile.Avatar = userprofileData.Avatar
	profile.Update_time = time.Now().Unix()
	profile.Nick_name = userprofileData.Nick_name
	profile.Password = encriptPass(userprofileData.Password)
	profile.Ext = userprofileData.Ext
	profile.Active = userprofileData.Active
	res, err := o.Insert(profile)

	if err != nil {
		logs.Warning("add user fail " + err.Error())
		return res, err
	}else{
		logs.Trace("add user succ " + string(res))
	}
	return res, nil
}



func (userDb *UserProfileDb) UpdateNewUserProfile(userprofileData UserProfile) (UserProfile, error) {
	o := userDb.DbConnect.Orm
	o.Using(userDb.DbConnect.DbName)

	requireUpdate := false

	var profile UserProfile
	profile.Id = userprofileData.Id
	//profile.Passid = userprofileData.Passid
	//profile.Email = userprofileData.Email
	//profile.Phone = userprofileData.Phone

	if o.Read(&profile) == nil {

		//update fields
		profile.Update_time = time.Now().Unix()
		if len(userprofileData.Nick_name) > 0 {
			requireUpdate = true
			profile.Nick_name = userprofileData.Nick_name
		}
		if len(userprofileData.Password) > 0 {
			requireUpdate = true
			profile.Password = encriptPass(userprofileData.Password)
		}
		if len(userprofileData.Ext) > 0 {
			requireUpdate = true
			profile.Ext = userprofileData.Ext
		}
		if len(userprofileData.Avatar) > 0 {
			requireUpdate = true
			profile.Avatar = userprofileData.Avatar
		}
		if len(userprofileData.Thirdid) > 0 {
			requireUpdate = true
			profile.Thirdid = userprofileData.Thirdid
		}
		if !requireUpdate {
			return profile, nil
		}
		if num, err := o.Update(&profile); err == nil && num == 1{
			logs.Trace(string(profile.Id) + " update userprofile succ " + string(num))
			return profile, nil
		}else{
			logs.Warning(string(profile.Id) + " update userprofile fail " + err.Error())
			return profile, err
		}

	}else{
		logs.Warning("get user fail " + string(profile.Id))
		return profile, errors.New("获取用户失败")
	}
	return profile, errors.New("更新用户失败")
}

func (userDb *UserProfileDb) ActiveUserProfile(userprofileData UserProfile) (UserProfile, error) {
	o := userDb.DbConnect.Orm
	o.Using(userDb.DbConnect.DbName)

	var profile UserProfile
	profile.Id = userprofileData.Id


	if o.Read(&profile) == nil {

		//update fields
		profile.Update_time = time.Now().Unix()
		profile.Active = 1
		if num, err := o.Update(&profile); err == nil && num == 1{
			logs.Trace(string(profile.Id) + " update userprofile succ " + string(num))
			return profile, nil
		}else{
			logs.Warning(string(profile.Id) + " update userprofile fail " + err.Error())
			return profile, err
		}

	}else{
		logs.Warning("get user fail " + string(profile.Id))
		return profile, errors.New("获取用户失败")
	}
	return profile, errors.New("更新用户失败")
}



func (userDb *UserProfileDb) GetUserProfile() (err error) {
	o := userDb.DbConnect.Orm
	o.Using(userDb.DbConnect.DbName)

	var maps []orm.Params
	res, err := o.Raw("select * from user_profile where nick_name = ?", "oooook").Values(&maps)

	if err == nil && res > 0 {
		//data := maps[0]["email"]
		//logs.Warning(data)
		for key, v := range maps {
			logs.Warning(key)
			logs.Warning(v)
		}
	}

	return err
}

func (userDb *UserProfileDb)ClearProfileOut(userProfile UserProfile) (userPubProfile map[string]interface{}) {
	var pubUserInfo = make(map[string]interface{})
	pubUserInfo["Passid"] = userProfile.Passid
	pubUserInfo["Thirdid"] = userProfile.Thirdid
	pubUserInfo["Nick_name"] = userProfile.Nick_name
	return pubUserInfo
}

func GetPid(phone int64, email string) string {
	passIdEncode := md5.New()
	passIdEncode.Write([]byte(string(phone) + "_" + email))
	passId := hex.EncodeToString(passIdEncode.Sum(nil))
	return passId
}

func encriptPass(password string)  string{
	passWordEncode := md5.New()
	passWordEncode.Write([]byte(password))
	passWord := hex.EncodeToString(passWordEncode.Sum(nil))
	return passWord
}

/**
stringType passid openid
 */
func getUserCacheKey(prefix string, stringType string) string{
	switch stringType {
	case "passid":
		return "passid:"+prefix
	case "openid":
		return "openid"+prefix
	default:
		return "default"+prefix
	}
}

func SyncSetUserCache(userObj UserProfile, usingOpenId bool) (UserCache, bool) {

	var redisCacheKey string
	var cacheType string

	if usingOpenId {
		cacheType = "openid"
		redisCacheKey = getUserCacheKey(userObj.Openid, cacheType)
	}else{
		cacheType = "passid"
		redisCacheKey = getUserCacheKey(userObj.Passid, cacheType)
	}

	redsiConn := rediscli.RedisClient.Get()
	userObj.Password = ""

	var userCache UserCache
	userCache.UserProfile = userObj
	userCache.LastLogin = time.Now().Unix()
	jsonUser, err := json.Marshal(userCache)

	if err != nil {
		userCache.LastLogin = -1
		logs.Warn("SyncSetUserCache Fail" + userObj.Passid)
		return userCache, false
	}
	res, errCache := redsiConn.Do("SET", redisCacheKey, jsonUser)
	//expire user
	expireTime, confErr := beego.AppConfig.Int("redisuserexpire")
	if confErr != nil {
		expireTime = 60*60*24*30
	}
	redsiConn.Do("EXPIRE", userObj.Passid,expireTime)

	defer redsiConn.Close()
	if errCache == nil || res == "OK"{
		return userCache, true
	}else {
		userCache.LastLogin = -1
	}
	logs.Warn("SyncSetUserCache Fail" + errCache.Error())
	return userCache, false
}

func CleanUserCache(passId string) (bool, error) {

	redsiConn := rediscli.RedisClient.Get()
	cacheTypes := []string{"passid", "openid"}
	for _, cacheType := range cacheTypes  {
		redisCacheKey := getUserCacheKey(passId, cacheType)
		_, errCache := redsiConn.Do("DEL", redisCacheKey)
		if errCache != nil{
			return false, errCache
		}
	}

	return true, nil
}

func GetUserFromCache(passId string, actived bool) (UserCache, error) {

	var userCache UserCache
	redisCacheKey := getUserCacheKey(passId, "passid")
	redsiConn := rediscli.RedisClient.Get()
	res, errCache := redsiConn.Do("Get", redisCacheKey)
	defer redsiConn.Close()
	if errCache != nil || res == nil{
		//try to get user from db
		var newUserDb = NewUser()
		targetUser, err := newUserDb.GetUserProfileByPassId(passId)
		if err !=nil {
			logs.Warn("GetUser Fail " + passId + "  " + err.Error())
			return userCache, errCache
		}

		if actived && targetUser.Active != 1{
			logs.Warn("User is not active ")
			return userCache, errors.New("user unactivate")
		}
		if actived {
			SyncSetUserCache(targetUser, false)
		}
		userCache.UserProfile = targetUser
		return userCache, nil
	}

	if jsonRes, ok := res.([]byte); !ok {
		return userCache, errors.New("获取用户失败")
	} else {
		json.Unmarshal([]byte(jsonRes), &userCache)
		return userCache, nil
	}
}

func ClearOutputUserprofile(profile UserProfile) map[string]interface{} {
	var userInfo = make(map[string]interface{})
	userInfo["Nick_name"] = profile.Nick_name
	userInfo["Thirdid"] = profile.Thirdid
	userInfo["Avatar"] = profile.Avatar
	userInfo["Email"] = profile.Email
	userInfo["Active"] = profile.Active
	userInfo["Update_time"] = profile.Update_time
	return userInfo
}

