package account

import (
	. "github.com/fishedee/language"
	. "github.com/fishedee/web"
	. "mymanager/models/card"
	. "mymanager/models/category"
	. "mymanager/models/common"
	// "crypto/sha1"
	"fmt"
	"strconv"
	"time"
	// "io"
	"errors"
	// "math"
	// . "github.com/liudng/godump"
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

func (this *AccountAoModel) GetWeekTypeStatistic(userId int) []WeekTypeStatistic {
	var where Account
	var limit CommonPage
	accountSearchData := this.Search(userId, where, limit)
	enums := accountEnum.Datas()

	// fmt.Println("\naccountSearchData")
	// fmt.Printf("%+v", accountSearchData)

	yearWeekCardStatistic := QuerySelect(accountSearchData.Data, func(singleData Account) WeekTypeStatistic {
		t, err := time.Parse("2006-01-02 15:05:04", singleData.CreateTime)
		if err != nil {
			panic(err)
		}
		year, week := t.ISOWeek()
		return WeekTypeStatistic{
			Money: singleData.Money,
			Type:  singleData.Type,
			Week:  week,
			Year:  year,
		}
	}).([]WeekTypeStatistic)

	// fmt.Println("\nyearWeekCardStatistic")

	// godump.Dump(yearWeekCardStatistic)

	yearWeekStatistic := QueryGroup(yearWeekCardStatistic, "Year desc,Week desc,Type desc", func(list []WeekTypeStatistic) []WeekTypeStatistic {
		sum := QuerySum(QueryColumn(list, "Money").([]int))
		list[0].Money = sum.(int)
		return []WeekTypeStatistic{list[0]}
	}).([]WeekTypeStatistic)

	// fmt.Println("\nyearWeekStatistic")
	// godump.Dump(yearWeekStatistic)

	result := QueryGroup(yearWeekStatistic, "Year desc,Week desc", func(list []WeekTypeStatistic) []WeekTypeStatistic {
		single := list[0]
		result := QueryLeftJoin(enums, list, "Id = Type", func(left EnumData, right WeekTypeStatistic) WeekTypeStatistic {
			return WeekTypeStatistic{
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
		}).([]WeekTypeStatistic)
		return result
	}).([]WeekTypeStatistic)

	// fmt.Println("\nresult")
	// godump.Dump(result)

	return result

}
func (this *AccountAoModel) GetWeekDetailTypeStatistic(userId int, data WeekTypeStatistic) []WeekDetailTypeStatistic {

	//获取那个星期的时间范围
	t := time.Now()
	timeLocation := t.Location()
	thisWeekStartTime := firstDayOfISOWeek(data.Year, data.Week, timeLocation)
	thisWeekEndTime := thisWeekStartTime.AddDate(0, 0, 7)

	layout := "2006-01-02"
	startTimeString := thisWeekStartTime.Format(layout)
	endTimeString := thisWeekEndTime.Format(layout)
	fmt.Println(startTimeString)
	fmt.Println(endTimeString)

	timeRangeData := this.AccountDb.AccountJoinCategory(userId, data.Type, startTimeString, endTimeString)
	// Dump(timeRangeData)

	//总价格
	var totalPrice int

	result := QueryGroup(timeRangeData, "CategoryId desc", func(list []WeekDetailTypeStatistic) []WeekDetailTypeStatistic {
		sum := QuerySum(QueryColumn(list, "Money").([]int))

		totalPrice += sum.(int)

		list[0].Money = sum.(int)
		return []WeekDetailTypeStatistic{list[0]}
	}).([]WeekDetailTypeStatistic)

	totalPriceFloat, err := strconv.ParseFloat(strconv.Itoa(totalPrice), 32)
	if err != nil {
		panic(err)
	}

	result = QuerySelect(result, func(singleData WeekDetailTypeStatistic) WeekDetailTypeStatistic {

		categoryMoneyFloat, err := strconv.ParseFloat(strconv.Itoa(singleData.Money), 32)

		categoryMoneyScale := categoryMoneyFloat / totalPriceFloat * 100

		categoryMoneyScaleString := fmt.Sprintf("%.2f", categoryMoneyScale)

		if err != nil {
			panic(err)
		}

		singleData.Precent = categoryMoneyScaleString
		return singleData
	}).([]WeekDetailTypeStatistic)

	return result

}

func (this *AccountAoModel) GetWeekCardStatistic(userId int) []WeekCardStatistic {
	accountJoinCard := this.AccountDb.GetWeekCardStatistic(userId)
	// fmt.Printf("%+v", accountJoinCard)

	yearWeekCardStatistic := QuerySelect(accountJoinCard, func(singleData WeekCardStatistic) WeekCardStatistic {
		t, err := time.Parse("2006-01-02 15:05:04", singleData.CreateTime)
		if err != nil {
			panic(err)
		}
		year, week := t.ISOWeek()
		return WeekCardStatistic{
			CardId:       singleData.CardId,
			CardName:     singleData.CardName,
			CardMoney:    singleData.CardMoney,
			AccountMoney: singleData.AccountMoney,
			Type:         singleData.Type,
			Week:         week,
			Year:         year,
		}
	}).([]WeekCardStatistic)

	yearWeekCardStatistic = QueryGroup(yearWeekCardStatistic, "Year desc,Week desc,CardId desc", func(list []WeekCardStatistic) []WeekCardStatistic {

		// Dump(list)

		var cardTotalPice int

		listLanght := len(list)

		for i := 0; i < listLanght; i++ {
			switch list[i].Type {
			case accountEnum.INCOME, accountEnum.TRANSFER_INCOME, accountEnum.ACCOUNT_RECEIVABLE:
				cardTotalPice += list[i].AccountMoney
			case accountEnum.SPENDING, accountEnum.TRANSFER_SPENDING, accountEnum.ACCOUNTS_PAYABLE:
				cardTotalPice -= list[i].AccountMoney
			default:
				panic(errors.New("不存在该类型"))
			}
		}

		list[0].Money = cardTotalPice

		// Dump([]WeekCardStatistic{list[0]})
		return []WeekCardStatistic{list[0]}

	}).([]WeekCardStatistic)

	// Dump(yearWeekCardStatistic)

	var where Card
	var limit CommonPage
	cardSearch := this.CardAo.Search(userId, where, limit)
	cardData := cardSearch.Data

	initMoney := map[int]int{}

	cardDataLanght := len(cardData)
	for i := 0; i < cardDataLanght; i++ {
		initMoney[cardData[i].CardId] = cardData[i].Money
	}
	// fmt.Printf("%+v", initMoney)

	// Dump(yearWeekCardStatistic)

	result := QueryGroup(yearWeekCardStatistic, "Year asc,Week asc", func(list []WeekCardStatistic) []WeekCardStatistic {

		single := list[0]
		// fmt.Println("list")
		// Dump(list)

		result := QueryLeftJoin(cardData, list, "CardId = CardId", func(left Card, right WeekCardStatistic) WeekCardStatistic {
			// fmt.Println("left", "right")
			// Dump(initMoney[left.CardId])
			// Dump(right)

			initMoney[left.CardId] = initMoney[left.CardId] + right.Money
			// initMoney[left.CardId] += right.AccountMoney

			return WeekCardStatistic{
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

		}).([]WeekCardStatistic)
		// Dump(result)
		return result
	}).([]WeekCardStatistic)

	result = QueryGroup(result, "Year desc,Week desc", func(list []WeekCardStatistic) []WeekCardStatistic {
		return list
	}).([]WeekCardStatistic)

	return result

}

func (this *AccountAoModel) GetWeekDetailCardStatistic(userId int, data WeekCardStatistic) []WeekDetailTypeStatistic {

	accountEnumNames := accountEnum.Names()

	// fmt.Printf("%+v", accountEnumDatas)

	//获取那个星期的时间范围
	t := time.Now()
	timeLocation := t.Location()
	thisWeekStartTime := firstDayOfISOWeek(data.Year, data.Week, timeLocation)
	thisWeekEndTime := thisWeekStartTime.AddDate(0, 0, 7)

	// startMouth := string(thisWeekStartTime.Month())
	layout := "2006-01-02"
	startTimeString := thisWeekStartTime.Format(layout)
	endTimeString := thisWeekEndTime.Format(layout)
	fmt.Println(startTimeString)
	fmt.Println(endTimeString)

	timeRangeData := this.AccountDb.GetWeekDetailCardStatistic(userId, data.CardId, startTimeString, endTimeString)

	// fmt.Printf("%+v", timeRangeData)
	// Dump(timeRangeData)

	// //总价格
	var totalPrice int

	result := QueryGroup(timeRangeData, "Type asc", func(list []WeekDetailTypeStatistic) []WeekDetailTypeStatistic {
		sum := QuerySum(QueryColumn(list, "Money").([]int))

		totalPrice += sum.(int)

		list[0].Money = sum.(int)
		return []WeekDetailTypeStatistic{list[0]}
	}).([]WeekDetailTypeStatistic)

	totalPriceFloat, err := strconv.ParseFloat(strconv.Itoa(totalPrice), 32)
	if err != nil {
		panic(err)
	}

	result = QuerySelect(result, func(singleData WeekDetailTypeStatistic) WeekDetailTypeStatistic {

		categoryMoneyFloat, err := strconv.ParseFloat(strconv.Itoa(singleData.Money), 32)

		categoryMoneyScale := categoryMoneyFloat / totalPriceFloat * 100

		categoryMoneyScaleString := fmt.Sprintf("%.2f", categoryMoneyScale)

		if err != nil {
			panic(err)
		}

		singleData.TypeName = accountEnumNames[strconv.Itoa(singleData.Type)]
		singleData.Precent = categoryMoneyScaleString
		return singleData
	}).([]WeekDetailTypeStatistic)

	return result

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
