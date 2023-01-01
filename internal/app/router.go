package app

import (
	"github.com/Chendemo12/flaskgo/internal/constant"
	"github.com/Chendemo12/flaskgo/internal/swag"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"reflect"
	"strings"
)

// RouteSeparator 路由分隔符，用于分割路由方法和路径
const RouteSeparator = "_0#0_"

var ( // 记录创建的路由对象，用于其后的请求和响应校验
	MethodGetRoutes    = map[string]*Route{}
	MethodPostRoutes   = map[string]*Route{}
	MethodPatchRoutes  = map[string]*Route{}
	MethodPutRoutes    = map[string]*Route{}
	MethodDeleteRoutes = map[string]*Route{}
)

// APIRouter 创建一个路由组
func APIRouter(prefix string, tags []string) *Router {
	fgr := &Router{
		Prefix:     prefix,
		Tags:       tags,
		deprecated: false,
	}
	fgr.routes = make(map[string]*Route) // 初始化map,并保证为空
	return fgr
}

// Route 一个完整的路由对象，此对象会在程序启动时生成swagger文档
// 其中相对路径Path不能重复，否则后者会覆盖前者
type Route struct {
	RequestModel  *swag.RouteModel // 请求体模型, 此模型恒 != nil
	ResponseModel *swag.RouteModel // 响应模型,路由参数, 此模型恒 != nil
	Method        string           // 请求方法
	RelativePath  string           // 请求相对路由, 必定以/开头,路由参数
	Summary       string           // 路由摘要,路由参数
	Description   string           // 路由详细描述
	Tags          []string
	PathFields    []*swag.QModel  // 路径参数
	QueryModel    []*swag.QModel  // 查询参数
	Handlers      []fiber.Handler // 路由处理钩子
	Dependencies  []HandlerFunc
	deprecated    bool // 是否禁用此路由
}

func (f *Route) LowerMethod() string { return strings.ToLower(f.Method) }

// Deprecate 禁用路由
func (f *Route) Deprecate() *Route {
	f.deprecated = true
	return f
}

// AddDependency 添加依赖项，用于在执行路由函数前执行一个自定义操作，此操作作用于参数校验通过之后
// @param  fcs  HandlerFunc  依赖项
func (f *Route) AddDependency(fcs ...HandlerFunc) *Route {
	if len(fcs) > 0 {
		f.Dependencies = append(f.Dependencies, fcs...)
	}
	return f
}

// SetDescription 设置一个路由的详细描述信息
// @param  Description  string  详细描述信息
func (f *Route) SetDescription(description string) *Route {
	f.Description = description
	return f
}

// SetQueryParams 设置查询参数,此空struct的每一个字段都将作为一个单独的查询参数
// @param  m  any  查询参数对象
func (f *Route) SetQueryParams(m any) *Route {
	f.QueryModel = swag.QModelTransformer(m) // 转换为内部模型
	return f
}

// SetRequestModel 设置请求体对象,此model应为一个空struct实例,而非指针类型,且仅"GET",http.MethodDelete有效
// @param  m  any  请求体对象
func (f *Route) SetRequestModel(m swag.BaseModelIface) *Route {
	if f.Method != http.MethodGet && f.Method != http.MethodDelete {
		f.RequestModel = swag.RModelTransformer(m)
	}
	return f
}

// Path 合并路由
// @param  prefix  string  路由组前缀
func (f *Route) Path(prefix string) string { return CombinePath(prefix, f.RelativePath) }

// Router 一个独立的路由组，Prefix路由组前缀，其内部的子路由均包含此前缀
type Router struct {
	routes     map[string]*Route
	Prefix     string
	Tags       []string
	deprecated bool
}

// Routes 获取路由组内部定义的全部子路由信息
func (f *Router) Routes() map[string]*Route { return f.routes }

// Deprecate 禁用整个路由组路由
func (f *Router) Deprecate() *Router {
	f.deprecated = true
	return f
}

// Activate 激活整个路由组路由
func (f *Router) Activate() *Router {
	f.deprecated = false
	return f
}

// IncludeRouter 挂载一个子路由组,目前仅支持在子路由组初始化后添加
// @param  router  *Router  子路由组
func (f *Router) IncludeRouter(router *Router) *Router {
	for _, route := range router.Routes() {
		route.RelativePath = CombinePath(router.Prefix, route.RelativePath)
		f.routes[route.RelativePath] = route
	}

	return f
}

// Deprecated: SetDescription 设置一个路由的详细描述信息, 应使用 RouteModel.SetDescription()
// @param  relativePath  string  相对路由
// @param  Description   string  详细描述信息
func (f *Router) SetDescription(relativePath, description string) *Router {
	if _, ok := f.routes[relativePath]; ok {
		old := f.routes[relativePath]
		old.Description = description
		f.routes[relativePath] = old
	}

	return f
}

func (f *Router) method(
	method, relativePath, summary string,
	queryModel any, requestModel, responseModel swag.BaseModelIface,
	handler HandlerFunc,
	additions []any,
) *Route {
	// 路由处理函数，默认仅一个
	handlers := []fiber.Handler{routeHandler(handler)}
	deprecated := false // 是否禁用此路由

	for _, adt := range additions {
		rt := reflect.TypeOf(adt)
		switch rt.Kind() {
		case reflect.String:
			if adt == "deprecated" {
				deprecated = true
			}
		case reflect.Func:
			// 发现fiber.handler
			handlers = append(handlers, routeHandler(adt.(HandlerFunc)))
		}
	}

	if !f.deprecated { // 若路由组被禁用，则此路由必禁用
		deprecated = true
	}

	// 确保路径以/开头，若路由为空，则以路由组前缀为路由路径
	if len(relativePath) > 0 && !strings.HasPrefix(relativePath, constant.PathSeparator) {
		relativePath = constant.PathSeparator + relativePath
	}

	route := &Route{
		Method:        method,
		RelativePath:  relativePath,
		PathFields:    make([]*swag.QModel, 0),               // 路径参数
		QueryModel:    swag.QModelTransformer(queryModel),    // 查询参数
		RequestModel:  swag.RModelTransformer(requestModel),  // 请求体
		ResponseModel: swag.RModelTransformer(responseModel), // 响应体
		Summary:       summary,
		Handlers:      handlers,
		Dependencies:  make([]HandlerFunc, 0),
		Tags:          f.Tags,
		Description:   method + " " + summary,
		deprecated:    deprecated,
	}

	// 生成路径参数
	if pp, found := swag.DoesPathParamsFound(route.RelativePath); found {
		for name, required := range pp {
			qm := &swag.QModel{Name: name, Required: required, InPath: true}
			if required {
				qm.Tag = reflect.StructTag(`json:"` + name + `" validate:"required" binding:"required"`)
			} else {
				qm.Tag = reflect.StructTag(`json:"` + name + `"`)

			}
			route.PathFields = append(route.PathFields, qm)
		}
	}

	f.routes[relativePath+RouteSeparator+method] = route // 允许地址相同,方法不同的路由

	return route
}

// GET http get method
// @param  path           string         相对路径,必须以"/"开头
// @param  summary        string         路由摘要信息
// @param  queryModel     struct         查询参数，仅支持struct类型
// @param  responseModel  any            响应体对象,  此model应为一个空struct实例,而非指针类型
// @param  handler        []HandlerFunc  路由处理方法
// @param  addition       any            附加参数，如："deprecated"用于禁用此路由
func (f *Router) GET(
	path string, responseModel swag.BaseModelIface, summary string, handler HandlerFunc, addition ...any,
) *Route {
	// 对于查询参数仅允许struct类型
	return f.method(
		http.MethodGet, path, summary,
		nil, nil, responseModel,
		handler, addition,
	)
}

// DELETE http delete method
// @param  path           string         相对路径,必须以"/"开头
// @param  summary        string         路由摘要信息
// @param  queryModel     struct         查询参数，仅支持struct类型
// @param  responseModel  any            响应体对象,  此model应为一个空struct实例,而非指针类型
// @param  handler        []HandlerFunc  路由处理方法
// @param  addition       any            附加参数
func (f *Router) DELETE(
	path string, responseModel swag.BaseModelIface, summary string, handler HandlerFunc, addition ...any,
) *Route {
	// 对于查询参数仅允许struct类型
	return f.method(
		http.MethodDelete, path, summary,
		nil, nil, responseModel,
		handler, addition,
	)
}

// POST http post method
// @param  path           string         相对路径,必须以"/"开头
// @param  summary        string         路由摘要信息
// @param  requestModel   any            请求体对象,  此model应为一个空struct实例,而非指针类型
// @param  responseModel  any            响应体对象,  此model应为一个空struct实例,而非指针类型
// @param  handler        []HandlerFunc  路由处理方法
// @param  addition       any            附加参数，如："deprecated"用于禁用此路由
func (f *Router) POST(
	path string,
	requestModel, responseModel swag.BaseModelIface,
	summary string,
	handler HandlerFunc,
	addition ...any,
) *Route {
	return f.method(
		http.MethodPost, path, summary,
		nil, requestModel, responseModel,
		handler, addition,
	)
}

// PATCH http patch method
func (f *Router) PATCH(
	path string,
	requestModel, responseModel swag.BaseModelIface,
	summary string,
	handler HandlerFunc,
	addition ...any,
) *Route {
	return f.method(
		http.MethodPatch, path, summary,
		nil, requestModel, responseModel,
		handler, addition,
	)
}

// PUT http put method
func (f *Router) PUT(
	path string,
	requestModel, responseModel swag.BaseModelIface,
	summary string,
	handler HandlerFunc,
	addition ...any,
) *Route {
	return f.method(
		http.MethodPut, path, summary,
		nil, requestModel, responseModel,
		handler, addition,
	)
}

// CombinePath 合并路由
// @param  prefix  string  路由前缀
// @param  path    string  路由
func CombinePath(prefix, path string) string {
	if path == "" {
		return prefix
	}
	if !strings.HasPrefix(prefix, constant.PathSeparator) {
		prefix = constant.PathSeparator + prefix
	}

	if strings.HasSuffix(prefix, constant.PathSeparator) && strings.HasPrefix(path, constant.PathSeparator) {
		return prefix[:len(prefix)-1] + path
	}
	return prefix + path
}

// GetRoute 查询自定义路由
// @param   method  string  请求方法
// @param   path    string  请求路由
// @return  *Route 自定义路由对象
func GetRoute(method string, path string) (route *Route) {
	switch method {
	case http.MethodGet:
		r, ok := MethodGetRoutes[path]
		if !ok {
			route = nil
		} else {
			route = r
		}
	case http.MethodPut:
		r, ok := MethodPutRoutes[path]
		if !ok {
			route = nil
		} else {
			route = r
		}
	case http.MethodPatch:
		r, ok := MethodPatchRoutes[path]
		if !ok {
			route = nil
		} else {
			route = r
		}
	case http.MethodDelete:
		r, ok := MethodDeleteRoutes[path]
		if !ok {
			route = nil
		} else {
			route = r
		}
	case http.MethodPost:
		r, ok := MethodPostRoutes[path]
		if !ok {
			route = nil
		} else {
			route = r
		}
	default:
		route = nil
	}

	return
}
