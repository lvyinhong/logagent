package conf

type AppConf struct {
	KafkaConf	`ini:"kafka"`
	TailLogConf `ini:"taillog"`
	EtcdConf `ini:"ectd"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic string	`ini:"topic"`
}

type TailLogConf struct {
	Filename string	`ini:"filename"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Username string `ini:"test"`
	Password string `ini:"123456"`
}