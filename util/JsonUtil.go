package util

import (
	"github.com/astaxie/beego/orm"
)

type ResultA struct {
	Code   int
	Num    int64
	Reason *[]map[string]string //*map[string]map[string]string
}

type ResultB struct {
	Code   int
	Num    int64
	Reason *[]orm.Params //*map[string]map[string]string
}
type ResultC struct {
	Ret    int
	Reason string
	Data   *[]orm.Params //*map[string]map[string]string
}

type ResultD struct {
	Ret    int
	Reason string
	Url    string //*map[string]map[string]string
}

type ResultE struct {
	Code string
}

type ResultF struct {
	Ret    int
	Reason string
	Id     string
	Image  string
}
type Result struct {
	Ret      int
	Secrekey string
	Reason   string
}
