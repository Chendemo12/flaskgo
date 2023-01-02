package godantic

type Define struct {
	BaseModel
	Name string
}

var String = &Define{Name: "String"}
