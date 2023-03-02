package godantic

import (
	"reflect"
	"strings"
	"unicode"
)

var metaDataFactory = &MetaClass{data: make([]*MetaData, 0)}

func GetMetaDataFactory() *MetaClass   { return metaDataFactory }
func GetMetaData(pkg string) *MetaData { return metaDataFactory.Get(pkg) }
func SaveMetaData(data *MetaData)      { metaDataFactory.Set(data) }

type MetaField struct {
	Field
	Index     int          `json:"index" description:"当前字段所处的序号"`
	Exported  bool         `json:"exported" description:"是否是导出字段"`
	Anonymous bool         `json:"anonymous" description:"是否是嵌入字段"`
	RType     reflect.Type `description:"反射字段类型"`
}

type MetaData struct {
	names       []string     `description:"结构体名称,包名.结构体名称"`
	fields      []*MetaField `description:"结构体字段"`
	innerFields []*MetaField `description:"内部字段"`
}

// Name 获取结构体名称
func (m MetaData) Name() string { return m.names[0] }

// String 结构体全称：包名+结构体名称
func (m MetaData) String() string { return m.names[1] }

// Fields 结构体字段
func (m MetaData) Fields() []*MetaField { return m.fields }

// InnerFields 内部字段
func (m MetaData) InnerFields() []*MetaField { return m.innerFields }

// Id 获取结构体的唯一标识
func (m MetaData) Id() string { return m.String() }

type MetaClass struct {
	data []*MetaData
}

func (m *MetaClass) Query(pkg string) *MetaData {
	for i := 0; i < len(m.data); i++ {
		if m.data[i].String() == pkg {
			return m.data[i]
		}
	}
	return nil
}

// Save 保存一个元信息，存在则更新
func (m *MetaClass) Save(meta *MetaData) {
	for i := 0; i < len(m.data); i++ {
		if m.data[i].String() == meta.String() {
			m.data[i] = meta
			return
		}
	}
	m.data = append(m.data, meta)
}

func (m *MetaClass) Get(pkg string) *MetaData { return m.Query(pkg) }

func (m *MetaClass) Set(meta *MetaData) { m.Save(meta) }

func structFieldToMetaField(field reflect.StructField) *MetaField {
	mf := &MetaField{
		Field: Field{
			_pkg:        field.PkgPath,
			Title:       field.Name,
			Tag:         field.Tag,
			Description: QueryFieldTag(field.Tag, "description", field.Name),
			Default:     QueryFieldTag(field.Tag, "default", ""),
			ItemRef:     "",
			OType:       reflectKindToOType(field.Type.Kind()),
		},
		Index:     field.Index[0],
		Exported:  unicode.IsUpper(rune(field.Name[0])),
		Anonymous: field.Anonymous,
		RType:     field.Type,
	}
	return mf
}

// 处理数组元素
// @param elemType reflect.Type 子元素类型
// @param metadata *MetaData 根模型元信息
// @param metaField *MetaField 字段元信息
func parseArrayField(elemType reflect.Type, metadata *MetaData, fieldMeta *MetaField, no int) {
	if elemType.Kind() == reflect.Pointer { // 数组元素为指针结构体
		elemType = elemType.Elem()
	}

	// 处理数组的子元素
	switch elemType.Kind() {

	case reflect.Array, reflect.Slice, reflect.Chan: // [][]*Student
		fieldMeta.OType = ArrayType
		fieldMeta.ItemRef = elemType.String()

		mf := &MetaField{
			Field: Field{
				_pkg:        elemType.String(),
				Title:       elemType.String(),
				Tag:         "",
				Description: "",
				Default:     "",
				ItemRef:     "",
				OType:       ArrayType,
			},
			Index:     0,
			Exported:  true,
			Anonymous: false,
			RType:     elemType.Elem(),
		}
		no += 1
		parseArrayField(elemType.Elem(), metadata, mf, no)

	case reflect.Struct:
		fieldMeta.ItemRef = elemType.String()
		no += 1
		for i := 0; i < elemType.NumField(); i++ { // 此时必不是指针
			field := elemType.Field(i)
			extractField(field, metadata, no) // 递归
		}

	default:
		fieldMeta.ItemRef = elemType.String()
	}
}

// 提取结构体字段信息并添加到元信息中
func extractField(field reflect.StructField, meta *MetaData, no int) {
	if field.Anonymous && (field.Name == "BaseModel" || field.Name == "Field") { // 过滤模型基类
		return
	}
	if strings.HasPrefix(field.Name, "_") { // 过滤约定的匿名字段
		return
	}

	// ---------------------------------- 获取字段信息 ----------------------------------

	fieldMeta := structFieldToMetaField(field)
	if no < 1 { // 根模型字段
		meta.fields = append(meta.fields, fieldMeta)
	} else {
		meta.innerFields = append(meta.innerFields, fieldMeta)
	}

	switch fieldMeta.OType {
	case IntegerType, NumberType, BoolType, StringType:
		return // 基本类型,无需继续递归处理

	case ArrayType: // 字段为数组
		no += 1
		parseArrayField(field.Type.Elem(), meta, fieldMeta, no)

	case ObjectType: // 字段为结构体，指针，接口，map等
		if field.Type.Kind() == reflect.Interface || field.Type.Kind() == reflect.Map {
			return // 接口或map无需继续向下递归
		}

		// 结构体或结构体指针
		var nextFieldType reflect.Type
		if field.Type.Kind() == reflect.Pointer {
			nextFieldType = field.Type.Elem()
		} else {
			nextFieldType = field.Type
		}

		fieldMeta.ItemRef = nextFieldType.String() // 关联模型
		no += 1
		for i := 0; i < nextFieldType.NumField(); i++ { // 此时必不是指针
			field := nextFieldType.Field(i)
			extractField(field, meta, no) // 递归
		}
	}
}

// Reflect 反射建立任意类型的元信息
func (m *MetaClass) Reflect(model SchemaIface) *MetaData {
	if nm, ok := model.(*Field); ok { // 接口处定义了基本数据类型, 或List
		if GetMetaData(nm._pkg) != nil {
			return nm.MetaData()
		}
	}

	rt := reflect.TypeOf(model)
	if rt.Kind() == reflect.Pointer { // 由于接口定义，此处全部为结构体指针
		rt = rt.Elem()
	}

	meta := &MetaData{
		names:       []string{rt.Name(), rt.String()}, // 获取包名
		fields:      make([]*MetaField, 0),
		innerFields: make([]*MetaField, 0),
	}

	for i := 0; i < rt.NumField(); i++ { // 此时必不是指针
		field := rt.Field(i)
		extractField(field, meta, 0) // 0 根起点
	}

	return meta
}
