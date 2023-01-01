package swag

import (
	"bytes"
	"github.com/Chendemo12/flaskgo/internal/mode"
	"github.com/Chendemo12/functools/helper"
	"github.com/gofiber/fiber/v2"
)

const (
	ModelsSelectorName = "schemas"
	ModelsRefPrefix    = "#/components/schemas/"
)

var (
	template      = dict{}            // swag文档模板
	pathsMap      = map[string]dict{} // swag的路由文档
	modelsDocMap  = map[string]dict{} // swag的模型描述
	templateBytes = make([]byte, 0)   // swag文档模板字节表示
)

type dict map[string]any

// ------------------------------------------- 创建基础路由 -------------------------------------------

func clearCacheMap() {
	if !mode.IsDebug() {
		pathsMap = make(map[string]dict, 0)
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

	// TODO：implement
	// redoc 纯文档页面
	f.Get("/redoc", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(makeSwaggerUiHtml(title, openapiUrl, swaggerJsUrl, swaggerCssUrl, swaggerFaviconUrl))
	})

	// openapi 获取路由定义
	f.Get("/openapi.json", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusOK).SendStream(bytes.NewReader(templateBytes))
		//return c.Status(fiber.StatusOK).SendString(templateString)
		//return c.Status(fiber.StatusOK).JSON(template)
	})
}

// ------------------------------------------- 创建默认基础路由 end -------------------------------------------

func AddModelDoc(model *RouteModel) {
	if model != nil && model.Model != nil {
		modelsDocMap[model.FullName()] = model.Doc()
	}
}

func AddPathDoc(model *RouteInsGroup) {
	if model != nil {
		pathsMap[FastApiRoutePath(model.String())] = model.Swag()
	}
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
	// ---------------------------- swagger routes ------------------------
	modelsDocMap[String.String()] = String.Swag()

	// ---------------------------- swagger base info ------------------------
	template["openapi"] = "3.0.2"
	template["info"] = map[string]any{
		"description": description,
		"title":       title,
		"license":     license,
		"version":     version,
	}

	// ---------------------------- swagger descriptions ------------------------
	template["paths"] = pathsMap
	template["components"] = dict{ModelsSelectorName: modelsDocMap}

	// 序列化文档后返回字节流
	templateBytes, _ = helper.DefaultJsonMarshal(template)

	createSwaggerRoutes(f, title)
	clearCacheMap()
}
