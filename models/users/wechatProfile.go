package users

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
	"encoding/json"
	"onesteam/services/rediscli"
	"errors"
	"fmt"
)

func (userDb *UserProfileDb) GetUserByOpenIdOrCreate(openId string) (targetUser UserProfile, err error) {

	cacheUser, errCache := userDb.GetUserFromCacheForOpenId(openId)
	if errCache == nil {
		return cacheUser.UserProfile, nil
	}

	if errCache == orm.ErrNoRows {

		//create new
		userData := UserProfile{
			Passid:      openId,
			Email:       "",
			Phone:       0,
			Openid:      openId,
			Thirdid:      "",
			Password:    "",
			Update_time: time.Now().Unix(),
			Nick_name:   "",
			Ext:         "",
		}
		uid, errAdd := userDb.AddNewUserProfile(userData)
		if errAdd != nil {
			return targetUser, errAdd
		}
		userData.Id = int(uid)
		SyncSetUserCache(userData, true)
		return userData, nil
	}

	if err != nil {
		return targetUser, err
	}

	return targetUser, nil
}

func (userDb *UserProfileDb) GetUserProfileByOpenId(openId string) (targetUser UserProfile, err error) {
	o := userDb.DbConnect.Orm
	o.Using(userDb.DbConnect.DbName)
	targetUser = UserProfile{Openid: openId}
	err = o.Read(&targetUser, "openid")
	if err != nil {
		logs.Warning("get user fail " + err.Error())
		return targetUser, err
	}
	return targetUser, nil
}


func (userDb *UserProfileDb)GetUserFromCacheForOpenId(openId string) (UserCache, error) {

	var userCache UserCache
	redisCacheKey := getUserCacheKey(openId, "openid")

	redsiConn := rediscli.RedisClient.Get()
	res, errCache := redsiConn.Do("Get", redisCacheKey)
	defer redsiConn.Close()

	if errCache != nil || res == nil{
		//try to get user from db

		targetUser, err := userDb.GetUserProfileByOpenId(openId)
		if err !=nil {
			logs.Warn("GetUser Fail" + err.Error())
			return userCache, err
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

func (userDb *UserProfileDb)BindSteamAccount(openId, steamId string) (UserProfile, error) {
	ch := make(chan interface{})
	defer close(ch)
	var user UserProfile
	go GetSteamInfoById(steamId, ch)
	wechatUser, errUser := userDb.GetUserFromCacheForOpenId(openId)
	if errUser != nil{
		return user, errUser
	}
	user = wechatUser.UserProfile
	steamUser := <-ch
	if steamUser == nil{
		return user, fmt.Errorf("get steam info fail")
	}
	steamUserDe, _ := steamUser.(SteamUser)
	user.Thirdid = steamUserDe.Steamid
	_, errUpdate := userDb.UpdateNewUserProfile(user)
	if errUpdate != nil {
		return user, errUpdate
	}
	return user, nil
}

