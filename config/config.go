package config

// Kafka 定义Kafka结构体
type KafkaConfig struct {
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
}


//// LoadConfig 加载配置文件
//func LoadConfig(cfg *AppConfig) (cfg *AppConfig,err error){
//	err = ini.MapTo(cfg,"./config/config.ini")
//	if err != nil {
//		return nil,err
//	}
//	return cfg,nil
//}
