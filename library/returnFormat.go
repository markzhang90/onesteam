package library

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"reflect"
)

type (
	returnFormat struct {
		ErrNo  int
		ErrMsg string
		Data   interface{}
	}
)
/**
return format
 */
func ReturnJsonWithError(errNo int, errMsg string, data interface{}) (res string, err error) {

	if data == nil || !reflect.ValueOf(data).IsValid(){
		data = ""
	}
	if errMsg == "ref" {
		errMsg = CodeString(errNo)
	}

	formatter := new(returnFormat)
	formatter.ErrNo = errNo
	formatter.ErrMsg = errMsg
	formatter.Data = data

	result, err := json.Marshal(formatter)

	if err != nil {
		logs.Warn(err)
		return "", err
	}

	return string(result), nil
}
