package card

import (
	. "github.com/fishedee/language"
	. "mymanager/models/common"
	"strconv"
)

type CardDbModel struct {
	BaseModel
}

func (this *CardDbModel) Search(where Card, limit CommonPage) Cards {
	db := this.DB.NewSession()
	defer db.Close()

	if limit.PageSize == 0 && limit.PageIndex == 0 {
		return Cards{
			Count: 0,
			Data:  []Card{},
		}
	}

	if where.Name != "" {
		db = db.And("name like ?", "%"+where.Name+"%")
	}

	if where.Remark != "" {
		db = db.And("remark like ?", "%"+where.Remark+"%")
	}

	data := []Card{}
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

	count, err := db.And("userId = ?", where.UserId).Count(&Card{})
	if err != nil {
		panic(err)
	}

	return Cards{
		Count: int(count),
		Data:  data,
	}

}

func (this *CardDbModel) Get(id int) Card {
	var card []Card
	err := this.DB.Where("cardId = ?", id).Find(&card)
	if err != nil {
		panic(err)
	}

	if len(card) == 0 {
		Throw(1, "该"+strconv.Itoa(card[0].CardId)+"银行卡不存在")
	}

	return card[0]
}

func (this *CardDbModel) Add(cards Card) int {
	_, err := this.DB.Insert(&cards)
	if err != nil {
		panic(err)
	}

	return cards.CardId
}

func (this *CardDbModel) Mod(card Card) {
	_, err := this.DB.Where("cardId = ?", card.CardId).Update(&card)
	if err != nil {
		panic(err)
	}
}

func (this *CardDbModel) Del(cardId int) {
	_, err := this.DB.Where("cardId = ?", cardId).Delete(&Card{})
	if err != nil {
		panic(err)
	}
}
