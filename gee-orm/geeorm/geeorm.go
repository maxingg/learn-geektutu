package geeorm

import (
	"database/sql"
	"fmt"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/session"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

type TxFunc func(*session.Session) (interface{}, error)

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
		log.Errorf("dialect %s not found", driver)
		return
	}
	e = &Engine{
		db:      db,
		dialect: dial,
	}
	log.Info("Connect database success")
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			defer func() {
				if err != nil {
					_ = s.Rollback()
				}
			}()
			err = s.Commit()
		}
	}()
	return f(s)
}

// return a - b
func difference(a, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}

	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// 数据库中表结构与table对应的结构体是不一样的
func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		// Model是加载到内存的，此时已经更新了表的结构体
		if !s.Model(value).HasTable() {
			log.Infof("table %s does not exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Info("added cols %v, deleted cols %v", addCols, delCols)
		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}

		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}
