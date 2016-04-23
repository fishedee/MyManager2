package card

import (
	. "github.com/fishedee/language"
	. "mymanager/models/common"
)

type CardAoModel struct {
	BaseModel
	CardDb CardDbModel
}

func (this *CardAoModel) Search(userId int, where Card, pageInfo CommonPage) Cards {

	where.UserId = userId

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

	this.Queue.Publish(CardQueueEnum.EVENT_DEL, cardId)

}
