package controllers

import (
	// "fmt"
	. "mymanager/models/category"
	. "mymanager/models/common"
	. "mymanager/models/user"
)

type CategoryController struct {
	BaseController
	CategoryAo  CategoryAoModel
	UserLoginAo UserLoginAoModel
}

func (this *CategoryController) Search_Json() interface{} {

	//检查输入参数
	var where Category
	this.CheckGet(&where)

	var limit CommonPage
	this.CheckGet(&limit)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//执行业务逻辑
	return this.CategoryAo.Search(user.UserId, where, limit)
}

func (this *CategoryController) Get_Json() Category {
	//检查输入
	var data struct {
		CategoryId int
	}
	this.CheckGet(&data)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	return this.CategoryAo.Get(user.UserId, data.CategoryId)
}

func (this *CategoryController) Add_Json() {
	//检查输入
	category := Category{}
	this.CheckPost(&category)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	this.CategoryAo.Add(user.UserId, category)
}

func (this *CategoryController) Mod_Json() {

	category := Category{}
	this.CheckPost(&category)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	this.CategoryAo.Mod(user.UserId, category)
}

func (this *CategoryController) Del_Json() {
	//检查输入
	var data struct {
		CategoryId int
	}
	this.CheckPost(&data)

	//检查权限
	user := this.UserLoginAo.CheckMustLogin()

	//业务逻辑
	this.CategoryAo.Del(user.UserId, data.CategoryId)
}
