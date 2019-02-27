package controllers

//controll是控制器，处理各种请求
import (
	"fmt"
	"myBeegoWeb/models"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//MainController golang的类是通过Struct来实现的
//里面包含了beego.Controller这个结构体
//其实也就是继承了beego.Controller的所有的方法和变量
type MainController struct {
	beego.Controller
}

//MyController 用于login和用户界面
type MyController struct {
	beego.Controller
}

//ManagerController 用于管理员界面
type ManagerController struct {
	beego.Controller
}

//ClRegistController 用于乘客注册页面
type ClRegistController struct {
	beego.Controller
}

//DrRegistController 用于司机注册页面
type DrRegistController struct {
	beego.Controller
}

//SkipController 用于自动跳转页面
type SkipController struct {
	beego.Controller
}

//Get 获取beego的初始页面
func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

//Get login页面获取
func (c *MyController) Get() {
	c.TplName = "login.html"
}

//Post login页面的表单处理
func (c *MyController) Post() {
	//GetString(name) 可获得表单内相应name的框内信息
	userName := c.GetString("user_name")
	password := c.GetString("password")
	if userName == "" || password == "" {
		//控制台上报错
		beego.Error("用户名或密码获取失败！")
		//Session是保存在服务器上的信息
		//cookies是保存在用户浏览器的信息
		//我倾向于把日志信息保存在服务器，其他全扔给浏览器
		c.SetSession("/skipMessage", "用户名或密码获取失败!")
		c.Ctx.SetCookie("skipHtml", "/login")
		c.Redirect("/skip", 302)
		return
	}
	//判断是否为管理员，我不想有多个管理员，所以干脆固定下来
	if userName == "root" && password == "123" {
		c.Data["cal_name"] = "你好啊！管理员"
		c.Redirect("/manager", 302)
		return
	}
	//user是一个暂存信息的结构体
	user := models.Client{}
	user.ClName = userName
	//o为Orm对象，用于数据库操作
	o := orm.NewOrm()
	//进行查找，默认为靠主键查找
	err1 := o.Read(&user)
	if err1 != nil {
		//不为乘客再看是否为司机
		driver := models.Driver{}
		driver.DrName = userName
		err2 := o.Read(&driver)
		if err2 != nil {
			beego.Error("不存在此用户！")
			c.SetSession("/skipMessage", "不存在此用户")
			c.Ctx.SetCookie("skipHtml", "/login")
			c.Redirect("/skip", 302)
			return
		}
		if driver.DrPassword != password {
			beego.Error("密码错误！")
			c.SetSession("/skipMessage", "密码错误")
			c.Ctx.SetCookie("skipHtml", "/login")
			c.Redirect("/skip", 302)
			return
		}
		c.Ctx.SetCookie("dr_name", userName)
		c.SetSession("/skipMessage", "成功登入")
		c.Ctx.SetCookie("skipHtml", "/driver")
		c.Redirect("/skip", 302)
		return
	}
	if user.ClPassword != password {
		beego.Error("密码错误！")
		c.SetSession("/skipMessage", "密码错误")
		c.Ctx.SetCookie("skipHtml", "/login")
		c.Redirect("/skip", 302)
		return
	}
	c.Ctx.SetCookie("cl_name", userName)
	c.SetSession("/skipMessage", "成功登入")
	c.Ctx.SetCookie("skipHtml", "/client")
	c.Redirect("/skip", 302)
}

//Get 管理员页面获取
func (c *ManagerController) Get() {
	o := orm.NewOrm()
	//查询Client表
	var clients []models.Client
	_, err := o.QueryTable("client").All(&clients)
	if err != nil {
		beego.Error("查询所有乘客信息出错")
		return
	}
	c.Data["Client"] = clients
	//查询Driver表
	var drivers []models.Driver
	_, err = o.QueryTable("driver").All(&drivers)
	if err != nil {
		beego.Error("查询所有司机信息出错")
		return
	}
	c.Data["Driver"] = drivers
	//查询未完成订单，由于这类订单的cost默认为0，所以以此查询
	var waitOrder []models.Order
	qs := o.QueryTable("order")
	_, err = qs.Filter("cost", 0).All(&waitOrder)
	if err != nil {
		beego.Error("查询所有未接订单信息出错")
		return
	}
	c.Data["Order1"] = waitOrder
	//查询已完成订单
	var order []models.Order
	_, err = o.Raw("select * from `order` where cost!=0").QueryRows(&order)
	if err != nil {
		beego.Error("查询所有完成单信息出错")
		return
	}
	c.Data["Order2"] = order
	//查询留言信息
	var message []models.Message
	_, err = o.QueryTable("message").All(&message)
	if err != nil {
		beego.Error("查询所有留言信息出错")
		return
	}
	c.Data["Message"] = message
	beego.Info("管理员登入成功")
	c.TplName = "manager.html"
}

//Post 进行给管理员留言
func (c *ManagerController) Post() {
	//获取留言信息和留言者信息
	message := c.GetString("messageText")
	var name string
	name = c.Ctx.GetCookie("cl_name")
	var node models.Message
	o := orm.NewOrm()
	if name == "" {
		name = c.Ctx.GetCookie("dr_name")
		if name == "" {
			beego.Error("获取留言者的信息失败！")
			c.SetSession("/skipMessage", "失去了登入状态,请重新登入")
			c.Ctx.SetCookie("skipHtml", "/login")
			c.Redirect("/skip", 302)
			return
		}
		node.WriteBy = name
		node.Node = message
		_, err := o.Insert(&node)
		if err != nil {
			beego.Error("留言数据库操作失败！")
			c.SetSession("/skipMessage", "留言失败！")
			c.Ctx.SetCookie("skipHtml", "/driver")
			c.Redirect("/skip", 302)
			return
		}
		beego.Info(name, "留言成功")
		c.SetSession("/skipMessage", "留言成功")
		c.Ctx.SetCookie("skipHtml", "driver")
		c.Redirect("/skip", 302)
	}
	node.Node = message
	node.WriteBy = name
	//进行数据库插入操作
	_, err := o.Insert(&node)
	if err != nil {
		beego.Error("留言数据库操作失败！")
		c.SetSession("/skipMessage", "留言失败！")
		c.Ctx.SetCookie("skipHtml", "/client")
		c.Redirect("/skip", 302)
		return
	}
	beego.Info(name, "留言成功")
	c.SetSession("/skipMessage", "留言成功")
	c.Ctx.SetCookie("skipHtml", "client")
	c.Redirect("/skip", 302)
}

//Get 乘客注册页面获取
func (c *ClRegistController) Get() {
	c.TplName = "cl_regist.html"
}

//Post 乘客注册表单处理
func (c *ClRegistController) Post() {
	cl_name := c.GetString("Cl_Name")
	cl_password := c.GetString("Cl_password")
	if cl_password != c.GetString("Cl_password2") {
		beego.Error("两次密码输入不相同！")
		c.SetSession("/skipMessage", "两次密码输入不相同！")
		c.Ctx.SetCookie("skipHtml", "/cl_regist")
		c.Redirect("/skip", 302)
		return
	}
	cl_phone := c.GetString("Cl_Phone")
	cl_id := c.GetString("User_id")
	//3.插入数据库
	o := orm.NewOrm()
	client := models.Client{}
	client.ClName = cl_name
	client.ClPassword = cl_password
	client.ClPhone = cl_phone
	client.ClIdentificationCard = cl_id
	_, err := o.Insert(&client)
	if err != nil {
		c.Ctx.WriteString("")
		beego.Error("注册失败！数据库操作错误")
		c.SetSession("/skipMessage", "注册失败！输入不合法")
		c.Ctx.SetCookie("skipHtml", "/cl_regist")
		c.Redirect("/skip", 302)
		return
	}
	beego.Info(cl_name, " 注册成功")
	c.SetSession("/skipMessage", "注册成功")
	c.Ctx.SetCookie("skipHtml", "/login")
	c.Redirect("/skip", 302)
}

//Get 司机注册页面获取
func (c *DrRegistController) Get() {
	c.TplName = "dr_regist.html"
}

//Post 司机注册表单的处理
func (c *DrRegistController) Post() {
	dr_phone := c.GetString("DrPhone")
	dr_name := c.GetString("DrName")
	dr_password := c.GetString("Drpwd")
	if dr_password != c.GetString("Drpwd2") {
		beego.Error("两次密码输入不一致！")
		c.SetSession("/skipMessage", "两次密码输入不一致！")
		c.Ctx.SetCookie("skipHtml", "/dr_regist")
		c.Redirect("/skip", 302)
		return
	}
	dr_id := c.GetString("UserId")
	car_type := c.GetString("CarType")
	car_load, err2 := c.GetInt("CarLoad")
	if err2 != nil {
		beego.Error("获取信息失败", err2)
		c.SetSession("/skipMessage", "获取表内信息失败！")
		c.Ctx.SetCookie("skipHtml", "/dr_regist")
		c.Redirect("/skip", 302)
	}

	driver := models.Driver{
		DrName:               dr_name,
		DrPassword:           dr_password,
		DrPhone:              dr_phone,
		DrIdentificationCard: dr_id,
		CarBrand:             car_type,
		CarLoad:              car_load,
	}
	o := orm.NewOrm()
	_, err := o.Insert(&driver)
	if err != nil {
		beego.Error("注册司机失败！数据库操作错误", err)
		c.SetSession("/skipMessage", "注册失败，数据库操作错误！")
		c.Ctx.SetCookie("skipHtml", "/cl_regist")
		c.Redirect("/skip", 302)
		return
	}
	beego.Info(dr_name, "注册成功")
	c.SetSession("/skipMessage", "注册成功！")
	c.Ctx.SetCookie("skipHtml", "/login")
	c.Redirect("/skip", 302)
}

//ClientGet 乘客页面获取
func (c *MyController) ClientGet() {
	var username string
	username = c.Ctx.GetCookie("cl_name")
	c.Data["cl_name"] = username
	client := models.Client{}
	client.ClName = username
	o := orm.NewOrm()
	//获取用户信息
	err := o.Read(&client)
	if err != nil {
		beego.Error("用户信息获取失败", err)
		return
	}
	//获取已完成订单信息
	var order []models.Order
	_, err = o.Raw("select * from `order` where client_id = ? and cost != 0", username).QueryRows(&order)
	if err != nil {
		beego.Error("查询所有完成单信息出错")
		return
	}
	c.Data["Order2"] = order
	c.Data["ClName"] = client.ClName
	c.Data["ClPassword"] = client.ClPassword
	c.Data["ClPhone"] = client.ClPhone
	c.Data["ClIdentificationCard"] = client.ClIdentificationCard
	beego.Info(username, "登入成功")
	c.TplName = "client.html"
}

//ClientPost 下单打车操作
func (c *MyController) ClientPost() {
	username := c.Ctx.GetCookie("cl_name")
	client := models.Client{}
	client.ClName = username
	o := orm.NewOrm()
	//获取用户信息
	err := o.Read(&client)
	if err != nil {
		beego.Error("用户信息获取失败", err)
		c.SetSession("/skipMessage", "用户信息获取失败！")
		c.Ctx.SetCookie("skipHtml", "/client")
		c.Redirect("/skip", 302)
		return
	}
	//获取打车的信息
	orgin := c.GetString("orgin")
	destination := c.GetString("destination")
	node := c.GetString("nodeText")
	order := models.Order{
		Orgin:       orgin,
		Destination: destination,
		BeginTime:   time.Now(),
		Client:      &client,
		Node:        node,
	}
	//进行插入操作
	_, err = o.Insert(&order)
	if err != nil {
		beego.Error("订单插入数据库失败", err)
		c.SetSession("/skipMessage", "订单插入数据库失败！")
		c.Ctx.SetCookie("skipHtml", "/client")
		c.Redirect("/skip", 302)
		return
	}
	beego.Info("下单成功！")
	c.SetSession("/skipMessage", "下单成功！")
	c.Ctx.SetCookie("skipHtml", "/client")
	c.Redirect("/skip", 302)
}

//DriverGet 司机页面获取
func (c *MyController) DriverGet() {
	var drName string
	drName = c.Ctx.GetCookie("dr_name")
	//beego.Info(c.Ctx.GetCookie("dr_name"))
	c.Data["dr_name"] = drName
	driver := models.Driver{}
	driver.DrName = drName
	o := orm.NewOrm()
	//获取司机信息
	err := o.Read(&driver)
	if err != nil {
		beego.Error("司机信息获取失败", err)
		return
	}
	c.Data["DrName"] = driver.DrName
	c.Data["DrPassword"] = driver.DrPassword
	c.Data["DrPhone"] = driver.DrPhone
	c.Data["DrIdentificationCard"] = driver.DrIdentificationCard
	c.Data["CarBrand"] = driver.CarBrand
	c.Data["CarLoad"] = driver.CarLoad
	//获取待接订单信息
	qs := o.QueryTable("order")
	var wait []models.Order
	qs2 := qs.Filter("cost", 0)
	_, err = qs2.All(&wait)
	if err != nil {
		beego.Error("获取未接订单错误", err)
		return
	}
	c.Data["Order"] = wait
	for i, v := range wait {
		cookieName := fmt.Sprintf("client%d", i)
		c.Ctx.SetCookie(cookieName, v.Client.ClName)
	}
	var order []models.Order
	//查询历史订单记录
	_, err = o.Raw("select * from `order` where driver_id = ?", drName).QueryRows(&order)
	if err != nil {
		beego.Info("查询所有完成单信息出错")
		return
	}
	c.Data["Order2"] = order
	beego.Info(drName, "登入成功")
	c.TplName = "driver.html"
}

//DriverPost 司机接单操作
func (c *MyController) DriverPost() {
	//获取司机信息
	dr_name := c.Ctx.GetCookie("dr_name")
	if dr_name == "" {
		beego.Error("司机信息获取失败！")
		c.SetSession("/skipMessage", "司机信息获取失败！")
		c.Ctx.SetCookie("skipHtml", "/driver")
		c.Redirect("/skip", 302)
		return
	}
	cost, err := c.GetFloat("cost", 10.01)
	if err != nil {
		beego.Error("获取价格出错！", err)
		c.SetSession("/skipMessage", "获取价格出错！")
		c.Ctx.SetCookie("skipHtml", "/driver")
		c.Redirect("/skip", 302)
		return
	}
	driver := models.Driver{}
	driver.DrName = dr_name
	o := orm.NewOrm()
	err = o.Read(&driver)
	if err != nil {
		beego.Error("从数据库获取司机信息出错！", err)
		c.SetSession("/skipMessage", "从数据库获取司机信息出错！")
		c.Ctx.SetCookie("skipHtml", "/driver")
		c.Redirect("/skip", 302)
		return
	}
	//获取订单信息
	orderId := c.GetString("OrId")
	if orderId == "" {
		beego.Error("订单信息获取失败！")
		c.SetSession("/skipMessage", "订单信息获取失败！")
		c.Ctx.SetCookie("skipHtml", "/driver")
		c.Redirect("/skip", 302)
		return
	}
	order := models.Order{}
	order.OrId, _ = strconv.ParseInt(orderId, 10, 64)
	err = o.Read(&order)
	if err != nil {
		beego.Error("从数据库获取订单失败！")
		c.SetSession("/skipMessage", "从数据库获取订单失败！")
		c.SetSession("skipHtml", "/driver")
		c.Redirect("/skip", 302)
		return
	}
	order.Driver = &driver
	order.Cost = cost
	//更新订单，状态为已完成
	_, err = o.Update(&order)
	if err != nil {
		beego.Error("更新数据失败！")
		c.SetSession("/skipMessage", "更新数据失败！")
		c.SetSession("skipHtml", "/driver")
		c.Redirect("/skip", 302)
		return
	}
	beego.Info("接单成功！")
	c.SetSession("/skipMessage", "接单成功！")
	c.SetSession("skipHtml", "/driver")
	c.Redirect("/skip", 302)
}

//SkipGet 信息提醒跳转页面
func (c *SkipController) SkipGet() {
	skipMessage := c.GetSession("/skipMessage")
	skipHtml := c.Ctx.GetCookie("skipHtml")
	c.Data["skipMessage"] = skipMessage
	c.Data["skipHtml"] = skipHtml
	c.TplName = "skip.html"
}

//PasswordChangePost 进行密码修改操作
func (c *SkipController) PasswordChangePost() {
	//获取表单信息
	pwd := c.GetString("pwd")
	newPwd := c.GetString("newPwd")
	newPwd2 := c.GetString("newPwd2")
	var user_name string
	//获取当前用户
	user_name = c.Ctx.GetCookie("cl_name")
	beego.Info(user_name)
	o := orm.NewOrm()
	if user_name == "" {
		user_name = c.Ctx.GetCookie("dr_name")
		beego.Info(user_name)
		if user_name == "" {
			beego.Error("获取当前用户出错！")
			c.SetSession("/skipMessage", "账户已退出")
			c.Ctx.SetCookie("skip", "/driver")
			c.Redirect("/skip", 302)
			return
		}
		if pwd == "" || newPwd == "" || newPwd2 == "" {
			beego.Error("获取更改密码表格信息出错！")
			c.SetSession("/skipMessage", "获取更改密码表格信息出错！")
			c.Ctx.SetCookie("skip", "/driver")
			c.Redirect("/skip", 302)
			return
		}
		driver := models.Driver{}
		driver.DrName = user_name
		o.Read(&driver)
		if driver.DrPassword != pwd {
			beego.Error("原密码输入错误！")
			c.SetSession("/skipMessage", "原密码输入错误")
			c.Ctx.SetCookie("skipHtml", "/driver")
			c.Redirect("/skip", 302)
			return
		}
		if newPwd != newPwd2 {
			beego.Error("两次新密码输入不一致！")
			c.SetSession("/skipMessage", "两次新密码输入不一致")
			c.Ctx.SetCookie("skipHtml", "/driver")
			c.Redirect("/skip", 302)
			return
		}
		driver.DrPassword = newPwd
		_, err := o.Update(&driver)
		if err != nil {
			beego.Error(user_name, "数据库更新密码失败！")
			c.SetSession("/skipMessage", "数据库更新密码失败！")
			c.Ctx.SetCookie("skipHtml", "/driver")
			c.Redirect("/skip", 302)
			return
		}
		beego.Info(user_name, "更新密码成功！")
		c.SetSession("/skipMessage", "更新密码成功")
		c.Ctx.SetCookie("skipHtml", "/driver")
		c.Redirect("/skip", 302)
		return
	} else {
		client := models.Client{}
		client.ClName = user_name
		o.Read(&client)
		if client.ClPassword != pwd {
			beego.Error("原密码输入错误！")
			c.SetSession("/skipMessage", "原密码输入错误")
			c.Ctx.SetCookie("skipHtml", "/client")
			c.Redirect("/skip", 302)
			return
		}
		if newPwd != newPwd2 {
			beego.Error("两次新密码输入不一致！")
			c.SetSession("/skipMessage", "两次新密码输入不一致")
			c.Ctx.SetCookie("skipHtml", "/client")
			c.Redirect("/skip", 302)
			return
		}
		client.ClPassword = newPwd
		_, err := o.Update(&client)
		if err != nil {
			beego.Error(user_name, "数据库更新密码失败！")
			c.SetSession("/skipMessage", "数据库更新密码失败！")
			c.Ctx.SetCookie("skipHtml", "/client")
			c.Redirect("/skip", 302)
			return
		}
		beego.Info(user_name, "更新密码成功！")
		c.SetSession("/skipMessage", "更新密码成功")
		c.Ctx.SetCookie("skipHtml", "/client")
		c.Redirect("/skip", 302)
		return
	}
}
