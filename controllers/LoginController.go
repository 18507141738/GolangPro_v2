package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"log"
	"strings"
)

// 登录页面
type LoginPage struct {
	beego.Controller
}

func (this *LoginPage) Get() {
	this.Data["SystemTitle"] = SelSystemTitleUtil()
	this.Data["version"] = GetVer()
	this.Data["updateTime"] = GetUpdateTime()
	this.Data["buildTime"] = GetBuildTime()
	this.TplName = "login.html"
}

// 登录验证
type LoginCheckHandler struct {
	beego.Controller
}

func (this *LoginCheckHandler) Post() {
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
			uid := uuid.Must(uuid.NewV4()).String()
			uid = strings.Replace(uid, "-", "", -1) // 生成SessionID
			var sessions []orm.Params
			num, err = o.Raw("select * from ss_sessions where userid=?", maps[0]["admin_id"]).Values(&sessions)
			if err != nil {
				this.Data["json"] = map[string]interface{}{"code": -1, "msg": "登录失败"}
			} else {
				if num == 0 {
					res, err := o.Raw("insert into ss_sessions (id, userid, sessionid) values (?,?,?)",
						strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1), maps[0]["admin_id"], uid).Exec()
					log.Println(err, res)
				} else {
					res, err := o.Raw("update ss_sessions set sessionid=? where userid=?", uid, maps[0]["admin_id"]).Exec()
					log.Println(err, res)
				}

				if maps[0]["type"] == "0" {
					this.Data["json"] = map[string]interface{}{"code": 0, "msg": "登录成功"} // 配置管理员
					// 保存cookie
					this.Ctx.SetCookie("admin_session", uid)
					this.Ctx.SetCookie("admin_id", maps[0]["admin_id"].(string))
					this.Ctx.SetCookie("admin_type", maps[0]["type"].(string))
				} else {
					if maps[0]["jurisdiction"] == nil || maps[0]["jurisdiction"] == "" {
						this.Data["json"] = map[string]interface{}{"code": -1, "msg": "账号无权限"}
					} else {
						this.Data["json"] = map[string]interface{}{"code": 1, "msg": "登录成功"} //
						// 保存cookie
						this.Ctx.SetCookie("user_session", uid)
						this.Ctx.SetCookie("user_id", maps[0]["admin_id"].(string))
						this.Ctx.SetCookie("user_type", maps[0]["type"].(string))
					}
				}
			}
		}
	} else {
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": "没有此用户"}
	}
	this.ServeJSON()
}

// 退出登录
type LogoutUserHandler struct {
	beego.Controller
}

func (this *LogoutUserHandler) Post() {
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	res, err := o.Raw("delete from ss_sessions where userid=?", admin_id).Exec()
	if err != nil {
		log.Println(res, err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "操作失败"}
	} else {
		this.Ctx.SetCookie("user_id", "")
		this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "修改成功"}
	}
	this.ServeJSON()
}
