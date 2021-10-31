package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	C         Conf
	viperConf *viper.Viper
)

type Conf struct {
	JobEventChan  int      `mapstructure:"job-event-chan" json:"job-event-chan" yaml:"job-event-chan"`
	AgentEndpoint string   `mapstructure:"agent-endpoint" json:"agent-endpoint" yaml:"agent-endpoint"`
	NoCalculation []string `mapstructure:"no-calculation" json:"no-calculation" yaml:"no-calculation"`

	Etcd  Etcd  `mapstructure:"etcd" json:"etcd" yaml:"etcd"`
	Mongo Mongo `mapstructure:"mongo" json:"mongo" yaml:"mongo"`
	Redis Redis `mapstructure:"redis" json:"redis" yaml:"redis"`
	Zap   Zap   `mapstructure:"zap" json:"zap" yaml:"zap"`
}

// System 系统配置
type System struct {
	Addr string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Env  string `mapstructure:"env" json:"env" yaml:"env"`
}

// Etcd etcd配置
type Etcd struct {
	Endpoints []string `mapstructure:"endpoints" json:"endpoints" yaml:"endpoints"`
	Timeout   int      `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
}

// Mongo mongodb配置
type Mongo struct {
	Endpoints string `mapstructure:"endpoints" json:"endpoints" yaml:"endpoints"`
	Timeout   int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	Db        string `mapstructure:"db" json:"db" yaml:"db"`
}

// Redis 配置
type Redis struct {
	DB             int    `mapstructure:"db" json:"db" yaml:"db"`
	Addr           string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password       string `mapstructure:"password" json:"password" yaml:"password"`
	MaxRetryTimes  int    `mapstructure:"max-retry-times" json:"MaxRetryTimes" yaml:"max-retry-times"`
	MaxIdle        int    `mapstructure:"max-idle" json:"maxIdle" yaml:"max-idle"`
	MaxActive      int    `mapstructure:"max-active" json:"maxActive" yaml:"max-active"`
	MaxIdleTimeout int    `mapstructure:"max-idle-timeout" json:"maxIdleTimeout" yaml:"max-idle-timeout"`
}

// Zap 日志配置
type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`
	Format        string `mapstructure:"format" json:"format" yaml:"format"`
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`
	LinkName      string `mapstructure:"link-name" json:"linkName" yaml:"link-name"`
	ShowLine      bool   `mapstructure:"show-line" json:"showLine" yaml:"showLine"`
	EncodeLevel   string `mapstructure:"encode-level" json:"encodeLevel" yaml:"encode-level"`
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key"`
	LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console"`
}

// InitConf 初始化配置
func InitConf(confPath string) (err error) {
	v := viper.New()
	v.SetConfigFile(confPath)
	err = v.ReadInConfig()
	if err != nil {
		return
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.Unmarshal(&C); err != nil {
			fmt.Println(err)
		}
	})

	if err = v.Unmarshal(&C); err != nil {
		return
	}
	viperConf = v
	return
}
