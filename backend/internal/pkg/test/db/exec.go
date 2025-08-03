package test_db

import (
	"fmt"
	"slices"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
)

type GenerateOption[T any] func(entity *T)

// GenerateEntity creates random entity and returns it.
func GenerateEntity[T any](opts ...GenerateOption[T]) (T, error) {
	var entity T

	if err := gofakeit.Struct(&entity); err != nil {
		return entity, nil
	}

	for _, opt := range opts {
		opt(&entity)
	}

	return entity, nil
}

// GenerateEntities creates n random entities and returns them.
func GenerateEntities[T any](n int, opts ...GenerateOption[T]) ([]T, error) {
	entities := make([]T, n)

	var err error
	for i := range entities {
		entities[i], err = GenerateEntity(opts...)
		if err != nil {
			return nil, err
		}
	}

	return entities, nil
}

func insertEntities[T, U any](c *Container, table string, entities []T, col string, idSkip bool) ([]U, error) {
	if len(entities) == 0 {
		return make([]U, 0), nil
	}

	fieldsSlice := lo.Map(entities, func(item T, _ int) []field {
		fields := toFields(item)
		if idSkip {
			fields = slices.DeleteFunc(fields, func(field field) bool {
				return field.key == col
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

	builder = builder.Suffix("RETURNING " + col)

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

// InsertEntitiesById skips fields with tag "id", inserts entities into table and returns ids from column "id".
func InsertEntitiesById[T, U any](c *Container, table string, entities []T) ([]U, error) {
	return insertEntities[T, U](c, table, entities, "id", true)
}

// InsertEntitiesByCol skips fields with tag сol, inserts entities into table and returns ids from column сol.
func InsertEntitiesByCol[T, U any](c *Container, table string, entities []T, сol string) ([]U, error) {
	return insertEntities[T, U](c, table, entities, сol, true)
}

// InsertEntitiesWithId inserts entities into table and returns ids from column "id".
func InsertEntitiesWithId[T, U any](c *Container, table string, entities []T) ([]U, error) {
	return insertEntities[T, U](c, table, entities, "id", false)
}

// InsertEntitiesWithCol inserts entities into table and returns ids from column сol.
func InsertEntitiesWithCol[T, U any](c *Container, table string, entities []T, col string) ([]U, error) {
	return insertEntities[T, U](c, table, entities, col, false)
}

func insertEntity[T, U any](c *Container, table string, entity T, col string, idSkip bool) (U, error) {
	ids, err := insertEntities[T, U](c, table, []T{entity}, col, idSkip)
	if err != nil {
		var id U
		return id, err
	}
	return ids[0], nil
}

// InsertEntityById skips field with tag "id", inserts entity into table and returns id from column "id".
func InsertEntityById[T, U any](c *Container, table string, entity T) (U, error) {
	return insertEntity[T, U](c, table, entity, "id", true)
}

// InsertEntityByCol skips field with tag сol, inserts entity into table and returns id from column сol.
func InsertEntityByCol[T, U any](c *Container, table string, entity T, сol string) (U, error) {
	return insertEntity[T, U](c, table, entity, сol, true)
}

// InsertEntityWithId inserts entity into table and returns id from column "id".
func InsertEntityWithId[T, U any](c *Container, table string, entity T) (U, error) {
	return insertEntity[T, U](c, table, entity, "id", false)
}

// InsertEntityWithCol inserts entity into table and returns id from column сol.
func InsertEntityWithCol[T, U any](c *Container, table string, entity T, сol string) (U, error) {
	return insertEntity[T, U](c, table, entity, сol, false)
}

// DeleteEntitiesByCol deletes entities from table by ids in column сol.
func DeleteEntitiesByCol[U any](c *Container, table string, сol string, ids []U) error {
	builder := c.builder.
		Delete(table).
		Where(sq.Eq{сol: ids})

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

// DeleteEntitiesById deletes entities from table by ids in column "id".
func DeleteEntitiesById[U any](c *Container, table string, ids []U) error {
	return DeleteEntitiesByCol(c, table, "id", ids)
}

// DeleteEntityByCol deletes entity from table by id in column сol.
func DeleteEntityByCol[U any](c *Container, table string, сol string, id U) error {
	return DeleteEntitiesByCol(c, table, сol, []U{id})
}

// DeleteEntityById deletes entity from table by id in column "id".
func DeleteEntityById[U any](c *Container, table string, id U) error {
	return DeleteEntitiesById(c, table, []U{id})
}

// SelectEntitiesByCol selects entities from table by ids in column сol
func SelectEntitiesByCol[T, U any](c *Container, table string, сol string, ids []U) ([]T, error) {
	var entity T
	fields := toFields(entity)
	keys := lo.Map(fields, func(item field, _ int) string { return item.key })

	builder := c.builder.
		Select(keys...).
		From(table).
		Where(sq.Eq{сol: ids})

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

// SelectEntitiesById selects entities from table by ids in column "id".
func SelectEntitiesById[T, U any](c *Container, table string, ids []U) ([]T, error) {
	return SelectEntitiesByCol[T, U](c, table, "id", ids)
}
