package controllers

import (
	// "fmt"
	. "mymanager/models/account"
	. "mymanager/models/common"
	. "mymanager/models/user"
	// "time"
)

type AccountController struct {
	BaseController
	AccountAo   AccountAoModel
	UserLoginAo UserLoginAoModel
}

func (this *AccountController) Search_Json() interface{} {

	//检查输入参数
	var where Account
	this.CheckGet(&where)

	var limit CommonPage
	this.CheckGet(&limit)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//执行业务逻辑
	return this.AccountAo.Search(user.UserId, where, limit)
}

func (this *AccountController) Get_Json() Account {
	//检查输入
	var data struct {
		AccountId int
	}
	this.CheckGet(&data)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	return this.AccountAo.Get(user.UserId, data.AccountId)
}

func (this *AccountController) Add_Json() {
	//检查输入
	account := Account{}
	this.CheckPost(&account)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	this.AccountAo.Add(user.UserId, account)
}

func (this *AccountController) Mod_Json() {

	account := Account{}
	this.CheckPost(&account)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	this.AccountAo.Mod(user.UserId, account)
}

func (this *AccountController) Del_Json() {
	//检查输入
	var data struct {
		AccountId int
	}
	this.CheckPost(&data)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	this.AccountAo.Del(user.UserId, data.AccountId)
}

func (this *AccountController) GetWeekTypeStatistic_Json() interface{} {

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	return this.AccountAo.GetWeekTypeStatistic(user.UserId)

}

func (this *AccountController) GetWeekDetailTypeStatistic_Json() interface{} {

	//检查输入
	var data WeekTypeStatistic
	this.CheckGet(&data)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	return this.AccountAo.GetWeekDetailTypeStatistic(user.UserId, data)
}
