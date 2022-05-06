package conf

type AppConf struct {
	KafkaConf `ini:"kafka"` //如果是结构体，则指明ini中的分区
	EtcdConf  `ini:"etcd"`
}

type KafkaConf struct {
	Address     string `ini:"address"` //ini配置文件分区中的key
	ChanMaxSize int    `ini:"chan_max_siz"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Key     string `ini:"collect_log_key"`
	Timeout int    `ini:"timeOut"`
}
