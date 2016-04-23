package main

import (
	"github.com/astaxie/beego"
	_ "mymanager/routers"
)

//go:generate fishgen ^./models/.*(ao|db)\.go$
func main() {
	beego.Run()
}
