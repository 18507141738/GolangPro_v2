package controllers

import (
	"Artifice_V2.0.0/conf/MacCameraNum"
	"Artifice_V2.0.0/models"
	"Artifice_V2.0.0/tool"
	_ "Artifice_V2.0.0/tool"
	"Artifice_V2.0.0/util"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"log"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ClientDistributeCamera_Handler struct {
	beego.Controller
}

func (this *ClientDistributeCamera_Handler) Get() {
	host_ip := this.GetString("ip")
	o := O
	//查询该主机
	var host models.Host

	err := o.QueryTable(new(models.Host)).Filter("host_ip__contains", host_ip).One(&host)
	if err != nil {
		//log.Println("获取主机信息异常:", err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "获取主机信息异常", "Data": nil}
		this.ServeJSON()
		return
	}
	host.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	host.Status = "0"
	_, err = o.Update(&host, "updateTime", "host_status")
	if err != nil {
		//log.Println("更新数据库异常:", err)
		//logs.Info("更新主机请求时间异常：",err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "更新主机时间异常", "Data": nil}
		this.ServeJSON()
		return
	}

	var hostList []models.Host

	_, err = o.QueryTable(new(models.Host)).All(&hostList)
	if err != nil {
		//log.Println("获取主机列表失败", err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "获取主机列表失败", "Data": nil}
		this.ServeJSON()
		return
	}
	var ON_Hosts []models.Host
	var OFF_Hosts []models.Host
	for _, h := range hostList {
		if h.Status == "0" { //在线
			ON_Hosts = append(ON_Hosts, h)
		} else { //离线
			OFF_Hosts = append(OFF_Hosts, h)
		}
	}

	var resultMaps []models.FuncCamera
	_, err = o.QueryTable(new(models.FuncCamera)).Filter("Camera__Host__ID", host.ID).Filter("switch", "on").RelatedSel(4).All(&resultMaps)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "获取摄像机配置异常", "Data": nil}
		this.ServeJSON()
		return
	}

	var OFFHostIDS []string
	for _, v := range OFF_Hosts {
		OFFHostIDS = append(OFFHostIDS, v.ID)
	}

	var OFFHostCamera []models.FuncCamera
	_, err = o.QueryTable(new(models.FuncCamera)).Filter("Camera__Host__ID__in", OFFHostIDS).Filter("switch", "on").RelatedSel(4).All(&OFFHostCamera)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "获取离线摄像机配置异常", "Data": nil}
		this.ServeJSON()
		return
	}

	if len(OFFHostCamera) > 0 {
		//计算每台主机分配的摄像机数量
		step := len(OFFHostCamera) / len(ON_Hosts)
		step_y := len(OFFHostCamera) % len(ON_Hosts)
		if len(OFFHostCamera)%len(ON_Hosts) != 0 && step > 0 {
			step += 1
		} else if step == 0 {
			step = 1
		}

		//根据主机下坐标计算分配的数量
		for index, value := range ON_Hosts {
			//截取http://
			v_index := strings.Index(host_ip, "//")
			if v_index != -1 {
				host_ip = host_ip[v_index+2:]
			}

			var IP = value.IP
			i_index := strings.Index(IP, "//")
			if i_index != -1 {
				IP = IP[i_index+2:]
			}

			//设置最大分配上限
			maxCameraNum := MacCameraNum.GetMaxCameraNum() - len(resultMaps)

			if IP == host_ip && index < len(OFFHostCamera) {
				var maps []models.FuncCamera
				//主机下坐标小于余数、余数为0，使用step+1
				if index < step_y || step_y == 0 {
					start := index * step
					end := (index + 1) * step
					DValue := end - start - maxCameraNum
					if DValue > 0 {
						end = end - DValue
					}
					maps = OFFHostCamera[start:end]
				} else if index == step_y { //主机下坐标等于余数，后位值使用step-1
					start := index * step
					end := (index+1)*step - 1
					DValue := end - start - maxCameraNum
					if DValue > 0 {
						end = end - DValue
					}
					maps = OFFHostCamera[start:end]
				} else {
					start := index*(step-1) + step_y
					end := (index+1)*(step-1) + step_y
					DValue := end - start - maxCameraNum
					if DValue > 0 {
						end = end - DValue
					}
					maps = OFFHostCamera[start:end]
				}

				for _, v := range maps {
					resultMaps = append(resultMaps, v)
				}
			}
		}
	}
	if len(resultMaps) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "未查询到摄像机配置", "Data": nil}
		this.ServeJSON()
		return
	}

	var data []orm.Params
	for _, v := range resultMaps {
		data = append(data, ModelFCForInterface(v))
	}

	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "success", "Data": &data}
	this.ServeJSON()
	return
}

type ClientDistributeCamera2_Handler struct {
	beego.Controller
}

func (this *ClientDistributeCamera2_Handler) Get() {
	host_ip := this.GetString("ip")
	o := O
	//查询该主机
	var host models.Host
	LogsInfo("请求主机ip:", host_ip, ",请求时间:", time.Now().Format("2006-01-02 15:04:05"))
	err := o.QueryTable(new(models.Host)).Filter("host_ip__contains", host_ip).One(&host)
	if err != nil {
		//log.Println("获取主机信息异常:", err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "获取主机信息异常", "Data": nil}
		this.ServeJSON()
		return
	}
	host.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	host.Status = "0"
	_, err = o.Update(&host, "updateTime", "host_status")
	if err != nil {
		//log.Println("更新数据库异常:", err)
		//logs.Info("更新主机请求时间异常：",err)
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "更新主机时间异常", "Data": nil}
		this.ServeJSON()
		return
	}

	dict := ReadFile()
	var data []orm.Params

	for v := range dict {
		if v != host.ID || dict[v] == nil {
			continue
		}

		hostfc := dict[v].(map[string]interface{})
		var AMaps []interface{}
		var BMaps []interface{}
		//AMaps := hostfc["AMaps"].([]interface{})
		//BMaps := hostfc["BMaps"].([]interface{})

		if hostfc["AMaps"] != nil {
			AMaps = hostfc["AMaps"].([]interface{})
		}
		if hostfc["BMaps"] != nil {
			BMaps = hostfc["BMaps"].([]interface{})
		}

		for _, v := range AMaps {
			var fc models.FuncCamera
			jsonStr, err := json.Marshal(v)
			if err != nil {
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据转换失败", "Data": nil}
				this.ServeJSON()
				return
			}
			json.Unmarshal(jsonStr, &fc)
			data = append(data, ModelFCForInterface(fc))
		}
		for _, v := range BMaps {
			var fc models.FuncCamera
			jsonStr, err := json.Marshal(v)
			if err != nil {
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据转换失败", "Data": nil}
				this.ServeJSON()
				return
			}
			json.Unmarshal(jsonStr, &fc)
			data = append(data, ModelFCForInterface(fc))
		}
		break
	}
	if len(data) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "未查询到摄像机配置", "Data": nil}
	} else {
		this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "success", "Data": &data}
	}
	this.ServeJSON()
}

type GetOnfuncCaemra struct {
	beego.Controller
}

func (this *GetOnfuncCaemra) Get() {
	dict := ReadFile()

	var data []orm.Params

	for v := range dict {

		var host_data []orm.Params

		hostfc := dict[v].(map[string]interface{})
		var AMaps []interface{}
		var BMaps []interface{}
		//AMaps := hostfc["AMaps"].([]interface{})
		//BMaps := hostfc["BMaps"].([]interface{})

		if hostfc["AMaps"] != nil {
			AMaps = hostfc["AMaps"].([]interface{})
		}
		if hostfc["BMaps"] != nil {
			BMaps = hostfc["BMaps"].([]interface{})
		}

		for _, v := range AMaps {
			var fc models.FuncCamera
			jsonStr, err := json.Marshal(v)
			if err != nil {
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据转换失败", "Data": nil}
				this.ServeJSON()
				return
			}
			json.Unmarshal(jsonStr, &fc)
			fc_data := ModelFCForInterface(fc)
			fc_dic := orm.Params{}
			fc_dic["cameraId"] = fc_data["cameraId"]
			fc_dic["name"] = fc_data["name"]
			fc_dic["cdconfig"] = fc_data["cdconfig"]
			host_data = append(host_data, fc_dic)
		}
		for _, v := range BMaps {
			var fc models.FuncCamera
			jsonStr, err := json.Marshal(v)
			if err != nil {
				this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据转换失败", "Data": nil}
				this.ServeJSON()
				return
			}
			json.Unmarshal(jsonStr, &fc)
			fc_data := ModelFCForInterface(fc)
			fc_dic := orm.Params{}
			fc_dic["cameraId"] = fc_data["cameraId"]
			fc_dic["name"] = fc_data["name"]
			fc_dic["cdconfig"] = fc_data["cdconfig"]
			host_data = append(host_data, fc_dic)
		}

		host_dic := orm.Params{}
		host_dic["host_id"] = v
		host_dic["data"] = host_data
		host_dic["count"] = len(host_data)
		data = append(data, host_dic)
	}
	if len(data) == 0 {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "未查询到摄像机配置", "Data": nil}
	} else {
		this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "success", "Data": &data}
	}
	this.ServeJSON()
}

func ModelFCForInterface(fc models.FuncCamera) (maps orm.Params) {

	rType := reflect.TypeOf(fc)
	rValue := reflect.ValueOf(fc)

	var funcCamera = make(map[string]interface{})

	for i := 0; i < rType.NumField(); i++ {
		k := rType.Field(i)
		v := rValue.Field(i).String()

		if k.Name == "Camera" {
			var camer = fc.Camera
			camera := ModelCameraForInterface(*camer)
			for key, _ := range camera {
				funcCamera[key] = camera[key]
			}
		} else if k.Name == "DistanceMode" || k.Name == "ID" {
			funcCamera[k.Name] = v
		} else if k.Name == "Type" {
			if v == "leakage" {
				v = "liquid_leak"
			}
			funcCamera["function_type"] = v
		} else if k.Name == "Frequency" { // && (fc.Type == "sleep_count" || fc.Type == "queue_count")
			int, _ := strconv.Atoi(v)
			if int == 0 {
				int = 15
			}
			v = strconv.Itoa(int * 60)
			funcCamera[tool.Lowercase(k.Name)] = v
		} else {
			funcCamera[tool.Lowercase(k.Name)] = v
		}
	}

	return funcCamera
}

func ModelCameraForInterface(c models.Camera) (maps orm.Params) {
	rType := reflect.TypeOf(c)
	rValue := reflect.ValueOf(c)

	var camera = make(map[string]interface{})
	for i := 0; i < rType.NumField(); i++ {
		k := rType.Field(i)
		v := rValue.Field(i).String()
		if k.Name == "ID" {
			camera["cameraId"] = v
		} else if k.Name == "Place" {
			camera["location"] = c.Place.Name
		} else if k.Name == "Host" {
			camera["host_id"] = c.Host.ID
		} else if k.Name == "Alarm" || k.Name == "FuncCamera" {
			continue
		} else if k.Name == "VideoCode" && v == "" {
			camera["videoCode"] = "0"
		} else {
			camera[tool.Lowercase(k.Name)] = v
		}
	}

	return camera
}

func ValueNilStr(v reflect.Value) interface{} {
	log.Println("value:", v)
	if v.Len() == 0 {
		return ""
	}
	return v.String()
}

//func ModelOrganizeForInterface(o *models.Organize) (maps orm.Params){
//	rType := reflect.TypeOf(o)
//	rValue := reflect.ValueOf(o)
//
//	var org = make(map[string]interface{})
//	for i:=0;i<rType.NumField();i++ {
//		k := rType.Field(i)
//		v := rValue.Field(i)
//		if k.Name == "User" || k.Name == "Place" {
//			continue
//		}else{
//			org[tool.Capitalize(k.Name)] = v
//		}
//	}
//	return org
//}

//func ModelHostForInterface(h *models.Host) (maps orm.Params){
//	rType := reflect.TypeOf(h)
//	rValue := reflect.ValueOf(h)
//
//	var host = make(map[string]interface{})
//	for i:=0;i<rType.NumField();i++ {
//		k := rType.Field(i)
//		v := rValue.Field(i)
//		if k.Name == "Camera" {
//			continue
//		}else{
//			host[tool.Capitalize(k.Name)] = v
//		}
//	}
//	return host
//}

//负样本存储
type ClientCameraFuncNegative struct {
	beego.Controller
}

func (this *ClientCameraFuncNegative) Get() {
	cameraID, funcType, negative := this.GetString("cameraID"), this.GetString("funcType"), this.GetString("negative")
	if negative == "" {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "负样本为空，保存失败"}
		this.ServeJSON()
		return
	}
	if cameraID == "" || funcType == "" {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "摄像机参数错误"}
		this.ServeJSON()
		return
	}

	o := O
	var maps []models.FuncCamera
	_, err := o.QueryTable(new(models.FuncCamera)).Filter("function_type", funcType).Filter("Camera__ID", cameraID).All(&maps)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "操作失败"}
		this.ServeJSON()
		return
	}

	for _, funcCamer := range maps {
		if funcCamer.Negative != "" && !strings.Contains(funcCamer.Negative, negative) {
			funcCamer.Negative = funcCamer.Negative + negative + ";"
		} else if funcCamer.Negative == "" {
			funcCamer.Negative = negative + ";"
		}
		o.Update(&funcCamer, "negative")
	}

	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "操作成功"}
	this.ServeJSON()
	return

}

// 报警保存
type ClientAlarmSave_Handler struct {
	beego.Controller
}

func (this *ClientAlarmSave_Handler) Post() {
	alarm_detial, alarm_type, alarm_head, alarm_place, alarm_place_type, alarm_file,
		alarm_time, alarm_level, host_id, camera_id, alarm_stream, alarm_video, imageBase, videoBase, cameraNub :=
		this.GetString("alarm_detial"), this.GetString("alarm_type"), this.GetString("alarm_head"), this.GetString("alarm_place"),
		this.GetString("alarm_place_type"), this.GetString("alarm_file"), time.Now().Format("2006-01-02 15:04:05"), this.GetString("alarm_level"),
		this.GetString("host_id"), this.GetString("camera_id"), this.GetString("alarm_stream"), this.GetString("alarm_video"), this.GetString("imageBase"), this.GetString("videoBase"), this.GetString("cameraNub")
	if alarm_detial == "" || alarm_type == "" || alarm_head == "" || camera_id == "" || alarm_file == "" {
		mystruct := &util.ResultC{0, "缺少必要参数", nil}
		this.Data["json"] = mystruct
		this.ServeJSON()
		return
	}
	if alarm_type == "liquid_leak" {
		alarm_type = "leakage"
	}
	alarmNumber := ""
	o := O
	if host_id != "" {
		var maps []orm.Params
		_, err := o.Raw("select * from ss_func_host where host_id=?", host_id).Values(&maps)
		if err != nil {
			mystruct := &util.ResultC{0, "插入数据失败", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}
		if len(maps) > 0 && maps[0]["alarmNumber"] != nil {
			alarmNumber = maps[0]["alarmNumber"].(string)
		}

	}

	id := alarmNumber + time.Now().Format("20060102150405") + CreateRandomNumber(6)
	sqlStr := ""

	sqlStr = "insert into ss_alarm (alarm_id,alarm_detial,alarm_type,alarm_head,alarm_place,alarm_place_type,alarm_file,alarm_time,alarm_level,host_id,camera_id,alarm_stream,alarm_video) values(?,?,?,?,?,?,?,?,?,?,?,?,?)"

	var alarmD models.AlarmDetail
	err := o.QueryTable(new(models.AlarmDetail)).Filter("type", alarm_type).One(&alarmD)
	if alarmD.Detail != "" {
		alarm_detial = alarmD.Detail
	}
	_, err = o.Raw(sqlStr,
		id, alarm_detial, alarm_type, alarm_head, alarm_place, alarm_place_type, alarm_file,
		alarm_time, alarm_level, host_id, camera_id, alarm_stream, alarm_video).Exec()
	if err != nil {
		mystruct := &util.ResultC{0, "插入数据失败", nil}
		this.Data["json"] = mystruct
		this.ServeJSON()
		return
	} else {

		//var alarms []orm.Params
		//_, err = o.Raw("select * from ss_alarm where alarm_id=?", id).Values(&alarms)
		//if err == nil && len(alarms) > 0 {
		//	SendAlarm(alarms[0])
		//}

		//_dir := "static/alarm/" + time.Now().Format("2006-01-02") + "/" + alarm_type
		//exist, err := PathExists(_dir)
		//if err != nil {
		//	fmt.Printf("get dir error![%v]\n", err)
		//	return
		//}
		//if exist {
		//	fmt.Printf("has dir![%v]\n", _dir)
		//} else {
		//	fmt.Printf("no dir![%v]\n", _dir)
		//	// 创建文件夹
		//	err := os.MkdirAll(_dir, os.ModePerm)
		//	if err != nil {
		//		fmt.Printf("mkdir failed![%v]\n", err)
		//	} else {
		//		fmt.Printf("mkdir success!\n")
		//	}
		//}
		//videoBase_url := ""
		//imageBase_url := ""
		var videoBaseDist []byte
		var imageBaseDist []byte
		//timeUnixNano := time.Now().UnixNano()
		if videoBase != "" {
			//videoBase_url = _dir + "/" + "video_" + strconv.FormatInt(timeUnixNano, 10) + ".mp4"
			dist, _ := base64.StdEncoding.DecodeString(videoBase)
			//f, _ := os.OpenFile(videoBase_url, os.O_RDWR|os.O_CREATE, os.ModePerm)
			//defer f.Close()
			//f.Write(dist)
			videoBaseDist = dist
		}
		if imageBase != "" {
			//imageBase_url = _dir + "/" + "pic_" + strconv.FormatInt(timeUnixNano, 10) + ".png"
			dist, _ := base64.StdEncoding.DecodeString(imageBase)
			//f, _ := os.OpenFile(imageBase_url, os.O_RDWR|os.O_CREATE, os.ModePerm)
			//defer f.Close()
			//f.Write(dist)
			imageBaseDist = dist
		}
		uid := uuid.Must(uuid.NewV4()).String()
		uid = strings.Replace(uid, "-", "", -1) // 生成SessionID
		mapsInfo := map[string]string{
			"systemInfo":    SelSystemTitleUtil(),
			"deviceId":      uid,
			"locationId":    alarm_place,
			"eventType":     alarm_type,
			"host_id":       host_id,
			"eventId":       id,
			"cameraId":      cameraNub,
			"priority":      alarm_level,
			"repeatId":      id,
			"info":          alarm_detial,
			"evidenceImg":   string(imageBaseDist[:]),
			"evidenceVideo": string(videoBaseDist[:]),
			"timestamp":     alarm_time,
		}
		//logs.Info(mapsInfo["systemInfo"])
		PostAlarmInfo(mapsInfo)

		mystruct := &util.ResultC{1, "插入成功", nil}
		this.Data["json"] = mystruct
		this.ServeJSON()
		return
	}
}

// 本地存储短视频
type ClientAlarmSaveController struct {
	beego.Controller
}

func (this *ClientAlarmSaveController) Post() {
	alarm_detial, alarm_type, alarm_place, alarm_place_type,
		alarm_time, alarm_level, host_id, camera_id, alarm_stream, imageBase, videoBase, cameraNub, commands, uuids :=
		this.GetString("alarm_detial"),
		this.GetString("alarm_type"),
		this.GetString("alarm_place"),
		this.GetString("alarm_place_type"),
		time.Now().Format("2006-01-02 15:04:05"),
		this.GetString("alarm_level"),
		this.GetString("host_id"),
		this.GetString("camera_id"),
		this.GetString("alarm_stream"),
		this.GetString("imageBase"),
		this.GetString("videoBase"),
		this.GetString("cameraNub"),
		this.GetString("commands"),
		this.GetString("uuids")
	//logs.Info("告警数据参数：", "alarm_detial=", alarm_detial,",alarm_type=",alarm_type,",camera_id=",camera_id,",imageBase=",imageBase,"videoBase=",videoBase)
	//LogsInfo("告警数据参数", "alarm_detial=", alarm_detial, ",alarm_type=", alarm_type, ",camera_id=", camera_id, ",imageBase=", imageBase, "videoBase=", videoBase)

	if len(commands) == 0 {
		beego.Info("预警指令长度为空，则默认为第一次插入")
		commands = "1"
		//mystruct := &util.ResultC{0, "缺少参数-commands", nil}
		//this.Data["json"] = mystruct
		//this.ServeJSON()
		//return
	}

	if len(uuids) == 0 {
		mystruct := &util.ResultC{0, "缺少参数-uuids", nil}
		this.Data["json"] = mystruct
		this.ServeJSON()
		return
	}

	tablename := "ss_alarm"
	if FilterStatus {
		tablename = "ss_alarm_filter"
	}

	beego.Info("预警指令是>>> ", commands)
	if commands == "2" {
		beego.Info(fmt.Sprintf("更新%s报警(新增短视频)！！！", uuids))

		//todo:----------------------------------------------------------------------------------------------------start

		if len(camera_id) == 0 {
			mystruct := &util.ResultC{0, "缺少参数-camera_id", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}

		if len(alarm_detial) == 0 {
			mystruct := &util.ResultC{0, "缺少参数-alarm_detial", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}
		if len(alarm_type) == 0 {
			mystruct := &util.ResultC{0, "缺少参数-alarm_type", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}
		if len(camera_id) == 0 {
			mystruct := &util.ResultC{0, "缺少参数-camera_id", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}
		if videoBase == "" {
			mystruct := &util.ResultC{0, "缺少参数-videoBase", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}

		if alarm_type == "liquid_leak" {
			alarm_type = "leakage"
		}

		alarmNumber := ""
		o := O
		if host_id != "" {
			var maps []orm.Params
			_, err := o.Raw("select * from ss_func_host where host_id=?", host_id).Values(&maps)
			if err != nil {
				LogsError("获取主机告警编号失败，error:", err)
				mystruct := &util.ResultC{0, "插入数据失败", nil}
				this.Data["json"] = mystruct
				this.ServeJSON()
				return
			}
			if len(maps) > 0 && maps[0]["alarmNumber"] != nil {
				alarmNumber = maps[0]["alarmNumber"].(string)
			}

		}

		var camera models.Camera
		err := o.QueryTable(new(models.Camera)).Filter("camera_id", camera_id).RelatedSel().One(&camera)
		if err != nil {
			LogsError("获取摄像机失败，error:", err)
			mystruct := &util.ResultC{0, "获取摄像机失败", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}
		id := alarmNumber + uuids

		// 生成告警目录
		_dir := "static/alarm/" + time.Now().Format("2006-01-02") + "/" + alarm_type + "/" + camera.Place.ID + "/" + camera.ID + "/" + id
		// 查询告警路径是否存在
		exist, err := PathExists(_dir)
		if err != nil {
			LogsError("获取存储路径异常，error:", err)
			mystruct := &util.ResultC{0, "获取存储路径异常", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}

		timeUnixNano := uuids
		videoBase_url := ""
		var videoBaseDist []byte
		if videoBase != "" {
			videoBase_url = _dir + "/" + "video_" + timeUnixNano + ".mp4"
			dist, _ := base64.StdEncoding.DecodeString(videoBase)
			videoBaseDist = dist
		}

		var mapsaA []orm.Params
		num, erra := o.Raw("select * from "+tablename+" where uuids=?;", uuids).Values(&mapsaA)
		if erra != nil {
			beego.Info(uuids, "查询该报警记录失败 >>> result: ", num)
			mystruct := &util.ResultC{0, "更新告警数据失败", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		} else {
			if len(mapsaA) >= 1 {
				if !exist {
					// 创建文件夹
					err := os.MkdirAll(_dir, os.ModePerm)
					if err != nil {
						LogsError("创建存储路径异常，error:", err)
						mystruct := &util.ResultC{0, "创建存储路径异常", nil}
						this.Data["json"] = mystruct
						this.ServeJSON()
						return
					}
				}
			} else {
				beego.Info(fmt.Sprintf("没有该条%s报警记录~此次更新失败！", uuids))
				mystruct := &util.ResultC{0, "更新告警数据失败", nil}
				this.Data["json"] = mystruct
				this.ServeJSON()
				return
			}
		}
		//beego.Info("继续？？？")
		sqlStr := "update " + tablename + " set alarm_video=? where uuids=?"

		_, err = o.Raw(sqlStr,
			"/"+videoBase_url, uuids).Exec()
		if err != nil {
			LogsError("更新告警数据失败，error:", err)
			mystruct := &util.ResultC{0, "更新告警数据失败", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		} else {
			if videoBase != "" {
				f, _ := os.OpenFile(videoBase_url, os.O_RDWR|os.O_CREATE, os.ModePerm)
				defer f.Close()
				_, err := f.Write(videoBaseDist)
				if err != nil {
					LogsError("存储告警视频本地文件失败，error:", err)
					mystruct := &util.ResultC{0, "存储告警视频本地文件失败", nil}
					this.Data["json"] = mystruct
					this.ServeJSON()
					return
				}
			}
		}
		//todo:------------------------------------------------------------------------------------------------------end
		//uid := uuid.Must(uuid.NewV4()).String()
		//uid = strings.Replace(uid, "-", "", -1) // 生成SessionID
		//filePath := "http://" + this.Ctx.Request.Host + "/"+videoBase_url
		//imageBaseDist := FileGetByte64(filePath)
		//mapsInfo := map[string]string{
		//	"systemInfo":    SelSystemTitleUtil(),
		//	"deviceId":      uid,
		//	"locationId":    alarm_place,
		//	"eventType":     alarm_type,
		//	"host_id":       host_id,
		//	"eventId":       id,
		//	"cameraId":      cameraNub,
		//	"priority":      alarm_level,
		//	"repeatId":      id,
		//	"info":          alarm_detial,
		//	"evidenceImg":   string(imageBaseDist[:]),
		//	"evidenceVideo": string(videoBaseDist[:]),
		//	"timestamp":     alarm_time,
		//}
		//mapsInfo2 := map[string]string{
		//	"systemInfo":    SelSystemTitleUtil(),
		//	"deviceId":      uid,
		//	"locationId":    alarm_place,
		//	"eventType":     alarm_type,
		//	"host_id":       host_id,
		//	"eventId":       alarm_stream,
		//	"cameraId":      camera_id,
		//	"priority":      alarm_level,
		//	"repeatId":      id,
		//	"info":          alarm_detial,
		//	"evidenceImg":   string(imageBaseDist[:]),
		//	"evidenceVideo": string(videoBaseDist[:]),
		//	"timestamp":     alarm_time,
		//}
		//PostAlarmInfo(mapsInfo)
		//AlarmToPhoneX(mapsInfo2)
		mystruct := &util.ResultC{1, "更新成功", nil}
		this.Data["json"] = mystruct
		this.ServeJSON()
		return

	} else if commands == "1" {
		beego.Info(fmt.Sprintf("新增%s报警（新增图片）！！！", uuids))

		if len(alarm_detial) == 0 {
			mystruct := &util.ResultC{0, "缺少参数-alarm_detial", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}
		if len(alarm_type) == 0 {
			mystruct := &util.ResultC{0, "缺少参数-alarm_type", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}
		if len(camera_id) == 0 {
			mystruct := &util.ResultC{0, "缺少参数-camera_id", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}
		if imageBase == "" {
			mystruct := &util.ResultC{0, "缺少参数-imageBase", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}

		if videoBase == "" {
			beego.Warning(">>> WARNNING >>> 指令:", commands, "videoBase为空！！！")
		}

		if alarm_type == "liquid_leak" {
			alarm_type = "leakage"
		}

		alarmNumber := ""
		o := O
		if host_id != "" {
			var maps []orm.Params
			_, err := o.Raw("select * from ss_func_host where host_id=?", host_id).Values(&maps)
			if err != nil {
				LogsError("获取主机告警编号失败，error:", err)
				mystruct := &util.ResultC{0, "插入数据失败", nil}
				this.Data["json"] = mystruct
				this.ServeJSON()
				return
			}
			if len(maps) > 0 && maps[0]["alarmNumber"] != nil {
				alarmNumber = maps[0]["alarmNumber"].(string)
			}

		}

		var camera models.Camera
		err := o.QueryTable(new(models.Camera)).Filter("camera_id", camera_id).RelatedSel().One(&camera)
		if err != nil {
			LogsError("获取摄像机失败，error:", err)
			mystruct := &util.ResultC{0, "获取摄像机失败", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}
		id := alarmNumber + uuids
		//id := alarmNumber + time.Now().Format("20060102150405") + CreateRandomNumber(6)

		// 生成告警目录
		_dir := "static/alarm/" + time.Now().Format("2006-01-02") + "/" + alarm_type + "/" + camera.Place.ID + "/" + camera.ID + "/" + id
		// 查询告警路径是否存在
		exist, err := PathExists(_dir)
		if err != nil {
			LogsError("获取存储路径异常，error:", err)
			mystruct := &util.ResultC{0, "获取存储路径异常", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		}
		if !exist {
			// 创建文件夹
			err := os.MkdirAll(_dir, os.ModePerm)
			if err != nil {
				LogsError("创建存储路径异常，error:", err)
				mystruct := &util.ResultC{0, "创建存储路径异常", nil}
				this.Data["json"] = mystruct
				this.ServeJSON()
				return
			}
		}

		imageBase_url := ""

		var imageBaseDist []byte
		timeUnixNano := uuids //time.Now().UnixNano()
		videoBase_url := ""
		var videoBaseDist []byte
		if videoBase != "" {
			videoBase_url = _dir + "/" + "video_" + timeUnixNano + ".mp4"
			dist, _ := base64.StdEncoding.DecodeString(videoBase)
			videoBaseDist = dist
		}
		if imageBase != "" {
			imageBase_url = _dir + "/" + "pic_" + timeUnixNano + ".jpg"
			dist, _ := base64.StdEncoding.DecodeString(imageBase)
			imageBaseDist = dist
		}

		sqlStr := "insert into " + tablename + " (alarm_id,alarm_detial,alarm_type,alarm_head,alarm_place,alarm_place_type,alarm_file,alarm_time,alarm_level,host_id,camera_id,alarm_stream,alarm_video,uuids) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

		//查询预设值告警文案
		var alarmD models.AlarmDetail
		err = o.QueryTable(new(models.AlarmDetail)).Filter("type", alarm_type).One(&alarmD)
		if err != nil && alarmD.Detail != "" {
			alarm_detial = alarmD.Detail
		}

		// 告警存储主机ip
		alarm_head := this.Ctx.Request.Host

		_, err = o.Raw(sqlStr,
			id, alarm_detial, alarm_type, alarm_head, alarm_place, alarm_place_type, "/"+imageBase_url,
			alarm_time, alarm_level, host_id, camera_id, alarm_stream, "/"+videoBase_url, uuids).Exec()
		if err != nil {
			LogsError("插入告警数据失败，error:", err)
			mystruct := &util.ResultC{0, "插入告警数据失败", nil}
			this.Data["json"] = mystruct
			this.ServeJSON()
			return
		} else {
			if videoBase != "" {
				f, _ := os.OpenFile(videoBase_url, os.O_RDWR|os.O_CREATE, os.ModePerm)
				defer f.Close()
				_, err := f.Write(videoBaseDist)
				if err != nil {
					o.Raw("delete from "+tablename+" where alarm_id=?", id).Exec()
					LogsError("存储告警视频本地文件失败，error:", err)
					mystruct := &util.ResultC{0, "存储告警视频本地文件失败", nil}
					this.Data["json"] = mystruct
					this.ServeJSON()
					return
				}
				//uid := uuid.Must(uuid.NewV4()).String()
				//uid = strings.Replace(uid, "-", "", -1) // 生成SessionID
				//mapsInfo := map[string]string{
				//	"systemInfo":    SelSystemTitleUtil(),
				//	"deviceId":      uid,
				//	"locationId":    alarm_place,
				//	"eventType":     alarm_type,
				//	"host_id":       host_id,
				//	"eventId":       id,
				//	"cameraId":      cameraNub,
				//	"priority":      alarm_level,
				//	"repeatId":      id,
				//	"info":          alarm_detial,
				//	"evidenceImg":   string(imageBaseDist[:]),
				//	"evidenceVideo": string(videoBaseDist[:]),
				//	"timestamp":     alarm_time,
				//}
				//mapsInfo2 := map[string]string{
				//	"systemInfo":    SelSystemTitleUtil(),
				//	"deviceId":      uid,
				//	"locationId":    alarm_place,
				//	"eventType":     alarm_type,
				//	"host_id":       host_id,
				//	"eventId":       alarm_stream,
				//	"cameraId":      camera_id,
				//	"priority":      alarm_level,
				//	"repeatId":      id,
				//	"info":          alarm_detial,
				//	"evidenceImg":   string(imageBaseDist[:]),
				//	"evidenceVideo": string(videoBaseDist[:]),
				//	"timestamp":     alarm_time,
				//}
				//PostAlarmInfo(mapsInfo)
				//AlarmToPhoneX(mapsInfo2)
			}
			if imageBase != "" {
				f, _ := os.OpenFile(imageBase_url, os.O_RDWR|os.O_CREATE, os.ModePerm)
				defer f.Close()
				_, err := f.Write(imageBaseDist)
				if err != nil {
					o.Raw("delete from "+tablename+" where alarm_id=?", id).Exec()
					LogsError("存储告警图片本地文件失败，error:", err)
					mystruct := &util.ResultC{0, "存储告警图片本地文件失败", nil}
					this.Data["json"] = mystruct
					this.ServeJSON()
					return
				}
			}

		}

		uid := uuid.Must(uuid.NewV4()).String()
		uid = strings.Replace(uid, "-", "", -1) // 生成SessionID
		mapsInfo := map[string]string{
			"systemInfo":    SelSystemTitleUtil(),
			"deviceId":      uid,
			"locationId":    alarm_place,
			"eventType":     alarm_type,
			"host_id":       host_id,
			"eventId":       id,
			"cameraId":      cameraNub,
			"priority":      alarm_level,
			"repeatId":      id,
			"info":          alarm_detial,
			"evidenceImg":   string(imageBaseDist[:]),
			"evidenceVideo": string(videoBaseDist[:]),
			"timestamp":     alarm_time,
		}
		mapsInfo2 := map[string]string{
			"systemInfo":    SelSystemTitleUtil(),
			"deviceId":      uid,
			"locationId":    alarm_place,
			"eventType":     alarm_type,
			"host_id":       host_id,
			"eventId":       alarm_stream,
			"cameraId":      camera_id,
			"priority":      alarm_level,
			"repeatId":      id,
			"info":          alarm_detial,
			"evidenceImg":   string(imageBaseDist[:]),
			"evidenceVideo": string(videoBaseDist[:]),
			"timestamp":     alarm_time,
		}
		PostAlarmInfo(mapsInfo)
		AlarmToPhoneX(mapsInfo2)

		mystruct := &util.ResultC{1, "插入成功", nil}
		this.Data["json"] = mystruct
		this.ServeJSON()
		return

	} else {
		mystruct := &util.ResultC{0, fmt.Sprintf("报警指令%s有误！！！", commands), nil}
		this.Data["json"] = mystruct
		this.ServeJSON()
		return
	}
}

//type ClientAlarmSaveController2 struct {
//	beego.Controller
//}
//
//func (this *ClientAlarmSaveController2) Post() {
//	alarm_detial, alarm_type, alarm_place, alarm_place_type,
//	alarm_time, alarm_level, host_id, camera_id, alarm_stream, imageBase, videoBase, cameraNub, commands, uuids :=
//		this.GetString("alarm_detial"),
//		this.GetString("alarm_type"),
//		this.GetString("alarm_place"),
//		this.GetString("alarm_place_type"),
//		time.Now().Format("2006-01-02 15:04:05"),
//		this.GetString("alarm_level"),
//		this.GetString("host_id"),
//		this.GetString("camera_id"),
//		this.GetString("alarm_stream"),
//		this.GetString("imageBase"),
//		this.GetString("videoBase"),
//		this.GetString("cameraNub"),
//		this.GetString("commands"),
//		this.GetString("uuids")
//
//	if len(commands) == 0 {
//		beego.Info("预警指令长度为空，则默认为第一次插入")
//		commands = "1"
//
//	}
//
//	if len(uuids) == 0 {
//		mystruct := &util.ResultC{0, "缺少参数-uuids", nil}
//		this.Data["json"] = mystruct
//		this.ServeJSON()
//		return
//	}
//
//	if len(camera_id) == 0 {
//		mystruct := &util.ResultC{0, "缺少参数-camera_id", nil}
//		this.Data["json"] = mystruct
//		this.ServeJSON()
//		return
//	}
//
//	if len(alarm_detial) == 0 {
//		mystruct := &util.ResultC{0, "缺少参数-alarm_detial", nil}
//		this.Data["json"] = mystruct
//		this.ServeJSON()
//		return
//	}
//	if len(alarm_type) == 0 {
//		mystruct := &util.ResultC{0, "缺少参数-alarm_type", nil}
//		this.Data["json"] = mystruct
//		this.ServeJSON()
//		return
//	}
//	if len(camera_id) == 0 {
//		mystruct := &util.ResultC{0, "缺少参数-camera_id", nil}
//		this.Data["json"] = mystruct
//		this.ServeJSON()
//		return
//	}
//
//	if commands == "1" {
//		if imageBase == "" {
//			mystruct := &util.ResultC{0, "缺少参数-imageBase", nil}
//			this.Data["json"] = mystruct
//			this.ServeJSON()
//			return
//		}
//
//		if alarm_type == "liquid_leak" {
//			alarm_type = "leakage"
//		}
//
//	}else if commands == "2" {
//		if videoBase == "" {
//			mystruct := &util.ResultC{0, "缺少参数-videoBase", nil}
//			this.Data["json"] = mystruct
//			this.ServeJSON()
//			return
//		}
//	}else {
//		mystruct := &util.ResultC{0, fmt.Sprintf("报警指令%s有误！！！", commands), nil}
//		this.Data["json"] = mystruct
//		this.ServeJSON()
//		return
//	}
//}

func CreateRandomNumber(len int) string {
	var numbers = []byte{1, 2, 3, 4, 5, 7, 8, 9}
	var container string
	length := bytes.NewReader(numbers).Len()

	for i := 1; i <= len; i++ {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {

		}
		container += fmt.Sprintf("%d", numbers[random.Int64()])
	}
	return container
}

func SetBoundaryStr(boundaryStr string) string {
	kv := strings.Split(boundaryStr, ";")
	result := ""
	for i := 0; i < len(kv)-1; i++ {
		boundary := kv[i]
		boundary = strings.Replace(boundary, "[[", "", -1)
		boundary = strings.Replace(boundary, "]]", "", -1)
		boundarys := strings.Split(boundary, "],[")

		if boundarys[0] != boundarys[len(boundarys)-1] {
			boundarys = append(boundarys, boundarys[0])
			msg := "["
			str := ""
			for i := 0; i < len(boundarys); i++ {
				if i == len(boundarys)-1 {
					str += "[" + boundarys[i] + "]"
				} else {
					str += "[" + boundarys[i] + "],"
				}
			}
			result = msg + str + "];" + result
			//return result
		} else {
			result = boundaryStr
		}
	}
	return result
}
