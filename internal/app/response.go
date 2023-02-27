package app

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type ResponseType int

const (
	CustomResponseType ResponseType = iota + 1
	JsonResponseType
	StringResponseType
	StreamResponseType
	FileResponseType
	ErrResponseType
	HtmlResponseType
	AdvancedResponseType
)

var responseHeaders []*ResponseHeader

// ResponseHeader 自定义响应头
type ResponseHeader struct {
	Key   string `json:"key" Description:"Key" binding:"required"`
	Value string `json:"value" Description:"Value" binding:"required"`
}

// Response 路由返回值
type Response struct {
	Content     any          `json:"content"`     // 响应体
	ContentType string       `json:"contentType"` // 响应类型,默认为 application/json
	Type        ResponseType `json:"type"`        // 返回体类型
	StatusCode  int          `json:"status_code"` // 响应状态码

}

// ValidationError 参数校验错误
type ValidationError struct {
	Ctx  map[string]any `json:"service" Description:"Service"`
	Msg  string         `json:"msg" Description:"Message" binding:"required"`
	Type string         `json:"type" Description:"Error Type" binding:"required"`
	Loc  []string       `json:"loc" Description:"Location" binding:"required"`
}

func (v ValidationError) Doc__() string { return "Validation Error" }

type HTTPValidationError struct {
	Detail []*ValidationError `json:"detail" Description:"Detail" binding:"required"`
}

func (v HTTPValidationError) Doc__() string { return "HTTPValidationError" }

type CustomError404 struct {
	ErrorCode string `json:"error_code" Description:"Error Code" binding:"required"`
	ValidationError
}

func (v CustomError404) Doc__() string { return "CustomError404" }

// ValidationErrorResponse 参数校验错误返回值
func ValidationErrorResponse(ves ...*ValidationError) *Response {
	return &Response{
		StatusCode: http.StatusUnprocessableEntity,
		Content:    &HTTPValidationError{Detail: ves},
		Type:       ErrResponseType,
	}
}

// AnyResponse 自定义响应体,响应体可是任意类型
// @param   statusCode   int     响应状态码
// @param   content      any     响应体
// @param   contentType  string  响应头MIME
// @return  resp *Response response返回体
func AnyResponse(statusCode int, content any, contentType string) *Response {
	return &Response{
		StatusCode: statusCode, Content: &content, ContentType: contentType,
		Type: CustomResponseType,
	}
}

// JSONResponse 仅支持可以json序列化的响应体
// @param   statusCode  int  响应状态码
// @param   content     any  可以json序列化的类型
// @return  resp *Response response返回体
func JSONResponse(statusCode int, content any) *Response {
	return &Response{
		StatusCode: statusCode, Content: &content, Type: JsonResponseType,
	}
}

// StringResponse 返回值为字符串对象
// @param   content  string  字符串文本
// @return  resp *Response response返回体
func StringResponse(content string) *Response {
	return &Response{
		StatusCode: http.StatusOK, Content: content, Type: StringResponseType,
	}
}

// StreamResponse 返回值为字节流对象
// @param   statusCode  int     响应状态码
// @param   content     []byte  字节流
// @return  resp *Response response返回体
func StreamResponse(statusCode int, content []byte) *Response {
	return &Response{
		StatusCode: statusCode, Content: &content, Type: StreamResponseType,
	}
}

// FileResponse 返回值为文件对象，如：照片视频文件流等, 若文件不存在，则状态码置为404
// @param   filepath  string  文件路径
// @return  resp *Response response返回体
func FileResponse(filepath string) *Response {
	return &Response{
		StatusCode: http.StatusOK, Content: filepath, Type: FileResponseType,
	}
}

// ErrorResponse 返回一个服务器错误
// @param   content  any  错误消息
// @return  resp *Response response返回体
func ErrorResponse(content any) *Response {
	return JSONResponse(http.StatusInternalServerError, content)
}

// HTMLResponse 返回一段HTML文本
// @param   statusCode  int     响应状态码
// @param   content     string  HTML文本字符串
// @return  resp *Response response返回体
func HTMLResponse(statusCode int, context string) *Response {
	return &Response{
		Type:        HtmlResponseType,
		StatusCode:  statusCode,
		Content:     context,
		ContentType: fiber.MIMETextHTML,
	}
}

// OKResponse 返回状态码为200的 JSONResponse
// @param   content  any  可以json序列化的类型
// @return  resp *Response response返回体
func OKResponse(content any) *Response { return JSONResponse(http.StatusOK, content) }

// AdvancedResponse 高级返回值，允许返回一个函数，以实现任意类型的返回
// @param   statusCode  int            响应状态码
// @param   content     fiber.Handler  钩子函数
// @return  resp *Response response返回体
func AdvancedResponse(statusCode int, content fiber.Handler) *Response {
	return &Response{
		Type:        AdvancedResponseType,
		StatusCode:  statusCode,
		Content:     content,
		ContentType: "",
	}
}
