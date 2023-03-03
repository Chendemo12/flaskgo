package godantic

func init() {
	// 初始化基本类型
	String.SetId("godantic.str")
	String.Description = QueryFieldTag(String.Tag, "description", String.SchemaName())
	String.Default = QueryFieldTag(String.Tag, "default", "")
	SaveMetadata(newEmptyMeta("str"))

	Bool.SetId("godantic.bool")
	Bool.Description = QueryFieldTag(Bool.Tag, "description", Bool.SchemaName())
	Bool.Default = QueryFieldTag(Bool.Tag, "default", "")
	SaveMetadata(newEmptyMeta("bool"))

	// integer
	Int.SetId("godantic.int")
	Int.Description = QueryFieldTag(Int.Tag, "description", Int.SchemaName())
	Int.Default = QueryFieldTag(Int.Tag, "default", "")
	SaveMetadata(newEmptyMeta("int"))

	Int8.SetId("godantic.int8")
	Int8.Description = QueryFieldTag(Int8.Tag, "description", Int8.SchemaName())
	Int8.Default = QueryFieldTag(Int8.Tag, "default", "")
	SaveMetadata(newEmptyMeta("int8"))

	Int16.SetId("godantic.int16")
	Int16.Description = QueryFieldTag(Int16.Tag, "description", Int16.SchemaName())
	Int16.Default = QueryFieldTag(Int16.Tag, "default", "")
	SaveMetadata(newEmptyMeta("int16"))

	Int32.SetId("godantic.int32")
	Int32.Description = QueryFieldTag(Int32.Tag, "description", Int32.SchemaName())
	Int32.Default = QueryFieldTag(Int32.Tag, "default", "")
	SaveMetadata(newEmptyMeta("int32"))

	Int64.SetId("godantic.int64")
	Int64.Description = QueryFieldTag(Int64.Tag, "description", Int64.SchemaName())
	Int64.Default = QueryFieldTag(Int64.Tag, "default", "")
	SaveMetadata(newEmptyMeta("int64"))

	Uint8.SetId("godantic.uint8")
	Uint8.Description = QueryFieldTag(Uint8.Tag, "description", Uint8.SchemaName())
	Uint8.Default = QueryFieldTag(Uint8.Tag, "default", "")
	SaveMetadata(newEmptyMeta("uint8"))

	Uint16.SetId("godantic.uint16")
	Uint16.Description = QueryFieldTag(Uint16.Tag, "description", Uint16.SchemaName())
	Uint16.Default = QueryFieldTag(Uint16.Tag, "default", "")
	SaveMetadata(newEmptyMeta("uint16"))

	Uint32.SetId("godantic.uint32")
	Uint32.Description = QueryFieldTag(Uint32.Tag, "description", Uint32.SchemaName())
	Uint32.Default = QueryFieldTag(Uint32.Tag, "default", "")
	SaveMetadata(newEmptyMeta("uint32"))

	Uint64.SetId("godantic.uint64")
	Uint64.Description = QueryFieldTag(Uint64.Tag, "description", Uint64.SchemaName())
	Uint64.Default = QueryFieldTag(Uint64.Tag, "default", "")
	SaveMetadata(newEmptyMeta("uint64"))

	Float.SetId("godantic.float")
	Float.Description = QueryFieldTag(Float.Tag, "description", Float.SchemaName())
	Float.Default = QueryFieldTag(Float.Tag, "default", "")
	SaveMetadata(newEmptyMeta("float64"))

	Float32.SetId("godantic.float32")
	Float32.Description = QueryFieldTag(Float32.Tag, "description", Float32.SchemaName())
	Float32.Default = QueryFieldTag(Float32.Tag, "default", "")
	SaveMetadata(newEmptyMeta("float32"))

	Float64.SetId("godantic.float64")
	Float64.Description = QueryFieldTag(Float64.Tag, "description", Float64.SchemaName())
	Float64.Default = QueryFieldTag(Float64.Tag, "default", "")
	SaveMetadata(newEmptyMeta("float64"))
}

func newEmptyMeta(name string) *Metadata {
	return &Metadata{
		names:       []string{name, "godantic." + name},
		fields:      make([]*MetaField, 0),
		innerFields: make([]*MetaField, 0),
	}
}

var (
	// ------------------------------------- int ---------------------------------------

	Int8 = &Field{
		_pkg:        "godantic.int8",
		Title:       "Int8",
		Tag:         `json:"int8" gte:"-128" lte:"127" description:"int8" default:"0"`,
		OType:       IntegerType,
		Description: "",
		Default:     "",
	}
	Int16 = &Field{
		_pkg:        "godantic.int16",
		Title:       "Int16",
		Tag:         `json:"int16" gte:"-32768" lte:"32767" description:"int16" default:"0"`,
		OType:       IntegerType,
		Description: "",
		Default:     "",
	}
	Int32 = &Field{
		_pkg:        "godantic.int32",
		Title:       "Int32",
		Tag:         `json:"int32" gte:"-2147483648" lte:"2147483647" description:"int32" default:"0"`,
		OType:       IntegerType,
		Description: "",
		Default:     "",
	}
	Int64 = &Field{
		_pkg:        "godantic.int64",
		Title:       "Int64",
		Tag:         `json:"int64" gte:"-9223372036854775808" lte:"9223372036854775807" description:"int64" default:"0"`,
		OType:       IntegerType,
		Description: "",
		Default:     "",
	}
	Int = &Field{
		_pkg:        "godantic.int",
		Title:       "Int",
		Tag:         `json:"int" gte:"-9223372036854775808" lte:"9223372036854775807" description:"int" default:"0"`,
		OType:       IntegerType,
		Description: "",
		Default:     "",
	}

	// ------------------------------------- uint ---------------------------------------

	Uint8 = &Field{
		_pkg:        "godantic.uint8",
		Title:       "Uint8",
		Tag:         `json:"uint8" gte:"0" lte:"255" description:"uint8"`,
		OType:       IntegerType,
		Description: "",
		Default:     "",
	}
	Uint16 = &Field{
		_pkg:        "godantic.uint16",
		Title:       "Uint16",
		Tag:         `json:"uint16" gte:"0" lte:"65535" description:"uint16" default:"0"`,
		OType:       IntegerType,
		Description: "",
		Default:     "",
	}
	Uint32 = &Field{
		_pkg:        "godantic.uint32",
		Title:       "Uint32",
		Tag:         `json:"uint32" gte:"0" lte:"4294967295" description:"uint32" default:"0"`,
		OType:       IntegerType,
		Description: "",
		Default:     "",
	}
	Uint64 = &Field{
		_pkg:        "godantic.uint64",
		Title:       "Uint64",
		Tag:         `json:"uint64" gte:"0" lte:"18446744073709551615" description:"uint64" default:"0"`,
		OType:       IntegerType,
		Description: "",
		Default:     "",
	}

	// ------------------------------------- Float ---------------------------------------

	Float32 = &Field{
		_pkg:        "godantic.float32",
		Title:       "Float32",
		Tag:         `json:"float32" description:"float32" default:"0.0"`,
		OType:       NumberType,
		Description: "",
		Default:     "",
	}
	Float64 = &Field{
		_pkg:        "godantic.float64",
		Title:       "Float64",
		Tag:         `json:"float64" description:"float64" default:"0.0"`,
		OType:       NumberType,
		Description: "",
		Default:     "",
	}
	Float = &Field{
		_pkg:        "godantic.float",
		Title:       "Float",
		Tag:         `json:"float" description:"float" default:"0.0"`,
		OType:       NumberType,
		Description: "",
		Default:     "",
	}

	// ------------------------------------- other ---------------------------------------

	String = &Field{
		_pkg:        "godantic.string",
		Title:       "String",
		Tag:         `json:"string" min:"0" max:"255" description:"string" default:""`,
		OType:       StringType,
		Description: "",
		Default:     "",
	}
	Bool = &Field{
		_pkg:        "godantic.bool",
		Title:       "Bool",
		Tag:         `json:"boolean" oneof:"true false" description:"boolean" default:"false"`,
		OType:       BoolType,
		Description: "",
		Default:     "",
	}
)

func List(model SchemaIface) *Field {
	meta := GetMetadataFactory().Reflect(model)
	SaveMetadata(meta)
	model.SetId(meta.Id())

	return &Field{
		_pkg:        meta.Id(),
		Title:       meta.String(),
		Tag:         "",
		Description: model.SchemaDesc(),
		Default:     "",
		ItemRef:     model.SchemaName(),
		OType:       ArrayType,
	}
}
