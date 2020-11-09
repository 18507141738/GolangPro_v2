package controllers

import (
	"Artifice_V2.0.0/models"
	"Artifice_V2.0.0/util"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
)

// 周界
type BoundaryConfigPageHandler struct {
	beego.Controller
}

func (this *BoundaryConfigPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/boundaryConfig.html"
	}
}

type BoundaryConfigOCXPageHandler struct {
	beego.Controller
}

func (this *BoundaryConfigOCXPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/ie/boundaryConfig.html"
	}
}

type BoundaryConfRTMPPageHandler struct {
	beego.Controller
}

func (this *BoundaryConfRTMPPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/ruleconfig/rtmp/boundaryConfig.html"
	}
}

// 着装
type ClothsConfigPageHandler struct {
	beego.Controller
}

func (this *ClothsConfigPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/clothsConfig.html"
	}
}

type ClothsConfigOCXPageHandler struct {
	beego.Controller
}

func (this *ClothsConfigOCXPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/ie/clothsConfig.html"
	}
}

type ClothsConfigRTMPPageHandler struct {
	beego.Controller
}

func (this *ClothsConfigRTMPPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/ruleconfig/rtmp/clothsConfig.html"
	}
}

// 睡岗
type PlatformSleepCountPageHandler struct {
	beego.Controller
}

func (this *PlatformSleepCountPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/sleep_count_Config.html"
	}
}

type PlatformSleepCountOCXPageHandler struct {
	beego.Controller
}

func (this *PlatformSleepCountOCXPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/ie/sleep_count_Config.html"
	}
}

// 睡岗
type PlatformSleepCountRTMPPageHandler struct {
	beego.Controller
}

func (this *PlatformSleepCountRTMPPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/ruleconfig/rtmp/sleep_count_Config.html"
	}
}

// 离岗
type PlatformLeaveCountPageHandler struct {
	beego.Controller
}

func (this *PlatformLeaveCountPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/leave_count_Config.html"
	}
}

type PlatformLeaveCountOCXPageHandler struct {
	beego.Controller
}

func (this *PlatformLeaveCountOCXPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/ie/leave_count_Config.html"
	}
}

type PlatformLeaveCountRTMPPageHandler struct {
	beego.Controller
}

func (this *PlatformLeaveCountRTMPPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/ruleconfig/rtmp/leave_count_Config.html"
	}
}

// 烟雾
type PlatformSmokepagetPageHandler struct {
	beego.Controller
}

func (this *PlatformSmokepagetPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/smokeConfig.html"
	}
}

type PlatformSmokepagetOCXPageHandler struct {
	beego.Controller
}

func (this *PlatformSmokepagetOCXPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/ie/smokeConfig.html"
	}
}

type PlatformSmokepagetRTMPPageHandler struct {
	beego.Controller
}

func (this *PlatformSmokepagetRTMPPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/ruleconfig/rtmp/smokeConfig.html"
	}
}

// 火焰
type PlatformFireworkpagePageHandler struct {
	beego.Controller
}

func (this *PlatformFireworkpagePageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/fireConfig.html"
	}
}

type PlatformFireworkpageOCXPageHandler struct {
	beego.Controller
}

func (this *PlatformFireworkpageOCXPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/ie/fireConfig.html"
	}
}

type PlatformFireworkpageRTMPPageHandler struct {
	beego.Controller
}

func (this *PlatformFireworkpageRTMPPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/ruleconfig/rtmp/fireConfig.html"
	}
}

// 泄露
type PlatformLeakagePageHandler struct {
	beego.Controller
}

func (this *PlatformLeakagePageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/leakConfig.html"
	}
}

type PlatformLeakageOCXPageHandler struct {
	beego.Controller
}

func (this *PlatformLeakageOCXPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/rtspconfig/ie/leakConfig.html"
	}
}

type PlatformLeakageRTMPPageHandler struct {
	beego.Controller
}

func (this *PlatformLeakageRTMPPageHandler) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["bmk"] = "ssss"
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.Data["organize_name"] = user.Organize.Name
		this.Data["jurisdiction"] = user.Jurisdiction
		this.TplName = "platform/ruleconfig/rtmp/leakConfig.html"
	}
}

//获取配置详情
type GetAlgorithmListHandler struct {
	beego.Controller
}

func (this *GetAlgorithmListHandler) Post() {
	camara_id := this.GetString("camara_id")
	function_type := this.GetString("function_type")
	log.Println("<----------------请求的算法类型是--------------->", function_type, camara_id)
	var maps []orm.Params
	o := O
	num, err := o.Raw("select * from ss_func_camera,ss_camera where ss_func_camera.function_type = ? and ss_func_camera.cameraId=?  and ss_func_camera.cameraId = ss_camera.camera_id", function_type, camara_id).Values(&maps)
	if err != nil {
		log.Println("摄像机列表：", err)
		return
	}
	mystruct := &util.ResultB{0, num, &maps}
	this.Data["json"] = mystruct
	this.ServeJSON()

	return
}

//保存配置
type AddGrpHandler struct {
	beego.Controller
}

func (this *AddGrpHandler) Post() {
	o := O
	if this.GetString("function_type") == "boundary" || this.GetString("function_type") == "cloths" {
		function_type, camera_id, location, threshold, frequency, boundary, DistanceMode, wallOn, wallrate, rateup, ratedown, area_shape, match_boundary0, match_boundary1, match_boundary2, hatcolor, uppercolor, lowercolor, whole, alarm_mode, timepoint, alarm_time_second, video_switch := this.GetString("function_type"),
			this.GetString("id"),
			this.GetString("location"),
			this.GetString("threshold"),
			this.GetString("frequency"),
			this.GetString("boundary"),
			this.GetString("DistanceMode"),
			this.GetString("wallOn"),
			this.GetString("wallrate"),
			this.GetString("rateup"),
			this.GetString("ratedown"),
			this.GetString("area_shape"),
			this.GetString("match_boundary0"),
			this.GetString("match_boundary1"),
			this.GetString("match_boundary2"),
			this.GetString("hatcolor"),
			this.GetString("uppercolor"),
			this.GetString("lowercolor"),
			this.GetString("whole"),
			this.GetString("alarm_mode"),
			this.GetString("timepoint"),
			this.GetString("alarm_time_second"),
			this.GetString("video_switch")
		log.Println("入侵检测 & 着装监控布防设备参数------------------------", function_type, camera_id, location, threshold, frequency, boundary, DistanceMode, wallOn, wallrate, rateup, ratedown, area_shape, match_boundary0, match_boundary1, match_boundary2, hatcolor, uppercolor, lowercolor, whole, alarm_mode, timepoint)
		if boundary != "[]" {
			boundary = SetBoundaryStr(boundary)
		}
		if wallOn == "0" {
			match_boundary0 = boundary
			match_boundary1 = ""
			match_boundary2 = ""
		}
		if wallOn == "1" {
			match_boundary0 = ""
			match_boundary1 = boundary
			match_boundary2 = ""
		}
		if wallOn == "2" {
			match_boundary0 = ""
			match_boundary1 = ""
			match_boundary2 = boundary
		}
		if wallOn == "" {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "报警模式不能为空"}
			this.ServeJSON()
			return
		}
		res, err := o.Raw("update ss_func_camera set function_type = ?,DistanceMode = ?,threshold = ?,frequency = ?,wallOn = ?,wallrate = ?,rateup = ?,ratedown = ?,area_shape=?,match_boundary0=?,match_boundary1=?,match_boundary2=?,hatcolor=?,uppercolor=?,lowercolor=?,whole=?,alarm_mode=?,timepoint=?,alarm_time_second=?,video_switch=? where id = ?;",
			function_type, string(DistanceMode), string(threshold), string(frequency), string(wallOn), string(wallrate), string(rateup), string(ratedown), string(area_shape), string(match_boundary0), string(match_boundary1), string(match_boundary2), string(hatcolor), string(uppercolor), string(lowercolor), string(whole), string(alarm_mode), string(timepoint), string(alarm_time_second), string(video_switch), string(camera_id)).Exec()
		log.Println(res, err)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "坐标数据添加失败"}
			this.ServeJSON()
			return
		}
		this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "坐标数据添加成功"}
		this.ServeJSON()

	} else if this.GetString("function_type") == "sleep_count" {
		function_type, camera_id, location, threshold, frequency, boundary, real_sleep_minutes, DistanceMode, area_shape, alarm_mode, alarm_time_second, video_switch, bg_thr, wallrate, distance, timepoint := this.GetString("function_type"),
			this.GetString("id"),
			this.GetString("location"),
			this.GetString("threshold"),
			this.GetString("frequency"),
			this.GetString("boundary"),
			//this.GetString("almost_sleep_seconds"),
			this.GetString("real_sleep_minutes"),
			this.GetString("DistanceMode"),
			this.GetString("area_shape"),
			this.GetString("alarm_mode"),
			this.GetString("alarm_time_second"),
			this.GetString("video_switch"),
			this.GetString("bg_thr"),
			this.GetString("wallrate"),
			this.GetString("distance"),
			this.GetString("timepoint")
		if boundary != "[]" {
			boundary = SetBoundaryStr(boundary)
		}
		log.Println("睡岗检测布防设备参数>>>", function_type, camera_id, location, threshold, frequency, boundary, real_sleep_minutes, DistanceMode, "形状:"+area_shape, alarm_mode)
		o := O
		res, err := o.Raw("update ss_func_camera set function_type = ?,location = ?,boundary = ?,threshold = ?,frequency = ?,real_sleep_minutes = ?,DistanceMode = ?,area_shape = ?,alarm_mode = ?,alarm_time_second = ?,video_switch = ?,bg_thr=?,wallrate=?,distance=?,timepoint=? where id = ?;",
			function_type, location, string(boundary), string(threshold), string(frequency), string(real_sleep_minutes), string(DistanceMode), string(area_shape), string(alarm_mode), string(alarm_time_second), string(video_switch), string(bg_thr), string(wallrate), string(distance), string(timepoint), string(camera_id)).Exec()
		log.Println(res, err)

		if err != nil {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "坐标数据添加失败"}
			this.ServeJSON()
			return
		}
		this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "坐标数据添加成功"}
		this.ServeJSON()
	} else if this.GetString("function_type") == "queue_count" {
		function_type,
			camera_id,
			location,
			threshold,
			frequency,
			boundary,
			DistanceMode,
			area_shape,
			alarm_mode,
			detect_mode,
			departure_minute,
			num_person_threshold,
			alarm_time_second,
			video_switch,
			distance,
			timepoint :=
			this.GetString("function_type"),
			this.GetString("id"),
			this.GetString("location"),
			this.GetString("threshold"),
			this.GetString("frequency"),
			this.GetString("boundary"),
			this.GetString("DistanceMode"),
			this.GetString("area_shape"),
			this.GetString("alarm_mode"),
			"0",
			this.GetString("departure_minute"),
			this.GetString("num_person_threshold"),
			this.GetString("alarm_time_second"),
			this.GetString("video_switch"),
			this.GetString("distance"),
			this.GetString("timepoint")
		if boundary != "[]" {
			boundary = SetBoundaryStr(boundary)
		}
		res, err := o.Raw("update ss_func_camera set function_type = ?,location = ?,boundary = ?,threshold = ?,frequency = ?,DistanceMode = ?,area_shape = ?,alarm_mode = ?,detect_mode = ?,num_person_threshold = ?,alarm_time_second=?,video_switch=?,distance=?,departure_minute=?,timepoint=? where id = ?;",
			function_type, location, string(boundary), string(threshold), string(frequency), string(DistanceMode), string(area_shape), string(alarm_mode), string(detect_mode), string(num_person_threshold), string(alarm_time_second), string(video_switch), string(distance), departure_minute, string(timepoint), string(camera_id)).Exec()
		log.Println(res, err)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "坐标数据添加失败"}
			this.ServeJSON()
			return
		}
		this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "坐标数据添加成功"}
		this.ServeJSON()
	} else if this.GetString("function_type") == "fire" || this.GetString("function_type") == "smoke" {
		function_type,
			camera_id,
			location,
			threshold,
			frequency,
			boundary,
			boundary2,
			bg_threshold,
			perimeterThreshold,
			DistanceMode,
			area_shape, alarm_mode,
			alarm_time_second,
			video_switch,
			distance,
			timepoint :=
			this.GetString("function_type"),
			this.GetString("id"),
			this.GetString("location"),
			this.GetString("threshold"),
			this.GetString("frequency"),
			this.GetString("boundary"),
			this.GetString("boundary2"),
			this.GetString("bg_threshold"),
			this.GetString("perimeterThreshold"),
			this.GetString("DistanceMode"),
			this.GetString("area_shape"),
			this.GetString("alarm_mode"),
			this.GetString("alarm_time_second"),
			this.GetString("video_switch"),
			this.GetString("distance"),
			this.GetString("timepoint")

		log.Println("火焰检测检测 & 烟雾检测布防设备参数>>>", function_type, camera_id, location, threshold, frequency, boundary, bg_threshold, perimeterThreshold, DistanceMode, area_shape, alarm_mode)
		if boundary != "[]" {
			boundary = SetBoundaryStr(boundary)
		}
		_, err := o.Raw("update ss_func_camera set function_type = ?,location = ?,boundary = ?,threshold = ?,frequency = ?,bg_threshold = ?,perimeterThreshold = ?,DistanceMode = ?,area_shape = ?,alarm_mode = ?,alarm_time_second=?,video_switch=?, boundary2=?,distance=?,timepoint=?  where id = ?;",
			function_type, location, string(boundary), string(threshold), string(frequency), string(bg_threshold), string(perimeterThreshold), string(DistanceMode), string(area_shape), string(alarm_mode), string(alarm_time_second), string(video_switch), boundary2, string(distance), string(timepoint), string(camera_id)).Exec()

		if err != nil {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "坐标数据添加失败"}
			this.ServeJSON()
			return
		}

		this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "坐标数据添加成功"}
		this.ServeJSON()
	} else if this.GetString("function_type") == "leakage" {
		function_type, camera_id, location, threshold, frequency, boundary, DistanceMode, match_boundary0, match_boundary1, match_boundary2, alarm_mode, perimeterThreshold, detectNum, warning_number, leakThreshold, video_switch, timepoint :=
			this.GetString("function_type"),
			this.GetString("id"),
			this.GetString("location"),
			this.GetString("threshold"),
			this.GetString("frequency"),
			this.GetString("boundary"),
			this.GetString("DistanceMode"),
			this.GetString("match_boundary0"),
			this.GetString("match_boundary1"),
			this.GetString("match_boundary2"),
			this.GetString("alarm_mode"),
			this.GetString("perimeterThreshold"),
			this.GetString("detect_number"),
			this.GetString("warning_number"),
			this.GetString("leakThreshold"),
			this.GetString("video_switch"),
			this.GetString("timepoint")
		//c.GetString("fireExtinguisher_boundary"),
		//c.GetString("tubing_boundary"),
		//c.GetString("groundwire_boundary")
		log.Println("跑冒滴漏检测布防设备参数>>>", function_type, camera_id, location, threshold, perimeterThreshold, frequency, boundary, DistanceMode, match_boundary0, match_boundary1, match_boundary2, alarm_mode)
		if boundary != "[]" {
			boundary = SetBoundaryStr(boundary)
		}
		res, err := o.Raw("update ss_func_camera set function_type = ?,boundary = ?,threshold = ?,frequency = ?,DistanceMode = ?,boundary= ?,match_boundary0=?,match_boundary1=?,match_boundary2=?,alarm_mode=?,perimeterThreshold=?,detect_number=?,warning_number=?,leakThreshold=?,video_switch=?,timepoint=?  where id = ?;",
			function_type, string(boundary), string(threshold), string(frequency), string(DistanceMode), string(boundary), string(match_boundary0), string(match_boundary1), string(match_boundary2), string(alarm_mode), string(perimeterThreshold), string(detectNum), string(warning_number), string(leakThreshold), string(video_switch), string(timepoint), string(camera_id)).Exec()
		if err == nil {
			num, _ := res.RowsAffected()
			fmt.Println("mysql row affected nums: ", num)
			this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "坐标数据添加成功~"}
			this.ServeJSON()
		} else {
			log.Println(err)
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "坐标数据添加错误~"}
			this.ServeJSON()
		}
	} else {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "查询类型不存在"}
		this.ServeJSON()
	}
	ID := this.GetString("id")
	var funcCamera models.FuncCamera
	err := o.QueryTable(new(models.FuncCamera)).Filter("id", ID).RelatedSel().One(&funcCamera)
	if err == nil && funcCamera.Camera.Status == "0" {
		var camera models.Camera
		err := o.QueryTable(new(models.Camera)).Filter("camera_id", funcCamera.Camera.ID).One(&camera)
		camera.Status = "1"
		_, err = o.Update(&camera, "status")
		if err != nil {
			this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "配置失败"}
			this.ServeJSON()
		}
	}
	return
}

/**
删除指定摄像机配置算法生成的负样本区域
*/
type DelCameraFuncNegative struct {
	beego.Controller
}

func (this *DelCameraFuncNegative) Post() {
	cameraID := this.GetString("cameraID")
	if cameraID == "" {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "摄像机参数错误"}
		this.ServeJSON()
		return
	}

	o := O
	var cameraFunc models.FuncCamera

	o.QueryTable(new(models.FuncCamera)).Filter("Camera__ID", cameraID).One(&cameraFunc)
	if cameraFunc.ID != "" {
		cameraFunc.Negative = ""
		_, err := o.Update(&cameraFunc, "negative")
		if err != nil {
			LogsError("清空摄像机Negative字段异常:", err)
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "操作失败"}
			this.ServeJSON()
			return
		}
		this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "操作成功"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "操作失败, 未查询到摄像机"}
		this.ServeJSON()
		return
	}

}
