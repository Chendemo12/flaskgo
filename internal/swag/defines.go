package swag

import (
	"gitlab.cowave.com/gogo/flaskgo/internal/constant"
	"reflect"
)

type BaseModelKind int

const (
	RModelKind BaseModelKind = iota + 1
	RFieldKind
	QModelKind
	RouteRespKind
	RouteInsKind
	RouteInsGroupKind
)

// 用于swagger的一些静态文件，来自FastApi
const (
	swaggerCssUrl     = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@4/swagger-ui.css"
	swaggerFaviconUrl = "https://fastapi.tiangolo.com/img/favicon.png"
	swaggerJsUrl      = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@4/swagger-ui-bundle.js"
	openapiUrl        = "openapi.json"
)

var (
	// ------------------------------------- int ---------------------------------------

	Int8 = RModelField{
		Name:        "Int8",
		Tag:         `json:"int8" gte:"-128" lte:"127" description:"int8" default:"0"`,
		Type:        constant.IntegerName,
		ReflectKind: reflect.Int8,
	}
	Int16 = RModelField{
		Name:        "Int16",
		Tag:         `json:"int16" gte:"-32768" lte:"32767" description:"int16" default:"0"`,
		Type:        constant.IntegerName,
		ReflectKind: reflect.Int16,
	}
	Int32 = RModelField{
		Name:        "Int32",
		Tag:         `json:"int32" gte:"-2147483648" lte:"2147483647" description:"int32" default:"0"`,
		Type:        constant.IntegerName,
		ReflectKind: reflect.Int32,
	}
	Int64 = RModelField{
		Name:        "Int64",
		Tag:         `json:"int64" gte:"-9223372036854775808" lte:"9223372036854775807" description:"int64" default:"0"`,
		Type:        constant.IntegerName,
		ReflectKind: reflect.Int64,
	}

	// ------------------------------------- uint ---------------------------------------

	Uint8 = RModelField{
		Name:        "Uint8",
		Tag:         `json:"uint8" gte:"0" lte:"255" description:"uint8"`,
		Type:        constant.IntegerName,
		ReflectKind: reflect.Uint8,
	}
	Uint16 = RModelField{
		Name:        "Uint16",
		Tag:         `json:"uint16" gte:"0" lte:"65535" description:"uint16" default:"0"`,
		Type:        constant.IntegerName,
		ReflectKind: reflect.Uint16,
	}
	Uint32 = RModelField{
		Name:        "Uint32",
		Tag:         `json:"uint32" gte:"0" lte:"4294967295" description:"uint32" default:"0"`,
		Type:        constant.IntegerName,
		ReflectKind: reflect.Uint32,
	}
	Uint64 = RModelField{
		Name:        "Uint64",
		Tag:         `json:"uint64" gte:"0" lte:"18446744073709551615" description:"uint64" default:"0"`,
		Type:        constant.IntegerName,
		ReflectKind: reflect.Uint64,
	}

	// ------------------------------------- Float ---------------------------------------

	Float32 = RModelField{
		Name:        "Float32",
		Tag:         `json:"float32" description:"float32" default:"0.0"`,
		Type:        constant.NumberName,
		ReflectKind: reflect.Float32,
	}
	Float64 = RModelField{
		Name:        "Float64",
		Tag:         `json:"float64" description:"float64" default:"0.0"`,
		Type:        constant.NumberName,
		ReflectKind: reflect.Float64,
	}

	// ------------------------------------- other ---------------------------------------

	String = RModelField{
		Name:        "String",
		Tag:         `json:"string" min:"0" max:"255" description:"string" default:""`,
		Type:        "string",
		ReflectKind: reflect.String,
	}
	Boolean = RModelField{
		Name:        constant.BooleanName,
		Tag:         `json:"boolean" oneof:"true false" description:"boolean" default:"false"`,
		Type:        constant.BooleanName,
		ReflectKind: reflect.Bool,
	}
	Mapping = RModelField{
		Name:        "mapping",
		Tag:         `json:"mapping"`,
		Type:        constant.ObjectName,
		ReflectKind: reflect.Map,
	}
)

type BaseModelIface interface {
	Doc__() string // 模型描述
}

type BaseModel struct{}

func (b BaseModel) Doc__() string { return "" }
