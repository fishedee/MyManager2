package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	. "github.com/fishedee/web"
	"reflect"
)

type BaseController struct {
	BeegoValidateController
}

func InitRoute(namespace string, target beego.ControllerInterface) {
	InitBeegoVaildateControllerRoute(namespace, target)
}

type baseControllerResult struct {
	Code int
	Data interface{}
	Msg  string
}

func (this *BaseController) jsonRender(result baseControllerResult) {
	resultString, err := EncodeJson(result)
	if err != nil {
		panic(err)
	}
	this.Ctx.WriteString(string(resultString))
}

func (this *BaseController) redirectRender(result baseControllerResult) {
	//FIXME 没有做更多的容错尝试
	if result.Code == 0 {
		url := result.Data.(string)
		this.Ctx.Redirect(302, url)
	} else {
		this.Ctx.WriteString("跳转不成功 " + result.Msg)
	}
}

func (this *BaseController) excelRender(result baseControllerResult) {
	//获取excel的导出配置
	var excelArgs struct {
		ViewTitle  string `validate:"_viewTitle"`
		ViewFormat string `validate:"_viewFormat"`
	}
	this.Check(&excelArgs)

	excelTitle := excelArgs.ViewTitle
	excelFormat, err := DecodeUrl(excelArgs.ViewFormat)
	if err != nil {
		panic(err)
	}
	jsonFormat := map[string]string{}
	err = json.Unmarshal([]byte(excelFormat), &jsonFormat)
	if err != nil {
		panic(err)
	}
	jsonData := reflect.ValueOf(result.Data).FieldByName("Data").Interface()
	tableData := ArrayColumnTable(jsonFormat, jsonData)

	//写入数据
	resultByte, err := EncodeXlsx(tableData)
	if err != nil {
		panic(err)
	}
	this.WriteMimeHeader("xlsx", excelTitle)
	this.Write(resultByte)
}

func (this *BaseController) AutoRender(returnValue interface{}, viewname string) {
	result := baseControllerResult{}
	resultError, ok := returnValue.(Exception)
	if ok {
		//带错误码的error
		result.Code = resultError.GetCode()
		result.Msg = resultError.GetMessage()
		result.Data = nil
	} else {
		//正常返回
		result.Code = 0
		result.Data = returnValue
		result.Msg = ""
	}
	this.Ctx.Output.Header("Cache-Control", "private, no-store, no-cache, must-revalidate, max-age=0")
	this.Ctx.Output.Header("Cache-Control", "post-check=0, pre-check=0")
	this.Ctx.Output.Header("Pragma", "no-cache")
	this.Ctx.Output.Header("Access-Control-Allow-Origin", this.Ctx.Input.Header("Origin"))
	this.Ctx.Output.Header("Access-Control-Allow-Credentials", "true")

	var inputViewName struct {
		View string `validate:"_view"`
	}
	this.Check(&inputViewName)

	if inputViewName.View == "excel" {
		this.excelRender(result)
	} else if viewname == "json" {
		this.jsonRender(result)
	} else if viewname == "redirect" {
		this.redirectRender(result)
	} else {
		panic("不合法的viewName " + viewname)
	}
}
