package signs

import (
	"Artifice_V2.0.0/controllers"
	"Artifice_V2.0.0/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/satori/go.uuid"
	"github.com/wxnacy/wgo/arrays"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type HlsPhoneController struct {
	beego.Controller
}

type OGR struct {
	Code    string                 `json:"code"`
	OrgList map[string]interface{} `json:"orgList"`
}

type MNN struct {
	DisplayName string `json:"displayName"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Username    string `json:"username"`
}

func (this *HlsPhoneController) Post() {
	apps := make(map[string]interface{})
	json.Unmarshal([]byte(string((this.Ctx.Input.RequestBody))), &apps)
	//beego.Info("+++++++++++\n",apps,"\n++++++++++")

	//beego.Info(reflect.TypeOf(apps["enabled"]),apps["enabled"])
	enable := apps["enabled"]
	//beego.Info(reflect.TypeOf(apps["status"]),apps["status"])
	status := apps["status"]
	//beego.Info(reflect.TypeOf(apps["accountName"]),apps["accountName"])
	accountName := apps["accountName"].(string)
	beego.Info("操作账号:", accountName, reflect.TypeOf(accountName))

	//解析账号映射用户信息
	mnn := apps["mappingAttr"]
	jsonStr, err := json.Marshal(mnn)
	if err != nil {
		beego.Info("MapToJsonDemo err: ", err)
	}
	kkk := string(jsonStr)
	//fmt.Println(kkk)

	var someM MNN
	if err := json.Unmarshal([]byte(kkk), &someM); err == nil {
		man_names := someM.DisplayName
		phone := someM.Mobile
		beego.Info("账号用户:", man_names)
		beego.Info("手机号码:", phone)
	} else {
		beego.Info(err)
	}
	man_names := someM.DisplayName
	phone := someM.Mobile
	//解析账号组织结构ID
	o := controllers.O
	kpis := apps["orgList"]
	var DOCKER_TRID []string

	//adding a new docker,include many someOne.Code and it not belong to
	special_docker := []string{
		"00000000000000000000000000208423", //长庆石化分公司运行一部
		"00000000000000000000000000208424", //长庆石化分公司运行二部
		"00000000000000000000000000208425", //长庆石化分公司运行三部
		"00000000000000000000000065125637", //长庆石化分公司运行四部
		"00000000000000000000000000208415", //长庆石化分公司质量检验部（检测中心）
		"00000000000000000000000065125186", //长庆石化分公司机关附属机构（调度中心）
		"00000000000000000000000000208426", //长庆石化分公司油品运行部
		"00000000000000000000000000208414", //长庆石化分公司运行保障部
	}
	//beego.Info("此处没有领导视角的班组ID们是：",special_docker)

	switch reflect.TypeOf(kpis).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(kpis)
		for i := 0; i < s.Len(); i++ {
			ppp := s.Index(i).Interface()
			var someOne OGR
			if err := json.Unmarshal([]byte(ppp.(string)), &someOne); err == nil {
				uuids := someOne.Code
				beego.Info("所属班组-组织结构编号:", uuids, reflect.TypeOf(uuids))

				var mapsA []orm.Params
				//num, erra := o.Raw("select * from ss_organize where ori_id=?;",uuids).Values(&mapsA)
				num, erra := o.Raw("select a.id,a.ori_id,b.codes,b.parentcodes from ss_organize a left join "+
					"ss_bu_ban b on a.ori_id=b.parentcodes where a.ori_id != '' and a.id != 'cqsh' and b.codes=?;", uuids).Values(&mapsA)
				if erra != nil {
					beego.Info("查询部门id失败 >>> result: ", num)
				} else {
					if len(mapsA) >= 1 {
						beego.Info("查询部门uuid成功 >>> result: ", num, mapsA[0]["ori_id"].(string), reflect.TypeOf(mapsA))

						existsCode_root := arrays.Contains(special_docker, mapsA[0]["ori_id"].(string))
						if existsCode_root == -1 {
							beego.Info("其他部门~")
							DOCKER_TRID = append(DOCKER_TRID, "cqsh")
						} else {
							DOCKER_TRID = append(DOCKER_TRID, mapsA[0]["id"].(string))
						}
						beego.Info("DOCKER_TRID:", DOCKER_TRID)
					} else {
						beego.Info("部门ori_id未找到:", mapsA)
						this.Data["json"] = map[string]interface{}{"code": "error", "message": "组织结构ori_id未找到"}
						this.ServeJSON()
						return
					}
				}
			} else {
				beego.Info("获取部门uuid失败:", err)
			}
		}
	}

	mini_dock := []string{}
	if enable == true {
		mini_dock = append(mini_dock, "1")
	} else {
		mini_dock = append(mini_dock, "0")
	}

	if len(accountName) >= 1 {
		//此处添加账号信息判断

		if status == "CREATE" {
			//开始创建用户

			var userMaps []*models.User
			qs := o.QueryTable(new(models.User))
			_, err := qs.Filter("admin_user", accountName).All(&userMaps)
			if err != nil {
				log.Println(err)
				this.Data["json"] = map[string]interface{}{"code": "error", "message": "验证账号信息异常"}
				this.ServeJSON()
				return
			}
			for _, u := range userMaps {
				if u.Acount == accountName {
					this.Data["json"] = map[string]interface{}{"code": "error", "message": "添加员工账号已存在"}
					this.ServeJSON()
					return
				}
			}

			uid := uuid.Must(uuid.NewV4()).String()
			uid = strings.Replace(uid, "-", "", -1) // 生成UUID

			_, errx := o.Raw("insert into ss_user(admin_password,admin_id,admin_user,admin_name,organize_id,update_time,type,jurisdiction,tele,allowd) values(?,?,?,?,?,?,?,?,?,?)",
				"e10adc3949ba59abbe56e057f20f883e", uid, accountName, man_names, DOCKER_TRID[0], time.Now().Format("2006-01-02 15:04:05"), 2, "0,1,2,3,4,5,6", phone, mini_dock[0]).Exec()
			if errx != nil {
				beego.Info("添加账号成功", errx)
				this.Data["json"] = map[string]interface{}{"code": "error", "message": "添加账号异常"}
				this.ServeJSON()
				return
			} else {
				beego.Info("添加账号成功")
			}
		} else if status == "DELETE" {
			//开始删除用户
			o := controllers.O
			_, err = o.Raw("DELETE FROM ss_user where admin_user=?", accountName).Exec()
			if err != nil {
				log.Println(err)
				this.Data["json"] = map[string]interface{}{"code": "error", "message": "删除账号异常"}
				this.ServeJSON()
				return
			} else {
				beego.Info("删除账号成功")
			}

		} else if status == "UPSERT" {
			//开始修改账号
			beego.Info("开始修改账号")
			o := controllers.O
			_, err = o.Raw("UPDATE ss_user set admin_name=?,update_time=?,type=?,jurisdiction=?,tele=?,allowd=? where admin_user=? and organize_id=?;",
				man_names, time.Now().Format("2006-01-02 15:04:05"), 2, "0,1,2,3,4,5,6", phone, mini_dock[0], accountName, DOCKER_TRID[0]).Exec()
			if err != nil {
				log.Println(err)
				this.Data["json"] = map[string]interface{}{"code": "error", "message": "修改(编辑)账号异常"}
				this.ServeJSON()
				return
			} else {
				beego.Info("修改(编辑)账号成功")
			}
		} else {
			//首先判断status参数的值。若发现无此参数，则去判断enable参数的值
			beego.Info(">>>开始判断账号是否被禁用>>>", enable, reflect.TypeOf(enable))
			//if enable == true {
			beego.Info("设置账号是否禁用开关", mini_dock[0])
			o := controllers.O
			_, errr := o.Raw("UPDATE ss_user set allowd=? where admin_user=?",
				mini_dock[0], accountName).Exec()
			if errr != nil {
				beego.Info("设置账号是否禁用状态失败!", errr)
				this.Data["json"] = map[string]interface{}{"code": "error", "message": "设置账号可使用异常"}
				this.ServeJSON()
				return
			} else {
				beego.Info("设置账号是否禁用状态成功!")
			}
		}
		beego.Info("***HAPPY TO END***")
		this.Data["json"] = map[string]interface{}{"code": "success"}
		this.ServeJSON()
		return
	} else {
		beego.Info("没有获取到账号信息(账号长度为零)")
		this.Data["json"] = map[string]interface{}{"code": "error", "message": "没有获取到账号信息(账号长度为零)"}
		this.ServeJSON()
		return
	}
	beego.Info("***HAPPY TO END***")
	this.Data["json"] = map[string]interface{}{"code": "success"}
	this.ServeJSON()
	return
}

func (this *HlsPhoneController) SSO() {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(GetToken(this.GetString("code"), 3)), &m)
	accessToken := m["accessToken"].(string)
	//-------------------------------------------------------------------------
	n := make(map[string]interface{})
	json.Unmarshal([]byte(GetLogInfo(accessToken, 3)), &n)
	LoginInfo := n["accountName"].(string)
	beego.Info("the LoginInfo is:", LoginInfo)
	o := controllers.O
	var mapsB []orm.Params
	num, erra := o.Raw("select * from ss_user where admin_user=?;", LoginInfo).Values(&mapsB)
	if erra != nil {
		beego.Info("查询是否用户uuid >>> error is: ", erra)
	}
	admin_ids := mapsB[0]["admin_id"]
	beego.Info("该用户唯一ID是:", admin_ids)

	this.Data["SystemTitle"] = controllers.SelSystemTitleUtil()
	this.Data["bmk"] = "homemenu"
	this.Data["Time"] = controllers.GetTimeSE()

	uid := uuid.Must(uuid.NewV4()).String()
	uid = strings.Replace(uid, "-", "", -1) // 生成SessionID
	var sessions []orm.Params
	num, errA := o.Raw("select * from ss_sessions where userid=?", admin_ids).Values(&sessions)
	if errA != nil {
		this.Data["json"] = map[string]interface{}{"code": -1, "msg": "登录失败"}
	} else {
		if num == 0 {
			//beego.Info("新增session")
			resa, erra := o.Raw("insert into ss_sessions (id, userid, sessionid) values (?,?,?)",
				strings.Replace(uid, "-", "", -1), admin_ids, admin_ids).Exec()
			if erra != nil {
				beego.Info("新增session失败", erra, resa)
			}
			beego.Info("新增session成功", admin_ids, uid)
			this.Ctx.SetCookie("user_session", admin_ids.(string), "/")
			this.Ctx.SetCookie("user_id", admin_ids.(string), "/")
			this.Ctx.SetCookie("admin_type", "2", "/")

			//初次单点登录
			var user models.User
			admin_idsA := admin_ids.(string)
			err := o.QueryTable(new(models.User)).Filter("admin_id", admin_idsA).RelatedSel().One(&user)
			if err != nil {
				beego.Info(">>>>>", err)
				this.TplName = "login.html"
			} else {
				this.Data["admin_name"] = user.Name
				this.Data["admin_user"] = user.Acount
				this.Data["organize_name"] = user.Organize.Name
				this.Data["jurisdiction"] = user.Jurisdiction
				//this.TplName = "platform/index.html"
				this.Ctx.Redirect(302, "/platform/homepage")
			}

		} else {
			res, errb := o.Raw("update ss_sessions set sessionid=?  where userid=?", admin_ids, admin_ids).Exec()
			if errb != nil {
				beego.Info("更新session失败", errb, res)
			}
			beego.Info("更新session成功", admin_ids, uid)
			this.Ctx.SetCookie("user_session", admin_ids.(string), "/")
			this.Ctx.SetCookie("user_id", admin_ids.(string), "/")
			this.Ctx.SetCookie("admin_type", "2", "/")

			var user models.User
			admin_idsA := admin_ids.(string)
			err := o.QueryTable(new(models.User)).Filter("admin_id", admin_idsA).RelatedSel().One(&user)
			if err != nil {
				beego.Info(">>>>>", err)
				this.TplName = "login.html"
			} else {
				this.Data["admin_name"] = user.Name
				this.Data["admin_user"] = user.Acount
				this.Data["organize_name"] = user.Organize.Name
				this.Data["jurisdiction"] = user.Jurisdiction
				//this.TplName = "platform/index.html"
				this.Ctx.Redirect(302, "/platform/homepage")
			}

		}
	}
}

func GetToken(codes string, timeout int) (results string) {
	beego.Info("the ask code is:", codes)
	sign_code := fmt.Sprintf(beego.AppConfig.String("Getstoken"),
		codes)

	if timeout < 0 {
		timeout = 5
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //设置建立连接超时
				if err != nil {
					beego.Info(err)
				}
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))) //设置发送接受数据超时
				return conn, nil
			},
		},
	}

	resp, err := http.NewRequest("GET", sign_code, nil)
	resp.Header.Add("Authorization", beego.AppConfig.String("Authors"))

	if err != nil {
		fmt.Println("http.Get err=", err)
		return
	}

	response, nil := client.Do(resp)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		beego.Info("请求错误是:", err)
	}
	insk := string(body)
	//{"access_token":"XXXXX","exprise_in":"60"}

	return insk
}

func GetLogInfo(tokens string, timeout int) (results string) {
	beego.Info("the ask token is:", tokens)
	sign_token := fmt.Sprintf(beego.AppConfig.String("UserInfos"),
		tokens)

	if timeout < 0 {
		timeout = 5
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //设置建立连接超时
				if err != nil {
					beego.Info(err)
				}
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))) //设置发送接受数据超时
				return conn, nil
			},
		},
	}

	resp, err := http.NewRequest("GET", sign_token, nil)
	resp.Header.Add("Authorization", beego.AppConfig.String("Authors"))

	if err != nil {
		fmt.Println("http.Get err=", err)
		return
	}

	response, nil := client.Do(resp)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		beego.Info("请求错误是:", err)
	}
	insk := string(body)
	//{"access_token":"XXXXX","exprise_in":"60"}

	return insk
}
