package database

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBExecutor interface {
	WithContext(ctx context.Context) DBExecutor

	Model(value any) DBExecutor

	Where(query any, args ...any) DBExecutor

	Clauses(conds ...clause.Expression) DBExecutor

	Limit(limit int) DBExecutor

	Offset(offset int) DBExecutor

	Group(query string) DBExecutor

	Order(value any) DBExecutor

	First(dest any, conds ...any) DBExecutor

	Find(dest any, conds ...any) DBExecutor

	Preload(query string, args ...any) DBExecutor

	Joins(query string, args ...any) DBExecutor

	Select(query any, args ...any) DBExecutor

	Scan(dest any) DBExecutor

	Create(value any) DBExecutor

	Updates(values any) DBExecutor

	Save(value any) DBExecutor

	Delete(value any) DBExecutor

	Association(column string) *gorm.Association

	Raw(query string, args ...any) DBExecutor

	Count(count *int64) DBExecutor

	Begin() DBExecutor

	Commit() error

	Rollback() error

	Error() error

	RowsAffected() int64
}
