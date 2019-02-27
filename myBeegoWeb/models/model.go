package models

//model函数包涵对数据形式和数据库的初始化
//注意对表的操作是以相应的结构体为主
//我使用的数据库软件为MySQl，用户名root，密码123
//请在使用前先开一个testdidi数据库
//则本程序可以自动生成响应的数据表

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//Client 是乘客类
type Client struct {
	ClName               string   `orm:"pk;size(20);unique"` //客户用户名，主键
	ClPassword           string   `orm:"size(20)"`           //密码
	ClPhone              string   `orm:"size(11);unique"`    //电话
	ClIdentificationCard string   `orm:"size(18);unique"`    //身份证号
	ClOrder              []*Order `orm:"reverse(many)"`      //一个用户有多个完成的订单
}

//Driver 是司机类
type Driver struct {
	//DrId                 int64    `orm:"pk;auto"`         //司机编号
	DrName               string   `orm:"pk;size(20);unique"` //司机用户名，主键
	DrPassword           string   `orm:"size(20)"`           //密码
	DrPhone              string   `orm:"size(11);unique"`    //电话
	DrIdentificationCard string   `orm:"size(18);unique"`    //身份证
	CarBrand             string   `orm:"size(20)"`           //车辆品牌
	CarLoad              int      //车载量
	DrOrder              []*Order `orm:"reverse(many)"` //一个司机有多个完成的订单
}

//Order 是完成后的订单存储
type Order struct {
	OrId        int64     `orm:"pk;auto"`                          //订单编号，自动增加确定
	Orgin       string    `orm:"size(30)"`                         //起始地
	Destination string    `orm:"size(30)"`                         //目的地
	BeginTime   time.Time `orm:"auto_now_add;type(datetime)"`      //开始时间，我设定为下单时自动填入
	Endtime     time.Time `orm:"auto_now_add;type(datetime);null"` //结束时间，我设定为接单时自动填入
	Node        string    `orm:"null"`                             //备注
	Cost        float64   `orm:"null"`                             //车费
	Client      *Client   `orm:"rel(fk)"`                          //乘客类，是外键
	Driver      *Driver   `orm:"rel(fk);null"`                     //司机类，是外键
}

//Message 是留言信息表
type Message struct {
	MessageId int64  `orm:"pk;auto"`  //留言编号
	WriteBy   string `orm:"size(20)"` //留言者信息
	Node      string `orm:"null"`     //留言信息
}

//DBA 是管理员信息，默认账户root,密码123
//出于安简便，管理员就只有查看权，账户密码也是固定的
type DBA struct {
	Name     string `orm:"pk"`
	Password string `orm:"size(20)"`
}

func init() {
	// 设置数据库基本信息
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/testDidi?charset=utf8&loc=Local")
	// 映射model数据
	orm.RegisterModel(new(Client), new(Driver), new(Order), new(Message))
	// 生成表
	orm.RunSyncdb("default", false, true)
}
