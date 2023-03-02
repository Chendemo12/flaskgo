package godantic

import (
	"reflect"
	"strings"
	"unicode"
)

var metaDataFactory = &MetaClass{data: make([]*MetaData, 0)}

func GetMetaDataFactory() *MetaClass   { return metaDataFactory }
func GetMetaData(pkg string) *MetaData { return metaDataFactory.Get(pkg) }
func SetMetaData(data *MetaData)       { metaDataFactory.Set(data) }

type MetaField struct {
	Field
	Index     int          `json:"index" description:"当前字段所处的序号"`
	Exported  bool         `json:"exported" description:"是否是导出字段"`
	Anonymous bool         `json:"anonymous" description:"是否是嵌入字段"`
	RType     reflect.Type `description:"反射字段类型"`
}

type MetaData struct {
	names       []string     `description:"结构体名称,结构体名称,包名.结构体名称"`
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

func (m *MetaClass) structFieldToMetaField(field reflect.StructField) *MetaField {
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
	//meta.innerFields = append(meta.innerFields, mf)
	return mf
}

// 反射模型获取字段信息
func (m *MetaClass) reflectModel(rtype reflect.Type, meta *MetaData, no int) {
	if rtype.Kind() == reflect.Interface {
		return
	}
	if rtype.Kind() == reflect.Pointer {
		rtype = rtype.Elem()
	}

	switch reflectKindToOType(rtype.Kind()) {
	case IntegerType:
		//tianjia
		return
	case StringType:
		return
	case BoolType:
		return
	case NumberType:
		return
	}

	for i := 0; i < rtype.NumField(); i++ { // 此时rtype必须不是指针
		field := rtype.Field(i)

		if field.Anonymous && (field.Name == "BaseModel" || field.Name == "Field") { // 过滤模型基类
			continue
		}
		if strings.HasPrefix(field.Name, "_") { // 过滤约定的匿名字段
			continue
		}

		// ---------------------------------- 获取字段信息 ----------------------------------

		fieldMeta := m.structFieldToMetaField(field)

		if no > 0 { // 递归获取的内部字段
			meta.innerFields = append(meta.innerFields, fieldMeta)
		} else {
			meta.fields = append(meta.fields, fieldMeta)
		}

		switch fieldMeta.OType {
		case ArrayType:
			fieldMeta.ItemRef = field.Name
			no += 1
			m.reflectModel(field.Type.Elem(), meta, no)

		case ObjectType:
			no += 1
			m.reflectModel(field.Type, meta, no)

		default:
			continue // 获取下一个字段信息
		}
	}
}

// Reflect 反射建立任意类型的元信息
func (m *MetaClass) Reflect(model SchemaIface) *MetaData {
	if nm, ok := model.(*Field); ok {
		if GetMetaData(nm._pkg) != nil {
			return nm.MetaData()
		}
	}

	at := reflect.TypeOf(model)
	if at.Kind() == reflect.Pointer { // 全部为指针类型
		at = at.Elem()
	}

	meta := &MetaData{
		names:       []string{at.Name(), at.String()}, // 获取包名
		fields:      make([]*MetaField, 0),
		innerFields: make([]*MetaField, 0),
	}

	//m.reflectModel(at, meta, 0)

	for i := 0; i < at.NumField(); i++ { // 此时rtype必须不是指针
		field := at.Field(i)

		if field.Anonymous && (field.Name == "BaseModel" || field.Name == "Field") { // 过滤模型基类
			continue
		}
		if strings.HasPrefix(field.Name, "_") { // 过滤约定的匿名字段
			continue
		}

		// ---------------------------------- 获取字段信息 ----------------------------------

		fieldMeta := m.structFieldToMetaField(field)
		meta.fields = append(meta.fields, fieldMeta)

		switch fieldMeta.OType {
		case ArrayType:
			fieldMeta.ItemRef = field.Name
			m.reflectModel(field.Type.Elem(), meta, 1)

		case ObjectType:
			m.reflectModel(field.Type, meta, 1)

		default:
			continue // 获取下一个字段信息
		}
	}

	return meta
}
