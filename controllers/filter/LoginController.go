package filter

import (
	"Artifice_V2.0.0/controllers"
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type FilterLoginPage struct {
	beego.Controller
}

func (this *FilterLoginPage) Get() {
	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	this.Data["version"] = controllers.GetVer()
	this.Data["updateTime"] = controllers.GetUpdateTime()
	this.Data["buildTime"] = controllers.GetBuildTime()
	this.TplName = "filter/login.html"
}

type FilterLoginCheckHandler struct {
	beego.Controller
}

func (this *FilterLoginCheckHandler) Post() {
	username, password := this.GetString("Username"), this.GetString("Password")
	h := md5.New()
	h.Write([]byte(password))
	md5pro := hex.EncodeToString(h.Sum(nil))
	o := controllers.O
	var maps []orm.Params
	num, err := o.Raw("select * from ss_filter_user where username=?", username).Values(&maps)
	if err != nil {
		controllers.LogsError("用户登录过滤网平台报错：", err)
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": "登录失败！"}
	} else {
		if num == 0 {
			this.Data["json"] = map[string]interface{}{"code": -1, "msg": "登录失败，账号不存在！"}
		} else {
			if maps[0]["password"].(string) != md5pro {
				this.Data["json"] = map[string]interface{}{"code": -1, "msg": "登录失败，密码错误！"}
			} else {
				uid := uuid.Must(uuid.NewV4()).String()
				uid = strings.Replace(uid, "-", "", -1) // 生成SessionID
				var sessions []orm.Params
				num, err = o.Raw("select * from ss_sessions where userid=?", maps[0]["id"]).Values(&sessions)
				if err != nil {
					controllers.LogsError("用户登录过滤网平台Session报错：", err)
					this.Data["json"] = map[string]interface{}{"code": -1, "msg": "登录失败！"}
				} else {
					if num == 0 {
						_, err := o.Raw("insert into ss_sessions (id, userid, sessionid) values (?,?,?)",
							strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1), maps[0]["id"], uid).Exec()
						if err != nil {
							controllers.LogsError("用户登录过滤网平台Session报错：", err)
							this.Data["json"] = map[string]interface{}{"code": -1, "msg": "登录失败！"}
						} else if username == "admin" { // admin登录
							this.Data["json"] = map[string]interface{}{"code": 1, "msg": "登录成功"}
						} else {
							this.Data["json"] = map[string]interface{}{"code": 0, "msg": "登录成功"}
						}
					} else {
						_, err := o.Raw("update ss_sessions set sessionid=? where userid=?", uid, maps[0]["id"]).Exec()
						if err != nil {
							controllers.LogsError("用户登录过滤网平台Session报错：", err)
							this.Data["json"] = map[string]interface{}{"code": -1, "msg": "登录失败！"}
						} else if username == "admin" { // admin登录
							this.Data["json"] = map[string]interface{}{"code": 1, "msg": "登录成功"}
						} else {
							this.Data["json"] = map[string]interface{}{"code": 0, "msg": "登录成功"}
						}
					}
					this.Ctx.SetCookie("filter_session_id", uid)
					this.Ctx.SetCookie("filter_admin_id", maps[0]["id"].(string))
				}
			}
		}
	}
	this.ServeJSON()
}
