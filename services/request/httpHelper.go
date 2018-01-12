package request

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"fmt"
	"strings"
	"net"
	"time"
	"github.com/astaxie/beego/logs"
)

var (
	client *http.Client
)

func init()  {

	client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*2)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 2))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 2,
		},
	}

}

// 简单直接的GET请求
func HttpGet(urlStr string, queryList map[string]string) (string, error){
	var temp = make([]string, 0, len(queryList))

	for key, value := range queryList{
		stringQuery :=  key + "=" + value
		temp = append(temp, stringQuery)
	}

	queryVar := strings.Join(temp, "&")
	fullQuery := urlStr + "?" + queryVar
	resp, err := client.Get(fullQuery)
	if err != nil {
		return "", err
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Warning("request fail for " + fullQuery + " " + err.Error())
		return "", err
		// handle error
	}
	return string(body), nil
}

// POST请求 -- 使用http.Post()方法
//Tips：使用这个方法的话，第二个参数要设置成”application/x-www-form-urlencoded”，否则post参数无法传递。

func HttpPost(urlStr string, formList map[string]string) string{

	var temp []string
	for key, value := range formList {
		temp = append(temp, key + "=" + value)
	}
	implodedStr := strings.Join(temp, "&")
	logs.Warn(implodedStr)
	resp, err := client.Post(urlStr,
		"application/x-www-form-urlencoded",
		strings.NewReader(implodedStr))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	return string(body)
}

// POST请求 -- 使用http.PostForm()方法
func HttpPostForm(urlStr string) {
	resp, err := client.PostForm(urlStr,
		url.Values{"key": {"Value"}, "id": {"123"}})

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

// 复杂的请求（设置头参数、cookie之类的数据），可以使用http.Client的Do()方法</strong>
func HttpDo() {

	req, err := http.NewRequest("POST", "http://www.baidu.com", strings.NewReader("name=cjb"))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}