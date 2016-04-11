package card

import (
	. "github.com/fishedee/language"
	. "mymanager/models/common"
	// "crypto/sha1"
	// "fmt"
	// "io"
)

type CardAoModel struct {
	BaseModel
	CardDb CardDbModel
}

func (this *CardAoModel) Search(userId int, where Card, pageInfo CommonPage) Cards {

	// wheres := Card{
	// 	UserId: userId,
	// 	Name:   data.Name,
	// 	Remark: data.Remark,
	// }

	where.UserId = userId

	// fmt.Printf("%+v", wheres)

	return this.CardDb.Search(where, pageInfo)

}

func (this *CardAoModel) Get(userId, cardId int) Card {
	card := this.CardDb.Get(cardId)
	if card.UserId != userId {
		Throw(1, "你没有权利查看或编辑等操作")
	}
	return card
}

func (this *CardAoModel) Add(userId int, card Card) {

	// cards := Card{
	// 	UserId: userId,
	// 	Name:   card.Name,
	// 	Bank:   card.Bank,
	// 	Card:   card.Card,
	// 	Money:  card.Money,
	// 	Remark: card.Remark,
	// }
	card.UserId = userId

	this.CardDb.Add(card)
}

func (this *CardAoModel) Mod(userId int, card Card) {

	//检查该类型是不是属于他本人
	this.Get(userId, card.CardId)

	this.CardDb.Mod(card)
}
func (this *CardAoModel) Del(userId, cardId int) {

	//检查该类型是不是属于他本人
	this.Get(userId, cardId)

	this.CardDb.Del(cardId)
}
