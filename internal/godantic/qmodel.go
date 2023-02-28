package godantic

import (
	"github.com/Chendemo12/functools/helper"
	"reflect"
	"unicode"
)

// QModel 查询参数或路径参数模型, 此类型会进一步转换为 openapi.Parameter
type QModel struct {
	Name     string            `json:"names,omitempty" description:"字段名称"`
	Required bool              `json:"required,omitempty" description:"是否必须"`
	InPath   bool              `json:"in_path,omitempty" description:"是否是路径参数"`
	Tag      reflect.StructTag `json:"tag,omitempty" description:"TAG"`
	OType    OpenApiDataType   `json:"otype,omitempty" description:"openaapi 数据类型"`
}

// Schema 输出为OpenAPI文档模型,字典格式
//
//	{
//		"required": true,
//		"schema": {
//			"title": "names",
//			"type": "string",
//			"default": "jack"
//		},
//		"names": "names",
//		"in": "query"/"path"
//	}
func (q *QModel) Schema() (m map[string]any) {
	return
}

// SchemaRef 模型引用文档
func (q *QModel) SchemaRef() (m map[string]any) {
	m[RefName] = RefPrefix + q.SchemaName()
	return
}

// SchemaName 获取结构体的名称,默认包含包名
func (q *QModel) SchemaName(exclude ...bool) string { return q.Name }

// SchemaDesc 结构体文档注释
func (q *QModel) SchemaDesc() string {
	if q.InPath {
		return "field in path"
	} else {
		return "field in query"
	}
}

// SchemaType 模型类型
func (q *QModel) SchemaType() OpenApiDataType { return StringType }

// SchemaJson 输出为OpenAPI文档模型,字符串格式
func (q *QModel) SchemaJson() string {
	bytes, err := helper.DefaultJsonMarshal(q.Schema())
	if err != nil {
		return string(bytes)
	} else {
		return ""
	}
}

// InnerSchema 内部字段模型文档, 全名:文档
func (q *QModel) InnerSchema() (m map[string]map[string]any) {
	m = make(map[string]map[string]any)
	return
}

// IsRequired 是否必须
func (q *QModel) IsRequired() bool { return q.Required }

// QueryModel 查询参数基类
type QueryModel struct{}

func (q *QueryModel) Fields() []*QModel {
	rt := reflect.TypeOf(q)
	// 当此model作为查询参数时，此struct的每一个字段都将作为一个查询参数
	m := make([]*QModel, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		if unicode.IsLower(rune(field.Name[0])) {
			continue
		}

		// 仅导出字段可用
		m[i] = &QModel{
			Name:     QueryJsonName(field.Tag, field.Name), // 以json字段为准
			Required: IsFieldRequired(field.Tag),
			InPath:   false,
			Tag:      field.Tag,
			OType:    ObjectType, // 无意义
		}
	}
	return m
}
