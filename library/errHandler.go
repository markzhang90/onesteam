package library

import (
	"github.com/astaxie/beego/logs"
)

func RecoverPanic(){

	//捕获 and 恢复 panic.
	/*注意defer延迟执行的函数可以修改外围函数“testPanic2Error”的命名返回值.
	 *通过调用recover捕获panic并转化为error. 也许有人打算在main函数中放一个下面这个defer语句,用于捕获程序中的一切panic异常,
	 *建立最后一道防火墙,从而使程序避免崩溃运行下去, 但很不幸, 当main函数的defer延迟函数被执行时,也就意味着main函数要退出了,
	 *此时再捕获panic恢复程序,意义还有多大呢;不过我们有办法克服, 建立一个像“testPanic2Error”这样的一个外围封装函数,在这个函数中,
	 *建立最后一道防火墙,就像此处例子代码中所做的一样, 将panic封闭在自己的包内, 不允许蔓延传染给其它包, 包与包之间只通过error传递
	 *结果状态.
	*/

	defer func(){ // 必须要先声明defer，否则不能捕获到panic异常
		logs.Warning("c")

		if err:=recover();err!=nil{
			logs.Warning(err)
		}
		logs.Warning("c")
	}()
	//defer func() {
	//	logs.Warn("in defer")
	//	r := recover()
	//	logs.Warn(r)
	//	if r != nil {
	//		//check exactly what the panic was and create error.
	//		switch x := r.(type) {
	//		case string:
	//			err = errors.New(x)
	//		case error:
	//			err = x
	//		default:
	//			err = errors.New("Unknow panic")
	//		}
	//		logs.Warn(123)
	//
	//		logs.Warn(err)
	//		logs.Warn(456)
	//
	//	}
	//
	//}()

}
