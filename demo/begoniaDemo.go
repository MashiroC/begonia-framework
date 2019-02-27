package main

import (
	"begonia"
)

func main1() {
	app := begonia.Init()
	app.Get("/hello", func(c *begonia.Context) {
		c.String("hello "+c.Param["name"])
	})
	app.Post("/welcome", func(c *begonia.Context) {
		c.String("welcome!!! "+c.Param["name"])
	})
	app.Start(1234)
}


