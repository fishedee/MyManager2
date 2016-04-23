package card

import (
	. "github.com/fishedee/language"
)

var CardQueueEnum struct {
	EnumStructString
	EVENT_DEL string `enum:"/card/_del,类目被删除"`
}

func init() {
	InitEnumStructString(&CardQueueEnum)
}
