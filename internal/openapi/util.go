package openapi

import (
	"errors"
	"github.com/Chendemo12/flaskgo/internal/constant"
	"github.com/Chendemo12/flaskgo/internal/godantic"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// basicTypeToModel 转换基本的数据类型为 RouteModel , 此返回值应仅用于判断，不应直接修改返回
func basicTypeToModel(rt reflect.Type) (RModelField, error) {
	switch rt.Kind() {

	// 此处类型用于实际接口返回为内置类型的常量或者变量
	case reflect.String:
		return String, nil

	case reflect.Bool: // bool型变量
		return Boolean, nil

	case reflect.Int64: // 接口定义返回一个内置int型的变量
		return Int64, nil
	case reflect.Int, reflect.Int32:
		return Int32, nil
	case reflect.Int16:
		return Int16, nil
	case reflect.Int8:
		return Int8, nil

	case reflect.Uint64:
		return Uint64, nil
	case reflect.Uint, reflect.Uint32:
		return Uint32, nil
	case reflect.Uint16:
		return Uint16, nil
	case reflect.Uint8:
		return Uint8, nil

	case reflect.Float32: // 接口定义返回一个内置float型的变量
		return Float32, nil
	case reflect.Float64, reflect.Complex64, reflect.Complex128:
		return Float64, nil

	case reflect.Map: // 此类型不应用于接口传参
		return Mapping, nil

	default: // 非基本类型, 此处String不重要
		return String, errors.New("unknown type")
	}
}

// arrayToModel 从数组类型中生成数据模型, 支持数组字段嵌套数组
// @param  rt  reflect.Type  数组类型
func arrayToModel(rt reflect.Type) *RModel {
	// 直接生成子元素的文档描述，并挂载到 modelsDocMap
	elemType := rt.Elem()
	if elemType.Kind() == reflect.Pointer {
		elemType = elemType.Elem()
	}
	return innerStructSchema(elemType)
}

// innerStructSchema 递归生成对象的文档树并直接添加到 modelsDocMap
// @param   rt  reflect.Type  数组内部元素的类型  或  结构体的字段类型
// @return  RouteModelIface 返回模型的关联文档
func innerStructSchema(rt reflect.Type) *RModel {
	if rt.Kind() == reflect.Pointer { // 结构体字段为指针类型
		rt = rt.Elem()
	}

	rm := &RModel{Name: rt.String(), Description: rt.String()}

	// 此类型为基本数据类型，或递归调用的字段为一个基本类型
	bt, err := basicTypeToModel(rt)
	if err == nil {
		modelField := &RModelField{Name: rt.String(), Tag: bt.Tag, Type: bt.Type}
		rm.Fields = append(rm.Fields, modelField)
		return rm
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if unicode.IsLower(rune(field.Name[0])) {
			continue
		}

		// 处理匿名结构体
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			// 跳过空匿名结构体，比如 BaseModel
			if field.Type.NumField() == 0 {
				continue
			}
		}

		// 仅导出字段可用
		modelField := &RModelField{Name: field.Name, Tag: field.Tag, Type: reflectKindToName(field.Type.Kind())}

		bt, err := basicTypeToModel(field.Type)
		// 字段为自定义类型，直接挂载字段后并关联类型
		if err != nil {
			switch field.Type.Kind() {
			case reflect.Pointer:
				modelField.ItemRef = ModelsRefPrefix + field.Type.Elem().String()
				b := innerStructSchema(field.Type)
				AddModelDoc(&RouteModel{Model: b, Struct: nil, RetArray: true})

			case reflect.Struct:
				modelField.ItemRef = ModelsRefPrefix + field.Type.String()
				b := innerStructSchema(field.Type)
				AddModelDoc(&RouteModel{Model: b, Struct: nil, RetArray: true})

			case reflect.Array, reflect.Slice: // 字段为数组类型，递归调用
				b := innerStructSchema(field.Type.Elem())
				elemType := reflectKindToName(field.Type.Elem().Kind())

				switch elemType {
				case godantic.ObjectName, godantic.ArrayName:
					// 对于 []struct{} 类型的字段，关联其模型连接
					modelField.ItemRef = ModelsRefPrefix + field.Type.Elem().String()
				default:
					// 对于 []string 类型的字段，直接标注内部元素的基本类型
					modelField.ItemRef = elemType
				}

				// 记录模型文档
				AddModelDoc(&RouteModel{Model: b, Struct: nil, RetArray: true})
			}

		} else { // 字段为 int, string 等基本类型
			modelField.Type = bt.Type
		}

		rm.Fields = append(rm.Fields, modelField)
	}
	return rm
}

// reflectKindToName 转换reflect.Kind为swagger类型说明
// @param  ReflectKind  reflect.Kind  反射类型
func reflectKindToName(kind reflect.Kind) (name string) {
	switch kind {

	case reflect.Array, reflect.Slice, reflect.Chan:
		name = godantic.ArrayName
	case reflect.String:
		name = godantic.StringName
	case reflect.Bool:
		name = godantic.BooleanName
	default:
		if reflect.Bool < kind && kind <= reflect.Uint64 {
			name = godantic.IntegerName
		} else if reflect.Float32 <= kind && kind <= reflect.Complex128 {
			name = godantic.NumberName
		} else {
			name = godantic.ObjectName
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

func GetDefaultV(tag reflect.StructTag, swagType string) (v any) {
	defaultV := QueryFieldTag(tag, "default", "")

	if defaultV == "" {
		v = nil
	} else { // 存在默认值
		switch swagType {

		case "string":
			v = defaultV
		case godantic.IntegerName:
			v, _ = strconv.Atoi(defaultV)
		case godantic.NumberName:
			v, _ = strconv.ParseFloat(defaultV, 64)
		case godantic.BooleanName:
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
	return reflectKindToName(reflect.TypeOf(object).Kind()) == godantic.ArrayName
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

// QueryJsonName 查询字段名
func QueryJsonName(tag reflect.StructTag, undefined string) string {
	if tag == "" {
		return undefined
	}
	if v := tag.Get("json"); v != "" {
		return strings.TrimSpace(strings.Split(v, ",")[0])
	}
	return undefined
}

// AnyToBaseModel 【核心方法】用于将 <自定义类型> 等转换成 RModel 或 RModelField
// @param  object  any  数据模型，只会是  struct,  map  或  数组类型,  不应该是  **指针类型**
func AnyToBaseModel(object BaseModelIface) RouteModelIface {
	rt := reflect.TypeOf(object)

	bt, err := basicTypeToModel(rt)
	if err == nil { // 基本数据类型
		bt.Name = rt.Name()
		return bt
	}

	// 不是基本数据类型，此时只能是 struct 或 数组类型
	switch rt.Kind() {

	case reflect.Array, reflect.Slice: // 接口显式返回数组类型, 形如:[]byte
		rModel := arrayToModel(rt)
		rModel.Description = object.Doc__()
		return rModel
	case reflect.Struct: // 接口处不应该返回struct的指针类型
		rModel := innerStructSchema(rt)
		rModel.Description = object.Doc__()
		return rModel

	default: // 对于无法识别的类型，统一处理为字符串类型
		return String
	}
}

func RModelTransformer(m BaseModelIface) *RouteModel {
	var rm *RouteModel

	switch m.(type) {

	case nil:
		rm = &RouteModel{Model: String, Struct: nil, RetArray: false}

	case *RouteModel: // 接口直接返回 FlaskGo.List()
		rm = m.(*RouteModel)

	case RModelField: // 接口返回 FlaskGo 基本数据类型 Uint64 等
		rm = &RouteModel{
			Model:    m.(RModelField),
			Struct:   m,
			RetArray: false, // RModelField基本数据类型 不存在数组型
		}

	default:
		rm = &RouteModel{ // 接口返回 自定义struct, map 或 数组
			Model:    AnyToBaseModel(m),
			Struct:   m,
			RetArray: IsArray(m),
		}
	}
	AddModelDoc(rm) // 记录主模型

	return rm
}

// QModelTransformer 将表示查询参数的struct或map转换成查询参数模型
// @param  object  any  查询参数模型，此类别必须为  map[string]bool  或  struct
func QModelTransformer(object any) (m []*QModel) {
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

// StringsReverse 数组倒序, 就地修改
func StringsReverse(arr *[]string) {
	length := len(*arr)
	var temp string
	for i := 0; i < length/2; i++ {
		temp = (*arr)[i]
		(*arr)[i] = (*arr)[length-1-i]
		(*arr)[length-1-i] = temp
	}
}

// FastApiRoutePath 将 fiber.App 格式的路径转换成 FastApi 格式的路径
//
//	Example:
//	必选路径参数：
//		Input: "/api/rcst/:no"
//		Output: "/api/rcst/{no}"
//	可选路径参数：
//		Input: "/api/rcst/:no?"
//		Output: "/api/rcst/{no}"
//	常规路径：
//		Input: "/api/rcst/no"
//		Output: "/api/rcst/no"
func FastApiRoutePath(path string) string {
	paths := strings.Split(path, constant.PathSeparator) // 路径字符
	// 查找路径中的参数
	for i := 0; i < len(paths); i++ {
		if strings.HasPrefix(paths[i], constant.PathParamPrefix) {
			// 识别到路径参数
			if strings.HasSuffix(paths[i], constant.OptionalPathParamSuffix) {
				// 可选路径参数
				paths[i] = "{" + paths[i][1:len(paths[i])-1] + "}"
			} else {
				paths[i] = "{" + paths[i][1:] + "}"
			}
		}
	}

	return strings.Join(paths, constant.PathSeparator)
}
