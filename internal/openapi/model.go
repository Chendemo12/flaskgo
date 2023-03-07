package openapi

import (
	"github.com/Chendemo12/flaskgo/internal/godantic"
	"github.com/Chendemo12/functools/helper"
)

// Contact 联系方式, 显示在 info 字段内部
// 无需重写序列化方法
type Contact struct {
	Name  string `json:"name" description:"姓名/名称"`
	Url   string `json:"url" description:"链接"`
	Email string `json:"email" description:"联系方式"`
}

// License 权利证书, 显示在 info 字段内部
// 无需重写序列化方法
type License struct {
	Name string `json:"name" description:"名称"`
	Url  string `json:"url" description:"链接"`
}

// Info 文档说明信息
// 无需重写序列化方法
type Info struct {
	Title          string  `json:"title" description:"显示在文档顶部的标题"`
	Version        string  `json:"version" description:"显示在标题右上角的程序版本号"`
	Description    string  `json:"description" description:"显示在标题下方的说明"`
	Contact        Contact `json:"contact" description:"联系方式"`
	License        License `json:"license" description:"许可证"`
	TermsOfService string  `json:"termsOfService" description:"服务条款(不常用)"`
}

// Reference 引用模型,用于模型字段和路由之间互相引用
type Reference struct {
	// 关联模型, 取值为 godantic.RefPrefix + modelName
	Name string `json:"-" description:"关联模型"`
}

func (r *Reference) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m[godantic.RefName] = godantic.RefPrefix + r.Name

	return helper.DefaultJsonMarshal(m)
}

// ComponentScheme openapi 的模型文档部分
type ComponentScheme struct {
	Name  string               `json:"name" description:"模型名称，包含包名"`
	Model godantic.SchemaIface `json:"model" description:"模型定义"`
}

// Components openapi 的模型部分
// 需要重写序列化方法
type Components struct {
	Scheme []*ComponentScheme `json:"scheme" description:"模型文档"`
}

// MarshalJSON 重载序列化方法
func (c *Components) MarshalJSON() ([]byte, error) {
	m := make(map[string]map[string]any, 0)
	for _, v := range c.Scheme {
		m[v.Name] = v.Model.Schema() // 记录根模型
		// 从 BaseModel 生成模型，处理嵌入类型
		for _, innerF := range v.Model.Metadata().InnerFields() {
			if innerM := innerF.InnerModel(); innerM != nil { // 发现子模型
				m[innerM.SchemaName()] = innerM.Schema()
			}
		}
	}

	// 记录内置错误类型文档
	m["ValidationError"] = validationErrorDefinition
	m["HTTPValidationError"] = validationErrorResponseDefinition
	// delete
	//m["CustomValidationError"] = customErrorDefinition

	return helper.DefaultJsonMarshal(map[string]any{"schemas": m})
}

// AddModel 添加一个模型文档
func (c *Components) AddModel(m godantic.SchemaIface) {
	c.Scheme = append(c.Scheme, &ComponentScheme{
		Name:  m.SchemaName(),
		Model: m,
	})
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
	Description string                   `json:"description,omitempty" description:"说明"`
	Required    bool                     `json:"required" description:"是否必须"`
	Deprecated  bool                     `json:"deprecated" description:"是否禁用"`
	Type        godantic.OpenApiDataType `json:"type" description:"数据类型"`
}

// Parameter 路径参数或者查询参数
type Parameter struct {
	ParameterBase
	Title   string          `json:"title"`
	Name    string          `json:"name" description:"名称"`
	In      ParameterInType `json:"in" description:"参数位置"`
	Default any             `json:"default,omitempty" description:"默认值"`
}

type ModelContentSchema interface {
	OType() godantic.OpenApiDataType
	Schema() map[string]any
}

// BaseModelContentSchema 适用于请求体是基本数据类型的类型
type BaseModelContentSchema struct {
	Title string                   `json:"title,omitempty" description:"标题"`
	Type  godantic.OpenApiDataType `json:"type" description:"模型类型"`
}

func (s BaseModelContentSchema) OType() godantic.OpenApiDataType { return s.Type }

func (s BaseModelContentSchema) Schema() map[string]any {
	return map[string]any{
		"title": s.Title,
		"type":  s.Type,
	}
}

// ArrayModelContentSchema 适用于请求体是数组的类型
type ArrayModelContentSchema struct {
	BaseModelContentSchema
	Items Reference `json:"items" description:"子项目"`
}

func (s ArrayModelContentSchema) Schema() map[string]any {
	return map[string]any{
		"title": s.Title,
		"type":  godantic.ArrayType,
		"items": s.Items,
	}
}

// ObjectModelContentSchema 适用于请求体是struct的类型
type ObjectModelContentSchema struct {
	BaseModelContentSchema
	Reference
}

func (s ObjectModelContentSchema) Schema() map[string]any {
	return map[string]any{
		"title":          s.Title,
		"type":           godantic.ObjectType,
		godantic.RefName: godantic.RefPrefix + s.Name,
	}
}

// RequestBody 路由 请求体模型文档
type RequestBody struct {
	Required bool              `json:"required" description:"是否必须"`
	Content  *PathModelContent `json:"content,omitempty" description:"请求体模型"`
}

// PathModelContent 路由中请求体 RequestBody 和 响应体中返回值 Responses 模型
type PathModelContent struct {
	MIMEType ApplicationMIMEType `json:"-"`
	Schema   ModelContentSchema  `json:"schema" description:"模型引用文档"`
}

// MarshalJSON 自定义序列化
func (p *PathModelContent) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m[string(p.MIMEType)] = map[string]any{"schema": p.Schema.Schema()}

	return helper.DefaultJsonMarshal(m)
}

// Response 路由返回体，包含了返回状态码，状态码说明和返回值模型
type Response struct {
	StatusCode  int               `json:"-" description:"状态码"`
	Description string            `json:"description" description:"说明"`
	Content     *PathModelContent `json:"content" description:"返回值模型"`
}

// Operation 路由HTTP方法: Get/Post/Patch/Delete 等操作方法
type Operation struct {
	Tags        []string `json:"tags" description:"路由标签"`
	Summary     string   `json:"summary" description:"摘要描述"`
	Description string   `json:"description" description:"说明"`
	OperationId string   `json:"operationId,omitempty" description:"唯一ID"` // no use, keep
	// 路径参数和查询参数, 对于路径相同，方法不同的路由来说，其查询参数可以不一样，但其路径参数都是一样的
	Parameters []*Parameter `json:"parameters,omitempty" description:"路径参数和查询参数"`
	// 请求体，通过 MakeOperationRequestBody 构建
	RequestBody *RequestBody `json:"requestBody,omitempty" description:"请求体"`
	// 响应文档，对于任一个路由，均包含2个响应实例：200 + 422， 通过函数 MakeOperationResponses 构建
	Responses  []*Response `json:"responses" description:"响应体"`
	Deprecated bool        `json:"deprecated" description:"是否禁用"`
}

// MarshalJSON 重写序列化方法，修改 Responses 和 RequestBody 字段
func (o *Operation) MarshalJSON() ([]byte, error) {
	type OperationWithResponseMap struct {
		Operation
		Responses map[int]*Response `json:"responses" description:"响应体"`
	}

	orm := OperationWithResponseMap{}
	orm.Tags = o.Tags
	orm.Summary = o.Summary
	orm.Description = o.Description
	orm.OperationId = o.OperationId
	orm.Parameters = o.Parameters
	orm.RequestBody = o.RequestBody // TODO:
	orm.Deprecated = o.Deprecated

	orm.Responses = make(map[int]*Response)
	for _, r := range o.Responses {
		orm.Responses[r.StatusCode] = r
	}

	return helper.DefaultJsonMarshal(orm)
}

// PathItem 路由选项，由于同一个路由可以存在不同的操作方法，因此此选项可以存在多个 Operation
type PathItem struct {
	Path string `json:"-" description:"请求绝对路径"` // 无需包含此字段
	// 路由下存在的多种方法, 若字段无内容，则忽略
	Get    *Operation `json:"get,omitempty" description:"GET方法"`
	Put    *Operation `json:"put,omitempty" description:"PUT方法"`
	Post   *Operation `json:"post,omitempty" description:"POST方法"`
	Patch  *Operation `json:"patch,omitempty" description:"PATCH方法"`
	Delete *Operation `json:"delete,omitempty" description:"DELETE方法"`
	Head   *Operation `json:"head,omitempty" description:"header方法"`
	Trace  *Operation `json:"trace,omitempty" description:"trace方法"`
}

// Paths openapi 的路由部分
// 需要重写序列化方法
type Paths struct {
	Paths []*PathItem
}

func (p *Paths) AddItem(item *PathItem) {
	p.Paths = append(p.Paths, item)
}

// MarshalJSON 重载序列化方法
func (p *Paths) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	for _, v := range p.Paths {
		m[v.Path] = v
	}

	return helper.DefaultJsonMarshal(m)
}

// OpenApi 模型类, 移除 FastApi 中不常用的属性
type OpenApi struct {
	Version     string      `json:"openapi" description:"Open API版本号"`
	Info        *Info       `json:"info,omitempty" description:"联系信息"`
	Components  *Components `json:"components" description:"模型文档"`
	Paths       *Paths      `json:"paths" description:"路由列表,同一路由存在多个方法文档"`
	initialized bool
	cache       []byte
}

// AddDefinition 添加一个模型文档
func (o *OpenApi) AddDefinition(model godantic.SchemaIface) *OpenApi {
	o.Components.AddModel(model)
	return o
}

// AddPathItem 添加一个路由对象
func (o *OpenApi) AddPathItem(item *PathItem) {
	// 修改路径格式
	item.Path = FastApiRoutePath(item.Path)
	o.Paths.AddItem(item)
}

// QueryPathItem 查询路由对象
func (o *OpenApi) QueryPathItem(path string) *PathItem {
	path = FastApiRoutePath(path) // 修改路径格式

	for _, item := range o.Paths.Paths {
		if item.Path == path {
			return item
		}
	}
	return nil
}

// RecreateDocs 重建Swagger 文档
func (o *OpenApi) RecreateDocs() *OpenApi {
	bs, err := helper.DefaultJson.Marshal(o)
	if err == nil {
		o.cache = bs
	}

	o.initialized = true
	return o
}

// Schema Swagger 文档, 并非完全符合 OpenApi 文档规范
func (o *OpenApi) Schema() []byte {
	if !o.initialized {
		o.RecreateDocs()
	}

	return o.cache
}
