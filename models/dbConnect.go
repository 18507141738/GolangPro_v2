package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)
import _ "github.com/go-sql-driver/mysql"

type Place struct {
	ID         string    `orm:"column(place_id);pk"`
	Name       string    `orm:"column(place_name)"`
	Createtime string    `orm:"column(createtime)"`
	Organize   *Organize `orm:"column(organize_id);rel(fk)"`
	Camera     []*Camera `orm:"reverse(many)"`
}

type Camera struct {
	ID           string        `orm:"column(camera_id);pk"`
	Name         string        `orm:"column(camera_title)"`
	Cdconfig     string        `orm:"column(cdconfig)"`
	WRate        string        `orm:"column(wRate)"`
	HRate        string        `orm:"column(hRate)"`
	AlgorithmURL string        `orm:"column(algorithm_url)"`
	Width        string        `orm:"column(algorithm_width)"`
	Height       string        `orm:"column(algorithm_height)"`
	VideoCode    string        `orm:"column(videoCode)"`
	CameraNub    string        `orm:"column(cameraNub)"`
	Status       string        `orm:"column(status)"`
	Place        *Place        `orm:"column(place_id);rel(fk)"`
	Alarm        []*Alarm      `orm:"reverse(many)"`
	FuncCamera   []*FuncCamera `orm:"reverse(many)"`
	Host         *Host         `orm:"column(host_id);rel(fk)"`
}

type Alarm struct {
	ID         string `orm:"column(alarm_id);pk"`
	AType      string `orm:"column(alarm_type)"`
	APlaceType string `orm:"column(alarm_place_type)"`
	APlace     string `orm:"column(alarm_place)"`
	ADetial    string `orm:"column(alarm_detial)"`
	AFile      string `orm:"column(alarm_file)"`
	AVideo     string `orm:"column(alarm_video)"`
	AStream    string `orm:"column(alarm_stream)"`
	AHead      string `orm:"column(alarm_head)"`
	Atime      string `orm:"column(alarm_time)"`
	Astatus    string `orm:"column(alarm_status)"`
	ALevel     string `orm:"column(alarm_level)"`
	Syn        string `orm:"column(syn)"`
	PageRead   string `orm:"column(page_read_type)"`
	Handler    string `orm:"column(alarmhandler)"`
	Hostid     string `orm:"column(host_id)"`
	//Host       *Host   `orm:"column(host_id);rel(fk)"`
	Camera   *Camera     `orm:"column(camera_id);rel(fk)"`
	Notation []*Notation `orm:"reverse(many)"`
}

type User struct {
	ID     string `orm:"column(admin_id);pk"`
	Acount string `orm:"column(admin_user)"`
	//Number       string    `orm:"column(number)"`
	Pass         string      `orm:"column(admin_password)"`
	Name         string      `orm:"column(admin_name)"`
	UpdateTime   string      `orm:"column(update_time)"`
	Type         string      `orm:"column(type)"`
	Jurisdiction string      `orm:"column(jurisdiction)"`
	Organize     *Organize   `orm:"column(organize_id);rel(fk)"`
	Notation     []*Notation `orm:"reverse(many)"`
	Tele         string      `orm:"column(tele)"`
	Allowd       string      `orm:"column(allowd)"`
}

type Organize struct {
	ID    string  `orm:"column(id);pk"`
	PID   string  `orm:"column(pId)"`
	Name  string  `orm:"column(name)"`
	Level string  `orm:"column(level)"`
	PName string  `orm:"column(pname)"`
	Phone string  `orm:"column(phone)"`
	User  []*User `orm:"reverse(many)"`
	//Host  []*Host  `orm:"reverse(many)"`
	Place []*Place `orm:"reverse(many)"`
}

type Host struct {
	ID         string `orm:"column(host_id);pk"`
	Name       string `orm:"column(host_name)"`
	UpdateTime string `orm:"column(updateTime);"`
	Status     string `orm:"column(host_status);"`
	IP         string `orm:"column(host_ip)"`
	//Organize *Organize `orm:"column(org_id);rel(fk)"`
	Camera []*Camera `orm:"reverse(many)"`
	//Alarm    []*Alarm  `orm:"reverse(many)"`
}

type FuncCamera struct {
	ID                   string  `orm:"column(id);pk"`
	Type                 string  `orm:"column(function_type)"`
	Boundary             string  `orm:"column(boundary)"`
	DistanceMode         string  `orm:"column(DistanceMode)"`
	PerimeterThreshold   string  `orm:"column(perimeterThreshold)"`
	Num_person_threshold string  `orm:"column(num_person_threshold)"`
	Real_sleep_minutes   string  `orm:"column(real_sleep_minutes)"`
	Threshold            string  `orm:"column(threshold)"`
	Frequency            string  `orm:"column(frequency)"`
	WallOn               string  `orm:"column(wallOn)"`
	Wallrate             string  `orm:"column(wallrate)"`
	Rateup               string  `orm:"column(rateup)"`
	Ratedown             string  `orm:"column(ratedown)"`
	Match_boundary0      string  `orm:"column(match_boundary0)"`
	Match_boundary1      string  `orm:"column(match_boundary1)"`
	Match_boundary2      string  `orm:"column(match_boundary2)"`
	Hatcolor             string  `orm:"column(hatcolor)"`
	Alarm_mode           string  `orm:"column(alarm_mode)"`
	Detect_number        string  `orm:"column(detect_number)"`
	Warning_number       string  `orm:"column(warning_number)"`
	Departure_minute     string  `orm:"column(departure_minute)"`
	Timepoint            string  `orm:"column(timepoint)"`
	Alarm_time_second    string  `orm:"column(alarm_time_second)"`
	Video_switch         string  `orm:"column(video_switch)"`
	Bg_thr               string  `orm:"column(bg_thr)"`
	LeakThreshold        string  `orm:"column(leakThreshold)"`
	Swicth               string  `orm:"column(switch)"`
	Distance             string  `orm:"column(distance)"`
	Boundary2            string  `orm:"column(boundary2)"`
	Negative             string  `orm:"column(negative)"`
	Camera               *Camera `orm:"column(cameraId);rel(fk)"`
}

type Notation struct {
	ID         string `orm:"column(id);pk"`
	Createtime string `orm:"column(createtime);"`
	Notation   string `orm:"column(notation);"`
	//UserID       string    `orm:"column(userid);"`
	//AlarmID       string    `orm:"column(alarm_id);"`
	User  *User  `orm:"column(userid);rel(fk);"`   //设置一对多关系
	Alarm *Alarm `orm:"column(alarm_id);rel(fk);"` //设置一对多关系
}

type AlarmDetail struct {
	ID     string `orm:"column(id);pk"`
	Type   string `orm:"column(type)"`
	Detail string `orm:"column(detail)"`
}

func init() {
	orm.RegisterModelWithPrefix("ss_", new(Camera))
	orm.RegisterModelWithPrefix("ss_", new(Place))
	orm.RegisterModelWithPrefix("ss_", new(Alarm))
	orm.RegisterModelWithPrefix("ss_", new(Organize))
	orm.RegisterModelWithPrefix("ss_", new(User))
	orm.RegisterModelWithPrefix("ss_func_", new(Host))
	orm.RegisterModelWithPrefix("ss_", new(FuncCamera))
	orm.RegisterModelWithPrefix("ss_", new(AlarmDetail))
	orm.RegisterModelWithPrefix("ss_", new(Notation))
}

func DBBase() {
	username := beego.AppConfig.String("mysqluser")
	password := beego.AppConfig.String("mysqlpass")
	url := beego.AppConfig.String("mysqlurls")
	dbName := beego.AppConfig.String("mysqldb")
	maxConn := 300 // 设置最大打开的连接数,0表示不限制
	maxIdle := 30  // 设置闲置的连接数
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", username+":"+password+"@tcp("+url+")/"+dbName+"?charset=utf8", maxIdle, maxConn)

}
