package router

import (
	"chat_room_server/src/express"
	"chat_room_server/src/model"
	"chat_room_server/src/utils/transfer"
	"encoding/json"
	"fmt"
	"net"
)

type Sms struct {
	Source model.User `json:"source"`
	Data   string     `json:"data"`
}

func init() {
	routerCollector("sms/add", func(context *express.Context, conn *net.Conn, data string, next func()) {
		fmt.Println("--enter sms/add---")
		mes := express.Message{}
		_ = json.Unmarshal([]byte(data), &mes)
		sms := Sms{}
		fmt.Println("mes.Data --->", mes.Data)
		_ = json.Unmarshal([]byte(mes.Data), &sms)

		t := transfer.Transfer{
			Conn: conn,
		}
		//回复
		t.Write("sms/add", sms)
		//转发
		for _, v := range ConnUserArray {
			if v.User.Username == sms.Source.Username {
				continue
			}
			t := transfer.Transfer{
				Conn: v.Conn,
			}
			t.Write("sms/receive", sms)
		}
	})
}
