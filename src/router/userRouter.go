package router

import (
	"chat_room_server/src/express"
	"chat_room_server/src/model"
	"chat_room_server/src/redis"
	"chat_room_server/src/utils/transfer"
	"encoding/json"
	"fmt"
	"net"
)

func init() {
	//用户添加路由
	routerCollector("user/add", func(context *express.Context, conn *net.Conn, data string, next func()) {
		fmt.Println("user/add ---> : ", data)
		dbConn := db.GetRedisConn()
		defer dbConn.Close()

		mes := express.Message{}
		_ = json.Unmarshal([]byte(data), &mes)
		var req model.User
		err := json.Unmarshal([]byte(mes.Data), &req)
		if err != nil {
			fmt.Println("user/add unmarshal err -->", err)
			return
		}
		t := &transfer.Transfer{
			Conn: conn,
		}
		var user model.User = model.User{}
		user, err = user.GetUser(dbConn, req.Username)
		if err != nil {
			t.WriteError("user/add", 500, "the users already exist")
			return
		}

		next()
	})

	//用户登录
	routerCollector("user/login", func(context *express.Context, conn *net.Conn, data string, next func()) {
		fmt.Println("enter user/login data -->:", data)
		dbConn := db.GetRedisConn()
		defer dbConn.Close()

		mes := express.Message{}
		_ = json.Unmarshal([]byte(data), &mes)

		var req model.User
		err := json.Unmarshal([]byte(mes.Data), &req)
		t := &transfer.Transfer{
			Conn: conn,
		}
		if err != nil {
			fmt.Println("user/login unmarshal err -->", err)
			return
		}
		var user model.User = model.User{}
		fmt.Println("req.Username:", req.Username)
		user, err = user.GetUser(dbConn, req.Username)
		if err != nil {
			fmt.Println("user/login err --->", err)
			return
		}
		if req.Userpwd != user.Userpwd {
			fmt.Println("user/login pwd is err")
			return
		}
		(*context)["user"] = user
		t.Write("user/login", user)
		next()
		return
	})

	//用户登录后保存有效用户链接和用户信息,并广播给客户端
	routerCollector("user/login", func(context *express.Context, conn *net.Conn, data string, next func()) {
		user, _ := (*context)["user"].(model.User)
		cu := ConnectUser{
			User: user,
			Conn: conn,
		}

		var usersArray []model.User = make([]model.User, 0)

		for _, v := range ConnUserArray {
			usersArray = append(usersArray, v.User)
		}
		for _, v := range ConnUserArray {
			t := &transfer.Transfer{
				Conn: v.Conn,
			}
			t.Write("user/lists/adduser", user)
		}
		ConnUserArray = append(ConnUserArray, cu)
		t := &transfer.Transfer{
			Conn: conn,
		}
		t.Write("user/lists", usersArray)
		next()
		return
	})
}
