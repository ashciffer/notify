package main

import (
	"net/http"
	"time"

	"github.com/go-ini/ini"

	. "git.ishopex.cn/matrix/lazy-notify/module"
	"git.ishopex.cn/matrix/lazy-notify/schedule"
	"github.com/robfig/cron"
)

func main() {
	T.Info("%s Runing ~~~ ", time.Now().Format(FORMATTIME))
	cfg, err := ini.LooseLoad("conf/app.conf", "")
	if err != nil {
		T.Error("lazy-notfiy conf init failed:%s", err)
		return
	}

	mongourl := cfg.Section("mongo").Key("url").String()
	mongodb := cfg.Section("mongo").Key("db").String()
	mongocol := cfg.Section("mongo").Key("collection").String()

	redisurl := cfg.Section("redis").Key("url").String()
	redisdb, _ := cfg.Section("redis").Key("db").Int()

	var mo schedule.MongoOptions
	err = mo.Init(mongourl, mongodb, mongocol, redisurl, redisdb)
	if err != nil {
		T.Error("lazy-notfiy mongo option init failed:%s", err)
		return
	}
	cronapp := cron.New()
	cronapp.AddFunc("@every 10m", mo.CreateTask)
	cronapp.Start()
	err = WebServer("7788")
	if err != nil {
		T.Error("lazy-notfiy webserver init failed:%s", err)
		return
	}
}

//启动服务
func WebServer(port string) error {
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./"))))
	return http.ListenAndServe(":"+port, nil)
}
