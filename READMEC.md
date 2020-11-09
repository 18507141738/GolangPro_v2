## Artifice_V2.0.0

炼化升级版

---------------------------------------------mysql数据库操作-------------------------------------------------- 
注：当程序自动插入字段异常时，使用下列方法手动插入数据 
1.数据库所在设备上打开终端输入"mysql"回车登录数据库 
2.终端输"use Artifice;"进入数据库 
3.运行对应"数据库更新操作"文件命令

---------------------------------------------基础数据库文件---------------------------------- 
Artifice.sql基础数据库文件

------------------------------------------------基本功能---------------------------------- 

着装监测：cloths
烟雾监测：smoke
火苗监测：fire
区域入侵：boundary
离岗监测：queue_count V1.0.0版(person_count)
睡岗监测：sleep_count
泄漏监测：leakage

--------------------------------------------设置服务器分配摄像头上限---------------------------------- 
1.修改tests/test.cfg maxCameraNum值即可，无需重启程序

***以下内容详见有道云笔记(https://note.youdao.com/group/) --------------------------------------------Golang环境安装准则----------------------------------------------

(a)用之前应该用vlc插件进行测试RTMP流是否正常拉流(异常的摄像机自行记录并上报管理员),并记录摄像机尺寸大小,正常的流后面添加播放

(b)golang环境安装: vim ~/.bashrc shift+G 尾部追加: export GOROOT=/usr/local/go export GOPATH=~/workspace:~/goproject export GOBIN=~/gobin export PATH=$PATH:$GOROOT/bin:$GOBIN 终端运行go version go version go1.12 linux/amd64即安装成功.

(c)GO代码配置 进入~目录 mkdir -p workspace/src | cd workspace/src cp Artifice_V1.0 . 即可 进入 Artifice_V1.0 vim util/DbList.go 相关参数修改完成即可

(d)代码运行 主目录下执行 make

--------------------------------------------WEB前端应用资源及说明----------------------------------------------

(a)运行goweb.sql文件,将备份的数据库结构导入 ***注意:线下数据库会出现group by 不支持的现象,添加以下语句到mysql.conf/mysqld区域 sql_mode=STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION

(b)前端可播放rtmp列表是(只在测试环境下有效): rtmp://192.168.10.120/live/livestreamsb rtmp://192.168.10.120/live/livestreamsc rtmp://192.168.10.120/live/livestreamsd rtmp://192.168.10.120/live/livestreamse 门口 rtmp://192.168.10.120/live/livestreamsf rtmp://192.168.10.120/live/livestreamsg rtmp://192.168.10.120/live/livestreamsh rtmp://192.168.10.120/live/livestreamshab 办公区 (谨慎使用,经常报错) rtmp://192.168.10.120/live/livestreamshac 会议室

平台功能(英文)总览: 中文 function_type 员工岗位 man_work 着装监控 cloths* 课堂考勤 attence_class 安全防盗 Thief_defend 烟雾检测 warning_part 人员超限 count_limit 身份核实 identify
徘徊监测 walks 区域入侵 boundary* 火焰检测 fire* 烟雾检测 smoke* 睡岗检测 sleep_count* 离岗检测 leave_count* //person_count 人头计数 person_count
泄漏检测 leakage 黑白名单 blacklist 无感考勤 attendance 客流量 flush_count 排队人数 queue_count 卸油区检测 unload_area 抽烟打电话 phone_smoke 跑冒滴漏 leakage

炼化需求: 系统异常断开情况处理(假设算法服务器与控制器分开): 1.检测磁盘 2.断电数据自动转移 a,主控制器断电,数据库同步(主从复制) b,备控制器(算法服务器)断电,信号传递(socket通信,异常报错) 3.主,备控制器之间互相以socket连接,主断电后启用备控制器上的web服务,摄像机设备重新分布到其余存活的算法服务器上 炼化厂需求以及额外可调参数（除坐标,阈值外）: 1.烟雾检测 perimeterThreshold 轮廓周长 bg_threshold 背景剪除 5.火焰检测 perimeterThreshold bg_threshold 同上 2.区域入侵 DistanceMode 聚焦模式 3.着装监控 无额外可调参数 4.离岗检测 almost_sleep_seconds 疑似睡眠秒数 real_sleep_minutes 真实睡眠分钟数 num_person_threshold 人数上限 6.睡岗检测 almost_sleep_seconds real_sleep_minutes num_person_threshold 同上

