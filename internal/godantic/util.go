package godantic

import (
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// reflectKindToOType 转换reflect.Kind为swagger类型说明
// @param  ReflectKind  reflect.Kind  反射类型
func reflectKindToOType(kind reflect.Kind) (name OpenApiDataType) {
	switch kind {

	case reflect.Array, reflect.Slice, reflect.Chan:
		name = ArrayType
	case reflect.String:
		name = StringType
	case reflect.Bool:
		name = BoolType
	default:
		if reflect.Bool < kind && kind <= reflect.Uint64 {
			name = IntegerType
		} else if reflect.Float32 <= kind && kind <= reflect.Complex128 {
			name = NumberType
		} else {
			name = ObjectType
		}
	}

	return
}

// IsFieldRequired 从tag中判断此字段是否是必须的
func IsFieldRequired(tag reflect.StructTag) bool {
	for _, name := range []string{"binding", "validate"} {
		bindings := strings.Split(QueryFieldTag(tag, name, ""), ",") // binding 存在多个值
		for i := 0; i < len(bindings); i++ {
			if strings.TrimSpace(bindings[i]) == "required" {
				return true
			}
		}
	}

	return false
}

// GetDefaultV 从Tag中提取字段默认值
func GetDefaultV(tag reflect.StructTag, otype OpenApiDataType) (v any) {
	defaultV := QueryFieldTag(tag, "default", "")

	if defaultV == "" {
		v = nil
	} else { // 存在默认值
		switch otype {

		case StringType:
			v = defaultV
		case IntegerType:
			v, _ = strconv.Atoi(defaultV)
		case NumberType:
			v, _ = strconv.ParseFloat(defaultV, 64)
		case BoolType:
			v, _ = strconv.ParseBool(defaultV)
		default:
			v = defaultV
		}
	}
	return
}

// IsArray 判断一个对象是否是数组类型
func IsArray(object any) bool {
	if object == nil {
		return false
	}
	return reflectKindToOType(reflect.TypeOf(object).Kind()) == ArrayType
}

// QueryFieldTag 查找struct字段的Tag
// @param   tag        reflect.StructTag  字段的Tag
// @param   label      string             要查找的标签
// @param   undefined  string             当查找的标签不存在时返回的默认值
// @return  string 查找到的标签值, 不存在则返回提供的默认值
func QueryFieldTag(tag reflect.StructTag, label string, undefined string) string {
	if tag == "" {
		return undefined
	}
	if v := tag.Get(label); v != "" {
		return v
	}
	return undefined
}

// QueryJsonName 查询字段定义的json名称
func QueryJsonName(tag reflect.StructTag, undefined string) string {
	if tag == "" {
		return undefined
	}
	if v := tag.Get("json"); v != "" {
		return strings.TrimSpace(strings.Split(v, ",")[0])
	}
	return undefined
}

// Deprecated: AnyToQModel 将表示查询参数的struct或map转换成查询参数模型
// @param  object  any  查询参数模型，此类别必须为  map[string]bool  或  struct
func AnyToQModel(object any) (m []*QModel) {
	if object == nil {
		return
	}

	rt := reflect.TypeOf(object)
	switch rt.Kind() {

	case reflect.Map: // 若为map类型，则key标识字段名称，value表示参数是否必须
		if fields, ok := object.(map[string]bool); ok {
			m = make([]*QModel, 0)

			for field, require := range fields {
				qm := &QModel{Name: field, Required: require, InPath: false}
				if require {
					qm.Tag = reflect.StructTag(`json:"` + field + `" binding:"required" validate:"required"`)
				} else {
					qm.Tag = reflect.StructTag(`json:"` + field + `"`)
				}
				m = append(m, qm)
			}
		}

	case reflect.Struct:
		// 当此model作为查询参数时，此struct的每一个字段都将作为一个查询参数
		m = make([]*QModel, rt.NumField())
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)

			if unicode.IsLower(rune(field.Name[0])) {
				continue
			}

			// 仅导出字段可用
			m[i] = &QModel{
				Name:     QueryJsonName(field.Tag, field.Name), // 以json字段为准
				Tag:      field.Tag,
				Required: IsFieldRequired(field.Tag),
			}
		}
	}

	return
}
