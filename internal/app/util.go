package app

import (
	"github.com/Chendemo12/flaskgo/internal/constant"
	"github.com/Chendemo12/flaskgo/internal/core"
	"github.com/Chendemo12/flaskgo/internal/mode"
	"strings"
)

// innerOutput 允许上层禁用FlaskGo的输出，但开启dev模式时则忽略禁用开关
func innerOutput(level, message string) {
	if !mode.IsDebug() && core.PrintDisabled {
		return
	}

	switch level {
	case "DEBUG":
		console.SDebug(message)
	case "INFO":
		console.SInfo(message)
	case "WARN":
		console.SWarn(message)
	case "ERROR":
		console.SError(message)
	default:
		console.SInfo(message)
	}
}

// resetRunMode 重设运行时环境
// @param  md  string  开发环境
func resetRunMode(md string) {
}

// DoesPathParamsFound 是否查找到路径参数
// @param  path  string  路由
func DoesPathParamsFound(path string) (map[string]bool, bool) {
	pathParameters := make(map[string]bool, 0)
	// 查找路径中的参数
	for _, p := range strings.Split(path, constant.PathSeparator) {
		if strings.HasPrefix(p, constant.PathParamPrefix) {
			// 识别到路径参数
			if strings.HasSuffix(p, constant.OptionalPathParamSuffix) {
				// 可选路径参数
				pathParameters[p[1:len(p)-1]] = false
			} else {
				pathParameters[p[1:]] = true
			}
		}
	}
	return pathParameters, len(pathParameters) > 0
}
