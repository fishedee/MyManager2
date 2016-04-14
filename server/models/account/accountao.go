package account

import (
	"errors"
	"fmt"
	. "github.com/fishedee/language"
	. "github.com/fishedee/web"
	. "mymanager/models/card"
	. "mymanager/models/category"
	. "mymanager/models/common"
	"time"
)

type AccountAoModel struct {
	BaseModel
	AccountDb  AccountDbModel
	CategoryAo CategoryAoModel
	CardAo     CardAoModel
}

func (this *AccountAoModel) Search(userId int, where Account, pageInfo CommonPage) Accounts {

	where.UserId = userId

	return this.AccountDb.Search(where, pageInfo)

}

func (this *AccountAoModel) Get(userId, accountId int) Account {
	account := this.AccountDb.Get(accountId)
	if account.UserId != userId {
		Throw(1, "你没有权利查看或编辑等操作")
	}
	return account
}

func (this *AccountAoModel) Add(userId int, account Account) {

	account.UserId = userId

	this.AccountDb.Add(account)
}

func (this *AccountAoModel) Mod(userId int, account Account) {

	//检查该类型是不是属于他本人
	this.Get(userId, account.AccountId)

	this.AccountDb.Mod(account)
}
func (this *AccountAoModel) Del(userId, accountId int) {

	//检查该类型是不是属于他本人
	this.Get(userId, accountId)

	this.AccountDb.Del(accountId)
}

func (this *AccountAoModel) GetWeekTypeStatistic(userId int) []WeekStatistic {
	var where Account
	var limit CommonPage
	accountSearchData := this.Search(userId, where, limit)
	enums := AccountTypeEnum.Datas()

	//时间转换星期
	yearWeekCardStatistic := QuerySelect(accountSearchData.Data, func(singleData Account) WeekStatistic {
		t, err := time.Parse("2006-01-02 15:05:04", singleData.CreateTime)
		if err != nil {
			panic(err)
		}
		year, week := t.ISOWeek()
		return WeekStatistic{
			Money: singleData.Money,
			Type:  singleData.Type,
			Week:  week,
			Year:  year,
		}
	}).([]WeekStatistic)

	yearWeekStatistic := QueryGroup(yearWeekCardStatistic, "Year desc,Week desc", func(list []WeekStatistic) []WeekStatistic {
		//合并单个type下面的多个account
		result := QueryGroup(list, "Type", func(list []WeekStatistic) []WeekStatistic {
			sum := QuerySum(QueryColumn(list, "Money"))
			list[0].Money = sum.(int)
			return []WeekStatistic{list[0]}
		}).([]WeekStatistic)

		//left join所有type
		single := list[0]
		result = QueryLeftJoin(enums, result, "Id = Type", func(left EnumData, right WeekStatistic) WeekStatistic {
			return WeekStatistic{
				Year: single.Year,
				Week: single.Week,
				Name: fmt.Sprintf(
					"%4d年%02d周",
					single.Year,
					single.Week,
				),
				Type:     left.Id,
				TypeName: left.Name,
				Money:    right.Money,
			}
		}).([]WeekStatistic)

		return result
	}).([]WeekStatistic)

	//排序
	result := QuerySort(yearWeekStatistic, "Year desc,Week desc,Type asc").([]WeekStatistic)

	return result

}
func (this *AccountAoModel) GetWeekDetailTypeStatistic(userId int, data WeekStatistic) []WeekDetailStatistic {
	//获取那个星期的时间范围
	t := time.Now()
	timeLocation := t.Location()
	thisWeekStartTime := firstDayOfISOWeek(data.Year, data.Week, timeLocation)
	thisWeekEndTime := thisWeekStartTime.AddDate(0, 0, 7)

	layout := "2006-01-02"
	startTimeString := thisWeekStartTime.Format(layout)
	endTimeString := thisWeekEndTime.Format(layout)
	timeRangeData := this.AccountDb.GetWeekDetailTypeStatistic(userId, data.Type, startTimeString, endTimeString)
	categoryData := this.CategoryAo.Search(userId, Category{}, CommonPage{})

	//合并category下的account
	result := QueryGroup(timeRangeData, "CategoryId desc", func(list []Account) []Account {
		sum := QuerySum(QueryColumn(list, "Money"))
		list[0].Money = sum.(int)
		return []Account{list[0]}
	}).([]Account)

	//计算比例
	totalPriceFloat := (float64)(QuerySum(QueryColumn(result, "Money")).(int))
	resultStatistic := QueryLeftJoin(result, categoryData.Data, "CategoryId = CategoryId", func(left Account, right Category) WeekDetailStatistic {
		categoryMoneyFloat := (float64)(left.Money)
		categoryMoneyScale := categoryMoneyFloat / totalPriceFloat * 100
		categoryMoneyScaleString := fmt.Sprintf("%.2f", categoryMoneyScale)
		return WeekDetailStatistic{
			CategoryId:   left.CategoryId,
			CategoryName: right.Name,
			Precent:      categoryMoneyScaleString,
			Money:        left.Money,
		}
	}).([]WeekDetailStatistic)

	return resultStatistic

}

func (this *AccountAoModel) GetWeekCardStatistic(userId int) []WeekStatistic {
	accountSearchData := this.Search(userId, Account{}, CommonPage{})

	//时间转换星期
	yearWeekCardStatistic := QuerySelect(accountSearchData.Data, func(singleData Account) WeekStatistic {
		t, err := time.Parse("2006-01-02 15:05:04", singleData.CreateTime)
		if err != nil {
			panic(err)
		}
		year, week := t.ISOWeek()
		return WeekStatistic{
			Money:  singleData.Money,
			Type:   singleData.Type,
			CardId: singleData.CardId,
			Week:   week,
			Year:   year,
		}
	}).([]WeekStatistic)

	//取出各个card的初始money
	cardSearch := this.CardAo.Search(userId, Card{}, CommonPage{})
	initMoney := map[int]int{}
	for _, singleCard := range cardSearch.Data {
		initMoney[singleCard.CardId] = singleCard.Money
	}

	this.Log.Debug("%#v", yearWeekCardStatistic)

	yearWeekCardStatistic = QueryGroup(yearWeekCardStatistic, "Year asc,Week asc", func(list []WeekStatistic) []WeekStatistic {
		//合并单个card下面的多个account
		result := QueryGroup(list, "CardId", func(list []WeekStatistic) []WeekStatistic {
			cardTotalPice := QuerySum(QuerySelect(list, func(single WeekStatistic) int {
				switch single.Type {
				case AccountTypeEnum.INCOME, AccountTypeEnum.TRANSFER_INCOME, AccountTypeEnum.ACCOUNT_RECEIVABLE:
					return single.Money
				case AccountTypeEnum.SPENDING, AccountTypeEnum.TRANSFER_SPENDING, AccountTypeEnum.ACCOUNTS_PAYABLE:
					return -single.Money
				default:
					panic(errors.New("不存在该类型"))
				}
			})).(int)
			list[0].Money = cardTotalPice
			return []WeekStatistic{list[0]}
		}).([]WeekStatistic)

		this.Log.Debug("%#v", result)

		//left join所有card
		single := list[0]
		result = QueryLeftJoin(cardSearch.Data, result, "CardId = CardId", func(left Card, right WeekStatistic) WeekStatistic {
			initMoney[left.CardId] = initMoney[left.CardId] + right.Money
			return WeekStatistic{
				Year: single.Year,
				Week: single.Week,
				Name: fmt.Sprintf(
					"%4d年%02d周",
					single.Year,
					single.Week,
				),
				Money:    initMoney[left.CardId],
				CardId:   left.CardId,
				CardName: left.Name,
			}
		}).([]WeekStatistic)

		return result
	}).([]WeekStatistic)

	//排序
	result := QuerySort(yearWeekCardStatistic, "Year desc,Week desc,CardId asc").([]WeekStatistic)

	return result

}

func (this *AccountAoModel) GetWeekDetailCardStatistic(userId int, data WeekStatistic) []WeekDetailStatistic {
	//获取那个星期的时间范围
	t := time.Now()
	timeLocation := t.Location()
	thisWeekStartTime := firstDayOfISOWeek(data.Year, data.Week, timeLocation)
	thisWeekEndTime := thisWeekStartTime.AddDate(0, 0, 7)

	layout := "2006-01-02"
	startTimeString := thisWeekStartTime.Format(layout)
	endTimeString := thisWeekEndTime.Format(layout)
	timeRangeData := this.AccountDb.GetWeekDetailCardStatistic(userId, data.CardId, startTimeString, endTimeString)
	enums := AccountTypeEnum.Datas()

	//合并type下的account
	result := QueryGroup(timeRangeData, "Type asc", func(list []Account) []Account {
		sum := QuerySum(QueryColumn(list, "Money"))
		list[0].Money = sum.(int)
		return []Account{list[0]}
	}).([]Account)

	//计算比例
	totalPriceFloat := (float64)(QuerySum(QueryColumn(result, "Money")).(int))
	resultStatistic := QueryLeftJoin(result, enums, "Type = Id", func(left Account, right EnumData) WeekDetailStatistic {
		categoryMoneyFloat := (float64)(left.Money)
		categoryMoneyScale := categoryMoneyFloat / totalPriceFloat * 100
		categoryMoneyScaleString := fmt.Sprintf("%.2f", categoryMoneyScale)
		return WeekDetailStatistic{
			Type:     left.Type,
			TypeName: right.Name,
			Precent:  categoryMoneyScaleString,
			Money:    left.Money,
		}
	}).([]WeekDetailStatistic)

	return resultStatistic

}

func (this *AccountAoModel) whenCategoryDel(categoryId int) {
	this.AccountDb.updateCategoryIdByZero(categoryId)
}

func firstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()

	// iterate back to Monday
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the first week
	for isoYear < year {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the given week
	for isoWeek < week {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	return date
}

func init() {
	InitDaemon(func(this *AccountAoModel) {
		this.Queue.Subscribe("/category/_del", this.whenCategoryDel)
	})
}
