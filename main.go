package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ryouaki/koa"
	"time"
)

var rdb *redis.Client = nil

var (
	_1kb   string = ""
	_10kb  string = ""
	_100kb string = ""
	_1mb   string = ""
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "42.192.194.38:6000",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := rdb.Ping().Result()
	fmt.Println(pong, err)

	var i = 0
	for ; i < 1000; i++ {
		_1kb += "A"
	}
	i = 0
	for ; i < 10; i++ {
		_10kb += _1kb
	}
	i = 0
	for ; i < 10; i++ {
		_100kb += _10kb
	}
	i = 0
	for ; i < 10; i++ {
		_1mb += _100kb
	}
}

func main() {
	app := koa.New() // 初始化服务对象

	// 设置api路由，其中var为url传参
	app.Get("/set/:count", func(err error, ctx *koa.Context, next koa.NextCb) {
		params := ctx.Params
		var v = ""
		switch params["count"] {
		case "1":
			v = _1kb
		case "10":
			v = _10kb
		case "100":
			v = _100kb
		case "1000":
			v = _1mb
		}
		e := rdb.Set(params["count"], v, 1000*time.Second).Err()
		if e != nil {
			ctx.Status = 500
			ctx.Write([]byte("ERROR"))
		} else {
			ctx.Write([]byte("OK"))
		}
	})
	app.Get("/get/:count", func(err error, ctx *koa.Context, next koa.NextCb) {
		params := ctx.Params
		_, e := rdb.Get(params["count"]).Result()
		if e != nil {
			ctx.Status = 500
			ctx.Write([]byte("ERROR"))
		} else {
			ctx.Write([]byte("OK"))
		}
	})

	err := app.Run(3000) // 启动
	if err != nil {      // 是否发生错误
		fmt.Println(err)
	}
}
