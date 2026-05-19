// Команда seed заливает в Postgres готовый набор данных для локальной
// разработки: пользователей, библиотеку стандартных шаблонов по ГОСТ,
// проекты с кастомными шаблонами, версиями и задачами. Перед заливкой все
// таблицы очищаются через TRUNCATE ... RESTART IDENTITY CASCADE, поэтому
// сидер можно повторно запускать на той же БД.
//
// Запуск:
//
//	cd backend && \
//	PGHOST=localhost PGPORT=15432 PGDATABASE=master \
//	PGUSER=devuser PGPASSWORD=devpassword \
//	go run ./cmd/seed
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/qsoulior/tech-generator/backend/internal/pkg/postgres"
)

const (
	demoUserName     = "demo"
	demoUserEmail    = "demo@example.com"
	demoUserPassword = "Demo1234!"

	// Второй пользователь — нужен, чтобы в данных были примеры межпользовательского
	// доступа: владелец одного из проектов и автор одного из шаблонов
	// в чужом проекте.
	architectUserName     = "platform_lead"
	architectUserEmail    = "platform.lead@example.com"
	architectUserPassword = "Atlas2026!"

	demoProjectName      = "Сервис уведомлений NotifyHub"
	demoProjectName2     = "Платёжный шлюз PayBridge"
	architectProjectName = "API-Gateway Atlas"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	if err := godotenv.Overload(); err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Error("overload env", slog.String("err", err.Error()))
		return 1
	}

	db, err := postgres.Connect(ctx)
	if err != nil {
		logger.Error("connect postgres", slog.String("err", err.Error()))
		return 1
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("close postgres connection", slog.String("err", err.Error()))
		}
	}()

	s := &seeder{db: db, logger: logger}

	if err := s.run(ctx); err != nil {
		logger.Error("seed failed", slog.String("err", err.Error()))
		return 1
	}

	logger.Info("seed completed",
		slog.String("user", demoUserName),
		slog.String("password", demoUserPassword),
		slog.String("project", demoProjectName),
	)
	return 0
}

type seeder struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func (s *seeder) run(ctx context.Context) error {
	if err := s.truncate(ctx); err != nil {
		return fmt.Errorf("truncate: %w", err)
	}
	s.logger.Info("truncated tables")

	userID, err := s.createUser(ctx, demoUserName, demoUserEmail, demoUserPassword)
	if err != nil {
		return fmt.Errorf("create demo user: %w", err)
	}
	s.logger.Info("created user", slog.String("name", demoUserName), slog.Int64("id", userID))

	architectID, err := s.createUser(ctx, architectUserName, architectUserEmail, architectUserPassword)
	if err != nil {
		return fmt.Errorf("create architect user: %w", err)
	}
	s.logger.Info("created user", slog.String("name", architectUserName), slog.Int64("id", architectID))

	for _, t := range defaultTemplates() {
		id, err := s.createTemplate(ctx, templateInput{
			Name:      t.Name,
			IsDefault: true,
			ProjectID: nil,
			AuthorID:  nil,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
			Version:   t.Version,
		})
		if err != nil {
			return fmt.Errorf("create default template %q: %w", t.Name, err)
		}
		s.logger.Info("created default template", slog.String("name", t.Name), slog.Int64("id", id))
	}

	projectID, err := s.createProject(ctx, demoProjectName, userID)
	if err != nil {
		return fmt.Errorf("create project NotifyHub: %w", err)
	}
	s.logger.Info("created project", slog.String("name", demoProjectName), slog.Int64("id", projectID))

	if _, err := s.createProject(ctx, demoProjectName2, userID); err != nil {
		return fmt.Errorf("create project PayBridge: %w", err)
	}
	s.logger.Info("created project", slog.String("name", demoProjectName2))

	// Проект второго пользователя, в котором demo участвует с ролью read,
	// чтобы продемонстрировать механизм совместного доступа к проектам.
	atlasID, err := s.createProject(ctx, architectProjectName, architectID)
	if err != nil {
		return fmt.Errorf("create project Atlas: %w", err)
	}
	if err := s.shareProject(ctx, atlasID, userID, "read"); err != nil {
		return fmt.Errorf("share project Atlas with demo: %w", err)
	}
	s.logger.Info("created project", slog.String("name", architectProjectName), slog.Int64("id", atlasID))

	custom := customTemplates()
	templateIDs := make(map[string]int64, len(custom))
	for _, t := range custom {
		// Один из шаблонов (ADR) принадлежит второму пользователю и расшарен
		// автору демо-данных через template_user, чтобы в данных был пример
		// шеринга шаблона между пользователями.
		authorID := userID
		if t.Key == "adr" {
			authorID = architectID
		}

		id, err := s.createTemplate(ctx, templateInput{
			Name:      t.Name,
			IsDefault: false,
			ProjectID: &projectID,
			AuthorID:  &authorID,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
			Version:   t.Version,
		})
		if err != nil {
			return fmt.Errorf("create custom template %q: %w", t.Name, err)
		}
		templateIDs[t.Key] = id

		if t.Key == "adr" {
			if err := s.shareTemplate(ctx, id, userID, "read"); err != nil {
				return fmt.Errorf("share template ADR with demo: %w", err)
			}
		}

		s.logger.Info("created custom template", slog.String("name", t.Name), slog.Int64("id", id))
	}

	passportTemplateID := templateIDs["passport"]
	versionID, err := s.lastVersionID(ctx, passportTemplateID)
	if err != nil {
		return fmt.Errorf("get last version for passport: %w", err)
	}

	for _, tk := range passportTasks() {
		taskID, err := s.createAndProcessTask(ctx, userID, versionID, tk)
		if err != nil {
			return fmt.Errorf("create task %q: %w", tk.Title, err)
		}
		s.logger.Info("created task", slog.String("title", tk.Title), slog.Int64("id", taskID))
	}

	approbationTasks := []struct {
		templateKey string
		task        taskInput
	}{
		{"approbation_tz_as", approbationTZTask()},
		{"approbation_pmi", approbationPMITask()},
		{"approbation_op", approbationOPTask()},
	}

	for _, a := range approbationTasks {
		tid, ok := templateIDs[a.templateKey]
		if !ok {
			return fmt.Errorf("template %q was not created", a.templateKey)
		}
		vid, err := s.lastVersionID(ctx, tid)
		if err != nil {
			return fmt.Errorf("get last version for %q: %w", a.templateKey, err)
		}
		taskID, err := s.createAndProcessTask(ctx, userID, vid, a.task)
		if err != nil {
			return fmt.Errorf("create task %q: %w", a.task.Title, err)
		}
		s.logger.Info("created task", slog.String("title", a.task.Title), slog.Int64("id", taskID))
	}

	return nil
}

func (s *seeder) truncate(ctx context.Context) error {
	const stmt = `
		TRUNCATE TABLE
			task, result,
			variable_constraint, variable,
			template_version, template_user, template,
			project_user, project,
			usr
		RESTART IDENTITY CASCADE
	`
	_, err := s.db.ExecContext(ctx, stmt)
	return err
}

func (s *seeder) createUser(ctx context.Context, name, email, password string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var id int64
	err = s.db.GetContext(ctx,
		&id,
		`INSERT INTO usr (name, email, password) VALUES ($1, $2, $3) RETURNING id`,
		name, email, hash,
	)
	return id, err
}

func (s *seeder) createProject(ctx context.Context, name string, authorID int64) (int64, error) {
	var id int64
	err := s.db.GetContext(ctx,
		&id,
		`INSERT INTO project (name, author_id) VALUES ($1, $2) RETURNING id`,
		name, authorID,
	)
	return id, err
}

func (s *seeder) shareProject(ctx context.Context, projectID, userID int64, role string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO project_user (project_id, user_id, role) VALUES ($1, $2, $3)`,
		projectID, userID, role,
	)
	return err
}

func (s *seeder) shareTemplate(ctx context.Context, templateID, userID int64, role string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO template_user (template_id, user_id, role) VALUES ($1, $2, $3)`,
		templateID, userID, role,
	)
	return err
}

func (s *seeder) lastVersionID(ctx context.Context, templateID int64) (int64, error) {
	var id int64
	err := s.db.GetContext(ctx,
		&id,
		`SELECT last_version_id FROM template WHERE id = $1`,
		templateID,
	)
	return id, err
}
