package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	version_get_domain "github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
	task_process_domain "github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
	data_process_service "github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/service/data_process"
	variable_process_service "github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/service/variable_process"
)

type templateInput struct {
	Name      string
	IsDefault bool
	ProjectID *int64
	AuthorID  *int64
	CreatedAt time.Time
	UpdatedAt *time.Time
	Version   *versionInput
}

type versionInput struct {
	Data      string
	CreatedAt time.Time
	Variables []variableInput
}

type variableInput struct {
	Name        string
	Title       string
	Type        variable_domain.Type
	IsInput     bool
	Expression  string // пустая строка == NULL
	Constraints []constraintInput
}

type constraintInput struct {
	Name       string
	Expression string
	IsActive   bool
}

type taskInput struct {
	Title     string // только для логов
	Payload   map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time // также используется как момент завершения обработки
}

// createTemplate вставляет template, при наличии — version, variable
// и variable_constraint, после чего обновляет last_version_id и возвращает
// id шаблона. Если CreatedAt не задан, используется now() — это нормально
// для шаблонов, к которым не привязана история.
func (s *seeder) createTemplate(ctx context.Context, in templateInput) (int64, error) {
	createdAt := in.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	var templateID int64
	err := s.db.GetContext(ctx,
		&templateID,
		`INSERT INTO template (name, is_default, project_id, author_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id`,
		in.Name, in.IsDefault, in.ProjectID, in.AuthorID, createdAt, in.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("insert template: %w", err)
	}

	if in.Version == nil {
		return templateID, nil
	}

	versionID, err := s.createVersion(ctx, templateID, in.AuthorID, *in.Version)
	if err != nil {
		return 0, fmt.Errorf("create version: %w", err)
	}

	// last_version_id выставляется без затрагивания updated_at — иначе
	// «бесшовно созданный» шаблон сразу получил бы updated_at = now().
	_, err = s.db.ExecContext(ctx,
		`UPDATE template SET last_version_id = $1 WHERE id = $2`,
		versionID, templateID,
	)
	if err != nil {
		return 0, fmt.Errorf("update template last_version_id: %w", err)
	}

	return templateID, nil
}

func (s *seeder) createVersion(ctx context.Context, templateID int64, authorID *int64, v versionInput) (int64, error) {
	createdAt := v.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	var versionID int64
	err := s.db.GetContext(ctx,
		&versionID,
		`INSERT INTO template_version (number, template_id, author_id, data, created_at)
		 VALUES (1, $1, $2, $3, $4)
		 RETURNING id`,
		templateID, authorID, []byte(v.Data), createdAt,
	)
	if err != nil {
		return 0, fmt.Errorf("insert version: %w", err)
	}

	for _, variable := range v.Variables {
		var variableID int64
		var expr any
		if variable.Expression != "" {
			expr = variable.Expression
		}
		err := s.db.GetContext(ctx,
			&variableID,
			`INSERT INTO variable (version_id, name, title, type, expression, is_input)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 RETURNING id`,
			versionID, variable.Name, variable.Title, string(variable.Type), expr, variable.IsInput,
		)
		if err != nil {
			return 0, fmt.Errorf("insert variable %q: %w", variable.Name, err)
		}

		if len(variable.Constraints) == 0 {
			continue
		}

		b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
			Insert("variable_constraint").
			Columns("variable_id", "name", "expression", "is_active")
		for _, c := range variable.Constraints {
			b = b.Values(variableID, c.Name, c.Expression, c.IsActive)
		}
		query, args, err := b.ToSql()
		if err != nil {
			return 0, fmt.Errorf("build constraint query: %w", err)
		}
		if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
			return 0, fmt.Errorf("insert constraints for %q: %w", variable.Name, err)
		}
	}

	return versionID, nil
}

// createAndProcessTask вставляет задачу и сразу запускает рендеринг через
// те же сервисы рендеринга шаблона и переменных, что и продакшен-воркер,
// сохраняя либо result, либо error напрямую в БД — без обращения к
// RabbitMQ.
func (s *seeder) createAndProcessTask(ctx context.Context, creatorID, versionID int64, in taskInput) (int64, error) {
	payloadJSON, err := json.Marshal(in.Payload)
	if err != nil {
		return 0, fmt.Errorf("marshal payload: %w", err)
	}

	createdAt := in.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	updatedAt := in.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = createdAt.Add(7 * time.Second)
	}

	var taskID int64
	err = s.db.GetContext(ctx,
		&taskID,
		`INSERT INTO task (version_id, status, payload, creator_id, created_at)
		 VALUES ($1, 'created', $2::jsonb, $3, $4)
		 RETURNING id`,
		versionID, payloadJSON, creatorID, createdAt,
	)
	if err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}

	version, err := s.loadVersion(ctx, versionID)
	if err != nil {
		return 0, fmt.Errorf("load version: %w", err)
	}

	result, processErr := s.processTask(ctx, version, in.Payload)

	switch {
	case processErr != nil:
		errJSON, err := json.Marshal(processErr)
		if err != nil {
			return 0, fmt.Errorf("marshal process error: %w", err)
		}
		_, err = s.db.ExecContext(ctx,
			`UPDATE task SET status = 'failed', error = $1::jsonb, updated_at = $2
			 WHERE id = $3`,
			errJSON, updatedAt, taskID,
		)
		if err != nil {
			return 0, fmt.Errorf("update task with error: %w", err)
		}

	default:
		var resultID int64
		err = s.db.GetContext(ctx,
			&resultID,
			`INSERT INTO result (data) VALUES ($1) RETURNING id`,
			result,
		)
		if err != nil {
			return 0, fmt.Errorf("insert result: %w", err)
		}

		_, err = s.db.ExecContext(ctx,
			`UPDATE task SET status = 'succeed', result_id = $1, updated_at = $2
			 WHERE id = $3`,
			resultID, updatedAt, taskID,
		)
		if err != nil {
			return 0, fmt.Errorf("update task with result: %w", err)
		}
	}

	return taskID, nil
}

func (s *seeder) processTask(ctx context.Context, version version_get_domain.Version, payload map[string]string) ([]byte, *task_domain.ProcessError) {
	values, err := variable_process_service.New().Handle(ctx, task_process_domain.VariableProcessIn{
		Variables: version.Variables,
		Payload:   payload,
	})
	if err != nil {
		var pe *task_domain.ProcessError
		if errors.As(err, &pe) {
			return nil, pe
		}
		return nil, &task_domain.ProcessError{Message: err.Error()}
	}

	data, err := data_process_service.New().Handle(ctx, task_process_domain.DataProcessIn{
		Values: values,
		Data:   version.Data,
	})
	if err != nil {
		var pe *task_domain.ProcessError
		if errors.As(err, &pe) {
			return nil, pe
		}
		return nil, &task_domain.ProcessError{Message: err.Error()}
	}

	return data, nil
}

// loadVersion вытаскивает версию шаблона со всеми переменными и ограничениями,
// в виде, совместимом с task_process сервисами.
func (s *seeder) loadVersion(ctx context.Context, versionID int64) (version_get_domain.Version, error) {
	var v struct {
		ID         int64  `db:"id"`
		TemplateID int64  `db:"template_id"`
		Number     int64  `db:"number"`
		Data       []byte `db:"data"`
	}
	err := s.db.GetContext(ctx, &v,
		`SELECT id, template_id, number, data FROM template_version WHERE id = $1`, versionID)
	if err != nil {
		return version_get_domain.Version{}, fmt.Errorf("select version: %w", err)
	}

	var variables []struct {
		ID         int64   `db:"id"`
		Name       string  `db:"name"`
		Title      string  `db:"title"`
		Type       string  `db:"type"`
		Expression *string `db:"expression"`
		IsInput    bool    `db:"is_input"`
	}
	err = s.db.SelectContext(ctx, &variables,
		`SELECT id, name, title, type, expression, is_input
		 FROM variable WHERE version_id = $1 ORDER BY id`, versionID)
	if err != nil {
		return version_get_domain.Version{}, fmt.Errorf("select variables: %w", err)
	}

	out := version_get_domain.Version{
		ID:         v.ID,
		TemplateID: v.TemplateID,
		Number:     v.Number,
		Data:       v.Data,
		Variables:  make([]version_get_domain.Variable, 0, len(variables)),
	}

	for _, vv := range variables {
		var constraints []struct {
			ID         int64  `db:"id"`
			Name       string `db:"name"`
			Expression string `db:"expression"`
			IsActive   bool   `db:"is_active"`
		}
		err = s.db.SelectContext(ctx, &constraints,
			`SELECT id, name, expression, is_active
			 FROM variable_constraint WHERE variable_id = $1 ORDER BY id`, vv.ID)
		if err != nil {
			return version_get_domain.Version{}, fmt.Errorf("select constraints: %w", err)
		}

		domainConstraints := make([]version_get_domain.Constraint, 0, len(constraints))
		for _, c := range constraints {
			domainConstraints = append(domainConstraints, version_get_domain.Constraint{
				ID:         c.ID,
				VariableID: vv.ID,
				Name:       c.Name,
				Expression: c.Expression,
				IsActive:   c.IsActive,
			})
		}

		out.Variables = append(out.Variables, version_get_domain.Variable{
			ID:          vv.ID,
			Name:        vv.Name,
			Title:       vv.Title,
			Type:        variable_domain.Type(vv.Type),
			Expression:  vv.Expression,
			IsInput:     vv.IsInput,
			Constraints: domainConstraints,
		})
	}

	return out, nil
}
