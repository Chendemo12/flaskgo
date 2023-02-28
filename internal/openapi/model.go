package openapi

import (
	"github.com/Chendemo12/flaskgo/internal/godantic"
	"github.com/Chendemo12/functools/helper"
)

// Contact 联系方式, 显示在 info 字段内部
type Contact struct {
	Name  string `json:"name" description:"姓名/名称"`
	Url   string `json:"url" description:"链接"`
	Email string `json:"email" description:"联系方式"`
}

// License 权利证书, 显示在 info 字段内部
type License struct {
	Name string `json:"name" description:"名称"`
	Url  string `json:"url" description:"链接"`
}

// Info 文档说明信息
type Info struct {
	Title          string  `json:"title" description:"显示在文档顶部的标题"`
	Version        string  `json:"version" description:"显示在标题右上角的程序版本号"`
	Description    string  `json:"description" description:"显示在标题下方的说明"`
	Contact        Contact `json:"contact" description:"联系方式"`
	License        License `json:"license" description:"许可证"`
	TermsOfService string  `json:"termsOfService" description:"服务条款(不常用)"`
}

type RefIface interface {
	Alias() string
}

// Reference 引用模型,用于模型字段和路由之间互相引用
type Reference struct {
	// 关联模型, 取值为 ModelRefPrefix + modelName
	Ref string `json:"$ref" description:"关联模型"`
}

func (r *Reference) Schema() map[string]any {
	return map[string]any{ModelRefName: r.Ref}
}

type ParameterInType string

const (
	InQuery  ParameterInType = "query"
	InHeader ParameterInType = "header"
	InPath   ParameterInType = "path"
	InCookie ParameterInType = "cookie"
)

// ParameterBase 各种参数的基类
type ParameterBase struct {
	Description string    `json:"description" description:"说明"`
	Required    bool      `json:"required" description:"是否必须"`
	Deprecated  bool      `json:"deprecated" description:"是否禁用"`
	Schema      Reference `json:"schema" description:"模型引用信息"`
}

// Parameter 路径参数或者查询参数
type Parameter struct {
	ParameterBase
	Name string          `json:"name" description:"名称"`
	In   ParameterInType `json:"in" description:"参数位置"`
}

type RequestBodyContentSchema interface {
	OType() godantic.OpenApiDataType
	Schema() map[string]any
}

// BaseRequestBodyContentSchema 适用于请求体是基本数据类型的类型
type BaseRequestBodyContentSchema struct {
	Title string                   `json:"title" description:"标题"`
	Type  godantic.OpenApiDataType `json:"type" description:"模型类型"`
}

func (s BaseRequestBodyContentSchema) OType() godantic.OpenApiDataType { return s.Type }

func (s BaseRequestBodyContentSchema) Schema() map[string]any {
	return map[string]any{
		"title": s.Title,
		"type":  s.Type,
	}
}

// ArrayRequestBodyContentSchema 适用于请求体是数组的类型
type ArrayRequestBodyContentSchema struct {
	BaseRequestBodyContentSchema
	Items Reference `json:"items" description:"子项目"`
}

func (s ArrayRequestBodyContentSchema) Schema() map[string]any {
	return map[string]any{
		"title": s.Title,
		"type":  godantic.ArrayType,
		"items": s.Items.Schema(),
	}
}

// ObjectRequestBodyContentSchema 适用于请求体是struct的类型
type ObjectRequestBodyContentSchema struct {
	BaseRequestBodyContentSchema
	Reference
}

func (s ObjectRequestBodyContentSchema) Schema() map[string]any {
	return map[string]any{
		"title":      s.Title,
		"type":       godantic.ObjectType,
		ModelRefName: s.Ref,
	}
}

// PathDataModel 适用于路由请求体和响应体的数据模型, 通过函数 MakeDataModelContent 构建
type PathDataModel struct {
	// 形如:
	//	{
	//		"application/json": {
	//			"schema": {
	//				"title": "Response Set Svn Macs Api Rcst Network Svnmacs Post",
	//				"type": "array",
	//				"items": {
	//					"$ref": "#/components/schemas/SvnMac"
	//				}
	//			}
	//		}
	//	}
	Content map[ApplicationMIMEType]any `json:"content" description:"请求体内容"`
}

// RequestBody 路由 请求体模型文档
type RequestBody struct {
	PathDataModel
	Required bool `json:"required" description:"是否必须"`
}

type Response struct {
	PathDataModel
	Description string `json:"description" description:"说明"`
}

// Operation 路由HTTP方法: Get/Post/Patch/Delete 等操作方法
type Operation struct {
	Tags        []string `json:"tags" description:"路由标签"`
	Summary     string   `json:"summary" description:"摘要描述"`
	Description string   `json:"description" description:"说明"`
	OperationId string   `json:"operationId" description:"唯一ID"` // no use, keep
	// 路径参数和查询参数, 对于路径相同，方法不同的路由来说，其查询参数可以不一样，但其路径参数都是一样的
	Parameters []*Parameter `json:"parameters" description:"路径参数和查询参数"`
	// 请求体，通过 MakeDataModelContent 构建
	RequestBody map[ApplicationMIMEType]any `json:"requestBody" description:"请求体"`
	// 响应文档，对于任一个路由，均包含3个响应实例：200 + 422 + 500/404， 通过函数 MakeOperationResponses 构建
	Responses  map[int]*Response `json:"responses" description:"相应体"`
	Deprecated bool              `json:"deprecated" description:"是否禁用"`
}

type PathItemParameterUnion interface {
	Parameter | Reference
}

// PathItem 路由选项，由于同一个路由可以存在不同的操作方法，因此此选项可以存在多个 Operation
type PathItem struct {
	Path string `json:"path" description:"请求绝对路径"`
	// 路由下存在的多种方法
	Get    *Operation `json:"get" description:"GET方法"`
	Put    *Operation `json:"put" description:"PUT方法"`
	Post   *Operation `json:"post" description:"POST方法"`
	Patch  *Operation `json:"patch" description:"PATCH方法"`
	Delete *Operation `json:"delete" description:"DELETE方法"`
	Head   *Operation `json:"head" description:"header方法"`
	Trace  *Operation `json:"trace" description:"trace方法"`
	//Ref  *Reference `json:"ref" description:"关联的模型"` // TODO: 删除, 多余
	//Summary     string     `json:"summary" description:"摘要描述"`   // TODO: 多余
	//Description string     `json:"description" description:"说明"` // TODO: 多余
	//Servers []*Server  `json:"servers" description:"服务器配置信息"`
	// Parameters 路径参数, 对于路径相同，方法不同的路由来说，其查询参数可以不一样，但其路径参数都是一样的
	//Parameters []*Parameter `json:"parameters" description:"路径参数"`
}

func (p PathItem) Scheme() (m map[string]any) {

	return
}

type Components struct {
	Responses     map[string]RefIface
	Parameters    map[string]Reference
	RequestBodies map[string]Reference
	Headers       map[string]Reference
	Links         map[string]Reference
}

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (t Tag) Schema() map[string]string {
	return map[string]string{
		"name":        t.Name,
		"description": t.Description,
	}
}

// OpenApi 模型类
type OpenApi struct {
	Openapi     string
	Info        *Info                  `json:"info,omitempty" description:"联系信息"`
	Tags        []Tag                  `json:"tags" description:"标签"`
	Servers     map[string]string      `json:"servers" description:""`
	Definitions []godantic.SchemaIface `json:"definitions" description:"模型文档"`
	Routes      []*PathItem            `json:"routes" description:"路由列表,同一路由存在多个方法文档"`
	cache       []byte
}

// AddDefinition 添加一个模型文档
func (o *OpenApi) AddDefinition(model godantic.Iface) *OpenApi {
	o.Definitions = append(o.Definitions, model)
	return o
}

func (o *OpenApi) AddPathItem(item *PathItem) {
	o.Routes = append(o.Routes, item)
}

func (o *OpenApi) QueryPathItem(path string) *PathItem {
	for _, item := range o.Routes {
		if item.Path == path {
			return item
		}
	}
	return nil
}

func (o *OpenApi) Components() (m map[string]map[string]any) {
	schemas := make(map[string]any, len(o.Definitions)+3)

	for _, define := range o.Definitions {
		// TODO: 从 BaseModel 生成模型，处理嵌入类型
		schemas[define.SchemaName()] = define.Schema()
	}
	// 记录内置错误类型文档
	schemas["ValidationError"] = validationErrorDefinition
	schemas["HTTPValidationError"] = validationErrorResponseDefinition
	schemas["CustomValidationError"] = customErrorDefinition

	m["schemas"] = schemas
	return
}

func (o *OpenApi) Paths() (m map[string]map[string]any) {
	for i := 0; i < len(o.Routes); i++ {
		// TODO: 从 Route 生成模型
		m[o.Routes[i].Path] = o.Routes[i].Scheme()
	}
	return
}

// CreateDocs 创建文档
func (o *OpenApi) CreateDocs() map[string]any {
	output := map[string]any{"openapi": o.Openapi, "info": o.Info}
	tags := make([]map[string]string, len(o.Tags))
	if len(o.Servers) > 0 {
		output["servers"] = o.Servers
	}

	for i := 0; i < len(o.Tags); i++ {
		tags[i] = o.Tags[i].Schema()
	}

	output["tags"] = tags
	output["components"] = o.Components()
	output["paths"] = o.Paths()

	return output
}

// RecreateDocs 重建Swagger 文档
func (o *OpenApi) RecreateDocs() *OpenApi {
	bytes, err := helper.DefaultJson.Marshal(o.CreateDocs())
	if err != nil {
		o.cache = []byte{}
	} else {
		o.cache = bytes
	}

	return o
}

// Schema Swagger 文档, 并非完全符合 OpenApi 文档规范
func (o *OpenApi) Schema() []byte {
	if len(o.cache) < 1 {
		o.RecreateDocs()
	}

	return o.cache
}
