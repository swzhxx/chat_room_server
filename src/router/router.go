package router

import (
	"chat_room_server/src/express"
	"chat_room_server/src/model"
	"chat_room_server/src/utils/transfer"
	"fmt"
	"net"
)

var routerCollection = make([]map[string]func(*express.Context, *net.Conn, string, func()), 0)

type ConnectUser struct {
	User model.User
	Conn *net.Conn
}

var (
	ConnUserArray = make([]ConnectUser, 0)
)

//路由收集器，通过该方法收集路由,等待express初始化
func routerCollector(path string, callback func(_ *express.Context, _ *net.Conn, _ string, _ func())) {
	fmt.Println("router collector : ", path)
	m := map[string]func(*express.Context, *net.Conn, string, func()){}
	m[path] = callback
	slice := append(routerCollection, m)
	routerCollection = slice
}

//将通过routerCollector收集到的route 传入express
func RouterInit(App express.App) {
	fmt.Println("---------------routerinit--------------- routerCollection:", len(routerCollection))
	var globalApp express.App = App

	for _, v := range routerCollection {
		for path, callback := range v {
			globalApp.Use(path, callback)
		}
	}

	App.RegisterDisconnection(func(conn *net.Conn) {

		fmt.Println("enter disconnection handler ----")
		_arr := make([]ConnectUser, 0)

		var disConnUser ConnectUser

		fmt.Println("conn --->", conn)

		for _, v := range ConnUserArray {
			fmt.Println("v.conn --->", v.Conn)
			if conn != v.Conn {
				continue
			}
			disConnUser = v
			break
		}
		if disConnUser.User.Username == "" {
			return
		}

		for _, v := range ConnUserArray {
			if conn != v.Conn {
				t := transfer.Transfer{
					Conn: v.Conn,
				}
				t.Write("user/lists/deluser", disConnUser.User)
				_arr = append(_arr, v)
			}
		}
	})
}
