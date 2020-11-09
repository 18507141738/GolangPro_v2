package models

import (
	"github.com/astaxie/beego/orm"
)

type FilterUser struct {
	ID           string    `orm:"column(id);pk"`
	Acount       string    `orm:"column(username)"`
	Pass         string    `orm:"column(password)"`
	Jurisdiction string    `orm:"column(jurisdiction)"`
	Organize     *Organize `orm:"column(organize_id);rel(fk)"`
}

func (m *FilterUser) TableName() string {
	return "ss_filter_user"
}

func init() {
	orm.RegisterModel(new(FilterUser))
}

/**
获取指定组织下的过滤账户
orgid: 组织id
excid：过滤当前账户id
*/
func GetFilterUserByOrgid(o orm.Ormer, orgid string, excid string) (f []*FilterUser, err error) {
	_, err = o.QueryTable(new(FilterUser)).Exclude("ID", excid).Filter("Organize__ID", orgid).RelatedSel("Organize").All(&f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// 添加用户
func AddFilterUser(o orm.Ormer, u FilterUser) (err error) {
	_, err = o.Insert(&u)
	return err
}

// 更新用户
func UpdateFilterUser(o orm.Ormer, u FilterUser) (err error) {
	if u.Pass == "" || len(u.Pass) == 0 {
		_, err = o.Update(&u, "username")
	} else {
		_, err = o.Update(&u, "username", "password")
	}
	return err
}

// 删除用户
func DelFilterUser(o orm.Ormer, u FilterUser) (err error) {
	_, err = o.Delete(&u)
	return err
}

func SelFilterUserByID(o orm.Ormer, id string) (u FilterUser, err error) {
	err = o.QueryTable(new(FilterUser)).Filter("ID", id).One(&u)
	return u, err
}

// 账户名是否存在
func FilterExist(o orm.Ormer, acount string) (b bool) {
	c, err := o.QueryTable(new(FilterUser)).Filter("username", acount).Count()
	if err != nil {
		return true
	}
	if c > 0 {
		return true
	}
	return false
}
