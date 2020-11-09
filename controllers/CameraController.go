package controllers

import (
	"Artifice_V2.0.0/models"
	"Artifice_V2.0.0/util"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 摄像机管理页面
type PlatformCameraManagerPageController struct {
	beego.Controller
}

func (this *PlatformCameraManagerPageController) Get() {
	hostID := this.GetString("hostID")
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["hostID"] = hostID
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
		this.TplName = "platform/cameralist.html"
	}
}

// 摄像机管理页面
type PlatformCameraManagerIEPageController struct {
	beego.Controller
}

func (this *PlatformCameraManagerIEPageController) Get() {
	hostID := this.GetString("hostID")
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["hostID"] = hostID
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
		this.TplName = "platform/cameralistIE.html"
	}
}

type SelCameraListHandler struct {
	beego.Controller
}

func (this *SelCameraListHandler) Post() {
	limits, pages, hostId, cameraName := this.GetString("limit"), this.GetString("page"),
		this.GetString("hostId"), this.GetString("cameraName")
	start_page, err := strconv.Atoi(pages)
	limitss, err := strconv.Atoi(limits)
	start_pages := (start_page - 1) * limitss
	if err != nil {
		log.Println("摄像机列表：", err)
		return
	}
	if cameraName == "" {
		sqlStartStr := "select * from ss_camera where 1=1"
		var maps []orm.Params
		var mapsr []orm.Params
		o := O
		var num int64
		if hostId != "" {
			sqlStartStr += " and host_id=?"
			_, err = o.Raw(sqlStartStr+" order by create_time desc limit ?,?", hostId, int(start_pages), int(limitss)).Values(&maps)
			num, err = o.Raw(sqlStartStr, hostId).Values(&mapsr)
		} else {
			_, err = o.Raw(sqlStartStr+" order by create_time desc limit ?,?", int(start_pages), int(limitss)).Values(&maps)
			num, err = o.Raw(sqlStartStr).Values(&mapsr)
		}

		if err != nil {
			log.Println("摄像机列表：", err)
			mystruct := &util.ResultB{0, 0, nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}
		for i := 0; i < len(maps); i++ {
			maps[i]["switch"] = maps[i]["switch_c"]
		}

		mystruct := &util.ResultB{0, num, &maps}
		this.Data["json"] = mystruct
		this.ServeJSON()
	} else {
		sqlStartStr := "select * from ss_camera where 1=1"
		var maps []orm.Params
		o := O
		num, err := o.Raw(sqlStartStr+"host_id=? and camera_title like ? order by camera_id asc limit ?,?", hostId, cameraName, int(start_pages), int(limitss)).Values(&maps)
		if err != nil {
			log.Println("摄像机列表：", err)
			return
		}
		for i := 0; i < len(maps); i++ {
			maps[i]["switch"] = maps[i]["switch_c"]
		}
		var mapsr []orm.Params
		num, err = o.Raw(sqlStartStr).Values(&mapsr)
		mystruct := &util.ResultB{0, num, &maps}
		this.Data["json"] = mystruct
		this.ServeJSON()
	}
}

type SelPlaceHandler struct {
	beego.Controller
}

func (this *SelPlaceHandler) Post() {
	var maps []orm.Params
	o := O
	num, err := o.Raw("select * from ss_place").Values(&maps)
	if err != nil {
		log.Println("分页查询主机列表：", err)
		return
	}
	mystruct := &util.ResultB{0, num, &maps}
	this.Data["json"] = mystruct
	this.ServeJSON()
}

//添加修改相机
type SaveCameraHandler struct {
	beego.Controller
}

func (this *SaveCameraHandler) Post() {
	camera_id, host_id, titleName, cdConfig, locaTion, function_type, code, width, height, algorithm_url, algorithm_width, algorithm_height, videoCode, place_id, cameraNub, screen := this.GetString("camera_id"), this.GetString("host_id"), this.GetString("TitleName"), this.GetString("CdConfig"),
		this.GetString("LocaTion"), this.GetString("function_type"), this.GetString("mcode"), this.GetString("width"), this.GetString("height"), this.GetString("algorithm_url"), this.GetString("algorithm_width"), this.GetString("algorithm_height"),
		this.GetString("videoCode"), this.GetString("place_id"), this.GetString("CameraNub"), this.GetString("screen")

	//去空格
	cameraNub = strings.Replace(cameraNub, " ", "", -1)
	titleName = strings.Replace(titleName, " ", "", -1)
	cdConfig = strings.Replace(cdConfig, " ", "", -1)
	locaTion = strings.Replace(locaTion, " ", "", -1)
	width = strings.Replace(width, " ", "", -1)
	height = strings.Replace(height, " ", "", -1)

	//摄像机注册时间戳
	widths, err := strconv.Atoi(width)
	heights, err := strconv.Atoi(height)
	algorithm_widths, err := strconv.Atoi(algorithm_width)
	algorithm_heights, err := strconv.Atoi(algorithm_height)
	if err != nil {
		return
	}
	fmt.Println(camera_id, host_id, titleName, cdConfig, locaTion, function_type, code, width, height)

	if code == "0" {
		uid := uuid.Must(uuid.NewV4()).String()
		uid = strings.Replace(uid, "-", "", -1) // 生成SessionID
		var maps []orm.Params
		o := O
		log.Println(function_type, " "+cdConfig+" "+host_id)
		num, err := o.Raw("select * from ss_camera where camera_function_type = ? and (cdconfig = ? or algorithm_url=?)", function_type, cdConfig, algorithm_url).Values(&maps)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态设置异常！"}
			this.ServeJSON()
			return
		}
		if num > 0 {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "配置和类型相同的摄像机不能被重复添加哦"}
			this.ServeJSON()
			return
		} else {

			if screen == "1" {
				screen = "on"
				o := O
				var maps []orm.Params
				num, err := o.Raw("select * from ss_camera").Values(&maps)
				log.Println(num, err)
				for i := 0; i < len(maps); i++ {
					num, err := o.Raw("update ss_camera set screen ='off' where camera_id=?", maps[i]["camera_id"].(string)).Exec()
					if err != nil {
						log.Println(num)
						this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "修改默认状态失败"}
						this.ServeJSON()
						return
					}
				}
			} else {
				screen = "off"
			}
			var sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,DistanceMode,timepoint) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','1','[[00:00,23:59],[],[],[]]');"
			if function_type == "boundary" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,DistanceMode,wallOn,match_boundary0,timepoint,alarm_mode,frequency,rateup,ratedown) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','1','0','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[[00:00,23:59],[],[],[]]','1','15','0.8','0.5');"
			} else if function_type == "cloths" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,DistanceMode,wallOn,match_boundary0,timepoint,alarm_mode,frequency,hatcolor,rateup,ratedown) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','1','0','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[[00:00,23:59],[],[],[]]','1','15','1','0.8','0.5');"
			} else if function_type == "sleep_count" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,timepoint,alarm_mode,frequency,real_sleep_minutes,distance) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','[[00:00,23:59],[],[],[]]','1','15','3','1.0');"
			} else if function_type == "queue_count" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,timepoint,alarm_mode,frequency,num_person_threshold,departure_minute,distance) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','[[00:00,23:59],[],[],[]]','1','15','3','5','1.0');"
			} else if function_type == "fire" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,timepoint,alarm_mode,frequency,distance,perimeterThreshold) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','[[00:00,23:59],[],[],[]]','1','15','1.0','10');"
			} else if function_type == "smoke" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,timepoint,alarm_mode,frequency,distance,perimeterThreshold) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','[[00:00,23:59],[],[],[]]','1','15','1.0','10');"
			} else if function_type == "leakage" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,DistanceMode,wallOn,match_boundary0,timepoint,alarm_mode,frequency,perimeterThreshold,detect_number,warning_number) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','0','1','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[[00:00,23:59],[],[],[]]','1','15','5','10','8');"
			}

			_, err := o.Raw(sql2, uid, host_id, titleName, cdConfig, locaTion, function_type, time.Now().Format("2006-01-02 15:04:05")).Exec()
			if err != nil {
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态设置异常！"}
				this.ServeJSON()
				return
			}
			_, err = o.Raw("insert into ss_camera(camera_id,host_id,camera_title,cdconfig,camera_location,camera_function_type,create_time, width, height,wRate,hRate,algorithm_url,algorithm_width,algorithm_height,videoCode,place_id,cameraNub,screen) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);", uid, host_id, string(titleName),
				string(cdConfig), string(locaTion), string(function_type), time.Now().Format("2006-01-02 15:04:05"), widths, heights, "480", "270", algorithm_url, algorithm_widths, algorithm_heights, videoCode, place_id, cameraNub, screen).Exec()

			if err != nil {
				log.Println("添加摄像机异常:", err)
				o.Raw("delete from ss_func_camera where cameraId =?;", uid).Exec()
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "添加摄像机异常"}
				this.ServeJSON()
				return
			}

		}

		this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "添加成功！"}
		this.ServeJSON()
	} else if code == "1" {
		var maps []orm.Params
		var smaps []orm.Params
		o := O

		_, err := o.Raw("select * from ss_camera where camera_id=?", camera_id).Values(&smaps)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "更新摄像机配置失败"}
			this.ServeJSON()
			return
		}
		if len(smaps) == 0 {
			log.Println("未查询到该摄像机")
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "未查询到该摄像机"}
			this.ServeJSON()
			return
		}
		if smaps[0]["switch_c"] == "on" {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "摄像头启动中无法操作!"}
			this.ServeJSON()
			return
		}

		if smaps[0]["camera_function_type"] != function_type {

			num, err := o.Raw("select * from ss_camera where camera_function_type = ? and (cdconfig = ? or algorithm_url=?)", function_type, cdConfig, algorithm_url).Values(&maps)
			if err != nil {
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态设置异常！"}
				this.ServeJSON()
				return
			}
			if num > 0 {
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "配置和类型相同的摄像机不能被重复添加哦"}
				this.ServeJSON()
				return
			}

			_, err = o.Raw("delete from ss_func_camera where cameraId =? and function_type=?;", camera_id, smaps[0]["camera_function_type"]).Exec()
			if err != nil {
				log.Println("状态设置异常~\n", err)
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态设置异常"}
				this.ServeJSON()
				return
			}
			var sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,DistanceMode,timepoint) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','1','[[00:00,23:59],[],[],[]]');"
			if function_type == "boundary" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,DistanceMode,wallOn,match_boundary0,timepoint,alarm_mode,frequency,rateup,ratedown) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','1','0','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[[00:00,23:59],[],[],[]]','1','15','0.8','0.5');"
			} else if function_type == "cloths" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,DistanceMode,wallOn,match_boundary0,timepoint,alarm_mode,frequency,hatcolor,rateup,ratedown) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','1','0','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[[00:00,23:59],[],[],[]]','1','15','1','0.8','0.5');"
			} else if function_type == "sleep_count" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,timepoint,alarm_mode,frequency,real_sleep_minutes,distance) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','[[00:00,23:59],[],[],[]]','1','15','3','1.0');"
			} else if function_type == "queue_count" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,timepoint,alarm_mode,frequency,num_person_threshold,departure_minute,distance) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','[[00:00,23:59],[],[],[]]','1','15','3','5','1.0');"
			} else if function_type == "fire" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,timepoint,alarm_mode,frequency,distance,perimeterThreshold) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','[[00:00,23:59],[],[],[]]','1','15','1.0','10');"
			} else if function_type == "smoke" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,timepoint,alarm_mode,frequency,distance,perimeterThreshold) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','[[00:00,23:59],[],[],[]]','1','15','1.0','10');"
			} else if function_type == "leakage" {
				sql2 = "insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,boundary2,threshold,video_switch,DistanceMode,wallOn,match_boundary0,timepoint,alarm_mode,frequency,perimeterThreshold,detect_number,warning_number) values(?,?, ?,?,?,?,?,'off','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[]','0.4','0','0','1','[[0,0],[480,0],[480,270],[0,270],[0,0]];','[[00:00,23:59],[],[],[]]','1','15','5','10','8');"
			}
			_, err = o.Raw(sql2, camera_id, host_id, titleName, cdConfig, locaTion, function_type, time.Now().Format("2006-01-02 15:04:05")).Exec()
			if err != nil {
				log.Println("状态设置异常~\n", err)
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态设置异常"}
				this.ServeJSON()
				return
			}
		} else {
			_, err = o.Raw("update ss_func_camera set title_name=?, cdconfig=?, location=? where cameraId=? and function_type=? and host_id=?;", titleName, cdConfig, locaTion, camera_id, function_type, host_id).Exec()
			if err != nil {
				log.Println("状态设置异常~\n", err)
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态设置异常"}
				this.ServeJSON()
				return
			}
		}

		//设置投屏功能处理
		if screen == "1" {
			screen = "on"
			var maps []orm.Params
			num, err := o.Raw("select * from ss_camera").Values(&maps)
			log.Println(num, err)
			for i := 0; i < len(maps); i++ {
				num, err := o.Raw("update ss_camera set screen ='off' where camera_id=?", maps[i]["camera_id"].(string)).Exec()
				if err != nil {
					log.Println(num)
					this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "修改默认状态失败"}
					this.ServeJSON()
					return
				}
			}
		} else {
			screen = "off"
		}

		_, err = o.Raw("update ss_camera set camera_title = ? ,cameraNub = ? , camera_location = ? , camera_function_type = ? ,"+
			" cdconfig = ? , width = ?, height = ?, wRate = ?, hRate = ?,screen=?,algorithm_url = ? ,"+
			"algorithm_width = ?,algorithm_height = ?, videoCode=?, place_id=? where camera_id = ?;",
			titleName, cameraNub, locaTion, function_type,
			cdConfig, widths, heights, "480", "270", screen,
			algorithm_url, algorithm_widths, algorithm_heights,
			videoCode, place_id, camera_id).Exec()
		if err != nil {
			log.Println("状态设置异常\n", err)
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "更新摄像机配置失败"}
			this.ServeJSON()
			return
		}

		//num, err = o.Raw("select * from ss_func_camera where cameraId = ?", camera_id).Values(&maps)
		//var camera_function_type = ""
		//if num > 0 {
		//	camera_function_type = maps[0]["function_type"].(string)
		//}
		//res, err := o.Raw("update ss_camera set camera_title = ? , camera_location = ? , camera_function_type = ? , create_time = ? , cdconfig = ? , width = ?, height = ?, wRate = ?, hRate = ?,algorithm_url = ? ,algorithm_width = ?,algorithm_height = ?, videoCode=?, place_id=? where camera_id = ?;",
		//	titleName, locaTion, function_type, time.Now().Format("2006-01-02 15:04:05"), cdConfig, widths, heights, widths/4, heights/4, algorithm_url, algorithm_widths, algorithm_heights, videoCode, place_id, camera_id).Exec()
		//log.Println(res, err)
		//log.Println(camera_id, titleName, cdConfig, locaTion, function_type, code, widths/4, heights/4, "<<<<<-------------")
		//if err != nil {
		//	this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "更新摄像机配置失败！"}
		//	this.ServeJSON()
		//	return
		//}
		//if function_type == camera_function_type {
		//	////todo:更新分组摄像机初始状态
		//	res, err := o.Raw("update ss_func_camera set title_name = ?, cdconfig = ? , location = ?  where cameraId = ?;",
		//		titleName, cdConfig, locaTion, string(camera_id)).Exec()
		//	log.Println(res, err)
		//	if err != nil {
		//		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态设置异常！"}
		//		this.ServeJSON()
		//		return
		//	}
		//} else {
		//	res, err := o.Raw("delete from ss_func_camera where cameraId =?", camera_id).Exec()
		//	log.Println(res, err)
		//	if err != nil {
		//		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态设置异常！"}
		//		this.ServeJSON()
		//		return
		//	}
		//	/*res, err = o.Raw("insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,threshold) values(?,?, ?,?,?,?,?,'off','[]','0.4');",
		//	camera_id, host_id, titleName, cdConfig, locaTion, function_type, time.Now().Format("2006-01-02 15:04:05")).Exec()*/
		//	if string(function_type) == "face_detection" { // 人脸识别
		//		res, err = o.Raw("insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,threshold,DistanceMode,wallrate,topk) values(?,?, ?,?,?,?,?,'off','[]','0.6','0','0.4','5');", camera_id, host_id, titleName, cdConfig, locaTion, function_type, time.Now().Format("2006-01-02 15:04:05")).Exec()
		//	} else if string(function_type) == "mask_detection" { // 口罩识别
		//		res, err = o.Raw("insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,threshold,DistanceMode,wallrate,topk) values(?,?, ?,?,?,?,?,'off','[]','0.6','0','0.4','5');", camera_id, host_id, titleName, cdConfig, locaTion, function_type, time.Now().Format("2006-01-02 15:04:05")).Exec()
		//	} else if string(function_type) == "sleep_count" { // 睡岗
		//		res, err = o.Raw("insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,threshold,bg_thr,wallrate,real_sleep_minutes) values(?,?,?,?,?,?,?,'off','[]','0.4','50','0.5','10');", camera_id, host_id, titleName, cdConfig, locaTion, function_type, time.Now().Format("2006-01-02 15:04:05")).Exec()
		//	} else {
		//		res, err = o.Raw("insert into ss_func_camera(cameraId,host_id,title_name,cdconfig,location,function_type,create_date,switch,boundary,threshold) values(?,?, ?,?,?,?,?,'off','[]','0.4');", camera_id, host_id, titleName, cdConfig, locaTion, function_type, time.Now().Format("2006-01-02 15:04:05")).Exec()
		//	}
		//	log.Println(res, err)
		//	if err != nil {
		//		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态设置异常！"}
		//		this.ServeJSON()
		//		return
		//	}
		//}
		this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "修改成功！"}
		this.ServeJSON()
	}
}

type SetCameraStatusHandler struct {
	beego.Controller
}

func (this *SetCameraStatusHandler) Post() {
	//hostId := this.GetString("hostId")
	ids := this.GetString("ids")
	code := this.GetString("code")
	log.Println("id+code:", ids, code)
	paramsa := strings.Split(ids, ",")
	//var maps []orm.Params
	o := O
	for _, v := range paramsa {
		log.Println("v", v, code)
		if code == "0" {
			res, err := o.Raw("update ss_camera set switch_c='off' where camera_id=?", v).Exec()
			log.Println(res, err)
			if err != nil {
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态修改失败！"}
				this.ServeJSON()
				return
			}
			res, err = o.Raw("update ss_func_camera set switch='off' where cameraId=?", v).Exec()
			if err != nil {
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态修改失败！"}
				this.ServeJSON()
				return
			}
			this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "停用成功！"}
			this.ServeJSON()
		} else if code == "1" {
			//
			//num, err := o.Raw("select * from ss_camera where switch_c='on' and host_id=?;", hostId).Values(&maps)
			//if err != nil {
			//	this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态修改失败！"}
			//	this.ServeJSON()
			//	return
			//}
			//if num >= 8 {
			//	this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "开启失败，最多开启8路摄像头！"}
			//	this.ServeJSON()
			//	return
			//}
			res, err := o.Raw("update ss_camera set switch_c='on' where camera_id=?", v).Exec()
			if err != nil {
				log.Println(res, err)
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态修改失败！"}
				this.ServeJSON()
				return
			}
			res, err = o.Raw("update ss_func_camera set switch='on' where cameraId=?", v).Exec()
			if err != nil {
				log.Println(res, err)
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "状态修改失败！"}
				this.ServeJSON()
				return
			}
			this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "开启成功！"}
			this.ServeJSON()
		}
	}
}

type DelCameraDataHandler struct {
	beego.Controller
}

func (this *DelCameraDataHandler) Post() {
	params := this.GetString("param")
	paramsa := strings.Split(params, ",")
	log.Println("paramsa:", paramsa)
	log.Println(reflect.TypeOf(paramsa), len(paramsa), paramsa)
	o := O
	//： 执行删除语句
	for _, mains := range paramsa {
		res, err := o.Raw("delete from ss_func_camera where cameraId = ?;", mains).Exec()
		log.Println(res, err)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "删除失败！"}
			this.ServeJSON()
			return
		}
		res, err = o.Raw("delete from ss_camera where camera_id = ?;", mains).Exec()
		if err != nil {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "删除失败！"}
			this.ServeJSON()
			return
		}

	}
	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "删除成功！"}
	this.ServeJSON()
}

type SelCameraForFuncType struct {
	beego.Controller
}

func (this *SelCameraForFuncType) Post() {
	funcType := this.GetString("funcType")
	o := O
	var cameraList []models.Camera
	_, err := o.QueryTable(new(models.Camera)).Filter("FuncCamera__Type", funcType).All(&cameraList)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "查询失败！", "Data": nil}
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "查询成功！", "Data": cameraList}
	this.ServeJSON()
	return
}
