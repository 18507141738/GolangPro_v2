package config

import (
	"Artifice_V2.0.0/controllers"
	"Artifice_V2.0.0/models"
	"Artifice_V2.0.0/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"log"
	"strings"
)

type AlarmSetPage struct {
	beego.Controller
}

func (this *AlarmSetPage) Get() {
	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	this.Data["bmk"] = "alarm_menu"
	this.TplName = "config/alarmSet.html"
}

type SelAlarmDetailController struct {
	beego.Controller
}

func (this *SelAlarmDetailController) Post() {
	//limits, pages :=
	//this.GetString("limit"),
	//	this.GetString("page")
	//start_page, err := strconv.Atoi(pages)
	//limitss, err := strconv.Atoi(limits)
	//start_pages := (start_page - 1) * limitss
	o := controllers.O
	var maps []orm.Params

	num, err := o.QueryTable(new(models.AlarmDetail)).Values(&maps)
	//num,err := o.Raw("select * from ss_alarm_detail").Values(&maps)
	if err != nil {
		//this.Data["json"] = &util.ResultB{0, 0, nil}
		//this.ServeJSON()
		//return
		log.Print("告警文案查询失败")
	}

	this.Data["json"] = &util.ResultB{0, num, &maps}
	this.ServeJSON()

}

type AddAlarmDetailController struct {
	beego.Controller
}

func (this *AddAlarmDetailController) Post() {
	funcType, detail := this.GetString("type"), this.GetString("detail")

	if len(funcType) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "功能类型不能为空"}
		this.ServeJSON()
		return
	}

	if len(detail) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "告警文案不能为空"}
		this.ServeJSON()
		return
	}

	o := controllers.O

	uid := uuid.Must(uuid.NewV4()).String()
	uid = strings.Replace(uid, "-", "", -1) // 生成SessionID

	var maps []models.AlarmDetail
	_, err := o.QueryTable(new(models.AlarmDetail)).Filter("type", funcType).All(&maps)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "保存失败"}
		this.ServeJSON()
		return
	}
	if len(maps) > 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "该功能告警文案已配置"}
		this.ServeJSON()
		return
	}

	_, err = o.Raw("insert into ss_alarm_detail(id, type, detail) values(?,?,?)", uid, funcType, detail).Exec()
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "保存失败"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "添加成功"}
	this.ServeJSON()

}

type EditAlarmDetailController struct {
	beego.Controller
}

func (this *EditAlarmDetailController) Post() {
	id, funcType, detail := this.GetString("id"), this.GetString("type"), this.GetString("detail")

	if len(id) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "参数错误"}
		this.ServeJSON()
		return
	}

	if len(funcType) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "功能类型不能为空"}
		this.ServeJSON()
		return
	}

	if len(detail) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "告警文案不能为空"}
		this.ServeJSON()
		return
	}

	o := controllers.O
	var alarmDetail []orm.Params

	num, err := o.Raw("select * from ss_alarm_detail where id = ?", id).Values(&alarmDetail)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "修改失败"}
		this.ServeJSON()
		return
	}

	if num > 0 && alarmDetail[0]["type"] == funcType {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "修改失败，该功能告警文案已配置"}
		this.ServeJSON()
		return
	}
	_, err = o.Raw("update ss_alarm_detail type=?, detail=? where id=?", funcType, detail, id).Exec()
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "修改失败"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "修改成功"}
	this.ServeJSON()
}

type DelAlarmDetailController struct {
	beego.Controller
}

func (this *DelAlarmDetailController) Post() {
	ids := this.GetString("ids")

	idArr := strings.Split(ids, ",")
	if len(idArr) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "参数错误"}
		this.ServeJSON()
		return
	}

	o := controllers.O
	_, err := o.QueryTable(new(models.AlarmDetail)).Filter("id__in", idArr).Delete()
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "删除失败"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "删除成功"}
	this.ServeJSON()
}
