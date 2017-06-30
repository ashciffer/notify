package module

import (
	"time"
)

//
type LazyTask interface {
	Run() error  //执行任务
	Stop() error //停止任务
	Next() int64 //返回任务下次执行的时间戳
}

//根据等级获取时间
func DurationForLevel(level int) time.Duration {
	var t time.Duration
	switch level {
	case 0:
		t = time.Duration(5 * time.Second)
	case 1:
		t = time.Duration(2 * time.Minute)
	case 2:
		t = time.Duration(10 * time.Minute)
	case 3:
		t = time.Duration(15 * time.Minute)
	case 4:
		t = time.Duration(1 * time.Hour)
	case 5:
		t = time.Duration(2 * time.Hour)
	case 6:
		t = time.Duration(6 * time.Hour)
	case 7:
		t = time.Duration(15 * time.Hour)
	default:
		t = time.Duration(5 * time.Second)
	}
	return t
}
