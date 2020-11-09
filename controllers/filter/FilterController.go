package filter

import (
	"Artifice_V2.0.0/controllers"
	"Artifice_V2.0.0/models"
	"Artifice_V2.0.0/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"log"
	"strconv"
	"strings"
	"time"
)

type FilterHomePage struct {
	beego.Controller
}

func (this *FilterHomePage) Get() {
	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	this.Data["bmk"] = "homemenu"
	this.Data["Time"] = controllers.GetTimeSE()
	admin_id := this.Ctx.GetCookie("filter_admin_id")
	fu, err := models.SelFilterUserByID(controllers.O, admin_id)
	if err != nil {
		log.Println(err)
		this.TplName = "filter/login.html"
	} else {
		this.Data["admin_user"] = fu.Acount
		this.Data["organize_name"] = fu.Organize.Name
		this.Data["jurisdiction"] = fu.Jurisdiction
		this.TplName = "filter/index.html"
	}
}

type FilterFuncPage struct {
	beego.Controller
}

func (this *FilterFuncPage) Get() {
	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	t := this.GetString("type")
	this.Data["bmk"] = titleTypeByParamType(t)
	this.Data["Time"] = controllers.GetTimeSE()
	admin_id := this.Ctx.GetCookie("filter_admin_id")
	fu, err := models.SelFilterUserByID(controllers.O, admin_id)
	if err != nil {
		log.Println(err)
		this.TplName = "filter/login.html"
	} else {
		this.Data["admin_user"] = fu.Acount
		this.Data["organize_name"] = fu.Organize.Name
		this.Data["jurisdiction"] = fu.Jurisdiction
		this.TplName = "filter/funcPage.html"
	}
}

type FilterFuncPageIE struct {
	beego.Controller
}

func (this *FilterFuncPageIE) Get() {
	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	t := this.GetString("type")
	this.Data["bmk"] = titleTypeByParamType(t)
	this.Data["Time"] = controllers.GetTimeSE()
	admin_id := this.Ctx.GetCookie("filter_admin_id")
	fu, err := models.SelFilterUserByID(controllers.O, admin_id)
	if err != nil {
		log.Println(err)
		this.TplName = "filter/login.html"
	} else {
		this.Data["admin_user"] = fu.Acount
		this.Data["organize_name"] = fu.Organize.Name
		this.Data["jurisdiction"] = fu.Jurisdiction
		this.TplName = "filter/funcPageIE.html"
	}
}

type FilterHomeDataHandler struct {
	beego.Controller
}

func (this *FilterHomeDataHandler) Post() {
	admin_id := this.Ctx.GetCookie("filter_admin_id")

	var startTime = time.Now().AddDate(0, 0, -30).Format("2006-01-02") + " 00:00:00"
	var endTime = time.Now().Format("2006-01-02") + " 23:59:59"

	fu, err := models.SelFilterUserByID(controllers.O, admin_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "查询用户信息异常"}
		this.ServeJSON()
		return
	}
	if fu.Organize.ID == "" {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户组织架构异常"}
		this.ServeJSON()
		return
	}

	if fu.Jurisdiction == "" {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户账号无权限"}
		this.ServeJSON()
		return
	}

	jurs := strings.Split(fu.Jurisdiction, ",")

	if len(jurs) == 0 {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户账号无权限"}
		this.ServeJSON()
		return
	}

	orgIds, err := controllers.SelOrgIDBelongUserOrg(fu.Organize)
	if err != nil {
		controllers.LogsError("查询用户下属组织异常:", err)

	}
	var ids = []string{fu.Organize.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}

	countInfo := make(orm.Params)
	jurc := [...]string{"cloths", "smoke", "fire", "boundary", "queue_count", "sleep_count", "leakage"}
	var z_num = []string{}
	var z_key = []string{}
	var c_num = []string{}
	for _, v := range jurs {
		if v == "7" || v == "8" || v == "9" {
			continue
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Println(err)
			this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "用户权限异常"}
			this.ServeJSON()
			return
		}

		//查询用户下的主机特定区域的告警统计
		cnt, err := models.GetFilterAlarmCountByType(controllers.O, jurc[n], ids)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "分类统计异常"}
			this.ServeJSON()
			return
		}

		switch v {
		case "0":
			countInfo["zh"] = cnt
			z_num = append(z_num, "IFNULL(a.num_zh,0) as num_zh")
			c_num = append(c_num, "COUNT(alarm_type = 'cloths' OR NULL) AS num_zh")
			z_key = append(z_key, "num_zh")
		case "1":
			countInfo["yw"] = cnt
			z_num = append(z_num, "IFNULL(a.num_yw,0) as num_yw")
			c_num = append(c_num, "COUNT(alarm_type = 'smoke' OR NULL) AS num_yw")
			z_key = append(z_key, "num_yw")
		case "2":
			countInfo["hm"] = cnt
			z_num = append(z_num, "IFNULL(a.num_hm,0) as num_hm")
			c_num = append(c_num, "COUNT(alarm_type = 'fire' OR NULL) AS num_hm")
			z_key = append(z_key, "num_hm")
		case "3":
			countInfo["qy"] = cnt
			z_num = append(z_num, "IFNULL(a.num_qy,0) as num_qy")
			c_num = append(c_num, "COUNT(alarm_type = 'boundary' OR NULL) AS num_qy")
			z_key = append(z_key, "num_qy")
		case "4":
			countInfo["lg"] = cnt
			z_num = append(z_num, "IFNULL(a.num_lg,0) as num_lg")
			c_num = append(c_num, "COUNT(alarm_type = 'queue_count' OR NULL) AS num_lg")
			z_key = append(z_key, "num_lg")
		case "5":
			countInfo["sg"] = cnt
			z_num = append(z_num, "IFNULL(a.num_sg,0) as num_sg")
			c_num = append(c_num, "COUNT(alarm_type = 'sleep_count' OR NULL) AS num_sg")
			z_key = append(z_key, "num_sg")
			//case "6":
			//	countInfo["xl"] = cnt
			//	z_num = append(z_num, "IFNULL(a.num_xl,0) as num_xl")
			//	c_num = append(c_num, "COUNT(alarm_type = 'leakage' OR NULL) AS num_xl")
			//	z_key = append(z_key, "num_xl")
		}
	}

	var z_numstr = ""
	for i, v := range z_num {
		if i == 0 {
			z_numstr = ","
		}
		z_numstr += v
		if i < len(z_num)-1 {
			z_numstr += ","
		}
	}
	var z_keystr = ""
	for i, v := range z_key {
		z_keystr += v
		if i < len(z_key)-1 {
			z_keystr += "+"
		}
	}
	if len(z_keystr) > 0 {
		z_keystr = ",IFNULL(" + z_keystr + ",0) as num "
	}

	var c_numstr = ""
	for i, v := range c_num {
		if i == 0 {
			c_numstr = ","
		}
		c_numstr += v
		if i < len(c_num)-1 {
			c_numstr += ","
		}
	}

	var sql = "SELECT mm.aa as brief_time " +
		z_numstr + " " +
		z_keystr +
		" FROM  " +
		" ( SELECT DATE_FORMAT(DATE_SUB(CONVERT(?,datetime),INTERVAL xc-1 DAY),'%m.%d') as aa," + //获取xc-1天的日期“月日”
		" DATE_FORMAT( DATE_SUB(CONVERT(?,datetime), INTERVAL xc - 1 DAY), '%Y-%m-%d' ) AS bb " + //获取xc-1天的日期“年月日”
		" FROM" +
		" ( SELECT @xi:=@xi+1 as xc " + //设置xi自加1 赋值给xc
		" FROM " +
		" (SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6) xc1, " + //查询6次
		" (SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6) xc2, " + //查询6次 在上6次中，共6*6=36次
		" (SELECT @xi:=0) xc0 ) xcxc " + //xi 默认0
		" WHERE " +
		" xc<=DATEDIFF(CONVERT(?,datetime),CONVERT(?,datetime))+1 )mm " + //判断条件 计算两个时间的相差天数
		" LEFT OUTER JOIN " + //向左加入
		"( SELECT " +
		" DATE_FORMAT(alarm_time, '%m.%d') as brief_time " +
		c_numstr +
		" FROM ss_alarm_filter ss " +
		" INNER JOIN ss_camera sc ON sc.camera_id = ss.camera_id " + //关联摄像头 查询对应告警下的摄像头
		" INNER JOIN ss_place sp ON sp.place_id = sc.place_id " +
		" INNER JOIN ss_organize so ON so.id = sp.organize_id " +
		" WHERE " +
		"so.id in ("
	var sqlc = ""
	for i := 0; i < len(ids); i++ {
		sqlc += ",?"
	}
	sqlc = sqlc[1:len(sqlc)]
	sql += sqlc

	sql += ") AND alarm_time>=? and alarm_time<=? and ss.`status`=0 " +
		"GROUP BY brief_time)a " +
		"ON a.brief_time = mm.aa ORDER BY mm.bb asc"
	var maps []orm.Params
	o := controllers.O
	_, err = o.Raw(sql, endTime, endTime, endTime, startTime, ids, startTime, endTime).Values(&maps)
	if err != nil {
		log.Println(err)
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "查询统计信息异常"}
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"Code": 1, "Jur": jurs, "AC": countInfo, "Chart": maps}
	this.ServeJSON()
}

type FilterFuncLineHandler struct {
	beego.Controller
}

func (this *FilterFuncLineHandler) Post() {
	zId, funcType := this.GetString("zId"), this.GetString("type")

	if funcType == "" || len(funcType) == 0 {
		this.Data["json"] = &util.ResultB{0, 0, nil}
		this.ServeJSON()
		return
	}

	userid := this.Ctx.GetCookie("filter_admin_id")
	var startTime = time.Now().AddDate(0, 0, -30).Format("2006-01-02") + " 00:00:00"
	var endTime = time.Now().Format("2006-01-02") + " 23:59:59"

	o := controllers.O

	fu, err := models.SelFilterUserByID(o, userid)
	if err != nil {
		this.Data["json"] = &util.ResultB{0, 0, nil}
		this.ServeJSON()
		return
	}

	if fu.Organize.ID == "" {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Num": 0, "Reason": "用户组织架构异常"}
		this.ServeJSON()
		return
	}

	var org = fu.Organize
	var organize models.Organize

	if len(zId) > 0 && fu.Organize.ID != zId {
		err = o.QueryTable(new(models.Organize)).Filter("id", zId).One(&organize)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "组织信息异常"}
			this.ServeJSON()
			return
		}
		if controllers.InOrganize(&organize, fu.Organize) {
			org = &organize
		}
	}

	orgIds, err := controllers.SelOrgIDBelongUserOrg(org)
	if err != nil {
		controllers.LogsError("查询用户下属组织异常:", err)
	}
	var ids = []string{fu.Organize.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}

	var sql = "SELECT mm.aa as brief_time,IFNULL(a.num,0) as num " +
		"FROM  " +
		"( SELECT DATE_FORMAT(DATE_SUB(CONVERT(?,datetime),INTERVAL xc-1 DAY),'%m%d') as aa," +
		" DATE_FORMAT( DATE_SUB(CONVERT(?,datetime), INTERVAL xc - 1 DAY), '%Y-%m-%d' ) AS bb " +
		" FROM" +
		"( SELECT @xi:=@xi+1 as xc " +
		"FROM " +
		"(SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6) xc1, " +
		"(SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6) xc2, " +
		"(SELECT @xi:=0) xc0 ) xcxc " +
		"WHERE " +
		"xc<=DATEDIFF(CONVERT(?,datetime),CONVERT(?,datetime))+1 )mm " +
		"LEFT OUTER JOIN " +
		"(SELECT " +
		"DATE_FORMAT(alarm_time, '%m%d') as brief_time," +
		"COUNT(*) AS num " +
		"FROM ss_alarm_filter ss " +
		"INNER JOIN ss_camera sc ON sc.camera_id = ss.camera_id " +
		" INNER JOIN ss_place sp ON sp.place_id = sc.place_id " +
		" INNER JOIN ss_organize so ON so.id = sp.organize_id " +
		"WHERE " +
		"so.id in ("
	var sqlc = ""
	for i := 0; i < len(ids); i++ {
		sqlc += ",?"
	}
	sqlc = sqlc[1:len(sqlc)]
	sql += sqlc

	sql += ") AND alarm_time>=? and alarm_time<=? and ss.`status`=0 " +
		"AND alarm_type = ? " +
		"GROUP BY brief_time)a " +
		"ON a.brief_time = mm.aa ORDER BY mm.bb asc"
	var maps []orm.Params
	_, err = o.Raw(sql, endTime, endTime, endTime, startTime, ids, startTime, endTime, funcType).Values(&maps)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "查询信息异常"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = &util.ResultB{1, 0, &maps}
	this.ServeJSON()
}

type FilterFuncAlarmsHandler struct {
	beego.Controller
}

func (this *FilterFuncAlarmsHandler) Post() {
	limits, pages, dateTime, zId, funcType, cameraID, status :=
		this.GetString("limit"),
		this.GetString("page"),
		this.GetString("dateTime"),
		this.GetString("zId"),
		this.GetString("type"),
		this.GetString("cameraID"),
		this.GetString("status")
	if funcType == "" || len(funcType) == 0 {
		log.Println("算法类型为空")
		this.Data["json"] = &util.ResultB{0, 0, nil}
		this.ServeJSON()
		return
	}
	start_page, err := strconv.Atoi(pages)
	limitss, err := strconv.Atoi(limits)
	start_pages := (start_page - 1) * limitss
	var startTime = ""
	var endTime = ""

	times := strings.Split(dateTime, " - ")
	if len(dateTime) > 0 || len(times) > 1 {
		startTime = times[0] + " 00:00:00"
		endTime = times[1] + " 23:59:59"

		day := controllers.TimeDifferDay(times[0], times[1])

		//超过31天或比对异常
		if day < 0 || day > 31 {
			loc, _ := time.LoadLocation("Local") //获取本地时区
			eTime, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, loc)
			if err != nil {
				//时间转换异常从当前时间查询
				startTime = time.Now().AddDate(0, 0, -30).Format("2006-01-02") + " 00:00:00"
				endTime = time.Now().Format("2006-01-02") + " 23:59:59"
			} else {
				//超过的时间不显示
				startTime = eTime.AddDate(0, 0, -30).Format("2006-01-02") + " 00:00:00"
			}
		}
	} else { // 默认查询当天记录
		startTime = time.Now().AddDate(0, 0, -30).Format("2006-01-02") + " 00:00:00"
		endTime = time.Now().Format("2006-01-02") + " 23:59:59"
	}

	userid := this.Ctx.GetCookie("filter_admin_id")
	o := controllers.O

	fu, err := models.SelFilterUserByID(o, userid)
	if err != nil {
		this.Data["json"] = &util.ResultB{0, 0, nil}
		this.ServeJSON()
		return
	}

	if fu.Organize.ID == "" {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Num": 0, "Reason": "用户组织架构异常"}
		this.ServeJSON()
		return
	}

	var org = fu.Organize
	var organize models.Organize

	if len(zId) > 0 && fu.Organize.ID != zId {
		err = o.QueryTable(new(models.Organize)).Filter("id", zId).One(&organize)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "组织信息异常"}
			this.ServeJSON()
			return
		}
		if controllers.InOrganize(&organize, fu.Organize) {
			org = &organize
		}
	}

	orgIds, err := controllers.SelOrgIDBelongUserOrg(org)
	if err != nil {
		controllers.LogsError("查询用户下属组织异常:", err)
	}
	var ids = []string{fu.Organize.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}

	maps, num, err := models.FilterAlarmPage(controllers.O, cameraID, funcType, ids, startTime, endTime, limitss, start_pages, status)
	if err != nil {
		this.Data["json"] = &util.ResultB{0, 0, nil}
		this.ServeJSON()
		return
	}
	this.Data["json"] = map[string]interface{}{"Code": 0, "Num": num, "Reason": &maps}
	this.ServeJSON()
}

type FilterSelNelAlarmHandler struct {
	beego.Controller
}

func (this *FilterSelNelAlarmHandler) Post() {
	admin_id := this.Ctx.GetCookie("filter_admin_id")
	fu, err := models.SelFilterUserByID(controllers.O, admin_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "查询失败"}
		this.ServeJSON()
		return
	}
	orgIds, err := controllers.SelOrgIDBelongUserOrg(fu.Organize)
	if err != nil {
		controllers.LogsError("查询用户下属组织异常:", err)
	}
	var ids = []string{fu.Organize.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}

	maps, num, err := models.FilterNewAlarm(controllers.O, ids)
	if err != nil {
		mystruct := map[string]interface{}{"Code": 1, "Num": 0, "Reason": "处理失败"}
		this.Data["json"] = mystruct
		this.ServeJSON()
		return
	}

	if num == 0 || maps == nil {
		mystruct := map[string]interface{}{"Code": 1, "Num": 0, "Reason": "无数据"}
		this.Data["json"] = mystruct
		this.ServeJSON()
		return
	}

	mystruct := map[string]interface{}{"Code": 0, "Num": num, "Reason": &maps}
	this.Data["json"] = mystruct
	this.ServeJSON()
}

type FilterInsertAlarmHandler struct {
	beego.Controller
}

func (this *FilterInsertAlarmHandler) Post() {
	alarmId := this.GetString("alarmId")
	fa, err := models.FilterAlarmByID(controllers.O, alarmId)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败!,再来一次"}
		this.ServeJSON()
		return
	}

	//if fa.AVideo == "" {
	//	this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "短视频不存在提交失败"}
	//	this.ServeJSON()
	//	return
	//}

	alarm := models.Alarm{
		ID:         fa.ID,
		AType:      fa.AType,
		APlaceType: fa.APlaceType,
		APlace:     fa.APlace,
		ADetial:    fa.ADetial,
		AFile:      fa.AFile,
		AVideo:     fa.AVideo,
		AStream:    fa.AStream,
		AHead:      fa.AHead,
		Atime:      fa.Atime,
		ALevel:     fa.ALevel,
		Hostid:     fa.Hostid,
		Astatus:    "0",
		Syn:        "0",
		PageRead:   "0",
		Handler:    "0",
		Camera:     fa.Camera,
	}

	o := controllers.O

	_, err = o.Insert(&alarm)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "数据操作失败!,再来一次"}
		this.ServeJSON()
		return
	}

	fa.Status = "1"
	models.FilterUpdateAlarmStatus(o, fa)

	if beego.AppConfig.String("webserviceHD") == "0" || beego.AppConfig.String("webserviceSW") == "0" {
		var imageBaseDist []byte
		var videoBaseDist []byte
		filePath := "http://" + alarm.AHead + alarm.AFile
		imageBaseDist = controllers.FileGetByte64(filePath)
		if alarm.AVideo != "" {
			videoPath := "http://" + alarm.AHead + alarm.AVideo
			videoBaseDist = controllers.FileGetByte64(videoPath)
		}

		uid := uuid.Must(uuid.NewV4()).String()
		uid = strings.Replace(uid, "-", "", -1) // 生成SessionID
		mapsInfo := map[string]string{
			"systemInfo":    controllers.SelSystemTitleUtil(),
			"deviceId":      uid,
			"locationId":    alarm.APlace,
			"eventType":     alarm.AType,
			"host_id":       alarm.Hostid,
			"eventId":       alarm.ID,
			"cameraId":      alarm.Camera.CameraNub,
			"priority":      alarm.ALevel,
			"repeatId":      alarm.ID,
			"info":          alarm.ADetial,
			"evidenceImg":   string(imageBaseDist[:]),
			"evidenceVideo": string(videoBaseDist[:]),
			"timestamp":     alarm.Atime,
		}
		mapsInfo2 := map[string]string{
			"systemInfo":    controllers.SelSystemTitleUtil(),
			"deviceId":      uid,
			"locationId":    alarm.APlace,
			"eventType":     alarm.AType,
			"host_id":       alarm.Hostid,
			"eventId":       alarm.AStream,
			"cameraId":      alarm.Camera.ID,
			"priority":      alarm.ALevel,
			"repeatId":      alarm.ID,
			"info":          alarm.ADetial,
			"evidenceImg":   string(imageBaseDist[:]),
			"evidenceVideo": string(videoBaseDist[:]),
			"timestamp":     alarm.Atime,
		}
		controllers.PostAlarmInfo(mapsInfo)
		controllers.AlarmToPhoneX(mapsInfo2)
	}

	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "插入成功"}
	this.ServeJSON()
}

type FilterMisinformationHandler struct {
	beego.Controller
}

func (this *FilterMisinformationHandler) Post() {
	alarmId := this.GetString("alarmId")
	fa := models.FilterAlarm{ID: alarmId, Status: "2"}
	_, err := models.FilterUpdateAlarmStatus(controllers.O, fa)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "修改失败"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "修改成功"}
	this.ServeJSON()
}

// 组织下拉
type FilterOrgainzeByUserHandler struct {
	beego.Controller
}

func (this *FilterOrgainzeByUserHandler) Post() {
	admin_id := this.Ctx.GetCookie("filter_admin_id")
	fu, err := models.SelFilterUserByID(controllers.O, admin_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "查询失败"}
		this.ServeJSON()
		return
	}
	orgIds, err := controllers.SelOrgIDBelongUserOrg(fu.Organize)
	if err != nil {
		controllers.LogsError("查询用户下属组织异常:", err)
	}
	var ids = []string{fu.Organize.ID}
	for _, v := range orgIds {
		ids = append(ids, v["id"].(string))
	}
	o := controllers.O
	var orgs []models.Organize
	_, err = o.QueryTable(new(models.Organize)).Filter("id__in", ids).OrderBy("level").All(&orgs)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Code": 0, "Reason": "查询失败"}
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"Code": 1, "Reason": orgs}
	this.ServeJSON()
}

type SelCameraForFuncType struct {
	beego.Controller
}

func (this *SelCameraForFuncType) Post() {
	funcType := this.GetString("funcType")
	o := controllers.O
	var cameraList []models.Camera
	_, err := o.QueryTable(new(models.Camera)).Filter("FuncCamera__Type", funcType).All(&cameraList)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"Ret": 0, "Reason": "查询失败！", "Data": nil}
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]interface{}{"Ret": 1, "Reason": "查询成功！", "Data": cameraList}
	this.ServeJSON()
	return
}

func titleTypeByParamType(t string) string {
	switch t {
	case "cloths":
		return "zhjc"
	case "fire":
		return "hmjc"
	case "boundary":
		return "qyjc"
	case "smoke":
		return "ywjc"
	case "queue_count":
		return "lgjc"
	case "leakage":
		return "xljc"
	case "sleep_count":
		return "sgjc"

	}
	return ""
}
