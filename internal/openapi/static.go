package openapi

import (
	"github.com/Chendemo12/flaskgo/internal/godantic"
)

const ApiVersion = "3.0.2"

// 用于swagger的一些静态文件，来自FastApi
const (
	SwaggerCssUrl     = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@4/swagger-ui.css"
	SwaggerFaviconUrl = "https://fastapi.tiangolo.com/img/favicon.png"
	SwaggerJsUrl      = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@4/swagger-ui-bundle.js"
	RedocJsUrl        = "https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"
	RedocFaviconUrl   = "https://fastapi.tiangolo.com/img/favicon.png"
	JsonUrl           = "openapi.json"
)

const ModelSelectorName = "schemas"
const (
	PathParamPrefix         = ":" // 路径参数起始字符
	PathSeparator           = "/" // 路径分隔符
	OptionalPathParamSuffix = "?" // 可选路径参数结束字符
)

type ApplicationMIMEType string

const (
	MIMETextXML                    ApplicationMIMEType = "text/xml"
	MIMETextHTML                   ApplicationMIMEType = "text/html"
	MIMETextPlain                  ApplicationMIMEType = "text/plain"
	MIMETextJavaScript             ApplicationMIMEType = "text/javascript"
	MIMEApplicationXML             ApplicationMIMEType = "application/xml"
	MIMEApplicationJSON            ApplicationMIMEType = "application/json"
	MIMEApplicationForm            ApplicationMIMEType = "application/x-www-form-urlencoded"
	MIMEOctetStream                ApplicationMIMEType = "application/octet-stream"
	MIMEMultipartForm              ApplicationMIMEType = "multipart/form-data"
	MIMETextXMLCharsetUTF8         ApplicationMIMEType = "text/xml; charset=utf-8"
	MIMETextHTMLCharsetUTF8        ApplicationMIMEType = "text/html; charset=utf-8"
	MIMETextPlainCharsetUTF8       ApplicationMIMEType = "text/plain; charset=utf-8"
	MIMETextJavaScriptCharsetUTF8  ApplicationMIMEType = "text/javascript; charset=utf-8"
	MIMEApplicationXMLCharsetUTF8  ApplicationMIMEType = "application/xml; charset=utf-8"
	MIMEApplicationJSONCharsetUTF8 ApplicationMIMEType = "application/json; charset=utf-8"
)

type dict map[string]any

const (
	ValidationErrorName       string = "ValidationError"
	HttpValidationErrorName   string = "HTTPValidationError"
	CustomValidationErrorName string = "CustomValidationError"
)

// 422 表单验证错误模型
var validationErrorDefinition = dict{
	"title": ValidationErrorName,
	"type":  godantic.ObjectType,
	"properties": dict{
		"loc": dict{
			"title": "Location",
			"type":  "array",
			"items": dict{"anyOf": []map[string]string{{"type": "string"}, {"type": "integer"}}},
		},
		"msg":  dict{"title": "Message", "type": "string"},
		"type": dict{"title": "Error RType", "type": "string"},
	},
	"required": []string{"loc", "msg", "type"},
}

// 请求体相应体错误消息
var validationErrorResponseDefinition = dict{
	"title":    HttpValidationErrorName,
	"type":     godantic.ObjectType,
	"required": []string{"detail"},
	"properties": dict{
		"detail": dict{
			"title": "Detail",
			"type":  "array",
			"items": dict{"$ref": godantic.RefPrefix + "ValidationError"},
		},
	},
}

// 自定义错误消息
var customErrorDefinition = dict{
	"title":    CustomValidationErrorName,
	"required": []string{"error_code"},
	"type":     godantic.ObjectType,
	"properties": dict{
		"error_code": dict{
			"title":       "ErrorCode",
			"type":        "string",
			"required":    true,
			"description": "ErrorCode",
		},
		"ValidationError": dict{
			"$ref":        "#/components/schemas/ValidationError",
			"title":       "ValidationError",
			"type":        "object",
			"required":    false,
			"description": "ValidationError",
		},
		"description": "CustomValidationError",
	},
}
