package controllers

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"net"
	"os"
	"time"
)

func SelSystemTitleUtil() string {
	o := O
	var maps []orm.Params
	num, err := o.Raw("select * from ss_system").Values(&maps)
	if err != nil {
		return ""
	} else {
		if num == 0 {
			return ""
		} else {
			return maps[0]["websitename"].(string)
		}
	}
}

func SelSystemFilter() bool {
	o := O
	var maps []orm.Params
	num, err := o.Raw("select * from ss_system").Values(&maps)
	if err != nil {
		return true
	}
	if num == 0 {
		return true
	}
	if maps[0]["filterstatus"] == "0" {
		return false
	}

	return true
}

func GetTimeSE() string {
	var startTime = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	var endTime = time.Now().Format("2006-01-02")
	return startTime + " - " + endTime
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getLocalIp() (ip string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}
