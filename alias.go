// Package flaskgo 是一个基于fiber封装了常用方法的软件包
//
// 其提供了类似于FastAPI的API设计，并提供了接口文档自动生成、请求体自动校验和返回值自动序列化等使用功能；
package flaskgo

import (
	"github.com/Chendemo12/flaskgo/internal/app"
	"github.com/Chendemo12/flaskgo/internal/core"
	"github.com/Chendemo12/flaskgo/internal/godantic"
	"github.com/Chendemo12/flaskgo/internal/mode"
	"time"
)

//goland:noinspection GoUnusedGlobalVariable
const (
	DevMode        = mode.DevMode
	ProdMode       = mode.ProdMode
	RouteSeparator = app.RouteSeparator
	Version        = app.Version
)

//goland:noinspection GoUnusedGlobalVariable
var (
	// Deprecated: NewDefaultFlaskGo use NewFlaskGo instead.
	NewDefaultFlaskGo = app.NewFlaskGo
	NewFlaskGo        = app.NewFlaskGo
	APIRouter         = app.APIRouter
	CombinePath       = app.CombinePath
)

//goland:noinspection GoUnusedGlobalVariable
var (
	IsDebug = mode.IsDebug
	SetMode = mode.SetMode
	GetMode = mode.GetMode
)

type Field = godantic.Field
type BaseModel = godantic.BaseModel
type BaseModelIface = godantic.Iface

//goland:noinspection GoUnusedGlobalVariable
//var ( // types
//	Int8    = openapi.Int8
//	Int16   = openapi.Int16
//	Int32   = openapi.Int32
//	Int64   = openapi.Int64
//	Uint8   = openapi.Uint8
//	Uint16  = openapi.Uint16
//	Uint32  = openapi.Uint32
//	Uint64  = openapi.Uint64
//	Float32 = openapi.Float32
//	Float64 = openapi.Float64
//	String  = openapi.String
//	Boolean = openapi.Boolean
//	Bool    = openapi.Boolean
//	Mapping = openapi.Mapping
//
//	Int     = Int32
//	Byte    = Uint8
//	Uint    = Uint32
//	Float   = Float64
//	Array   = openapi.List
//	List    = openapi.List
//	Ints    = &openapi.RouteModel{Model: Int32, Struct: Int32, RetArray: true}
//	Bytes   = &openapi.RouteModel{Model: Uint8, Struct: Uint8, RetArray: true}
//	Strings = &openapi.RouteModel{Model: String, Struct: String, RetArray: true}
//	Floats  = &openapi.RouteModel{Model: Float64, Struct: Float64, RetArray: true}
//)

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

//goland:noinspection GoUnusedGlobalVariable
//var (
//	ValidationErrorResponse = app.ValidationErrorResponse
//	AnyResponse             = app.AnyResponse
//	JSONResponse            = app.JSONResponse
//	StringResponse          = app.StringResponse
//	StreamResponse          = app.StreamResponse
//	FileResponse            = app.FileResponse
//	ErrorResponse           = app.ErrorResponse
//	HTMLResponse            = app.HTMLResponse
//	OKResponse              = app.OKResponse
//	ResourceNotFound        = app.ResourceNotFound
//	AdvancedResponse        = app.AdvancedResponse
//	StringsReverse          = openapi.StringsReverse
//)

// DisableBaseRoutes 禁用基础路由
func DisableBaseRoutes() { core.BaseRoutesDisabled = true }

// DisableSwagAutoCreate 禁用文档自动生成
func DisableSwagAutoCreate() { core.SwaggerDisabled = true }

// DisableDefaultOutput 禁用默认输出
func DisableDefaultOutput() { core.PrintDisabled = true }

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
