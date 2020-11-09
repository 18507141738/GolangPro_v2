package config

import (
	"Artifice_V2.0.0/controllers"
	"Artifice_V2.0.0/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type SysPageHandler struct {
	beego.Controller
}

func (this *SysPageHandler) Get() {
	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	this.Data["bmk"] = "system"
	this.TplName = "config/system.html"
}

// 疫情防控系统查询系统配置
type SelSysTemInfoHandler struct {
	beego.Controller
}

func (this *SelSysTemInfoHandler) Post() {
	var maps []orm.Params
	o := controllers.O
	num, err := o.Raw("select * from ss_system").Values(&maps)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"ret": 0, "reason": "数据操作失败!,再来一次"}
	} else {
		mystruct := &util.ResultB{0, num, &maps}
		this.Data["json"] = mystruct
	}
	this.ServeJSON()
}

// 疫情防控平台更新系统配置
type UpdateConfigHandler struct {
	beego.Controller
}

func (this *UpdateConfigHandler) Post() {
	websitename, mediaserverIP, spare_switch := this.GetString("websitename"), this.GetString("mediaserverIP"), this.GetString("spare_switch")
	o := controllers.O
	_, err := o.Raw("update ss_system set websitename=?,mediaserverIP=?,spare_switch=?", websitename, mediaserverIP, spare_switch).Exec()
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败!,再来一次"}
	} else {
		this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "操作成功"}
	}
	this.ServeJSON()
}
