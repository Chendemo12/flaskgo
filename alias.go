// Package flaskgo 是一个基于fiber封装了常用方法的软件包
//
// 其提供了类似于FastAPI的API设计，并提供了接口文档自动生成、请求体自动校验和返回值自动序列化等使用功能；
package flaskgo

import (
	"github.com/Chendemo12/flaskgo/internal/app"
	"github.com/Chendemo12/flaskgo/internal/core"
	"github.com/Chendemo12/flaskgo/internal/godantic"
	"time"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	// Deprecated: NewDefaultFlaskGo use NewFlaskGo instead.
	NewDefaultFlaskGo = app.NewFlaskGo
	NewFlaskGo        = app.NewFlaskGo
	APIRouter         = app.APIRouter
	CombinePath       = app.CombinePath
)

type Field = godantic.Field
type BaseModel = godantic.BaseModel
type BaseModelIface = godantic.Iface

//goland:noinspection GoUnusedGlobalVariable
var ( // types
	String = godantic.String
	Str    = godantic.String

	Bool    = godantic.Bool
	Boolean = godantic.Bool

	Int    = godantic.Int
	Byte   = godantic.Uint8
	Int8   = godantic.Int8
	Int16  = godantic.Int16
	Int32  = godantic.Int32
	Int64  = godantic.Int64
	Uint8  = godantic.Uint8
	Uint16 = godantic.Uint16
	Uint32 = godantic.Uint32
	Uint64 = godantic.Uint64
	// Float32 = openapi.Float32
	// Float64 = openapi.Float64

	// Mapping = openapi.Mapping

	// Float   = Float64
	// Array   = openapi.List
	// List    = openapi.List
	// Ints    = &openapi.RouteModel{Model: Int32, Struct: Int32, RetArray: true}
	// Bytes   = &openapi.RouteModel{Model: Uint8, Struct: Uint8, RetArray: true}
	// Strings = &openapi.RouteModel{Model: String, Struct: String, RetArray: true}
	// Floats  = &openapi.RouteModel{Model: Float64, Struct: Float64, RetArray: true}
)

//goland:noinspection GoUnusedGlobalVariable
type Dict = map[string]any // python dict
type SDict = map[string]string
type H = map[string]any // gin.H
type M = map[string]any // fiber.M
type SM = map[string]string
type Map = map[string]any // python map
type ServiceContextIface = app.CustomContextIface
type CustomContextIface = app.CustomContextIface
type ServiceContext = app.Context
type Context = app.Context
type Service = app.Service

//goland:noinspection GoUnusedGlobalVariable
type FlaskGo = app.FlaskGo
type HandlerFunc = app.HandlerFunc
type StackTraceHandlerFunc = app.StackTraceHandlerFunc
type Route = app.Route
type Router = app.Router
type Response = app.Response
type ResponseHeader = app.ResponseHeader
type ValidationError = app.ValidationError

type CronJob = app.CronJob
type Scheduler = app.Scheduler
type CronJobFunc = app.CronJobFunc

//goland:noinspection GoUnusedGlobalVariable
var (
	OKResponse              = app.OKResponse
	JSONResponse            = app.JSONResponse
	ErrorResponse           = app.ErrorResponse
	ValidationErrorResponse = app.ValidationErrorResponse
	AnyResponse             = app.AnyResponse
	StringResponse          = app.StringResponse
	StreamResponse          = app.StreamResponse
	FileResponse            = app.FileResponse
	HTMLResponse            = app.HTMLResponse
	AdvancedResponse        = app.AdvancedResponse
)

// DisableBaseRoutes 禁用基础路由
func DisableBaseRoutes() { core.BaseRoutesDisabled = true }

// DisableSwagAutoCreate 禁用文档自动生成
func DisableSwagAutoCreate() { core.SwaggerDisabled = true }

// DisableRequestValidate 禁用请求体自动验证
func DisableRequestValidate() { core.RequestValidateDisabled = true }

// DisableResponseValidate 禁用返回体自动验证
func DisableResponseValidate() { core.ResponseValidateDisabled = true }

// DisableMultipleProcess 禁用多进程
func DisableMultipleProcess() { core.MultipleProcessDisabled = true }

// ShutdownWithTimeout 指定关机前最大等待时间
func ShutdownWithTimeout() time.Duration { return core.ShutdownWithTimeout * time.Second }

// SetShutdownTimeout 修改关机前最大等待时间
// @param  timeout  in  修改关机前最大等待时间,  单位秒
func SetShutdownTimeout(timeout int) { core.ShutdownWithTimeout = time.Duration(timeout) }
