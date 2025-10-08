package test_db

import (
	"fmt"
	"slices"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
)

// GenerateEntity creates random entity and returns it.
func GenerateEntity[T any](opts ...func(entity *T)) T {
	var entity T

	_ = gofakeit.Struct(&entity)

	for _, opt := range opts {
		opt(&entity)
	}

	return entity
}

// GenerateEntities creates n random entities and returns them.
func GenerateEntities[T any](n int, opts ...func(entity *T, i int)) []T {
	entities := make([]T, n)

	for i := range entities {
		_ = gofakeit.Struct(&entities[i])

		for _, opt := range opts {
			opt(&entities[i], i)
		}
	}

	return entities
}

func insertEntities[U, T any](c *Container, table string, entities []T, column string, columnSkip bool) ([]U, error) {
	if len(entities) == 0 {
		return make([]U, 0), nil
	}

	fieldsSlice := lo.Map(entities, func(item T, _ int) []field {
		fields := toFields(item)
		if columnSkip {
			fields = slices.DeleteFunc(fields, func(field field) bool {
				return field.key == column
			})
		}
		return fields
	})

	keys := lo.Map(fieldsSlice[0], func(item field, _ int) string { return item.key })

	builder := c.builder.
		Insert(table).
		Columns(keys...)

	for _, fields := range fieldsSlice {
		values := lo.Map(fields, func(item field, _ int) any { return item.value })
		builder = builder.Values(values...)
	}

	builder = builder.Suffix("RETURNING " + column)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for %q: %w", table, err)
	}

	var ids []U
	err = c.db.Select(&ids, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query for %q: %w", table, err)
	}

	return ids, nil
}

// InsertEntitiesSkipID skips fields with tag "id", inserts entities into table and returns ids from column "id".
func InsertEntitiesSkipID[U, T any](c *Container, table string, entities []T) ([]U, error) {
	return insertEntities[U](c, table, entities, "id", true)
}

// InsertEntitiesSkipColumn skips fields with tag column, inserts entities into table and returns values from column.
func InsertEntitiesSkipColumn[U, T any](c *Container, table string, entities []T, column string) ([]U, error) {
	return insertEntities[U](c, table, entities, column, true)
}

// InsertEntitiesWithID inserts entities into table and returns ids from column "id".
func InsertEntitiesWithID[U, T any](c *Container, table string, entities []T) ([]U, error) {
	return insertEntities[U](c, table, entities, "id", false)
}

// InsertEntitiesWithColumn inserts entities into table and returns values from column.
func InsertEntitiesWithColumn[U, T any](c *Container, table string, entities []T, column string) ([]U, error) {
	return insertEntities[U](c, table, entities, column, false)
}

func insertEntity[U, T any](c *Container, table string, entity T, column string, columnSkip bool) (U, error) {
	ids, err := insertEntities[U](c, table, []T{entity}, column, columnSkip)
	if err != nil {
		var id U
		return id, err
	}
	return ids[0], nil
}

// InsertEntitySkipID skips field with tag "id", inserts entity into table and returns id from column "id".
func InsertEntitySkipID[U, T any](c *Container, table string, entity T) (U, error) {
	return insertEntity[U](c, table, entity, "id", true)
}

// InsertEntitySkipColumn skips field with tag column, inserts entity into table and returns value from column.
func InsertEntitySkipColumn[U, T any](c *Container, table string, entity T, column string) (U, error) {
	return insertEntity[U](c, table, entity, column, true)
}

// InsertEntityWithID inserts entity into table and returns id from column "id".
func InsertEntityWithID[U, T any](c *Container, table string, entity T) (U, error) {
	return insertEntity[U](c, table, entity, "id", false)
}

// InsertEntityWithColumn inserts entity into table and returns value from column.
func InsertEntityWithColumn[U, T any](c *Container, table string, entity T, column string) (U, error) {
	return insertEntity[U](c, table, entity, column, false)
}

// DeleteEntitiesByColumn deletes entities from table by ids in column.
func DeleteEntitiesByColumn[U any](c *Container, table string, column string, values []U) error {
	builder := c.builder.
		Delete(table).
		Where(sq.Eq{column: values})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build query for %q: %w", table, err)
	}

	_, err = c.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("exec query for %q: %w", table, err)
	}

	return nil
}

// DeleteEntitiesByID deletes entities from table by ids in column "id".
func DeleteEntitiesByID[U any](c *Container, table string, ids []U) error {
	return DeleteEntitiesByColumn(c, table, "id", ids)
}

// DeleteEntityByColumn deletes entity from table by value in column.
func DeleteEntityByColumn[U any](c *Container, table string, column string, value U) error {
	return DeleteEntitiesByColumn(c, table, column, []U{value})
}

// DeleteEntityByID deletes entity from table by id in column "id".
func DeleteEntityByID[U any](c *Container, table string, id U) error {
	return DeleteEntitiesByID(c, table, []U{id})
}

// SelectEntitiesByColumn selects entities from table by values in column.
func SelectEntitiesByColumn[T, U any](c *Container, table string, column string, values []U) ([]T, error) {
	var entity T
	fields := toFields(entity)
	keys := lo.Map(fields, func(item field, _ int) string { return item.key })

	builder := c.builder.
		Select(keys...).
		From(table).
		Where(sq.Eq{column: values})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for %q: %w", table, err)
	}

	var entities []T
	err = c.db.Select(&entities, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query for %q: %w", table, err)
	}

	return entities, nil
}

// SelectEntitiesByID selects entities from table by ids in column "id".
func SelectEntitiesByID[T, U any](c *Container, table string, ids []U) ([]T, error) {
	return SelectEntitiesByColumn[T](c, table, "id", ids)
}
