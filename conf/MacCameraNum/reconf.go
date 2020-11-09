package MacCameraNum

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	filename       string
	data           map[string]string
	lastModifyTime int64
	rwLock         sync.RWMutex
	notifyList     []Notifyer
}

func NewConfig(file string) (conf *Config, err error) {
	//初始化config
	conf = &Config{
		filename: file,
		data:     make(map[string]string, 1024),
	}

	m, err := conf.parse()
	if err != nil {
		fmt.Printf("parse conf error:%v\n", err)
		return
	}

	//将解析配置文件后的数据更新到结构体的map中，写锁
	conf.rwLock.Lock()
	conf.data = m
	conf.rwLock.Unlock()

	//启一个后台线程去检测配置文件是否更改
	go conf.reload()
	return
}

/*
	解析函数，读取配置文件，一行行解析存放到map中
*/
func (c *Config) parse() (m map[string]string, err error) {
	//如果在parse()中定义一个map，这样就是一个新的map不用加锁
	m = make(map[string]string, 1024)

	f, err := os.Open(c.filename)
	if err != nil {
		println(err)
		return
	}

	reader := bufio.NewReader(f)
	//声明一个变量存放读取行数
	var lineNo int
	for {
		line, errRet := reader.ReadString('\n')
		if errRet == io.EOF {
			//最后一行必须添加\n，否则会漏读
			lineParse(&lineNo, &line, &m)
			break
		}
		if errRet != nil {
			err = errRet
			return
		}

		lineParse(&lineNo, &line, &m)
	}

	return
}

func lineParse(lineNO *int, line *string, m *map[string]string) {
	*lineNO++

	l := strings.TrimSpace(*line)
	//如果空行 或者 是注释 跳过
	if len(l) == 0 || l[0] == '\n' || l[0] == '#' || l[0] == ';' {
		return
	}

	itemSlice := strings.Split(l, "=")

	if len(itemSlice) == 0 {
		fmt.Println("invalid config,line:%d", lineNO)
		return
	}

	key := strings.TrimSpace(itemSlice[0])
	if len(key) == 0 {
		fmt.Println("invalid config,line:%d", lineNO)
		return
	}
	if len(key) == 1 {
		(*m)[key] = ""
		return
	}

	value := strings.TrimSpace(itemSlice[1])
	(*m)[key] = value

	return
}

func (c *Config) reload() {
	//定时 用time.NewTicker每隔5秒去检查一下配置文件
	ticker := time.NewTicker(time.Second * 5)

	for _ = range ticker.C {
		// 打开文件
		// 为什么使用匿名函数？ 当匿名函数退出时可用defer去关闭文件
		// 如果不用匿名函数，在循环中不好关闭文件，一不小心就内存泄露
		func() {
			f, err := os.Open(c.filename)
			if err != nil {
				fmt.Printf("open file error:%s\n", err)
				return
			}
			defer f.Close()

			fileInfo, err := f.Stat()
			if err != nil {
				fmt.Printf("stat file error:%s\n", err)
				return
			}
			//获取当前文件修改时间
			curModifyTime := fileInfo.ModTime().Unix()
			//如果配置文件的修改时间比上一次修改时间大，我们认为配置文件更新了。
			// 那么我们调用parse()解析配置文件，并更新conf实例中的数据。
			// 并且更新配置文件的修改时间
			if curModifyTime > c.lastModifyTime {
				//重新解析时，要考虑应用程序正在读取这个配置因此应该加锁
				m, err := c.parse()
				if err != nil {
					fmt.Printf("parse config  error:%v\n", err)
					return
				}

				c.rwLock.Lock()
				c.data = m
				c.rwLock.Unlock()

				c.lastModifyTime = curModifyTime

				//配置更新通知所有观察者
				for _, n := range c.notifyList {
					n.Callback(c)
				}
			}

		}()

	}

}

//封装 返回int值
func (c *Config) GetInt(key string) (value int, err error) {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	str, ok := c.data[key]
	if !ok {
		err = fmt.Errorf("key [%s] not found", key)
	}
	value, err = strconv.Atoi(str)
	return
}

func (c *Config) GetString(key string) (value string, err error) {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	value, ok := c.data[key]
	if !ok {
		err = fmt.Errorf("key [%s] not found", key)
	}
	return
}
