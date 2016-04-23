package category

import . "mymanager/models/common"
import . "github.com/fishedee/language"

func (this *CategoryAoModel) Search_WithError(userId int, where Category, pageInfo CommonPage) (_fishgen1 Categorys, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Search(userId, where, pageInfo)
	return
}

func (this *CategoryAoModel) Get_WithError(userId int, categoryId int) (_fishgen1 Category, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Get(userId, categoryId)
	return
}

func (this *CategoryAoModel) Add_WithError(userId int, category Category) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Add(userId, category)
	return
}

func (this *CategoryAoModel) Mod_WithError(userId int, category Category) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Mod(userId, category)
	return
}

func (this *CategoryAoModel) Del_WithError(userId int, categoryId int) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Del(userId, categoryId)
	return
}

func (this *CategoryDbModel) Search_WithError(where Category, limit CommonPage) (_fishgen1 Categorys, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Search(where, limit)
	return
}

func (this *CategoryDbModel) Get_WithError(id int) (_fishgen1 Category, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Get(id)
	return
}

func (this *CategoryDbModel) Add_WithError(categorys Category) (_fishgen1 int, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Add(categorys)
	return
}

func (this *CategoryDbModel) Mod_WithError(category Category) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Mod(category)
	return
}

func (this *CategoryDbModel) Del_WithError(categoryId int) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Del(categoryId)
	return
}
