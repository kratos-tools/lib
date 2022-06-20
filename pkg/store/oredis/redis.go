package oredis

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/kratos-tools/lib/pkg/util/ostring"
	klog "github.com/silenceper/log"
)

type Redis struct {
	config *Config
	pool   *redis.Pool
}

const HotKeyNum int = 16 // 热key散列数量

// ping 测试链接
func (r *Redis) ping() bool {
	if result, err := redis.String(r.Do("PING")); err != nil || result != "PONG" {
		return false
	}
	return true
}

// GetConfig 查询配置
func (r *Redis) GetConfig() *Config {
	return r.config
}

// GetConfig 获取格式化的名称
func (r *Redis) GetAddrDB() string {
	return fmt.Sprintf("%s %s/%d", r.config.Name, r.config.Addr, r.config.DB)
}

// Do 执行命令
func (r *Redis) Do(command string, args ...interface{}) (interface{}, error) {
	var pool = r.pool.Get()
	slime := time.Now()
	defer func() {
		if err := pool.Close(); err != nil {
			klog.Errorf("%s: %s/%d pool.Close err:%s", r.config.Name, r.config.Addr, r.config.DB, err)
		}
		// 打印慢查询
		if t := time.Now().Sub(slime); t > r.config.SlowThreshold {
			klog.Warnf("redis: %s do command over %s use: %s command: %s %s", r.config.Name, r.config.SlowThreshold.String(), t.String(), command, args)
		}
	}()
	return pool.Do(command, args...)
}

// Lock redis 加锁
func (r *Redis) Lock(key string, value string, expire int64) bool {
	if result, err := r.Do("SET", key, value, "EX", expire, "NX"); err != nil || result == nil {
		return false
	}
	return true
}

// RmLock redis 移除锁
//     value 为空则直接删除;否则按照value值删除
func (r *Redis) RmLock(key string, value string) (bool, error) {
	if value == "" {
		if _, err := r.Do("DEL", key); err != nil {
			return false, err
		}
		return true, nil
	}
	result, err := redis.Int64(r.Do("EVAL", "if redis.call(\"get\", KEYS[1]) == ARGV[1] then return redis.call(\"del\", KEYS[1]) else return 0 end", 1, key, value))
	if err != nil {
		return false, err
	}
	if result == 0 {
		return false, nil
	}
	return true, nil
}

// GetConnect 获取连接&close方法
func (r *Redis) GetConnect() (redis.Conn, func() error) {
	conn := r.pool.Get()
	return conn, func() error {
		return conn.Close()
	}
}

// GetJSON 查询缓存json数据
func (r *Redis) GetJSON(key string, result interface{}) error {
	b, err := redis.Bytes(r.Do("GET", key))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, result)
}

// GetHotKeyJSON 查询热key缓存json数据
func (r *Redis) GetHotKeyJSON(key string, result interface{}) error {
	rand.Seed(time.Now().Unix())
	key = fmt.Sprintf("%s:%d", key, rand.Intn(HotKeyNum))
	b, err := redis.Bytes(r.Do("GET", key))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, result)
}

// SetJSON 写入json数据
func (r *Redis) SetJSON(key string, result interface{}, expire int64) error {
	if ostring.InterfaceIsNil(result) {
		return errors.New("set result is nil")
	}

	b, err := json.Marshal(result)
	if err != nil {
		return err
	}
	_, err = r.Do("SET", key, string(b), "EX", expire)
	return err
}

// SetHotKeyJSON 写入热key json数据
func (r *Redis) SetHotKeyJSON(key string, result interface{}, expire int64) error {
	if ostring.InterfaceIsNil(result) {
		return errors.New("set result is nil")
	}

	b, err := json.Marshal(result)
	if err != nil {
		return err
	}
	for i := 0; i < HotKeyNum; i++ {
		key := fmt.Sprintf("%s:%d", key, i)
		_, err = r.Do("SET", key, string(b), "EX", expire)
		if err != nil {
			return err
		}
	}
	return err
}

// DelKey 删除热key
func (r *Redis) DelKey(key string) error {
	keys := make([]interface{}, HotKeyNum)
	for i := 0; i < HotKeyNum; i++ {
		key := fmt.Sprintf("%s:%d", key, i)
		keys[i] = key
	}
	_, err := r.Do("DEL", keys...)
	return err
}

// GetAndSetJSON 查询缓存json数据并返回写入方法
func (r *Redis) GetAndSetJSON(key string, result interface{}) (error, func(cacheResult interface{}, expire int64) error) {
	if err := r.GetJSON(key, result); err != nil {
		return err, func(cacheResult interface{}, expire int64) error {
			return r.SetJSON(key, cacheResult, expire)
		}
	}
	return nil, nil
}

// GetAndSetHotKeyJSON 查询热key缓存json数据并返回写入方法
func (r *Redis) GetAndSetHotKeyJSON(key string, result interface{}) (error, func(cacheResult interface{}, expire int64) error) {
	if err := r.GetHotKeyJSON(key, result); err != nil {
		return err, func(cacheResult interface{}, expire int64) error {
			return r.SetHotKeyJSON(key, cacheResult, expire)
		}
	}
	return nil, nil
}
