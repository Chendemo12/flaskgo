package godantic

type Define struct {
	BaseModel
	Name string
}

var String = &Define{Name: "String"}

//
//var (
//	// ------------------------------------- int ---------------------------------------
//
//	Int8 = RModelField{
//		names:        "Int8",
//		Tag:         `json:"int8" gte:"-128" lte:"127" description:"int8" default:"0"`,
//		RType:        godantic.IntegerType,
//		ReflectKind: reflect.Int8,
//	}
//	Int16 = RModelField{
//		names:        "Int16",
//		Tag:         `json:"int16" gte:"-32768" lte:"32767" description:"int16" default:"0"`,
//		RType:        godantic.IntegerType,
//		ReflectKind: reflect.Int16,
//	}
//	Int32 = RModelField{
//		names:        "Int32",
//		Tag:         `json:"int32" gte:"-2147483648" lte:"2147483647" description:"int32" default:"0"`,
//		RType:        godantic.IntegerType,
//		ReflectKind: reflect.Int32,
//	}
//	Int64 = RModelField{
//		names:        "Int64",
//		Tag:         `json:"int64" gte:"-9223372036854775808" lte:"9223372036854775807" description:"int64" default:"0"`,
//		RType:        godantic.IntegerType,
//		ReflectKind: reflect.Int64,
//	}
//
//	// ------------------------------------- uint ---------------------------------------
//
//	Uint8 = RModelField{
//		names:        "Uint8",
//		Tag:         `json:"uint8" gte:"0" lte:"255" description:"uint8"`,
//		RType:        godantic.IntegerType,
//		ReflectKind: reflect.Uint8,
//	}
//	Uint16 = RModelField{
//		names:        "Uint16",
//		Tag:         `json:"uint16" gte:"0" lte:"65535" description:"uint16" default:"0"`,
//		RType:        godantic.IntegerType,
//		ReflectKind: reflect.Uint16,
//	}
//	Uint32 = RModelField{
//		names:        "Uint32",
//		Tag:         `json:"uint32" gte:"0" lte:"4294967295" description:"uint32" default:"0"`,
//		RType:        godantic.IntegerType,
//		ReflectKind: reflect.Uint32,
//	}
//	Uint64 = RModelField{
//		names:        "Uint64",
//		Tag:         `json:"uint64" gte:"0" lte:"18446744073709551615" description:"uint64" default:"0"`,
//		RType:        godantic.IntegerType,
//		ReflectKind: reflect.Uint64,
//	}
//
//	// ------------------------------------- Float ---------------------------------------
//
//	Float32 = RModelField{
//		names:        "Float32",
//		Tag:         `json:"float32" description:"float32" default:"0.0"`,
//		RType:        godantic.NumberType,
//		ReflectKind: reflect.Float32,
//	}
//	Float64 = RModelField{
//		names:        "Float64",
//		Tag:         `json:"float64" description:"float64" default:"0.0"`,
//		RType:        godantic.NumberType,
//		ReflectKind: reflect.Float64,
//	}
//
//	// ------------------------------------- other ---------------------------------------
//
//	String = RModelField{
//		names:        "String",
//		Tag:         `json:"string" min:"0" max:"255" description:"string" default:""`,
//		RType:        "string",
//		ReflectKind: reflect.String,
//	}
//	Boolean = RModelField{
//		names:        godantic.BoolType,
//		Tag:         `json:"boolean" oneof:"true false" description:"boolean" default:"false"`,
//		RType:        godantic.BoolType,
//		ReflectKind: reflect.Bool,
//	}
//	Mapping = RModelField{
//		names:        "mapping",
//		Tag:         `json:"mapping"`,
//		RType:        godantic.ObjectType,
//		ReflectKind: reflect.Map,
//	}
//)
//
//type BaseModelIface interface {
//	Doc__() string // 模型描述
//}
//
//type BaseModel struct{}
//
//func (b BaseModel) Doc__() string { return "" }
