package category

import (
	. "github.com/fishedee/language"
	. "mymanager/models/common"
	"strconv"
)

type CategoryDbModel struct {
	BaseModel
}

func (this *CategoryDbModel) Search(where Category, limit CommonPage) Categorys {
	db := this.DB.NewSession()
	defer db.Close()

	if limit.PageSize == 0 {
		limit.PageSize = 50
	}

	if where.Name != "" {
		db = db.And("name like ?", "%"+where.Name+"%")
	}

	if where.Remark != "" {
		db = db.And("remark like ?", "%"+where.Remark+"%")
	}

	data := []Category{}
	err := db.And("userId = ?", where.UserId).OrderBy("createTime desc").Limit(limit.PageSize, limit.PageIndex).Find(&data)
	if err != nil {
		panic(err)
	}

	if where.Name != "" {
		db = db.And("name like ?", "%"+where.Name+"%")
	}

	if where.Remark != "" {
		db = db.And("remark like ?", "%"+where.Remark+"%")
	}

	count, err := db.And("userId = ?", where.UserId).Count(&Category{})
	if err != nil {
		panic(err)
	}

	return Categorys{
		Count: int(count),
		Data:  data,
	}

}

func (this *CategoryDbModel) Get(id int) Category {
	var category []Category
	err := this.DB.Where("categoryId = ?", id).Find(&category)
	if err != nil {
		panic(err)
	}

	if len(category) == 0 {
		Throw(1, "该"+strconv.Itoa(category[0].CategoryId)+"类型不存在")
	}

	return category[0]
}

func (this *CategoryDbModel) Add(categorys Category) int {
	_, err := this.DB.Insert(&categorys)
	if err != nil {
		panic(err)
	}

	return categorys.CategoryId
}

func (this *CategoryDbModel) Mod(category Category) {
	_, err := this.DB.Where("categoryId = ?", category.CategoryId).Update(&category)
	if err != nil {
		panic(err)
	}
}

func (this *CategoryDbModel) Del(categoryId int) {
	_, err := this.DB.Where("categoryId = ?", categoryId).Delete(&Category{})
	if err != nil {
		panic(err)
	}
}
