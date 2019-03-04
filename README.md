# Begonia Framework

A golang web framework for efficient and concise.

## User Guide

### Prerequisites

- Golang Version >= 1.11.2

### Installation

```sh
$ go get github.com/MashiroC/begonia
```

### Example

```golang
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



## Run

```
go build xxxx
./xxxx
```

## Release History

* 0.1.0

    - Proof-of-concept code