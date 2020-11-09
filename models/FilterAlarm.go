package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type FilterAlarm struct {
	ID         string  `orm:"column(alarm_id);pk"`
	AType      string  `orm:"column(alarm_type)"`
	APlaceType string  `orm:"column(alarm_place_type)"`
	APlace     string  `orm:"column(alarm_place)"`
	ADetial    string  `orm:"column(alarm_detial)"`
	AFile      string  `orm:"column(alarm_file)"`
	AVideo     string  `orm:"column(alarm_video)"`
	AStream    string  `orm:"column(alarm_stream)"`
	AHead      string  `orm:"column(alarm_head)"`
	Atime      string  `orm:"column(alarm_time)"`
	Status     string  `orm:"column(status)"`
	Hostid     string  `orm:"column(host_id)"`
	ALevel     string  `orm:"column(alarm_level)"`
	PageRead   string  `orm:"column(page_read_type)"`
	Camera     *Camera `orm:"column(camera_id);rel(fk)"`
}

func (t *FilterAlarm) TableName() string {
	return "ss_alarm_filter"
}

func init() {
	orm.RegisterModel(new(FilterAlarm))
}

func GetFilterAlarmCountByType(ormer orm.Ormer, atype string, orgids []string) (count int64, err error) {
	var endtime = time.Now().Format("2006-01-02") + " 23:59:59"
	var starttime = time.Now().AddDate(0, 0, -30).Format("2006-01-02") + " 00:00:00"
	count, err = ormer.QueryTable(new(FilterAlarm)).Filter("Status", "0").Filter("AType", atype).Filter("alarm_time__gte", starttime).Filter("alarm_time__lte", endtime).Filter("Camera__Place__Organize__ID__in", orgids).RelatedSel().Count()
	return count, err
}

func FilterAlarmByID(ormer orm.Ormer, id string) (f FilterAlarm, err error) {
	err = ormer.QueryTable(new(FilterAlarm)).Filter("alarm_id", id).One(&f)
	return f, err
}

func FilterAlarmPage(ormer orm.Ormer, cameraid string, funcType string, orgids []string, starttime string, endtime string, limit int, page int, status string) (maps []FilterAlarm, zong int64, err error) {
	qs := ormer.QueryTable(new(FilterAlarm)).Filter("alarm_type", funcType).Filter("Camera__Place__Organize__ID__in", orgids).Filter("alarm_time__gte", starttime).Filter("alarm_time__lte", endtime).OrderBy("-alarm_time").RelatedSel(3)
	if len(cameraid) > 0 && cameraid != "0" {
		qs = qs.Filter("Camera__ID", cameraid)
	}
	if status == "" {
		qs = qs.Filter("Status", "0")
	} else if status != "" && status != "-1" {
		qs = qs.Filter("Status", status)
	}
	_, err = qs.Limit(limit, page).All(&maps)
	if err != nil {
		return nil, 0, err
	}

	zong, err = qs.Count()
	if err != nil {
		return nil, 0, nil
	}

	return maps, zong, nil
}

func FilterNewAlarm(ormer orm.Ormer, orgids []string) (maps []FilterAlarm, num int64, err error) {
	qs := ormer.QueryTable(new(FilterAlarm)).Filter("page_read_type", "0").Filter("Camera__Place__Organize__ID__in", orgids).RelatedSel("Camera__Place__Organize")
	num, err = qs.All(&maps)
	if err != nil {
		return nil, 0, err
	}
	qs.Update(orm.Params{
		"page_read_type": "1",
	})
	if err != nil {
		return nil, 0, err
	}
	return maps, num, nil
}

func FilterUpdateAlarmStatus(ormer orm.Ormer, alarm FilterAlarm) (id int64, err error) {
	return ormer.Update(&alarm, "status")
}
