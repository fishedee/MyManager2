package category

import (
	. "github.com/fishedee/language"
	. "mymanager/models/common"
	// "crypto/sha1"
	// "fmt"
	// "io"
)

type CategoryAoModel struct {
	BaseModel
	CategoryDb CategoryDbModel
}

func (this *CategoryAoModel) Search(userId int, where Category, pageInfo CommonPage) Categorys {

	where.UserId = userId

	return this.CategoryDb.Search(where, pageInfo)

}

func (this *CategoryAoModel) Get(userId, categoryId int) Category {
	category := this.CategoryDb.Get(categoryId)
	if category.UserId != userId {
		Throw(1, "你没有权利查看或编辑等操作")
	}
	return category
}

func (this *CategoryAoModel) Add(userId int, category Category) {

	category.UserId = userId
	this.CategoryDb.Add(category)
}

func (this *CategoryAoModel) Mod(userId int, category Category) {

	//检查该类型是不是属于他本人
	this.Get(userId, category.CategoryId)

	this.CategoryDb.Mod(category)
}
func (this *CategoryAoModel) Del(userId, categoryId int) {

	//检查该类型是不是属于他本人
	this.Get(userId, categoryId)

	this.CategoryDb.Del(categoryId)
	this.Queue.Publish("/category/_del", categoryId)
}
