package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/qsoulior/tech-generator/backend/internal/config"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
	"github.com/qsoulior/tech-generator/backend/internal/pkg/ed25519key"
	"github.com/qsoulior/tech-generator/backend/internal/pkg/httpserver"
	"github.com/qsoulior/tech-generator/backend/internal/pkg/postgres"
	"github.com/qsoulior/tech-generator/backend/internal/pkg/rabbitmq"
	"github.com/qsoulior/tech-generator/backend/internal/transport/http"
	error_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/error"
	project_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_create"
	project_delete_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_delete"
	project_get_by_id_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_get_by_id"
	project_list_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_list"
	project_update_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_update"
	project_update_users_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_update_users"
	project_users_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_users"
	task_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/task_create"
	task_get_by_id_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/task_get_by_id"
	task_list_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/task_list"
	template_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_create"
	template_create_from_default_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_create_from_default"
	template_default_list_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_default_list"
	template_delete_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_delete"
	template_get_by_id_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_get_by_id"
	template_get_meta_by_id_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_get_meta_by_id"
	template_import_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_import"
	template_list_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_list"
	template_update_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_update"
	template_update_users_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_update_users"
	template_users_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_users"
	user_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/user_create"
	user_list_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/user_list"
	user_get_by_id_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/user_get_by_id"
	user_token_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/user_token_create"
	user_token_delete_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/user_token_delete"
	version_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/version_create"
	version_create_from_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/version_create_from"
	version_list_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/version_list"
	auth_middleware "github.com/qsoulior/tech-generator/backend/internal/transport/http/middleware/auth"
	project_create_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/project_create"
	project_delete_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/project_delete"
	project_get_by_id_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/project_get_by_id"
	project_list_by_user_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/project_list_by_user"
	project_update_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/project_update"
	project_user_list_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_list"
	project_user_update_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/project_user_update"
	task_create_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/task_create"
	task_get_by_id_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/task_get_by_id"
	task_list_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/task_list"
	template_create_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_create"
	template_create_from_default_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_create_from_default"
	template_delete_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_delete"
	template_get_by_id_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id"
	template_get_meta_by_id_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_meta_by_id"
	template_import_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_import"
	template_list_by_user_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_by_user"
	template_list_default_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_default"
	template_update_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_update"
	template_user_list_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_list"
	template_user_update_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/template_user_update"
	user_create_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/user_create"
	user_list_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/user_list"
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

	cfg, err := config.New()
	if err != nil {
		logger.Error("init config", slog.String("err", err.Error()))
		return 1
	}

	db, err := postgres.Connect(ctx)
	if err != nil {
		logger.Error("connect postgres", slog.String("err", err.Error()))
		return 1
	}
	defer func() {
		err := db.Close()
		if err != nil {
			logger.Error("close postgres connection", slog.String("err", err.Error()))
			code = 1
		}
	}()

	conn, err := rabbitmq.Connect()
	if err != nil {
		logger.Error("connect rabbitmq", slog.String("err", err.Error()))
		return 1
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Error("close rabbitmq connection", slog.String("err", err.Error()))
			code = 1
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("connect rabbitmq channel", slog.String("err", err.Error()))
		return 1
	}
	defer func() {
		err := ch.Close()
		if err != nil {
			logger.Error("close rabbitmq channel", slog.String("err", err.Error()))
			code = 1
		}
	}()

	_, err = ch.QueueDeclare("task_created", true, false, false, false, nil)
	if err != nil {
		logger.Error("declare queue", slog.String("err", err.Error()))
		return 1
	}

	privateKey, err := ed25519key.LoadPrivateKey(cfg.Ed25519PrivateKeyPath)
	if err != nil {
		logger.Error("load ed25519 private key", slog.String("err", err.Error()))
		return 1
	}

	publicKey, err := ed25519key.LoadPublicKey(cfg.Ed25519PublicKeyPath)
	if err != nil {
		logger.Error("load ed25519 public key", slog.String("err", err.Error()))
		return 1
	}

	projectCreateUsecase := project_create_usecase.New(db)
	projectDeleteUsecase := project_delete_usecase.New(db)
	projectGetByIDUsecase := project_get_by_id_usecase.New(db)
	projectListUsecase := project_list_by_user_usecase.New(db)
	projectUpdateUsecase := project_update_usecase.New(db)
	projectUserListUsecase := project_user_list_usecase.New(db)
	projectUserUpdateUsecase := project_user_update_usecase.New(db)
	taskCreateUsecase := task_create_usecase.New(db, ch)
	taskGetByIDUsecase := task_get_by_id_usecase.New(db)
	taskListUsecase := task_list_usecase.New(db)
	templateCreateUsecase := template_create_usecase.New(db)
	templateCreateFromDefaultUsecase := template_create_from_default_usecase.New(db)
	templateDefaultListUsecase := template_list_default_usecase.New(db)
	templateDeleteUsecase := template_delete_usecase.New(db)
	templateGetByIDUsecase := template_get_by_id_usecase.New(db)
	templateGetMetaByIDUsecase := template_get_meta_by_id_usecase.New(db)
	templateImportUsecase := template_import_usecase.New(db)
	templateListUsecase := template_list_by_user_usecase.New(db)
	templateUpdateUsecase := template_update_usecase.New(db)
	templateUserListUsecase := template_user_list_usecase.New(db)
	templateUserUpdateUsecase := template_user_update_usecase.New(db)
	userCreateUsecase := user_create_usecase.New(db)
	userGetByIDUsecase := user_get_by_id_usecase.New(db)
	userListUsecase := user_list_usecase.New(db)
	userTokenCreateUsecase := user_token_create_usecase.New(db, privateKey, cfg)
	userTokenParseUsecase := user_token_parse_usecase.New(publicKey)
	versionCreateUsecase := version_create_usecase.New(db)
	versionCreateFromUsecase := version_create_from_usecase.New(db)
	versionListUsecase := version_list_usecase.New(db)

	apiHandler := &http.Handler{
		ProjectCreateHandler:             project_create_handler.New(projectCreateUsecase),
		ProjectDeleteHandler:             project_delete_handler.New(projectDeleteUsecase),
		ProjectGetByIDHandler:            project_get_by_id_handler.New(projectGetByIDUsecase),
		ProjectListHandler:               project_list_handler.New(projectListUsecase),
		ProjectUpdateHandler:             project_update_handler.New(projectUpdateUsecase),
		ProjectUpdateUsersHandler:        project_update_users_handler.New(projectUserUpdateUsecase),
		ProjectUsersHandler:              project_users_handler.New(projectUserListUsecase),
		TaskCreateHandler:                task_create_handler.New(taskCreateUsecase),
		TaskGetByIDHandler:               task_get_by_id_handler.New(taskGetByIDUsecase),
		TaskListHandler:                  task_list_handler.New(taskListUsecase),
		TemplateCreateHandler:            template_create_handler.New(templateCreateUsecase),
		TemplateCreateFromDefaultHandler: template_create_from_default_handler.New(templateCreateFromDefaultUsecase),
		TemplateDefaultListHandler:       template_default_list_handler.New(templateDefaultListUsecase),
		TemplateDeleteHandler:            template_delete_handler.New(templateDeleteUsecase),
		TemplateGetByIDHandler:           template_get_by_id_handler.New(templateGetByIDUsecase),
		TemplateGetMetaByIDHandler:       template_get_meta_by_id_handler.New(templateGetMetaByIDUsecase),
		TemplateImportHandler:            template_import_handler.New(templateImportUsecase),
		TemplateListHandler:              template_list_handler.New(templateListUsecase),
		TemplateUpdateHandler:            template_update_handler.New(templateUpdateUsecase),
		TemplateUpdateUsersHandler:       template_update_users_handler.New(templateUserUpdateUsecase),
		TemplateUsersHandler:             template_users_handler.New(templateUserListUsecase),
		UserCreateHandler:                user_create_handler.New(userCreateUsecase),
		UserGetByIDHandler:               user_get_by_id_handler.New(userGetByIDUsecase),
		UserListHandler:                  user_list_handler.New(userListUsecase),
		UserTokenCreateHandler:           user_token_create_handler.New(userTokenCreateUsecase),
		UserTokenDeleteHandler:           user_token_delete_handler.New(),
		VersionCreateHandler:             version_create_handler.New(versionCreateUsecase),
		VersionCreateFromHandler:         version_create_from_handler.New(versionCreateFromUsecase),
		VersionListHandler:               version_list_handler.New(versionListUsecase),
	}

	apiServer, err := api.NewServer(apiHandler,
		api.WithErrorHandler(error_handler.New(logger).Handle),
	)
	if err != nil {
		logger.Error("create api server", slog.String("err", err.Error()))
		return 1
	}

	authMiddleware := auth_middleware.New(userTokenParseUsecase, logger)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   cfg.ServiceAllowedOrigins,
		AllowedHeaders:   []string{"Content-Type", "Accept"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "DELETE"},
		AllowCredentials: true,
	})

	httpHandler := authMiddleware.Handle(apiServer)
	httpHandler = corsMiddleware.Handler(httpHandler)

	server := httpserver.New(httpHandler, logger)
	if err := server.Run(ctx); err != nil {
		logger.Error("fail server", slog.String("err", err.Error()))
		return 1
	}

	logger.Info("stop server")
	return 0
}
