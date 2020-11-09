package filter

import (
	"Artifice_V2.0.0/controllers"
	"github.com/astaxie/beego"
)

type FilterStationPage struct {
	beego.Controller
}

func (this *FilterStationPage) Get() {
	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	this.Data["version"] = controllers.GetVer()
	this.Data["updateTime"] = controllers.GetUpdateTime()
	this.Data["buildTime"] = controllers.GetBuildTime()
	this.Data["sysfilter"] = controllers.SelSystemFilter()
	this.Data["bmk"] = "filterStation"
	this.TplName = "filter/admin/station.html"

}

type FilterSetStatus struct {
	beego.Controller
}

func (this *FilterSetStatus) Post() {
	status := this.GetString("filterstatus")
	o := controllers.O
	_, err := o.Raw("update ss_system set filterstatus=?", status).Exec()
	if err != nil {
		controllers.LogsError("跟新系统过滤状态异常", err)
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": "操作失败"}
	} else {
		if status == "0" {
			controllers.FilterStatus = false
		} else {
			controllers.FilterStatus = true
		}
		this.Data["json"] = map[string]interface{}{"code": 1, "msg": "操作成功"}
	}
	this.ServeJSON()
}
