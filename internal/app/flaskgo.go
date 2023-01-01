// Package app
//
// hooks.go: FlaskGo 路由中间件, 参数校验和返回值序列化方法;
//
// response.go: 返回体;
//
// svc: FlaskGo 全局服务依赖 Service, Context，自定义服务上下文需实现 ServiceContextIface 接口
//
// router.go: FlaskGo 路由组;
package app

import (
	"github.com/Chendemo12/flaskgo/internal/core"
	"github.com/Chendemo12/flaskgo/internal/mode"
	"github.com/Chendemo12/flaskgo/internal/swag"
	"github.com/Chendemo12/functools/python"
	"github.com/Chendemo12/functools/zaplog"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type EventKind string

const (
	startupEvent  EventKind = "startup"
	shutdownEvent EventKind = "shutdown"
)

var (
	once               = sync.Once{}
	console            = zaplog.ConsoleLogger{}
	appEngine *FlaskGo = nil // 单例模式
)

type Dict = map[string]any

type Event struct {
	Fc   func()
	Type EventKind // 事件类型：startup 或 shutdown
}

type FlaskGo struct {
	console     zaplog.ConsoleLogger // 控制台日志
	logger      zaplog.Iface         // 日志对象，通常=Sugar(*zap.SugaredLogger已实现此接口)
	isStarted   chan struct{}        // 标记程序是否完成启动
	service     *Service             // 全局服务依赖
	engine      *fiber.App           // fiber.App
	version     string
	host        string
	port        string
	description string
	title       string
	routers     []*Router // FlaskGo 路由组 Router
	events      []*Event
	middlewares []any
}

// Title 应用程序名和日志文件名
func (f *FlaskGo) Title() string { return f.title }

// Version 版本号
func (f *FlaskGo) Version() string { return f.version }

// Description 描述信息，同时会显示在Swagger文档上
func (f *FlaskGo) Description() string { return f.description }
func (f *FlaskGo) Host() string        { return f.host }
func (f *FlaskGo) Port() string        { return f.port }

// mountBaseRoutes 创建基础路由
func (f *FlaskGo) mountBaseRoutes() {
	// 注册最基础的路由
	router := APIRouter("/api/base", []string{"Base"})
	{
		router.GET("/title", swag.String, "获取软件名", func(c *Context) any {
			return StringResponse(appEngine.title)
		})
		router.GET("/Description", swag.String, "获取软件描述信息", func(c *Context) any {
			return StringResponse(appEngine.Description())
		})
		router.GET("/version", swag.String, "获取软件版本号", func(c *Context) any {
			return StringResponse(appEngine.version)
		})
		router.GET("/heartbeat", swag.String, "心跳检测", func(c *Context) any {
			return StringResponse("pong")
		})
		router.GET("/debug", DebugMode{}, "获取调试开关", func(c *Context) any {
			return OKResponse(DebugMode{Mode: mode.GetMode()})
		})
	}
	f.routers = append(f.routers, router)
}

// mountUserRoutes 挂载并记录自定义路由
func (f *FlaskGo) mountUserRoutes() {
	for _, router := range f.routers {
		rtr := f.engine.Group(router.Prefix)
		for _, route := range router.Routes() {
			switch route.Method {
			case http.MethodGet:
				rtr.Get(route.RelativePath, route.Handlers...)
				// 记录自定义路由
				MethodGetRoutes[route.Path(router.Prefix)] = route

			case http.MethodPost:
				rtr.Post(route.RelativePath, route.Handlers...)
				MethodPostRoutes[route.Path(router.Prefix)] = route

			case http.MethodDelete:
				rtr.Delete(route.RelativePath, route.Handlers...)
				MethodDeleteRoutes[route.Path(router.Prefix)] = route

			case http.MethodPatch:
				rtr.Patch(route.RelativePath, route.Handlers...)
				MethodPatchRoutes[route.Path(router.Prefix)] = route

			case http.MethodPut:
				rtr.Put(route.RelativePath, route.Handlers...)
				MethodPutRoutes[route.Path(router.Prefix)] = route

			case "ANY", "ALL":
				rtr.All(route.RelativePath, route.Handlers...)
				MethodGetRoutes[route.Path(router.Prefix)] = route
			}
		}
	}
}

// initialize 初始化FlaskGo,并完成服务依赖的建立
// FlaskGo启动前，必须显式的初始化FlaskGo的基本配置，若初始化中发生异常则panic
//  1. 记录工作地址： host:Port
//  2. 创建fiber.App createFiberApp
//  3. 挂载中间件
//  4. 按需挂载基础路由 mountBaseRoutes
//  5. 挂载自定义路由 mountUserRoutes
//  6. 初始化日志logger logger.NewLogger
//  7. 安装创建swagger文档 makeSwaggerDocs
func (f *FlaskGo) initialize() *FlaskGo {
	f.console.SDebug("Run mode: " + mode.GetMode())

	// 创建 fiber.App
	f.engine = createFiberApp(f.title, f.version)
	// 注册中间件
	for i := 0; i < len(f.middlewares); i++ {
		f.engine.Use(f.middlewares[i])
	}

	//挂载基础路由
	if python.Any(mode.IsDebug(), !core.BaseRoutesDisabled) {
		f.mountBaseRoutes()
	}

	// 挂载自定义路由
	f.mountUserRoutes()

	// 配置日志
	if f.logger == nil {
		lc := &zaplog.Config{
			Filename:   f.title,
			Level:      0,
			Rotation:   2,
			Retention:  7,
			MaxBackups: 5,
			Compress:   true,
		}
		if mode.IsDebug() {
			lc.Level = zaplog.DEBUG
		} else {
			lc.Level = zaplog.WARNING
		}

		f.logger = zaplog.NewLogger(lc).Sugar()
		innerOutput("DEBUG", "Logger initialized.")
	}

	// 创建swag文档, 必须等上层注册完路由之后才能调用
	makeSwaggerDocs(f)

	return f
}

// Service 获取FlaskGo全局服务上下文
func (f *FlaskGo) Service() *Service { return f.service }

// CustomServiceContext 获取上层自定义服务依赖
func (f *FlaskGo) CustomServiceContext() CustomContextIface { return f.service.ctx }

// APIRouters 获取全部注册的路由组
// @return  []*Router 路由组
func (f *FlaskGo) APIRouters() []*Router { return f.routers }

// FiberApp 获取fiber引擎
// @return  *fiber.App fiber引擎
func (f *FlaskGo) FiberApp() *fiber.App { return f.engine }

// OnEvent 添加启动事件
// @param  kind  事件类型，取值需为  "startup"  /  "shutdown"
// @param  fs    func()     事件
func (f *FlaskGo) OnEvent(kind EventKind, fc func()) *FlaskGo {
	switch kind {
	case startupEvent:
		f.events = append(f.events, &Event{
			Type: startupEvent,
			Fc:   fc,
		})
	case shutdownEvent:
		f.events = append(f.events, &Event{
			Type: shutdownEvent,
			Fc:   fc,
		})
	default:
	}
	return f
}

// AddResponseHeader 添加一个响应头
// @param  key    string  键
// @param  value  string  值
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
// @param  key  string  键
func (f *FlaskGo) DeleteResponseHeader(key string) *FlaskGo {
	for i := 0; i < len(responseHeaders); i++ {
		if responseHeaders[i].Key == key {
			responseHeaders[i].Value = ""
		}
	}
	return f
}

// ReplaceCtx 替换自定义服务上下文
// @param  service  CustomContextIface  服务上下文
func (f *FlaskGo) ReplaceCtx(ctx CustomContextIface) *FlaskGo {
	f.service.SetServiceContext(ctx)
	return f
}

// ReplaceLogger 替换日志句柄，此操作必须在run之前进行
// @param  logger  LoggerIface  日志句柄
func (f *FlaskGo) ReplaceLogger(logger zaplog.Iface) *FlaskGo {
	f.logger = logger
	return f
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

// SetDescription 设置APP的详细描述信息
// @param  Description  string  详细描述信息
func (f *FlaskGo) SetDescription(description string) *FlaskGo {
	f.description = description
	return f
}

// IncludeRouter 注册一个路由组
// @param  router  *Router  路由组
func (f *FlaskGo) IncludeRouter(router *Router) *FlaskGo {
	f.routers = append(f.routers, router)
	return f
}

// Use 添加中间件
func (f *FlaskGo) Use(middleware ...any) *FlaskGo {
	f.middlewares = append(f.middlewares, middleware...)
	return f
}

// ActivateHotSwitch 创建一个热开关，监听信号量30，用来改变程序调试开关状态
func (f *FlaskGo) ActivateHotSwitch() *FlaskGo {
	swt := make(chan os.Signal)
	signal.Notify(swt, syscall.Signal(core.HotSwitchSigint))

	go func() {
		for range swt {
			if mode.IsDebug() {
				resetRunMode(mode.ProdMode)
			} else {
				resetRunMode(mode.DevMode)
			}
		}
	}()

	return f
}

// RunCronjob 启动定时任务, 此函数内部通过创建一个协程来执行任务，并且阻塞至flaskgo完成初始化
// @param  tasker   func(service CustomContextIface)  error  定时任务
// @param  service  CustomContextIface                服务依赖
func (f *FlaskGo) RunCronjob(tasker func(ctx *Service) error) *FlaskGo {
	go func() {
		defer close(f.isStarted)
		for {
			select {
			case <-f.isStarted:
				if err := tasker(f.Service()); err != nil {
					f.console.Error("cronjob error: ", err.Error())
				} else {
					innerOutput("DEBUG", "Tasker started.")
				}
				return
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
	return f
}

// serve 初始化服务
func (f *FlaskGo) serve() *FlaskGo {
	f.initialize().ActivateHotSwitch()

	// 执行启动前事件
	for _, event := range f.events {
		if event.Type == startupEvent {
			event.Fc()
		}
	}

	f.isStarted <- struct{}{} // 解除阻塞上层的任务
	f.console.SInfo("HTTP server listening on: " + net.JoinHostPort(f.host, f.port))

	return f
}

// Run 启动服务, 此方法会阻塞运行，因此必须放在main函数结尾
func (f *FlaskGo) Run(host, port string) {
	if !fiber.IsChild() {
		f.host = host
		f.port = port
		f.serve()
	}
	// 关闭开关
	quit := make(chan os.Signal, core.ShutdownSigint)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Fatal(f.engine.Listen(net.JoinHostPort(f.host, f.port)))
	}()

	<-quit // 阻塞进程，直到接收到停止信号,准备关闭程序
	// 执行关机前事件
	for _, event := range f.events {
		if event.Type == shutdownEvent {
			event.Fc()
		}
	}

	_ = f.logger.Sync()
	// TODO：implement 平滑关机
	innerOutput("INFO", "Server exiting")
}

func GetFlaskGo() *FlaskGo { return appEngine }

// NewFlaskGo 创建一个WEB服务
// @param   title    string              Application  name
// @param   version  string              Version
// @param   debug    bool                是否开启调试模式
// @param   service  CustomContextIface  custom  ServiceContext
// @return  *FlaskGo flaskgo对象
func NewFlaskGo(title, version string, debug bool, ctx CustomContextIface) *FlaskGo {
	if debug {
		mode.SetMode(mode.DevMode)
	} else {
		mode.SetMode(mode.ProdMode)
	}

	once.Do(func() {
		appEngine = &FlaskGo{
			title:       title,
			version:     version,
			console:     console,
			description: title + " Micro Context",
			service:     &Service{ctx: ctx, validate: validator.New()},
			isStarted:   make(chan struct{}, 1),
			middlewares: make([]any, 0),
			events:      make([]*Event, 0),
		}
		appEngine.service.app = appEngine
	})

	return appEngine
}
