package category

import (
	. "github.com/fishedee/language"
	. "mymanager/models/common"
	// "crypto/sha1"
	"fmt"
	// "io"
)

type CategoryAoModel struct {
	BaseModel
	CategoryDb CategoryDbModel
}

func (this *CategoryAoModel) Search(userId int, data Category, pageInfo CommonPage) Categorys {

	wheres := Category{
		UserId: userId,
		Name:   data.Name,
		Remark: data.Remark,
	}

	fmt.Printf("%+v", wheres)

	return this.CategoryDb.Search(wheres, pageInfo)

}

func (this *CategoryAoModel) Get(userId, categoryId int) Category {
	category := this.CategoryDb.Get(categoryId)
	if category.UserId != userId {
		Throw(1, "你没有权利查看或编辑等操作")
	}
	return category
}

func (this *CategoryAoModel) Add(userId int, category Category) {

	categorys := Category{
		UserId: userId,
		Name:   category.Name,
		Remark: category.Remark,
	}

	this.CategoryDb.Add(categorys)
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
}
