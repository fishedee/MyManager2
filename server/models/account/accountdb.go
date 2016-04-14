package account

import (
	"fmt"
	. "github.com/fishedee/language"
	. "mymanager/models/common"
	"strconv"
)

type AccountDbModel struct {
	BaseModel
}

func (this *AccountDbModel) Search(where Account, limit CommonPage) Accounts {
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
	if where.CategoryId != 0 {
		db = db.And("categoryId = ?", where.CategoryId)
	}
	if where.CardId != 0 {
		db = db.And("cardId = ?", where.CardId)
	}
	if where.Type != 0 {
		db = db.And("type = ?", where.Type)
	}

	data := []Account{}
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
	if where.CategoryId != 0 {
		db = db.And("categoryId = ?", where.CategoryId)
	}
	if where.CardId != 0 {
		db = db.And("cardId = ?", where.CardId)
	}
	if where.Type != 0 {
		db = db.And("type = ?", where.Type)
	}

	count, err := db.And("userId = ?", where.UserId).Count(&Account{})
	if err != nil {
		panic(err)
	}

	return Accounts{
		Count: int(count),
		Data:  data,
	}

}

func (this *AccountDbModel) Get(id int) Account {
	var account []Account
	err := this.DB.Where("accountId = ?", id).Find(&account)
	if err != nil {
		panic(err)
	}

	if len(account) == 0 {
		Throw(1, "该"+strconv.Itoa(account[0].AccountId)+"类型不存在")
	}

	return account[0]
}

func (this *AccountDbModel) Add(accounts Account) int {
	_, err := this.DB.Insert(&accounts)
	if err != nil {
		panic(err)
	}

	return accounts.AccountId
}

func (this *AccountDbModel) Mod(account Account) {
	_, err := this.DB.Where("accountId = ?", account.AccountId).Update(&account)
	if err != nil {
		panic(err)
	}
}

func (this *AccountDbModel) updateCategoryIdByZero(categoryId int) {
	var account Account
	fmt.Println(categoryId)
	fmt.Println("update")
	_, err := this.DB.Where("categoryId = ?", categoryId).Cols("categoryId").Update(&account)
	if err != nil {
		panic(err)
	}
}

func (this *AccountDbModel) Del(accountId int) {
	_, err := this.DB.Where("accountId = ?", accountId).Delete(&Account{})
	if err != nil {
		panic(err)
	}
}

func (this *AccountDbModel) AccountJoinCategory(userId int, thisType int, startTime string, endTime string) []WeekDetailTypeStatistic {
	var data []WeekDetailTypeStatistic
	// err := this.DB.Sql("select * from t_account where userId=? AND type=? AND createTime>=? AND createTime<=?;", userId, thisType, startTime, endTime).Find(&data)
	err := this.DB.Sql("select ca.categoryId,ca.name as categoryName,ac.money,ac.createTime from t_account as ac inner join t_category as ca ON ac.categoryId=ca.categoryId where ac.userId=? AND ac.type=? AND ac.createTime>=? AND ac.createTime<=?", userId, thisType, startTime, endTime).Find(&data)

	if err != nil {
		panic(err)
	}

	if len(data) == 0 {
		Throw(1, "你所寻找的资料不存在")
	}

	return data
}

func (this *AccountDbModel) GetWeekCardStatistic(userId int) []WeekCardStatistic {
	var data []WeekCardStatistic
	err := this.DB.Sql("select ac.cardId,card.name as cardName,card.money as cardMoney, ac.money as accountMoney,ac.name,ac.type,ac.CreateTime,ac.type from t_account as ac inner join t_card card  on ac.cardId=card.cardId where ac.userId=?", userId).Find(&data)

	if err != nil {
		panic(err)
	}

	if len(data) == 0 {
		Throw(1, "你所寻找的资料不存在")
	}
	return data
}

func (this *AccountDbModel) GetWeekDetailCardStatistic(userId int, cardId int, startTime string, endTime string) []WeekDetailTypeStatistic {
	var data []WeekDetailTypeStatistic
	err := this.DB.Sql("select * from t_account where userId=? AND cardId=? AND createTime>=? AND createTime<=?", userId, cardId, startTime, endTime).Find(&data)

	if err != nil {
		panic(err)
	}

	if len(data) == 0 {
		Throw(1, "你所寻找的资料不存在")
	}
	return data
}
