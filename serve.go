package main

import (
	"github.com/NSMichelJ/go-todo-list/config"
	"github.com/NSMichelJ/go-todo-list/app"
)

func main(){
	config := config.GetDBConfig()
	app := app.App{}
	app.InitDB(config)
	app.Run(":8080")
}
