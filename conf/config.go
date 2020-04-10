package conf

type AppConf struct {
	KafkaConf	`ini:"kafka"`
	TailLogConf `ini:"taillog"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic string	`ini:"topic"`
}

type TailLogConf struct {
	Filename string	`ini:"filename"`
}