package third

import (
	"onestory/library"
	"time"
	"sort"
	"strings"
	"qiniupkg.com/x/url.v7"
	"crypto/hmac"
	"crypto/sha1"
	"onestory/services/request"
	"encoding/base64"
	"github.com/astaxie/beego/logs"
	"net/smtp"
)

var (
	accessKeyId = "LTAIGhsbZyLlSFvk"
	accessKeySecret = "KlPlZ7NTKj9X7WSCM6QH9JrcbC6OcI"
	user = "service@mail.onestory.cn"
	password = "Zyy45612301Mark"
	host = "smtpdm.aliyun.com:80"
)

func AliyunApiCommon() map[string]string {
	var requestParam  = make(map[string]string)
	requestParam["Format"] = "JSON"
	requestParam["Version"] = "2015-11-23"
	requestParam["RegionId"] = "cn-hangzhou"
	requestParam["SignatureNonce"] = library.RandSeq(14)
	requestParam["SignatureVersion"] = "1.0"
	requestParam["SignatureMethod"] = "HMAC-SHA1"
	requestParam["AccessKeyId"] = accessKeyId
	requestParam["Timestamp"] = time.Now().UTC().String()
	return requestParam
}
func AliyunSigniture(requestVars map[string]string) string {

	var keys []string
	for key, _ := range requestVars{
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var temp []string

	for _, keyValue := range keys {
		temp = append(temp, url.QueryEscape(keyValue) + "=" + url.QueryEscape(requestVars[keyValue]))
	}

	implodedStr := strings.Join(temp, "&")
	StringToSign := "POST&" + url.QueryEscape("/") + "&" + implodedStr
	StringToSign = strings.Replace(StringToSign, "+", "%20", -1)
	StringToSign = strings.Replace(StringToSign, "*", "%2A", -1)
	StringToSign = strings.Replace(StringToSign, "%7E", "~", -1)
	logs.Warn(StringToSign)
	encodeKey := []byte(accessKeySecret+"&")
	mac := hmac.New(sha1.New, encodeKey)
	mac.Write([]byte(StringToSign))
	encodeRes := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeRes)
}

func SingleSendMail() string {
	var requestVar = AliyunApiCommon()
	requestVar["Action"] = "SingleSendMail"
	requestVar["AccountName"] = "service@mail.onestory.cn"
	requestVar["ReplyToAddress"] = "true"
	requestVar["Subject"] = "hahahah"
	requestVar["AddressType"] = "1"
	requestVar["ToAddress"] = "e930300047@163.com"
	requestVar["TextBody"] = "text"
	sign := AliyunSigniture(requestVar)
	requestVar["Signature"] = sign
	return request.HttpPost("https://dm.aliyuncs.com/", requestVar)
}


func SendToMail(to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}