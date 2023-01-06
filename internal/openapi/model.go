// Package openapi
//
// defines.go: 自定义类型, 用于辅助自动生成swagger文档及参数校验;
//
// model.go: schema文档模型;
package openapi

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// RouteModelIface 模型接口
type RouteModelIface interface {
	BaseModelIface
	Title() string        // 显示在swagger文档上的字段标题
	String() string       // 模型的名称，格式：packageName.ModelName
	IsRequired() bool     // 检查一个字段或模型是否是必须的
	Swag() map[string]any // swagger 文档信息
	Kind() BaseModelKind  // 模型类型
}

// RModel 请求体和响应体struct模型, RequestModel & ResponseModel(不包含路径参数和查询参数)
type RModel struct {
	Name        string         // 模型名称, 非空
	Description string         // 模型描述
	Fields      []*RModelField // 字段列表
}

func (m RModel) Title() string {
	sp := strings.Split(m.Name, ".")
	StringsReverse(&sp)
	return sp[0]
}
func (m RModel) String() string      { return m.Name }
func (m RModel) IsRequired() bool    { return true }
func (m RModel) Kind() BaseModelKind { return RModelKind }
func (m RModel) Doc__() string       { return m.Description }

// Swag 生成模型的详细描述信息
//
//	// 返回值为struct模型
//
//	{
//		"description": "example.mytimeslot",
//		"properties": {
//			"control_timeslot": {
//				"description": "控制时隙编号数组",
//				"items": {
//					"type": "integer"
//				},
//				"required": false,
//				"title": "control_timeslot",
//				"type": "array"
//			},
//			"superframe_count": {
//				"description": "超帧计数",
//				"required": false,
//				"title": "superframe_count",
//				"type": "integer"
//			},
//
//		"required": [],
//		"title": "example.MyTimeslot",
//		"type": "object"
//	},
func (m RModel) Swag() (mp map[string]any) {
	required := make([]string, 0) // 必选字段
	fields := make(dict, 0)       // 属性字段
	// 字段详细描述信息中的 "required"属性 为 openapi 规范要求，用于在字段胖显式标明 ”字段必须“
	// 与此对应的为 模型详细信息中的 "required"数组为 FastApi 文档的特殊标记，用于在必须字段未填写便提交请求时，弹出 “未填提示”
	mp = dict{"title": m.Title(), "required": &required, "type": godantic.ObjectName, "properties": &fields}

	if m.Description == "" {
		mp["description"] = strings.ToLower(m.Name)
	} else {
		mp["description"] = m.Description
	}

	for _, field := range m.Fields { // object类型不存在field不合理
		// 【swagger规定】若定义了json标签，则其名称为json字段定义的值
		fieldName := QueryJsonName(field.Tag, field.Title())
		fields[fieldName] = field.Swag()
		if field.IsRequired() {
			required = append(required, fieldName)
		}
	}

	return
}

// RModelField RModel 的字段描述
type RModelField struct {
	Name        string            // 字段名称,非空
	Tag         reflect.StructTag // binding:"required" 标记一个字段是必须的
	Type        string            // swag数据类型
	ItemRef     string            // 子元素类型, 仅Type=array/object时有效
	ReflectKind reflect.Kind      // 数据类型
}

func (p RModelField) Title() string       { return p.Name }
func (p RModelField) String() string      { return p.Name }
func (p RModelField) IsRequired() bool    { return IsFieldRequired(p.Tag) }
func (p RModelField) Kind() BaseModelKind { return RFieldKind }
func (p RModelField) Doc__() string       { return QueryFieldTag(p.Tag, "description", p.String()) }

// Swag 生成字段的详细描述信息
//
//	// 字段为结构体类型
//
//	"position_sat": {
//		"$ref": "#/components/schemas/example.PositionGeo",
//		"description": "position_sat",
//		"required": false,
//		"title": "position_sat",
//		"type": "object"
//	}
//
//	// 字段为数组类型, 数组元素为基本类型
//
//	"traffic_timeslot": {
//		"description": "业务时隙编号数组",
//		"items": {
//			"type": "integer"
//		},
//		"required": false,
//		"title": "traffic_timeslot",
//		"type": "array"
//	}
//
//	// 字段为数组类型, 数组元素为自定义结构体类型
//
//	"Detail": {
//		"description": "Detail",
//		"items": {
//			"$ref": "#/components/schemas/ValidationError"
//		},
//		"required": true,
//		"title": "Detail",
//		"type": "array"
//	}
func (p RModelField) Swag() (mp map[string]any) {
	// 最基础的属性，必须
	mp = dict{
		"title":       p.Title(),
		"type":        p.Type,
		"required":    p.IsRequired(),
		"description": QueryFieldTag(p.Tag, "description", p.Name),
	}
	// 生成默认值
	if v := GetDefaultV(p.Tag, p.Type); v != nil {
		mp["default"] = v
	}
	// 生成字段的枚举值
	if es := QueryFieldTag(p.Tag, "oneof", ""); es != "" {
		mp["enum"] = strings.Split(es, " ")
	}

	// 为不同的字段类型生成相应的描述
	switch p.Type {

	case godantic.IntegerName, godantic.NumberName: // 生成数字类型的最大最小值
		if lt := QueryFieldTag(p.Tag, "lt", ""); lt != "" {
			mp["maximum"], _ = strconv.Atoi(lt)
		}
		if gt := QueryFieldTag(p.Tag, "gt", ""); gt != "" {
			mp["minimum"], _ = strconv.Atoi(gt)
		}

		if lt := QueryFieldTag(p.Tag, "lte", ""); lt != "" {
			mp["exclusiveMaximum"], _ = strconv.Atoi(lt)
		}
		if gt := QueryFieldTag(p.Tag, "gte", ""); gt != "" {
			mp["exclusiveMinimum"], _ = strconv.Atoi(gt)
		}
		// 存在多个标记
		if lt := QueryFieldTag(p.Tag, "max", ""); lt != "" {
			mp["maximum"], _ = strconv.Atoi(lt)
		}
		if gt := QueryFieldTag(p.Tag, "min", ""); gt != "" {
			mp["minimum"], _ = strconv.Atoi(gt)
		}

	case godantic.StringName: // 生成字符串类型的最大最小长度
		if lt := QueryFieldTag(p.Tag, "max", ""); lt != "" {
			mp["maxLength"], _ = strconv.Atoi(lt)
		}
		if gt := QueryFieldTag(p.Tag, "min", ""); gt != "" {
			mp["minLength"], _ = strconv.Atoi(gt)
		}

	case godantic.ArrayName:
		// 为数组类型生成子类型描述
		if p.ItemRef != "" {
			if strings.HasPrefix(p.ItemRef, ModelsRefPrefix) { // 数组子元素为关联类型
				mp["items"] = map[string]string{"$ref": p.ItemRef}
			} else { // 子元素为基本数据类型
				mp["items"] = map[string]string{"type": p.ItemRef}
			}
		} else { // 缺省为string
			mp["items"] = map[string]string{"type": godantic.StringName}
		}
		// 限制数组的长度
		if lt := QueryFieldTag(p.Tag, "max", ""); lt != "" {
			mp["maxLength"], _ = strconv.Atoi(lt)
		}
		if gt := QueryFieldTag(p.Tag, "min", ""); gt != "" {
			mp["minLength"], _ = strconv.Atoi(gt)
		}

	case godantic.ObjectName:
		if p.ItemRef != "" { // 字段类型为自定义结构体，生成关联类型，此内部结构体已注册
			mp["$ref"] = p.ItemRef
		}

	default:
	}

	return
}

// QModel 查询参数, QModelForm 的详细字段
type QModel struct {
	Name     string            // 字段名称
	Tag      reflect.StructTag // binding:"required" 标记一个字段是必须的
	Required bool              // 是否必须
	InPath   bool              // 是否是路径参数
}

func (q QModel) Title() string       { return q.Name }
func (q QModel) String() string      { return q.Name }
func (q QModel) IsRequired() bool    { return q.Required }
func (q QModel) Kind() BaseModelKind { return QModelKind }
func (q QModel) Doc__() string       { return QueryFieldTag(q.Tag, "description", q.String()) }

// Swag
//
//	{
//		"required": true,
//		"schema": {
//			"title": "Name",
//			"type": "string",
//			"default": "jack"
//		},
//		"name": "name",
//		"in": "query"/"path"
//	}
func (q QModel) Swag() (mp map[string]any) {
	if q.InPath { // 作为路径参数
		mp = dict{
			"name":     q.Name,
			"schema":   dict{"title": q.Title(), "type": godantic.StringName},
			"required": q.Required,
			"in":       "path",
		}
	} else {
		mp = dict{
			"name":        QueryJsonName(q.Tag, q.Title()),
			"schema":      dict{"title": q.Title(), "type": godantic.StringName},
			"required":    q.IsRequired(),
			"description": QueryFieldTag(q.Tag, "description", q.Name),
			"in":          "query",
		}

		// 获取默认值
		if v := QueryFieldTag(q.Tag, "default", ""); v != "" {
			mp["default"] = v
		}
	}

	return
}

// RouteModel 复合类型，包含了请求体、响应体、或查询参数模型
type RouteModel struct {
	Model    RouteModelIface // 数据模型 RModel , RModelField 或 QModel 其具体类型由Kind决定
	Struct   BaseModelIface  // 原始数据类型
	RetArray bool            // 是否为数组类型, 用于标识是否生成关于内部模型的数组类型
}

func (r RouteModel) Title() string            { return r.Model.Title() }
func (r RouteModel) FullName() string         { return r.Model.String() }
func (r RouteModel) IsArray() bool            { return r.RetArray }
func (r RouteModel) Doc__() string            { return r.Model.Doc__() }
func (r RouteModel) Doc() (mp map[string]any) { return r.Model.Swag() }

// Schema 关联数据模型，支持数组类型
func (r RouteModel) Schema() (mp map[string]any) {
	if r.RetArray {
		mp = dict{
			"type":  godantic.ArrayName,
			"items": map[string]string{"$ref": ModelsRefPrefix + r.Model.String()},
		}
	} else {
		mp = dict{"$ref": ModelsRefPrefix + r.Model.String()}
	}
	return
}

// RouteResp 接口响应完整实例
type RouteResp struct {
	Body       *RouteModel // 响应体模型, 此模型恒 != nil
	StatusCode int         // 状态码
}

func (r RouteResp) Title() string       { return r.Body.Title() }
func (r RouteResp) String() string      { return r.Body.FullName() }
func (r RouteResp) IsRequired() bool    { return true }
func (r RouteResp) Kind() BaseModelKind { return RouteRespKind }

// Doc 路由响应体模型
//
//	{
//		"200": {
//	    	"description": "OK",
//	        "content": {
//	        	"application/json": {
//	           		"schema": {
//	            		"title": "string",
//	                	"type": "string",
//	                	"required": false,
//	                	"description": "string",
//	                	"maxLength": 255,
//	                	"minLength": 0
//	        		}
//	    		}
//			}
//		}
//	}
//
//	{
//		"422": {
//			"description": "Validation Error",
//			"content": {
//				"application/json": {
//					"schema": {
//						"$ref": "#/components/schemas/HTTPValidationError"
//					}
//				}
//			}
//		}
//	}
func (r RouteResp) Doc() (mp map[string]any) {
	if r.Body != nil {
		mp = MakeResponseSchema(
			fiber.MIMEApplicationJSON,
			http.StatusText(r.StatusCode),
			r.Body.Schema(),
		)
	} else {
		mp = dict{}
	}

	return
}

// RouteInstance 路由实例
// 一个路由实例是包含了路由信息、参数模型（查询参数）、响应模型和响应状态码的描述集合
type RouteInstance struct {
	Method       string       // 请求方法
	Path         string       // 路由绝对路径
	Summary      string       // 路由摘要信息
	Description  string       // 路由详细描述
	Tags         []string     // 标签组
	RequestModel *RouteModel  // 请求体,只能有一个, 此模型恒 != nil
	PathFields   []*QModel    // 路径参数
	QueryFields  []*QModel    // 查询参数
	RespGroup    []*RouteResp // 接口全部响应实例
}

func (s RouteInstance) Title() string       { return s.Path }
func (s RouteInstance) String() string      { return s.Path }
func (s RouteInstance) IsRequired() bool    { return true }
func (s RouteInstance) Kind() BaseModelKind { return RouteInsKind }
func (s RouteInstance) Doc__() string       { return s.Description }

// Swag 生成 schema 文档
//
// 对于 QueryFields 和 RequestModel 非必须，RespGroup 则是必须的，且 RequestModel 对于 GET/DELETE 方法无效
//
// 通常 GET/DELETE 接口示例如下，不包含请求体：
//
//	{
//	 "tags": [
//	   "Debug"
//	 ],
//	 "summary": "获取调试模式开关",
//	 "operationId": "read_debug_switch_api_rcst_debug_switch_get",
//	 "parameters": [
//	   {
//	     "type": "string",
//	     "required": true,
//	     "description": "姓名",
//	     "name": "name",
//	     "in": "query",
//	     "title": "Name"
//	   },
//	   {
//	     "name": "age",
//	     "in": "query",
//	     "title": "Age",
//	     "type": "string",
//	     "required": false,
//	     "description": "年龄",
//	     "default": "23"
//	   }
//	 ],
//	 "responses": {
//	   "200": {
//	     "description": "Successful Response",
//	     "content": {
//	       "application/json": {
//	         "schema": {
//	           "title": "Response Read Debug Switch Api Rcst Debug Switch Get",
//	           "type": "boolean"
//	         }
//	       }
//	     }
//	   },
//	   "404": {
//	     "description": "Not Found",
//	     "content": {
//	       "application/json": {
//	         "schema": {
//	           "title": "string",
//	           "type": "string",
//	           "required": false,
//	           "description": "string",
//	           "maxLength": 255,
//	           "minLength": 0
//	         }
//	       }
//	     }
//	   }
//	 }
//	}
//
// 对于 RequestModel 的渲染示例：其中 "requestBody" 与 "responses" 处于同一级别
//
//	{
//	 "requestBody": {
//	   "content": {
//	     "application/json": {
//	       "schema": {
//	         "$ref": "#/components/schemas/TDMParams"
//	       }
//	     }
//	   },
//	   "required": true
//	 }
//	}
func (s RouteInstance) Swag() (mp map[string]any) {
	mp = dict{"tags": s.Tags, "summary": s.Summary, "description": s.Description}
	required := make([]string, 0) // 用于标明 查询参数或路径参数请求体是否必须【FastApi】

	// 请求参数, 包含路径参数和查询参数
	parameters := make([]dict, 0)
	// 生成路径参数
	for i := 0; i < len(s.PathFields); i++ {
		parameters = append(parameters, s.PathFields[i].Swag())
		if s.PathFields[i].IsRequired() { // 标明路径参数必须
			required = append(required, s.PathFields[i].Name)
		}
	}

	// 生成query参数
	for i := 0; i < len(s.QueryFields); i++ {
		parameters = append(parameters, s.QueryFields[i].Swag())
		if s.QueryFields[i].IsRequired() {
			required = append(required, s.QueryFields[i].Name)
		}
	}

	// 对于请求体和响应体，数据模型处已生成详细Schema，此处仅需关联模型即可
	if s.RequestModel.Model != nil && s.Method != http.MethodGet && s.Method != http.MethodDelete { // POST, PATCH and PUT methods
		mp["requestBody"] = MakeResponseSchema(
			fiber.MIMEApplicationJSON,
			s.RequestModel.FullName(),
			s.RequestModel.Schema(), // 关联数据模型
		)
	}

	responses := dict{} // 响应实例
	// 生成响应模型
	for _, r := range s.RespGroup {
		responses[strconv.Itoa(r.StatusCode)] = r.Doc()
	}

	mp["parameters"] = parameters
	mp["required"] = required
	// responses 是必须存在的
	mp["responses"] = responses

	return
}

// RouteInsGroup 同一个路由可以有多个不同方法的路由实例
type RouteInsGroup struct {
	Path     string           // 绝对路由
	InsArray []*RouteInstance // 绝对路由：[]路由实例
}

func (s RouteInsGroup) Title() string       { return s.Path }
func (s RouteInsGroup) String() string      { return s.Path }
func (s RouteInsGroup) IsRequired() bool    { return true }
func (s RouteInsGroup) Kind() BaseModelKind { return RouteInsGroupKind }
func (s RouteInsGroup) Doc__() string       { return s.Path }

func (s RouteInsGroup) Swag() (mp map[string]any) {
	mp = dict{}
	for _, ins := range s.InsArray {
		mp[strings.ToLower(ins.Method)] = ins.Swag()
	}

	return
}

// List 数组类型
// @param  m  any  数组内部元素
func List(m BaseModelIface) *RouteModel {

	switch m.(type) {
	case nil:
		return &RouteModel{Model: String, Struct: nil, RetArray: true}

	case RModelField:
		return &RouteModel{
			Model:    m.(RModelField),
			Struct:   m,
			RetArray: true,
		}

	default:
		return &RouteModel{ // m 为 常量, map, struct, array, slice 等
			Model:    AnyToBaseModel(m),
			Struct:   m,
			RetArray: true,
		}
	}
}

// MakeResponseSchema 响应体文档
//
//	{
//		"description": "Validation Error",
//		"content": {
//			"application/json": {
//				"schema": {
//					"$ref": "#/components/schemas/HTTPValidationError"
//				}
//			}
//		}
//	}
func MakeResponseSchema(mimeType, description string, schemaMap dict) map[string]any {
	return dict{
		"description": description,
		"content": dict{
			mimeType: dict{
				"schema": schemaMap,
			},
		},
	}
}
