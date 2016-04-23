package user

import (
	// . "github.com/fishedee/language"
	. "github.com/fishedee/language"
	. "mymanager/models/common"
	"strconv"
)

type UserDbModel struct {
	BaseModel
}

func (this *UserDbModel) Search(where User, limit CommonPage) Users {
	db := this.DB.NewSession()
	defer db.Close()

	if limit.PageSize == 0 && limit.PageIndex == 0 {
		return Users{
			Count: 0,
			Data:  []User{},
		}
	}

	if where.Name != "" {
		db = db.And("name like ?", "%"+where.Name+"%")
	}

	if where.Type != 0 {
		db = db.And("type = ?", where.Type)
	}

	data := []User{}
	err := db.OrderBy("createTime desc").Limit(limit.PageSize, limit.PageIndex).Find(&data)
	if err != nil {
		panic(err)
	}

	if where.Name != "" {
		db = db.And("name like ?", "%"+where.Name+"%")
	}

	if where.Type != 0 {
		db = db.And("type = ?", where.Type)
	}

	count, err := db.Count(&User{})
	if err != nil {
		panic(err)
	}

	return Users{
		Count: int(count),
		Data:  data,
	}

}

func (this *UserDbModel) GetByName(name string) []User {
	var users []User
	err := this.DB.Where("name = ?", name).Find(&users)
	if err != nil {
		panic(err)
	}
	return users
}

func (this *UserDbModel) Get(id int) User {
	var users []User
	err := this.DB.Where("userId = ?", id).Find(&users)
	if err != nil {
		panic(err)
	}

	if len(users) == 0 {
		Throw(1, "该"+strconv.Itoa(id)+"用户不存在")
	}

	return users[0]
}
func (this *UserDbModel) Add(user User) int {
	_, err := this.DB.Insert(&user)
	if err != nil {
		panic(err)
	}

	return user.UserId
}

func (this *UserDbModel) Mod(id int, user User) {
	_, err := this.DB.Where("userId = ?", id).Update(&user)
	if err != nil {
		panic(err)
	}

}

func (this *UserDbModel) Del(id int) {
	_, err := this.DB.Where("userId = ?", id).Delete(&User{})
	if err != nil {
		panic(err)
	}
}
