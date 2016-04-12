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

	// wheres := Account{
	// 	UserId:     userId,
	// 	Name:       data.Name,
	// 	Remark:     data.Remark,
	// 	CategoryId: data.CategoryId,
	// 	CardId:     data.CardId,
	// 	Type:       data.Type,
	// }

	where.UserId = userId

	// fmt.Printf("%+v", wheres)

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

	// accounts := Account{
	// 	UserId:     userId,
	// 	Name:       account.Name,
	// 	Money:      account.Money,
	// 	Remark:     account.Remark,
	// 	CategoryId: account.CategoryId,
	// 	CardId:     account.CardId,
	// 	Type:       account.Type,
	// }

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

// func (this *AccountAoModel) GetWeekTypeStatistic(userId int) []WeekTypeStatistic {

// 	// var classify map[int]map[int][]WeekTypeStatistic
// 	classify := map[int]map[int]map[int][]WeekTypeStatistic{}
// 	weekState := this.AccountDb.AccountJoinCard(userId)
// 	accountEnums := accountEnum.Names()
// 	weekStateLangut := len(weekState)

// 	for i := 0; i < weekStateLangut; i++ {

// 		t, err := time.Parse("2006-01-02 15:05:04", weekState[i].CreateTime)
// 		if err != nil {
// 			panic(err)
// 		}

// 		year, week := t.ISOWeek()
// 		weekState[i].Name = strconv.Itoa(year) + "年" + strconv.Itoa(week) + "周"
// 		weekState[i].Year = year
// 		weekState[i].Week = week

// 		weekState[i].TypeName = accountEnums[strconv.Itoa(weekState[i].Type)]

// 		// fmt.Printf("%+v", weekState[i])
// 		// fmt.Printf("%+v", classify[year][week])

// 		// fmt.Println("\nstart\n")

// 		_, ok := classify[year]
// 		if ok != true {
// 			classify[year] = map[int]map[int][]WeekTypeStatistic{}
// 		}

// 		_, oks := classify[year][week]
// 		if oks != true {
// 			classify[year][week] = map[int][]WeekTypeStatistic{}
// 		}

// 		theType := weekState[i].Type
// 		classify[year][week][theType] = append(classify[year][week][theType], weekState[i])

// 		// fmt.Println(append(classify[year][week], weekState[i]))

// 		// fmt.Println("\nend\n")

// 	}

// 	// this.Log.Debug(classify)

// 	// fmt.Printf("%+v", classify)

// 	for year, weeks := range classify {
// 		for week, Types := range weeks {
// 			for TypeNum, Type := range Types {
// 				for _, single := range Type {
// 					fmt.Println("\nyear\n")
// 					fmt.Printf("%+v", year)
// 					fmt.Println("\nweek\n")
// 					fmt.Printf("%+v", week)
// 					fmt.Println("\nTypeNum\n")
// 					fmt.Printf("%+v", TypeNum)
// 					fmt.Println("\nsingle\n")
// 					fmt.Printf("%+v", single)
// 				}

// 			}

// 		}
// 	}

// 	return weekState

// }

func (this *AccountAoModel) GetWeekTypeStatistic(userId int) []WeekTypeStatistic {
	var where Account
	var limit CommonPage
	accountSearchData := this.Search(userId, where, limit)
	accountEnums := accountEnum.Names()
	accountData := accountSearchData.Data

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
		_, ok := classify[year]
		if ok != true {
			classify[year] = map[int]map[int][]WeekTypeStatistic{}
		}

		_, oks := classify[year][week]
		if oks != true {
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

	// startMouth := string(thisWeekStartTime.Month())
	startTimeString := strconv.Itoa(thisWeekStartTime.Year()) + "-" + mouthNum(thisWeekStartTime.Month().String()) + "-" + strconv.Itoa(thisWeekStartTime.Day())
	endTimeString := strconv.Itoa(thisWeekEndTime.Year()) + "-" + mouthNum(thisWeekEndTime.Month().String()) + "-" + strconv.Itoa(thisWeekEndTime.Day())
	fmt.Println(startTimeString)
	fmt.Println(endTimeString)

	timeRangeData := this.AccountDb.TimerangeOfData(userId, data.Type, startTimeString, endTimeString)

	//分类摆放
	classify := map[int][]WeekDetailTypeStatistic{}

	//统计总价格
	var totalPrice int

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

func mouthNum(mouth string) string {
	switch mouth {
	case "January":
		return "01"
	case "February":
		return "02"
	case "March":
		return "03"
	case "April":
		return "04"
	case "May":
		return "05"
	case "June":
		return "06"
	case "July":
		return "07"
	case "August":
		return "08"
	case "September":
		return "09"
	case "October":
		return "10"
	case "November":
		return "11"
	case "December":
		return "12"
	default:
		panic(errors.New("没有这个月份"))
	}
}
