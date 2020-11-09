package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/wxnacy/wgo/arrays"
	"io/ioutil"
	"net/http"
	"net/url"
	//"reflect"
	"strings"
)

func AlarmToPhoneX(maps map[string]string) {
	DX_Stop := beego.AppConfig.String("DX_Stop")
	if DX_Stop == "1" {
		systems := maps["systemInfo"]
		alarmTypes := maps["eventType"]
		if alarmTypes == "fire" {
			alarmTypes = "厂区火苗监测"
		} else if alarmTypes == "smoke" {
			alarmTypes = "厂区烟雾监测"
		} else if alarmTypes == "sleep_count" {
			alarmTypes = "员工睡岗监测"
		} else if alarmTypes == "queue_count" {
			alarmTypes = "员工离岗监测"
		} else if alarmTypes == "boundary" {
			alarmTypes = "高危区域监测"
		} else if alarmTypes == "cloths" {
			alarmTypes = "员工上岗监测"
		}
		times := maps["timestamp"]
		Infos := maps["info"]
		alarmLocateion := maps["locationId"]
		camera_titles := FindCameraTitle(maps["cameraId"])

		if camera_titles != "" {
			mobiles := FindPromissDX(alarmLocateion) //短信发送号码

			if len(mobiles) > 0 {
				beego.Info(">>> 开始给对应负责人发送短信 >>>")
				beego.Info(camera_titles)
				beego.Info(times)
				beego.Info(systems)
				beego.Info(alarmTypes)
				beego.Info(Infos)
				beego.Info(alarmLocateion)
				for index_xx, xx := range mobiles {
					beego.Info("发送给第", index_xx+1, "人")
					beego.Info(xx)
					DxSend_work(xx, systems, alarmTypes, alarmLocateion, camera_titles, times, Infos)
				}
			} else {
				beego.Info("未找到允许发送短信人员")
			}
		} else {
			beego.Info("未找到该摄像头名称")
		}
	} else {
		beego.Info("短信功能关闭")
	}
}

func FindCameraTitle(camera_ids string) (camera_title string) {
	beego.Info("预警摄像机ID:", camera_ids)
	o := O
	var mapsC []orm.Params
	num, erra := o.Raw("select camera_title from ss_camera where camera_id=?;", camera_ids).Values(&mapsC)
	if erra != nil {
		beego.Info("查询当前camera_title失败 >>> result: ", erra, num)
		return ""
	} else {
		if num >= 1 {
			return mapsC[0]["camera_title"].(string)
		} else {
			return ""
		}
	}
}

func FindPromissDX(localtion_id string) (phones_docker []string) {
	//beego.Info(localtion_id)

	//todo:收集负责人手机号码的过程

	//根据地点名找到组织结构
	o := O
	var mapsA []orm.Params
	num, erra := o.Raw("select organize_id from ss_place where place_name=?;", localtion_id).Values(&mapsA)
	if erra != nil {
		beego.Info("查询当前报警地点ID失败 >>> result: ", num)
	} else {
		if len(mapsA) >= 1 {
			beego.Info("地点组织结构ID:", mapsA[0]["organize_id"])
		} else {
			beego.Info("报警地点:", localtion_id, "没有依附的组织结构!!!")
		}
	}
	var DOCKERS []string
	var mapsB []orm.Params
	nums, errB := o.Raw("select tele,allowd,jurisdiction from ss_user where organize_id=?;", mapsA[0]["organize_id"]).Values(&mapsB)
	if errB != nil {
		beego.Info("查询组织结构下所有用户失败 >>> result: ", nums)
	} else {
		if len(mapsB) >= 1 {

			//beego.Info("当前组织结构下用户数量:",nums)
			//mapsB[0]["tele"],mapsB[0]["allowd"],mapsB[0]["jurisdiction"])
			for _, humans := range mapsB {
				//beego.Info("筛选第",indexB+1,"人")
				//beego.Info(humans["allowd"],humans["tele"])
				//beego.Info(">>>",humans["jurisdiction"])
				existsCode := arrays.Contains(strings.Split(humans["jurisdiction"].(string), ","), "9")
				if (humans["allowd"] == "0") || (existsCode == -1) {
					beego.Info("该用户没有短信权限或被禁用或未填写手机号码")
				} else {
					if humans["tele"] != "" {
						DOCKERS = append(DOCKERS, humans["tele"].(string))
						//beego.Info("此时DOCKER:",DOCKERS)
					}
				}
			}
		} else {
			beego.Info("当前组织结构下无用户!!!")
		}
	}
	return DOCKERS
}

func DxSend_work(mobiles string, systems string, alarmTypes string, alarmLocateion string, videos string, times string, Infos string) {

	DX_ADD := beego.AppConfig.String("DX_ADD")
	u, _ := url.Parse(DX_ADD)
	q := u.Query()
	q.Set("mobile", mobiles)
	paramsBoss := fmt.Sprintf("%s发现一条%s预警,请您及时处理! \n预警摄像头:%s \n预警时间:%s \n预警地点:%s \n预警描述:%s",
		systems, alarmTypes, videos, times, alarmLocateion, Infos)
	q.Set("content", paramsBoss)
	u.RawQuery = q.Encode()
	res, erra := http.Get(u.String())
	if erra != nil {
		beego.Info(erra)
	}
	result, errb := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if errb != nil {
		beego.Info(errb)
	}

	badocker := []byte{}
	for _, b := range result {
		badocker = append(badocker, byte(b))
	}
	beego.Info("短信请求结果:", string(badocker))
	return
}
