package conf

type AppConf struct {
	KafkaConf	`ini:"kafka"`
	TailLogConf `ini:"taillog"`
	EtcdConf `ini:"etcd"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	//Topic string	`ini:"topic"`
	MaxSize int `ini:"chan_max_size"`
}

type TailLogConf struct {
	Filename string	`ini:"filename"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Key string `ini:"collect_log_key"`
	Username string `ini:"username"`
	Password string `ini:"password"`
	Timeout int `init:"timeout"`
}