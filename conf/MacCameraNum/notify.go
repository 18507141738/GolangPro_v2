package MacCameraNum

/*定义这样一个Notifyer接口，只要实现了Callback方法的对象，就都实现了这个Notifyer接口。
实现了这个接口的对象，如果都需要被通知配置文件更新，那这些对象都可以加入到前面定义的
notifyList []Notifyer这个切片中，等待被通知配置文件更新。*/
type Notifyer interface {
	Callback(*Config)
}

func (c *Config) AddObserver(n Notifyer) {
	c.notifyList = append(c.notifyList, n)
}
