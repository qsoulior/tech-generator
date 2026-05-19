package main

import variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"

// approbationTZASTemplate — Техническое задание на АС по ГОСТ 34.602-2020.
func approbationTZASTemplate() customTemplate {
	return customTemplate{
		Key:       "approbation_tz_as",
		Name:      "ТЗ на создание АС NotifyHub (ГОСТ 34.602-2020)",
		CreatedAt: dateApprobationTZCreated,
		Version: &versionInput{
			CreatedAt: dateApprobationTZVersion,
			Variables: []variableInput{
				{Name: "system_name", Title: "Полное наименование системы", Type: variable_domain.TypeString, IsInput: true},
				{Name: "system_code", Title: "Условное обозначение", Type: variable_domain.TypeString, IsInput: true},
				{Name: "document_version", Title: "Версия документа", Type: variable_domain.TypeString, IsInput: true},
				{Name: "document_date", Title: "Дата документа", Type: variable_domain.TypeString, IsInput: true},
				{Name: "customer", Title: "Заказчик", Type: variable_domain.TypeString, IsInput: true},
				{Name: "customer_owner", Title: "Ответственный от заказчика", Type: variable_domain.TypeString, IsInput: true},
				{Name: "developer", Title: "Исполнитель", Type: variable_domain.TypeString, IsInput: true},
				{Name: "developer_owner", Title: "Ответственный от исполнителя", Type: variable_domain.TypeString, IsInput: true},
				{Name: "tz_basis", Title: "Основание для разработки", Type: variable_domain.TypeString, IsInput: true},
				{Name: "start_date", Title: "Дата начала разработки", Type: variable_domain.TypeString, IsInput: true},
				{Name: "finish_date", Title: "Срок ввода в действие", Type: variable_domain.TypeString, IsInput: true},
				{Name: "system_purpose", Title: "Назначение системы", Type: variable_domain.TypeString, IsInput: true},
				{Name: "system_goals", Title: "Цели создания (через `;`)", Type: variable_domain.TypeString, IsInput: true},
				{Name: "users_endusers", Title: "Кол-во конечных получателей", Type: variable_domain.TypeString, IsInput: true},
				{Name: "users_operators", Title: "Кол-во операторов", Type: variable_domain.TypeString, IsInput: true},
				{Name: "systems_sources", Title: "Кол-во систем-источников событий", Type: variable_domain.TypeString, IsInput: true},
				{Name: "sla_target", Title: "Целевой SLA", Type: variable_domain.TypeString, IsInput: true},
				{Name: "max_downtime_min", Title: "Допустимый downtime, мин/мес", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "subsystems", Title: "Перечень подсистем (через `;`)", Type: variable_domain.TypeString, IsInput: true},
				{Name: "runtime", Title: "Среда выполнения", Type: variable_domain.TypeString, IsInput: true},
				{Name: "k8s_version", Title: "Версия Kubernetes", Type: variable_domain.TypeString, IsInput: true},
				{Name: "db_version", Title: "Версия СУБД", Type: variable_domain.TypeString, IsInput: true},
				{Name: "min_replicas", Title: "Минимум реплик сервиса", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "training_hours", Title: "Часы обучения операторов", Type: variable_domain.TypeInteger, IsInput: true},

				{
					Name: "document_id", Title: "Шифр документа",
					Type: variable_domain.TypeString, IsInput: false,
					Expression: `"ТЗ-" + system_code + "-v" + document_version`,
				},
			},
			Data: readTemplate("approbation_tz_as.md"),
		},
	}
}

// approbationPMITemplate — Программа и методика испытаний по ГОСТ 19.301-79.
func approbationPMITemplate() customTemplate {
	return customTemplate{
		Key:       "approbation_pmi",
		Name:      "Программа и методика испытаний NotifyHub (ГОСТ 19.301-79)",
		CreatedAt: dateApprobationPMICreated,
		Version: &versionInput{
			CreatedAt: dateApprobationPMIVersion,
			Variables: []variableInput{
				{Name: "program_name", Title: "Наименование программы", Type: variable_domain.TypeString, IsInput: true},
				{Name: "program_version", Title: "Версия программы", Type: variable_domain.TypeString, IsInput: true},
				{Name: "document_code", Title: "Шифр документа", Type: variable_domain.TypeString, IsInput: true},
				{Name: "test_date_start", Title: "Дата начала испытаний", Type: variable_domain.TypeString, IsInput: true},
				{Name: "test_date_end", Title: "Дата окончания испытаний", Type: variable_domain.TypeString, IsInput: true},
				{Name: "test_location", Title: "Место проведения", Type: variable_domain.TypeString, IsInput: true},
				{Name: "test_stand", Title: "Стенд", Type: variable_domain.TypeString, IsInput: true},
				{Name: "responsible", Title: "Ответственный", Type: variable_domain.TypeString, IsInput: true},
				{Name: "approver", Title: "Утверждающий", Type: variable_domain.TypeString, IsInput: true},
				{Name: "tz_reference", Title: "Ссылка на ТЗ", Type: variable_domain.TypeString, IsInput: true},
				{Name: "min_replicas", Title: "Минимум реплик сервиса", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "db_replicas", Title: "Кол-во реплик БД", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "team_size", Title: "Размер команды испытаний", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "functional_tests_count", Title: "Кол-во функциональных тест-кейсов", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "regression_tests_count", Title: "Кол-во регрессионных тест-кейсов", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "load_target_rps", Title: "Целевой RPS", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "load_peak_rps", Title: "Пиковый RPS", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "target_p95_ms", Title: "Целевая p95, мс", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "target_p99_ms", Title: "Целевая p99, мс", Type: variable_domain.TypeInteger, IsInput: true},
			},
			Data: readTemplate("approbation_pmi.md"),
		},
	}
}

// approbationOPTemplate — Описание программы по ГОСТ 19.402-78.
func approbationOPTemplate() customTemplate {
	return customTemplate{
		Key:       "approbation_op",
		Name:      "Описание программы NotifyHub (ГОСТ 19.402-78)",
		CreatedAt: dateApprobationOPCreated,
		Version: &versionInput{
			CreatedAt: dateApprobationOPVersion,
			Variables: []variableInput{
				{Name: "program_name", Title: "Наименование программы", Type: variable_domain.TypeString, IsInput: true},
				{Name: "program_designation", Title: "Условное обозначение", Type: variable_domain.TypeString, IsInput: true},
				{Name: "program_version", Title: "Версия программы", Type: variable_domain.TypeString, IsInput: true},
				{Name: "document_date", Title: "Дата документа", Type: variable_domain.TypeString, IsInput: true},
				{Name: "document_code", Title: "Шифр документа", Type: variable_domain.TypeString, IsInput: true},
				{Name: "runtime", Title: "Среда выполнения", Type: variable_domain.TypeString, IsInput: true},
				{Name: "target_os", Title: "Целевая ОС", Type: variable_domain.TypeString, IsInput: true},
				{Name: "container_image", Title: "Тег контейнера", Type: variable_domain.TypeString, IsInput: true},
				{Name: "language", Title: "Язык программирования", Type: variable_domain.TypeString, IsInput: true},
				{Name: "source_repo", Title: "Репозиторий", Type: variable_domain.TypeString, IsInput: true},
				{Name: "cpu_arch", Title: "Архитектура CPU", Type: variable_domain.TypeString, IsInput: true},
				{Name: "min_cpu_per_node", Title: "Мин. vCPU на узел", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "min_memory_gb_per_node", Title: "Мин. RAM (ГБ) на узел", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "min_disk_gb_per_node", Title: "Мин. диск (ГБ) на узел", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "min_network_gbps", Title: "Мин. сеть, Гбит/с", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "k8s_version", Title: "Версия Kubernetes", Type: variable_domain.TypeString, IsInput: true},
				{Name: "helm_repo", Title: "Адрес Helm-репозитория", Type: variable_domain.TypeString, IsInput: true},
				{Name: "daily_events", Title: "Событий в сутки", Type: variable_domain.TypeString, IsInput: true},
				{Name: "daily_push", Title: "Push-сообщений в сутки", Type: variable_domain.TypeString, IsInput: true},
				{Name: "daily_sms", Title: "SMS в сутки", Type: variable_domain.TypeString, IsInput: true},
				{Name: "daily_email", Title: "Email-писем в сутки", Type: variable_domain.TypeString, IsInput: true},
			},
			Data: readTemplate("approbation_op.md"),
		},
	}
}
