package main

import (
	"Artifice_V2.0.0/routers"
	"github.com/astaxie/beego/orm"
)

func main() {
	//debug模式打印信息
	orm.Debug = true

	routers.Start()
}
