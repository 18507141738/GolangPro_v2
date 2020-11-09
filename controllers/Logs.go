package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func LogsInfo(f interface{}, v ...interface{}) {
	if beego.AppConfig.String("openlogs") == "0" {
		logs.Info(f, v)
	}
}

func LogsError(f interface{}, v ...interface{}) {
	if beego.AppConfig.String("openlogs") == "0" {
		logs.Info(f, v)
	}
}

func LogsAlert(f interface{}, v ...interface{}) {
	if beego.AppConfig.String("openlogs") == "0" {
		logs.Alert(f, v)
	}
}
