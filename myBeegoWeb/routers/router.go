package routers

//routers这里处理路由选择
//beego.Router, 这个函数的功能是映射 URL 到 controller
//第一个参数是 URL (用户请求的地址)

import (
	"myBeegoWeb/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//兜底路由"/"，出于对beego框架开发者的敬意我设置为beego页面，会自动跳转到login页面
	beego.Router("/", &controllers.MainController{})
	//login是注册页面，由MyController的Get方法和Post方法处理
	beego.Router("/login", &controllers.MyController{})
	//manager是管理员页面，注意他的POST方法被我重写用来写留言处理了
	beego.Router("/manager", &controllers.ManagerController{})
	//cl_regist是客户注册页面
	beego.Router("/cl_regist", &controllers.ClRegistController{})
	//dr_regist是司机注册页面
	beego.Router("/dr_regist", &controllers.DrRegistController{})
	//skip是自动跳转页面，用于处理信息的提醒，注意我用到的Get:SkipGet的这些
	//这将skip对Get()方法的寻找是转换为找SkipGet()
	//本质上就是函数重命名，但一旦这样做就不会去找原本函数了(原本函数就是我前面/login写的Get和Post方法)
	beego.Router("/skip", &controllers.SkipController{}, "Get:SkipGet;Post:PasswordChangePost")
	//client是用户界面，ClientPost处理下单
	beego.Router("/client", &controllers.MyController{}, "Get:ClientGet;Post:ClientPost")
	//driver是司机页面，DriverPost处理接单
	beego.Router("/driver", &controllers.MyController{}, "Get:DriverGet;Post:DriverPost")
}
