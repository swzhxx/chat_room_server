package express

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type App interface {
	Use(string, func(*Context, *net.Conn, string, func()))
	Listen(string, string, func()) (err error)
	RegisterMessageValiator(func(string) bool)
	RegisterDisconnection(func(*net.Conn))
}

//上下文
type Context map[string]interface{}

type app struct {
	context               Context
	middlewares           []middleware
	validator             func(string) bool
	disconnectioncallback func(*net.Conn)
}

//中间件
type middleware struct {
	path     string
	callback func(*Context, *net.Conn, string, func())
}

type Message struct {
	Path string `json:"path"`
	Data string `json:"data"`
	Code int    `json:"code"`
}

//注册消息验证函数
func (this *app) RegisterMessageValiator(validator func(string) bool) {
	this.validator = validator
	return
}

//断线函数
func (this *app) RegisterDisconnection(callback func(*net.Conn)) {
	this.disconnectioncallback = callback
	return
}

//开启服务监听
//address 监听ip地址
//port 监听端口
//callback 监听成功后返回的回调函数
func (this *app) Listen(address string, port string, callback func()) (err error) {
	listen, err := net.Listen("tcp", address+":"+port)
	if err != nil {

		fmt.Println("app listen error : ", err)
		return
	}
	defer listen.Close()
	callback()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("lisen Accept error : ", err)
			return err
		}
		fmt.Printf("Accept client ip=%v \n", conn.RemoteAddr().String())

		go this.processor(&conn)
	}
}

//创建中间件方法
// params : path? , callback
//path 不一定传入 ,callback 必须传入
func (this *app) Use(path string, callback func(context *Context, conn *net.Conn, data string, _ func())) {
	var _middleware middleware
	_middleware = middleware{
		path:     path,
		callback: callback,
	}
	this.middlewares = append(this.middlewares, _middleware)
	return
}

//连接的实际执行函数
//在这里进行消息的分发,分给对应的消息的中间件
//
func (this *app) processor(conn *net.Conn) (err error) {
	if conn == nil {
		info := "express processor error :conn is nil"
		err = errors.New(info)
		fmt.Println(info)
		return
	}
	defer func() {
		(*conn).Close()
		if this.disconnectioncallback == nil {
			return
		}
		this.disconnectioncallback(conn)
	}()
	for {
		buf := make([]byte, 1024*4)
		n, err2 := (*conn).Read(buf)
		fmt.Println("n--->", n)
		if err2 != nil {
			fmt.Println("conn Read err : ", err2)
			err = err2
			if n == 0 {
				return
			}
			return
		}
		validator := this.validator
		if validator == nil {
			validator = func(string) bool {
				return true
			}
		}
		str := string(buf[:n])
		fmt.Printf("str - > :%v \n", str)
		var validate = validator(str)
		if !validate {
			//验证不通过消息 直接丢弃
			fmt.Printf("validate fail ")
			continue
		}
		//将消息反系列化
		var mes Message
		err = json.Unmarshal(buf[:n], &mes)
		fmt.Printf("message - > :%v \n", mes)
		if err != nil {
			fmt.Println("json Unmarshal message  err", err)
			continue
		}
		//查找需要执行的中间件
		i := 0
		var next func()
		var context = Context{}
		next = func() {
			length := len(this.middlewares)
			if i >= length {
				return
			}
			ware := this.middlewares[i]
			i++

			if ware.path == "" || ware.path == mes.Path {
				ware.callback(&context, conn, str, next)
			} else {
				next()
			}
		}
		next()
		fmt.Println("middleware complete")
	}

	fmt.Println("some one leave")
	return
}

//创建一个app实例,
//返回一个带有监听（创建tcp）服务器的方法
//一个Use中间件方法
func NewApp() App {

	context := Context{}
	_app := &app{
		middlewares: make([]middleware, 0),
		context:     context,
	}
	var _interface App = _app
	return _interface
}
