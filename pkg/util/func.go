package util

import "os"

var defaultGoEnv = "production"

// GetGoEnv 查询golang运行环境,默认生产
func GetGoEnv() string {
	goEnv := os.Getenv("GO_ENV")
	if goEnv == "" {
		goEnv = defaultGoEnv
	}
	return goEnv
}
