package app

import (
	"errors"
	"fmt"
	"github.com/Chendemo12/flaskgo/internal/core"
	"github.com/gofiber/fiber/v2"
	"os"
	"runtime"
	"strconv"
	"time"
)

// AcquireCtx 申请一个 Context 并初始化
func (f *FlaskGo) AcquireCtx(fctx *fiber.Ctx) *Context {
	c := f.pool.Get().(*Context)
	// 初始化各种参数
	c.ec = fctx
	c.RequestBody = int64(1) // 初始化为1，避免访问错误
	c.PathFields = map[string]string{}
	c.QueryFields = map[string]string{}
	return c
}

// ReleaseCtx 释放并归还 Context
func (f *FlaskGo) ReleaseCtx(ctx *Context) {
	ctx.ec = nil
	ctx.RequestBody = int64(1)
	ctx.PathFields = nil
	ctx.QueryFields = nil

	f.pool.Put(ctx)
}

// ReplaceErrorHandler 替换fiber错误处理方法，是 请求错误处理方法
func (f *FlaskGo) ReplaceErrorHandler(fc fiber.ErrorHandler) *FlaskGo {
	fiberErrorHandler = fc
	return f
}

// ReplaceStackTraceHandler 替换错误堆栈处理函数，即 recover 方法
func (f *FlaskGo) ReplaceStackTraceHandler(fc StackTraceHandlerFunc) *FlaskGo {
	recoverHandler = fc
	return f
}

// ReplaceRecover 重写全局 recover 方法
func (f *FlaskGo) ReplaceRecover(fc StackTraceHandlerFunc) *FlaskGo {
	return f.ReplaceStackTraceHandler(fc)
}

// AddResponseHeader 添加一个响应头
//
//	@param	key		string	键
//	@param	value	string	值
func (f *FlaskGo) AddResponseHeader(key, value string) *FlaskGo {
	// 首先判定是否已经存在
	for i := 0; i < len(responseHeaders); i++ {
		if responseHeaders[i].Key == key {
			responseHeaders[i].Value = value
			return f
		}
	}
	// 不存在，新建
	responseHeaders = append(responseHeaders, &ResponseHeader{
		Key:   key,
		Value: value,
	})
	return f
}

// DeleteResponseHeader 删除一个响应头
//
//	@param	key	string	键
func (f *FlaskGo) DeleteResponseHeader(key string) *FlaskGo {
	for i := 0; i < len(responseHeaders); i++ {
		if responseHeaders[i].Key == key {
			responseHeaders[i].Value = ""
		}
	}
	return f
}

// Deprecated: RunCronjob 启动定时任务, 此函数内部通过创建一个协程来执行任务，并且阻塞至 FlaskGo 完成初始化
//
//	@param	tasker	func(service CustomContextIface)	error	定时任务
//	@param	service	CustomContextIface					服务依赖
func (f *FlaskGo) RunCronjob(_ func(ctx *Service) error) *FlaskGo {
	return f
}

// SetShutdownTimeout 修改关机前最大等待时间
//
//	@param	timeout	in	修改关机前最大等待时间,	单位秒
func (f *FlaskGo) SetShutdownTimeout(timeout int) *FlaskGo {
	core.ShutdownWithTimeout = time.Duration(timeout)
	return f
}

// DisableBaseRoutes 禁用基础路由
func (f *FlaskGo) DisableBaseRoutes() *FlaskGo {
	core.BaseRoutesDisabled = true
	return f
}

// DisableSwagAutoCreate 禁用文档自动生成
func (f *FlaskGo) DisableSwagAutoCreate() *FlaskGo {
	core.SwaggerDisabled = true
	return f
}

// DisableRequestValidate 禁用请求体自动验证
func (f *FlaskGo) DisableRequestValidate() *FlaskGo {
	core.RequestValidateDisabled = true
	return f
}

// DisableResponseValidate 禁用返回体自动验证
func (f *FlaskGo) DisableResponseValidate() *FlaskGo {
	core.ResponseValidateDisabled = true
	return f
}

// DisableMultipleProcess 禁用多进程
func (f *FlaskGo) DisableMultipleProcess() *FlaskGo {
	core.MultipleProcessDisabled = true
	return f
}

// ShutdownWithTimeout 关机前最大等待时间
func (f *FlaskGo) ShutdownWithTimeout() time.Duration {
	return core.ShutdownWithTimeout * time.Second
}

// EnableDumpPID 启用PID存储
func (f *FlaskGo) EnableDumpPID() *FlaskGo {
	core.DumpPIDEnabled = true
	return f
}

// DumpPID 存储PID, 文件权限为0775
// 对于 Windows 其文件为当前运行目录下的pid.txt;
// 对于 类unix系统，其文件为/run/{Title}.pid, 若无读写权限则改为当前运行目录下的pid.txt;
func (f *FlaskGo) DumpPID() {
	var path string
	switch runtime.GOOS {
	case "darwin", "linux":
		path = fmt.Sprintf("/run/%s.pid", f.title)
	case "windows":
		path = "pid.txt"
	}

	pid := []byte(strconv.Itoa(f.PID()))
	err := os.WriteFile(path, pid, 755)
	if errors.Is(err, os.ErrPermission) {
		_ = os.WriteFile("pid.txt", pid, 0775)
	}
}
