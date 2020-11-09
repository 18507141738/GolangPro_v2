package filter

import (
	"Artifice_V2.0.0/controllers"
	"Artifice_V2.0.0/models"
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"log"
	"strings"
)

type FilterUserPage struct {
	beego.Controller
}

func (this *FilterUserPage) Get() {
	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	this.Data["version"] = controllers.GetVer()
	this.Data["updateTime"] = controllers.GetUpdateTime()
	this.Data["buildTime"] = controllers.GetBuildTime()
	this.Data["bmk"] = "filterUser"
	this.TplName = "filter/admin/user.html"
}

type FilterSelOrgHandler struct {
	beego.Controller
}

func (this *FilterSelOrgHandler) Post() {
	var maps []orm.Params
	o := controllers.O
	num, err := o.Raw("select * from ss_organize;").Values(&maps)
	if err != nil {
		log.Println(err)
		this.Data["json"] = map[string]interface{}{"Code": 0, "reason": "数据操作失败!,再来一次~"}
	} else {

		for _, v := range maps {
			v["isParent"] = "true"
			v["icon"] = "/static/css/zTreeStyle/img/diy/homeIcon.png"
			v["Level"] = v["level"]
		}

		mystruct := map[string]interface{}{"Code": 1, "Num": num, "Reason": &maps}
		this.Data["json"] = mystruct
	}
	this.ServeJSON()
}

type FilterSelUserHandler struct {
	beego.Controller
}

func (this *FilterSelUserHandler) Post() {
	orgid, userid := this.GetString("organize_id"), this.Ctx.GetCookie("admin_id")
	fs, err := models.GetFilterUserByOrgid(controllers.O, orgid, userid)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "查询组织成员失败"}
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"Ret": 0, "num": len(fs), "Reason": &fs}
	this.ServeJSON()

}

type FilterAddUserHandler struct {
	beego.Controller
}

func (this *FilterAddUserHandler) Post() {
	oid, acount, pass, jur :=
		this.GetString("member_organizeId"),
		this.GetString("member_user"),
		this.GetString("member_pass"),
		"0,1,2,3,4,5"

	if len(oid) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "未选择组织架构"}
		this.ServeJSON()
		return
	}
	if len(acount) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "登录名不能为空"}
		this.ServeJSON()
		return
	}
	if len(pass) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "密码不能为空"}
		this.ServeJSON()
		return
	}

	if models.FilterExist(controllers.O, acount) {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "添加员工账号已存在"}
		this.ServeJSON()
		return
	}

	uid := uuid.Must(uuid.NewV4()).String()
	uid = strings.Replace(uid, "-", "", -1) // 生成SessionID

	h := md5.New()
	h.Write([]byte(pass))
	md5Pass := hex.EncodeToString(h.Sum(nil)) //加密密码

	var fu models.FilterUser
	fu.ID = uid
	fu.Acount = acount
	fu.Pass = md5Pass
	fu.Jurisdiction = jur
	var org models.Organize
	org.ID = oid
	fu.Organize = &org

	err := models.AddFilterUser(controllers.O, fu)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "添加用户异常"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "添加成功"}
	this.ServeJSON()
}

type FilterUpdateUserHandler struct {
	beego.Controller
}

func (this *FilterUpdateUserHandler) Post() {
	id, acount, pass :=
		this.GetString("member_id"),
		this.GetString("member_user"),
		this.GetString("member_pass")

	if len(acount) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "登录名不能为空"}
		this.ServeJSON()
		return
	}

	if models.FilterExist(controllers.O, acount) {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "添加员工账号已存在"}
		this.ServeJSON()
		return
	}

	fu := models.FilterUser{ID: id}
	fu.Acount = acount
	fu.Pass = pass
	err := models.UpdateFilterUser(controllers.O, fu)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "更新用户信息异常"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "更新用户信息成功"}
	this.ServeJSON()
}

type FilterDelUserHandler struct {
	beego.Controller
}

func (this *FilterDelUserHandler) Post() {
	id := this.GetString("member_id")
	fu := models.FilterUser{ID: id}
	err := models.DelFilterUser(controllers.O, fu)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "删除该用户用户信息异常"}
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "删除成功"}
	this.ServeJSON()
}
