# Begonia Framework

> 这个框架是一个Golang语言的简易框架，初衷是原生的net/http包写起来并不爽，想和springboot一样来写接口，现在的完成度还比较低

一个简易的Golang Web框架。

##  使用指南

### Prerequisites 项目使用条件

基于Go 1.11.2

无第三方包

### Installation 安装

先拉取这个包

```sh
mashiroc@ubuntu:$ go get github.com/MashiroC/begonia
```

然后import

```go
import "begonia/framework/application"
```



### Usage example 使用示例

```
package main

import "begonia/framework/application"

func main(){
	app:= application.Init()
	app.Get("/test",func(ctx *application.Context){
		ctx.ResponseString("hello world")
	})

	h:=application.Handle{Uri:"/welcome",Method:"GET",Fun:func(ctx *application.Context){
		ctx.ResponseString("welcome")
	}}

	app.AddHandle(h)

	app.Start(1234)
}
```



## 部署方法

```
go build xxxx
```

上传至服务器并运行

## Release History 版本历史

* 0.1.0

    上传初始代码

## Authors 关于作者

null

电脑没电了