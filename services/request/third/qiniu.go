package third

import (
	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
	"qiniupkg.com/api.v7/kodocli"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/config"
	"qiniupkg.com/x/errors.v7"
)

var (
	// 设置上传到的空间
	bucket = "onestory"
	buckeOpen = "onestory-pub"
	domain = "ormnna350.bkt.clouddn.com"
	domainOpen = "orrxp85k4.bkt.clouddn.com"
)

// 构造返回值字段
type (
	PutRet struct {
		Hash string `json:"hash"`
		Key  string `json:"key"`
	}

	Qiniu struct {
		bucket string
		domain string
	}
)

func NewQiNiu(open bool) *Qiniu {
	mybucket := bucket
	mydomain := domain
	if open {
		mybucket = buckeOpen
		mydomain = domainOpen
	}

	return &Qiniu{mybucket, mydomain}
}

func (ins *Qiniu) DownloadUrl(key string)string {
	// 调用MakeBaseUrl()方法将domain,key处理成http://domain/key的形式
	baseUrl := kodo.MakeBaseUrl(ins.domain, key)

	if ins.bucket == buckeOpen {
		return baseUrl
	}
	policy := kodo.GetPolicy{}
	// 生成一个client对象
	c := kodo.New(0, nil)

	// 调用MakePrivateUrl方法返回url
	return c.MakePrivateUrl(baseUrl, &policy)
}

func (ins *Qiniu) Upoloader(filepath string) (string, error){
	// 初始化AK，SK

	thirdConf, err := config.NewConfig("ini", "conf/third.conf")

	if err != nil {
		logs.Warn(err)
		panic(err.Error())
	}

	accessKey := thirdConf.String("qiniu::accesskey")
	secretKey := thirdConf.String("qiniu::secretkey")

	conf.ACCESS_KEY = accessKey
	conf.SECRET_KEY = secretKey

	// 创建一个Client
	c := kodo.New(0, nil)
	// 设置上传的策略
	policy := &kodo.PutPolicy{
		Scope: ins.bucket,
		//设置Token过期时间
		Expires: 3600,
	}
	// 生成一个上传token
	token := c.MakeUptoken(policy)
	// 构建一个uploader
	zone := 0
	uploader := kodocli.NewUploader(zone, nil)

	var ret PutRet
	// 设置上传文件的路径
	//filepath := "/Users/dxy/sync/sample2.flv"
	// 调用PutFileWithoutKey方式上传，没有设置saveasKey以文件的hash命名
	res := uploader.PutFileWithoutKey(nil, &ret, token, filepath, nil)
	// 打印返回的信息
	fmt.Println(ret)
	// 打印出错信息
	if res != nil {
		fmt.Println("io.Put failed:", res.Error())
		return "", errors.New(res.Error())
	}

	key := ret.Key
	return key, nil
}
