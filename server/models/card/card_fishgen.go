package card

import . "github.com/fishedee/language"
import . "mymanager/models/common"

func (this *CardAoModel) Search_WithError(userId int, where Card, pageInfo CommonPage) (_fishgen1 Cards, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Search(userId, where, pageInfo)
	return
}

func (this *CardAoModel) Get_WithError(userId int, cardId int) (_fishgen1 Card, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Get(userId, cardId)
	return
}

func (this *CardAoModel) Add_WithError(userId int, card Card) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Add(userId, card)
	return
}

func (this *CardAoModel) Mod_WithError(userId int, card Card) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Mod(userId, card)
	return
}

func (this *CardAoModel) Del_WithError(userId int, cardId int) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Del(userId, cardId)
	return
}

func (this *CardDbModel) Search_WithError(where Card, limit CommonPage) (_fishgen1 Cards, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Search(where, limit)
	return
}

func (this *CardDbModel) Get_WithError(id int) (_fishgen1 Card, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Get(id)
	return
}

func (this *CardDbModel) Add_WithError(cards Card) (_fishgen1 int, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Add(cards)
	return
}

func (this *CardDbModel) Mod_WithError(card Card) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Mod(card)
	return
}

func (this *CardDbModel) Del_WithError(cardId int) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Del(cardId)
	return
}
