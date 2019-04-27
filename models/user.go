package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	myutils "wac/wacUtils"
)

type User struct {
	Id       int
	Username string
	Password string
	Email    string
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterModel(new(User))
	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("mysql"))
}

func Login(username string, password string) (*User, string) {
	if len(username) == 0 || len(password) == 0 {
		return &User{}, "缺少字段，登录失败"
	}
	o := orm.NewOrm()

	user := &User{Username: username}
	err := o.Read(user, "username")
	if err != nil {
		fmt.Println(err)
		return &User{}, "用户名不存在"
	}
	if myutils.CheckPasswordHash(password, user.Password) {
		user.Password = ""
		return user, ""
	}
	return &User{}, "用户名或密码错误"
}

func GetUser(username string) (u *User, err error) {
	o := orm.NewOrm()
	//o.("wac")
	user := &User{Username: username}

	if o.Read(user) != nil && user.Id != 0 {
		return user, nil
	}

	return user, errors.New("获取失败")
}
