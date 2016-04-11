package account

import (
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

func (this *AccountDbModel) Del(accountId int) {
	_, err := this.DB.Where("accountId = ?", accountId).Delete(&Account{})
	if err != nil {
		panic(err)
	}
}

func (this *AccountDbModel) AccountJoinCard(userId int) []WeekTypeStatistic {
	var data []WeekTypeStatistic
	err := this.DB.Sql("select ac.cardId,card.name as cardName,ac.money,ac.name,ac.type,ac.CreateTime,ac.ModifyTime  from t_account as ac inner join t_card card  on ac.cardId=card.cardId where ac.userId=?", userId).Find(&data)

	if err != nil {
		panic(err)
	}
	return data
}
