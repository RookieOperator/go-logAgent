package config

// Kafka 定义Kafka结构体
type KafkaConfig struct {
	Address string `ini:"address"`
	//Topic string `ini:"topic"`
	MaxSize int `ini:"chan_max_size"`
}

// Etcd 定义Etcd结构体
type EtcdConfig struct {
	Address string `ini:"address"`
	Topic string `ini:"topic"`
}

// TailLog 定义TailLog结构体
type TailLogConfig struct {
	FileName string `ini:"filename"`
}

//
type AppConfig struct {
	KafkaConfig `ini:"kafka"`
	TailLogConfig `ini:"tailLog"`
	EtcdConfig `ini:"etcd"`
}


//// LoadConfig 加载配置文件
//func LoadConfig(cfg *AppConfig) (cfg *AppConfig,err error){
//	err = ini.MapTo(cfg,"./config/config.ini")
//	if err != nil {
//		return nil,err
//	}
//	return cfg,nil
//}
