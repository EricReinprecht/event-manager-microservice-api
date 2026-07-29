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

func (g *GormExecutor) Model(value any) DBExecutor {
	return &GormExecutor{
		db: g.db.Model(value),
	}
}

func (g *GormExecutor) Where(query any, args ...any) DBExecutor {
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

func (g *GormExecutor) Order(value any) DBExecutor {
	return &GormExecutor{
		db: g.db.Order(value),
	}
}

func (g *GormExecutor) Group(
	query string,
) DBExecutor {

	return &GormExecutor{
		db: g.db.Group(query),
	}
}

func (g *GormExecutor) First(dest any, conds ...any) DBExecutor {
	return &GormExecutor{
		db: g.db.First(dest, conds...),
	}
}

func (g *GormExecutor) Find(dest any, conds ...any) DBExecutor {
	return &GormExecutor{
		db: g.db.Find(dest, conds...),
	}
}

func (g *GormExecutor) Preload(query string, args ...any) DBExecutor {
	return &GormExecutor{
		db: g.db.Preload(query, args...),
	}
}

func (g *GormExecutor) Joins(query string, args ...any) DBExecutor {
	return &GormExecutor{
		db: g.db.Joins(query, args...),
	}
}

func (g *GormExecutor) Select(query any, args ...any) DBExecutor {
	return &GormExecutor{
		db: g.db.Select(query, args...),
	}
}

func (g *GormExecutor) Scan(dest any) DBExecutor {
	return &GormExecutor{
		db: g.db.Scan(dest),
	}
}

func (g *GormExecutor) Create(value any) DBExecutor {
	return &GormExecutor{
		db: g.db.Create(value),
	}
}

func (e *GormExecutor) Updates(values any) DBExecutor {
	return &GormExecutor{
		db: e.db.Updates(values),
	}
}

func (g *GormExecutor) Save(value any) DBExecutor {
	return &GormExecutor{
		db: g.db.Save(value),
	}
}

func (g *GormExecutor) Delete(value any) DBExecutor {
	return &GormExecutor{
		db: g.db.Delete(value),
	}
}

func (g *GormExecutor) Association(
	column string,
) *gorm.Association {

	return g.db.Association(column)
}

func (g *GormExecutor) Raw(
	query string,
	args ...any,
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
