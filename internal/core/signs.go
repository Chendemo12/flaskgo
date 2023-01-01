// Package core 内部标志量
package core

import "time"

const (
	HotSwitchSigint = 30 // 热调试开关
	ShutdownSigint  = 1  // 关机信号
)

var (
	BaseRoutesDisabled       = false            // 禁用基础路由
	SwaggerDisabled          = false            // 禁用文档自动生成
	RequestValidateDisabled  = true             // 禁用请求体自动验证
	ResponseValidateDisabled = false            // 禁用返回体自动验证
	PrintDisabled            = false            // 禁用默认输出
	MultipleProcessDisabled  = true             // 禁用多进程
	ShutdownWithTimeout      = 20 * time.Second // 关机前的等待时间
)
