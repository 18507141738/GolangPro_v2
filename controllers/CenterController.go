package controllers

import (
	"Artifice_V2.0.0/models"
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
)

//个人中心页面
type PlatCenterPage struct {
	beego.Controller
}

func (this *PlatCenterPage) Get() {
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
		this.TplName = "platform/centerPage.html"
	}
}

// 设备管理页面
type PlatformHostManagePage struct {
	beego.Controller
}

func (this *PlatformHostManagePage) Get() {
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
		this.TplName = "platform/hostlist.html"
	}
}

//成员管理
type PlatformUserMangerPageController struct {
	beego.Controller
}

func (this *PlatformUserMangerPageController) Get() {
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
		this.TplName = "platform/userlist.html"
	}
}

// 修改密码
type UpdateUserPasswordHandler struct {
	beego.Controller
}

func (this *UpdateUserPasswordHandler) Post() {
	admin_password, admin_newPwd := this.GetString("oldPwd"), this.GetString("newPwd")
	h := md5.New()
	h.Write([]byte(admin_password))
	md5pro := hex.EncodeToString(h.Sum(nil))
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var maps []orm.Params
	num, err := o.Raw("select * from ss_user where admin_id=?", admin_id).Values(&maps)
	log.Println(num)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "操作失败，请重新登录后再修改"}
	} else {
		if maps[0]["admin_password"] != md5pro {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "操作失败，密码错误"}
		} else {
			h = md5.New()
			h.Write([]byte(admin_newPwd))
			newmd5pro := hex.EncodeToString(h.Sum(nil))
			res, err := o.Raw("update ss_user set admin_password=? where admin_id=?", newmd5pro, admin_id).Exec()
			if err != nil {
				log.Println(res, err)
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "操作失败"}
			} else {
				this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "修改成功"}
			}
		}
	}
	this.ServeJSON()
}
