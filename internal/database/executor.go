package database

import "context"

type DBExecutor interface {
	WithContext(ctx context.Context) DBExecutor

	Model(value interface{}) DBExecutor

	Where(query interface{}, args ...interface{}) DBExecutor

	First(dest interface{}, conds ...interface{}) DBExecutor

	Find(dest interface{}, conds ...interface{}) DBExecutor

	Preload(query string, args ...interface{}) DBExecutor

	Joins(query string, args ...interface{}) DBExecutor

	Select(query interface{}, args ...interface{}) DBExecutor

	Scan(dest interface{}) DBExecutor

	Create(value interface{}) DBExecutor

	Save(value interface{}) DBExecutor

	Delete(value interface{}) DBExecutor

	Raw(query string, args ...interface{}) DBExecutor

	Count(count *int64) DBExecutor

	Begin() DBExecutor

	Commit() error

	Rollback() error

	Error() error
}
