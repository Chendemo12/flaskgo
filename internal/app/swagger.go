package app

import (
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

	f.createDefines()
	f.createPaths()
	f.createSwaggerRoutes()
}

// 注册 swagger 的文档路由
func (f *FlaskGo) createSwaggerRoutes() {
	// docs 在线调试页面
	f.engine.Get("/docs", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(openapi.MakeSwaggerUiHtml(
			f.title,
			openapi.OpenapiUrl,
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
			openapi.OpenapiUrl,
			openapi.RedocJsUrl,
			openapi.RedocFaviconUrl,
		))
	})

	// openapi 获取路由定义
	f.engine.Get("/openapi.json", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		//return c.Status(fiber.StatusOK).SendString(templateString)
		return c.Status(fiber.StatusOK).JSON(f.service.openApi.Schema())
	})
}

// 生成模型定义
func (f *FlaskGo) createDefines() {
	for _, router := range f.APIRouters() {
		for _, route := range router.Routes() {
			if route.RequestModel != nil {
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
			ab := route.Path(router.Prefix)
			// 存在相同路由，不同方法的路由选项
			item := f.service.openApi.QueryPathItem(ab)
			if item == nil {
				item = &openapi.PathItem{
					Path:        ab,
					Ref:         nil, // TODO: what is the ref
					Summary:     route.Summary,
					Description: route.Description,
					Get:         nil,
					Put:         nil,
					Post:        nil,
					Patch:       nil,
					Delete:      nil,
					Head:        nil,
					Trace:       nil,
					Servers:     nil,
					Parameters:  make([]*openapi.Parameter, len(route.PathFields)),
				}
				f.service.openApi.AddPathItem(item)
			}

			// 构造路径参数
			for no, q := range route.PathFields {
				item.Parameters[no].Name = q.Name
				item.Parameters[no].In = openapi.InPath
				item.Parameters[no].Description = q.SchemaName()
				item.Parameters[no].Deprecated = route.deprecated
				item.Parameters[no].Required = q.IsRequired()
			}

			// 构造查询参数
			queryParams := make([]*openapi.Parameter, len(route.QueryFields))
			for no, q := range route.QueryFields {
				queryParams[no].Name = q.Name
				queryParams[no].In = openapi.InQuery
				queryParams[no].Description = q.SchemaName()
				queryParams[no].Deprecated = route.deprecated
				queryParams[no].Required = q.IsRequired()
			}

			// 构造操作符
			operation := &openapi.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				OperationId: "", // TODO: keep
				Parameters:  queryParams,
				RequestBody: &openapi.RequestBody{},
				Response:    openapi.NewResponse(route.ResponseModel.String()), // TODO: ?
				Deprecated:  route.deprecated,
				Servers:     nil,
			}
			operation.RequestBody.Deprecated = route.deprecated
			operation.RequestBody.Description = route.RequestModel.SchemaDesc()
			operation.RequestBody.Required = route.RequestModel.IsRequired()
			operation.RequestBody.Content = route.RequestModel.Ref()

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
	}
}

//func MakeDefaultRespGroup(method string, model *openapi.RouteResp) []*openapi.RouteResp {
//	group := []*openapi.RouteResp{
//		{
//			StatusCode: http.StatusNotFound,
//			Body: &openapi.RouteModel{
//				Model:  openapi.String,
//				Struct: openapi.RModelField{},
//			},
//		},
//	}
//
//	switch method {
//	case http.MethodGet, http.MethodDelete:
//		group = append(group, model)
//
//	case http.MethodPost, http.MethodPatch, http.MethodPut:
//		group = append(group, model, &openapi.RouteResp{ // 422请求参数校验错误返回实例
//			StatusCode: http.StatusUnprocessableEntity,
//			Body:       validationErrorModel,
//		})
//	}
//
//	return group
//}
//
//func makeSwaggerDocs(f *FlaskGo) {
//	f.service.openApi = openapi.NewOpenApi(f.Title(), f.Version(), f.Description())
//
//	// 挂载模型文档
//	for _, router := range f.APIRouters() {
//		for _, route := range router.Routes() {
//			// 挂载 请求体模型文档
//			openapi.AddModelDoc(route.RequestModel.SchemaName(), route.RequestModel.Schema())
//			openapi.AddModelDoc(route.ResponseModel.SchemaName(), route.ResponseModel.Schema())
//
//			// 挂载内部模型
//			for name, model := range route.RequestModel.InnerSchema() {
//				openapi.AddModelDoc(name, model)
//			}
//			for name, model := range route.ResponseModel.InnerSchema() {
//				openapi.AddModelDoc(name, model)
//			}
//		}
//	}
//
//	for _, router := range f.APIRouters() {
//		for _, route := range router.Routes() {
//			path := CombinePath(router.Prefix, route.RelativePath)
//			openapi.AddPathDoc(path, route)
//		}
//	}
//
//	// 存储全部的路由
//	routeInsGroups := make([]*openapi.RouteInsGroup, 0)
//
//	routesMap := make(map[string][]*Route)
//	for _, router := range f.APIRouters() {
//		for _, route := range router.Routes() {
//			// 挂载 请求体模型文档
//			openapi.AddModelDoc(route.RequestModel.SchemaName(), route.RequestModel.Schema())
//
//			path := CombinePath(router.Prefix, route.RelativePath)
//			group := &openapi.RouteInsGroup{Path: path}
//			routeInsGroups = append(routeInsGroups, group)
//			routesMap[path] = append(routesMap[path], route)
//		}
//	}
//
//	// 扫描注册全部的请求头和响应体模型，以及路由对象
//	for path, routes := range routesMap {
//		for _, route := range routes {
//			ins := &openapi.RouteInstance{
//				Method:       route.Method,
//				Path:         strings.Split(path, RouteSeparator)[0],
//				Summary:      route.Summary,
//				Description:  route.Description,
//				Tags:         route.Tags,
//				PathFields:   route.PathFields,
//				QueryFields:  route.QueryFields,
//				RequestModel: route.RequestModel,
//				RespGroup: MakeDefaultRespGroup(route.Method, &openapi.RouteResp{
//					Body:       route.ResponseModel,
//					StatusCode: http.StatusOK,
//				}),
//			}
//
//			for i := 0; i < len(routeInsGroups); i++ {
//				if routeInsGroups[i].Path == path {
//					routeInsGroups[i].InsArray = append(routeInsGroups[i].InsArray, ins)
//				}
//			}
//		}
//	}
//
//	for _, group := range routeInsGroups {
//		openapi.AddPathDoc(group.Path, group.Swag())
//	}
//
//	// 不允许创建swag文档
//	if python.All(!mode.IsDebug(), core.SwaggerDisabled) {
//		return
//	}
//
//	openapi.RegisterSwagger(
//		f.engine, f.Title(),
//		f.Description(),
//		f.version+" | FlaskGo: "+Version,
//		map[string]string{"name": Copyright, "url": Website},
//	)
//}
