package controllers

import (
	// . "github.com/fishedee/language"
	// . "github.com/fishedee/web"
	// "fmt"

	// . "mymanager/models/common"
	// . "github.com/fishedee/language"
	// "fmt"
	. "mymanager/models/user"
)

type LoginController struct {
	BaseController
	UserLoginAo UserLoginAoModel
	// ClientLoginAo ClientLoginAoModel
}

//检查是否登录
func (this *LoginController) Islogin_Json() {

	this.UserLoginAo.CheckMustLogin()

}

//登录操作
func (this *LoginController) Checkin_Json() {

	//检查输入参数
	var user User
	this.CheckPost(&user)

	this.UserLoginAo.Login(user.Name, user.Password)
}

func (this *LoginController) Checkout_Json() {

	this.UserLoginAo.Logout()
}
