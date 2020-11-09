package controllers

import (
	"Artifice_V2.0.0/conf/MacCameraNum"
	"Artifice_V2.0.0/models"
	"Artifice_V2.0.0/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	"log"
	"strconv"
	"sync"
	"time"
)

var lock sync.Mutex

// 分页查询主机列表
type PlatformSelHostPageHandler struct {
	beego.Controller
}

func (this *PlatformSelHostPageHandler) Post() {
	admin_id := this.Ctx.GetCookie("user_id")
	o := O
	var user models.User
	var organize models.Organize
	var org *models.Organize
	// 获取登录用户详情
	err := o.QueryTable(new(models.User)).Filter("admin_id", admin_id).RelatedSel().One(&user)
	org = user.Organize
	var mapss []orm.Params
	num, err := o.Raw("select id from (select t1.id, if(find_in_set(pId, @pids) > 0, @pids := concat(@pids, ',', id), 0) as ischild from (select id,pId from ss_organize t order by `level` ASC) t1, (select @pids := ?) t2) t3 where ischild != '0'", org.ID).Values(&mapss)
	log.Println(num)
	org_id, hostStatus := this.GetString("org_id"), this.GetString("hostStatus")
	bools := false
	for i := 0; i < len(mapss); i++ {
		if mapss[i]["id"].(string) == org_id || org_id == org.ID {
			bools = true // 参数组织ID包含在adminID的组织内
		}
	}
	if bools {
		err := o.QueryTable(new(models.Organize)).Filter("id", org_id).RelatedSel().One(&organize)
		log.Println(err)
		org = &organize
	}

	if err != nil {
		// 报错或者无站
		mystruct := &util.ResultB{0, 0, nil}
		this.Data["json"] = mystruct
		this.ServeJSON()
	} else {
		limits, pages := this.GetString("limit"), this.GetString("page")
		start_page, err := strconv.Atoi(pages)
		limitss, err := strconv.Atoi(limits)
		start_pages := (start_page - 1) * limitss
		sqlStartStr := "SELECT b.host_name as name,b.host_id as hostID, b.org_id, b.host_ip as hostIP" +
			", b.alarmNumber as alarmNumber, b.syn_center as sync_center, b.host_status as state FROM ss_func_host as b"
		if hostStatus != "-1" {
			sqlStartStr += " where 1=1 and b.host_status = " + hostStatus
		}
		var maps []orm.Params
		num, err := o.Raw(sqlStartStr+" order by host_id desc limit ?,?", int(start_pages), int(limitss)).Values(&maps)
		if err != nil {
			log.Println("分页查询主机列表：", err)
			return
		}

		for j := 0; j < len(maps); j++ {
			var sumON []orm.Params
			num, err := o.Raw("SELECT * FROM `ss_camera` where host_id=? and switch_c='on'", maps[j]["hostID"]).Values(&sumON)
			maps[j]["switchOn"] = num
			var sumOFF []orm.Params
			num1, err := o.Raw("SELECT * FROM `ss_camera` where host_id=? and switch_c='off'", maps[j]["hostID"]).Values(&sumOFF)
			maps[j]["switchOFF"] = num1

			if err != nil {
				log.Println("分页查询主机列表：", err)
				return
			}
		}

		var mapsr []orm.Params
		num, err = o.Raw(sqlStartStr).Values(&mapsr)
		mystruct := &util.ResultB{0, num, &maps}
		this.Data["json"] = mystruct
		this.ServeJSON()
	}
}

// 新增、修改分析主机
type HostManageHandler struct {
	beego.Controller
}

func (this *HostManageHandler) Post() {
	hostID, name, hostIP, sync_center, alarmNumber := this.GetString("hostID"),
		this.GetString("name"), this.GetString("hostIP"), this.GetString("sync_center"),
		this.GetString("alarmNumber")

	if len(alarmNumber) > 12 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "主机编号长度不可超过12位"}
		this.ServeJSON()
		return
	}

	o := O
	if hostID != "" {
		var maps []orm.Params
		num, err := o.Raw("select * from ss_func_host where host_ip=?", hostIP).Values(&maps)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败"}
			this.ServeJSON()
			return
		}
		if num > 0 && maps[0]["host_id"] != hostID {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "已存在相同主机IP"}
			this.ServeJSON()
			return
		}

		res, err := o.Raw("update ss_func_host set host_name=?,host_ip=?,syn_center=?,alarmNumber=? where host_id=?",
			name, hostIP, sync_center, alarmNumber, hostID).Exec()
		if err != nil {
			log.Println(res, err)
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败"}
		} else {
			this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "操作成功"}
		}
	} else {
		var maps []orm.Params
		num, err := o.Raw("select * from ss_func_host where host_ip=?", hostIP).Values(&maps)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "获取主机信息失败"}
		} else {
			if num > 0 {
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "已存在相同主机IP"}
			} else {
				res, err := o.Raw("insert into ss_func_host(host_name,createtime,host_ip,syn_center,host_status,func_status,alarmNumber) values(?,?,?,?,1,0,?)",
					name, time.Now().Format("2006-01-02 15:04:05"), hostIP, sync_center, alarmNumber).Exec()
				log.Println(res, err)
				if err != nil {
					this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败"}
				} else {
					this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "操作成功"}
				}
			}
		}
	}
	this.ServeJSON()
}

// 批量删除分析主机
type DelBatchHostHandler struct {
	beego.Controller
}

func (this *DelBatchHostHandler) Post() {
	hostIDS := this.GetString("hostIDS")
	var maps []orm.Params
	o := O
	num, err := o.Raw("select * from ss_camera where host_id in (" + hostIDS + ")").Values(&maps)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败"}
	} else {
		log.Println(num)
		if num > 0 {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "删除失败，分析主机中包含摄像头"}
		} else {
			res, err := o.Raw("delete from ss_func_host where host_id in (" + hostIDS + ")").Exec()
			if err != nil {
				log.Println(res, err)
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败"}
			} else {
				this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "操作成功"}
			}
		}
	}
	this.ServeJSON()
}

// 判断主机是否超时
type PlatformCheckHostOutHandler struct {
	beego.Controller
}

func (this *PlatformCheckHostOutHandler) Post() {
	id := this.GetString("id")
	var maps []orm.Params
	o := O
	num, err := o.Raw("select * from  ss_func_host where host_id = ?", id).Values(&maps)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "主机已超时"}
		return
	}
	if num > 0 {
		if maps[0]["overtime"] == "1" {
			this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "主机已超时"}
			return
		} else {
			this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "主机未超时"}
		}
	}
	this.ServeJSON()
}

func TestHostStatus2() (b bool) {

	o := O
	hostDB := new(models.Host)
	var hosts []orm.Params

	_, err := o.QueryTable(hostDB).Values(&hosts, "host_id", "updateTime", "host_status")
	if err != nil {
		log.Println("获取主机列表失败", err)
		return false
	}
	//logs.Info("host",hosts)
	//0在线 1离线
	for _, value := range hosts {
		lock.Lock()
		//log.Println(value)
		var status = "1"
		hostID := value["ID"].(string)
		if value["UpdateTime"] != nil {
			//logs.Info("hostID:",hostID,",updatetime:",value["UpdateTime"].(string),",status:",value["Status"].(string))
			updateTime := value["UpdateTime"].(string)

			nowTime := time.Now()
			loc, _ := time.LoadLocation("Local") //获取本地时区
			uTime, _ := time.ParseInLocation("2006-01-02 15:04:05", updateTime, loc)

			subM := nowTime.Sub(uTime)

			//计算时差
			if subM.Seconds() < 60 { //设置长时间离线时间、用来缓冲意外原因造成的重启等问题，可以设置小些时间
				status = "0"

			} else if value["Status"] == "1" { //超时了，也是离线的跳过数据修改，没有意义
				lock.Unlock()
				continue
			}
		} else if value["Status"] == "1" {
			lock.Unlock()
			continue
		}

		_, err = o.QueryTable(hostDB).Filter("host_id", hostID).Update(orm.Params{
			"host_status": status,
		})
		if err != nil {
			log.Println("更新主机状态失败:", err)
		}
		//logs.Info("更新主机",hostID,"状态:", status)
		lock.Unlock()
	}

	return true
}

func setHostFuncCameras() {

	var hostList []models.Host
	o := O
	_, err := o.QueryTable(new(models.Host)).All(&hostList) //查询所有主机
	if err != nil {
		LogsError("--分布式接口 统计所有主机异常--", err)
		return
	}
	var ON_Hosts []models.Host  //在线主机
	var OFF_Hosts []models.Host //离线主机
	for _, h := range hostList {
		if h.Status == "0" { //在线
			ON_Hosts = append(ON_Hosts, h)
		} else { //离线
			OFF_Hosts = append(OFF_Hosts, h)
		}
	}
	dict := orm.Params{}            //存放在线主机，启用摄像机的的字典
	for _, host := range ON_Hosts { //遍历在线主机，查询摄像机配置表
		var resultMaps []models.FuncCamera //.Filter("switch", "on")
		_, err = o.QueryTable(new(models.FuncCamera)).Filter("Camera__Host__ID", host.ID).RelatedSel(4).All(&resultMaps)
		if err != nil {
			LogsError("--分布式接口 查询在线主机摄像头异常--", err)
			return
		}
		dictf := orm.Params{}
		dictf["AMaps"] = resultMaps //AMaps存放主机下的摄像头配置
		dict[host.ID] = dictf       //以主机ID作为key存储
	}

	var OFFHostIDS []string //离线主机ID数组
	for _, v := range OFF_Hosts {
		OFFHostIDS = append(OFFHostIDS, v.ID)
	}

	if len(OFFHostIDS) > 0 { //查询离线主机下所有已启用的摄像机配置表，并以BMaps存储
		var OFFHostCamera []models.FuncCamera
		_, err = o.QueryTable(new(models.FuncCamera)).Filter("Camera__Host__ID__in", OFFHostIDS).Filter("switch", "on").RelatedSel(4).All(&OFFHostCamera)
		if err != nil {
			LogsError("--分布式接口 查询离线主机摄像头异常--", err)
			return
		}

		var maxcamera = MacCameraNum.GetMaxCameraNum() //最大分配数量
		var lastIndex = 0
		for v := range dict {
			dictf := dict[v].(orm.Params)
			arr := dictf["AMaps"]
			//arr := dict[v]
			var Dvalue = maxcamera - len(arr.([]models.FuncCamera)) //还可以分配的数量
			if Dvalue > 0 {
				//如果分配数量超过本地数据
				if Dvalue+lastIndex > len(OFFHostCamera) {
					Dvalue = len(OFFHostCamera) - lastIndex
				}
				if Dvalue < 1 { //可分配为0 跳过分配
					continue
				}
				arrfc := OFFHostCamera[lastIndex : Dvalue+lastIndex]
				dictf["BMaps"] = arrfc
				dict[v] = dictf
				//log.Println("切数组:%s",arrfc," - startindex:",lastIndex," - dvalue:",Dvalue)
				lastIndex = Dvalue + lastIndex

			} else {
				continue
			}
		}
	}

	WriteFile(dict)
}

func task() error {

	TestHostStatus2()
	setHostFuncCameras()

	return nil
}

func InitHostTask() {
	//没个n分钟执行一次 0 */n * * * *
	//没个n秒执行一次 0/n * * * * *
	tk := toolbox.NewTask("testHostStatus", "0/10 * * * * *", task)
	//err := tk.Run()
	//if err != nil {
	//	fmt.Println("启动主机验证异常：",err)
	//}
	//加入全局的计划任务列表
	toolbox.AddTask("testHostStatus", tk)
	//开始执行全局的任务
	toolbox.StartTask()
	//defer toolbox.StopTask()
}
