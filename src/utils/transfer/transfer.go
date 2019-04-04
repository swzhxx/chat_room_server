package transfer

import (
	"chat_room_server/src/express"
	"encoding/json"
	"fmt"
	"net"
)

type Transfer struct {
	Conn *net.Conn
}

func (this *Transfer) Read() {

}

func (this *Transfer) Write(path string, data interface{}) {
	var m = make(map[string]interface{})
	m["path"] = path
	m["data"] = data
	bytes, err := json.Marshal(m)
	if err != nil {
		fmt.Println("transfer write error --->", err)
		return
	}

	endFlag, _ := json.Marshal("||||||")
	bytes = append(bytes, endFlag...)
	conn := *(this.Conn)
	conn.Write(bytes)

	return
}

func (this *Transfer) WriteError(path string, code int, err interface{}) {
	var data string
	switch err.(type) {
	case string:
		data = err.(string)
	case error:
		data = fmt.Sprintf("%s", err)
	}

	message := &express.Message{
		Path: path,
		Data: data,
		Code: code,
	}
	pkg, err := json.Marshal(message)
	if err != nil {
		fmt.Println("transfer WriteError -->", err)
		return
	}
	conn := *(this.Conn)
	conn.Write(pkg)
	return
}
