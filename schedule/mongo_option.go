package schedule

import (
	"runtime/debug"
	"strconv"

	mongo "git.ishopex.cn/matrix/kaola/db"
	"git.ishopex.cn/matrix/lazy-notify/lib"
	. "git.ishopex.cn/matrix/lazy-notify/module"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2/bson"
)

type MongoOptions struct {
	Mop       mongo.Mop
	DB        string //数据库
	Colletion string //集合
	rp        *redis.Pool
}

func ExitRecovery() {
	if re := recover(); re != nil {
		T.Error("runtime error:%s", re)
		T.Error("stack :%s", string(debug.Stack()))
	}

}

func (mo *MongoOptions) Init(mongourl, mongodb, mongocol, redisurl string, redisdb int) error {
	mo.DB = mongodb
	mo.Colletion = mongocol
	mo.Mop.Url = mongourl

	mo.rp = lib.NewRedisPoll(redisurl, redisdb)
	return nil
}

//遍历集合，查找status为0的response
func (mo *MongoOptions) Traverse() ([]map[string]string, error) {
	var resutlt []map[string]string
	var ret []map[string]string

	session, err := mo.Mop.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	col := session.DB(mo.DB).C(mo.Colletion)
	query := bson.M{"status": 0}
	err = col.Find(query).All(&resutlt)
	if err != nil {
		return nil, err
	}

	defer ExitRecovery()
	for _, v := range resutlt {
		temp := make(map[string]string)
		if id, ok := v["msg_id"]; ok {
			temp["URL"] = id
		}

		if u, ok := v["callback_url"]; ok {
			temp["URL"] = u
		}

		if p, ok := v["params"]; ok {
			temp["params"] = p
		}
		ret = append(ret, temp)
	}
	return ret, nil
}

//修改mongo中数据状态
func (mo *MongoOptions) ChangeStatus(id string) {
	change := bson.M{"$set": bson.M{"status": 1}}
	selector := bson.M{"msg_id": id}

	session, err := mo.Mop.GetSession()
	if err != nil {
		T.Error("mongooptions change status ,get session failed,id:%s,error:%s", id, err)
		return
	}
	defer session.Close()

	col := session.DB(mo.DB).C(mo.Colletion)
	err = col.Update(selector, change)
	if err != nil {
		T.Error("mongooptions change status ,update data failed,id:%s,error:%s", id, err)
		return
	}
}

//从redis中获取id
func (mo *MongoOptions) GetLevel(id string) int {
	defer ExitRecovery()
	conn := mo.rp.Get()
	T.Info("get level id ->%s", id)
	reply, err := conn.Do("GET", id)
	if err != nil {
		T.Error("mongooptions  getlevel,command do failed,id:%s,error:%s", id, err)
		return 0
	}

	if reply == nil {
		return 0
	}
	T.Error("reply -> %+v", string(reply.([]byte)))
	b, _ := strconv.Atoi(string(reply.([]byte)))
	return b
}

func (mo *MongoOptions) SetLevel(id string) {
	conn := mo.rp.Get()
	T.Info("set level id ->%s", id)
	_, err := conn.Do("INCRBY", id, 1)
	if err != nil {
		T.Error("mongooptions  setlevel,command do failed,id:%s,error:%s", id, err)
	}
}

//创建http任务
func (mo *MongoOptions) CreateTask() {
	travers, err := mo.Traverse()
	if err != nil {
		T.Error("mongooptions create task ,traverse failed,error:%s", err)
		return
	}

	for _, v := range travers {
		var t HttpTask
		t.Final = mo.ChangeStatus
		t.GetLevel = mo.GetLevel
		t.SetLevel = mo.SetLevel
		t.Id = v["id"]
		t.Input = v
		go t.Run()
	}
}
