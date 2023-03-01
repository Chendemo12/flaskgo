package openapi

import (
	"github.com/Chendemo12/flaskgo/internal/godantic"
	"net/http"
	"strings"
)

// MakeOperationRequestBody 将路由中的 godantic.SchemaIface 转换成 openapi 的请求体 RequestBody
func MakeOperationRequestBody(model godantic.SchemaIface) *RequestBody {
	if model == nil {
		return &RequestBody{}
	}

	r := &RequestBody{
		Required: model.IsRequired(),
		Content: &PathModelContent{
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
	if model == nil { // 若返回值为空，则设置为字符串
		model = godantic.String
	}

	m := make([]*Response, 2) // 200 + 422
	// 200 接口处注册的返回值
	m[0] = &Response{
		StatusCode:  http.StatusOK,
		Description: http.StatusText(http.StatusOK),
		Content: &PathModelContent{
			MIMEType: MIMEApplicationJSON,
			Schema:   nil,
		},
	}
	// 422 所有接口默认携带的请求体校验错误返回值
	m[1] = Resp422

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

func QModelToParameter(model *godantic.QModel) *Parameter {
	p := &Parameter{
		ParameterBase: ParameterBase{
			Description: model.SchemaDesc(),
			Required:    model.IsRequired(),
			Deprecated:  false,
		},
		Title:   model.Title,
		Name:    model.SchemaName(),
		In:      InQuery,
		Default: godantic.GetDefaultV(model.Tag, model.SchemaType()),
		Type:    model.SchemaType(),
	}

	if model.InPath {
		p.In = InPath
	}

	return p
}
