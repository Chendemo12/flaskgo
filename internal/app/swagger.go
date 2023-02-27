package app

import (
	"github.com/Chendemo12/flaskgo/internal/core"
	"github.com/Chendemo12/flaskgo/internal/godantic"
	"github.com/Chendemo12/flaskgo/internal/mode"
	"github.com/Chendemo12/flaskgo/internal/openapi"
	"github.com/Chendemo12/functools/python"
	"net/http"
	"strings"
)

type DebugMode struct {
	godantic.BaseModel
	Mode string `json:"mode" oneof:"prod dev testDev" Description:"调试模式"`
}

func (d DebugMode) SchemaDesc() string { return "调试模式模型" }

func MakeDefaultRespGroup(method string, model *openapi.RouteResp) []*openapi.RouteResp {
	group := []*openapi.RouteResp{
		{
			StatusCode: http.StatusNotFound,
			Body: &openapi.RouteModel{
				Model:  openapi.String,
				Struct: openapi.RModelField{},
			},
		},
	}

	switch method {
	case http.MethodGet, http.MethodDelete:
		group = append(group, model)

	case http.MethodPost, http.MethodPatch, http.MethodPut:
		group = append(group, model, &openapi.RouteResp{ // 422请求参数校验错误返回实例
			StatusCode: http.StatusUnprocessableEntity,
			Body:       validationErrorModel,
		})
	}

	return group
}

func makeSwaggerDocs(f *FlaskGo) {
	f.service.openApi = openapi.NewOpenApi(f.Title(), f.Version(), f.Description())

	// 挂载模型文档
	for _, router := range f.APIRouters() {
		for _, route := range router.Routes() {
			// 挂载 请求体模型文档
			openapi.AddModelDoc(route.RequestModel.SchemaName(), route.RequestModel.Schema())
			openapi.AddModelDoc(route.ResponseModel.SchemaName(), route.ResponseModel.Schema())

			// 挂载内部模型
			for name, model := range route.RequestModel.InnerSchema() {
				openapi.AddModelDoc(name, model)
			}
			for name, model := range route.ResponseModel.InnerSchema() {
				openapi.AddModelDoc(name, model)
			}
		}
	}

	for _, router := range f.APIRouters() {
		for _, route := range router.Routes() {
			path := CombinePath(router.Prefix, route.RelativePath)
			openapi.AddPathDoc(path, route)
		}
	}

	// 存储全部的路由
	routeInsGroups := make([]*openapi.RouteInsGroup, 0)

	routesMap := make(map[string][]*Route)
	for _, router := range f.APIRouters() {
		for _, route := range router.Routes() {
			// 挂载 请求体模型文档
			openapi.AddModelDoc(route.RequestModel.SchemaName(), route.RequestModel.Schema())

			path := CombinePath(router.Prefix, route.RelativePath)
			group := &openapi.RouteInsGroup{Path: path}
			routeInsGroups = append(routeInsGroups, group)
			routesMap[path] = append(routesMap[path], route)
		}
	}

	// 扫描注册全部的请求头和响应体模型，以及路由对象
	for path, routes := range routesMap {
		for _, route := range routes {
			ins := &openapi.RouteInstance{
				Method:       route.Method,
				Path:         strings.Split(path, RouteSeparator)[0],
				Summary:      route.Summary,
				Description:  route.Description,
				Tags:         route.Tags,
				PathFields:   route.PathFields,
				QueryFields:  route.QueryFields,
				RequestModel: route.RequestModel,
				RespGroup: MakeDefaultRespGroup(route.Method, &openapi.RouteResp{
					Body:       route.ResponseModel,
					StatusCode: http.StatusOK,
				}),
			}

			for i := 0; i < len(routeInsGroups); i++ {
				if routeInsGroups[i].Path == path {
					routeInsGroups[i].InsArray = append(routeInsGroups[i].InsArray, ins)
				}
			}
		}
	}

	for _, group := range routeInsGroups {
		openapi.AddPathDoc(group.Path, group.Swag())
	}

	// 不允许创建swag文档
	if python.All(!mode.IsDebug(), core.SwaggerDisabled) {
		return
	}

	openapi.RegisterSwagger(
		f.engine, f.Title(),
		f.Description(),
		f.version+" | FlaskGo: "+Version,
		map[string]string{"name": Copyright, "url": Website},
	)
}
