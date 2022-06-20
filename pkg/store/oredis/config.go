package oredis

import (
	klog "github.com/silenceper/log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/kratos-tools/lib/pkg/util/otime"
)

// Config for oredis, contains RedisStubConfig and RedisClusterConfig
type Config struct {
	Name string `json:"name"`
	// Addr 实例配置地址
	Addr string `json:"addr"`
	// Password 密码
	Password string `json:"password"`
	// DB，默认为0, 一般应用不推荐使用DB分片
	DB int `json:"db"`
	// Maximum number of idle connections in the pool.
	MaxIdle int `json:"maxIdle"`
	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int `json:"maxActive"`
	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout int `json:"idleTimeout"`
	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait bool `json:"wait"`
	// 慢日志门限值，超过该门限值的请求，将被记录到慢日志中
	SlowThreshold time.Duration `json:"slowThreshold"`
}

// New 根据外部config构建独立config
func New(conf *Config) *Config {
	defaultConfig := DefaultRedisConfig()
	if conf.MaxIdle == 0 {
		conf.MaxIdle = defaultConfig.MaxIdle
	}
	if conf.MaxActive == 0 {
		conf.MaxActive = defaultConfig.MaxActive
	}
	if conf.SlowThreshold == 0 {
		conf.SlowThreshold = defaultConfig.SlowThreshold
	}
	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = defaultConfig.IdleTimeout
	}
	return conf
}

// DefaultRedisConfig default config ...
func DefaultRedisConfig() Config {
	return Config{
		DB:            0,
		MaxIdle:       10,
		MaxActive:     500,
		Wait:          true,
		IdleTimeout:   120,
		SlowThreshold: otime.Duration("50ms"),
	}
}

// Build 构建redis客户端
func (c Config) Build() *Redis {
	redisClient := &Redis{
		config: &c,
	}
	redisClient.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			client, err := redis.Dial("tcp", c.Addr)
			if err != nil {
				return nil, err
			}

			if c.Password != "" {
				if _, err := client.Do("AUTH", c.Password); err != nil {
					_ = client.Close()
					return nil, err
				}
			}

			if _, err := client.Do("SELECT", c.DB); err != nil {
				_ = client.Close()
				return nil, err
			}
			return client, nil
		},
		MaxIdle:     c.MaxIdle,
		MaxActive:   c.MaxActive,
		IdleTimeout: time.Second * time.Duration(int64(c.IdleTimeout)),
	}

	if isPong := redisClient.ping(); !isPong {
		klog.Errorf("redis addr: %s/%d Not connect", c.Addr, c.DB)
		return nil
	}

	klog.Infof("redis: %s addr: %s/%d  MaxIdle:%d MaxActive:%d IdleTimeout:%ds connect", c.Name, c.Addr, c.DB, c.MaxIdle, c.MaxActive, c.IdleTimeout)
	return redisClient
}
