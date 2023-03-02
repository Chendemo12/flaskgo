package godantic

var (
	String  = &str{}
	Bool    = &boolean{}
	Int     = &integer{}
	Int8    = &integer{}
	Int16   = &integer{}
	Int32   = &integer{}
	Int64   = &integer{}
	Uint8   = &integer{}
	Uint16  = &integer{}
	Uint32  = &integer{}
	Uint64  = &integer{}
	Float32 = &float{}
	Float64 = &float{}

	// Mapping = openapi.Mapping
	//

	// Array   = openapi.List
	// List    = openapi.List
	// Ints    = &openapi.RouteModel{Model: Int32, Struct: Int32, RetArray: true}
	// Bytes   = &openapi.RouteModel{Model: Uint8, Struct: Uint8, RetArray: true}
	// Strings = &openapi.RouteModel{Model: String, Struct: String, RetArray: true}
	// Floats  = &openapi.RouteModel{Model: Float64, Struct: Float64, RetArray: true}
)

func init() {
	// 初始化基本类型
	String.SetId("godantic.str")
	SetMetaData(newEmptyMeta("str"))
	Bool.SetId("godantic.bool")
	SetMetaData(newEmptyMeta("bool"))
	// integer
	Int.SetId("godantic.int")
	SetMetaData(newEmptyMeta("int"))
	Int8.SetId("godantic.int8")
	SetMetaData(newEmptyMeta("int8"))
	Int16.SetId("godantic.int16")
	SetMetaData(newEmptyMeta("int16"))
	Int32.SetId("godantic.int32")
	SetMetaData(newEmptyMeta("int32"))
	Int64.SetId("godantic.int64")
	SetMetaData(newEmptyMeta("int64"))

	Uint8.SetId("godantic.uint8")
	SetMetaData(newEmptyMeta("uint8"))
	Uint16.SetId("godantic.uint16")
	SetMetaData(newEmptyMeta("uint16"))
	Uint32.SetId("godantic.uint32")
	SetMetaData(newEmptyMeta("uint32"))
	Uint64.SetId("godantic.uint64")
	SetMetaData(newEmptyMeta("uint64"))

	Float32.SetId("godantic.float32")
	SetMetaData(newEmptyMeta("float32"))
	Float64.SetId("godantic.float64")
	SetMetaData(newEmptyMeta("float64"))
}

func newEmptyMeta(name string) *MetaData {
	return &MetaData{
		names:       []string{name, "godantic." + name},
		fields:      make([]*Field, 0),
		innerFields: make([]*Field, 0),
	}
}

type str struct {
	BaseModel
}

// SchemaName 获取结构体的名称,默认包含包名
func (d *str) SchemaName(exclude ...bool) string { return string(StringType) }
func (d *str) SchemaDesc() string                { return string(StringType) }
func (d *str) SchemaType() OpenApiDataType       { return StringType }
func (d *str) IsRequired() bool                  { return true }
func (d *str) Schema() (m map[string]any) {
	m = make(map[string]any)
	m["title"] = StringType
	return
}

type boolean struct {
	BaseModel
}

func (d *boolean) SchemaName(exclude ...bool) string { return string(BoolType) }
func (d *boolean) SchemaDesc() string                { return string(BoolType) }
func (d *boolean) SchemaType() OpenApiDataType       { return BoolType }
func (d *boolean) IsRequired() bool                  { return true }
func (d *boolean) Schema() (m map[string]any) {
	m = make(map[string]any)
	m["title"] = BoolType
	return
}

type integer struct {
	BaseModel
}

func (d *integer) SchemaName(exclude ...bool) string { return string(IntegerType) }
func (d *integer) SchemaDesc() string                { return string(IntegerType) }
func (d *integer) SchemaType() OpenApiDataType       { return IntegerType }
func (d *integer) IsRequired() bool                  { return true }
func (d *integer) Schema() (m map[string]any) {
	m = make(map[string]any)
	m["title"] = IntegerType
	return
}

type float struct {
	BaseModel
}

func (d *float) SchemaName(exclude ...bool) string { return string(NumberType) }
func (d *float) SchemaDesc() string                { return string(NumberType) }
func (d *float) SchemaType() OpenApiDataType       { return NumberType }
func (d *float) IsRequired() bool                  { return true }
func (d *float) Schema() (m map[string]any) {
	m = make(map[string]any)
	m["title"] = NumberType
	return
}

func List(model SchemaIface) *Field {
	return &Field{
		Title:     model.SchemaName(),
		Index:     0,
		Default:   nil,
		Exported:  true,
		Anonymous: false,
		Tag:       "",
		ItemRef:   model.SchemaName(),
		RType:     nil,
		_pkg:      "",
	}
}

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
//	SchemaDesc() string // 模型描述
//}
//
//type BaseModel struct{}
//
//func (b BaseModel) SchemaDesc() string { return "" }
