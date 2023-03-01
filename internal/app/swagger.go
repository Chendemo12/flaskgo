package app

import (
	"bytes"
	"github.com/Chendemo12/flaskgo/internal/core"
	"github.com/Chendemo12/flaskgo/internal/godantic"
	"github.com/Chendemo12/flaskgo/internal/mode"
	"github.com/Chendemo12/flaskgo/internal/openapi"
	"github.com/Chendemo12/functools/python"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type DebugMode struct {
	godantic.BaseModel
	Mode string `json:"mode" oneof:"prod dev testDev" Description:"调试模式"`
}

func (d DebugMode) SchemaDesc() string { return "调试模式模型" }

func (f *FlaskGo) createOpenApiDoc() {
	// 不允许创建swag文档
	if python.All(!mode.IsDebug(), core.SwaggerDisabled) {
		return
	}

	f.service.openApi = openapi.NewOpenApi(f.title, f.version, f.Description())
	f.createDefines()
	//f.createPaths()
	f.createSwaggerRoutes()
}

// 注册 swagger 的文档路由
func (f *FlaskGo) createSwaggerRoutes() {
	// docs 在线调试页面
	f.engine.Get("/docs", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(openapi.MakeSwaggerUiHtml(
			f.title,
			openapi.JsonUrl,
			openapi.SwaggerJsUrl,
			openapi.SwaggerCssUrl,
			openapi.SwaggerFaviconUrl,
		))
	})

	// redoc 纯文档页面
	f.engine.Get("/redoc", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(openapi.MakeRedocUiHtml(
			f.title,
			openapi.JsonUrl,
			openapi.RedocJsUrl,
			openapi.RedocFaviconUrl,
		))
	})

	// openapi 获取路由定义
	f.engine.Get("/openapi.json", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		//return c.Status(fiber.StatusOK).SendString(templateString)
		return c.Status(fiber.StatusOK).SendStream(bytes.NewReader(f.service.openApi.Schema()))
	})
}

// 生成模型定义
func (f *FlaskGo) createDefines() {
	for _, router := range f.APIRouters() {
		for _, route := range router.Routes() {
			if route.RequestModel != nil {
				// 内部会处理嵌入类型
				f.service.openApi.AddDefinition(route.RequestModel)
			}
			if route.ResponseModel != nil {
				f.service.openApi.AddDefinition(route.ResponseModel)
			}
		}
	}
}

// 生成路由定义
func (f *FlaskGo) createPaths() {
	for _, router := range f.APIRouters() {
		for _, route := range router.Routes() {
			routeToPathItem(router, route, f.service.openApi)
		}
	}
}

func routeToPathItem(router *Router, route *Route, api *openapi.OpenApi) {
	ab := route.Path(router.Prefix)
	// 存在相同路径，不同方法的路由选项
	item := api.QueryPathItem(ab)
	if item == nil {
		item = &openapi.PathItem{
			Path:   ab,
			Get:    nil,
			Put:    nil,
			Post:   nil,
			Patch:  nil,
			Delete: nil,
			Head:   nil,
			Trace:  nil,
		}
		api.AddPathItem(item)
	}

	// 构造路径参数
	pathParams := make([]*openapi.Parameter, len(route.PathFields))
	for no, q := range route.PathFields {
		p := &openapi.Parameter{
			ParameterBase: openapi.ParameterBase{
				Description: q.SchemaDesc(),
				Required:    q.IsRequired(),
				Deprecated:  route.deprecated,
				Schema: openapi.Reference{
					Name: q.SchemaName(),
				},
			},
			Name: q.Name,
			In:   openapi.InPath,
		}

		pathParams[no] = p
	}

	println(python.Repr(pathParams[0]))

	// 构造查询参数
	queryParams := make([]*openapi.Parameter, len(route.QueryFields))
	for no, q := range route.QueryFields {
		p := &openapi.Parameter{
			ParameterBase: openapi.ParameterBase{
				Description: q.SchemaDesc(),
				Required:    q.IsRequired(),
				Deprecated:  route.deprecated,
				Schema: openapi.Reference{
					Name: q.SchemaName(),
				},
			},
			Name: q.Name,
			In:   openapi.InQuery,
		}
		queryParams[no] = p
	}

	// 构造操作符
	operation := &openapi.Operation{
		Summary:     route.Summary,
		Description: route.Description,
		Tags:        route.Tags,
		Parameters:  append(pathParams, queryParams...),
		RequestBody: openapi.MakeOperationRequestBody(route.RequestModel),
		Responses:   openapi.MakeOperationResponses(route.ResponseModel),
		Deprecated:  route.deprecated,
	}

	// 绑定到操作方法
	switch route.Method {

	case http.MethodPost:
		item.Post = operation
	case http.MethodPut:
		item.Put = operation
	case http.MethodDelete:
		item.Delete = operation
	case http.MethodPatch:
		item.Patch = operation
	case http.MethodHead:
		item.Head = operation
	case http.MethodTrace:
		item.Trace = operation

	default:
		item.Get = operation
	}
}
