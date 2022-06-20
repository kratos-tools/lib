package opg

import (
	"time"

	"github.com/kratos-tools/lib/pkg/util/otime"
	"log"
)

// Config options
type Config struct {
	Name string
	// 连接地址
	Host string `json:"host"`
	// 端口
	Port int `json:"port"`
	// 用户名
	User string `json:"user"`
	// 密码
	Password string `json:"password"`
	// 数据库名称
	Db string `json:"db"`
	// 是否开启ssl连接
	SslMode string `json:"sslMode"`
	// 最大空闲连接数
	MaxIdleConns int `json:"maxIdleConns"`
	// 最大活动连接数
	MaxOpenConns int `json:"maxOpenConns"`
	// 连接的最大存活时间
	ConnMaxLifetime time.Duration `json:"connMaxLifetime"`
	// 慢日志阈值
	SlowThreshold time.Duration `json:"slowThreshold"`
	// 关闭指标采集
	DisableMetric bool `json:"disableMetric"`
	// Debug开关
	Debug         bool `json:"debug"`
	SingularTable bool `json:"singularTable"`

	// 日志
	//_logger *log.Logger

	// 自动更新时间字段
	UpdateColumn []string `json:"update_column"`
}

// New 根据外部config构建独立config
func New(conf *Config) *Config {
	defaultConfig := DefaultConfig()
	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = defaultConfig.MaxIdleConns
	}
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = defaultConfig.MaxOpenConns
	}
	if conf.SlowThreshold == 0 {
		conf.SlowThreshold = defaultConfig.SlowThreshold
	}
	if conf.ConnMaxLifetime == 0 {
		conf.ConnMaxLifetime = defaultConfig.ConnMaxLifetime
	}
	return conf
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Debug:           false,
		SslMode:         "disable",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: otime.Duration("300s"),
		SlowThreshold:   otime.Duration("500ms"),
		//_logger:         klog.NewLogger(),
	}
}

// Build ...
func (c *Config) Build() *DB {

	db, err := Open(c)
	if err != nil {
		log.Panicf("connect pg err %s", err)
	}

	if err := db.DB().Ping(); err != nil {
		log.Panicf("ping pg err %s", err)
	}
	return db
}
