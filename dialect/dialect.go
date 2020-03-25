package dialect

import "reflect"

type Dialect interface {
	DataTypeOf(typ reflect.Value) string                    // 用于将GO语言的类型转换为该数据库的数据类型
	TableExistSQL(tableName string) (string, []interface{}) // 返回某个表示否存在的sql语句，参数是表名(table)
}

var dialectsMap = map[string]Dialect{}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
