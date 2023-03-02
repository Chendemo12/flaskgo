package godantic

import (
	"github.com/Chendemo12/functools/helper"
	"github.com/Chendemo12/functools/structfuncs"
	"reflect"
	"strconv"
	"strings"
)

type dict map[string]any

// Field 基本数据模型, 此模型不可再分, 同时也是 BaseModel 的字段类型
// 但此类型不再递归记录,仅记录一个关联模型为基本
type Field struct {
	_pkg        string            `description:"包名.结构体名"`
	Title       string            `json:"title" description:"字段名称"`
	Tag         reflect.StructTag `json:"tag" description:"字段标签"`
	Description string            `json:"description,omitempty" description:"说明"`
	Default     string            `json:"default" description:"默认值"` // 暂时仅限 swagger 使用，后期也应在字段校验时使用
	ItemRef     string            `description:"子元素类型, 仅Type=array/object时有效"`
	OType       OpenApiDataType   `json:"otype,omitempty" description:"openaapi 数据类型"`
}

// Schema 生成字段的详细描述信息
//
//	// 字段为结构体类型
//
//	"position_sat": {
//		"title": "position_sat",
//		"type": "object"
//		"description": "position_sat",
//		"required": false,
//		"$ref": "#/comonents/schemas/example.PositionGeo",
//	}
//
//	// 字段为数组类型, 数组元素为基本类型
//
//	"traffic_timeslot": {
//		"title": "traffic_timeslot",
//		"type": "array"
//		"description": "业务时隙编号数组",
//		"required": false,
//		"items": {
//			"type": "integer"
//		},
//	}
//
//	// 字段为数组类型, 数组元素为自定义结构体类型
//
//	"Detail": {
//		"title": "Detail",
//		"type": "array"
//		"description": "Detail",
//		"required": true,
//		"items": {
//			"$ref": "#/comonents/schemas/ValidationError"
//		},
//	}
func (f *Field) Schema() (m map[string]any) {
	// 最基础的属性，必须
	m = dict{
		"title":       f.Title,
		"type":        f.OType,
		"required":    f.IsRequired(),
		"description": f.SchemaDesc(),
	}
	// 生成默认值
	if f.Default != "" {
		m["default"] = f.Default
	}
	// 生成字段的枚举值
	if es := QueryFieldTag(f.Tag, "oneof", ""); es != "" {
		m["enum"] = strings.Split(es, " ")
	}

	// 为不同的字段类型生成相应的描述
	switch f.OType {
	case IntegerType, NumberType: // 生成数字类型的最大最小值
		if lt := QueryFieldTag(f.Tag, "lt", ""); lt != "" {
			m["maximum"], _ = strconv.Atoi(lt)
		}
		if gt := QueryFieldTag(f.Tag, "gt", ""); gt != "" {
			m["minimum"], _ = strconv.Atoi(gt)
		}

		if lt := QueryFieldTag(f.Tag, "lte", ""); lt != "" {
			m["exclusiveMaximum"], _ = strconv.Atoi(lt)
		}
		if gt := QueryFieldTag(f.Tag, "gte", ""); gt != "" {
			m["exclusiveMinimum"], _ = strconv.Atoi(gt)
		}
		// 存在多个标记
		if lt := QueryFieldTag(f.Tag, "max", ""); lt != "" {
			m["maximum"], _ = strconv.Atoi(lt)
		}
		if gt := QueryFieldTag(f.Tag, "min", ""); gt != "" {
			m["minimum"], _ = strconv.Atoi(gt)
		}

	case StringType: // 生成字符串类型的最大最小长度
		if lt := QueryFieldTag(f.Tag, "max", ""); lt != "" {
			m["maxLength"], _ = strconv.Atoi(lt)
		}
		if gt := QueryFieldTag(f.Tag, "min", ""); gt != "" {
			m["minLength"], _ = strconv.Atoi(gt)
		}

	case ArrayType:
		// 为数组类型生成子类型描述
		if f.ItemRef != "" {
			if strings.HasPrefix(f.ItemRef, RefPrefix) { // 数组子元素为关联类型
				m["items"] = map[string]string{"$ref": RefPrefix + f.ItemRef}
			} else { // 子元素为基本数据类型
				m["items"] = map[string]string{"type": f.ItemRef}
			}
		} else { // 缺省为string
			m["items"] = map[string]OpenApiDataType{"type": StringType}
		}
		// 限制数组的长度
		if lt := QueryFieldTag(f.Tag, "max", ""); lt != "" {
			m["maxLength"], _ = strconv.Atoi(lt)
		}
		if gt := QueryFieldTag(f.Tag, "min", ""); gt != "" {
			m["minLength"], _ = strconv.Atoi(gt)
		}

	case ObjectType:
		if f.ItemRef != "" { // 字段类型为自定义结构体，生成关联类型，此内部结构体已注册
			m["$ref"] = RefPrefix + f.ItemRef
		}

	default:
	}

	return
}

// SchemaName swagger文档字段名
func (f *Field) SchemaName(exclude ...bool) string { return QueryJsonName(f.Tag, f.Title) }

// SchemaDesc 字段注释说明
func (f *Field) SchemaDesc() string { return f.Description }

// SchemaType 模型类型
func (f *Field) SchemaType() OpenApiDataType { return f.OType }

// SchemaJson swagger文档字符串格式
func (f *Field) SchemaJson() string {
	bytes, err := helper.DefaultJsonMarshal(f.Schema())
	if err != nil {
		return string(bytes)
	} else {
		return ""
	}
}

// IsRequired 字段是否必须
func (f *Field) IsRequired() bool { return IsFieldRequired(f.Tag) }

// IsArray 字段是否是数组类型
func (f *Field) IsArray() bool { return f.OType == ArrayType }

// InnerSchema 内部字段模型文档, 全名:文档
func (f *Field) InnerSchema() (m map[string]map[string]any) {
	m = make(map[string]map[string]any)
	return
}

func (f *Field) MetaData() *MetaData { return GetMetaData(f._pkg) }

func (f *Field) SetId(id string) { f._pkg = id }

// BaseModel 基本数据模型, 对于上层的 app.Route 其请求和相应体都应为继承此结构体的结构体
// 在 OpenApi 文档模型中,此模型的类型始终为 "object";
// 对于 BaseModel 其字段仍然可能会 BaseModel
type BaseModel struct {
	_pkg string `description:"包名.结构体名"`
}

// Schema 输出为OpenAPI文档模型,字典格式
//
//	{
//		"title": "examle.MyTimeslot",
//		"type": "object"
//		"description": "examle.mytimeslot",
//		"required": [],
//		"properties": {
//			"control_timeslot": {
//				"title": "control_timeslot",
//				"type": "array"
//				"description": "控制时隙编号数组",
//				"required": false,
//				"items": {
//					"type": "integer"
//				},
//			},
//			"superframe_count": {
//				"title": "superframe_count",
//				"type": "integer"
//				"description": "超帧计数",
//				"required": false,
//			},
//		},
//	},
func (b *BaseModel) Schema() (m map[string]any) {
	// 模型标题排除包名
	m = dict{"title": b.SchemaName(true), "type": b.SchemaType(), "description": b.SchemaDesc()}

	meta := GetMetaData(b._pkg)
	required := make([]string, 0, len(meta.fields))
	properties := make(map[string]any, len(meta.fields))

	for i := 0; i < len(meta.fields); i++ {
		if !meta.fields[i].Exported || meta.fields[i].Anonymous { // 非导出字段
			continue
		}

		properties[meta.fields[i].SchemaName()] = meta.fields[i].Schema()
		if meta.fields[i].IsRequired() {
			required = append(required, meta.fields[i].SchemaName())
		}
	}

	m["required"], m["properties"] = required, properties

	return
}

// SchemaName 获取结构体的名称,默认包含包名
// @param  exclude  []bool  是否排除包名LL
func (b *BaseModel) SchemaName(exclude ...bool) string {
	meta := GetMetaData(b._pkg)
	if len(exclude) > 0 { // 排除包名
		return meta.names[0]
	} else {
		return meta.names[1]
	}
}

// SchemaDesc 结构体文档注释
func (b *BaseModel) SchemaDesc() string { return "BaseModel" }

// SchemaType 模型类型
func (b *BaseModel) SchemaType() OpenApiDataType { return ObjectType }

// SchemaJson 输出为OpenAPI文档模型,字符串格式
func (b *BaseModel) SchemaJson() string {
	bytes, err := helper.DefaultJsonMarshal(b.Schema())
	if err != nil {
		return string(bytes)
	} else {
		return ""
	}
}

// InnerSchema 内部字段模型文档
func (b *BaseModel) InnerSchema() (m map[string]map[string]any) {
	meta := GetMetaData(b._pkg)
	for i := 0; i < len(meta.innerFields); i++ {
		// TODO: error
		if meta.innerFields[i].Exported && !strings.HasPrefix(meta.innerFields[i].SchemaName(), "_") {
			m[meta.innerFields[i].SchemaName()] = meta.innerFields[i].Schema()
		}
	}

	return
}

func (b *BaseModel) IsRequired() bool { return true }

// Map 将结构体转换为字典视图
func (b *BaseModel) Map() (m map[string]any) {
	m = structfuncs.GetFieldsValue(b)
	return
}

// Dict 将结构体转换为字典视图，并允许过滤一些字段或添加一些键值对到字典中
func (b *BaseModel) Dict(exclude []string, include map[string]any) (m map[string]any) {

	excludeMap := make(map[string]string, len(exclude))
	for i := 0; i < len(exclude); i++ {
		excludeMap[exclude[i]] = exclude[i]
	}

	// 实时反射取值
	v := reflect.Indirect(reflect.ValueOf(b))
	meta := GetMetaData(b._pkg)

	for i := 0; i < len(meta.Fields()); i++ {
		if !meta.fields[i].Exported || meta.fields[i].Anonymous { // 非导出字段
			continue
		}

		if _, ok := excludeMap[meta.fields[i].Title]; ok { // 此字段被排除
			continue
		}

		switch meta.fields[i].RType.Kind() { // 获取字段定义的类型

		case reflect.Array, reflect.Slice:
			m[meta.fields[i].Title] = v.Field(meta.fields[i].Index).Bytes()

		case reflect.Uint8:
			m[meta.fields[i].Title] = byte(v.Field(meta.fields[i].Index).Uint())
		case reflect.Uint16:
			m[meta.fields[i].Title] = uint16(v.Field(meta.fields[i].Index).Uint())
		case reflect.Uint32:
			m[meta.fields[i].Title] = uint32(v.Field(meta.fields[i].Index).Uint())
		case reflect.Uint64, reflect.Uint:
			m[meta.fields[i].Title] = v.Field(meta.fields[i].Index).Uint()

		case reflect.Int8:
			m[meta.fields[i].Title] = int8(v.Field(meta.fields[i].Index).Int())
		case reflect.Int16:
			m[meta.fields[i].Title] = int16(v.Field(meta.fields[i].Index).Int())
		case reflect.Int32:
			m[meta.fields[i].Title] = int32(v.Field(meta.fields[i].Index).Int())
		case reflect.Int64, reflect.Int:
			m[meta.fields[i].Title] = v.Field(meta.fields[i].Index).Int()

		case reflect.Float32:
			m[meta.fields[i].Title] = float32(v.Field(meta.fields[i].Index).Float())
		case reflect.Float64:
			m[meta.fields[i].Title] = v.Field(meta.fields[i].Index).Float()

		case reflect.Struct, reflect.Interface, reflect.Map:
			m[meta.fields[i].Title] = v.Field(meta.fields[i].Index).Interface()

		case reflect.String:
			m[meta.fields[i].Title] = v.Field(meta.fields[i].Index).String()

		case reflect.Pointer:
			m[meta.fields[i].Title] = v.Field(meta.fields[i].Index).Pointer()
		case reflect.Bool:
			m[meta.fields[i].Title] = v.Field(meta.fields[i].Index).Bool()
		}

	}

	if include != nil {
		for k := range include {
			m[k] = include[k]
		}
	}

	return
}

// Exclude 将结构体转换为字典视图，并过滤一些字段
func (b *BaseModel) Exclude(exclude ...string) (m map[string]any) {
	return b.Dict(exclude, nil)
}

// Include 将结构体转换为字典视图，并允许添加一些键值对到字典中
func (b *BaseModel) Include(include map[string]any) (m map[string]any) {
	return b.Dict([]string{}, include)
}

// Validate 检验实例是否符合tag要求
func (b *BaseModel) Validate(stc any) []*ValidationError {
	// TODO: NotImplemented
	return nil
}

// ParseRaw 从原始字节流中解析结构体对象
func (b *BaseModel) ParseRaw(stc []byte) []*ValidationError {
	// TODO: NotImplemented
	return nil
}

// Copy 拷贝一个新的空实例对象
func (b *BaseModel) Copy() any {
	// TODO: NotImplemented
	return nil
}

// MetaData 获取反射后的字段元信息, 此字段应慎重使用
func (b *BaseModel) MetaData() *MetaData { return GetMetaData(b._pkg) }

// SetId 设置结构体的唯一标识
func (b *BaseModel) SetId(id string) { b._pkg = id }

// ValidationError 参数校验错误
type ValidationError struct {
	BaseModel
	Ctx  map[string]any `json:"service" Description:"Service"`
	Msg  string         `json:"msg" Description:"Message" binding:"required"`
	Type string         `json:"type" Description:"Error RType" binding:"required"`
	Loc  []string       `json:"loc" Description:"Location" binding:"required"`
}

func (v *ValidationError) SchemaDesc() string { return "参数校验错误" }

func (v *ValidationError) Map() (m map[string]any) {
	m["loc"], m["msg"] = v.Loc, v.Msg
	m["service"], m["type"] = v.Ctx, v.Type

	return
}
