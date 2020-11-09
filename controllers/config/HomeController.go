package config

import (
	"Artifice_V2.0.0/controllers"
	"Artifice_V2.0.0/models"
	"github.com/astaxie/beego"
	"log"
)

type ConfigHomePage struct {
	beego.Controller
}

func (this *ConfigHomePage) Get() {
	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	admin_id := this.Ctx.GetCookie("admin_id")
	o := controllers.O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).One(&user)
	if err != nil {
		log.Println(err)
		this.TplName = "login.html"
	} else {
		this.Data["admin_name"] = user.Name
		this.Data["admin_user"] = user.Acount
		this.TplName = "config/index.html"
	}
}
