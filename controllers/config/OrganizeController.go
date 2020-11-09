package config

import (
	"Artifice_V2.0.0/controllers"
	"Artifice_V2.0.0/models"
	"Artifice_V2.0.0/util"
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
	"strings"
	"time"
)

// 区域查询组织架构
type SelOrganizeListHandler struct {
	beego.Controller
}

func (this *SelOrganizeListHandler) Post() {
	o := controllers.O
	var maps []orm.Params
	num, err := o.Raw("select * from ss_organize;").Values(&maps)
	log.Println(">>>", num)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败!,再来一次~"}
		this.ServeJSON()
		return
	}
	for i := 0; i < len(maps); i++ {
		if maps[i]["level"] == "4" {
			maps[i]["isParent"] = "false"
		} else {
			maps[i]["isParent"] = "true"
		}
	}
	mystruct := &util.ResultB{1, num, &maps}
	this.Data["json"] = mystruct
	this.ServeJSON()
}

type OrganizePageHandler struct {
	beego.Controller
}

func (this *OrganizePageHandler) Get() {
	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	this.Data["bmk"] = "orgmenu"
	this.TplName = "config/organize.html"
}

type SelOrganizeTreeDataHandler struct {
	beego.Controller
}

func (c *SelOrganizeTreeDataHandler) Post() {
	o := controllers.O
	var maps []orm.Params
	num, err := o.Raw("select * from ss_organize;").Values(&maps)
	if err != nil {
		log.Println(">>>", num)
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败!,再来一次~"}
		c.ServeJSON()
		return
	}
	for i := 0; i < len(maps); i++ {
		if maps[i]["level"] == "4" {
			maps[i]["isParent"] = "false"
		} else {
			maps[i]["isParent"] = "true"
		}
	}
	mystruct := &util.ResultB{1, num, &maps}
	c.Data["json"] = mystruct
	c.ServeJSON()
}

type SelUserForOrgIdHandler struct {
	beego.Controller
}

func (c *SelUserForOrgIdHandler) Post() {
	organize_id := c.GetString("organize_id")
	log.Println(organize_id)
	o := controllers.O
	var maps []orm.Params
	num, err := o.Raw("select * from ss_user where organize_id = ?;", organize_id).Values(&maps)
	if err != nil {
		log.Println("查询成员异常", num, err)
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "查询成员异常"}
		c.ServeJSON()
		return
	}
	for i := 0; i < len(maps); i++ {
		var qmaps []orm.Params
		num, err := o.Raw("select * from ss_jurisdiction where id =?;", organize_id).Values(&qmaps)
		if err == nil && num > 0 {
			log.Println(num)
			maps[i]["jur_type"] = qmaps[0]["type"]
			maps[i]["jur_dev"] = qmaps[0]["device"]
			maps[i]["jur_timeslot"] = qmaps[0]["time_slot"]
		}
		maps[i]["name"] = maps[i]["admin_name"]
		maps[i]["user"] = maps[i]["admin_user"]
	}
	mystruct := &util.ResultB{0, num, &maps}
	c.Data["json"] = mystruct
	c.ServeJSON()
}

// 新增组织
type SaveOrganizeDataHandler struct {
	beego.Controller
}

func (c *SaveOrganizeDataHandler) Post() {
	name, pId, level, code, id, phone :=
		c.GetString("organize_name"), c.GetString("organize_pId"), c.GetString("organize_level"), c.GetString("code"), c.GetString("organize_id"), c.GetString("phone")
	if pId == "" {
		pId = "0"
	}
	pname := c.GetString("pname")
	if pname == "" {
		pname = c.GetString("organize_name")
	} else {
		pname = c.GetString("pname") + ">>" + c.GetString("organize_name")
	}
	log.Println(name, pId, level, code, id, pname)

	if code == "0" { //添加
		//创建一个节点号为1 节点的新节点 多台服务器不同范围1-1023
		o_id := uuid.Must(uuid.NewV4()).String()
		o_id = strings.Replace(o_id, "-", "", -1) // 生成SessionID

		level, err := strconv.Atoi(level)
		if err != nil {
			log.Println(err)
			c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织信息保存失败"}
			c.ServeJSON()
			return
		}
		level = level + 1

		o := controllers.O
		res, err := o.Raw("insert into ss_organize(id,name, pId, level,pname,phone) values(?, ?, ?, ?,?,?)", o_id, name, pId, level, pname, phone).Exec()
		if err != nil {
			log.Println("添加组织架构err:", res)
			c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "添加组织架构失败"}
			c.ServeJSON()
			return
		}
		c.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "添加组织架构成功"}
		c.ServeJSON()
	} else if code == "1" { //更新
		o := controllers.O
		res, err := o.Raw("update ss_organize set name = ?,phone=? where id = ?", name, phone, id).Exec()
		if err != nil {
			log.Println("更新数据异常", res)
			c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "修改组织架构失败"}
			c.ServeJSON()
			return
		}
		c.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "修改组织架构成功"}
		c.ServeJSON()
	}
}

type DelOrganizeDataHandler struct {
	beego.Controller
}

func (c *DelOrganizeDataHandler) Post() {
	Id := c.GetString("id")
	//查询是否包含字节点
	o := controllers.O
	var maps []orm.Params
	num, err := o.Raw("select * from ss_organize where pId=?;", Id).Values(&maps)
	if err != nil {
		log.Println("删除组织err:", err)
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "删除组织失败"}
		c.ServeJSON()
		return
	}

	if num > 0 {
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织包含子节点，请先删除子节点继续操作~"}
		c.ServeJSON()
		return
	}
	_, err = o.Raw("delete from ss_user where organize_id=?", Id).Exec()
	if err != nil {
		log.Println("组织成员删除失败", err)
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织成员删除失败"}
		c.ServeJSON()
		return
	}
	_, err = o.Raw("delete from ss_organize where id=?", Id).Exec()
	if err != nil {
		log.Println("组织删除失败", err)
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织删除失败"}
		c.ServeJSON()
		return
	}

	if len(Id) > 0 {
		//DelOrganize(Id)
	}
	c.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "数据删除成功"}
	c.ServeJSON()
}

type SaveUserOrganizeHandler struct {
	beego.Controller
}

func (c *SaveUserOrganizeHandler) Post() {
	oid, name, acount, pass, utype, jur, tele :=
		c.GetString("member_organizeId"),
		c.GetString("member_name"),
		c.GetString("member_user"),
		c.GetString("member_pass"),
		c.GetString("type"),
		c.GetString("jur"),
		c.GetString("member_tele")

	if len(oid) == 0 {
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "未选择组织架构"}
		c.ServeJSON()
		return
	}
	if len(name) == 0 {
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "员工名称不能为空"}
		c.ServeJSON()
		return
	}

	if len(acount) == 0 {
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "登录名不能为空"}
		c.ServeJSON()
		return
	}
	if len(pass) == 0 {
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "密码不能为空"}
		c.ServeJSON()
		return
	}
	if len(utype) == 0 {
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "账户类型不能为空"}
		c.ServeJSON()
		return
	}

	o := controllers.O

	var userMaps []*models.User
	qs := o.QueryTable(new(models.User))
	_, err := qs.Filter("admin_user", acount).All(&userMaps)
	if err != nil {
		log.Println(err)
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "验证用户信息异常"}
		c.ServeJSON()
		return
	}

	for _, u := range userMaps {
		if u.Acount == acount {
			c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "添加员工账号已存在"}
			c.ServeJSON()
			return
		}
	}

	uid := uuid.Must(uuid.NewV4()).String()
	uid = strings.Replace(uid, "-", "", -1) // 生成SessionID

	h := md5.New()
	h.Write([]byte(pass))
	md5Pass := hex.EncodeToString(h.Sum(nil)) //加密密码

	_, err = o.Raw("insert into ss_user(admin_id,admin_user,admin_password,admin_name,organize_id,update_time,type,jurisdiction,tele) values(?,?,?,?,?,?,?,?,?)", uid, acount, md5Pass, name, oid, time.Now().Format("2006-01-02 15:04:05"), utype, jur, tele).Exec()
	if err != nil {
		log.Println(err)
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "添加用户异常"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "添加成功"}
	c.ServeJSON()
}

// 组织修改成员
type EditUserOrganizeHandler struct {
	beego.Controller
}

func (c *EditUserOrganizeHandler) Post() {
	id, name, acount, pass, jur, utype, tele :=
		c.GetString("member_id"),
		c.GetString("member_name"),
		c.GetString("member_user"),
		c.GetString("member_pass"),
		c.GetString("jur"),
		c.GetString("type"),
		c.GetString("member_tele")

	if len(name) == 0 {
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "成员名称不能为空"}
		c.ServeJSON()
		return
	}
	if len(acount) == 0 {
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "登录名不能为空"}
		c.ServeJSON()
		return
	}
	if len(utype) == 0 {
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "账号类型不能为空"}
		c.ServeJSON()
		return
	}

	o := controllers.O
	var user models.User

	//查询编号，账号是否存在
	cont := orm.NewCondition()
	var userMaps []*models.User
	qs := o.QueryTable(new(models.User))
	_, err := qs.SetCond(cont.AndNot("admin_id", id)).All(&userMaps)
	if err != nil {
		log.Println(err)
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "验证用户信息异常"}
		c.ServeJSON()
		return
	}

	for _, u := range userMaps {
		if u.Acount == acount {
			c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "添加员工账号已存在"}
			c.ServeJSON()
			return
		}
	}

	err = o.QueryTable(new(models.User)).Filter("admin_id", id).One(&user)
	if err != nil {
		log.Println(err)
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "查询用户信息异常"}
		c.ServeJSON()
		return
	}

	user.Name = name
	user.Acount = acount
	user.Jurisdiction = jur
	user.Type = utype
	user.Tele = tele

	if len(pass) > 0 {
		h := md5.New()
		h.Write([]byte(pass))
		md5Pass := hex.EncodeToString(h.Sum(nil)) //加密密码
		user.Pass = md5Pass
		_, err = o.Update(&user, "admin_name", "admin_user", "jurisdiction", "type", "tele", "admin_password")
		if err != nil {
			log.Println(err)
			c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "更新用户信息异常"}
			c.ServeJSON()
			return
		}
	} else {
		_, err = o.Update(&user, "admin_name", "admin_user", "jurisdiction", "type", "tele")
		if err != nil {
			log.Println(err)
			c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "更新用户信息异常"}
			c.ServeJSON()
			return
		}
	}

	c.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "更新用户信息成功"}
	c.ServeJSON()
}

type DelUserOrganizeHandler struct {
	beego.Controller
}

func (c *DelUserOrganizeHandler) Post() {
	member_ids := c.GetString("member_ids")
	log.Println("member_ids", member_ids)
	paramsa := strings.Split(member_ids, ",")
	for index := range paramsa {
		member_id := paramsa[index]
		log.Println("member_id", member_id)
		o := controllers.O
		res, err := o.Raw("delete from ss_user where admin_id =?", member_id).Exec()
		if err != nil {
			log.Println("删除账号异常", res)
			c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "删除账号异常"}
			c.ServeJSON()
			return
		}
		if len(member_id) > 0 {
			//DelUser(member_id)
		}

	}
	c.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "删除账号成功"}
	c.ServeJSON()
}
