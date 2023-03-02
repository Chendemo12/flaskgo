package godantic

const (
	RefName   = "$ref"
	RefPrefix = "#/components/schemas/"
)

type OpenApiDataType string

const (
	IntegerType OpenApiDataType = "integer"
	NumberType  OpenApiDataType = "number"
	StringType  OpenApiDataType = "string"
	ArrayType   OpenApiDataType = "array"
	ObjectType  OpenApiDataType = "object"
	BoolType    OpenApiDataType = "boolean"
)
