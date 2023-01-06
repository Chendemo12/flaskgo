package app

import (
	"bytes"
	"github.com/Chendemo12/flaskgo/internal/core"
	"github.com/Chendemo12/functools/helper"
	"github.com/gofiber/fiber/v2"
	fiberu "github.com/gofiber/fiber/v2/utils"
	"strings"
)

var recoverHandler StackTraceHandlerFunc = nil
var fiberErrorHandler fiber.ErrorHandler = nil // 设置fiber自定义错误处理函数

// HandlerFunc 路由处理函数
type HandlerFunc = func(s *Context) *Response

// StackTraceHandlerFunc 错误堆栈处理函数, 即 recover 方法
type StackTraceHandlerFunc = func(c *fiber.Ctx, e any)

// routeHandler 路由处理方法(装饰器实现)，用于请求体校验和返回体序列化，同时注入全局服务依赖,
// 此方法接收一个业务层面的路由钩子方法handler，
// 该方法有且仅有1个参数：flaskgo.Context(), 有且必须有一个返回值 *Response
//
// routeHandler 方法首先会new一个新的 flaskgo.Context, 并初始化请求体、路由参数、fiber.Ctx 和 flaskgo.Service
// 之后会校验并绑定路由参数（包含路径参数和查询参数）是否正确，如果错误则直接返回422错误，反之会继续序列化并绑定请求体（如果存在）
// 序列化成功之后会校验请求参数正确性（开关控制），校验通过后会接着将ctx传入handler
// 执行handler之后将校验返回值（开关控制），并返回422或写入响应体。
//
// @return  fiber.Handler fiber路由处理方法
func routeHandler(f HandlerFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := &Context{
			PathFields:  make(map[string]string),
			QueryFields: make(map[string]string),
			RequestBody: int64(1), // 初始化为1，避免访问错误
			fs:          appEngine.Service(),
			ec:          c,
		}

		// 路由唯一标识: c.Method()+RouteSeparator+c.RelativePath()
		// c.Route().RelativePath 获取注册的路径，
		// c.RelativePath() 获取匹配后的请求路由
		route := GetRoute(c.Method(), c.Route().Path) // 获取请求路由
		if route != nil {                             // 存在路由信息
			resp := routeParamsValidate(ctx, route) // 路由参数校验
			if resp != nil {
				// 路径参数或查询参数校验未通过
				return c.Status(resp.StatusCode).JSON(resp.Content)
			}

			//resp = requestBodyMarshal(ctx, route) // 请求体序列化
			//if resp != nil {
			//	return c.Status(resp.StatusCode).JSON(resp.Content)
			//}

			if !core.RequestValidateDisabled { // 开启了请求体自动校验
				resp = ctx.Validate(ctx.RequestBody)
				if resp != nil {
					return c.Status(resp.StatusCode).JSON(resp.Content)
				}
			}

			// ------------------------------- 校验通过或禁用自动校验 -------------------------------
			// 处理依赖项
			resp = dependencyDone(ctx, route)
			if resp != nil {
				return responseWriter(c, resp) // 返回消息流
			}
		}
		//
		// 执行处理函数并获取返回值
		if resp := f(ctx); resp != nil { // 自定义函数存在返回值
			return responseWriter(c, resp) // 返回消息流
		}

		// 自定义函数无任何返回值
		return c.Status(fiber.StatusOK).SendString(fiberu.StatusMessage(fiber.StatusOK))
	}
}

// jsoniterUnmarshalErrorToValidationError 将jsoniter 的反序列化错误转换成 接口错误类型
func jsoniterUnmarshalErrorToValidationError(err error) *ValidationError {
	// jsoniter 的反序列化错误格式：
	//
	// jsoniter.iter.ReportError():224
	//
	// 	iter.Error = fmt.Errorf("%s: %s, error found in #%v byte of ...|%s|..., bigger context ...|%s|...",
	//		operation, msg, iter.head-peekStart, parsing, context)
	//
	// 	err.Error():
	//
	// 	main.SimpleForm.Name: ReadString: expects " or n, but found 2, error found in #10 byte of ...| "name": 23,
	//		"a|..., bigger context ...|{
	//		"name": 23,
	//		"age": "23",
	//		"sex": "F"
	// 	}|...
	msg := err.Error()
	ve := &ValidationError{Loc: []string{"body"}}
	for i := 0; i < len(msg); i++ {
		if msg[i:i+1] == ":" {
			ve.Loc = append(ve.Loc, msg[:i])
			break
		}
	}
	if msgs := strings.Split(msg, jsoniterUnmarshalErrorSeparator); len(msgs) > 0 {
		_ = helper.DefaultJsonUnmarshal([]byte(msgs[jsonErrorFormIndex]), &ve.Ctx)
		ve.Msg = msgs[jsonErrorFieldMsgIndex][len(ve.Loc[1])+2:]
		if s := strings.Split(ve.Msg, ":"); len(s) > 0 {
			ve.Type = s[0]
		}
	}

	return ve
}

// routeParamsValidate 路径参数和查询参数校验
func routeParamsValidate(ctx *Context, route *Route) *Response {
	// 路径参数校验
	for i := 0; i < len(route.PathFields); i++ {
		ctx.PathFields[route.PathFields[i].Name] = ctx.Context().Params(route.PathFields[i].Name)
		if route.PathFields[i].IsRequired() && ctx.PathFields[route.PathFields[i].Name] == "" {
			// 不存在此路径参数, 但是此路径参数设置为必选
			return ValidationErrorResponse(&ValidationError{
				Loc:  []string{"path", route.PathFields[i].Name},
				Msg:  "path must not be empty",
				Type: "string",
				Ctx:  nil,
			})
		}
	}

	// 查询参数校验
	for i := 0; i < len(route.QueryFields); i++ {
		ctx.QueryFields[route.QueryFields[i].Name] = ctx.Context().Query(route.QueryFields[i].Name)
		if route.QueryFields[i].IsRequired() && ctx.QueryFields[route.QueryFields[i].Name] == "" {
			// 但是此查询参数设置为必选
			return ValidationErrorResponse(&ValidationError{
				Loc:  []string{"query", route.QueryFields[i].Name},
				Msg:  "query must not be empty",
				Type: "string",
				Ctx:  nil,
			})
		}
	}

	return nil
}

func dependencyDone(ctx *Context, route *Route) *Response {
	for i := 0; i < len(route.Dependencies); i++ {
		if resp := route.Dependencies[i](ctx); resp != nil {
			return resp
		}
	}

	return nil
}

func responseWriter(c *fiber.Ctx, resp *Response) error {
	switch resp.Type {

	case JsonResponseType: // Json类型
		if core.ResponseValidateDisabled {
			return c.Status(resp.StatusCode).JSON(resp.Content)
		} else {
			// TODO: implement 响应体校验
			return c.Status(resp.StatusCode).JSON(resp.Content)
		}

	case StringResponseType:
		return c.Status(resp.StatusCode).SendString(resp.Content.(string))

	case HtmlResponseType: // 返回HTML页面
		// 设置返回类型
		c.Set(fiber.HeaderContentType, resp.ContentType)
		return c.Status(resp.StatusCode).SendString(resp.Content.(string))

	case ErrResponseType:
		return c.Status(resp.StatusCode).JSON(resp.Content)

	case StreamResponseType: // 返回字节流
		return c.SendStream(bytes.NewReader(resp.Content.([]byte)))

	case FileResponseType: // 返回一个文件
		return c.Download(resp.Content.(string))

	case AdvancedResponseType:
		return resp.Content.(fiber.Handler)(c)

	case CustomResponseType:
		c.Status(resp.StatusCode).Set(fiber.HeaderContentType, resp.ContentType)
		switch resp.ContentType {

		case fiber.MIMETextHTML, fiber.MIMETextHTMLCharsetUTF8:
			return c.SendString(resp.Content.(string))
		case fiber.MIMEApplicationJSON, fiber.MIMEApplicationJSONCharsetUTF8:
			return c.JSON(resp.Content)
		case fiber.MIMETextXML, fiber.MIMEApplicationXML, fiber.MIMETextXMLCharsetUTF8, fiber.MIMEApplicationXMLCharsetUTF8:
			return c.XML(resp.Content)
		case fiber.MIMETextPlain, fiber.MIMETextPlainCharsetUTF8:
			return c.SendString(resp.Content.(string))
		//case fiber.MIMETextJavaScript, fiber.MIMETextJavaScriptCharsetUTF8:
		//case fiber.MIMEApplicationForm:
		//case fiber.MIMEOctetStream:
		//case fiber.MIMEMultipartForm:
		default:
			return c.Status(resp.StatusCode).JSON(resp.Content)
		}
	default:
		return c.Status(resp.StatusCode).JSON(resp.Content)
	}
}
