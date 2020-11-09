package conf

import (
	"Artifice_V2.0.0/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
)

/*
	程序添加后续数据库字段操作
*/
func Sql() {

	addFuncCameraField()
	addHostField()
	addCameraField()
	addAlarmField()
	addTableSession()
	editUserField()
	addAlarmNumber()
	addOvertime()
	addAlgorithm_url()
	addAlgorithm_widthl()
	addAlgorithm_height()
	addVideoCode()
	addTimepoint()
	addAlarm_time_second()
	addVideo_switch()
	addBg_thr()
	addLeakThreshold()
	addAlarmDetail()
	editSystemDB()
	editPlace()
	edit_Orignize()
	addBuBan()
	addFilterUser()
	addAlarmFilter()
}
func addLeakThreshold() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//用户表添加jurisdition
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_func_camera", "leakThreshold").Values(&maps)
	if err != nil {
		log.Println("检测数据库表字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_func_camera ADD leakThreshold varchar(32) DEFAULT '' COMMENT '注释：背景周长';").Exec()
			if err != nil {
				log.Printf("数据库表字段插入异常", err)
			}
		}
	}
}

func addBg_thr() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//用户表添加jurisdition
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_func_camera", "bg_thr").Values(&maps)
	if err != nil {
		log.Println("检测数据库表字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_func_camera ADD bg_thr varchar(32) DEFAULT '' COMMENT '注释：背景减除的阈值默认50';").Exec()
			if err != nil {
				log.Printf("数据库表字段插入异常", err)
			}
		}
	}
}

//摄像机配置表
func addVideo_switch() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//用户表添加jurisdition
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_func_camera", "video_switch").Values(&maps)
	if err != nil {
		log.Println("检测数据库表字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_func_camera ADD video_switch varchar(32) DEFAULT '';").Exec()
			if err != nil {
				log.Printf("数据库表字段插入异常", err)
			}
		}
	}
}

//摄像机配置表
func addAlarm_time_second() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//用户表添加jurisdition
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_func_camera", "alarm_time_second").Values(&maps)
	if err != nil {
		log.Println("检测数据库表字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_func_camera ADD alarm_time_second varchar(255) DEFAULT '';").Exec()
			if err != nil {
				log.Printf("数据库表字段插入异常", err)
			}
		}
	}
}

//摄像机配置表
func addTimepoint() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//用户表添加jurisdition
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_func_camera", "timepoint").Values(&maps)
	if err != nil {
		log.Println("检测数据库表字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_func_camera ADD timepoint varchar(255) DEFAULT '';").Exec()
			if err != nil {
				log.Printf("数据库表字段插入异常", err)
			}
		}
	}
}

//摄像机表
func addVideoCode() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//摄像添加新参数updateTime
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_camera", "videoCode").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_camera,updateTime字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_camera ADD videoCode VARCHAR(255) DEFAULT NULL COMMENT '注释：0:H264/1:h265';").Exec()
			if err != nil {
				log.Printf("数据库表ss_camera,updateTime字段插入异常", err)
			}
		}
	}
}

//摄像机表
func addAlgorithm_widthl() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//摄像添加新参数updateTime
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_camera", "algorithm_width").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_camera,updateTime字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_camera ADD algorithm_width VARCHAR(255) DEFAULT NULL COMMENT '注释：计算流宽度';").Exec()
			if err != nil {
				log.Printf("数据库表ss_camera,updateTime字段插入异常", err)
			}
		}
	}
}

//摄像机表
func addAlgorithm_height() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//摄像添加新参数updateTime
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_camera", "algorithm_height").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_camera,updateTime字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_camera ADD algorithm_height VARCHAR(255) DEFAULT NULL COMMENT '注释：计算流高度';").Exec()
			if err != nil {
				log.Printf("数据库表ss_camera,updateTime字段插入异常", err)
			}
		}
	}
}

//摄像机表
func addAlgorithm_url() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//摄像添加新参数updateTime
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_camera", "algorithm_url").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_camera,updateTime字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_camera ADD algorithm_url VARCHAR(255) DEFAULT NULL COMMENT '注释：计算流';").Exec()
			if err != nil {
				log.Printf("数据库表ss_camera,updateTime字段插入异常", err)
			}
		}
	}
}

//摄像机表
func addCameraField() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//摄像添加新参数updateTime
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_camera", "updateTime").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_camera,updateTime字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_camera ADD updateTime datetime(0) DEFAULT NULL COMMENT '注释：更新时间';").Exec()
			if err != nil {
				log.Printf("数据库表ss_camera,updateTime字段插入异常", err)
			}
		}
	}

	_, err = o.Raw("SELECT data_type as type from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_camera", "place_id").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_camera异常", err)
	} else if len(maps) > 0 {
		if maps[0]["type"] == "int" {
			_, err = o.Raw("ALTER TABLE ss_camera MODIFY COLUMN place_id varchar(32);").Exec()
			if err != nil {
				log.Printf("数据库表字段插入异常", err)
			}
		}
	}

	//摄像机默认编辑状态
	_, err = o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =? AND column_name=?;", dbname, "ss_camera", "status").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_camera,status字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			_, err = o.Raw("ALTER TABLE ss_camera ADD `status` tinyint(0) NULL DEFAULT 0 COMMENT '这里是注释：默认编辑状态 0未编辑 1已编辑';").Exec()
			if err != nil {
				log.Printf("数据库表ss_camera,status字段插入异常", err)
			}
		}
	}
}

//摄像机配置表
func addFuncCameraField() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//摄像机配置添加新参数detect_number
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =? AND column_name=?;", dbname, "ss_func_camera", "detect_number").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_func_camera,detect_number字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入detect_number字段")
			_, err = o.Raw("ALTER TABLE ss_func_camera ADD detect_number VARCHAR(255) DEFAULT NULL COMMENT '这里是注释：检测阈值';").Exec()
			if err != nil {
				log.Printf("数据库表ss_func_camera,detect_number字段插入异常", err)
			}
		}
	}
	//摄像机配置添加新参数warning_number
	_, err = o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_func_camera", "warning_number").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_func_camera,warning_number字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入detect_number字段")
			_, err = o.Raw("ALTER TABLE ss_func_camera ADD warning_number VARCHAR(255) DEFAULT NULL COMMENT '这里是注释：告警阈值';").Exec()
			if err != nil {
				log.Printf("数据库表ss_func_camera,warning_number字段插入异常", err)
			}
		}
	}
	//摄像机配置添加新参数departure_minute
	_, err = o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =? AND column_name=?;", dbname, "ss_func_camera", "departure_minute").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_func_camera,departure_minute字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入detect_number字段")
			_, err = o.Raw("ALTER TABLE ss_func_camera ADD departure_minute INT(10) DEFAULT NULL COMMENT '这里是注释：离岗分钟数';").Exec()
			if err != nil {
				log.Printf("数据库表ss_func_camera,departure_minute字段插入异常", err)
			}
		}
	}
	//摄像机配置添加新参数distance
	_, err = o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =? AND column_name=?;", dbname, "ss_func_camera", "distance").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_func_camera,distance字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入detect_number字段")
			_, err = o.Raw("ALTER TABLE ss_func_camera ADD distance VARCHAR(20) DEFAULT 1.0 COMMENT '这里是注释：目标距离';").Exec()
			if err != nil {
				log.Printf("数据库表ss_func_camera,distance字段插入异常", err)
			}
		}
	}

	//摄像机配置添加新参数boundary2
	_, err = o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =? AND column_name=?;", dbname, "ss_func_camera", "boundary2").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_func_camera,boundary2字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			_, err = o.Raw("ALTER TABLE ss_func_camera ADD boundary2 TEXT CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL COMMENT '这里是注释：负样本坐标';").Exec()
			if err != nil {
				log.Printf("数据库表ss_func_camera,boundary2字段插入异常", err)
			}
		}
	}

	//摄像机配置添加标准负样本
	_, err = o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =? AND column_name=?;", dbname, "ss_func_camera", "negative").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_func_camera,negative字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			_, err = o.Raw("ALTER TABLE ss_func_camera ADD negative TEXT CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL COMMENT '这里是注释：由算法计算负样本坐标';").Exec()
			if err != nil {
				log.Printf("数据库表ss_func_camera,negative字段插入异常", err)
			}
		}
	}
}

//主机配置表
func addAlarmNumber() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//主机表添加新参数updateTime
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =? AND column_name=?;", dbname, "ss_func_host", "alarmNumber").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_func_host,updateTime字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_func_host ADD alarmNumber VARCHAR(255)  DEFAULT NULL COMMENT '注释：报警编号';").Exec()
			if err != nil {
				log.Printf("数据库表ss_func_host,updateTime字段异常", err)
			}
		}
	}
}

//主机配置表
func addOvertime() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//主机表添加新参数updateTime
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =? AND column_name=?;", dbname, "ss_func_host", "overtime").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_func_host,updateTime字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_func_host ADD overtime VARCHAR(32)  DEFAULT NULL COMMENT '注释：主机超时状态';").Exec()
			if err != nil {
				log.Printf("数据库表ss_func_host,updateTime字段异常", err)
			}
		}
	}
}

//主机配置表
func addHostField() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//主机表添加新参数updateTime
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =? AND column_name=?;", dbname, "ss_func_host", "updateTime").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_func_host,updateTime字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_func_host ADD updateTime datetime(0) DEFAULT NULL COMMENT '注释：更新时间';").Exec()
			if err != nil {
				log.Printf("数据库表ss_func_host,updateTime字段异常", err)
			}
		}
	}
}

//告警表
func addAlarmField() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//主机表添加新参数updateTime
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =? AND column_name=?;", dbname, "ss_alarm", "syn").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_alarm,syn字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_alarm ADD syn int(10) DEFAULT 0 COMMENT '注释：是否外网同步';").Exec()
			if err != nil {
				log.Printf("数据库表ss_alarm,syn字段异常", err)
			}
		}
	}
	_, erra := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =? AND column_name=?;", dbname, "ss_alarm", "uuids").Values(&maps)
	if erra != nil {
		log.Println("检测数据库表ss_alarm,uuids字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_alarm ADD uuids varchar(50) DEFAULT '' COMMENT '预警主序号;'").Exec()
			if err != nil {
				log.Printf("数据库表ss_alarm,uuids字段异常", err)
			} else {
				beego.Info("新增ss_alarm.uuids~")
			}
		}
	}

	//主机表添加新参数page_read_type
	_, err = o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =? AND column_name=?;", dbname, "ss_alarm", "page_read_type").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_alarm,page_read_type字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_alarm ADD page_read_type tinyint(1) DEFAULT 0 COMMENT '注释：0未读，1已读';").Exec()
			if err != nil {
				log.Printf("数据库表ss_alarm,page_read_type字段异常", err)
			}
		}
	}

	//主机表添加新参数alarmhandler
	_, err = o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =? AND column_name=?;", dbname, "ss_alarm", "alarmhandler").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_alarm,alarmhandler字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_alarm ADD alarmhandler tinyint(1) DEFAULT 0 COMMENT '注释：0未处理，1已处理';").Exec()
			if err != nil {
				log.Printf("数据库表ss_alarm,alarmhandler字段异常", err)
			}
		}
	}

}

func editUserField() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_user", "jurisdiction").Values(&maps)
	if err != nil {
		log.Println("检测数据库表字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			_, err = o.Raw("ALTER TABLE ss_user ADD jurisdiction varchar(255) DEFAULT NULL COMMENT '注释：权限：0员工劳保上岗监测 1厂区烟雾检测 2厂区火苗事件监测 3高危区域安全监测 4员工离岗行为监测 5员工睡岗行为监测 6泄漏监测 7设备管理 8人员管理';").Exec()
			if err != nil {
				log.Printf("数据库表字段插入异常(jurisdiction)", err)
			}
		}
	}

	_, erra := o.Raw("SELECT data_type as type from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_user", "tele").Values(&maps)
	if erra != nil {
		log.Println("检测数据库表ss_user异常", err)
	}
	if len(maps) > 0 {
		//beego.Info(222)
		//已存在此字段
	} else {
		//beego.Info(333)
		_, err = o.Raw("ALTER TABLE ss_user ADD tele varchar(255) DEFAULT '',ADD allowd varchar(255) DEFAULT '';").Exec()
		if err != nil {
			log.Printf("数据库表字段插入异常(tele,allowd)", err)
		} else {
			beego.Info("ss_user新增tele,allowd字段!")
		}
	}

	//_, err = o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_user", "number").Values(&maps)
	//if err != nil {
	//	log.Println("检测数据库表字段异常", err)
	//} else {
	//	if maps[0]["Num"] == "0" {
	//		_, err = o.Raw("ALTER TABLE ss_user ADD number varchar(255) DEFAULT NULL COMMENT '注释：员工编号;").Exec()
	//		if err != nil {
	//			log.Printf("数据库表字段插入异常", err)
	//		}
	//	}
	//}
}

func addTableSession() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//监测表是否存在
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =?;", dbname, "ss_sessions").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_sessions异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			_, err = o.Raw("CREATE TABLE `ss_sessions` ( " +
				"`id` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, " +
				"`userid` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"`sessionid` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"PRIMARY KEY (`id`) USING BTREE " +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;").Exec()
			if err != nil {
				log.Printf("数据库ss_sessions表创建异常", err)
			}
		}
	}
}

func addAlarmDetail() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//监测表是否存在
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =?;", dbname, "ss_alarm_detail").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_alarm_detail异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			_, err = o.Raw("CREATE TABLE `ss_alarm_detail` ( " +
				"`id` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, " +
				"`type` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"`detail` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"PRIMARY KEY (`id`) USING BTREE " +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;").Exec()
			if err != nil {
				log.Printf("数据库ss_alarm_detail表创建异常", err)
			}
		}
	}
}

func editSystemDB() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//用户表添加jurisdition
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_system", "spare_switch").Values(&maps)
	if err != nil {
		log.Println("检测数据库表字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			//log.Printf("开始插入updateTime字段")
			_, err = o.Raw("ALTER TABLE ss_system ADD spare_switch varchar(255) DEFAULT '' COMMENT '注释：备用开关';").Exec()
			if err != nil {
				log.Printf("数据库表字段插入异常", err)
			}
		}
	}

	_, err = o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_system", "filterstatus").Values(&maps)
	if err != nil {
		log.Println("检测数据库表字段异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			_, err = o.Raw("ALTER TABLE ss_system ADD filterstatus tinyint(1) DEFAULT 0 COMMENT '注释：1开：过滤、0关：不过滤';").Exec()
			if err != nil {
				log.Printf("数据库表字段插入异常", err)
			}
		}
	}
}

func editPlace() {
	o := controllers.O
	var maps []orm.Params

	dbname := beego.AppConfig.String("mysqldb")

	_, err := o.Raw("SELECT data_type as type from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_place", "place_id").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_place异常", err)
	} else if len(maps) > 0 {
		if maps[0]["type"] == "int" {
			_, err = o.Raw("ALTER TABLE ss_place MODIFY COLUMN place_id varchar(32);").Exec()
			if err != nil {
				log.Printf("数据库表字段插入异常", err)
			}
		}
	}
}

func edit_Orignize() {
	o := controllers.O
	var maps []orm.Params
	//beego.Info(111)

	dbname := beego.AppConfig.String("mysqldb")

	_, err := o.Raw("SELECT data_type as type from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_organize", "ori_id").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_organize异常", err)
	} else if len(maps) > 0 {
		//beego.Info(222)
		//已存在此字段
	} else {
		//beego.Info(333)

		_, err = o.Raw("ALTER TABLE ss_organize ADD ori_id varchar(255) DEFAULT '';").Exec()
		if err != nil {
			log.Printf("数据库表字段插入异常", err)
		} else {
			beego.Info("ss_organize新增ori_id字段!")
		}
	}
	//_, errc := o.Raw("SELECT data_type as type from information_schema.COLUMNS WHERE table_schema=? and table_name = ? AND column_name=?;", dbname, "ss_organize", "include_ban").Values(&maps)
	//if errc != nil {
	//	log.Println("检测数据库表ss_organize.include_ban异常", errc)
	//} else if len(maps) > 0 {
	//	//beego.Info(222)
	//	//已存在此字段
	//}else {
	//	//beego.Info(333)
	//
	//	_, err = o.Raw("ALTER TABLE ss_organize ADD include_ban varchar(255) DEFAULT '';").Exec()
	//	if err != nil {
	//		log.Printf("数据库表字段插入异常", err)
	//	}else {
	//		beego.Info("ss_organize新增include_ban字段!")
	//	}
	//}
}

func addBuBan() {
	o := controllers.O
	var maps []orm.Params
	dbname := beego.AppConfig.String("mysqldb")
	//监测表是否存在
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name =?;", dbname, "ss_bu_ban").Values(&maps)
	if err != nil {
		log.Println("检测数据库表ss_bu_ban异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			_, err = o.Raw("CREATE TABLE `ss_bu_ban` ( " +
				"`codes` varchar(50) DEFAULT '', " +
				"`names` varchar(50) DEFAULT '', " +
				"`parentcodes` varchar(50) DEFAULT '');").Exec()
			if err != nil {
				log.Printf("数据库ss_bu_ban表创建异常", err)
			} else {
				beego.Info("创建ss_bu_ban表成功~")
			}
		}
	}
}

func addFilterUser() {
	var maps []orm.Params
	dbName := beego.AppConfig.String("mysqldb")
	o := controllers.O
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ?;", dbName, "ss_filter_user").Values(&maps)
	if err != nil {
		controllers.LogsError("监测ss_filter_user异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			_, err = o.Raw("CREATE TABLE `ss_filter_user`  ( " +
				"`id` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, " +
				"`username` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"`password` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"`jurisdiction` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '0,1,2,3,4', " +
				"`organize_id` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"PRIMARY KEY (`id`) USING BTREE " +
				") ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;").Exec()

			_, err = o.Raw("INSERT INTO `ss_filter_user`(`id`, `username`, `password`, `jurisdiction`, `organize_id`) VALUES ('20200624142511289384', 'admin', 'ee10c315eba2c75b403ea99136f5b48d', '0,1,2,3,4', '31571093282099200');").Exec()
			if err != nil {
				log.Printf("数据库创建异常", err)
			}
		}
	}
}

func addAlarmFilter() {
	o := orm.NewOrm()
	var maps []orm.Params
	dbName := beego.AppConfig.String("mysqldb")
	//用户表添加字段
	_, err := o.Raw("SELECT COUNT(*) as Num from information_schema.COLUMNS WHERE table_schema=? and table_name = ?;", dbName, "ss_alarm_filter").Values(&maps)
	if err != nil {
		controllers.LogsError("检测数据库异常", err)
	} else {
		if maps[0]["Num"] == "0" {
			_, err = o.Raw("CREATE TABLE `ss_alarm_filter`  ( " +
				"`alarm_id` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, " +
				"`alarm_type` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"`alarm_place` varchar(80) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"`alarm_place_type` varchar(80) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"`alarm_head` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"`alarm_file` varchar(180) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"`alarm_detial` varchar(150) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"`alarm_time` timestamp(0) DEFAULT CURRENT_TIMESTAMP, " +
				"`alarm_level` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, " +
				"`host_id` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"`camera_id` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, " +
				"`alarm_stream` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"`alarm_video` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL, " +
				"`alarm_time_second` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '注释：报警短视频时长', " +
				"`video_switch` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '0' COMMENT '注释：短视频开关', " +
				"`alarm_stage` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '注释：报警阶段', " +
				"`asynch_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '注释：异步请求id', " +
				"`status` tinyint(1) DEFAULT 0 COMMENT '0未处理、1已插入、2误报', " +
				"`page_read_type` tinyint(1) DEFAULT 0 COMMENT '0未读，1已读', " +
				"`uuids` varchar(50) DEFAULT '' COMMENT '预警主序号', " +
				"PRIMARY KEY (`alarm_id`) USING BTREE " +
				") ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;").Exec()
			if err != nil {
				controllers.LogsError("数据库创建异常", err)
			}
		}
	}
}
