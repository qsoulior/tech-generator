package main

import (
	"context"
	"crypto/ed25519"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/qsoulior/tech-generator/backend/internal/config"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/pkg/httpserver"
	"github.com/qsoulior/tech-generator/backend/internal/pkg/postgres"
	"github.com/qsoulior/tech-generator/backend/internal/transport/http"
	error_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/error"
	project_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_create"
	project_delete_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_delete"
	project_list_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_list"
	task_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/task_create"
	task_get_by_id_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/task_get_by_id"
	task_list_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/task_list"
	template_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_create"
	template_delete_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_delete"
	template_get_by_id_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_get_by_id"
	template_list_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_list"
	user_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/user_create"
	user_get_by_id_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/user_get_by_id"
	user_token_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/user_token_create"
	version_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/version_create"
	version_create_from_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/version_create_from"
	version_list_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/version_list"
	auth_middleware "github.com/qsoulior/tech-generator/backend/internal/transport/http/middleware/auth"
	project_create_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/project_create"
	project_delete_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/project_delete"
	project_list_by_user_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user"
	task_create_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/task_create"
	task_get_by_id_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/task_get_by_id"
	task_list_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/task_list"
	template_create_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_create"
	template_delete_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_delete"
	template_get_by_id_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id"
	template_list_by_user_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_by_user"
	user_create_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/user_create"
	user_get_by_id_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/user_get_by_id"
	user_token_create_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_create"
	user_token_parse_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/user_token_parse"
	version_create_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/version_create"
	version_create_from_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/version_create_from"
	version_list_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/version_list"
)

func main() {
	os.Exit(run())
}

func run() (code int) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	err := godotenv.Overload()
	if err != nil {
		logger.Error("overload env", slog.String("err", err.Error()))
		return 1
	}

	db, err := postgres.Connect(ctx)
	if err != nil {
		logger.Error("connect postgres", slog.String("err", err.Error()))
		return 1
	}
	defer func() {
		err := db.Close()
		logger.Error("close postgres connection", slog.String("err", err.Error()))
		code = 1
	}()

	cfg := &config.Config{
		UserTokenExpiration: 30 * 24 * time.Hour,
	}

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		logger.Error("generate ed25519 key", slog.String("err", err.Error()))
		return 1
	}

	projectCreateUsecase := project_create_usecase.New(db)
	projectDeleteUsecase := project_delete_usecase.New(db)
	projectListUsecase := project_list_by_user_usecase.New(db)
	taskCreateUsecase := task_create_usecase.New(db)
	taskGetByIDUsecase := task_get_by_id_usecase.New(db)
	taskListUsecase := task_list_usecase.New(db)
	templateCreateUsecase := template_create_usecase.New(db)
	templateDeleteUsecase := template_delete_usecase.New(db)
	templateGetByIDUsecase := template_get_by_id_usecase.New(db)
	templateListUsecase := template_list_by_user_usecase.New(db)
	userCreateUsecase := user_create_usecase.New(db)
	userGetByIDUsecase := user_get_by_id_usecase.New(db)
	userTokenCreateUsecase := user_token_create_usecase.New(db, privateKey, cfg)
	userTokenParseUsecase := user_token_parse_usecase.New(publicKey)
	versionCreateUsecase := version_create_usecase.New(db)
	versionCreateFromUsecase := version_create_from_usecase.New(db)
	versionListUsecase := version_list_usecase.New(db)

	apiHandler := &http.Handler{
		ProjectCreateHandler:     project_create_handler.New(projectCreateUsecase),
		ProjectDeleteHandler:     project_delete_handler.New(projectDeleteUsecase),
		ProjectListHandler:       project_list_handler.New(projectListUsecase),
		TaskCreateHandler:        task_create_handler.New(taskCreateUsecase),
		TaskGetByIDHandler:       task_get_by_id_handler.New(taskGetByIDUsecase),
		TaskListHandler:          task_list_handler.New(taskListUsecase),
		TemplateCreateHandler:    template_create_handler.New(templateCreateUsecase),
		TemplateDeleteHandler:    template_delete_handler.New(templateDeleteUsecase),
		TemplateGetByIDHandler:   template_get_by_id_handler.New(templateGetByIDUsecase),
		TemplateListHandler:      template_list_handler.New(templateListUsecase),
		UserCreateHandler:        user_create_handler.New(userCreateUsecase),
		UserGetByIDHandler:       user_get_by_id_handler.New(userGetByIDUsecase),
		UserTokenCreateHandler:   user_token_create_handler.New(userTokenCreateUsecase),
		VersionCreateHandler:     version_create_handler.New(versionCreateUsecase),
		VersionCreateFromHandler: version_create_from_handler.New(versionCreateFromUsecase),
		VersionListHandler:       version_list_handler.New(versionListUsecase),
	}

	apiServer, err := api.NewServer(apiHandler,
		api.WithErrorHandler(error_handler.New(logger).Handle),
	)
	if err != nil {
		logger.Error("create api server", slog.String("err", err.Error()))
		return 1
	}

	authMiddleware := auth_middleware.New(userTokenParseUsecase, logger)

	httpHandler := authMiddleware.Handle()(apiServer)

	server := httpserver.New(httpHandler, logger)
	if err := server.Run(ctx); err != nil {
		logger.Error("fail server", slog.String("err", err.Error()))
		return 1
	}

	logger.Info("stop server")
	return 0
}
