package main

import (
	_ "myBeegoWeb/routers"

	"github.com/astaxie/beego"
)

//函数入口
//beego.Run()将监听app.conf所设置的httpport = 8080端口
func main() {
	beego.Run()
}
