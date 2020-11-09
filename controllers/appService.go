package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
)

/*
	定时发送告警信息
*/
func sendAlarm() {
	var alarmlist []orm.Params
	o := O
	_, err := o.Raw("select * from ss_alarm where DATE_SUB(CURDATE(), INTERVAL 30 DAY)<= date(alarm_time) and syn<>1 or syn is null order by alarm_time desc").Values(&alarmlist)
	if err != nil {
		log.Println("读取告警表异常:", err)
		return
	}

	//暂无告警信息
	if len(alarmlist) == 0 {
		return
	}

	//查询数据传送外网，修改状态
	for i := 0; i < len(alarmlist); i++ {
		filePath := "http://" + alarmlist[i]["alarm_head"].(string) + alarmlist[i]["alarm_file"].(string)

		img64 := fileGetByte64(filePath)
		alarmlist[i]["alarm_file"] = string(img64[:])

		if alarmlist[i]["alarm_video"] != "" {
			videoPath := "http://" + alarmlist[i]["alarm_head"].(string) + alarmlist[i]["alarm_video"].(string)
			mp464 := fileGetByte64(videoPath)
			alarmlist[i]["alarm_video"] = string(mp464[:])
		} else {
			alarmlist[i]["alarm_video"] = ""
		}

		jmap := orm.Params{"alarm": alarmlist[i]}
		jsonMap, err := json.Marshal(jmap)
		if err != nil {
			log.Println("请求同步接口异常，", err)
			return
		}
		//logs.Alert("发送数据-alarm：", bytes.NewBuffer(jsonMap))

		resp, err := http.Post(beego.AppConfig.String("appService")+"/syn/alarmdata", "application/json", bytes.NewBuffer(jsonMap))
		if err != nil {
			log.Println("请求同步接口异常，", err)
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		LogsInfo("alarmdata-成功返回：", string(body))

		var mapBody map[string]string
		err = json.Unmarshal(body, &mapBody)
		if err != nil {
			log.Println("转换json异常", err)
		} else if mapBody["code"] == "0" {
			o := O
			o.Raw("update ss_alarm set syn='1' where alarm_id=?", alarmlist[i]["alarm_id"]).Exec()
		}
	}
}

//都到告警计时发送告警信息
func SendAlarm(dit map[string]interface{}) {
	if beego.AppConfig.String("synchronize") == "1" {
		return
	}
	//filePath := "http://" + dit["alarm_head"].(string) + dit["alarm_file"].(string)
	//videoPath := "http://" + dit["alarm_head"].(string) + dit["alarm_video"].(string)
	//img64 := fileGetByte64(filePath)
	//mp464 := fileGetByte64(videoPath)
	//img64 := fileToBase64("./static/alarm/2020-06-02/test/1.jpg")
	//mp464 := fileToBase64("./static/alarm/2020-06-02/test/1.mp4")
	//dit["alarm_file"] = string(img64[:])
	//dit["alarm_video"] = string(mp464[:])
	dit["alarm_file"] = dit["imageBase"]
	dit["alarm_video"] = dit["videoBase"]

	jmap := orm.Params{"alarm": dit}
	jsonMap, err := json.Marshal(jmap)
	if err != nil {
		log.Println("处理同步数据json格式异常，", err)
		return
	}
	//logs.Alert("发送数据-SAlarm：", bytes.NewBuffer(jsonMap))

	resp, err := http.Post(beego.AppConfig.String("appService")+"/syn/alarmdata", "application/json", bytes.NewBuffer(jsonMap))
	if err != nil {
		log.Println("请求同步接口异常，", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	LogsInfo("alarm2-成功返回：", string(body))

	var mapBody map[string]string
	err = json.Unmarshal(body, &mapBody)
	if err != nil {
		log.Println("转换json异常", err)
	} else if mapBody["code"] == "0" {
		o := O
		o.Raw("update ss_alarm set syn_status='1' where alarm_id=?", dit["alarm_id"]).Exec()
	}
}

/*
	定时发送用户信息
*/
func sendUser() {
	o := O
	var userlist []orm.Params

	_, err := o.Raw("select * from ss_user where type != 0").Values(&userlist)
	if err != nil {
		LogsError("同步查询用户信息异常", err)
		return
	}
	if userlist == nil {
		return
	}
	jmap := orm.Params{"users": userlist}
	json, err := json.Marshal(jmap)

	if err != nil {
		LogsError("处理同步数据json格式异常，", err)
		return
	}
	LogsAlert("发送数据-user：", bytes.NewBuffer(json))

	resp, err := http.Post(beego.AppConfig.String("appService")+"/syn/userData", "application/json", bytes.NewBuffer(json))
	if err != nil {
		LogsError("请求同步接口异常，", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	LogsInfo("userdata-成功返回：", string(body))
}

/*
	推送删除用户
*/
func DelUser(userid string) {
	if beego.AppConfig.String("synchronize") == "1" {
		return
	}
	resp, err := http.PostForm(beego.AppConfig.String("appService")+"delEmployee",
		url.Values{"admin_id": {userid}})
	if err != nil {
		log.Println("发送删除成员异常:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	log.Println("成功返回：", string(body))
}

/*
	定时发送组织信息
*/
func sendOrganize() {
	o := O
	var orglist []orm.Params

	_, err := o.Raw("select * from ss_organize").Values(&orglist)
	if err != nil {
		LogsError("同步处理获取组织信息异常", err)
		return
	}
	if orglist == nil {
		return
	}
	jmap := orm.Params{"orgs": orglist}

	json, err := json.Marshal(jmap)
	if err != nil {
		LogsError("处理同步数据json格式异常")
		return
	}

	LogsAlert("发送数据-organize:", string(json))

	resp, err := http.Post(beego.AppConfig.String("appService")+"/syn/orgdata", "application/json", bytes.NewBuffer(json))
	if err != nil {
		LogsError("请求同步接口异常，", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	LogsInfo("orgdata-成功返回：", string(body))
}

func DelOrganize(orgid string) {
	if beego.AppConfig.String("synchronize") == "1" {
		return
	}
	resp, err := http.PostForm(beego.AppConfig.String("appService")+"delOrganize",
		url.Values{"id": {orgid}})
	if err != nil {
		log.Println("发送移除组织异常:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	log.Println("成功返回：", string(body))
}

// 主机
func sendHost() {
	o := O
	var hostlist []orm.Params

	_, err := o.Raw("select * from ss_func_host").Values(&hostlist)
	if err != nil {
		LogsError("同步查询主机信息异常", err)
		return
	}
	if hostlist == nil {
		return
	}
	jmap := orm.Params{"hosts": hostlist}
	json, err := json.Marshal(jmap)

	if err != nil {
		LogsError("处理同步数据json格式异常，", err)
		return
	}
	LogsAlert("发送数据-host：", bytes.NewBuffer(json))

	resp, err := http.Post(beego.AppConfig.String("appService")+"/syn/hostdata", "application/json", bytes.NewBuffer(json))
	if err != nil {
		LogsError("请求同步接口异常，", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	LogsInfo("hostdata-成功返回：", string(body))
}

//摄像头
func sendCamera() {
	o := O
	var cameralist []orm.Params

	_, err := o.Raw("select * from ss_camera").Values(&cameralist)
	if err != nil {
		LogsError("同步查询摄像机信息异常", err)
		return
	}
	if cameralist == nil {
		return
	}
	jmap := orm.Params{"cameras": cameralist}
	json, err := json.Marshal(jmap)

	if err != nil {
		LogsError("处理同步数据json格式异常，", err)
		return
	}
	LogsAlert("发送数据-camera：", bytes.NewBuffer(json))

	resp, err := http.Post(beego.AppConfig.String("appService")+"/syn/cameradata", "application/json", bytes.NewBuffer(json))
	if err != nil {
		LogsError("请求同步接口异常，", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	LogsInfo("cameradata-成功返回：", string(body))
}

// 区域
func sendPlace() {
	o := O
	var placelist []orm.Params

	_, err := o.Raw("select * from ss_place").Values(&placelist)
	if err != nil {
		LogsError("同步查询摄像机信息异常", err)
		return
	}
	if placelist == nil {
		return
	}
	jmap := orm.Params{"places": placelist}
	json, err := json.Marshal(jmap)

	if err != nil {
		LogsError("处理同步数据json格式异常，", err)
		return
	}
	LogsAlert("发送数据-place：", bytes.NewBuffer(json))

	resp, err := http.Post(beego.AppConfig.String("appService")+"/syn/placedata", "application/json", bytes.NewBuffer(json))
	if err != nil {
		LogsError("请求同步接口异常，", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	LogsInfo("place-成功返回：", string(body))
}

/*
	定时发送批注信息
*/
//func sendNotation() {
//	o := orm.NewOrm()
//	notationDB := new(models.Notation)
//	var notationList []models.Notation
//
//	_, err := o.QueryTable(notationDB).All(&notationList)
//	if err != nil {
//		log.Println("读取批注异常:", err)
//		return
//	}
//
//	for _, notation := range notationList {
//		resp, err := http.PostForm(beego.AppConfig.String("appService")+"synNotation",
//			url.Values{
//				"id":         {notation.ID},
//				"alarm_id":   {notation.Alarm.ID},
//				"notation":   {notation.Notation},
//				"userid":     {notation.User.ID},
//				"createtime": {notation.Createtime}})
//		if err != nil {
//			log.Println("发送成批注常:", err)
//			continue
//		}
//		defer resp.Body.Close()
//
//		body, err := ioutil.ReadAll(resp.Body)
//		log.Println("成功返回：", string(body))
//	}
//}

func FileGetByte64(url string) (b []byte) {
	return fileGetByte64(url)
}

func fileGetByte64(url string) (b []byte) {

	fileName := path.Base(url)
	println(fileName)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("获取文件异常")
		return
	}
	// defer后的为延时操作，通常用来释放相关变量
	defer res.Body.Close()

	//// 获得get请求响应的reader对象
	//reader := bufio.NewReaderSize(res.Body, 32 * 1024)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("读取文件异常!")
		return
	}

	n := base64.StdEncoding.EncodedLen(len(body)) //计算加密后数据长度
	bufstore := make([]byte, n)                   //数据缓存
	base64.StdEncoding.Encode(bufstore, body)     // 文件转base64

	//ddd, _ := base64.StdEncoding.DecodeString(string(bufstore[:])) //成图片文件并把文件写入到buffer
	//err = ioutil.WriteFile("./"+fileName, ddd, 0666)   //buffer输出到jpg文件中（不做处理，直接写到文件）
	//if err != nil{
	//	fmt.Println("生成文件异常!")
	//	return
	//}
	//defer os.Remove("./"+fileName)

	//file,err := ioutil.TempFile("tests",fileName)
	//if err != nil {
	//	fmt.Println("生成临时文件失败")
	//	return
	//}
	//
	////defer os.Remove(file.Name())
	//_,err = file.Write(bufstore)
	//if err != nil {
	//	fmt.Println("写入失败")
	//	return
	//}

	//err = ioutil.WriteFile("./output2.jpg.txt", bufstore, 0666) //直接写入到文件就ok完活了。
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	return bufstore
}

func AppServiceMain() {
	if beego.AppConfig.String("synchronize") == "0" {
		initTask()
	}
}

func taskSendAlarm() error {
	sendAlarm()
	return nil
}

func taskSendUser() error {
	sendUser()
	return nil
}

func taskSenOrganize() error {
	sendOrganize()
	return nil
}

func taskSendPlace() error {
	sendPlace()
	return nil
}

func taskSendHost() error {
	sendHost()
	return nil
}

func taskSendCamera() error {
	sendCamera()
	return nil
}

func initTask() {
	//没个n分钟执行一次 0 */n * * * *
	//没个n秒执行一次 0/n * * * * *
	//没个n分钟执行一次 0 */n * * * *
	taskSendUser := toolbox.NewTask("taskSendUser", "0 */3 * * * *", taskSendUser)
	taskSenOrganize := toolbox.NewTask("taskSenOrganize", "0 */3 * * * *", taskSenOrganize)
	taskSendPlace := toolbox.NewTask("taskSendPlace", "0 */3 * * * *", taskSendPlace)
	taskSendHost := toolbox.NewTask("taskSendHost", "0 */3 * * * *", taskSendHost)
	taskSendCamera := toolbox.NewTask("taskSendCamera", "0 */3 * * * *", taskSendCamera)
	taskSendAlarm := toolbox.NewTask("taskSendAlarm", "0 */3 * * * *", taskSendAlarm)

	//加入全局的计划任务列表
	toolbox.AddTask("taskSendUser", taskSendUser)
	toolbox.AddTask("taskSenOrganize", taskSenOrganize)
	toolbox.AddTask("taskSendPlace", taskSendPlace)
	toolbox.AddTask("taskSendHost", taskSendHost)
	toolbox.AddTask("taskSendCamera", taskSendCamera)
	toolbox.AddTask("taskSendAlarm", taskSendAlarm)
	//开始执行全局的任务
	toolbox.StartTask()
}
