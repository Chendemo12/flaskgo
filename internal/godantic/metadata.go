package godantic

import (
	"reflect"
	"strings"
	"unicode"
)

var metaDataFactory = &MetaClass{data: make([]*Metadata, 0)}

func GetMetadataFactory() *MetaClass   { return metaDataFactory }
func GetMetadata(pkg string) *Metadata { return metaDataFactory.Get(pkg) }
func SaveMetadata(data *Metadata)      { metaDataFactory.Set(data) }

type MetaField struct {
	Field
	Index     int          `json:"index" description:"当前字段所处的序号"`
	Exported  bool         `json:"exported" description:"是否是导出字段"`
	Anonymous bool         `json:"anonymous" description:"是否是嵌入字段"`
	RType     reflect.Type `description:"反射字段类型"`
}

type Metadata struct {
	names       []string     `description:"结构体名称,包名.结构体名称"`
	fields      []*MetaField `description:"结构体字段"`
	innerFields []*MetaField `description:"内部字段"`
}

// Name 获取结构体名称
func (m *Metadata) Name() string { return m.names[0] }

// String 结构体全称：包名+结构体名称
func (m *Metadata) String() string { return m.names[1] }

// Fields 结构体字段
func (m *Metadata) Fields() []*MetaField { return m.fields }

// InnerFields 内部字段
func (m *Metadata) InnerFields() []*MetaField { return m.innerFields }

// Id 获取结构体的唯一标识
func (m *Metadata) Id() string { return m.String() }

// AddField 添加字段记录
func (m *Metadata) AddField(field *MetaField)      { m.fields = append(m.fields, field) }
func (m *Metadata) AddInnerField(field *MetaField) { m.innerFields = append(m.innerFields, field) }

type MetaClass struct {
	data []*Metadata
}

func (m *MetaClass) Query(pkg string) *Metadata {
	for i := 0; i < len(m.data); i++ {
		if m.data[i].String() == pkg {
			return m.data[i]
		}
	}
	return nil
}

// Save 保存一个元信息，存在则更新
func (m *MetaClass) Save(meta *Metadata) {
	for i := 0; i < len(m.data); i++ {
		if m.data[i].String() == meta.String() {
			m.data[i] = meta
			return
		}
	}
	m.data = append(m.data, meta)
}

func (m *MetaClass) Get(pkg string) *Metadata { return m.Query(pkg) }

func (m *MetaClass) Set(meta *Metadata) { m.Save(meta) }

// Reflect 反射建立任意类型的元信息
func (m *MetaClass) Reflect(model SchemaIface) *Metadata {
	if nm, ok := model.(*Field); ok { // 接口处定义了基本数据类型, 或List
		if GetMetadata(nm._pkg) != nil {
			return nm.Metadata()
		}
	}

	rt := reflect.TypeOf(model)
	if rt.Kind() == reflect.Pointer { // 由于接口定义，此处全部为结构体指针
		rt = rt.Elem()
	}

	meta := &Metadata{ // 构造根模型元信息
		names:       []string{rt.Name(), rt.String()}, // 获取包名
		fields:      make([]*MetaField, 0),
		innerFields: make([]*MetaField, 0),
	}

	ref := ModelReflect{metadata: meta}
	for i := 0; i < rt.NumField(); i++ { // 此时肯定是个结构体
		field := rt.Field(i)
		ref.extractField(field, 0) // 0 根起点
	}

	return meta
}

type ModelReflect struct {
	metadata *Metadata
}

func (m *ModelReflect) Metadata() *Metadata { return m.metadata }

// 提取结构体字段信息并添加到元信息中
func (m *ModelReflect) extractField(structField reflect.StructField, no int) {
	// 过滤模型基类
	if structField.Anonymous && (structField.Name == "BaseModel" || structField.Name == "Field") {
		return
	}
	// 过滤约定的匿名字段
	if strings.HasPrefix(structField.Name, "_") {
		return
	}

	// ---------------------------------- 获取字段信息 ----------------------------------

	fieldMeta := m.structFieldToMetaField(structField)
	if no < 1 { // 根模型字段
		m.metadata.AddField(fieldMeta)
	} else {
		m.metadata.AddInnerField(fieldMeta)
	}

	switch fieldMeta.OType {
	case IntegerType, NumberType, BoolType, StringType:
		return // 基本类型,无需继续递归处理

	case ObjectType: // 字段为结构体，指针，接口，map等
		if structField.Type.Kind() == reflect.Interface || structField.Type.Kind() == reflect.Map {
			return // 接口或map无需继续向下递归
		}

		no += 1
		// 结构体或结构体指针
		fieldElemType := structField.Type
		if fieldElemType.Kind() == reflect.Ptr {
			fieldElemType = fieldElemType.Elem()
		}

		m.parseStructField(fieldElemType, fieldMeta, no)

	case ArrayType: // 字段为数组
		no += 1
		fieldElemType := structField.Type.Elem() // 子元素类型
		if fieldElemType.Kind() == reflect.Ptr {
			fieldElemType = fieldElemType.Elem()
		}

		m.parseArrayField(fieldElemType, fieldMeta, no)
	}
}

func (m *ModelReflect) structFieldToMetaField(field reflect.StructField) *MetaField {
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

	if field.PkgPath == "" { // 对于结构体字段，此值无意义
		mf._pkg = m.metadata.String() + "." + field.Name
	}

	return mf
}

// 处理字段是数组的元素
// @param elemType reflect.Type 子元素类型
// @param metadata *Metadata 根模型元信息
// @param metaField *MetaField 字段元信息
func (m *ModelReflect) parseArrayField(elemType reflect.Type, fieldMeta *MetaField, no int) {
	if elemType.Kind() == reflect.Pointer { // 数组元素为指针结构体
		elemType = elemType.Elem()
	}

	// 处理数组的子元素
	kind := elemType.Kind()
	switch kind {
	case reflect.String:
		fieldMeta.ItemRef = String.SchemaName()
	case reflect.Bool:
		fieldMeta.ItemRef = Bool.SchemaName()

	case reflect.Array, reflect.Slice, reflect.Chan: // [][]*Student
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
			RType:     elemType,
		}
		fieldMeta.ItemRef = mf.SchemaName()
		m.metadata.AddInnerField(mf)
		no += 1
		m.parseArrayField(elemType.Elem(), mf, no)

	case reflect.Struct:
		fieldMeta.ItemRef = elemType.String()
		no += 1
		for i := 0; i < elemType.NumField(); i++ { // 此时必不是指针
			field := elemType.Field(i)
			m.extractField(field, no) // 递归
		}

	default:
		if reflect.Bool < kind && kind <= reflect.Uint64 {
			fieldMeta.ItemRef = Int.SchemaName()
		}
		if reflect.Float32 <= kind && kind <= reflect.Complex128 {
			fieldMeta.ItemRef = Float.SchemaName()
		}
	}
}

// 处理字段是结构体的元素
// @param elemType reflect.Type 子元素类型
// @param metadata *Metadata 根模型元信息
// @param metaField *MetaField 字段元信息
func (m *ModelReflect) parseStructField(elemType reflect.Type, fieldMeta *MetaField, no int) {

	mf := &MetaField{
		Field: Field{
			_pkg:        elemType.String(),
			Title:       elemType.String(),
			Tag:         "",
			Description: "",
			Default:     "",
			ItemRef:     "",
			OType:       ObjectType,
		},
		Index:     0,
		Exported:  true,
		Anonymous: false,
		RType:     elemType,
	}
	fieldMeta.ItemRef = mf.SchemaName() // 关联模型
	m.metadata.AddInnerField(mf)

	for i := 0; i < elemType.NumField(); i++ { // 此时肯定是个结构体
		field := elemType.Field(i)
		m.extractField(field, no) // 递归
	}
}
