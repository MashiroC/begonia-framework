package main

import "begonia/begonia/framework/application"

func main(){
	app:= application.Init(":1234")
	app.Get("/test",func(ctx *application.Context){
		ctx.ResponseString("hello world")
	})

	app.Start(1234)
}