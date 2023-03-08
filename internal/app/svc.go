package app

import (
	"github.com/Chendemo12/flaskgo/internal/core"
	"github.com/Chendemo12/flaskgo/internal/godantic"
	"github.com/Chendemo12/flaskgo/internal/openapi"
	"github.com/Chendemo12/functools/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"reflect"
)

const ( // json序列化错误, 关键信息的序号
	jsoniterUnmarshalErrorSeparator = "|" // 序列化错误信息分割符, 定义于 validator/validator_instance.orSeparator
	jsonErrorFieldMsgIndex          = 0   // 错误原因
	jsonErrorFieldNameFormIndex     = 1   // 序列化错误的字段和值
	jsonErrorFormIndex              = 3   // 接收到的数据
)

var emptyMap = map[string]any{}
var emptyList = make([]string, 0)

type Dict = map[string]any

type Context struct {
	PathFields  map[string]string `json:"path_fields,omitempty"`  // 路径参数
	QueryFields map[string]string `json:"query_fields,omitempty"` // 查询参数
	RequestBody any               `json:"request_body,omitempty"` // 请求体，初始值为1
	app         *FlaskGo          `description:"flask-go application"`
	ec          *fiber.Ctx        `description:"engine context"`
	route       *Route            `description:"用于请求体和响应提校验"`
}

// Service 获取 FlaskGo 的 Service 服务依赖信息
//
//	@return	Service 服务依赖信息
func (c *Context) Service() *Service { return c.app.Service() }

// Context 获取web引擎的上下文 Service
//
//	@return	*fiber.Ctx fiber.App 的上下文信息
func (c *Context) Context() *fiber.Ctx { return c.ec }

// CustomContext 获取自定义服务上下文
func (c *Context) CustomContext() CustomService { return c.app.Service().CustomContext() }

func (c *Context) Config() any { return c.app.Service().Config() }

// Deprecated: App DO NOT DO THIS
func (c *Context) App() FlaskGo { return FlaskGo{} }

// Logger 获取日志句柄
func (c *Context) Logger() logger.Iface { return c.app.service.Logger() }

// Deprecated: Console use Logger instead
func (c *Context) Console() logger.Iface { return c.Logger() }

// Validator 获取请求体验证器
func (c *Context) Validator() *validator.Validate { return c.app.service.validate }

// Validate 结构体验证
func (c *Context) Validate(stc any) *Response { return c.app.service.Validate(stc) }

// BodyParser 序列化请求体
//
//	@param	c	*fiber.Ctx	fiber上下文
//	@param	a	any			请求体指针
//	@return	*Response 错误信息,若为nil 则序列化成功
func (c *Context) BodyParser(a any) *Response {
	if err := c.Context().BodyParser(a); err != nil { // 请求的表单序列化错误
		return ValidationErrorResponse(jsoniterUnmarshalErrorToValidationError(err))
	}

	return nil
}

// ShouldBindJSON 绑定并校验参数是否正确
func (c *Context) ShouldBindJSON(stc any) *Response {
	if err := c.BodyParser(stc); err != nil {
		return err
	}
	if resp := c.app.service.Validate(stc); resp != nil {
		return resp
	}
	return nil

}

func (c *Context) structResponseValidation(content any) *ValidationError {
	// 对于 struct 类型，允许缺省返回值以屏蔽返回值校验
	if core.ResponseValidateDisabled || c.route.ResponseModel == nil {
		return nil
	}

	rt := reflect.TypeOf(content)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	// 类型校验
	meta, err := c.route.ResponseModel.Metadata()
	if err != nil || meta.String() != rt.String() {
		v := &ValidationError{
			Ctx:  emptyMap,
			Msg:  ModelNotMatch,
			Type: string(godantic.StringType),
			Loc:  []string{rt.String()},
		}
		c.Logger().Error(v.SchemaDesc(), ":", v.String())
		return v
	}
	// 字段类型校验

	return nil
}

// OKResponse 返回状态码为200的 JSONResponse
//
//	@param	content	any	可以json序列化的类型
//	@return	resp *Response response返回体
func (c *Context) OKResponse(content any) *Response {
	if err := c.structResponseValidation(content); err != nil {
		return ValidationErrorResponse(err)
	}
	return OKResponse(content)
}

// JSONResponse 仅支持可以json序列化的响应体
//
//	@param	statusCode	int	响应状态码
//	@param	content		any	可以json序列化的类型
//	@return	resp *Response response返回体
func (c *Context) JSONResponse(statusCode int, content any) *Response {
	if err := c.structResponseValidation(content); err != nil {
		return ValidationErrorResponse(err)
	}
	return JSONResponse(statusCode, content)
}

// StringResponse 返回值为字符串对象
//
//	@param	content	string	字符串文本
//	@return	resp *Response response返回体
func (c *Context) StringResponse(content string) *Response {
	meta, err := c.route.ResponseModel.Metadata()
	if err != nil {
		return ValidationErrorResponse(&ValidationError{
			Ctx:  emptyMap,
			Msg:  ModelNotDefine,
			Type: string(godantic.StringType),
			Loc:  emptyList,
		})
	}

	if meta.SchemaType() == godantic.StringType {
		return StringResponse(content)
	} else {
		return ValidationErrorResponse(&ValidationError{
			Ctx:  emptyMap,
			Msg:  ModelNotString,
			Type: string(godantic.StringType),
			Loc:  emptyList,
		})
	}
}

// StreamResponse 返回值为字节流对象
//
//	@param	statusCode	int		响应状态码
//	@param	content		[]byte	字节流
//	@return	resp *Response response返回体
func (c *Context) StreamResponse(statusCode int, content []byte) *Response {
	return StreamResponse(statusCode, content)
}

// FileResponse 返回值为文件对象，如：照片视频文件流等, 若文件不存在，则状态码置为404
//
//	@param	filepath	string	文件路径
//	@return	resp *Response response返回体
func (c *Context) FileResponse(filepath string) *Response {
	return FileResponse(filepath)
}

// ErrorResponse 返回一个服务器错误
//
//	@param	content	any	错误消息
//	@return	resp *Response response返回体
func (c *Context) ErrorResponse(content any) *Response {
	return ErrorResponse(content)
}

// HTMLResponse 返回一段HTML文本
//
//	@param	statusCode	int		响应状态码
//	@param	content		string	HTML文本字符串
//	@return	resp *Response response返回体
func (c *Context) HTMLResponse(statusCode int, context string) *Response {
	return HTMLResponse(statusCode, context)
}

// AdvancedResponse 高级返回值，允许返回一个函数，以实现任意类型的返回
//
//	@param	statusCode	int				响应状态码
//	@param	content		fiber.Handler	钩子函数
//	@return	resp *Response response返回体
func (c *Context) AdvancedResponse(statusCode int, content fiber.Handler) *Response {
	return AdvancedResponse(statusCode, content)
}

// AnyResponse 自定义响应体,响应体可是任意类型
//
//	@param	statusCode	int		响应状态码
//	@param	content		any		响应体
//	@param	contentType	string	响应头MIME
//	@return	resp *Response response返回体
func (c *Context) AnyResponse(statusCode int, content any, contentType string) *Response {
	return AnyResponse(statusCode, content, contentType)
}

// ------------------------------------------------------------------------------------

// CustomService 自定义服务上下文信息
type CustomService interface {
	Config() any // 获取配置文件
}

// Service FlaskGo 全局服务依赖信息
// 此对象由FlaskGo启动时自动创建，此对象不应被修改，组合和嵌入，
// 但可通过SetServiceContext()接口设置自定义的上下文信息，并在每一个路由钩子函数中可得
type Service struct {
	logger   logger.Iface        `description:"日志对象"`
	ctx      CustomService       `description:"上层自定义服务依赖"`
	validate *validator.Validate `description:"请求体验证包"`
	openApi  *openapi.OpenApi    `description:"模型文档"`
	addr     string              `description:"绑定地址"`
}

// Config 获取自定义配置文件
func (s *Service) Config() any {
	if s.ctx == nil {
		return nil
	}
	return s.ctx.Config()
}

// CustomContext 获取自定义服务上下文
func (s *Service) CustomContext() CustomService { return s.ctx }

// Deprecated: ServiceContext 获取自定义服务上下文
func (s *Service) ServiceContext() CustomService { return s.CustomContext() }

// SetServiceContext 修改自定义服务上下文
func (s *Service) SetServiceContext(ctx CustomService) *Service {
	s.ctx = ctx
	return s
}

// Addr 绑定地址
//
//	@return	string 绑定地址
func (s *Service) Addr() string { return s.addr }

// Logger 获取日志句柄
func (s *Service) Logger() logger.Iface { return s.logger }

// ReplaceLogger 替换日志句柄
//
//	@param	logger	logger.Iface	日志句柄
func (s *Service) ReplaceLogger(logger logger.Iface) { s.logger = logger }

// Validator 获取请求体验证器
func (s *Service) Validator() *validator.Validate { return s.validate }

// Validate 结构体验证
func (s *Service) Validate(stc any) *Response {
	err := s.validate.Struct(stc)
	if err != nil { // 模型验证错误
		err, _ := err.(validator.ValidationErrors) // validator的校验错误信息

		if nums := len(err); nums == 0 {
			return ValidationErrorResponse()
		} else {
			ves := make([]*ValidationError, nums) // 自定义的错误信息
			for i := 0; i < nums; i++ {
				ves[i] = &ValidationError{
					Loc:  []string{"body", err[i].Field()},
					Msg:  err[i].Error(),
					Type: err[i].Type().String(),
					Ctx:  emptyMap,
				}
			}
			return ValidationErrorResponse(ves...)
		}
	}

	return nil
}
