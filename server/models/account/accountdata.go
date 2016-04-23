package account

import (
	"time"
)

type Account struct {
	AccountId  int `xorm:"autoincr"`
	UserId     int
	Name       string
	Money      int
	Remark     string
	CategoryId int
	CardId     int
	Type       int
	CreateTime time.Time `xorm:"created"`
	ModifyTime time.Time `xorm:"updated"`
}

type Accounts struct {
	Count int
	Data  []Account
}

type WeekStatistic struct {
	CardId     int
	CardName   string
	Money      int
	Name       string
	Type       int
	TypeName   string
	Week       int
	Year       int
	CreateTime string
	ModifyTime string
}

type WeekDetailStatistic struct {
	CategoryId   int
	CategoryName string
	Type         int
	TypeName     string
	Money        int
	Precent      string
}
