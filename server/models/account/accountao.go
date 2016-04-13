package account

import (
	. "github.com/fishedee/language"
	. "mymanager/models/common"
	// "crypto/sha1"
	"fmt"
	"strconv"
	"time"
	// "io"
	"errors"
	// "math"
)

type AccountAoModel struct {
	BaseModel
	AccountDb AccountDbModel
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
	accountEnums := accountEnum.Names()
	accountData := accountSearchData.Data

	//分类摆放 year/week/type
	classify := map[int]map[int]map[int][]WeekTypeStatistic{}

	accountLangut := len(accountData)

	for i := 0; i < accountLangut; i++ {

		singleData := accountData[i]

		//获取多小年的第几周
		t, err := time.Parse("2006-01-02 15:05:04", singleData.CreateTime)
		if err != nil {
			panic(err)
		}

		year, week := t.ISOWeek()

		//初始化map
		_, isExist := classify[year]
		if isExist != true {
			classify[year] = map[int]map[int][]WeekTypeStatistic{}
		}

		_, isExists := classify[year][week]
		if isExists != true {
			classify[year][week] = map[int][]WeekTypeStatistic{}
		}

		//按时间和类型组合
		theType := singleData.Type
		classify[year][week][theType] = append(classify[year][week][theType], WeekTypeStatistic{
			Money: singleData.Money,
			Type:  singleData.Type,
			Week:  week,
			Year:  year,
		})
	}

	var result []WeekTypeStatistic

	for _, weeks := range classify {
		for _, Types := range weeks {
			for _, Type := range Types {
				var tempMoney int
				for _, single := range Type {
					// fmt.Println("\nyear")
					// fmt.Printf("%+v", year)
					// fmt.Println("\nweek")
					// fmt.Printf("%+v", week)
					// fmt.Println("\nTypeNum")
					// fmt.Printf("%+v", TypeNum)
					// fmt.Println("\nsingle")
					// fmt.Printf("%+v", single)
					// fmt.Println("\n")

					tempMoney += single.Money
				}
				fmt.Println("tempMoney", tempMoney)
				result = append(result, WeekTypeStatistic{
					Money:    tempMoney,
					Type:     Type[0].Type,
					TypeName: accountEnums[strconv.Itoa(Type[0].Type)],
					Name:     strconv.Itoa(Type[0].Year) + "年" + strconv.Itoa(Type[0].Week) + "周",
					Week:     Type[0].Week,
					Year:     Type[0].Year,
				})

			}

		}
	}

	fmt.Printf("%+v", result)

	return result

}
func (this *AccountAoModel) GetWeekDetailTypeStatistic(userId int, data WeekTypeStatistic) []WeekDetailTypeStatistic {

	// var where Account
	// var limit CommonPage
	// accountSearchData := this.Search(userId, where, limit)
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

	//分类摆放
	classify := map[int][]WeekDetailTypeStatistic{}

	//统计总价格
	var totalPrice int

	//顺便按CategoryId分类
	for _, v := range timeRangeData {
		totalPrice += v.Money
		classify[v.CategoryId] = append(classify[v.CategoryId], v)
	}

	var result []WeekDetailTypeStatistic

	totalPriceFloat, err := strconv.ParseFloat(strconv.Itoa(totalPrice), 32)

	if err != nil {
		panic(err)
	}

	for categoryId, categoryData := range classify {
		var categoryMoney int
		for _, singleData := range categoryData {
			categoryMoney += singleData.Money
		}

		categoryMoneyFloat, err := strconv.ParseFloat(strconv.Itoa(categoryMoney), 32)

		if err != nil {
			panic(err)
		}

		// fmt.Println()

		categoryMoneyScale := categoryMoneyFloat / totalPriceFloat * 100

		categoryMoneyScaleString := fmt.Sprintf("%.2f", categoryMoneyScale)

		fmt.Println(categoryMoneyScaleString)

		result = append(result, WeekDetailTypeStatistic{
			CategoryId:   categoryId,
			CategoryName: categoryData[0].CategoryName,
			Money:        categoryMoney,
			Precent:      categoryMoneyScaleString,
		})
	}

	fmt.Printf("%+v", result)

	return result

}

func (this *AccountAoModel) GetWeekCardStatistic(userId int) []WeekCardStatistic {
	accountJoinCard := this.AccountDb.GetWeekCardStatistic(userId)
	fmt.Printf("%+v", accountJoinCard)
	//分类摆放 year/week/cardId/type
	classify := map[int]map[int]map[int]map[int][]WeekCardStatistic{}

	accountJoinCardLangut := len(accountJoinCard)

	for i := 0; i < accountJoinCardLangut; i++ {

		singleData := accountJoinCard[i]
		singleCardId := singleData.CardId
		singleType := singleData.Type

		//获取多小年的第几周
		t, err := time.Parse("2006-01-02 15:05:04", singleData.CreateTime)
		if err != nil {
			panic(err)
		}

		year, week := t.ISOWeek()

		//初始化map
		_, isExist := classify[year]
		if isExist != true {
			classify[year] = map[int]map[int]map[int][]WeekCardStatistic{}
		}

		_, isExists := classify[year][week]
		if isExists != true {
			classify[year][week] = map[int]map[int][]WeekCardStatistic{}
		}

		_, isExistss := classify[year][week][singleCardId]
		if isExistss != true {
			classify[year][week][singleCardId] = map[int][]WeekCardStatistic{}
		}

		//按时间和类型组合
		classify[year][week][singleCardId][singleType] = append(classify[year][week][singleCardId][singleType], WeekCardStatistic{
			CardId:       singleData.CardId,
			CardName:     singleData.CardName,
			CardMoney:    singleData.CardMoney,
			AccountMoney: singleData.AccountMoney,
			Type:         singleData.Type,
			Week:         week,
			Year:         year,
		})
	}
	fmt.Println("\n")
	fmt.Printf("%+v", classify)

	var result []WeekCardStatistic

	for _, weeks := range classify {
		for _, cardId := range weeks {
			for _, Types := range cardId {

				var tempMoney int
				var tempType int

				for typeNum, Type := range Types {

					for _, single := range Type {
						// fmt.Println("\nyear")
						// fmt.Printf("%+v", year)
						// fmt.Println("\nweek")
						// fmt.Printf("%+v", week)
						// fmt.Println("\nTypeNum")
						// fmt.Printf("%+v", TypeNum)
						// fmt.Println("\nsingle")
						// fmt.Printf("%+v", single)
						// fmt.Println("\n")

						// tempMoney += single.Money
						fmt.Println("\n")
						fmt.Printf("%+v", single)

						switch single.Type {
						case 1, 3, 5:
							tempMoney += single.AccountMoney
						case 2, 4, 6:
							tempMoney -= single.AccountMoney
						default:
							panic(errors.New("不存在该类型"))
						}

						tempType = typeNum
					}
					fmt.Println("tempMoney", tempMoney)

				}
				fmt.Println("\n")
				fmt.Println("Types")
				fmt.Printf("%+v", Types)
				result = append(result, WeekCardStatistic{
					Money:    Types[tempType][0].CardMoney + tempMoney,
					Name:     strconv.Itoa(Types[tempType][0].Year) + "年" + strconv.Itoa(Types[tempType][0].Week) + "周",
					CardId:   Types[tempType][0].CardId,
					CardName: Types[tempType][0].CardName,
					Week:     Types[tempType][0].Week,
					Year:     Types[tempType][0].Year,
				})
			}

		}
	}

	fmt.Printf("%+v", result)

	return result
}

func (this *AccountAoModel) GetWeekDetailCardStatistic(userId int, data WeekCardStatistic) []WeekDetailTypeStatistic {

	accountEnums := accountEnum.Names()

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
	fmt.Printf("%+v", timeRangeData)

	//分类摆放
	classify := map[int][]WeekDetailTypeStatistic{}

	//统计总价格
	var totalPrice int

	//顺便按CategoryId分类
	for _, v := range timeRangeData {
		fmt.Println("\n")
		totalPrice += v.Money
		fmt.Printf("%+v", v)
		classify[v.Type] = append(classify[v.Type], v)
	}

	var result []WeekDetailTypeStatistic

	totalPriceFloat, err := strconv.ParseFloat(strconv.Itoa(totalPrice), 32)

	if err != nil {
		panic(err)
	}

	for Type, categoryData := range classify {
		var categoryMoney int
		for _, singleData := range categoryData {
			categoryMoney += singleData.Money
		}

		categoryMoneyFloat, err := strconv.ParseFloat(strconv.Itoa(categoryMoney), 32)

		if err != nil {
			panic(err)
		}

		// fmt.Println()

		categoryMoneyScale := categoryMoneyFloat / totalPriceFloat * 100

		categoryMoneyScaleString := fmt.Sprintf("%.2f", categoryMoneyScale)

		fmt.Println(categoryMoneyScaleString)

		result = append(result, WeekDetailTypeStatistic{
			Type:     Type,
			TypeName: accountEnums[strconv.Itoa(Type)],
			// CategoryName: categoryData[0].CategoryName,
			Money:   categoryMoney,
			Precent: categoryMoneyScaleString,
		})
	}

	fmt.Printf("%+v", result)

	return result
	// return 0

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
