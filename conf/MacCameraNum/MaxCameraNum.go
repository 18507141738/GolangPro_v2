package MacCameraNum

import (
	"fmt"
	"sync/atomic"
	"time"
)

type TestConfig struct {
	maxCameraNum int
}

type TestConfigMgr struct {
	config atomic.Value
}

var testConfigMgr = &TestConfigMgr{}

func (t *TestConfigMgr) Callback(conf *Config) {
	testConfig := &TestConfig{}
	maxCameraNum, err := conf.GetInt("maxCameraNum")
	if err != nil {
		fmt.Printf("get maxCameraNum err:%v\n", err)
		return
	}
	testConfig.maxCameraNum = maxCameraNum

	testConfigMgr.config.Store(testConfig)
}

func initConfig(file string) {

	//打开配置文件
	conf, err := NewConfig(file)
	if err != nil {
		fmt.Printf("read config file err:%v\n", err)
		return
	}

	//添加观察者
	conf.AddObserver(testConfigMgr)

	//第一次读取配置文件
	var testConfig TestConfig
	testConfig.maxCameraNum, err = conf.GetInt("maxCameraNum")
	if err != nil {
		fmt.Printf("get maxCameraNum err:%v\n", err)
		return
	}

	//读取到的文件数据保存到atomic.value
	testConfigMgr.config.Store(&testConfig)

}

func run() {
	for {
		testConfig := testConfigMgr.config.Load().(*TestConfig)

		fmt.Println("maxCameraNum:", testConfig.maxCameraNum)
		fmt.Printf("%v\n", "--------------------")
		time.Sleep(5 * time.Second)
	}
}

//默认50
func GetMaxCameraNum() (value int) {
	tc := testConfigMgr.config.Load()
	if tc != nil {
		testConfig := tc.(*TestConfig)
		value = testConfig.maxCameraNum
	} else {
		value = 0
	}

	return
}

func Maingo() {
	confFile := "./tests/test.cfg"
	initConfig(confFile)
	//run()
}
