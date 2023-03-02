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

type MetaData struct {
	names       []string `description:"结构体名称,结构体名称,包名.结构体名称"`
	fields      []*Field `description:"结构体字段"`
	innerFields []*Field `description:"内部字段"`
}

// Name 获取结构体名称
func (m MetaData) Name() string { return m.names[0] }

// String 结构体全称：包名+结构体名称
func (m MetaData) String() string { return m.names[1] }

// Fields 结构体字段
func (m MetaData) Fields() []*Field { return m.fields }

// InnerFields 内部字段
func (m MetaData) InnerFields() []*Field { return m.innerFields }

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

// Reflect 反射建立任意类型的元信息
func (m *MetaClass) Reflect(model any) *MetaData {
	at := reflect.TypeOf(model) // 全部为指针类型
	if at.Kind() == reflect.Pointer {
		at = at.Elem()
	}

	meta := &MetaData{
		names:       []string{at.Name(), at.String()}, // 获取包名
		fields:      make([]*Field, 0),
		innerFields: make([]*Field, 0),
	}

	// 获取字段信息
	for i := 0; i < at.NumField(); i++ {
		field := at.Field(i)
		if field.Anonymous || strings.HasPrefix(field.Name, "_") {
			continue
		}
		// TODO: 处理嵌入式结构体
		meta.fields = append(meta.fields, &Field{
			Index:     i,
			Title:     field.Name,
			Tag:       field.Tag,
			RType:     field.Type,
			Exported:  unicode.IsUpper(rune(field.Name[0])),
			Anonymous: field.Anonymous,
		})
	}

	return meta
}
