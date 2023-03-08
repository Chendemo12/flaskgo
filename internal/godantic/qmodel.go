package godantic

import (
	"github.com/Chendemo12/functools/helper"
	"reflect"
	"unicode"
)

// QModel 查询参数或路径参数模型, 此类型会进一步转换为 openapi.Parameter
type QModel struct {
	Title  string            `json:"names,omitempty" description:"字段名称"`
	Tag    reflect.StructTag `json:"tag,omitempty" description:"TAG"`
	OType  OpenApiDataType   `json:"otype,omitempty" description:"openaapi 数据类型"`
	InPath bool              `json:"in_path,omitempty" description:"是否是路径参数"`
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

// SchemaName 获取名称,以json字段为准
func (q *QModel) SchemaName(exclude ...bool) string { return QueryJsonName(q.Tag, q.Title) }

// SchemaDesc 结构体文档注释
func (q *QModel) SchemaDesc() string { return QueryFieldTag(q.Tag, "description", q.Title) }

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
func (q *QModel) IsRequired() bool { return IsFieldRequired(q.Tag) }

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
			Title:  field.Name,
			InPath: false,
			Tag:    field.Tag,
			OType:  ObjectType, // 无意义
		}
	}
	return m
}
