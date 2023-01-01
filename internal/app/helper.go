package app

import (
	"gitlab.cowave.com/gogo/flaskgo/internal/core"
	"gitlab.cowave.com/gogo/flaskgo/internal/mode"
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
