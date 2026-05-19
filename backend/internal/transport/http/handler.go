package http

import (
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
)

type Handler struct {
	*ProjectCreateHandler
	*ProjectDeleteHandler
	*ProjectGetByIDHandler
	*ProjectListHandler
	*ProjectUpdateHandler
	*ProjectUpdateUsersHandler
	*ProjectUsersHandler
	*TaskCreateHandler
	*TaskGetByIDHandler
	*TaskListHandler
	*TemplateCreateHandler
	*TemplateCreateFromDefaultHandler
	*TemplateDefaultListHandler
	*TemplateDeleteHandler
	*TemplateGetByIDHandler
	*TemplateGetMetaByIDHandler
	*TemplateImportHandler
	*TemplateListHandler
	*TemplateUpdateHandler
	*TemplateUpdateUsersHandler
	*TemplateUsersHandler
	*UserCreateHandler
	*UserGetByIDHandler
	*UserListHandler
	*UserTokenCreateHandler
	*UserTokenDeleteHandler
	*VersionCreateHandler
	*VersionCreateFromHandler
	*VersionListHandler
}

type (
	ProjectCreateHandler             = project_create_handler.Handler
	ProjectDeleteHandler             = project_delete_handler.Handler
	ProjectGetByIDHandler            = project_get_by_id_handler.Handler
	ProjectListHandler               = project_list_handler.Handler
	ProjectUpdateHandler             = project_update_handler.Handler
	ProjectUpdateUsersHandler        = project_update_users_handler.Handler
	ProjectUsersHandler              = project_users_handler.Handler
	TaskCreateHandler                = task_create_handler.Handler
	TaskGetByIDHandler               = task_get_by_id_handler.Handler
	TaskListHandler                  = task_list_handler.Handler
	TemplateCreateHandler            = template_create_handler.Handler
	TemplateCreateFromDefaultHandler = template_create_from_default_handler.Handler
	TemplateDefaultListHandler       = template_default_list_handler.Handler
	TemplateDeleteHandler            = template_delete_handler.Handler
	TemplateGetByIDHandler           = template_get_by_id_handler.Handler
	TemplateGetMetaByIDHandler       = template_get_meta_by_id_handler.Handler
	TemplateImportHandler            = template_import_handler.Handler
	TemplateListHandler              = template_list_handler.Handler
	TemplateUpdateHandler            = template_update_handler.Handler
	TemplateUpdateUsersHandler       = template_update_users_handler.Handler
	TemplateUsersHandler             = template_users_handler.Handler
	UserCreateHandler                = user_create_handler.Handler
	UserGetByIDHandler               = user_get_by_id_handler.Handler
	UserListHandler                  = user_list_handler.Handler
	UserTokenCreateHandler           = user_token_create_handler.Handler
	UserTokenDeleteHandler           = user_token_delete_handler.Handler
	VersionCreateHandler             = version_create_handler.Handler
	VersionCreateFromHandler         = version_create_from_handler.Handler
	VersionListHandler               = version_list_handler.Handler
)
