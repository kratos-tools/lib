package odbcache

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gomodule/redigo/redis"
	"github.com/kratos-tools/lib/pkg/store/oredis"
	"github.com/kratos-tools/lib/pkg/util/ostring"
)

type Table interface {
	TableName() string
	PKField() string
	PKType() string
	PKName() string
	GetColumnByKey() map[string]string
}

func NewTableCache(client *oredis.Redis, table Table, lifeSpan int64, RegisterClear []*TableCacheRegisterClear) *TableCache {
	cache := &TableCache{
		Redis:     client,
		TableName: table.TableName(),
		PkName:    table.PKName(),
		PkType:    table.PKType(),
		PKField:   table.PKField(),
		Column:    table.GetColumnByKey(),
		LifeSpan:  lifeSpan,
	}

	cache.RegisterClear = cache.register(RegisterClear)
	return cache
}

// TableCacheRegisterClear 注册clear
type TableCacheRegisterClear struct {
	Params map[string]interface{}
	Action RegisterAction
}

type TableCache struct {
	*oredis.Redis
	TableName     string
	PkName        string
	PkType        string
	PKField       string
	Column        map[string]string
	LifeSpan      int64
	RegisterClear map[string]RegisterAction
}

func (t *TableCache) key(k string) string {
	return fmt.Sprintf("%s:%s:%s", t.GetConfig().Name, t.TableName, k)
}

// parseValue 解析存储的value
func (t *TableCache) parseValue(v []byte) interface{} {
	var result interface{}
	if err := json.Unmarshal(v, &result); err != nil {
		return nil
	}

	switch val := result.(type) {
	//case int: // 如果查询的类型是 map查询条件之后 总的条目数量
	//	fmt.Println("parseValue int", string(v))
	//	return val
	case int64: // 如果查询的类型是 map查询条件之后 总的条目数量
		var count int64
		if err := json.Unmarshal(v, &count); err != nil {
			log.Errorf("TableCache parseValue int64 err:%s v:%s", err, string(v))
			return nil
		}
		return count
	case []int64: // 如果查询的类型是 map查询条件之后 条目的主键id
		ids := make([]int64, 0)
		if err := json.Unmarshal(v, &ids); err != nil {
			log.Errorf("TableCache parseValue []int64 err:%s v:%s", err, string(v))
			return nil
		}
		return ids
	case []string: // 如果查询的类型是 map查询条件之后 条目的主键id
		ids := make([]string, 0)
		if err := json.Unmarshal(v, &ids); err != nil {
			log.Errorf("TableCache parseValue []string err:%s v:%s", err, string(v))
			return nil
		}
		return ids
	case []interface{}: // 如果查询的类型是 map查询条件之后 条目的主键id
		switch t.PkType {
		case "string":
			ids := make([]string, 0)
			if err := json.Unmarshal(v, &ids); err != nil {
				log.Errorf("TableCache parseValue []interface{} []string err:%s v:%s", err, string(v))
				return nil
			}
			return ids
		case "int64":
			ids := make([]int64, 0)
			if err := json.Unmarshal(v, &ids); err != nil {
				log.Errorf("TableCache parseValue []interface{} []int64 err:%s v:%s", err, string(v))
				return nil
			}
			return ids
		}
		return val
	default: // 如果查询的类型是 数据详情
		fmt.Println("parseValue interface", string(v))
		return val
	}
}

// Load 加载缓存结果
func (t *TableCache) Load(k string) interface{} {
	var err error
	var value []byte
	if value, err = redis.Bytes(t.Do("GET", t.key(k))); err != nil {
		return nil
	}

	return t.parseValue(value)
}

// LoadDetail 加载映射详情
func (t *TableCache) LoadDetail(k string, v interface{}) interface{} {
	var err error
	var value []byte
	if value, err = redis.Bytes(t.Do("GET", t.key(k))); err != nil {
		return nil
	}

	if err := json.Unmarshal(value, &v); err != nil {
		return nil
	}
	return v
}

// Put 重新设置缓存
func (t *TableCache) Put(k string, o interface{}) error {
	value := string(ostring.MustJson(o, false))
	if ostring.InterfaceIsNil(o) {
		return fmt.Errorf(" value is nil")
	}
	if _, err := t.Do("SET", t.key(k), value, "EX", t.LifeSpan); err != nil {
		return err
	}
	return nil
}

// Del 移除缓存
func (t *TableCache) Del(k string) error {
	_, err := t.Do("DEL", t.key(k))
	return err
}

// Exists 是否存在
func (t *TableCache) Exists(k string) bool {
	exists, err := t.Do("EXISTS", t.key(k))
	if err != nil {
		return false
	} else {
		return exists.(int64) == 1
	}
}

// OffSetParams 查询参数
type OffSetParams struct {
	Params  map[string]interface{} // 查询参数
	Raw     string                 // where参数 写 > < 参数
	PKOrder int64                  // 0 不走主键索引排序, -1 降序 +1 升序
	Order   string                 // 其他字段排序
	Offset  int64                  // 跳过offset个
	Limit   int64                  // 限制数量
}
