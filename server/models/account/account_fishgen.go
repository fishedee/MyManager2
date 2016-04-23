package account

import . "mymanager/models/common"
import . "github.com/fishedee/language"

func (this *AccountAoModel) Search_WithError(userId int, where Account, pageInfo CommonPage) (_fishgen1 Accounts, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Search(userId, where, pageInfo)
	return
}

func (this *AccountAoModel) Get_WithError(userId int, accountId int) (_fishgen1 Account, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Get(userId, accountId)
	return
}

func (this *AccountAoModel) Add_WithError(userId int, account Account) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Add(userId, account)
	return
}

func (this *AccountAoModel) Mod_WithError(userId int, account Account) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Mod(userId, account)
	return
}

func (this *AccountAoModel) Del_WithError(userId int, accountId int) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Del(userId, accountId)
	return
}

func (this *AccountAoModel) GetWeekTypeStatistic_WithError(userId int) (_fishgen1 []WeekStatistic, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.GetWeekTypeStatistic(userId)
	return
}

func (this *AccountAoModel) GetWeekDetailTypeStatistic_WithError(userId int, data WeekStatistic) (_fishgen1 []WeekDetailStatistic, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.GetWeekDetailTypeStatistic(userId, data)
	return
}

func (this *AccountAoModel) GetWeekCardStatistic_WithError(userId int) (_fishgen1 []WeekStatistic, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.GetWeekCardStatistic(userId)
	return
}

func (this *AccountAoModel) GetWeekDetailCardStatistic_WithError(userId int, data WeekStatistic) (_fishgen1 []WeekDetailStatistic, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.GetWeekDetailCardStatistic(userId, data)
	return
}

func (this *AccountDbModel) Search_WithError(where Account, limit CommonPage) (_fishgen1 Accounts, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Search(where, limit)
	return
}

func (this *AccountDbModel) Get_WithError(id int) (_fishgen1 Account, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Get(id)
	return
}

func (this *AccountDbModel) Add_WithError(accounts Account) (_fishgen1 int, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Add(accounts)
	return
}

func (this *AccountDbModel) Mod_WithError(account Account) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Mod(account)
	return
}

func (this *AccountDbModel) Del_WithError(accountId int) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Del(accountId)
	return
}

func (this *AccountDbModel) GetWeekDetailTypeStatistic_WithError(userId int, accountType int, startTime string, endTime string) (_fishgen1 []Account, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.GetWeekDetailTypeStatistic(userId, accountType, startTime, endTime)
	return
}

func (this *AccountDbModel) GetWeekDetailCardStatistic_WithError(userId int, cardId int, startTime string, endTime string) (_fishgen1 []Account, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.GetWeekDetailCardStatistic(userId, cardId, startTime, endTime)
	return
}
