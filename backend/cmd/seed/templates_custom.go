package main

import (
	"time"

	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
)

type customTemplate struct {
	Key       string // машинный ключ для адресации в коде сидера
	Name      string // отображаемое имя
	CreatedAt time.Time
	UpdatedAt *time.Time
	Version   *versionInput
}

// customTemplates возвращает кастомные шаблоны проекта NotifyHub: набор
// рабочих шаблонов команды (passport, releaseCard, ADR) и три полноформатных
// шаблона документации по ГОСТ с заполненными задачами.
func customTemplates() []customTemplate {
	return []customTemplate{
		passportTemplate(),
		releaseCardTemplate(),
		adrTemplate(),
		approbationTZASTemplate(),
		approbationPMITemplate(),
		approbationOPTemplate(),
	}
}

func passportTemplate() customTemplate {
	return customTemplate{
		Key:       "passport",
		Name:      "Технический паспорт микросервиса",
		CreatedAt: datePassportCreated,
		UpdatedAt: datePassportUpdated,
		Version: &versionInput{
			CreatedAt: datePassportVersion,
			Variables: []variableInput{
				{Name: "service_name", Title: "Название сервиса", Type: variable_domain.TypeString, IsInput: true},
				{Name: "service_owner", Title: "Команда-владелец", Type: variable_domain.TypeString, IsInput: true},
				{Name: "service_version", Title: "Версия сервиса", Type: variable_domain.TypeString, IsInput: true},
				{Name: "production_date", Title: "Дата ввода в эксплуатацию", Type: variable_domain.TypeString, IsInput: true},
				{Name: "runtime", Title: "Runtime", Type: variable_domain.TypeString, IsInput: true},
				{
					Name: "sla_target", Title: "Целевой SLA (0..1)",
					Type: variable_domain.TypeFloat, IsInput: true,
					Constraints: []constraintInput{
						{Name: "положительный_sla", Expression: "sla_target > 0", IsActive: true},
						{Name: "production_grade_sla", Expression: "sla_target >= 0.9", IsActive: true},
					},
				},
				{
					Name: "rps_target", Title: "Целевой RPS",
					Type: variable_domain.TypeInteger, IsInput: true,
					Constraints: []constraintInput{
						{Name: "положительный_rps", Expression: "rps_target > 0", IsActive: true},
					},
				},
				{
					Name: "latency_p50_ms", Title: "Латентность p50 (мс)",
					Type: variable_domain.TypeInteger, IsInput: true,
					Constraints: []constraintInput{
						{Name: "положительная_латентность", Expression: "latency_p50_ms > 0", IsActive: true},
					},
				},
				{Name: "latency_p95_ms", Title: "Латентность p95 (мс)", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "latency_p99_ms", Title: "Латентность p99 (мс)", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "dep_db", Title: "База данных", Type: variable_domain.TypeString, IsInput: true},
				{Name: "dep_queue", Title: "Брокер сообщений", Type: variable_domain.TypeString, IsInput: true},
				{Name: "dep_cache", Title: "Кеш", Type: variable_domain.TypeString, IsInput: true},
				{
					Name: "instances", Title: "Количество реплик",
					Type: variable_domain.TypeInteger, IsInput: true,
					Constraints: []constraintInput{
						{Name: "минимум_одна_реплика", Expression: "instances >= 1", IsActive: true},
					},
				},
				{Name: "cpu_per_instance", Title: "vCPU на реплику", Type: variable_domain.TypeFloat, IsInput: true},
				{Name: "memory_mb_per_instance", Title: "Память на реплику (МБ)", Type: variable_domain.TypeInteger, IsInput: true},

				// computed
				{
					Name: "availability_pct", Title: "Доступность, %",
					Type: variable_domain.TypeString, IsInput: false,
					Expression: "percent(sla_target, 3)",
				},
				{
					Name: "error_budget_min_per_month", Title: "Бюджет ошибок (мин/30 дн)",
					Type: variable_domain.TypeFloat, IsInput: false,
					Expression: "round((1 - sla_target) * 30 * 24 * 60, 1)",
				},
				{
					Name: "error_budget_sci", Title: "Бюджет ошибок (научная запись)",
					Type: variable_domain.TypeString, IsInput: false,
					Expression: "scientific(1 - sla_target, 2)",
				},
				{
					Name: "latency_jump_ms", Title: "Разрыв p99 − p95 (мс)",
					Type: variable_domain.TypeInteger, IsInput: false,
					Expression: "latency_p99_ms - latency_p95_ms",
				},
				{
					Name: "peak_rps", Title: "Пиковая нагрузка (RPS)",
					Type: variable_domain.TypeFloat, IsInput: false,
					Expression: "round(float(rps_target) * 1.5, 0)",
				},
				{
					Name: "total_cpu", Title: "Суммарный CPU (vCPU)",
					Type: variable_domain.TypeFloat, IsInput: false,
					Expression: "round(float(instances) * cpu_per_instance, 2)",
				},
				{
					Name: "total_memory_gb", Title: "Суммарная память (ГБ)",
					Type: variable_domain.TypeFloat, IsInput: false,
					Expression: "round(float(instances) * float(memory_mb_per_instance) / 1024.0, 2)",
				},
				{
					Name: "load_class", Title: "Класс нагрузки",
					Type: variable_domain.TypeString, IsInput: false,
					Expression: `rps_target > 1000 ? "высокая" : "стандартная"`,
				},
				{
					Name: "service_slug", Title: "Идентификатор для CI/CD",
					Type: variable_domain.TypeString, IsInput: false,
					Expression: `lower(replace(service_name, " ", "-"))`,
				},
			},
			Data: readTemplate("passport.md"),
		},
	}
}

func releaseCardTemplate() customTemplate {
	return customTemplate{
		Key:       "release",
		Name:      "Карточка релиза",
		CreatedAt: dateReleaseCardCreated,
		Version: &versionInput{
			CreatedAt: dateReleaseCardVersion,
			Variables: []variableInput{
				{Name: "release_version", Title: "Версия релиза", Type: variable_domain.TypeString, IsInput: true},
				{Name: "release_date", Title: "Дата релиза", Type: variable_domain.TypeString, IsInput: true},
				{Name: "release_owner", Title: "Релиз-мастер", Type: variable_domain.TypeString, IsInput: true},
				{
					Name: "features_added", Title: "Добавлено фич",
					Type: variable_domain.TypeInteger, IsInput: true,
					Constraints: []constraintInput{
						{Name: "не_отрицательное", Expression: "features_added >= 0", IsActive: true},
					},
				},
				{
					Name: "bugs_fixed", Title: "Исправлено багов",
					Type: variable_domain.TypeInteger, IsInput: true,
					Constraints: []constraintInput{
						{Name: "не_отрицательное", Expression: "bugs_fixed >= 0", IsActive: true},
					},
				},
				{Name: "has_breaking", Title: "Есть breaking changes (да/нет)", Type: variable_domain.TypeString, IsInput: true},

				{
					Name: "total_changes", Title: "Всего изменений",
					Type: variable_domain.TypeInteger, IsInput: false,
					Expression: "features_added + bugs_fixed",
				},
				{
					Name: "release_slug", Title: "Идентификатор релиза",
					Type: variable_domain.TypeString, IsInput: false,
					Expression: `"release-" + lower(replace(release_version, ".", "-"))`,
				},
			},
			Data: readTemplate("release_card.md"),
		},
	}
}

func adrTemplate() customTemplate {
	return customTemplate{
		Key:       "adr",
		Name:      "ADR — архитектурное решение",
		CreatedAt: dateADRCreated,
		Version: &versionInput{
			CreatedAt: dateADRVersion,
			Variables: []variableInput{
				{
					Name: "adr_number", Title: "Номер ADR",
					Type: variable_domain.TypeInteger, IsInput: true,
					Constraints: []constraintInput{
						{Name: "положительный_номер", Expression: "adr_number > 0", IsActive: true},
					},
				},
				{Name: "adr_title", Title: "Заголовок решения", Type: variable_domain.TypeString, IsInput: true},
				{Name: "adr_status", Title: "Статус", Type: variable_domain.TypeString, IsInput: true},
				{Name: "adr_date", Title: "Дата", Type: variable_domain.TypeString, IsInput: true},
				{Name: "decision_owner", Title: "Ответственный за решение", Type: variable_domain.TypeString, IsInput: true},
				{Name: "context", Title: "Контекст", Type: variable_domain.TypeString, IsInput: true},
				{Name: "decision", Title: "Решение", Type: variable_domain.TypeString, IsInput: true},
				{Name: "consequences", Title: "Последствия", Type: variable_domain.TypeString, IsInput: true},

				{
					Name: "adr_filename", Title: "Имя файла",
					Type: variable_domain.TypeString, IsInput: false,
					Expression: `"ADR-" + string(adr_number) + "-" + lower(replace(adr_title, " ", "-"))`,
				},
			},
			Data: readTemplate("adr.md"),
		},
	}
}
