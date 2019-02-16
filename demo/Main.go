package main

import (
	"begonia"
	"time"
)

func main() {
	app := begonia.Init()
	app.Get("/getQuestion", func(ctx *begonia.Context) {
		ctx.W.Header().Add("Access-Control-Allow-Origin", "*")
		time.Sleep(time.Duration(3)*time.Second)
		res:=""
		for i:=0;i<500;i++ {
			res+="css是魔鬼吧\n"
		}
		ctx.String(res)
	})
	app.Start(1234)
}
