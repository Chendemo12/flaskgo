package openapi

import (
	"github.com/Chendemo12/flaskgo/internal/constant"
	"github.com/Chendemo12/flaskgo/internal/godantic"
	"net/http"
	"strings"
)

// MakeDataModelContent 构建请求体文档内容
func MakeDataModelContent(mimeType ApplicationMIMEType, schema RequestBodyContentSchema) map[ApplicationMIMEType]any {
	m := make(map[ApplicationMIMEType]any)
	m[mimeType] = map[string]any{
		"schema": schema.Schema(),
	}
	return m
}

func MakeOperationResponses(schema RequestBodyContentSchema) map[int]*Response {
	m := make(map[int]*Response, 0)

	// 200 接口出注册的返回值
	m[http.StatusOK] = &Response{
		PathDataModel: PathDataModel{
			Content: MakeDataModelContent(MIMEApplicationJSON, schema),
		},
		Description: http.StatusText(http.StatusOK),
	}

	// 422 所有接口默认携带的请求体校验错误返回值
	m[http.StatusUnprocessableEntity] = &Response{
		PathDataModel: PathDataModel{
			Content: MakeDataModelContent(MIMEApplicationJSON, ObjectRequestBodyContentSchema{
				BaseRequestBodyContentSchema: BaseRequestBodyContentSchema{
					Title: HttpValidationErrorName,
					Type:  godantic.ObjectType,
				},
				Reference: Reference{
					Ref: ModelRefPrefix + HttpValidationErrorName,
				},
			}),
		},
		Description: http.StatusText(http.StatusUnprocessableEntity),
	}

	return m
}

func NewOpenApi(title, version, description string) *OpenApi {
	return &OpenApi{
		Openapi: ApiVersion,
		Info: &Info{
			Title:          title,
			Version:        version,
			Description:    description,
			TermsOfService: "",
			Contact: Contact{
				Name:  "FlaskGo",
				Url:   "github.com/Chendemo12/flaskgo",
				Email: "chendemo12@gmail.com",
			},
			License: License{
				Name: "FlaskGo",
				Url:  "github.com/Chendemo12/flaskgo",
			},
		},
		Tags:        make([]Tag, 0),
		Servers:     map[string]string{},
		Definitions: []godantic.SchemaIface{},
		Routes:      make([]*PathItem, 0),
		cache:       make([]byte, 0),
	}
}

// FastApiRoutePath 将 fiber.App 格式的路径转换成 FastApi 格式的路径
//
//	Example:
//	必选路径参数：
//		Input: "/api/rcst/:no"
//		Output: "/api/rcst/{no}"
//	可选路径参数：
//		Input: "/api/rcst/:no?"
//		Output: "/api/rcst/{no}"
//	常规路径：
//		Input: "/api/rcst/no"
//		Output: "/api/rcst/no"
func FastApiRoutePath(path string) string {
	paths := strings.Split(path, constant.PathSeparator) // 路径字符
	// 查找路径中的参数
	for i := 0; i < len(paths); i++ {
		if strings.HasPrefix(paths[i], constant.PathParamPrefix) {
			// 识别到路径参数
			if strings.HasSuffix(paths[i], constant.OptionalPathParamSuffix) {
				// 可选路径参数
				paths[i] = "{" + paths[i][1:len(paths[i])-1] + "}"
			} else {
				paths[i] = "{" + paths[i][1:] + "}"
			}
		}
	}

	return strings.Join(paths, constant.PathSeparator)
}
