package category

import (
	. "github.com/fishedee/language"
)

var CategoryQueueEnum struct {
	EnumStructString
	EVENT_DEL string `enum:"/category/_del,类目被删除"`
}

func init() {
	InitEnumStructString(&CategoryQueueEnum)
}
