package module

import "github.com/astaxie/beego/logs"

var T *logs.BeeLogger

func init() {
	T = logs.NewLogger(100)
	T.SetLogger("file", `{"filename":"logs/system.log","daily":true,"maxdays":7}`)
	T.SetLogger("console", `{"level":1}`)
	T.EnableFuncCallDepth(true)
	T.SetLogFuncCallDepth(2)
}
