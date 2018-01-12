package rediscli

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego"
	"time"
)

//connection pool
var (
	// 定义常量
	RedisClient *redis.Pool
	REDIS_HOST  string
)

func getFullName() string {
	return "redis_" + beego.AppConfig.String("appname")
}

func init() {
	// 从配置文件获取redis的ip以及db
	dbconf, err := config.NewConfig("ini", "conf/db.conf")

	if err != nil {
		logs.Warn(err)
		panic(err)
	}
	fullName := getFullName()
	host := dbconf.String(fullName + "::host")
	port := dbconf.String(fullName + "::port")
	if len(host) == 0 || len(port) == 0 {
		logs.Warn(err)
		panic(err)
	}

	var password = dbconf.String(fullName + "::password")

	REDIS_HOST = host + ":" + port
	// 建立连接池
	RedisClient = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     dbconf.DefaultInt(fullName+"::maxidle", 1),
		MaxActive:   dbconf.DefaultInt(fullName+"::maxactive", 10),
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {

			c, err := redis.Dial(
				"tcp",
				REDIS_HOST,
				redis.DialConnectTimeout(1*time.Second),
				redis.DialReadTimeout(1*time.Second),
				redis.DialWriteTimeout(1*time.Second),
			)
			if err != nil {
				logs.Warning("hahahhaqqqq" + err.Error())
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			// 选择db
			return c, nil
		},
	}
}

