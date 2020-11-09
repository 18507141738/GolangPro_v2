package controllers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func sendAlarmInfoWeb(weburl string, maps map[string]string) {
	timeout := time.Duration(30 * time.Second)
	client := http.Client{Timeout: timeout}

	buf, err := json.MarshalIndent(maps, "", "    ") //格式化编码
	if err != nil {
		LogsInfo("json转换错误", err)
		return
	}
	resp, err := client.Post(weburl, "application/json", strings.NewReader(string(buf)))

	//resp, err := client.PostForm(weburl, url.Values{
	//	"systemInfo": { maps["systemInfo"] },
	//	"deviceId": { maps["deviceId"] },
	//	"locationId": { maps["locationId"] },
	//	"eventType": { maps["eventType"] },
	//	"eventId": { maps["eventId"] },
	//	"cameraId": { maps["cameraId"] },
	//	"priority": { maps["priority"] },
	//	"repeatId": { maps["repeatId"] },
	//	"info": { maps["info"] },
	//	"evidenceImg": { maps["evidenceImg"] },
	//	"evidenceVideo": { maps["evidenceVideo"] },
	//	"timestamp": { maps["timestamp"] },
	//})

	if err != nil {
		LogsInfo("第三方对接失败,对接地址:", weburl, "错误信息:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	LogsInfo("对接第三方接口成功,返回值:", string(body))
}

func PostAlarmInfo(maps map[string]string) {
	if beego.AppConfig.String("webserviceHD") == "0" {
		str := beego.AppConfig.String("HDAlarmType")
		arr := strings.Split(str, ",")
		for _, t := range arr {
			if t == maps["eventType"] {
				sendAlarmInfoWeb(beego.AppConfig.String("HDService"), maps)
				continue
			}
		}
	}
	if beego.AppConfig.String("webserviceSW") == "0" {
		str := beego.AppConfig.String("SWAlarmType")
		arr := strings.Split(str, ",")
		for _, t := range arr {
			if t == maps["eventType"] {
				sendAlarmInfoWeb(beego.AppConfig.String("SWService"), maps)
				continue
			}
		}
	}
}

func TestPostImage() {
	imgByte := fileToBase64("./static/alarm/2020-06-02/test/1.jpg")
	videoByte := fileToBase64("./static/alarm/2020-06-02/test/1.mp4")

	mapsInfo := map[string]string{
		"systemInfo":    SelSystemTitleUtil(),
		"deviceId":      "123123123123",
		"locationId":    "炼化厂区",
		"eventType":     "cloths",
		"eventId":       "0987654322",
		"cameraId":      "003687",
		"priority":      "1",
		"repeatId":      "0987654322",
		"info":          "有员工未佩戴安全帽",
		"evidenceImg":   string(imgByte[:]),
		"evidenceVideo": string(videoByte[:]),
		"timestamp":     "2020-06-08 12:00:00",
	}
	PostAlarmInfo(mapsInfo)
}

func fileToBase64(filepath string) (base64str string) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	base64Str := base64.StdEncoding.EncodeToString(data)

	return base64Str
}
