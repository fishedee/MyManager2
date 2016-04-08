package routers

import (
	. "mymanager/controllers"
)

func init() {
	InitRoute("/login", &LoginController{})
	InitRoute("/user", &UserController{})
}
