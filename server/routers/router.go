package routers

import (
	. "mymanager/controllers"
)

func init() {
	InitRoute("/login", &LoginController{})
	InitRoute("/user", &UserController{})
	InitRoute("/category", &CategoryController{})
	InitRoute("/card", &CardController{})
}
