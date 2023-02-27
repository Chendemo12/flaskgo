package openapi

const ApiVersion = "3.0.2"

// 用于swagger的一些静态文件，来自FastApi
const (
	SwaggerCssUrl     = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@4/swagger-ui.css"
	SwaggerFaviconUrl = "https://fastapi.tiangolo.com/img/favicon.png"
	SwaggerJsUrl      = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@4/swagger-ui-bundle.js"
	RedocJsUrl        = "https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"
	RedocFaviconUrl   = "https://fastapi.tiangolo.com/img/favicon.png"
	OpenapiUrl        = "openapi.json"
)

const (
	ModelsSelectorName = "schemas"
	ModelsRefPrefix    = "#/components/schemas/"
)

const (
	PathParamPrefix         = ":" // 路径参数起始字符
	PathSeparator           = "/" // 路径分隔符
	OptionalPathParamSuffix = "?" // 可选路径参数结束字符
)

type dict map[string]any
