package handle_comment

import (
	"Artifice_V2.0.0/controllers"
	"Artifice_V2.0.0/models"
	"Artifice_V2.0.0/util"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	"strconv"
	"strings"
	"time"
)

type WorkshopAlarmController struct {
	beego.Controller
}
type WorkshopAlarmController2 struct {
	beego.Controller
}

type WorkshopUnAlarmListController struct {
	beego.Controller
}
type WorkshopUnAlarmListController2 struct {
	beego.Controller
}

type WorkshopUnAlarmController struct {
	beego.Controller
}
type WorkshopUnAlarmController2 struct {
	beego.Controller
}

type WorkshopDealAlarmController struct {
	beego.Controller
}

type WorkshopAlarmListController2 struct {
	beego.Controller
}

type WorkshopNotationListController2 struct {
	beego.Controller
}

func (c *WorkshopAlarmController) Get() {
	organizeId := c.Ctx.GetCookie("HCOrganize")
	log.Println("获取orgainzeId: ", organizeId)
	// 根据部门ID查询部门关联的所有CameraID
	o := controllers.O

	var system []orm.Params
	num, err := o.Raw("select * from ss_system").Values(&system)

	if err == nil && num > 0 {
		c.Data["SystemTitle"] = system[0]["websitename"]
		log.Println("system:", system[0]["websitename"])
	} else {
		c.Data["SystemTitle"] = "AI风险防控预警平台"
	}

	var maps []orm.Params
	num1, err1 := o.Raw("select * from ss_organize where id=?", organizeId).Values(&maps)
	if err1 == nil && num1 > 0 {
		log.Println("查询个数：", num1)
		c.Data["organizeName"] = maps[0]["name"]
	}

	num, err = o.Raw("select camera_id from ss_camera where organizeid=?", organizeId).Values(&maps)
	if err != nil || num == 0 {
		log.Println("num:", num)
		c.Data["sum"] = 0
		c.Data["num"] = 0
		c.Data["unnum"] = 0
	} else {
		ids := ""
		var idsv = []string{}
		for i := 0; i < len(maps); i++ {
			ids += ",?"
			idsv = append(idsv, maps[i]["camera_id"].(string))
		}
		ids = ids[1:len(ids)]
		//for i := 0; i < len(maps); i++ {
		//	if i != 0 {
		//		ids += ","
		//	}
		//	ids += "'"
		//	ids += maps[i]["camera_id"].(string)
		//	ids += "'"
		//}
		var summaps []orm.Params
		sql1 := "select * from ss_alarm where camera_id in (" + ids + ")"
		sum, err := o.Raw(sql1, idsv).Values(&summaps)
		log.Println("sum2:", sum)
		var nummaps []orm.Params
		sql2 := "select * from ss_alarm where alarm_status=1 and camera_id in (" + ids + ")"
		num, err := o.Raw(sql2, idsv).Values(&nummaps)
		log.Println("num2:", num)
		var unsummaps []orm.Params
		sql3 := "select * from ss_alarm where alarm_status=0 and camera_id in (" + ids + ")"
		num_un, err := o.Raw(sql3, idsv).Values(&unsummaps)
		log.Println("num_un2", num_un)
		if err == nil {
			c.Data["sum"] = sum
			c.Data["num"] = num
			c.Data["unnum"] = num_un
		}
	}
	c.Data["Name"] = c.Ctx.GetCookie("HCName")
	c.TplName = "handle_comment/handle/handle/workshop_alarm.html"
}

func (w *WorkshopAlarmController2) Get() {
	userId := w.Ctx.GetCookie("HCAdmin_id")
	organizeId := w.Ctx.GetCookie("HCOrganize")

	o := controllers.O

	var system []orm.Params
	num, err := o.Raw("select * from ss_system").Values(&system)

	if err == nil && num > 0 {
		w.Data["SystemTitle"] = system[0]["websitename"]
		log.Println("system:", system[0]["websitename"])
	} else {
		w.Data["SystemTitle"] = "AI风险防控预警平台"
	}

	var maps []orm.Params
	num1, err1 := o.Raw("select * from ss_organize where id=?", organizeId).Values(&maps)
	if err1 == nil && num1 > 0 {
		log.Println("查询个数：", num1)
		w.Data["organizeName"] = maps[0]["name"]
	}
	w.Data["Name"] = w.Ctx.GetCookie("HCName")

	var user models.User
	var cameras []models.Camera
	err = o.QueryTable(new(models.User)).Filter("admin_id", userId).RelatedSel().One(&user)
	if err != nil {
		log.Print("用户查询err:", err)
	}

	if user.Organize.ID == "" {
		w.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户组织架构异常"}
		w.ServeJSON()
		return
	}

	var org = user.Organize

	orgIds, err := controllers.SelOrgIDBelongUserOrg(user.Organize)
	if err != nil {
		controllers.LogsError("查询用户下属组织异常:", err)

	}
	var ids = []string{org.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}

	_, err = o.QueryTable(new(models.Camera)).Filter("Place__Organize__ID__in", ids).RelatedSel().All(&cameras)

	if err != nil {
		controllers.LogsError("查询设备信息异常:", err)
		w.Data["sum"] = 0
		w.Data["num"] = 0
		w.Data["unnum"] = 0
	} else {
		ids := ""
		idsv := []string{}
		for i := 0; i < len(maps); i++ {
			ids += ",?"
			idsv = append(idsv, cameras[i].ID)
		}
		ids = ids[1:len(ids)]
		//for i := 0; i < len(cameras); i++ {
		//	if i != 0 {
		//		ids += ","
		//	}
		//	ids += "'"
		//	ids += cameras[i].ID
		//	ids += "'"
		//}
		var summaps []orm.Params
		sql1 := "select * from ss_alarm where camera_id in (" + ids + ")"
		sum, err := o.Raw(sql1, idsv).Values(&summaps)
		log.Println("sum2:", sum)
		var nummaps []orm.Params
		sql2 := "select * from ss_alarm where alarm_status=1 and camera_id in (" + ids + ")"
		num, err := o.Raw(sql2, idsv).Values(&nummaps)
		log.Println("num2:", num)
		var unsummaps []orm.Params
		sql3 := "select * from ss_alarm where alarm_status=0 and camera_id in (" + ids + ")"
		num_un, err := o.Raw(sql3, idsv).Values(&unsummaps)
		log.Println("num_un2", num_un)
		if err == nil {
			w.Data["sum"] = sum
			w.Data["num"] = num
			w.Data["unnum"] = num_un
		}
	}

	//num, err = o.Raw("select cameraid from ss_map_permission where userid=?", userId).Values(&users)
	//if err == nil && num == 0 {
	//	w.Data["sum"] = 0
	//	w.Data["num"] = 0
	//	w.Data["unnum"] = 0
	//} else if err == nil {
	//	ids := ""
	//	for i := 0; i < len(users); i++ {
	//		if i != 0 {
	//			ids += ","
	//		}
	//		ids += "'"
	//		ids += users[i]["cameraid"].(string)
	//		ids += "'"
	//	}
	//	var summaps []orm.Params
	//	sql1 := "select * from ss_alarm where camera_id in (" + ids + ")"
	//	sum, err := o.Raw(sql1).Values(&summaps)
	//	log.Println("sum2:", sum)
	//	var nummaps []orm.Params
	//	sql2 := "select * from ss_alarm where alarm_status=1 and camera_id in (" + ids + ")"
	//	num, err := o.Raw(sql2).Values(&nummaps)
	//	log.Println("num2:", num)
	//	var unsummaps []orm.Params
	//	sql3 := "select * from ss_alarm where alarm_status=0 and camera_id in (" + ids + ")"
	//	num_un, err := o.Raw(sql3).Values(&unsummaps)
	//	log.Println("num_un2", num_un)
	//	if err == nil {
	//		w.Data["sum"] = sum
	//		w.Data["num"] = num
	//		w.Data["unnum"] = num_un
	//	}
	//}
	w.Data["Name"] = w.Ctx.GetCookie("HCName")
	w.TplName = "handle_comment/handle/workshop_alarm.html"
}

/**车间今日报警列表**/
func (c *WorkshopUnAlarmController) Get() {
	organizeId := c.Ctx.GetCookie("HCOrganize")
	log.Println("获取orgainzeId: ", organizeId)
	// 根据部门ID查询部门关联的所有CameraID
	o := controllers.O

	var system []orm.Params
	num, err := o.Raw("select * from ss_system").Values(&system)

	if err == nil && num > 0 {
		c.Data["SystemTitle"] = system[0]["websitename"]
		log.Println("system:", system[0]["websitename"])
	} else {
		c.Data["SystemTitle"] = "AI风险防控预警平台"
	}

	var maps []orm.Params
	num1, err1 := o.Raw("select * from ss_organize where id=?", organizeId).Values(&maps)
	if err1 == nil && num1 > 0 {
		log.Println("查询个数：", num1)
		c.Data["organizeName"] = maps[0]["name"]
	}
	num, err = o.Raw("select camera_id from ss_camera where organizeid=?", organizeId).Values(&maps)
	if err != nil || num == 0 {
		log.Println("num:", num)
		c.Data["sum"] = 0
		c.Data["num"] = 0
		c.Data["unnum"] = 0
	} else {
		ids := ""
		idsv := []string{}
		for i := 0; i < len(maps); i++ {
			ids += ",?"
			idsv = append(idsv, maps[i]["camera_id"].(string))
		}
		ids = ids[1:len(ids)]
		//for i := 0; i < len(maps); i++ {
		//	if i != 0 {
		//		ids += ","
		//	}
		//	ids += "'"
		//	ids += maps[i]["camera_id"].(string)
		//	ids += "'"
		//}

		now := time.Now().Format("2006-01-02")
		startTime := now + " 00:00:00"
		endTime := now + " 23:59:59"
		var summaps []orm.Params
		sql1 := "select * from ss_alarm where alarm_time>=? and alarm_time<=? and camera_id in (" + ids + ")"
		sum, err := o.Raw(sql1, startTime, endTime, idsv).Values(&summaps)
		log.Println("sum2:", sum)
		var nummaps []orm.Params
		sql2 := "select * from ss_alarm where alarm_status=1 and alarm_time>=? and alarm_time<=? and camera_id in (" + ids + ")"
		num, err := o.Raw(sql2, startTime, endTime, idsv).Values(&nummaps)
		log.Println("num2:", num)
		var unsummaps []orm.Params
		sql3 := "select * from ss_alarm where alarm_status=0 and alarm_time>=? and alarm_time<=? and camera_id in (" + ids + ")"
		num_un, err := o.Raw(sql3, startTime, endTime, idsv).Values(&unsummaps)
		log.Println("num_un2", num_un)
		if err == nil {
			c.Data["sum"] = sum
			c.Data["num"] = num
			c.Data["unnum"] = num_un
		}
	}
	c.Data["Name"] = c.Ctx.GetCookie("HCName")
	c.TplName = "handle_comment/handle/workshop_unalarm.html"
}

func (w *WorkshopUnAlarmController2) Get() {
	organizeId := w.Ctx.GetCookie("HCOrganize")
	userId := w.Ctx.GetCookie("HCAdmin_id")
	log.Println("获取orgainzeId: ", organizeId)
	// 根据部门ID查询部门关联的所有CameraID
	o := controllers.O

	var system []orm.Params
	num, err := o.Raw("select * from ss_system").Values(&system)

	if err == nil && num > 0 {
		w.Data["SystemTitle"] = system[0]["websitename"]
		log.Println("system:", system[0]["websitename"])
	} else {
		w.Data["SystemTitle"] = "AI风险防控预警平台"
	}

	var maps []orm.Params
	num, err = o.Raw("select * from ss_organize where id=?", organizeId).Values(&maps)

	if err == nil && num > 0 {
		w.Data["organizeName"] = maps[0]["name"]
	}
	w.Data["Name"] = w.Ctx.GetCookie("HCName")

	var user models.User
	var cameras []models.Camera
	err = o.QueryTable(new(models.User)).Filter("admin_id", userId).RelatedSel().One(&user)
	if err != nil {
		log.Print("用户查询err:", err)
	}

	if user.Organize.ID == "" {
		w.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户组织架构异常"}
		w.ServeJSON()
		return
	}

	var org = user.Organize

	orgIds, err := controllers.SelOrgIDBelongUserOrg(user.Organize)
	if err != nil {
		controllers.LogsError("查询用户下属组织异常:", err)

	}
	var ids = []string{org.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}

	_, err = o.QueryTable(new(models.Camera)).Filter("Place__Organize__ID__in", ids).RelatedSel().All(&cameras)

	if err != nil {
		controllers.LogsError("查询设备信息异常:", err)
		w.Data["sum"] = 0
		w.Data["num"] = 0
		w.Data["unnum"] = 0
	} else {
		ids := ""
		idsv := []string{}
		for i := 0; i < len(maps); i++ {
			ids += ",?"
			idsv = append(idsv, cameras[i].ID)
		}
		ids = ids[1:len(ids)]
		//for i := 0; i < len(cameras); i++ {
		//	if i != 0 {
		//		ids += ","
		//	}
		//	ids += "'"
		//	ids += cameras[i].ID
		//	ids += "'"
		//}
		now := time.Now().Format("2006-01-02")
		startTime := now + " 00:00:00"
		endTime := now + " 23:59:59"
		var summaps []orm.Params
		sql1 := "select * from ss_alarm where alarm_time>=? and alarm_time<=? and camera_id in (" + ids + ")"
		sum, err := o.Raw(sql1, startTime, endTime, idsv).Values(&summaps)
		log.Println("sum2:", sum)
		var nummaps []orm.Params
		sql2 := "select * from ss_alarm where alarm_status=1 and alarm_time>=? and alarm_time<=? and camera_id in (" + ids + ")"
		num, err := o.Raw(sql2, startTime, endTime, idsv).Values(&nummaps)
		log.Println("num2:", num)
		var unsummaps []orm.Params
		sql3 := "select * from ss_alarm where alarm_status=0 and alarm_time>=? and alarm_time<=? and camera_id in (" + ids + ")"
		num_un, err := o.Raw(sql3, startTime, endTime, idsv).Values(&unsummaps)
		log.Println("num_un2", num_un)
		if err == nil {
			w.Data["sum"] = sum
			w.Data["num"] = num
			w.Data["unnum"] = num_un
		}
	}
	w.TplName = "handle_comment/handle/workshop_unalarm.html"
}

/**查询今日未处理的告警列表**/
func (c *WorkshopUnAlarmListController) Get() {
	now := time.Now().Format("2006-01-02")
	startTime := now + " 00:00:00"
	endTime := now + " 23:59:59"
	// 根据部门ID查询该部门所有的告警信息
	organizeId := c.Ctx.GetCookie("HCOrganize")
	log.Println("获取orgainzeId: ", organizeId)
	// 根据部门ID查询部门关联的所有CameraID
	o := controllers.O
	var maps []orm.Params
	num, err := o.Raw("select camera_id from ss_camera where organizeid=?", organizeId).Values(&maps)
	if err != nil || num == 0 {
		log.Println("num:", num)
		c.Data["json"] = map[string]interface{}{"Code": 0, "Reason": &maps}
	} else {
		ids := ""
		idsv := []string{}
		for i := 0; i < len(maps); i++ {
			ids += ",?"
			idsv = append(idsv, maps[i]["camera_id"].(string))
		}
		ids = ids[1:len(ids)]
		//for i := 0; i < len(maps); i++ {
		//	if i != 0 {
		//		ids += ","
		//	}
		//	ids += "'"
		//	ids += maps[i]["camera_id"].(string)
		//	ids += "'"
		//}
		limits, pages := c.GetString("limit"), c.GetString("page")
		start_page, err := strconv.Atoi(pages)

		limitss, err := strconv.Atoi(limits)
		start_pages := (start_page - 1) * limitss
		log.Println("", err)
		var results []orm.Params
		/*var ManyStream = []map[string]string{}*/
		log.Println("cameraId:", maps[0]["camera_id"])
		sql1 := "select alarm_id as id,alarm_detial as content,alarm_time as createtime,alarm_status as state, alarm_head as head, alarm_file as file, alarm_video as video, alarm_type as type from ss_alarm where camera_id in (" + ids + ") and alarm_status=0 and alarm_time>=? and alarm_time<=? order by alarm_time desc limit ?,?"
		num, err := o.Raw(sql1, idsv, startTime, endTime, int(start_pages), int(limitss)).Values(&results)
		log.Println("num:", num)
		if err == nil {
			sql2 := "select * from ss_alarm where camera_id in (" + ids + ") and alarm_status=0 and alarm_time>=? and alarm_time<=?"
			sum, err := o.Raw(sql2, idsv, startTime, endTime).Values(&maps)
			log.Println("err", err)
			log.Println("sum", sum)
			mystruct := &util.ResultB{0, sum, &results}
			c.Data["json"] = mystruct
		} else {
			c.Data["json"] = map[string]interface{}{"Code": 0, "Reason": &results}
		}
	}
	c.ServeJSON()
	return
}
func (w *WorkshopUnAlarmListController2) Get() {
	userId := w.Ctx.GetCookie("HCAdmin_id")
	now := time.Now().Format("2006-01-02")
	startTime := now + " 00:00:00"
	endTime := now + " 23:59:59"
	o := controllers.O
	var results []orm.Params
	var user models.User
	var cameras []models.Camera
	err := o.QueryTable(new(models.User)).Filter("admin_id", userId).RelatedSel().One(&user)
	if err != nil {
		log.Print("用户查询err:", err)
	}

	if user.Organize.ID == "" {
		w.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户组织架构异常"}
		w.ServeJSON()
		return
	}

	var org = user.Organize

	orgIds, err := controllers.SelOrgIDBelongUserOrg(user.Organize)
	if err != nil {
		controllers.LogsError("查询用户下属组织异常:", err)

	}
	var ids = []string{org.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}

	_, err = o.QueryTable(new(models.Camera)).Filter("Place__Organize__ID__in", ids).RelatedSel().All(&cameras)

	if err != nil {
		controllers.LogsError("查询设备信息异常:", err)
		w.Data["json"] = map[string]interface{}{"Code": 0, "Reason": &results}
		w.Data["sum"] = 0
		w.Data["num"] = 0
		w.Data["unnum"] = 0
	} else {
		ids := ""
		idsv := []string{}
		for i := 0; i < len(cameras); i++ {
			ids += ",?"
			idsv = append(idsv, cameras[i].ID)
		}
		ids = ids[1:len(ids)]
		//for i := 0; i < len(cameras); i++ {
		//	if i != 0 {
		//		ids += ","
		//	}
		//	ids += "'"
		//	ids += cameras[i].ID
		//	ids += "'"
		//}
		limits, pages := w.GetString("limit"), w.GetString("page")
		start_page, err := strconv.Atoi(pages)

		limitss, err := strconv.Atoi(limits)
		start_pages := (start_page - 1) * limitss
		log.Println("", err)

		sql1 := "select a.alarm_id as id,a.alarm_detial as content,a.alarm_time as createtime,a.alarm_status as state, a.alarm_head as head, a.alarm_file as file, a.alarm_video as video, a.alarm_type as type, sc.videoCode as videoCode from ss_alarm a inner join ss_camera sc on sc.camera_id = a.camera_id where a.camera_id in (" + ids + ") and a.alarm_status=0 and a.alarm_time>=? and a.alarm_time<=? order by a.alarm_time desc limit ?,?"
		num, err := o.Raw(sql1, idsv, startTime, endTime, int(start_pages), int(limitss)).Values(&results)
		log.Println("num:", num)
		if err == nil {
			var maps []orm.Params
			sql2 := "select * from ss_alarm where camera_id in (" + ids + ") and alarm_status=0 and alarm_time>=? and alarm_time<=?"
			sum, err := o.Raw(sql2, idsv, startTime, endTime).Values(&maps)
			log.Println("err", err)
			log.Println("sum", sum)
			mystruct := &util.ResultB{0, sum, &results}
			w.Data["json"] = mystruct
		} else {
			w.Data["json"] = map[string]interface{}{"Code": 0, "Reason": &results}
		}
	}
	w.ServeJSON()
	return
}

/**处理告警信息**/
func (c *WorkshopDealAlarmController) Get() {
	userId := c.Ctx.GetCookie("HCAdmin_id")
	alarmId, dealContent, alarmStatus := c.GetString("alarmId"), c.GetString("dealContent"), c.GetString("alarm_status")
	println("参数信息alarmId：", alarmId, "dealContent:", dealContent, "alarmStatus:", alarmStatus)
	o := controllers.O
	res, err := o.Raw("update ss_alarm set alarm_status=?,readstatus=1,manage_desc=?,user_id=? where alarm_id=?;", alarmStatus, dealContent, userId, alarmId).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		fmt.Println("mysql row affected nums: ", num)
		c.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "处理信息提交成功~"}
		c.ServeJSON()
	} else {
		fmt.Println("告警处理异常: ", err)
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "处理信息提交失败~"}
		c.ServeJSON()
	}
	return
}

func (w *WorkshopAlarmListController2) Get() {
	userId := w.Ctx.GetCookie("HCAdmin_id")
	o := controllers.O
	var results []orm.Params
	var maps []orm.Params
	var user models.User
	var cameras []models.Camera
	err := o.QueryTable(new(models.User)).Filter("admin_id", userId).RelatedSel().One(&user)
	if err != nil {
		log.Print("用户查询err:", err)
	}

	if user.Organize.ID == "" {
		w.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户组织架构异常"}
		w.ServeJSON()
		return
	}

	var org = user.Organize

	orgIds, err := controllers.SelOrgIDBelongUserOrg(user.Organize)
	if err != nil {
		controllers.LogsError("查询用户下属组织异常:", err)

	}
	var ids = []string{org.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}

	_, err = o.QueryTable(new(models.Camera)).Filter("Place__Organize__ID__in", ids).RelatedSel().All(&cameras)

	if err != nil {
		controllers.LogsError("查询设备信息异常:", err)
		w.Data["json"] = map[string]interface{}{"Code": 0, "Reason": &results}
		w.Data["sum"] = 0
		w.Data["num"] = 0
		w.Data["unnum"] = 0
	} else {
		ids := ""
		idsv := []string{}
		for i := 0; i < len(cameras); i++ {
			ids += ",?"
			idsv = append(idsv, cameras[i].ID)
		}
		ids = ids[1:len(ids)]
		//for i := 0; i < len(cameras); i++ {
		//	if i != 0 {
		//		ids += ","
		//	}
		//	ids += "'"
		//	ids += cameras[i].ID
		//	ids += "'"
		//}
		limits, pages, dateTime, status := w.GetString("limit"), w.GetString("page"), w.GetString("dateTime"), w.GetString("status")
		start_page, err := strconv.Atoi(pages)

		limitss, err := strconv.Atoi(limits)
		start_pages := (start_page - 1) * limitss
		log.Println("", err)

		log.Println("dateTime:", dateTime)
		log.Println("status:", status)

		times := strings.Split(dateTime, " - ")
		if len(dateTime) > 0 || len(times) > 1 {
			startTime := times[0] + " 00:00:00"
			endTime := times[1] + " 23:59:59"
			sql1 := ""
			if status != "-1" {
				sql1 = "select a.manage_desc as manageDesc, a.alarm_id as id,a.alarm_detial as content,a.alarm_time as createtime,a.alarm_status as state, a.alarm_head as head, a.alarm_file as file, a.alarm_video as video,a.alarm_type as type, sc.videoCode as videoCode from ss_alarm a inner join ss_camera sc on sc.camera_id = a.camera_id  where a.camera_id in (" + ids + ") and a.alarm_time>=? and a.alarm_time<=? and a.alarm_status=" + status + " order by a.alarm_time desc limit ?,?"
			} else {
				sql1 = "select a.manage_desc as manageDesc, a.alarm_id as id,a.alarm_detial as content,a.alarm_time as createtime,a.alarm_status as state, a.alarm_head as head, a.alarm_file as file, a.alarm_video as video,a.alarm_type as type, sc.videoCode as videoCode from ss_alarm a inner join ss_camera sc on sc.camera_id = a.camera_id  where a.camera_id in (" + ids + ") and a.alarm_time>=? and a.alarm_time<=? order by a.alarm_time desc limit ?,?"
			}
			num, err := o.Raw(sql1, idsv, startTime, endTime, int(start_pages), int(limitss)).Values(&results)
			log.Println("num:", num)
			if err == nil {
				sql2 := ""
				if status != "-1" {
					sql2 = "select * from ss_alarm where camera_id in (" + ids + ") and alarm_time>=? and alarm_time<=? and alarm_status=" + status
				} else {
					sql2 = "select * from ss_alarm where camera_id in (" + ids + ") and alarm_time>=? and alarm_time<=?"
				}
				sum, err := o.Raw(sql2, idsv, startTime, endTime).Values(&maps)
				log.Println("err", err)
				log.Println("sum", sum)
				mystruct := &util.ResultB{0, sum, &results}
				w.Data["json"] = mystruct
			} else {
				w.Data["json"] = map[string]interface{}{"Code": 0, "Reason": &results}
			}
		} else {
			sql1 := ""
			if status != "-1" {
				sql1 = "select a.manage_desc as manageDesc, a.alarm_id as id,a.alarm_detial as content,a.alarm_time as createtime,a.alarm_status as state, a.alarm_head as head, a.alarm_file as file, a.alarm_video as video,a.alarm_type as type, sc.videoCode as videoCode from ss_alarm a inner join ss_camera sc on sc.camera_id = a.camera_id where a.camera_id in (" + ids + ") and a.alarm_status=" + status + " order by a.alarm_time desc limit ?,?"
			} else {
				sql1 = "select a.manage_desc as manageDesc, a.alarm_id as id,a.alarm_detial as content,a.alarm_time as createtime,a.alarm_status as state, a.alarm_head as head, a.alarm_file as file, a.alarm_video as video,a.alarm_type as type, sc.videoCode as videoCode from ss_alarm a inner join ss_camera sc on sc.camera_id = a.camera_id where a.camera_id in (" + ids + ") order by a.alarm_time desc limit ?,?"
			}

			num, err := o.Raw(sql1, idsv, int(start_pages), int(limitss)).Values(&results)
			log.Println("num:", num)
			if err == nil {
				sql2 := ""
				if status != "-1" {
					sql2 = "select * from ss_alarm where camera_id in (" + ids + ") and alarm_status=" + status
				} else {
					sql2 = "select * from ss_alarm where camera_id in (" + ids + ")"
				}
				sum, err := o.Raw(sql2, idsv).Values(&maps)
				log.Println("err", err)
				log.Println("sum", sum)
				mystruct := &util.ResultB{0, sum, &results}
				w.Data["json"] = mystruct
			} else {
				w.Data["json"] = map[string]interface{}{"Code": 0, "Reason": &results}
			}
		}
	}
	w.ServeJSON()
	return
}

func (w *WorkshopNotationListController2) Get() {
	limits, pages, dateTime := w.GetString("limit"), w.GetString("page"), w.GetString("dateTime")
	times := strings.Split(dateTime, " - ")
	start_page, err := strconv.Atoi(pages)

	limitss, err := strconv.Atoi(limits)
	start_pages := (start_page - 1) * limitss

	userId := w.Ctx.GetCookie("HCAdmin_id")
	o := controllers.O
	var results []orm.Params
	var maps []orm.Params
	var user models.User
	var cameras []models.Camera
	err = o.QueryTable(new(models.User)).Filter("admin_id", userId).RelatedSel().One(&user)
	if err != nil {
		log.Print("用户查询err:", err)
	}

	if user.Organize.ID == "" {
		w.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户组织架构异常"}
		w.ServeJSON()
		return
	}

	var org = user.Organize

	orgIds, err := controllers.SelOrgIDBelongUserOrg(user.Organize)
	if err != nil {
		controllers.LogsError("查询用户下属组织异常:", err)

	}
	var ids = []string{org.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}

	_, err = o.QueryTable(new(models.Camera)).Filter("Place__Organize__ID__in", ids).RelatedSel().All(&cameras)

	if err != nil {
		controllers.LogsError("查询设备信息异常:", err)
		w.Data["json"] = map[string]interface{}{"Code": 0, "Reason": &results}
		w.Data["sum"] = 0
		w.Data["num"] = 0
		w.Data["unnum"] = 0
	} else {
		ids := ""
		idsv := []string{}
		for i := 0; i < len(cameras); i++ {
			ids += ",?"
			idsv = append(idsv, cameras[i].ID)
		}
		ids = ids[1:len(ids)]
		//for i := 0; i < len(cameras); i++ {
		//	if i != 0 {
		//		ids += ","
		//	}
		//	ids += "'"
		//	ids += cameras[i].ID
		//	ids += "'"
		//}
		if len(dateTime) > 0 || len(times) > 1 {
			startTime := times[0] + " 00:00:00"
			endTime := times[1] + " 23:59:59"
			alarm_sql := "SELECT a.alarm_id AS id, a.alarm_detial AS content, a.alarm_time AS alarmtime, a.user_id AS dealid, a.manage_desc AS dealcontent, b.admin_name AS dealName, a.alarm_head as head, a.alarm_file as file, a.alarm_video as video, sc.videoCode as videoCode FROM ss_alarm a left join ss_user b ON a.user_id=b.admin_id inner join ss_camera sc on sc.camera_id = a.camera_id WHERE a.camera_id IN (" + ids + ") AND a.alarm_status=1 and notationstatus=1 and alarm_time>=? and alarm_time<=? limit ?,?"
			num, err := o.Raw(alarm_sql, idsv, startTime, endTime, int(start_pages), int(limitss)).Values(&maps)
			if err == nil {
				log.Println("已批注个数：", num)
				mystruct := &util.ResultB{0, num, &maps}
				w.Data["json"] = mystruct
			} else {
				w.Data["json"] = map[string]interface{}{"Code": 0, "Reason": &maps}
			}
		} else {
			alarm_sql := "SELECT a.alarm_id AS id, a.alarm_detial AS content, a.alarm_time AS alarmtime, a.user_id AS dealid, a.manage_desc AS dealcontent, b.admin_name AS dealName, a.alarm_head as head, a.alarm_file as file, a.alarm_video as video, sc.videoCode as videoCode FROM ss_alarm a left join ss_user b ON a.user_id=b.admin_id inner join ss_camera sc on sc.camera_id = a.camera_id WHERE a.camera_id IN (" + ids + ") AND a.alarm_status=1 and notationstatus=1 limit ?,?"
			num, err := o.Raw(alarm_sql, idsv, int(start_pages), int(limitss)).Values(&maps)
			if err == nil {
				log.Println("已批注个数：", num)
				mystruct := &util.ResultB{0, num, &maps}
				w.Data["json"] = mystruct
			} else {
				w.Data["json"] = map[string]interface{}{"Code": 0, "Reason": &maps}
			}
		}
	}
	w.ServeJSON()
	return
}

type WorkshopNotationController struct {
	beego.Controller
}
type WorkshopNotationController2 struct {
	beego.Controller
}

func (c *WorkshopNotationController) Get() {
	organizeId := c.Ctx.GetCookie("HCOrganize")
	o := controllers.O

	//var system []orm.Params
	//num, err := o.Raw("select * from ss_system").Values(&system)
	//
	//if err == nil && num > 0 {
	//	c.Data["SystemTitle"] = system[0]["websitename"]
	//	log.Println("system:", system[0]["websitename"])
	//} else {
	//	c.Data["SystemTitle"] = "AI风险防控预警平台"
	//}

	var maps []orm.Params
	num1, err1 := o.Raw("select * from ss_organize where id=?", organizeId).Values(&maps)
	if err1 == nil && num1 > 0 {
		log.Println("查询个数：", num1)
		c.Data["organizeName"] = maps[0]["name"]
	}
	c.Data["Name"] = c.Ctx.GetCookie("Name")
	c.TplName = "handle_comment/handle/workshop_notation.html"
}

/*告警批注列表*/
type WorkshopAlarmNotationListController struct {
	beego.Controller
}

func (c *WorkshopAlarmNotationListController) Post() {
	// 获取告警编号
	alarmId := c.GetString("alarmId")
	o := controllers.O
	var maps []orm.Params
	num, err := o.Raw("select o.notation as content,o.createtime,u.admin_name as notationName from ss_notation o left join ss_user u on o.userid=u.admin_id where alarm_id=? order by createtime desc", alarmId).Values(&maps)
	if err == nil {
		mystruct := &util.ResultB{0, num, &maps}
		c.Data["json"] = mystruct
	} else {
		c.Data["json"] = map[string]interface{}{"Code": -1, "Reason": &maps}
	}
	c.ServeJSON()
	return
}
