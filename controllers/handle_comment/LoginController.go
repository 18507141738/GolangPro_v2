package handle_comment

import (
	"Artifice_V2.0.0/controllers"
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/satori/go.uuid"
	"log"
	"strings"
)

//登录页
type HCLoginPage struct {
	beego.Controller
}

func (this *HCLoginPage) Get() {
	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	this.Data["version"] = controllers.GetVer()
	this.Data["updateTime"] = controllers.GetUpdateTime()
	this.Data["buildTime"] = controllers.GetBuildTime()
	this.TplName = "handle_comment/login.html"
}

//登录请求
type HCLoginCheck struct {
	beego.Controller
}

func (this *HCLoginCheck) Post() {
	username, password := this.GetString("Username"), this.GetString("Password")
	h := md5.New()
	h.Write([]byte(password))
	md5pro := hex.EncodeToString(h.Sum(nil))
	o := controllers.O
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

				// 保存cookie
				this.Ctx.SetCookie("HC_session", uid)
				this.Ctx.SetCookie("HCAdmin_id", maps[0]["admin_id"].(string))
				this.Ctx.SetCookie("HCAdmin_type", maps[0]["type"].(string))
				this.Ctx.SetCookie("HCOrganize", maps[0]["organize_id"].(string))
				this.Ctx.SetCookie("HCName", maps[0]["admin_name"].(string))

				if maps[0]["type"] == "1" {
					this.Data["json"] = map[string]interface{}{"code": 1, "msg": "登录成功"} //
				} else if maps[0]["type"] == "2" {
					this.Data["json"] = map[string]interface{}{"code": 2, "msg": "登录成功"} //
				} else {
					this.Data["json"] = map[string]interface{}{"code": -1, "msg": "登录失败"} //
				}
			}
		}
	} else {
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": "没有此用户"}
	}
	this.ServeJSON()
}
