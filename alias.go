// Package flaskgo 是一个基于fiber封装了常用方法的软件包
//
// 其提供了类似于FastAPI的API设计，并提供了接口文档自动生成、请求体自动校验和返回值自动序列化等使用功能；
package flaskgo

import (
	"github.com/Chendemo12/flaskgo/internal/app"
	"github.com/Chendemo12/flaskgo/internal/godantic"
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
	S      = godantic.String
	Str    = godantic.String
	String = godantic.String

	B       = godantic.Bool
	Bool    = godantic.Bool
	Boolean = godantic.Bool

	I      = godantic.Int
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

	Float   = godantic.Float
	Float32 = godantic.Float32
	Float64 = godantic.Float64

	L       = godantic.List
	List    = godantic.List
	Array   = godantic.List
	Ints    = godantic.List(Int)
	Bytes   = godantic.List(Byte)
	Strings = godantic.List(S)
	Floats  = godantic.List(Float)
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
