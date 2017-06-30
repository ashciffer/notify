package module

import (
	"bytes"
	"net/http"
	"time"

	"errors"
	"fmt"

	"git.ishopex.cn/matrix/lazy-notify/lib"
)

type HttpTask struct {
	Id       string            //任务名称
	Input    map[string]string // http参数,url,params
	Level    int
	Ticker   *time.Timer
	SetLevel func(id string)
	GetLevel func(id string) int //获取任务等级方法
	Final    func(id string)     //任务终结
}

func DefaultHttpTask() *HttpTask {
	ht := new(HttpTask)
	ht.Id = lib.UniqueId()
	ht.Level = 0
	ht.Ticker = &time.Timer{}
	return ht
}

func (ht *HttpTask) Run() {
	ht.Level = ht.GetLevel(ht.Id)
	ht.Ticker = time.NewTimer(DurationForLevel(ht.Level))
	for {
		select {
		case <-ht.Ticker.C: //
			err := ht.Do()
			if err != nil { //任务执行失败 ，level++
				ht.SetLevel(ht.Id)
				ht.Level++
				if ht.Level >= 8 { //任务执行完成
					ht.Final(ht.Id)
					return
				}
			} else {
				ht.Final(ht.Id)
				return
			}
			ht.Ticker.Reset(DurationForLevel(ht.Level))
		}
	}
}

//执行任务(post请求)
func (ht *HttpTask) Do() error {
	url := ht.Input["URL"]
	params := ht.Input["params"]
	if url == "" || params == "" {
		return errors.New(fmt.Sprintf("http task execute failed with missing required parameter ,url:%s,params:%s", url, params))
	}

	client := &http.Client{
		Transport: http.DefaultTransport,
	}

	url_param_reader := bytes.NewReader([]byte(params))
	req, _ := http.NewRequest(
		"POST",
		url,
		url_param_reader,
	)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	T.Info("req :%+v", req.URL)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("response status code isn't 200 ,status code :%d", resp.StatusCode))
	}
	return nil
}
