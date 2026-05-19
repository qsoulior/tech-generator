package main

// passportTasks возвращает две задачи для шаблона «Технический паспорт»:
// одну успешную и одну с тремя нарушениями constraints (два — на sla_target,
// одно — на rps_target).
func passportTasks() []taskInput {
	successful := taskInput{
		Title:     "Паспорт NotifyHub 1.4.2 (успех)",
		CreatedAt: dateTaskSuccessCreated,
		UpdatedAt: dateTaskSuccessUpdated,
		Payload: map[string]string{
			"service_name":           "NotifyHub",
			"service_owner":          "Команда платформы коммуникаций",
			"service_version":        "1.4.2",
			"production_date":        "15.03.2026",
			"runtime":                "Go 1.25 (alpine)",
			"sla_target":             "0.999",
			"rps_target":             "850",
			"latency_p50_ms":         "12",
			"latency_p95_ms":         "45",
			"latency_p99_ms":         "120",
			"dep_db":                 "PostgreSQL 17 (notify_db)",
			"dep_queue":              "RabbitMQ (notify-exchange)",
			"dep_cache":              "Redis (sessions)",
			"instances":              "4",
			"cpu_per_instance":       "0.5",
			"memory_mb_per_instance": "512",
		},
	}

	failing := taskInput{
		Title:     "Паспорт NotifyHub black-friday (ошибка)",
		CreatedAt: dateTaskFailCreated,
		UpdatedAt: dateTaskFailUpdated,
		Payload: map[string]string{
			"service_name":           "NotifyHub",
			"service_owner":          "Команда платформы коммуникаций",
			"service_version":        "1.5.0-rc1",
			"production_date":        "28.11.2026",
			"runtime":                "Go 1.25 (alpine)",
			"sla_target":             "-0.5", // нарушит положительный_sla и production_grade_sla
			"rps_target":             "-100", // нарушит положительный_rps
			"latency_p50_ms":         "15",
			"latency_p95_ms":         "60",
			"latency_p99_ms":         "140",
			"dep_db":                 "PostgreSQL 17 (notify_db)",
			"dep_queue":              "RabbitMQ (notify-exchange)",
			"dep_cache":              "Redis (sessions)",
			"instances":              "6",
			"cpu_per_instance":       "0.75",
			"memory_mb_per_instance": "768",
		},
	}

	return []taskInput{successful, failing}
}

// approbationTZTask — успешная задача под ТЗ NotifyHub с полным payload.
func approbationTZTask() taskInput {
	return taskInput{
		Title:     "ТЗ NotifyHub — апробация",
		CreatedAt: dateApprobationTZTask,
		UpdatedAt: dateApprobationTZTaskDone,
		Payload: map[string]string{
			"system_name":      "Сервис уведомлений NotifyHub",
			"system_code":      "ИС-УВЕД-2026",
			"document_version": "1.0",
			"document_date":    "12.05.2026",
			"customer":         "Дирекция платформы коммуникаций",
			"customer_owner":   "Е. М. Сорокина, директор платформы",
			"developer":        "Команда платформы коммуникаций",
			"developer_owner":  "А. И. Зайцев, тимлид",
			"tz_basis":         "Приказ от 12.03.2026 № 187/АС о развитии транзакционного канала",
			"start_date":       "2026-04-01",
			"finish_date":      "2027-01-31",
			"system_purpose": "Автоматизация доставки транзакционных и маркетинговых уведомлений по каналам push, SMS и email " +
				"с поддержкой подписок, расчётом политик повторных попыток и выдачей отчётности по доставляемости.",
			"system_goals": "сократить среднее время доставки транзакционных сообщений до 1 секунды; " +
				"повысить SLA отправки до 99,95 %; " +
				"снизить ручные операции операторов клиентского сервиса минимум на 60 %; " +
				"унифицировать интеграцию каналов push/SMS/email под единым API",
			"users_endusers":   "до 18 000 000",
			"users_operators":  "≈ 250",
			"systems_sources":  "более 40",
			"sla_target":       "99,95 %",
			"max_downtime_min": "22",
			"subsystems": "Шлюз API (api-gateway); " +
				"Движок политик (policy-engine); " +
				"Диспетчер отправки (dispatcher); " +
				"Адаптер push-канала (adapter-push); " +
				"Адаптер SMS-канала (adapter-sms); " +
				"Адаптер email-канала (adapter-email); " +
				"Административная панель (admin-ui)",
			"runtime":        "Go 1.25 (alpine)",
			"k8s_version":    "1.31",
			"db_version":     "PostgreSQL 17",
			"min_replicas":   "4",
			"training_hours": "24",
		},
	}
}

func approbationPMITask() taskInput {
	return taskInput{
		Title:     "ПМИ NotifyHub — апробация",
		CreatedAt: dateApprobationPMITask,
		UpdatedAt: dateApprobationPMITaskDone,
		Payload: map[string]string{
			"program_name":           "NotifyHub",
			"program_version":        "1.4.2",
			"document_code":          "ПМИ-NOTIFY-2026-04",
			"test_date_start":        "2026-04-25",
			"test_date_end":          "2026-05-15",
			"test_location":          "Лаборатория QA, ЦОД М-9",
			"test_stand":             "stage-2.notify.internal",
			"responsible":            "К. А. Литвиненко, тест-лид",
			"approver":               "А. И. Зайцев, тимлид",
			"tz_reference":           "ТЗ-ИС-УВЕД-2026-v1.0",
			"min_replicas":           "4",
			"db_replicas":            "3",
			"team_size":              "6",
			"functional_tests_count": "184",
			"regression_tests_count": "412",
			"load_target_rps":        "850",
			"load_peak_rps":          "1500",
			"target_p95_ms":          "60",
			"target_p99_ms":          "150",
		},
	}
}

func approbationOPTask() taskInput {
	return taskInput{
		Title:     "Описание программы NotifyHub — апробация",
		CreatedAt: dateApprobationOPTask,
		UpdatedAt: dateApprobationOPTaskDone,
		Payload: map[string]string{
			"program_name":           "NotifyHub",
			"program_designation":    "ПО-NOTIFY-2026",
			"program_version":        "1.4.2",
			"document_date":          "10.05.2026",
			"document_code":          "ОП-NOTIFY-2026-05",
			"runtime":                "Go 1.25 (alpine)",
			"target_os":              "Linux (Debian 12, RHEL 9)",
			"container_image":        "registry.platform/notifyhub:1.4.2",
			"language":               "Go",
			"source_repo":            "git@gitlab.platform:notify/notifyhub.git",
			"cpu_arch":               "amd64 / arm64",
			"min_cpu_per_node":       "4",
			"min_memory_gb_per_node": "8",
			"min_disk_gb_per_node":   "40",
			"min_network_gbps":       "1",
			"k8s_version":            "1.31",
			"helm_repo":              "https://charts.platform.internal",
			"daily_events":           "23 000 000",
			"daily_push":             "16 400 000",
			"daily_sms":              "1 800 000",
			"daily_email":            "4 800 000",
		},
	}
}
