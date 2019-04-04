package model

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type User struct {
	Username     string `json:"username"`
	Userpwd      string `json:"userpwd"`
	Usernickname string `json:"usernickname"`
}

func (this *User) GetUser(dbConn redis.Conn, field string) (user User, err error) {
	res, err := redis.String(dbConn.Do("hget", "users", field))
	if err != nil {
		if err == redis.ErrNil {
			return
		}
	}
	user = User{}
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Println("GetUser json.unmarshal err =", err)
	}
	return
}

func (this *User) Insert(dbConn redis.Conn) (err error) {
	str, err := json.Marshal(*this)
	_, err = dbConn.Do("hset", "users", this.Username, str)
	return
}
