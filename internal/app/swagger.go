package app

import (
	"github.com/Chendemo12/flaskgo/internal/core"
	"github.com/Chendemo12/flaskgo/internal/mode"
	"github.com/Chendemo12/flaskgo/internal/swag"
	"github.com/Chendemo12/functools/python"
	"net/http"
	"strings"
)

var (
	customError404Model      *swag.RouteModel = nil // 业务层面自定义的404 错误
	validationErrorModel     *swag.RouteModel = nil // 422 表单验证错误模型
	httpValidationErrorModel *swag.RouteModel = nil // 请求错误模型
)

func init() {
	validationErrorModel = swag.RModelTransformer(ValidationError{})
	customError404Model = swag.RModelTransformer(CustomError404{})
	httpValidationErrorModel = swag.RModelTransformer(HTTPValidationError{})
}

type DebugMode struct {
	Mode string `json:"mode" oneof:"prod dev testDev" Description:"调试模式"`
}

func (d DebugMode) Doc__() string { return "调试模式模型" }

func MakeDefaultRespGroup(method string, model *swag.RouteResp) []*swag.RouteResp {
	group := []*swag.RouteResp{
		{
			StatusCode: http.StatusNotFound,
			Body: &swag.RouteModel{
				Model:  swag.String,
				Struct: swag.RModelField{},
			},
		},
	}

	switch method {
	case http.MethodGet, http.MethodDelete:
		group = append(group, model)

	case http.MethodPost, http.MethodPatch, http.MethodPut:
		group = append(group, model, &swag.RouteResp{ // 422请求参数校验错误返回实例
			StatusCode: http.StatusUnprocessableEntity,
			Body:       validationErrorModel,
		})
	}

	return group
}

func makeSwaggerDocs(f *FlaskGo) {
	// 不允许创建swag文档
	if python.All(!mode.IsDebug(), core.SwaggerDisabled) {
		return
	}

	swag.AddModelDoc(validationErrorModel)
	swag.AddModelDoc(httpValidationErrorModel)
	swag.AddModelDoc(customError404Model)

	// 存储全部的路由
	routeInsGroups := make([]*swag.RouteInsGroup, 0)

	routesMap := make(map[string][]*Route)
	for _, router := range f.APIRouters() {
		for _, route := range router.Routes() {
			// 挂载 请求体模型文档
			swag.AddModelDoc(route.RequestModel)

			path := CombinePath(router.Prefix, route.RelativePath)
			group := &swag.RouteInsGroup{Path: path}
			routeInsGroups = append(routeInsGroups, group)
			routesMap[path] = append(routesMap[path], route)
		}
	}

	// 扫描注册全部的请求头和响应体模型，以及路由对象
	for path, routes := range routesMap {
		for _, route := range routes {
			ins := &swag.RouteInstance{
				Method:       route.Method,
				Path:         strings.Split(path, RouteSeparator)[0],
				Summary:      route.Summary,
				Description:  route.Description,
				Tags:         route.Tags,
				PathFields:   route.PathFields,
				QueryFields:  route.QueryModel,
				RequestModel: route.RequestModel,
				RespGroup: MakeDefaultRespGroup(route.Method, &swag.RouteResp{
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
		swag.AddPathDoc(group)
	}

	swag.RegisterSwagger(
		f.engine, f.Title(),
		f.Description(),
		f.version+" | FlaskGo: "+Version,
		map[string]string{"name": Copyright, "url": Website},
	)
}
