package database

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormExecutor struct {
	db *gorm.DB
}

func NewGormExecutor(db *gorm.DB) *GormExecutor {
	return &GormExecutor{
		db: db,
	}
}

func (g *GormExecutor) WithContext(ctx context.Context) DBExecutor {
	return &GormExecutor{
		db: g.db.WithContext(ctx),
	}
}

func (g *GormExecutor) Model(value interface{}) DBExecutor {
	return &GormExecutor{
		db: g.db.Model(value),
	}
}

func (g *GormExecutor) Where(query interface{}, args ...interface{}) DBExecutor {
	return &GormExecutor{
		db: g.db.Where(query, args...),
	}
}

func (g *GormExecutor) Clauses(
	conds ...clause.Expression,
) DBExecutor {

	return &GormExecutor{
		db: g.db.Clauses(conds...),
	}
}

func (g *GormExecutor) Limit(limit int) DBExecutor {
	return &GormExecutor{
		db: g.db.Limit(limit),
	}
}

func (g *GormExecutor) Offset(offset int) DBExecutor {
	return &GormExecutor{
		db: g.db.Offset(offset),
	}
}

func (g *GormExecutor) Order(value interface{}) DBExecutor {
	return &GormExecutor{
		db: g.db.Order(value),
	}
}

func (g *GormExecutor) First(dest interface{}, conds ...interface{}) DBExecutor {
	return &GormExecutor{
		db: g.db.First(dest, conds...),
	}
}

func (g *GormExecutor) Find(dest interface{}, conds ...interface{}) DBExecutor {
	return &GormExecutor{
		db: g.db.Find(dest, conds...),
	}
}

func (g *GormExecutor) Preload(query string, args ...interface{}) DBExecutor {
	return &GormExecutor{
		db: g.db.Preload(query, args...),
	}
}

func (g *GormExecutor) Joins(query string, args ...interface{}) DBExecutor {
	return &GormExecutor{
		db: g.db.Joins(query, args...),
	}
}

func (g *GormExecutor) Select(query interface{}, args ...interface{}) DBExecutor {
	return &GormExecutor{
		db: g.db.Select(query, args...),
	}
}

func (g *GormExecutor) Scan(dest interface{}) DBExecutor {
	return &GormExecutor{
		db: g.db.Scan(dest),
	}
}

func (g *GormExecutor) Create(value interface{}) DBExecutor {
	return &GormExecutor{
		db: g.db.Create(value),
	}
}

func (e *GormExecutor) Updates(values interface{}) DBExecutor {
	return &GormExecutor{
		db: e.db.Updates(values),
	}
}

func (g *GormExecutor) Save(value interface{}) DBExecutor {
	return &GormExecutor{
		db: g.db.Save(value),
	}
}

func (g *GormExecutor) Delete(value interface{}) DBExecutor {
	return &GormExecutor{
		db: g.db.Delete(value),
	}
}

func (g *GormExecutor) Raw(
	query string,
	args ...interface{},
) DBExecutor {

	return &GormExecutor{
		db: g.db.Raw(query, args...),
	}
}

func (g *GormExecutor) Count(count *int64) DBExecutor {
	return &GormExecutor{
		db: g.db.Count(count),
	}
}

func (g *GormExecutor) Begin() DBExecutor {
	return &GormExecutor{
		db: g.db.Begin(),
	}
}

func (g *GormExecutor) Commit() error {
	return g.db.Commit().Error
}

func (g *GormExecutor) Rollback() error {
	return g.db.Rollback().Error
}

func (g *GormExecutor) Error() error {
	return g.db.Error
}

func (e *GormExecutor) RowsAffected() int64 {
	return e.db.RowsAffected
}
