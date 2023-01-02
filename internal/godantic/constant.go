package godantic

var MethodsWithBody = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "PATCH"}

const RefPrefix = "#/components/schemas/"
const (
	PathParamPrefix         = ":" // 路径参数起始字符
	PathSeparator           = "/" // 路径分隔符
	OptionalPathParamSuffix = "?" // 可选路径参数结束字符
)

const (
	IntegerName = "integer"
	NumberName  = "number"
	StringName  = "string"
	ArrayName   = "array"
	ObjectName  = "object"
	BooleanName = "boolean"
)
