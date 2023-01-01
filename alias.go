// Package flaskgo 是一个基于fiber封装了常用方法的软件包
//
// 其提供了类似于FastAPI的API设计，并提供了接口文档自动生成、请求体自动校验和返回值自动序列化等使用功能；
package flaskgo

import (
	"gitlab.cowave.com/gogo/flaskgo/internal/app"
	"gitlab.cowave.com/gogo/flaskgo/internal/core"
	"gitlab.cowave.com/gogo/flaskgo/internal/mode"
	"gitlab.cowave.com/gogo/flaskgo/internal/swag"
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
	NewDefaultFlaskGo = app.NewFlaskGo
	NewFlaskGo        = app.NewFlaskGo
	GetFlaskGo        = app.GetFlaskGo
	APIRouter         = app.APIRouter
	CombinePath       = app.CombinePath
)

//goland:noinspection GoUnusedGlobalVariable
var (
	IsDebug = mode.IsDebug
	SetMode = mode.SetMode
	GetMode = mode.GetMode
)

//goland:noinspection GoUnusedGlobalVariable
var ( // types
	Int8    = swag.Int8
	Int16   = swag.Int16
	Int32   = swag.Int32
	Int64   = swag.Int64
	Uint8   = swag.Uint8
	Uint16  = swag.Uint16
	Uint32  = swag.Uint32
	Uint64  = swag.Uint64
	Float32 = swag.Float32
	Float64 = swag.Float64
	String  = swag.String
	Boolean = swag.Boolean
	Bool    = swag.Boolean
	Mapping = swag.Mapping

	Int     = Int32
	Byte    = Uint8
	Uint    = Uint32
	Float   = Float64
	Array   = swag.List
	List    = swag.List
	Ints    = &swag.RouteModel{Model: Int32, Struct: Int32, RetArray: true}
	Bytes   = &swag.RouteModel{Model: Uint8, Struct: Uint8, RetArray: true}
	Strings = &swag.RouteModel{Model: String, Struct: String, RetArray: true}
	Floats  = &swag.RouteModel{Model: Float64, Struct: Float64, RetArray: true}
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
type BaseModel = swag.BaseModel
type BaseModelIface = swag.BaseModelIface

//goland:noinspection GoUnusedGlobalVariable
var (
	ValidationErrorResponse = app.ValidationErrorResponse
	AnyResponse             = app.AnyResponse
	JSONResponse            = app.JSONResponse
	StringResponse          = app.StringResponse
	StreamResponse          = app.StreamResponse
	FileResponse            = app.FileResponse
	ErrorResponse           = app.ErrorResponse
	HTMLResponse            = app.HTMLResponse
	OKResponse              = app.OKResponse
	ResourceNotFound        = app.ResourceNotFound
	AdvancedResponse        = app.AdvancedResponse
	StringsReverse          = swag.StringsReverse
)

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
