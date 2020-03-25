package schema

import (
	"go/ast"
	"hao-go-orm/dialect"
	"reflect"
)

// 字段表示数据库的列
type Field struct {
	Name string // 字段名
	Type string // 类型
	Tag  string // 约束条件
}

// 表示数据库的表
type Schema struct {
	Model      interface{}       // 映射对象
	Name       string            // 表名
	Fields     []*Field          // 字段
	FieldNames []string          //包含所有的字段名(列名)
	fieldMap   map[string]*Field // 记录字段名和Field的映射关系
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// reflect.Indirect() 获取指针指向的实例。
// modelType.Name() 获取到结构体的名称作为表名。
// NumField() 获取实例的字段的个数，然后通过下标获取到特定字段 p := modelType.Field(i)。
// p.Name 即字段名，p.Type 即字段类型，通过 (Dialect).DataTypeOf() 转换为数据库的字段类型，p.Tag 即额外的约束条件。
// 将任意对象解析为Schema
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("horm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
