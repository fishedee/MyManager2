package account

import (
	. "github.com/fishedee/language"
)

var accountEnum struct {
	EnumStruct
	INCOME             int `enum:"1,收入"`
	SPENDING           int `enum:"2,支出"`
	TRANSFER_INCOME    int `enum:"3,转账收入"`
	TRANSFER_SPENDING  int `enum:"4,转账支出"`
	ACCOUNT_RECEIVABLE int `enum:"5,借还账收入"`
	ACCOUNTS_PAYABLE   int `enum:"5,借还账支出"`
}

func init() {
	InitEnumStruct(&accountEnum)
}
