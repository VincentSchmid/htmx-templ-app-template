package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type IGenericRepository[T any] interface {
	Create(entity *T) error
	GetById(id int) (*T, error)
	GetByUuid(uuid uuid.UUID) (*T, error)
	GetByField(field string, value interface{}) (*T, error)
	GetManyByField(field string, value interface{}) ([]T, error)
	Update(entity *T) error
	Delete(id int) error
	ListAll() ([]T, error)
}

type GenericRepository[T any] struct {
	db *bun.DB
}

func NewGenericRepository[T any](db *bun.DB) *GenericRepository[T] {
	return &GenericRepository[T]{db: db}
}

func (r *GenericRepository[T]) Create(entity *T) error {
	_, err := r.db.NewInsert().
		Model(entity).
		Exec(context.Background())

	if err != nil {
		return fmt.Errorf("failed to create entity: %w", err)
	}

	return nil
}

func (r *GenericRepository[T]) GetById(id int) (*T, error) {
	return r.GetByField("id", id)
}

func (r *GenericRepository[T]) GetByUuid(uuid uuid.UUID) (*T, error) {
	return r.GetByField("uuid", uuid)
}

func (r *GenericRepository[T]) GetByField(field string, value interface{}) (*T, error) {
	var entity T

	err := r.db.NewSelect().
		Model(&entity).
		Where(fmt.Sprintf("%s = ?", field), value).
		Scan(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to read entity: %w", err)
	}

	return &entity, nil
}

func (r *GenericRepository[T]) GetManyByField(field string, value interface{}) ([]T, error) {
	var entities []T

	err := r.db.NewSelect().
		Model(&entities).
		Where(fmt.Sprintf("%s = ?", field), value).
		Scan(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to read entities: %w", err)
	}

	return entities, nil
}

func (r *GenericRepository[T]) Update(entity *T) error {
	_, err := r.db.NewUpdate().
		Model(entity).
		WherePK().
		Exec(context.Background())

	if err != nil {
		return fmt.Errorf("failed to update entity: %w", err)
	}

	return nil
}

func (r *GenericRepository[T]) Delete(id int) error {
	var entity T

	_, err := r.db.NewDelete().
		Model(&entity).
		Where("id = ?", id).
		Exec(context.Background())

	return fmt.Errorf("failed to delete entity: %w", err)
}

func (r *GenericRepository[T]) ListAll() ([]T, error) {
	var entities []T

	err := r.db.NewSelect().
		Model(&entities).
		Scan(context.Background())

	if err != nil {
		return entities, fmt.Errorf("failed to list entities: %w", err)
	}

	return entities, nil
}
