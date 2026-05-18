package http

import (
	project_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_create"
	project_delete_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_delete"
	project_get_by_id_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_get_by_id"
	project_list_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_list"
	project_update_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/project_update"
	task_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/task_create"
	task_get_by_id_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/task_get_by_id"
	task_list_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/task_list"
	template_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_create"
	template_delete_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_delete"
	template_get_by_id_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_get_by_id"
	template_list_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_list"
	template_update_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/template_update"
	user_create_handler "github.com/qsoulior/tech-generator/backend/internal/transport/http/handler/user_create"
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
	*TaskCreateHandler
	*TaskGetByIDHandler
	*TaskListHandler
	*TemplateCreateHandler
	*TemplateDeleteHandler
	*TemplateGetByIDHandler
	*TemplateListHandler
	*TemplateUpdateHandler
	*UserCreateHandler
	*UserGetByIDHandler
	*UserTokenCreateHandler
	*UserTokenDeleteHandler
	*VersionCreateHandler
	*VersionCreateFromHandler
	*VersionListHandler
}

type (
	ProjectCreateHandler     = project_create_handler.Handler
	ProjectDeleteHandler     = project_delete_handler.Handler
	ProjectGetByIDHandler    = project_get_by_id_handler.Handler
	ProjectListHandler       = project_list_handler.Handler
	ProjectUpdateHandler     = project_update_handler.Handler
	TaskCreateHandler        = task_create_handler.Handler
	TaskGetByIDHandler       = task_get_by_id_handler.Handler
	TaskListHandler          = task_list_handler.Handler
	TemplateCreateHandler    = template_create_handler.Handler
	TemplateDeleteHandler    = template_delete_handler.Handler
	TemplateGetByIDHandler   = template_get_by_id_handler.Handler
	TemplateListHandler      = template_list_handler.Handler
	TemplateUpdateHandler    = template_update_handler.Handler
	UserCreateHandler        = user_create_handler.Handler
	UserGetByIDHandler       = user_get_by_id_handler.Handler
	UserTokenCreateHandler   = user_token_create_handler.Handler
	UserTokenDeleteHandler   = user_token_delete_handler.Handler
	VersionCreateHandler     = version_create_handler.Handler
	VersionCreateFromHandler = version_create_from_handler.Handler
	VersionListHandler       = version_list_handler.Handler
)
