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
	CollectLogInfo string `ini:"collect_log_info"`
}

// LogConfig 定义LogConfig结构体
type LogConfig struct {
	Path string `ini:"path"`
	MaxSize int
	MaxBackups int
	MaxAge int
	Compress bool
}

//
type AppConfig struct {
	KafkaConfig `ini:"kafka"`
	TailLogConfig `ini:"tailLog"`
	EtcdConfig `ini:"etcd"`
	LogConfig `ini:"log"`
}


