package handle_comment

import (
	"Artifice_V2.0.0/controllers"
	"Artifice_V2.0.0/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	"strconv"
	"strings"
	"time"
)

type EventPageHandler struct {
	beego.Controller
}

func (this *EventPageHandler) Get() {
	organizeId := this.Ctx.GetCookie("HCOrganize")
	o := controllers.O
	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	var maps []orm.Params
	num, err := o.Raw("select * from ss_organize where id=?", organizeId).Values(&maps)
	if err == nil && num > 0 {
		log.Println("查询个数：", num)
		this.Data["organizeName"] = maps[0]["name"]
	}
	this.Data["Name"] = this.Ctx.GetCookie("HCName")
	this.TplName = "handle_comment/eventpage.html"
}

type SelHcAlarmListHandler struct {
	beego.Controller
}

func (this *SelHcAlarmListHandler) Post() {
	limits, pages, dateTime := this.GetString("limit"), this.GetString("page"), this.GetString("dateTime")
	times := strings.Split(dateTime, " - ")
	start_page, err := strconv.Atoi(pages)
	limitss, err := strconv.Atoi(limits)
	start_pages := (start_page - 1) * limitss
	o := controllers.O
	organizeId := this.Ctx.GetCookie("HCOrganize")

	var maps []orm.Params
	org_sql := "select id from (select t1.id, if(find_in_set(pId, @pids) > 0, @pids := concat(@pids, ',', id), 0) as ischild from (select id,pId from ss_organize t order by pId, id) t1, (select @pids := " + organizeId + ") t2) t3 where ischild != 0"
	num, err := o.Raw(org_sql).Values(&maps)
	if err == nil {
		log.Print("下属部门个数：", num)
		var ids = []string{organizeId}
		for _, v := range maps {
			ids = append(ids, v["id"].(string))
		}
		orgids := ""
		for i := 0; i < len(ids); i++ {
			orgids += ",?"
		}
		orgids = orgids[1:len(orgids)]
		//orgids += "'" + organizeId + "'"
		//for i := 0; i < len(maps); i++ {
		//	orgids += ",'" + maps[i]["id"].(string) + "'"
		//}
		var placemaps []orm.Params
		num, err = o.Raw("select place_id from ss_place where organize_id in ("+orgids+")", ids).Values(&placemaps)
		if err != nil {
			// 组织查询区域失败
		} else {
			if num == 0 {
				// 组织下无区域
			} else {
				placeids := ""
				placeidsv := []string{}
				for j := 0; j < len(placemaps); j++ {
					//placeids += ",'" + placemaps[j]["place_id"].(string) + "'"
					placeids += ",?"
					placeidsv = append(placeidsv, placemaps[j]["place_id"].(string))
				}
				placeids = placeids[1:len(placeids)]
				var cameramaps []orm.Params
				num, err = o.Raw("select camera_id from ss_camera where place_id in ("+placeids+")", placeidsv).Values(&cameramaps)
				if err != nil {
					// 查询摄像头失败
				} else {
					if num == 0 {
						// 区域下无摄像头
					} else {
						cameraids := ""
						cameraidsv := []string{}
						for j := 0; j < len(cameramaps); j++ {
							//cameraids += ",'" + cameramaps[j]["camera_id"].(string) + "'"
							cameraids += ",?"
							cameraidsv = append(cameraidsv, cameramaps[j]["camera_id"].(string))
						}
						cameraids = cameraids[1:len(cameraids)]
						startTime := time.Now().Format("2006-01-02") + " 00:00:00"
						endTime := time.Now().Format("2006-01-02") + " 23:59:59"
						if len(dateTime) > 0 || len(times) > 1 {
							startTime = times[0] + " 00:00:00"
							endTime = times[1] + " 23:59:59"
						}
						alarm_sql := "SELECT a.alarm_id AS id, a.alarm_detial AS content, a.alarm_time AS alarmtime, a.user_id AS dealid, a.manage_desc AS dealcontent,a.notationstatus, a.alarm_head as head, a.alarm_file as file, a.alarm_video as video, b.admin_name AS dealName, sc.videoCode as videoCode FROM ss_alarm a left join ss_user b ON a.user_id=b.admin_id inner join ss_camera sc on sc.camera_id = a.camera_id WHERE a.camera_id IN (" + cameraids + ") AND a.alarm_status=1 and alarm_time>=? and alarm_time<=? limit ?,?"
						num, err := o.Raw(alarm_sql, cameraidsv, startTime, endTime, int(start_pages), int(limitss)).Values(&maps)
						if err == nil {
							log.Println("已批注个数：", num)
							mystruct := &util.ResultB{0, num, &maps}
							this.Data["json"] = mystruct
						} else {
							log.Println("获取批注信息异常：", err)
							this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": &maps}
						}
					}
				}
			}
		}
	}

	this.ServeJSON()
}

type SelAlarmNotationListHandler struct {
	beego.Controller
}

func (this *SelAlarmNotationListHandler) Post() {
	alarmId := this.GetString("alarmId")
	o := controllers.O
	var maps []orm.Params
	num, err := o.Raw("select o.notation as content,o.createtime,u.admin_name as notationName from ss_notation o left join ss_user u on o.userid=u.admin_id where alarm_id=? order by createtime desc", alarmId).Values(&maps)
	if err == nil {
		mystruct := &util.ResultB{0, num, &maps}
		this.Data["json"] = mystruct
	} else {
		this.Data["json"] = map[string]interface{}{"Code": -1, "Reason": &maps}
	}
	this.ServeJSON()
}

type SaveAlarmNotationHandler struct {
	beego.Controller
}

func (c *SaveAlarmNotationHandler) Post() {
	alarmId, notationContent := c.GetString("alarmId"), c.GetString("notationContent")
	userId := c.Ctx.GetCookie("HCAdmin_id")
	organizeId := c.Ctx.GetCookie("HCOrganize")
	log.Println("userId:", userId, " organizdeid:", organizeId)
	o := controllers.O
	var maps []orm.Params
	org_sql := "select id from (select t1.id, if(find_in_set(pId, @pids) > 0, @pids := concat(@pids, ',', id), 0) as ischild from (select id,pId from ss_organize t order by pId, id) t1, (select @pids := " + organizeId + ") t2) t3 where ischild != 0"
	num, err := o.Raw(org_sql).Values(&maps)
	if err == nil {
		log.Println("下属部门个数：", num)
		orgids := "" // 批注人的部门及下属部门
		orgids += "'" + organizeId + "'"
		for i := 0; i < len(maps); i++ {
			orgids += ",'" + maps[i]["id"].(string) + "'"
		}
		log.Println("下属部门:", maps, num)
		// 查询已批注的部门
		log.Println("查询已批注的部门前提:", maps)
		num, err := o.Raw("select o.id from ss_notation n left join ss_user u on n.userid=u.admin_id left join ss_organize o on u.organize_id=o.id where n.alarm_id=? order by n.createtime desc;", alarmId).Values(&maps)
		log.Println("查询已批注的部门:", maps, num)
		if err == nil {
			if num == 0 || (num == 1 && maps[0]["id"] == nil) { // 未批注
				res, err := o.Raw("insert into ss_notation(alarm_id,notation,userid,createtime) values(?,?,?,?);", alarmId, notationContent, userId, time.Now().Format("2006-01-02 15:04:05")).Exec()
				if err != nil {
					log.Println("stmt not error:", res)
					c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败"}
					c.ServeJSON()
					return
				}
				res1, err1 := o.Raw("update ss_alarm set notationstatus=1 where alarm_id=?", alarmId).Exec()
				if err1 != nil {
					log.Println("stmt not error:", res1)
					c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败"}
					c.ServeJSON()
					return
				}
				c.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "操作成功"}
				c.ServeJSON()
			} else {
				oid := maps[0]["id"].(string)
				if strings.Index(orgids, oid) >= 0 { // 包含
					res, err := o.Raw("insert into ss_notation(alarm_id,notation,userid,createtime) values(?,?,?,?);", alarmId, notationContent, userId, time.Now().Format("2006-01-02 15:04:05")).Exec()
					if err != nil {
						log.Println("stmt not error:", res)
						c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败"}
						c.ServeJSON()
						return
					}
					res1, err1 := o.Raw("update ss_alarm set notationstatus=1 where alarm_id=?", alarmId).Exec()
					if err1 != nil {
						log.Println("stmt not error:", res1)
						c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败"}
						c.ServeJSON()
						return
					}
					c.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "操作成功"}
					c.ServeJSON()
				} else {
					c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "批注失败，上级已做批注"}
					c.ServeJSON()
					return
				}
			}
		}
	}
}
