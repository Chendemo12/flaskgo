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
	BoolType    OpenApiDataType = "boolean"
	ObjectType  OpenApiDataType = "object"
	ArrayType   OpenApiDataType = "array"
)
