package odbcache

import (
	"errors"
	"sort"
	"strings"

	"github.com/fatih/structs"
	"github.com/kratos-tools/lib/pkg/util/ostring"
)

type RegisterAction int64

const (
	RegisterPut = RegisterAction(iota)
	RegisterDel
)

// paramsDecode 查询参数key序列化
func (t *TableCache) paramsDecode(params map[string]interface{}) string {
	var param []string
	var result strings.Builder

	for k := range params {
		param = append(param, k)
	}

	sort.Strings(param)
	for k, v := range param {
		if k != 0 {
			result.WriteString(",")
		}
		result.WriteString(v)
	}
	return result.String()
}

// paramsEncode 参数解析
func (t *TableCache) paramsEncode(key string, s interface{}) map[string]interface{} {
	data := structs.Map(s)
	result := map[string]interface{}{}
	params := strings.Split(key, ",")
	for _, key := range params {
		result[key] = data[t.Column[key]]
	}
	return result
}

func (t *TableCache) registerKey() string {
	return t.key("register")
}

// register 注册要清空的数据
func (t *TableCache) register(registers []*TableCacheRegisterClear) map[string]RegisterAction {
	registerAction := map[string]RegisterAction{}
	for _, register := range registers {
		registerAction[t.paramsDecode(register.Params)] = register.Action
	}
	return registerAction
}

// Clear 根据注册的数据清空现有的缓存
func (t *TableCache) Clear(data interface{}) (err error) {
	if ostring.InterfaceIsNil(data) {
		return errors.New(t.TableName + " clear data is nil")
	}
	for params, action := range t.RegisterClear {
		k := string(ostring.MustJson(t.paramsEncode(params, data), false))
		switch action {
		case RegisterDel:
			if err = t.Del(k); err != nil {
				return
			}
		case RegisterPut:
			if err = t.Put(k, data); err != nil {
				return
			}
		}
	}
	return nil
}
