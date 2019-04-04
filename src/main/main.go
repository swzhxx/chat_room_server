package main

import (
	"chat_room_server/src/express"
	"chat_room_server/src/router"
	"fmt"
)

func main() {
	var app express.App = express.NewApp()
	//定义router
	router.RouterInit(app)
	app.Listen("127.0.0.1", "8899", func() {
		fmt.Println("server is start")
	})
}
