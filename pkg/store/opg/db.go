package opg

import (
	"fmt"
	"github.com/kratos-tools/lib/pkg/util"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	klog "github.com/silenceper/log"
)

type (
	DB = gorm.DB
	// Scope ...
	Scope = gorm.Scope
)

// Open 连接数据库
func Open(c *Config) (*DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=%s dbname=%s",
		c.Host, c.Port, c.User, c.Password, c.SslMode, c.Db,
	)
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if util.GetGoEnv() == "development" {
		db.LogMode(true)
	} else {
		db.LogMode(c.Debug)
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	if c.ConnMaxLifetime != 0 {
		db.DB().SetConnMaxLifetime(c.ConnMaxLifetime)
	}

	klog.Infof("pg: %s connect: %s MaxIdleConns=%d MaxOpenConns=%d ConnMaxLifetime=%ds", c.Name, connStr, c.MaxIdleConns, c.MaxOpenConns, c.ConnMaxLifetime/time.Second)

	if len(c.UpdateColumn) > 0 {
		// 更新时间
		db.Callback().Update().Replace("gorm:update_time_stamp", func(scope *Scope) {
			if _, ok := scope.Get("gorm:update_column"); !ok {
				now := time.Now()
				for _, v := range c.UpdateColumn {
					err = scope.SetColumn(v, &now)
				}
			}
		})
	}

	// 修复回表时PK健不区分大小写的问题
	db.Callback().Create().Replace("gorm:force_reload_after_create", func(scope *gorm.Scope) {
		if blankColumnsWithDefaultValue, ok := scope.InstanceGet("gorm:blank_columns_with_default_value"); ok {
			db := scope.DB().New().Table(scope.TableName()).Select(blankColumnsWithDefaultValue.([]string))
			for _, field := range scope.Fields() {
				if field.IsPrimaryKey && !field.IsBlank {
					db = db.Where(fmt.Sprintf("%v = ?", scope.Dialect().Quote(field.DBName)), field.Field.Interface())
				}
			}
			db.Scan(scope.Value)
		}
	})
	return db, nil
}
