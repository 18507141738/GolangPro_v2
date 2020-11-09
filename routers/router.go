package routers

import (
	"Artifice_V2.0.0/conf"
	"Artifice_V2.0.0/controllers"
	"Artifice_V2.0.0/controllers/config"
	"Artifice_V2.0.0/controllers/filter"
	"Artifice_V2.0.0/controllers/handle_comment"
	"Artifice_V2.0.0/models"
	"Artifice_V2.0.0/signs"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strings"
)

func init() {

	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)

	//登录页面
	beego.Router("/login", &controllers.LoginPage{})
	//登录
	beego.Router("/loginCheck", &controllers.LoginCheckHandler{}, "get:Get;post:Post")
	//注销
	beego.Router("/logout", &controllers.LogoutUserHandler{}, "get:Get;post:Post")
	//展示首页
	beego.Router("/platform/homepage", &controllers.PlatHomePage{})
	//展示首页数据
	beego.Router("/platform/homedata", &controllers.PlatHomeController{}, "get:Get;post:Post")
	//算法统计页面
	beego.Router("/platform/funcpage", &controllers.PlatFuncPage{})
	beego.Router("/platform/funcpageie", &controllers.PlatFuncPageIE{})
	//展示统计页组织查询
	beego.Router("/platform/selorg", &controllers.PlatOrganizeBelongUser{}, "get:Get;post:Post")
	//展示统计页面图表数据查询
	beego.Router("/platform/funcChartLine", &controllers.PlatFuncLineController{}, "get:Get;post:Post")
	//展示统计页面告警数据查询
	beego.Router("/platform/funcAlarmList", &controllers.PlatFuncAlarmsController{}, "get:Get;post:Post")
	//展示个人中心
	beego.Router("/platform/centerpage", &controllers.PlatCenterPage{})
	//修改密码
	beego.Router("/platform/updatePass", &controllers.UpdateUserPasswordHandler{}, "get:Get;post:Post")
	//成员管理页面
	beego.Router("/platform/centerusermanagerpage", &controllers.PlatformUserMangerPageController{})
	//成员管理-组织架构查询
	beego.Router("/platform/selorgforadminid", &controllers.PlatSelOrgainze{}, "get:Get;post:Post")
	//成员管理-组织架构保存
	beego.Router("/platform/saveOrganize", &controllers.SaveOrganizeHandler{}, "get:Get;post:Post")
	//成员管理-组织架构修改
	beego.Router("/platform/editOrganize", &controllers.EditOrganizeHandler{}, "get:Get;post:Post")
	//成员管理-组织架构删除
	beego.Router("/platform/delOrganize", &controllers.DelOrganizeHandler{}, "get:Get;post:Post")
	//成员管理-成员查询
	beego.Router("/platform/selUserForOrg", &controllers.SelMemberByOrgHandler2{}, "get:Get;post:Post")
	//成员管理-成员添加
	beego.Router("/platform/saveOrgUser", &controllers.OrgAddUserHandler{}, "get:Get;post:Post")
	//成员管理-成员修改
	beego.Router("/platform/editOrgUser", &controllers.OrgEditUserHandler{}, "get:Get;post:Post")
	//成员管理-成员删除
	beego.Router("/platform/delOrgUser", &controllers.OrgDelUserHandler{}, "get:Get;post:Post")
	// 设备管理页面
	beego.Router("/platform/centerhostmanagerpage", &controllers.PlatformHostManagePage{})
	// 查询分析主机列表
	beego.Router("/platform/hostlist", &controllers.PlatformSelHostPageHandler{}, "get:Get;post:Post")
	// 新增、修改分析主机
	beego.Router("/platform/hostManage", &controllers.HostManageHandler{}, "get:Get;post:Post")
	// 批量删除分析主机
	beego.Router("/platform/delBatchHost", &controllers.DelBatchHostHandler{}, "get:Get;post:Post")
	// 判断主机是否超时
	//beego.Router("platform/checkhoststatus", &controllers.PlatformCheckHostOutHandler{}, "get:Get;post:Post")
	// 摄像机管理页面
	beego.Router("/platform/centercameramanagerpage", &controllers.PlatformCameraManagerPageController{})
	beego.Router("/platform/centercameramanageriepage", &controllers.PlatformCameraManagerIEPageController{})
	// 查询摄像机列表
	beego.Router("/platform/selCameraList", &controllers.SelCameraListHandler{}, "get:Get;post:Post")
	// 区域查询
	beego.Router("/platform/selPlace", &controllers.SelPlaceHandler{}, "get:Get;post:Post")
	// 新增、编辑摄像机
	beego.Router("/platform/saveCamera", &controllers.SaveCameraHandler{}, "get:Get;post:Post")
	// 启用、停用摄像机
	beego.Router("/platform/setCameraStatus", &controllers.SetCameraStatusHandler{}, "get:Get;post:Post")
	// 删除摄像机
	beego.Router("/platform/delCameraData", &controllers.DelCameraDataHandler{}, "get:Get;post:Post")
	// 查询指定类型下的摄像机列表
	beego.Router("/platform/selCameraForFuncType", &controllers.SelCameraForFuncType{}, "get:Get;post:Post")

	// 周界
	beego.Router("/platform/boundaryConfigPage", &controllers.BoundaryConfigPageHandler{})
	// 着装
	beego.Router("/platform/clothsConfigPage", &controllers.ClothsConfigPageHandler{})
	// 离岗页面
	beego.Router("/platform/sleep_count", &controllers.PlatformSleepCountPageHandler{})
	// 睡岗页面
	beego.Router("/platform/leave_count", &controllers.PlatformLeaveCountPageHandler{})
	// 烟雾页面
	beego.Router("/platform/smokepage", &controllers.PlatformSmokepagetPageHandler{})
	// 火焰页面
	beego.Router("/platform/fireworkpage", &controllers.PlatformFireworkpagePageHandler{})
	// 泄露
	beego.Router("/platform/leakage", &controllers.PlatformLeakagePageHandler{})
	// 配置详情页面
	beego.Router("/getAlgorithmList/", &controllers.GetAlgorithmListHandler{})
	// 保存配置
	beego.Router("/platform/add_grp", &controllers.AddGrpHandler{}, "get:Get;post:Post")
	//删除摄像机配置 算法端负样本
	beego.Router("/platform/delCameraFuncNegative", &controllers.DelCameraFuncNegative{})

	// ocx
	// 周界
	beego.Router("/platform/ocx/boundaryConfigPage", &controllers.BoundaryConfigOCXPageHandler{})
	// 着装
	beego.Router("/platform/ocx/clothsConfigPage", &controllers.ClothsConfigOCXPageHandler{})
	// 离岗页面
	beego.Router("/platform/ocx/sleep_count", &controllers.PlatformSleepCountOCXPageHandler{})
	// 睡岗页面
	beego.Router("/platform/ocx/leave_count", &controllers.PlatformLeaveCountOCXPageHandler{})
	// 烟雾页面
	beego.Router("/platform/ocx/smokepage", &controllers.PlatformSmokepagetOCXPageHandler{})
	// 火焰页面
	beego.Router("/platform/ocx/fireworkpage", &controllers.PlatformFireworkpageOCXPageHandler{})
	// 泄露
	beego.Router("/platform/ocx/leakage", &controllers.PlatformLeakageOCXPageHandler{})

	/*
		rtmp流配置页面
	*/
	//周界
	beego.Router("/platform/rtmp/boundaryConfigPage", &controllers.BoundaryConfRTMPPageHandler{})
	//着装
	beego.Router("/platform/rtmp/clothsConfigPage", &controllers.ClothsConfigRTMPPageHandler{})
	// 离岗页面
	beego.Router("/platform/rtmp/sleep_count", &controllers.PlatformSleepCountRTMPPageHandler{})
	// 睡岗页面
	beego.Router("/platform/rtmp/leave_count", &controllers.PlatformLeaveCountRTMPPageHandler{})
	// 烟雾页面
	beego.Router("/platform/rtmp/smokepage", &controllers.PlatformSmokepagetRTMPPageHandler{})
	// 火焰页面
	beego.Router("/platform/rtmp/fireworkpage", &controllers.PlatformFireworkpageRTMPPageHandler{})
	// 泄露
	beego.Router("/platform/rtmp/leakage", &controllers.PlatformLeakageRTMPPageHandler{})

	// 实时获取新告警
	beego.Router("/platform/selTaskHandler", &controllers.SelNewAlarmHandler{}, "get:Get;post:Post")
	// 告警处理
	beego.Router("/platform/alarmhandler", &controllers.AlarmHandle{}, "get:Get;post:Post")

	//分布式
	beego.Router("/distribute", &controllers.ClientDistributeCamera2_Handler{})
	// 报警信息保存
	beego.Router("/distribute/alarm_save", &controllers.ClientAlarmSaveController{}, "get:Get;post:Post")
	// 查询分配信息
	//beego.Router("/distributes", &controllers.GetOnfuncCaemra{})
	// 设置负样本
	beego.Router("/negative", &controllers.ClientCameraFuncNegative{}, "get:Get;post:Post")

	//投屏登录页
	beego.Router("/throwscreen/login", &controllers.TSLoginPage{})
	//投屏页登录
	beego.Router("/throwscreen/loginCheck", &controllers.TSLoginCheck{}, "get:Get;post:Post")
	//投屏页
	beego.Router("/throwscreen/homepage", &controllers.TSHomePage{})
	//投屏页算法统计
	beego.Router("/throwscreen/chartfuncdata", &controllers.TSFuncCountController{}, "get:Get;post:Post")
	//投屏摄像头
	beego.Router("/throwscreen/camera", &controllers.TSCamearController{}, "get:Get;post:Post")
	//投屏告警列表
	beego.Router("/throwscreen/alarmUntreated", &controllers.TSAlarmUntreatedController{}, "get:Get;post:Post")
	//投屏总统计
	beego.Router("/throwscreen/allFunc", &controllers.TSAllFuncAlarmController{}, "get:Get;post:Post")

	//设置页面登录
	// 修改密码
	beego.Router("/resetPwd", &config.UpdateUserPasswordHandler{}, "get:Get;post:Post")
	// 退出登录
	beego.Router("/logout", &config.LogoutUserHandler{}, "get:Get;post:Post")
	// 系统设置页面
	beego.Router("/config/systemPage", &config.SysPageHandler{})
	// 查询系统信息
	beego.Router("/config/selSysinfo", &config.SelSysTemInfoHandler{}, "get:Get;post:Post")
	// 修改系统设置
	beego.Router("/config/updateSys", &config.UpdateConfigHandler{}, "get:Get;post:Post")
	// 区域页面
	beego.Router("/config/loactionPage", &config.LocationPageHandler{})
	// 区域列表
	beego.Router("/config/sellocations", &config.SelLocationListHandler{}, "get:Get;post:Post")
	// 区域查询组织架构
	beego.Router("/config/selOrganizeList", &config.SelOrganizeListHandler{}, "get:Get;post:Post")
	// 添加、修改区域
	beego.Router("/config/locationAdd", &config.LocationAddController{}, "get:Get;post:Post")
	// 删除区域
	beego.Router("/config/locationDel", &config.DelLocationDataHandler{}, "get:Get;post:Post")
	// 组织架构页面
	beego.Router("/config/organizePage", &config.OrganizePageHandler{})
	// 组织架构树状数据
	beego.Router("/config/selOrgTree", &config.SelOrganizeTreeDataHandler{}, "get:Get;post:Post")
	// 根据组织架构查询用户
	beego.Router("/config/selUserForOrgId", &config.SelUserForOrgIdHandler{}, "get:Get;post:Post")
	// 保存修改组织架构
	beego.Router("/config/saveorgdata", &config.SaveOrganizeDataHandler{}, "get:Get;post:Post")
	// 删除组织
	beego.Router("/config/delOrgData", &config.DelOrganizeDataHandler{}, "get:Get;post:Post")
	// 组织架构新增用户
	beego.Router("/config/saveUser", &config.SaveUserOrganizeHandler{}, "get:Get;post:Post")
	// 组织架构修改用户
	beego.Router("/config/editUser", &config.EditUserOrganizeHandler{}, "get:Get;post:Post")
	// 删除用户
	beego.Router("/config/delUser", &config.DelUserOrganizeHandler{}, "get:Get;post:Post")
	//告警文案设置
	beego.Router("/config/alarmSetPage", &config.AlarmSetPage{})
	//告警文案添加
	beego.Router("/config/alarmDetailSel", &config.SelAlarmDetailController{}, "get:Get;post:Post")
	//告警文案添加
	beego.Router("/config/alarmDetailAdd", &config.AddAlarmDetailController{}, "get:Get;post:Post")
	//告警文案修改
	beego.Router("/config/alarmDetailEdit", &config.EditAlarmDetailController{}, "get:Get;post:Post")
	//告警文案删除
	beego.Router("/config/alarmDetailDel", &config.DelAlarmDetailController{}, "get:Get;post:Post")

	//批注处理页面登录
	beego.Router("/hc/login", &handle_comment.HCLoginPage{})
	beego.Router("/hc/loginCheck", &handle_comment.HCLoginCheck{}, "get:Get;post:Post")
	// 今日报警页面
	beego.Router("/hc/index", &handle_comment.EventPageHandler{})
	// 今日报警列表
	beego.Router("/hc/selAlarmList", &handle_comment.SelHcAlarmListHandler{}, "get:Get;post:Post")
	// 单条报警批注列表
	beego.Router("/hc/selAlarmNotations", &handle_comment.SelAlarmNotationListHandler{}, "get:Get;post:Post")
	// 添加批注
	beego.Router("/hc/saveNotation", &handle_comment.SaveAlarmNotationHandler{}, "get:Get;post:Post")

	/**车间告警页面**/
	beego.Router("/hc/workshop_alarm", &handle_comment.WorkshopAlarmController2{})
	/**车间告警页面**/
	beego.Router("/hc/workshop_unalarm", &handle_comment.WorkshopUnAlarmController2{})
	/**车间报警数量查询**/
	beego.Router("/hc/workshop_unalarmList", &handle_comment.WorkshopUnAlarmListController2{}, "get:Get;post:Post")
	/**车间处理告警消息**/
	beego.Router("/hc/workshop_dealalarm", &handle_comment.WorkshopDealAlarmController{}, "get:Get;post:Post")
	/**车间告警列表**/
	beego.Router("/hc/workshop_alarmList", &handle_comment.WorkshopAlarmListController2{}, "get:Get;post:Post")
	/**车间批注列表页面**/
	beego.Router("/hc/workshop_notationPage", &handle_comment.WorkshopNotationController{})
	/*本部门及下属部门批注列表*/
	beego.Router("/hc/query_notationList", &handle_comment.WorkshopNotationListController2{}, "get:Get;post:Post")
	/*单条告警批注列表*/
	beego.Router("/hc/alarm_notationlist", &handle_comment.WorkshopAlarmNotationListController{}, "get:Get;post:Post")
	//-------------------------------------------------------中油瑞飞sso模块-------------------------------------------------------
	beego.Router("/distribute/sso/code_recv", &signs.HlsPhoneController{}, "get:SSO;post:Post")
	beego.Router("/distribute/iam/user_update", &signs.HlsPhoneController{}, "get:Get;post:Post")
	//----------------------------------------------------------------------------------------------------------------------------

	// 过滤功能
	// 过滤登录页
	beego.Router("/filter/login", &filter.FilterLoginPage{})
	// 过滤登录
	beego.Router("/filter/loginCheck", &filter.FilterLoginCheckHandler{}, "post:Post")
	// 过滤设置页面
	beego.Router("/filter/station", &filter.FilterStationPage{})
	// 过滤设置
	beego.Router("/filter/setStatus", &filter.FilterSetStatus{}, "post:Post")
	// 过滤用户管理页面
	beego.Router("/filter/user", &filter.FilterUserPage{})
	// 过滤组织树查询
	beego.Router("/filter/orgtree", &filter.FilterSelOrgHandler{}, "post:Post")
	// 获取过滤组织下的用户
	beego.Router("/filter/seluser", &filter.FilterSelUserHandler{}, "post:Post")
	// 添加过滤组织账户
	beego.Router("/filter/adduser", &filter.FilterAddUserHandler{}, "post:Post")
	// 更新过滤组织账户
	beego.Router("/filter/edituser", &filter.FilterUpdateUserHandler{}, "post:Post")
	// 删除过滤组织账户
	beego.Router("/filter/deluser", &filter.FilterDelUserHandler{}, "post:Post")
	// 过滤首页
	beego.Router("/filter/homepage", &filter.FilterHomePage{})
	// 功能页
	beego.Router("/filter/funcpage", &filter.FilterFuncPage{})
	beego.Router("/filter/funcpageie", &filter.FilterFuncPageIE{})
	// 过滤首页统计
	beego.Router("/filter/homedata", &filter.FilterHomeDataHandler{}, "post:Post")
	// 过滤功能页组织
	beego.Router("/filter/selorg", &filter.FilterOrgainzeByUserHandler{}, "post:Post")
	// 过滤功能折线统计
	beego.Router("/filter/funcChartLine", &filter.FilterFuncLineHandler{}, "post:Post")
	// 过滤告警数据查询
	beego.Router("/filter/funcAlarmList", &filter.FilterFuncAlarmsHandler{}, "post:Post")
	// 查询指定类型下的摄像机列表
	beego.Router("/filter/selCameraForFuncType", &filter.SelCameraForFuncType{}, "post:Post")
	// 实时获取新告警
	beego.Router("/filter/selTaskHandler", &filter.FilterSelNelAlarmHandler{}, "post:Post")
	// 过滤告警插入
	beego.Router("/filter/insertAlarm", &filter.FilterInsertAlarmHandler{}, "post:Post")
	// 过滤告警误报
	beego.Router("/filter/misinformationHandler", &filter.FilterMisinformationHandler{}, "post:Post")
}

func FilterUser(ctx *context.Context) {
	if ctx.Request.RequestURI == "/login" || ctx.Request.RequestURI == "/loginCheck" || ctx.Request.RequestURI == "/hc/login" ||
		ctx.Request.RequestURI == "/hc/loginCheck" || ctx.Request.RequestURI == "/filter/login" || ctx.Request.RequestURI == "/filter/loginCheck" ||
		strings.Index(ctx.Request.RequestURI, "/distribute") != -1 {

	} else if strings.Index(ctx.Request.RequestURI, "/platform") != -1 {
		session_id, admin_id := ctx.GetCookie("user_session"), ctx.GetCookie("user_id")
		if session_id == "" || admin_id == "" {
			ctx.Redirect(302, "/login")
		} else {
			o := controllers.O
			var maps []orm.Params
			num, err := o.Raw("select * from ss_sessions where userid=?", admin_id).Values(&maps)
			if err != nil {
				ctx.Redirect(302, "/login")
			} else {
				if num == 0 {
					ctx.Redirect(302, "/login")
				} else {
					if maps[0]["sessionid"] != session_id {
						ctx.Redirect(302, "/login")
					}

					//这里可以做权限限制的验证

				}
			}
		}
	} else if strings.Index(ctx.Request.RequestURI, "/config") != -1 {
		session_id, admin_id := ctx.GetCookie("admin_session"), ctx.GetCookie("admin_id")
		if session_id == "" || admin_id == "" {
			ctx.Redirect(302, "/login")
		} else {
			o := controllers.O
			var maps []orm.Params
			num, err := o.Raw("select * from ss_sessions where userid=?", admin_id).Values(&maps)
			if err != nil {
				ctx.Redirect(302, "/login")
			} else {
				if num == 0 {
					ctx.Redirect(302, "/login")
				} else {
					if maps[0]["sessionid"] != session_id {
						ctx.Redirect(302, "/login")
					}

					//这里可以做权限限制的验证

				}
			}
		}
	} else if strings.Index(ctx.Request.RequestURI, "/hc") != -1 {
		session_id, admin_id := ctx.GetCookie("HC_session"), ctx.GetCookie("HCAdmin_id")
		if session_id == "" || admin_id == "" {
			ctx.Redirect(302, "/hc/login")
		} else {
			o := controllers.O
			var maps []orm.Params
			num, err := o.Raw("select * from ss_sessions where userid=?", admin_id).Values(&maps)
			if err != nil {
				ctx.Redirect(302, "/hc/login")
			} else {
				if num == 0 {
					ctx.Redirect(302, "/hc/login")
				} else {
					if maps[0]["sessionid"] != session_id {
						ctx.Redirect(302, "/hc/login")
					}

					//这里可以做权限限制的验证

				}
			}
		}
	} else if strings.Index(ctx.Request.RequestURI, "/filter") != -1 {
		session_id, admin_id := ctx.GetCookie("filter_session_id"), ctx.GetCookie("filter_admin_id")
		if session_id == "" || admin_id == "" {
			ctx.Redirect(302, "/filter/login")
		} else {
			o := controllers.O
			var maps []orm.Params
			num, err := o.Raw("select * from ss_sessions where userid=?", admin_id).Values(&maps)
			if err != nil {
				ctx.Redirect(302, "/filter/login")
			} else {
				if num == 0 {
					ctx.Redirect(302, "/filter/login")
				} else {
					if maps[0]["sessionid"] != session_id {
						ctx.Redirect(302, "/filter/login")
					}

					//这里可以做权限限制的验证

				}
			}
		}
	} else if strings.Index(ctx.Request.RequestURI, "/throwscreen") != -1 &&
		strings.Index(ctx.Request.RequestURI, "/throwscreen/login") == -1 &&
		strings.Index(ctx.Request.RequestURI, "/throwscreen/loginCheck") == -1 {
		admin_id := ctx.GetCookie("user_id")
		if admin_id == "" {
			ctx.Redirect(302, "/throwscreen/login")
		}
	}
}

func Start() {
	//日志输出

	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/test.log","maxdays":`+beego.AppConfig.String("logsdays")+`,"separate":[ "alert", "error", "warning", "info", "debug"]}`)
	//链接数据库
	models.DBBase()
	// 创建一个公共 Ormer
	// 需要 切换数据库 和 事务处理 的话，不要使用全局保存的 Ormer 对象。
	controllers.O = orm.NewOrm()
	controllers.FilterStatus = controllers.SelSystemFilter()
	//app数据同步
	controllers.AppServiceMain()

	//数据库验证
	conf.Sql()

	//主机状态验证
	controllers.InitHostTask()

	signs.Inits()

	//测试海顿 三维接口
	//controllers.TestPostImage()
	//暂停控制台log输出
	//beego.BeeLogger.DelLogger("console")
	//打开session
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run()
}
