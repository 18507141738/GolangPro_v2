package controllers

import (
	"Artifice_V2.0.0/models"
	"Artifice_V2.0.0/util"
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	"time"
)

//投屏登录页
type TSLoginPage struct {
	beego.Controller
}

func (this *TSLoginPage) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["version"] = GetVer()
	this.Data["updateTime"] = GetUpdateTime()
	this.Data["buildTime"] = GetBuildTime()
	this.TplName = "platform/throwScreen/login.html"
}

//投屏页
type TSHomePage struct {
	beego.Controller
}

func (this *TSHomePage) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.TplName = "platform/throwScreen/home.html"
}

//投屏登录
type TSLoginCheck struct {
	beego.Controller
}

func (this *TSLoginCheck) Post() {
	username, password := this.GetString("Username"), this.GetString("Password")
	h := md5.New()
	h.Write([]byte(password))
	md5pro := hex.EncodeToString(h.Sum(nil))
	o := O
	var maps []orm.Params
	num, err := o.Raw("select * from ss_user where admin_user=?", username).Values(&maps)
	if err == nil && num > 0 {
		pwd := maps[0]["admin_password"]
		if pwd != md5pro {
			this.Data["json"] = map[string]interface{}{"code": -1, "msg": "密码错误"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "msg": "登录成功"}
			this.Ctx.SetCookie("TSUser_id", maps[0]["admin_id"].(string))
		}
	} else {
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": "没有此用户"}
	}
	this.ServeJSON()
}

//单个算法内容统计
type TSFuncCountController struct {
	beego.Controller
}

func (this *TSFuncCountController) Post() {

	funcType := this.GetString("type")

	o := O
	var maps []orm.Params

	var nowTime = time.Now().Format("2006-01-02") + " 23:59:59"
	var startTime = time.Now().AddDate(0, 0, -6).Format("2006-01-02") + " 00:00:00"

	var sql = "SELECT mm.aa as brief_time,IFNULL(a.num,0) as num " +
		"FROM  " +
		"( SELECT DATE_FORMAT(DATE_SUB(CONVERT(?,datetime),INTERVAL xc-1 DAY),'%m%d') as aa," +
		" DATE_FORMAT( DATE_SUB(CONVERT(?,datetime), INTERVAL xc - 1 DAY), '%Y-%m-%d' ) AS bb " +
		" FROM" +
		"( SELECT @xi:=@xi+1 as xc " +
		"FROM " +
		"(SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6) xc1, " +
		"(SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6) xc2, " +
		"(SELECT @xi:=0) xc0 ) xcxc " +
		"WHERE " +
		"xc<=DATEDIFF(CONVERT(?,datetime),CONVERT(?,datetime))+1 )mm " +
		"LEFT OUTER JOIN " +
		"(SELECT " +
		"DATE_FORMAT(alarm_time, '%m%d') as brief_time," +
		"COUNT(*) AS num " +
		"FROM ss_alarm ss " +
		"INNER JOIN ss_camera sc ON sc.camera_id = ss.camera_id " +
		" INNER JOIN ss_place sp ON sp.place_id = sc.place_id " +
		" INNER JOIN ss_organize so ON so.id = sp.organize_id " +
		"WHERE " +
		"alarm_time>=? and alarm_time<=? " +
		"AND alarm_type = ? " +
		"GROUP BY brief_time)a " +
		"ON a.brief_time = mm.aa ORDER BY mm.bb asc"

	_, err := o.Raw(sql, nowTime, nowTime, nowTime, startTime, startTime, nowTime, funcType).Values(&maps)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "查询信息异常"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = &util.ResultB{1, 0, &maps}
	this.ServeJSON()
}

//展示摄像头
type TSCamearController struct {
	beego.Controller
}

func (this *TSCamearController) Post() {

	o := O
	var maps []orm.Params

	num, err := o.Raw("select * from ss_camera where screen='on';").Values(&maps)
	if err != nil {
		log.Println("获取展示摄像头异常", err)
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "获取展示摄像头异常"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = &util.ResultB{1, num, &maps}
	this.ServeJSON()

}

//所有未处理告警
type TSAlarmUntreatedController struct {
	beego.Controller
}

func (this *TSAlarmUntreatedController) Post() {
	o := O
	var alarms []models.Alarm

	num, err := o.QueryTable(new(models.Alarm)).Filter("alarm_status", "0").OrderBy("-alarm_time").Limit(100).RelatedSel(4).All(&alarms)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "获取告警信息异常"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = map[string]interface{}{"Code": 1, "num": num, "Reason": &alarms}
	this.ServeJSON()
}

//所有告警统计
type TSAllFuncAlarmController struct {
	beego.Controller
}

func (this *TSAllFuncAlarmController) Post() {
	o := O
	var maps []orm.Params

	var nowTime = time.Now().Format("2006-01-02") + " 23:59:59"
	var startTime = time.Now().AddDate(0, 0, -6).Format("2006-01-02") + " 00:00:00"

	var sql = "SELECT " +
		"COUNT(alarm_type='smoke' or null) AS num_smoke," +
		"COUNT(alarm_type='fire' or null) AS num_fire," +
		"COUNT(alarm_type='cloths' or null) AS num_cloths," +
		"COUNT(alarm_type='boundary' or null) AS num_boundary," +
		"COUNT(alarm_type='queue_count' or null) AS num_leave," +
		"COUNT(alarm_type='sleep_count' or null) AS num_sleep," +
		"COUNT(alarm_type='leakage' or null) AS num_leakage" +
		" FROM ss_alarm " +
		"WHERE " +
		"alarm_type IN ('smoke','fire','cloths','boundary','queue_count','sleep_count','leakage') " +
		"AND alarm_time>=? and alarm_time<=?"

	num, err := o.Raw(sql, startTime, nowTime).Values(&maps)
	if err != nil {
		log.Println("获取告警统计异常", err)
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "获取告警统计异常"}
		this.ServeJSON()
		return
	}

	this.Data["json"] = &util.ResultB{1, num, &maps}
	this.ServeJSON()

}
