package openapi

import (
	"bytes"
	"github.com/Chendemo12/flaskgo/internal/godantic"
	"github.com/Chendemo12/flaskgo/internal/mode"
	"github.com/Chendemo12/functools/helper"
	"github.com/gofiber/fiber/v2"
)

type dict map[string]any

const (
	ModelsSelectorName = "schemas"
	ModelsRefPrefix    = "#/components/schemas/"
)

var (
	template      = dict{}              // swag文档模板
	pathsMap      = map[string][]dict{} // swag的路由文档
	modelsDocMap  = map[string]dict{}   // swag的模型描述
	templateBytes = make([]byte, 0)     // swag文档模板字节表示
)

// 422 表单验证错误模型
var validationErrorDefinition = dict{
	"title": "ValidationError",
	"type":  "object",
	"properties": dict{
		"loc": dict{
			"title": "Location",
			"type":  "array",
			"items": dict{"anyOf": []map[string]string{{"type": "string"}, {"type": "integer"}}},
		},
		"msg":  dict{"title": "Message", "type": "string"},
		"type": dict{"title": "Error Type", "type": "string"},
	},
	"required": []string{"loc", "msg", "type"},
}

// 请求体相应体错误消息
var validationErrorResponseDefinition = dict{
	"title":    "HTTPValidationError",
	"type":     "object",
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
	"title":    "CustomValidationError",
	"required": []string{"error_code"},
	"type":     "object",
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

// ------------------------------------------- 创建基础路由 -------------------------------------------

func clearCacheMap() {
	if !mode.IsDebug() {
		pathsMap = make(map[string][]dict, 0)
		template = make(dict, 0)
		modelsDocMap = make(map[string]dict, 0)
	}
}

func createSwaggerRoutes(f *fiber.App, title string) {
	// docs 在线调试页面
	f.Get("/docs", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(makeSwaggerUiHtml(title, openapiUrl, swaggerJsUrl, swaggerCssUrl, swaggerFaviconUrl))
	})

	// redoc 纯文档页面
	f.Get("/redoc", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(makeRedocUiHtml(title, openapiUrl, redocJsUrl, redocFaviconUrl))
	})

	// openapi 获取路由定义
	f.Get("/openapi.json", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		//return c.Status(fiber.StatusOK).SendString(templateString)
		//return c.Status(fiber.StatusOK).JSON(template)
		return c.Status(fiber.StatusOK).SendStream(bytes.NewReader(templateBytes))
	})
}

// ------------------------------------------- 创建默认基础路由 end -------------------------------------------

func AddModelDoc(name string, schema map[string]any) {
	modelsDocMap[name] = schema
}

func AddPathDoc(path string, schema map[string]any) {
	path = FastApiRoutePath(path)
	pathsMap[path] = append(pathsMap[path], schema)
}

// RegisterSwagger 挂载swagger文档
//
//	swagger文档自动生成核心在于依次构建出以下对象：
//
//		RouteModel -> RouteResp -> RouteInstance -> RouteInsGroup
//
//	最终通过调用 RouteInsGroup.Doc() 方法获取完整的一个路由信息，并调用 AddPathDoc 添加到:
//
//		pathsMap -> template 对象中.
//
//	RouteModel 为以下任意对象： RModel , RModelField 或 QModel
//
// 其中 RouteModel 必须首先被生成。并通过 AddModelDoc 添加模型文档
func RegisterSwagger(f *fiber.App, title, description, version string, license map[string]string) {
	modelsDocMap[String.String()] = String.Swag()
	modelsDocMap["ValidationError"] = validationErrorDefinition
	modelsDocMap["HTTPValidationError"] = validationErrorResponseDefinition
	modelsDocMap["CustomValidationError"] = customErrorDefinition

	// ---------------------------- swagger base info ------------------------
	template["openapi"] = "3.0.2"
	template["info"] = map[string]any{
		"description": description,
		"title":       title,
		"license":     license,
		"version":     version,
	}

	// ---------------------------- swagger routes ------------------------
	m := dict{}
	for path, methods := range pathsMap {
		routes := make(map[string]any, len(methods))
		for i := 0; i < len(methods); i++ {
			for k, v := range methods[i] {
				routes[k] = v
			}
		}

		m[path] = routes
	}

	template["paths"] = m
	// ---------------------------- swagger descriptions ------------------------
	template["components"] = dict{ModelsSelectorName: modelsDocMap}

	// 序列化文档后返回字节流
	templateBytes, _ = helper.DefaultJsonMarshal(template)

	createSwaggerRoutes(f, title)
	clearCacheMap()
}
