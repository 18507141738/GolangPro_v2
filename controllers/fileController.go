package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"log"
	"os"
	"time"
)

type FileLock struct {
	dir         string
	f           *os.File
	writhOrRead string
}

func New(dir string) *FileLock {
	return &FileLock{
		dir: dir,
	}
}

//加锁 O_TRUNC覆盖插入
func (l *FileLock) Lock() error {
	f, err := os.OpenFile(l.dir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	l.f = f
	l.writhOrRead = "w"
	//err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	//if err != nil {
	//	return fmt.Errorf("cannot flock directory %s - %s", l.dir, err)
	//}
	return nil
}

//加锁
func (l *FileLock) Lock2() error {
	f, err := os.OpenFile(l.dir, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	l.f = f
	l.writhOrRead = "r"
	//err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	//if err != nil {
	//	return fmt.Errorf("cannot flock directory %s - %s", l.dir, err)
	//}
	return nil
}

//释放锁
func (l *FileLock) Unlock() error {
	defer l.f.Close()
	l.writhOrRead = ""
	//return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
	return nil
}

func WriteFile(param orm.Params) {
	f_lock := New("./conf/funcCamera.json")
	if f_lock.writhOrRead == "r" {
		LogsError("json文件正被读取中...:")
		return
	}
	err := f_lock.Lock()
	if err != nil {
		LogsError("lock_err:", err)
		return
	}
	//创建json编码器
	encode := json.NewEncoder(f_lock.f)
	err = encode.Encode(param)
	if err != nil {
		LogsError("写入funcCamera.json数据异常", err)
	}
	f_lock.Unlock()
}

func ReadFile() orm.Params {
	f_lock := New("./conf/funcCamera.json")
	if f_lock.writhOrRead == "w" {
		// 创建一个计时器
		timeTicker := time.NewTicker(time.Second * 2)
		i := 0
		for {
			if i > 2 {
				if f_lock.writhOrRead == "w" {
					LogsError("json文件正被写入中...:")
					return nil
				}
				break
			}

			i++
			<-timeTicker.C //计时暂停2s
		}
		// 清理计时器
		timeTicker.Stop()
	}
	err := f_lock.Lock2()
	if err != nil {
		//// 创建一个计时器
		//timeTicker := time.NewTicker(time.Second * 2)
		//i := 0
		//for {
		//	if i > 2 {
		//		if err != nil {
		//			return nil
		//		}
		//		break
		//	}
		//
		//	err = f_lock.Lock2()
		//	if err != nil {
		//		logs.Error("json文件写入中...:", err)
		//	}
		//	i++
		//	<-timeTicker.C
		//
		//}
		//// 清理计时器
		//timeTicker.Stop()
		LogsError("lock_err:", err)
		return nil
	}

	var person orm.Params
	//创建json解码器
	decoder := json.NewDecoder(f_lock.f)
	err = decoder.Decode(&person)
	if err != nil {
		LogsError("读取funcCamera.json数据异常", err)
	}
	f_lock.Unlock()
	return person
}

func readFile(path string) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	defer file.Close()
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	var person orm.Params
	//创建json解码器
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&person)
	if err != nil {
		fmt.Println("err2:", err)
	} else {
		fmt.Println("success:", person)
	}

}

func insertFile(path string, param orm.Params) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	//创建json编码器
	encode := json.NewEncoder(file)
	err = encode.Encode(param)
	if err != nil {
		log.Println("err3:", err)
	} else {
		log.Print("success")
	}

}
