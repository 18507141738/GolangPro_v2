package signs

import (
	"Artifice_V2.0.0/controllers"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	"github.com/satori/go.uuid"
	"github.com/wxnacy/wgo/arrays"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type CoreDocker struct {
	Names       string
	Codes       string
	ParentCodes string
}

func Inits() {
	go Insk1()
}

func Insk1() {
	rsyncStop := beego.AppConfig.String("rsyncStop")
	if rsyncStop == "1" {
		beego.Info("****应用系统从IAM更新最新的人员(组织结构)主数据*****")
		//每隔n时执行一次 * * */n * * *
		//每隔n分执行一次 0 */n * * * *
		//每隔n秒执行一次 0/n * * * * *
		//todo:这里采用每隔1分钟拉取一次的频率
		tk := toolbox.NewTask("primary_data", "0 */1 * * * *", primary_data_task)
		//加入全局的计划任务列表
		toolbox.AddTask("primary_data", tk)
		//开始执行全局的任务
		toolbox.StartTask()
		//defer toolbox.StopTask()
	} else {
		beego.Info("组织结构同步关闭!")
	}
}

func primary_data_task() error {
	cronExpress := beego.AppConfig.String("rsyncFre")
	if cronExpress == "" {
		beego.Info("未填写同步时数，初始化为早上4点")
		cronExpress = "04"
	}
	beego.Info(">>>同步时间为>>>Cron Express:" + cronExpress)
	ins := getIfTrueTime(cronExpress)
	beego.Info("是否更新主数据:", ins)
	if ins == true {
		beego.Info("用户主数据+组织结构同步开始 >>>")
		//addrs2 :=
		GetGrignizes(3)
		//GetAllUser(3) 	todo:注意：因为中油瑞飞只能通过赋予用户访问平台权限方式登录行为科技平台.所以必然会新建账号进入,因此此处注释!!!
		//AlarmToPhone2(2)
		//AlarmToPhone3(3) //第二种方式发送短信
	}
	return nil
}

func getIfTrueTime(params string) (x bool) {
	flags := false
	t4 := strconv.Itoa(time.Now().Hour())
	if t4 == params {
		flags = true
	} else {
	}

	return flags
}

func GetGrignizes(timeout int) (rseponses string, err error) {
	addr := beego.AppConfig.String("allOrignizeRoute")

	if timeout < 0 {
		timeout = 5
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))) //设置发送接受数据超时
				return conn, nil
			},
		},
	}

	resp, err := http.NewRequest("GET", addr, nil)
	//beego.Info(beego.AppConfig.String("Authors"))
	resp.Header.Add("Authorization", beego.AppConfig.String("Authors"))

	if err != nil {
		fmt.Println("http.Get err=", err)
		return "n", nil
	}

	response, erra := client.Do(resp)

	if erra != nil {
		beego.Info("请求错误是:", erra)
		return addr, erra
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	insk := string(body)

	//---------------->>> Method 2
	m := make(map[string]interface{})
	json.Unmarshal([]byte(insk), &m)
	kpis := m["data"].(map[string]interface{})["groupInfos"]

	//beego.Info("\n",m,"\n")

	if m["code"] != "SUCCESS" {
		beego.Info(m["code"], "IP未认证！")
		return
	}

	root_docker := []string{}
	ban_docker := []map[string]string{}
	DOCKERS := []string{}
	o := controllers.O
	switch reflect.TypeOf(kpis).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(kpis)
		// 第一次遍历，获取家目录信息
		//beego.Info(s.Len())
		for i := 0; i < s.Len(); i++ {
			kpisB := s.Index(i).Interface().(map[string]interface{})
			//codes,names	_,_
			_, _, parentCodes, parentNames := kpisB["code"].(string), kpisB["name"].(string), kpisB["parentCode"].(string), kpisB["parentName"].(string)
			//todo:	对方测试数据混乱,可能names相同,但是codes不会相同

			existsCode_root := arrays.Contains(root_docker, parentCodes) //codes)//parentCodes)
			if parentNames == "长庆石化分公司" {
				//beego.Info(codes,names,parentCodes,parentNames)
				if existsCode_root == -1 {
					root_docker = append(root_docker, parentCodes) //codes) //parentCodes)
				}
			}
		}
		// 获取家目录的ori_id
		cnpc_ori_id := root_docker[0]
		//beego.Info("一级组织ID为:",root_docker)

		//o_id := uuid.Must(uuid.NewV4()).String()
		//o_id = strings.Replace(o_id, "-", "", -1)
		//todo:此处长庆石化总节点id赋值为“cqsh”
		o_id := "cqsh"

		//SQL := fmt.Sprintf("INSERT INTO ss_organize(id,name,pId,level,pname,ori_id) SELECT '%s','%s','%s','%s','%s','%s' FROM DUAL WHERE NOT EXISTS" +
		//	"(SELECT * FROM ss_organize WHERE ori_id = '%s');",
		//	o_id,"长庆石化分公司","0","0","长庆石化分公司",cnpc_ori_id,cnpc_ori_id)
		_, err := o.Raw("INSERT INTO ss_organize(id,name,pId,level,pname,ori_id) SELECT ?,?,?,?,?,? FROM DUAL WHERE NOT EXISTS(SELECT * FROM ss_organize WHERE ori_id = ?)", o_id, "长庆石化分公司", "0", "0", "长庆石化分公司", cnpc_ori_id, cnpc_ori_id).Exec()
		if err != nil {
			beego.Info("更新组织架构 ERROR!!!:")
		} else {
			//beego.Info("更新组织结构成功 >> 长庆石化")
			UpdateDefaltUser(cnpc_ori_id)
			//go UpdateDefaltUser(cnpc_ori_id)
		}
		// 第二次遍历，获取二级组织信息
		var mapsA []orm.Params
		num, erra := o.Raw("select * from ss_organize where ori_id=?;", cnpc_ori_id).Values(&mapsA)
		if erra != nil {
			beego.Info("查询组织结构id失败 >>> result: ", num)
		}
		if num < 1 {
			beego.Info("没有该组织信息!")
		}
		insk := mapsA[0]["id"]

		for i := 0; i < s.Len(); i++ {
			kpisB := s.Index(i).Interface().(map[string]interface{})
			//codes,names	_,_
			codes, names, parentCodes, parentNames := kpisB["code"].(string), kpisB["name"].(string), kpisB["parentCode"].(string), kpisB["parentName"].(string)

			//todo:对方测试数据混乱，可能name相同，但是id不会相同
			existsCode := arrays.Contains(DOCKERS, parentCodes) //codes)
			if parentNames != "长庆石化分公司" && parentNames != "根节点" {
				if existsCode == -1 {
					//beego.Info(codes,names,parentCodes,parentNames)//todo:新的组织
					DOCKERS = append(DOCKERS, parentCodes) //codes)
					o_id := uuid.Must(uuid.NewV4()).String()
					o_id = strings.Replace(o_id, "-", "", -1)

					//if parentCodes == "00000000000000000000000065125186"{
					//	parentNames = "长庆石化分公司调度中心"
					//}
					//SQL := fmt.Sprintf("INSERT INTO ss_organize(id,name,pId,level,pname,ori_id) SELECT '%s','%s','%s','%s','%s','%s' FROM DUAL WHERE NOT EXISTS" +
					//	"(SELECT * FROM ss_organize WHERE ori_id = '%s');",
					//	o_id,parentNames,insk,"1",parentNames,parentCodes,parentCodes)//names,codes,codes)
					res, err := o.Raw("INSERT INTO ss_organize(id,name,pId,level,pname,ori_id) SELECT ?,?,?,?,?,? FROM DUAL WHERE NOT EXISTS(SELECT * FROM ss_organize WHERE ori_id = ?)", o_id, parentNames, insk, "1", parentNames, parentCodes, parentCodes).Exec()
					if err != nil {
						beego.Info("更新组织架构 ERROR!!!:", res, len(DOCKERS))
					} else {
						//beego.Info("更新组织结构成功~",len(DOCKERS))
					}
				}
				//todo:此处为部门下具体班组的信息
				Ban_detail := map[string]string{}
				exists_Ban := arrays.Contains(ban_docker, codes)
				if exists_Ban == -1 {
					Ban_detail["names"] = names
					Ban_detail["codes"] = codes
					Ban_detail["parentcodes"] = parentCodes
					ban_docker = append(ban_docker, Ban_detail)
				}
			}
		}
		//beego.Info(fmt.Sprintf("二级组织ID:%s,总量:%d",DOCKERS,len(DOCKERS)))
		//beego.Info(fmt.Sprintf("所有班组信息:%s-班组总数:%d",ban_docker,len(ban_docker)))
		beego.Info("***开始将班组信息归类到各自的部门中***")
		for _, x := range ban_docker {
			//beego.Info(fmt.Sprintf("开始归置班组%s", x))
			//SQL := fmt.Sprintf("INSERT INTO ss_bu_ban(codes,names,parentcodes) SELECT '%s','%s','%s' FROM DUAL WHERE NOT EXISTS" +
			//	"(SELECT * FROM ss_bu_ban WHERE codes = '%s');",
			//	x["codes"],x["names"],x["parentcodes"],x["codes"])
			//beego.Info("SQL is:", SQL)
			resa, erra := o.Raw("INSERT INTO ss_bu_ban(codes,names,parentcodes) SELECT ?,?,? FROM DUAL WHERE NOT EXISTS(SELECT * FROM ss_bu_ban WHERE codes = ?)", x["codes"], x["names"], x["parentcodes"], x["codes"]).Exec()
			if erra != nil {
				beego.Info("更新班组信息 ERROR!!!:", resa, len(x))
			} else {
				//beego.Info("更新组织结构成功~",len(DOCKERS))
			}
		}
	}
	beego.Info("本轮组织结构更新完成~")
	return
}

func UpdateDefaltUser(ori_id string) {
	o := controllers.O
	var mapsA []orm.Params
	num, erra := o.Raw("select * from ss_organize where ori_id=?;", ori_id).Values(&mapsA)
	if erra != nil {
		beego.Info("查询组织结构id失败 >>> result: ", num)
	}
	uid := uuid.Must(uuid.NewV4()).String()
	uid = strings.Replace(uid, "-", "", -1) // 生成UUID
	//beego.Info("更新默认用户(admin,范以云,范艳妮)  ID >>>  ",mapsA[0]["id"])
	//SQL := fmt.Sprintf("update ss_user set organize_id = '%s',tele = '%s' where admin_user in ('fan','admin','T0295112');",mapsA[0]["id"],"")
	res, err := o.Raw("INSERT INTO ss_user(admin_password,admin_id,admin_user,admin_name,organize_id,update_time,type,jurisdiction,tele,allowd) SELECT ?,?,?,?,?,?,?,?,?,? FROM DUAL WHERE NOT EXISTS"+
		"(SELECT * FROM ss_user WHERE organize_id = ?);",
		"e10adc3949ba59abbe56e057f20f883e", uid, "T0295112", "范艳妮", mapsA[0]["id"], time.Now().Format("2006-01-02 15:04:05"), 2, "0,1,2,3,4,5,6", "", "1", mapsA[0]["id"]).Exec()
	if err != nil {
		beego.Info("更新默认用户信息失败!!!:", res)
	} else {
		//beego.Info("更新默认用户信息成功~",res)
	}
	return
}

func GetAllUser(timeout int) (rseponses string, err error) {
	addr := beego.AppConfig.String("allUserRoute")
	if timeout < 0 {
		timeout = 5
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))) //设置发送接受数据超时
				return conn, nil
			},
		},
	}

	resp, err := http.NewRequest("GET", addr, nil)
	resp.Header.Add("Authorization", beego.AppConfig.String("Authors"))

	if err != nil {
		fmt.Println("http.Get err=", err)
		return
	}

	response, erra := client.Do(resp)

	if erra != nil {
		beego.Info("请求错误是:", erra)
		return addr, erra
	}

	//不采用下面的方式
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		beego.Info("请求错误是:", err)
	}
	insk := string(body)

	//---------------->>> Method 2
	m := make(map[string]interface{})
	json.Unmarshal([]byte(insk), &m)
	kpis := m["data"].(map[string]interface{})["userInfos"]
	DOCKERS := []string{}
	o := controllers.O
	switch reflect.TypeOf(kpis).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(kpis)
		for i := 0; i < s.Len(); i++ {
			kpisB := s.Index(i).Interface().(map[string]interface{})
			userNames, displayNames, groupss, createDates := kpisB["userName"].(string), kpisB["displayName"].(string), kpisB["groups"].([]interface{}), kpisB["createDate"].(string)
			//todo:对方测试数据混乱,可能name相同，但是id不会相同
			existsCode := arrays.Contains(DOCKERS, displayNames)
			if existsCode == -1 {
				DOCKERS = append(DOCKERS, userNames) //todo:	username为用户账号，而并非中文名
				o_id := uuid.Must(uuid.NewV4()).String()
				o_id = strings.Replace(o_id, "-", "", -1)
				for _, groupA_ID := range groupss {
					var mapsA []orm.Params
					//beego.Info(groupA_ID,reflect.TypeOf(groupA_ID))
					num, erra := o.Raw("select * from ss_organize where ori_id=?;", groupA_ID).Values(&mapsA)
					if erra != nil {
						beego.Info(displayNames, "查询组织结构id失败 >>> result: ", num)
					} else {
						if len(mapsA) >= 1 {
							//beego.Info(displayNames, "查询组织结构id成功 >>> result: ", num, mapsA[0], reflect.TypeOf(mapsA))
							trsform_id := mapsA[0]["id"]
							// 人员爬上组织结构树 ~
							//SQL := fmt.Sprintf("INSERT INTO ss_user(admin_password,admin_id,admin_user,update_time,admin_name,organize_id,jurisdiction,type,tele,allowd) " +
							//	"SELECT '%s', '%s','%s','%s','%s','%s','0,1,2,3,4,5,6,7,8',%d,'%s','%s' FROM DUAL WHERE NOT EXISTS" +
							//	"(SELECT * FROM ss_user WHERE admin_user = '%s');",
							//	"",o_id,userNames,createDates,displayNames,trsform_id,2,"","1",userNames)
							SQL := "INSERT INTO ss_user(admin_password,admin_id,admin_user,update_time,admin_name,organize_id,jurisdiction,type,tele,allowd) " +
								"SELECT ?, ?,?,?,?,?,'0,1,2,3,4,5,6,7,8',?,?,? FROM DUAL WHERE NOT EXISTS" +
								"(SELECT * FROM ss_user WHERE admin_user = ?)"
							_, err := o.Raw(SQL, "", o_id, userNames, createDates, displayNames, trsform_id, 2, "", "1", userNames).Exec()
							if err != nil {
								//beego.Info(displayNames,"更新人员信息失败~ERROR ~", groupA_ID)
							} else {
								//beego.Info(displayNames,"更新人员信息成功 ~", groupA_ID ,"组织结构ID:" ,trsform_id)
							}
						} else {
							beego.Info(displayNames, " >>> 人员无组织机构 >>> Remote  group ID is :", groupA_ID)
						}
					}
				}
			}
		}
	}
	beego.Info("本轮用户主数据更新完成 ~")
	return
}

//func AlarmToPhone2(timeout int) (rseponses string, err error) {
//	if timeout < 0 {
//		timeout = 5
//	}
//
//	client := &http.Client{
//		Transport: &http.Transport{
//			Dial: func(netw, addr string) (net.Conn, error) {
//				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //设置建立连接超时
//				if err != nil {
//					return nil, err
//				}
//				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))) //设置发送接受数据超时
//				return conn, nil
//			},
//		},
//	}
//	addr := "http://10.189.0.87:8080/api/send?token=6EBC3F3A3B8D4B2B89BABFD68C1E50C7"
//	//"http://10.189.0.87:8080/api/send?token=6EBC3F3A3B8D4B2B89BABFD68C1E50C7&mobile=13141005621&content=
//	// '【长庆石化公司视频行为分析预警平台】发现一条火焰检测预警，请您及时处理！<br/> 预警摄像头：rtsp://admin:q1w2e3r4:192.168.10.10:554 <br/>
//	// 预警时间：2020-08-08 12:11:56 <br/> 预警描述：此处产生火势，请前往扑灭！！！'"
//	datas := url.Values{}
//	datas.Set("mobile","13141005621")
//	datas.Set("content","11111")
//	b := bytes.NewBuffer([]byte(datas.Encode()))
//	beego.Info()
//	resp, err := http.NewRequest("GET", addr, b)
//
//	if err != nil {
//		fmt.Println("http.Get err=",err)
//		return "n",nil
//	}
//
//	response, nil := client.Do(resp)
//
//	defer response.Body.Close()
//
//	body, err := ioutil.ReadAll(response.Body)
//
//	if err != nil {
//		beego.Info("请求错误是:",err)
//	}
//
//	insk := string(body)
//	m := make(map[string]interface{})
//	json.Unmarshal([]byte(insk),&m)
//	msg := m["msg"]
//	beego.Info(msg,"&&&&&&&&&&&&&&&&&&&&&&")
//	return "y",nil
//}
//
//func AlarmToPhone3(timeout int) (results string){
//	//sign_token := fmt.Sprintf("http://rfiamtest.cnpc/ngiam-rst/oauth2/userinfo?appcode=cqshfxyj&secret=riqivvy471bsrdp8zvs3v5tsrhmd13sg&token=%s",
//	//	tokens)
//	sign_token := "http://10.189.0.87:8080/api/send?token=6EBC3F3A3B8D4B2B89BABFD68C1E50C7" +
//		"&mobile=13141005621" +
//		"&content=长庆石化公司视频行为分析预警平台发现一条火焰检测预警,请您及时处理预警摄像头:火焰检测camera001 预警时间:2020-08-11 12:11:56 预警描述:此处产生火势,请前往扑灭!!!"
//
//	if timeout < 0 {
//		timeout = 5
//	}
//
//	client := &http.Client{
//		Transport: &http.Transport{
//			Dial: func(netw, addr string) (net.Conn, error) {
//				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //设置建立连接超时
//				if err != nil {
//					beego.Info(err)
//				}
//				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))) //设置发送接受数据超时
//				return conn,nil
//			},
//		},
//	}
//
//	resp, err := http.NewRequest("GET", sign_token, nil)
//	resp.Header.Add("Authorization", "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCzW0spP6KEa15bqUsvuOIjxojmFox54QgziytBV" +
//		"b4Utt24Ix55TYq1tSt8H0XspTLCj1cgYgRRtRCaotbumyxXFHlu/eDrWfvmkO71nbFj9kGWQgPbQJ+AL1MbiN6GazLRteTIlCOJku1i9RYkKYuOKetjPJK2YqtcWKAeuxTVJQIDAQAB")
//
//	if err != nil {
//		fmt.Println("http.Get err=",err)
//		return
//	}
//
//	response, nil := client.Do(resp)
//
//	defer response.Body.Close()
//
//	body, err := ioutil.ReadAll(response.Body)
//
//	if err != nil {
//		beego.Info("请求错误是:",err)
//	}
//	insk := string(body)
//	beego.Info(insk)
//	//{"access_token":"XXXXX","exprise_in":"60"}
//
//	return insk
//}

func AlarmToPhone4() {
	DX_ADD := beego.AppConfig.String("DX_ADD")
	u, _ := url.Parse(DX_ADD)
	q := u.Query()
	q.Set("mobile", "13141005621")
	q.Set("content", "长庆石化公司视频行为分析预警平台发现一条火焰检测预警,请您及时处理!!!! 预警摄像头:火焰检测camera001 预警时间:2020-08-11 12:11:56 预警描述:此处产生火势,请前往扑灭!!!")
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

	beego.Info(result)
}
