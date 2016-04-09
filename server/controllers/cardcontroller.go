package controllers

import (
	// "fmt"
	. "mymanager/models/card"
	. "mymanager/models/common"
	. "mymanager/models/user"
)

type CardController struct {
	BaseController
	CardAo      CardAoModel
	UserLoginAo UserLoginAoModel
}

func (this *CardController) Search_Json() interface{} {

	//检查输入参数
	var where Card
	this.CheckGet(&where)

	var limit CommonPage
	this.CheckGet(&limit)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//执行业务逻辑
	return this.CardAo.Search(user.UserId, where, limit)
}

func (this *CardController) Get_Json() Card {
	//检查输入
	var data struct {
		CardId int
	}
	this.CheckGet(&data)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	return this.CardAo.Get(user.UserId, data.CardId)
}

func (this *CardController) Add_Json() {
	//检查输入
	card := Card{}
	this.CheckPost(&card)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	this.CardAo.Add(user.UserId, card)
}

func (this *CardController) Mod_Json() {
	//检查输入
	card := Card{}
	this.CheckPost(&card)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	this.CardAo.Mod(user.UserId, card)
}

func (this *CardController) Del_Json() {
	//检查输入
	var data struct {
		CardId int
	}
	this.CheckPost(&data)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	this.CardAo.Del(user.UserId, data.CardId)
}
