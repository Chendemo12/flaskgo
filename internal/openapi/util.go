package openapi

import (
	"github.com/Chendemo12/flaskgo/internal/godantic"
	"net/http"
	"strings"
)

// MakeOperationRequestBody 将路由中的 godantic.SchemaIface 转换成 openapi 的请求体 RequestBody
func MakeOperationRequestBody(model godantic.SchemaIface) *RequestBody {
	r := &RequestBody{
		Required: model.IsRequired(),
		Content: PathModelContent{
			MIMEType: MIMEApplicationJSON,
			Schema:   nil,
		},
	}

	bcs := BaseModelContentSchema{Title: model.SchemaName(), Type: model.SchemaType()}
	switch model.SchemaType() {
	case godantic.ArrayType:
		r.Content.Schema = ArrayModelContentSchema{
			BaseModelContentSchema: bcs,
			Items: Reference{
				Name: model.SchemaName(),
			},
		}

	case godantic.ObjectType:

		r.Content.Schema = ObjectModelContentSchema{
			BaseModelContentSchema: bcs,
			Reference: Reference{
				Name: model.SchemaName(),
			},
		}

	default:
		r.Content.Schema = bcs
	}

	return r
}

// MakeOperationResponses 将路由中的 godantic.SchemaIface 转换成 openapi 的返回体 []*Response
func MakeOperationResponses(model godantic.SchemaIface) []*Response {
	m := make([]*Response, 2) // 200 + 422

	// 200 接口出注册的返回值
	m[0] = &Response{
		StatusCode:  http.StatusOK,
		Description: http.StatusText(http.StatusOK),
		Content: PathModelContent{
			MIMEType: MIMEApplicationJSON,
			Schema:   nil,
		},
	}

	// 422 所有接口默认携带的请求体校验错误返回值
	m[1] = &Response{
		StatusCode:  http.StatusUnprocessableEntity,
		Description: http.StatusText(http.StatusUnprocessableEntity),
		Content: PathModelContent{
			MIMEType: MIMEApplicationJSON,
			Schema: ObjectModelContentSchema{
				BaseModelContentSchema: BaseModelContentSchema{
					Title: HttpValidationErrorName,
					Type:  godantic.ObjectType,
				},
				Reference: Reference{
					Name: HttpValidationErrorName,
				},
			},
		},
	}

	bcs := BaseModelContentSchema{Title: model.SchemaName(), Type: model.SchemaType()}
	switch model.SchemaType() {
	case godantic.ObjectType:
		m[0].Content.Schema = ObjectModelContentSchema{
			BaseModelContentSchema: bcs,
			Reference: Reference{
				Name: model.SchemaName(),
			},
		}
	case godantic.ArrayType:
		m[0].Content.Schema = ArrayModelContentSchema{
			BaseModelContentSchema: bcs,
			Items: Reference{
				Name: model.SchemaName(),
			},
		}
	default:
		m[0].Content.Schema = bcs
	}

	return m
}

// NewOpenApi 构造一个新的 OpenApi 文档
func NewOpenApi(title, version, description string) *OpenApi {
	return &OpenApi{
		Version: ApiVersion,
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
		Components:  &Components{Scheme: make([]*ComponentScheme, 0)},
		Paths:       &Paths{Paths: make([]*PathItem, 0)},
		initialized: false,
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
	paths := strings.Split(path, PathSeparator) // 路径字符
	// 查找路径中的参数
	for i := 0; i < len(paths); i++ {
		if strings.HasPrefix(paths[i], PathParamPrefix) {
			// 识别到路径参数
			if strings.HasSuffix(paths[i], OptionalPathParamSuffix) {
				// 可选路径参数
				paths[i] = "{" + paths[i][1:len(paths[i])-1] + "}"
			} else {
				paths[i] = "{" + paths[i][1:] + "}"
			}
		}
	}

	return strings.Join(paths, PathSeparator)
}
