package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gitlab.cowave.com/gogo/functools/zaplog"
)

const ( // json序列化错误, 关键信息的序号
	jsoniterUnmarshalErrorSeparator = "|" // 序列化错误信息分割符, 定义于 validator/validator_instance.orSeparator
	jsonErrorFieldMsgIndex          = 0   // 错误原因
	jsonErrorFieldNameFormIndex     = 1   // 序列化错误的字段和值
	jsonErrorFormIndex              = 3   // 接收到的数据
)

type Context struct {
	PathFields  map[string]string `json:"path_fields,omitempty"`  // 路径参数
	QueryFields map[string]string `json:"query_fields,omitempty"` // 查询参数
	RequestBody any               `json:"request_body,omitempty"` // 请求体，初始值为1
	fs          *Service
	ec          *fiber.Ctx
}

// Version 获取版本号
// @return  string 版本号
func (c *Context) Version() string { return Version }

// Service 获取 FlaskGo 的 Service 服务依赖信息
// @return  Service 服务依赖信息
func (c *Context) Service() *Service { return c.fs }

// Context 获取web引擎的上下文 Service
// @return  *fiber.Ctx fiber.App 的上下文信息
func (c *Context) Context() *fiber.Ctx { return c.ec }

// CustomContext 获取自定义服务上下文
func (c *Context) CustomContext() CustomContextIface { return c.fs.ServiceContext() }

func (c *Context) Config() any { return c.fs.Config() }

func (c *Context) App() FlaskGo { return *c.fs.App() }

func (c *Context) Console() zaplog.ConsoleLogger { return c.fs.Console() }

func (c *Context) Logger() zaplog.Iface { return c.fs.Logger() }

// Validator 获取请求体验证器
func (c *Context) Validator() *validator.Validate { return c.fs.validate }

// Validate 结构体验证
func (c *Context) Validate(stc any) *Response { return c.fs.Validate(stc) }

// BodyParser 序列化请求体
// @param   c  *fiber.Ctx  fiber上下文
// @param   a  any         请求体指针
// @return  *Response 错误信息,若为nil 则序列化成功
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
	if resp := c.fs.Validate(stc); resp != nil {
		return resp
	}
	return nil

}

func (c *Context) OKResponse(content any) *Response { return OKResponse(content) }

func (c *Context) JSONResponse(statusCode int, content any) *Response {
	return JSONResponse(statusCode, content)
}

func (c *Context) StringResponse(content string) *Response {
	return StringResponse(content)
}

func (c *Context) StreamResponse(statusCode int, content []byte) *Response {
	return StreamResponse(statusCode, content)
}

func (c *Context) FileResponse(filepath string) *Response {
	return FileResponse(filepath)
}

func (c *Context) ErrorResponse(content any) *Response {
	return ErrorResponse(content)
}

func (c *Context) HTMLResponse(statusCode int, context string) *Response {
	return HTMLResponse(statusCode, context)
}

func (c *Context) AdvancedResponse(statusCode int, content fiber.Handler) *Response {
	return AdvancedResponse(statusCode, content)
}

func (c *Context) AnyResponse(statusCode int, content any, contentType string) *Response {
	return AnyResponse(statusCode, content, contentType)
}

// ------------------------------------------------------------------------------------

// CustomContextIface 自定义服务上下文信息
type CustomContextIface interface {
	Config() any // 获取配置文件
}

// Service FlaskGo 全局服务依赖信息
// 此对象由FlaskGo启动时自动创建，此对象不应被修改，组合和嵌入，
// 但可通过SetServiceContext()接口设置自定义的上下文信息，并在每一个路由钩子函数中可得
type Service struct {
	app      *FlaskGo            // FlaskGo 对象
	ctx      CustomContextIface  // 上层自定义服务依赖
	validate *validator.Validate // 请求体验证包
}

// App 获取FlaskGo Application
func (s *Service) App() *FlaskGo { return s.app }

// Config 获取自定义配置文件
func (s *Service) Config() any {
	if s.ctx == nil {
		return nil
	}
	return s.ctx.Config()
}

// ServiceContext 获取自定义服务上下文
func (s *Service) ServiceContext() CustomContextIface { return s.ctx }

// CustomContext 获取自定义服务上下文
func (s *Service) CustomContext() CustomContextIface { return s.ServiceContext() }

// SetServiceContext 修改自定义服务上下文
func (s *Service) SetServiceContext(ctx CustomContextIface) *Service {
	s.ctx = ctx
	return s
}

// Logger 获取日志
func (s *Service) Logger() zaplog.Iface { return s.app.logger }

// Console 获取控制台日志
func (s *Service) Console() zaplog.ConsoleLogger { return s.app.console }

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
					Ctx:  Dict{},
				}
			}
			return ValidationErrorResponse(ves...)
		}
	}

	return nil
}
