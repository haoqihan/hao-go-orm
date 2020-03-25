package hao_go_orm

import (
	"database/sql"
	"hao-go-orm/dialect"
	"hao-go-orm/log"
	"hao-go-orm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// NewEngine 创建 Engine 实例时，获取 driver 对应的 dialect。
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Error("dialect %s Not Found", driver)
		return
	}

	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")

	}
	log.Info("Close database success")
}

// NewSession 创建 Session 实例时，传递 dialect 给构造函数 New。
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}
