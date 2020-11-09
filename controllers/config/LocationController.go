package config

import (
	"Artifice_V2.0.0/controllers"
	"Artifice_V2.0.0/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"log"
	"strings"
	"time"
)

// 区域管理页面
type LocationPageHandler struct {
	beego.Controller
}

func (this *LocationPageHandler) Get() {
	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	this.Data["bmk"] = "location_menu"
	this.TplName = "config/location.html"
}

// 查询区域列表
type SelLocationListHandler struct {
	beego.Controller
}

func (this *SelLocationListHandler) Post() {
	o := controllers.O
	var maps []orm.Params
	num, err := o.Raw("select p.place_id,p.place_name,p.createtime,o.name,o.id as org_id from ss_place p, ss_organize o where p.organize_id=o.id").Values(&maps)
	if err != nil {
		log.Println(err)
	}
	mystruct := &util.ResultB{0, num, &maps}
	this.Data["json"] = mystruct
	this.ServeJSON()
}

// 添加修改区域
type LocationAddController struct {
	beego.Controller
}

func (c *LocationAddController) Post() {
	code := c.GetString("code")
	place_id := c.GetString("place_id")
	if code == "0" {
		o := controllers.O
		uid := uuid.Must(uuid.NewV4()).String()
		uid = strings.Replace(uid, "-", "", -1)
		_, err := o.Raw("insert into ss_place(place_id,place_name,createtime,organize_id) values(?,?,?,?);", uid, c.GetString("placename"), time.Now().Format("2006-01-02 15:04:05"), c.GetString("organize_id")).Exec()
		if err == nil {
			c.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "保存成功"}
			c.ServeJSON()
		} else {
			log.Println("添加区域异常:", err)
			c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "保存失败"}
			c.ServeJSON()
		}
		return
	}
	if code == "1" {
		o := controllers.O
		_, err := o.Raw("update ss_place set place_name=?,organize_id=? where place_id=?", c.GetString("placename"), c.GetString("organize_id"), place_id).Exec()
		if err == nil {
			c.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "修改成功"}
			c.ServeJSON()
		} else {
			log.Println("修改失败", err)
			c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "修改失败"}
			c.ServeJSON()
		}
		return
	}
}

// 删除区域
type DelLocationDataHandler struct {
	beego.Controller
}

func (this *DelLocationDataHandler) Post() {
	params := this.GetString("param")
	paramsa := strings.Split(params, ",")
	//： 执行删除语句
	o := controllers.O
	var maps []orm.Params
	for _, mains := range paramsa {
		num, err := o.Raw("select * from ss_place where place_id=?", mains).Values(&maps)
		if err == nil && num > 0 {
			num, err = o.Raw("select * from ss_func_camera where location =?", maps[0]["place_name"]).Values(&maps)
			if err == nil && num > 0 {
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "请先移除该地点注册的所有相机,再来试试~"}
				this.ServeJSON()
				return
			} else {
				_, state := o.Raw("delete from ss_place where place_id =?", mains).Values(&maps)
				log.Println(state)

			}
		}
	}
	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "删除成功"}
	this.ServeJSON()
}
