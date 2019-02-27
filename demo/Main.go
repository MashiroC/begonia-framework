package main

import (
	"fmt"
)

//func main() {
//	app := begonia.Init()
//	app.Get("/getQuestion", func(ctx *begonia.Context) {
//		ctx.W.Header().Add("Access-Control-Allow-Origin", "*")
//		time.Sleep(time.Duration(3)*time.Second)
//		res:=""
//		for i:=0;i<500;i++ {
//			res+="css是魔鬼吧\n"
//		}
//		ctx.String(res)
//	})
//	app.Start(1234)
//}

func main() {
	n := Node{nodeType: Normal}
	n.AddChild("/test", nil)
	n.AddChild("/teet", nil)
	n.AddChild("/teat", nil)
	n.AddChild("/tezx",nil)
	n.AddChild("/tesa",nil)
	n.AddChild("/tesz",nil)
	n.AddChild("/testtest",nil)
	fmt.Println(n.Path)
	fmt.Println(len(n.children))
	for i := 0; i < len(n.children[0].children); i++ {
		fmt.Println("child",n.children[0].children[i].Path)
	}
}


