package controllers

import (
	"Artifice_V2.0.0/models"
	"Artifice_V2.0.0/util"
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"log"
	"strconv"
	"strings"
	"time"
)

type PlatOrganizeBelongUser struct {
	beego.Controller
}

//获取用户下的组织架构
func (this *PlatOrganizeBelongUser) Post() {
	userID := this.Ctx.GetCookie("user_id")

	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", userID).RelatedSel().One(&user)
	if err != nil {
		LogsError("获取用户信息异常:", err)
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "查询失败"}
		this.ServeJSON()
		return
	}

	orgIds, err := SelOrgIDBelongUserOrg(user.Organize)
	if err != nil {
		LogsError("查询用户下属组织异常:", err)

	}
	var ids = []string{user.Organize.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}

	var orgs []models.Organize
	_, err = o.QueryTable(new(models.Organize)).Filter("id__in", ids).OrderBy("level").All(&orgs)
	if err != nil {
		LogsError("查询用户下组织详细信息异常:", err)
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "查询失败"}
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"Code": 1, "Reason": orgs}
	this.ServeJSON()

	return
}

//获取用户下属组织id
func SelOrgIDBelongUserOrg(organize *models.Organize) ([]orm.Params, error) {
	o := O
	var maps []orm.Params

	sql := "select id from " +
		"(select tt1.id," +
		" if(find_in_set(pId, @pids) > 0, @pids := concat(@pids, ',', id), 0) as ischild " +
		"from " +
		"(select id,pId,`level` from ss_organize tt order by `level` ASC) tt1," +
		"(select @pids := ?) tt2" +
		") tt3 " +
		"where ischild != '0'"
	//sql := "select id from " +
	//	"(select tt1.id,tt1.`level`," +
	//	" if(find_in_set(pId, @pids) > 0, @pids := concat(@pids, ',', id), 0) as ischild " +
	//	"from " +
	//	"(select id,pId,`level` from ss_organize tt order by `level` ASC) tt1," +
	//	"(select @pids := ?) tt2" +
	//	") tt3 " +
	//	"where ischild != '0' ORDER BY `level` ASC"
	_, err := o.Raw(sql, organize.ID).Values(&maps)
	return maps, err
}

//parentorg是否包含childOrg
func InOrganize(childOrg *models.Organize, parentOrg *models.Organize) bool {
	maps, err := SelOrgIDBelongUserOrg(parentOrg)
	if err != nil {
		return false
	}

	for _, v := range maps {
		if v["id"] == childOrg.ID {
			return true
		}
	}

	return false
}

//展示页面获取组织架构
type PlatSelOrgainze struct {
	beego.Controller
}

func (this *PlatSelOrgainze) Post() {
	admin_id := this.Ctx.GetCookie("user_id")

	o := O
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户信息查询异常"}
		this.ServeJSON()
		return
	}

	org_id := user.Organize.ID // 用户组织架构
	var maps []orm.Params
	_, err = o.Raw("select id from (select t1.id, if(find_in_set(pId, @pids) > 0, @pids := concat(@pids, ',', id), 0) as ischild from (select id,pId from ss_organize t order by `level` ASC) t1, (select @pids := ?) t2) t3 where ischild != '0'", org_id).Values(&maps)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "组织架构查询异常"}
		this.ServeJSON()
		return
	}

	var ids = []string{}
	ids = append(ids, org_id)
	for _, v := range maps {
		ids = append(ids, v["id"].(string))
	}

	var orgs []orm.Params
	num, err := o.QueryTable("ss_organize").Filter("id__in", ids).Values(&orgs)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "组织架构查询异常"}
		this.ServeJSON()
		return
	}
	for i := 0; i < len(orgs); i++ {
		orgs[i]["icon"] = "/static/css/zTreeStyle/img/diy/homeIcon.png"
		orgs[i]["isParent"] = "true"
	}

	mystruct := &util.ResultB{1, num, &orgs}
	this.Data["json"] = mystruct
	this.ServeJSON()
	return
}

// 新增组织
type SaveOrganizeHandler struct {
	beego.Controller
}

func (this *SaveOrganizeHandler) Post() {
	name, phone, pid, pname, plevel :=
		this.GetString("organize_name"),
		this.GetString("phone"),
		this.GetString("organize_pId"),
		this.GetString("pname"),
		this.GetString("plevel")
	if len(name) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织名称不能为空"}
		this.ServeJSON()
		return
	}
	if len(phone) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "联系电话不能为空"}
		this.ServeJSON()
		return
	}
	o := O
	userDB := new(models.User)
	var organize models.Organize
	var user models.User
	adminid := this.Ctx.GetCookie("admin_id")
	err := o.QueryTable(userDB).Filter("admin_id", adminid).One(&user)
	if err != nil {
		log.Println("err1:", err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "查询账户权限异常"}
		this.ServeJSON()
		return
	}
	if len(user.Organize.ID) > 0 && len(pid) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "请选择组织架构"}
		this.ServeJSON()
		return
	}
	level, err := strconv.Atoi(plevel)
	if err != nil {
		log.Println(err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织信息保存失败"}
		this.ServeJSON()
		return
	}
	level = level + 1
	if level > 4 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "超出组织级别层次"}
		this.ServeJSON()
		return
	}
	uid := uuid.Must(uuid.NewV4()).String()
	uid = strings.Replace(uid, "-", "", -1) // 生成SessionID
	organize.ID = uid
	organize.PID = pid
	organize.Name = name
	organize.Level = strconv.Itoa(level)
	organize.Phone = phone
	organize.PName = pname

	id, err := o.Insert(&organize)
	log.Print(id)
	if err != nil {
		log.Println(err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织信息保存失败"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "操作成功"}
	this.ServeJSON()
}

// 新增组织
type PlatSaveOrganizeHandler struct {
	beego.Controller
}

func (this *PlatSaveOrganizeHandler) Post() {
	name, phone, pid, pname, plevel :=
		this.GetString("organize_name"),
		this.GetString("phone"),
		this.GetString("organize_pId"),
		this.GetString("pname"),
		this.GetString("plevel")
	if len(name) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织名称不能为空"}
		this.ServeJSON()
		return
	}
	if len(phone) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "联系电话不能为空"}
		this.ServeJSON()
		return
	}
	o := O
	userDB := new(models.User)
	var organize models.Organize
	var user models.User
	adminid := this.Ctx.GetCookie("user_id")
	err := o.QueryTable(userDB).Filter("admin_id", adminid).One(&user)
	if err != nil {
		log.Println("err1:", err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "查询账户权限异常"}
		this.ServeJSON()
		return
	}
	if len(user.Organize.ID) > 0 && len(pid) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "请选择组织架构"}
		this.ServeJSON()
		return
	}
	level, err := strconv.Atoi(plevel)
	if err != nil {
		log.Println(err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织信息保存失败"}
		this.ServeJSON()
		return
	}
	level = level + 1
	if level > 4 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "超出组织级别层次"}
		this.ServeJSON()
		return
	}
	uid := uuid.Must(uuid.NewV4()).String()
	uid = strings.Replace(uid, "-", "", -1) // 生成SessionID
	organize.ID = uid
	organize.PID = pid
	organize.Name = name
	organize.Level = strconv.Itoa(level)
	organize.Phone = phone
	organize.PName = pname

	id, err := o.Insert(&organize)
	log.Print(id)
	if err != nil {
		log.Println(err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织信息保存失败"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "操作成功"}
	this.ServeJSON()
}

// 修改组织
type EditOrganizeHandler struct {
	beego.Controller
}

func (this *EditOrganizeHandler) Post() {
	name, phone, id :=
		this.GetString("organize_name"),
		this.GetString("phone"),
		this.GetString("organize_id")

	if len(name) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织名称不能为空"}
		this.ServeJSON()
		return
	}

	if len(phone) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "联系电话不能为空"}
		this.ServeJSON()
		return
	}

	o := O
	organizeDB := new(models.Organize)
	var organize models.Organize

	err := o.QueryTable(organizeDB).Filter("id", id).One(&organize)
	if err != nil {
		log.Println(err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织信息保存失败"}
		this.ServeJSON()
		return
	}

	organize.Name = name
	organize.Phone = phone

	_, err = o.Update(&organize, "name", "phone")

	if err != nil {
		log.Println(err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织信息保存失败"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "操作成功"}
	this.ServeJSON()
}

// 删除组织
type DelOrganizeHandler struct {
	beego.Controller
}

func (this *DelOrganizeHandler) Post() {
	Id := this.GetString("id")
	//查询是否包含字节点
	o := O
	var maps []orm.Params
	num, err := o.Raw("select * from ss_organize where pId=?;", Id).Values(&maps)
	if err != nil {
		log.Println("删除组织err:", err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "删除组织失败"}
		this.ServeJSON()
		return
	}

	if num > 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织包含子节点，请先删除子节点继续操作~"}
		this.ServeJSON()
		return
	}
	_, err = o.Raw("delete from ss_user where organize_id=?", Id).Exec()
	if err != nil {
		log.Println("组织成员删除失败", err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织成员删除失败"}
		this.ServeJSON()
		return
	}
	_, err = o.Raw("delete from ss_organize where id=?", Id).Exec()
	if err != nil {
		log.Println("组织删除失败", err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "组织删除失败"}
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "数据删除成功"}
	this.ServeJSON()
}

// 查询组织用户
type SelMemberByOrgHandler struct {
	beego.Controller
}

func (this *SelMemberByOrgHandler) Post() {
	id := this.GetString("organize_id")

	o := O
	userDB := new(models.User)
	var users []*models.User

	userid := this.Ctx.GetCookie("admin_id")

	num, err := o.QueryTable(userDB).Exclude("admin_id", userid).Filter("organize_id", id).RelatedSel().All(&users)

	if err != nil {
		log.Println(err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "查询组织成员失败"}
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"Ret": 0, "num": num, "Reason": &users}
	this.ServeJSON()
}

// 查询组织用户-展示页
type SelMemberByOrgHandler2 struct {
	beego.Controller
}

func (this *SelMemberByOrgHandler2) Post() {
	id := this.GetString("organize_id")

	o := O
	userDB := new(models.User)
	var users []*models.User

	userid := this.Ctx.GetCookie("user_id")

	num, err := o.QueryTable(userDB).Exclude("admin_id", userid).Exclude("type", "0").Filter("organize_id", id).RelatedSel().All(&users)

	if err != nil {
		log.Println(err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "查询组织成员失败"}
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"Ret": 0, "num": num, "Reason": &users}
	this.ServeJSON()
}

// 组织添加成员
type OrgAddUserHandler struct {
	beego.Controller
}

func (c *OrgAddUserHandler) Post() {
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

	o := O

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
type OrgEditUserHandler struct {
	beego.Controller
}

func (c *OrgEditUserHandler) Post() {
	id, name, acount, tele, pass, jur, utype :=
		c.GetString("member_id"),
		c.GetString("member_name"),
		c.GetString("member_user"),
		c.GetString("member_tele"),
		c.GetString("member_pass"),
		c.GetString("jur"),
		c.GetString("type")

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

	o := O
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

// 组织删除成员
type OrgDelUserHandler struct {
	beego.Controller
}

func (c *OrgDelUserHandler) Post() {
	id := c.GetString("member_id")

	ids := strings.Split(id, ",")

	o := O
	var users []models.User

	_, err := o.QueryTable(new(models.User)).Filter("admin_id__in", ids).All(&users)
	if err != nil && err == orm.ErrNoRows {
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "未查到该用户用户信息"}
		c.ServeJSON()
		return
	} else if err != nil {
		log.Println(err)
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "查询该用户用户信息异常"}
		c.ServeJSON()
		return
	}
	_, err = o.QueryTable(new(models.User)).Filter("admin_id__in", ids).Delete()
	if err != nil {
		log.Println(err)
		c.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "删除该用户用户信息异常"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "删除成功"}
	c.ServeJSON()
}
