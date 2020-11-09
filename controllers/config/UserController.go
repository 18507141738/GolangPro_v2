package config

import (
	"Artifice_V2.0.0/controllers"
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
)

// 修改密码
type UpdateUserPasswordHandler struct {
	beego.Controller
}

func (this *UpdateUserPasswordHandler) Post() {
	admin_password, admin_newPwd := this.GetString("oldPwd"), this.GetString("newPwd")
	h := md5.New()
	h.Write([]byte(admin_password))
	md5pro := hex.EncodeToString(h.Sum(nil))
	admin_id := this.Ctx.GetCookie("admin_id")
	o := controllers.O
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

// 退出登录
type LogoutUserHandler struct {
	beego.Controller
}

func (this *LogoutUserHandler) Post() {
	admin_id := this.Ctx.GetCookie("admin_id")
	o := controllers.O
	res, err := o.Raw("delete from ss_sessions where userid=?", admin_id).Exec()
	if err != nil {
		log.Println(res, err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "操作失败"}
	} else {
		this.Ctx.SetCookie("admin_id", "")
		this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "修改成功"}
	}
	this.ServeJSON()
}
