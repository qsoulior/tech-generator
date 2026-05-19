package main

import (
	"time"

	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
)

type defaultTemplate struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
	Version   *versionInput
}

// defaultTemplates возвращает библиотеку стандартных ГОСТ-шаблонов, видимую
// всем пользователям через /template/default/list и клонируемую через
// /template/create_from_default.
func defaultTemplates() []defaultTemplate {
	return []defaultTemplate{
		gostTZAS(),
		gostTZProgram(),
		gostPMI(),
		gostOperatorManual(),
	}
}

// gostTZAS — Техническое задание на создание АС (ГОСТ 34.602-2020).
func gostTZAS() defaultTemplate {
	return defaultTemplate{
		Name:      "Техническое задание на создание АС (ГОСТ 34.602-2020)",
		CreatedAt: dateGOSTASCreated,
		UpdatedAt: dateGOSTASUpdated,
		Version: &versionInput{
			CreatedAt: dateGOSTASCreated,
			Variables: []variableInput{
				{Name: "system_name", Title: "Полное наименование системы", Type: variable_domain.TypeString, IsInput: true},
				{Name: "system_code", Title: "Шифр темы", Type: variable_domain.TypeString, IsInput: true},
				{Name: "customer", Title: "Заказчик", Type: variable_domain.TypeString, IsInput: true},
				{Name: "developer", Title: "Исполнитель (разработчик)", Type: variable_domain.TypeString, IsInput: true},
				{Name: "document_date", Title: "Дата документа", Type: variable_domain.TypeString, IsInput: true},
				{Name: "tz_basis", Title: "Основание для разработки", Type: variable_domain.TypeString, IsInput: true},
				{Name: "purpose", Title: "Назначение системы", Type: variable_domain.TypeString, IsInput: true},
				{Name: "goals", Title: "Цели создания (через `;`)", Type: variable_domain.TypeString, IsInput: true},
				{Name: "completion_date", Title: "Срок ввода в действие", Type: variable_domain.TypeString, IsInput: true},
			},
			Data: readTemplate("gost_tz_as.md"),
		},
	}
}

// gostTZProgram — ТЗ на программу или программное изделие (ГОСТ 19.201-78).
func gostTZProgram() defaultTemplate {
	return defaultTemplate{
		Name:      "Техническое задание на программу (ГОСТ 19.201-78)",
		CreatedAt: dateGOSTProgramCreated,
		Version: &versionInput{
			CreatedAt: dateGOSTProgramCreated,
			Variables: []variableInput{
				{Name: "program_name", Title: "Наименование программы", Type: variable_domain.TypeString, IsInput: true},
				{Name: "program_designation", Title: "Условное обозначение", Type: variable_domain.TypeString, IsInput: true},
				{Name: "application_area", Title: "Область применения", Type: variable_domain.TypeString, IsInput: true},
				{Name: "basis_doc", Title: "Документ-основание для разработки", Type: variable_domain.TypeString, IsInput: true},
				{Name: "purpose", Title: "Назначение разработки", Type: variable_domain.TypeString, IsInput: true},
				{Name: "user_qualification", Title: "Квалификация пользователей", Type: variable_domain.TypeString, IsInput: true},
				{Name: "stage_count", Title: "Кол-во этапов разработки", Type: variable_domain.TypeInteger, IsInput: true},
			},
			Data: readTemplate("gost_tz_program.md"),
		},
	}
}

// gostPMI — Программа и методика испытаний (ГОСТ 19.301-79).
func gostPMI() defaultTemplate {
	return defaultTemplate{
		Name:      "Программа и методика испытаний (ГОСТ 19.301-79)",
		CreatedAt: dateGOSTPMICreated,
		UpdatedAt: dateGOSTPMIUpdated,
		Version: &versionInput{
			CreatedAt: dateGOSTPMICreated,
			Variables: []variableInput{
				{Name: "program_name", Title: "Наименование программы (объекта испытаний)", Type: variable_domain.TypeString, IsInput: true},
				{Name: "version", Title: "Версия объекта испытаний", Type: variable_domain.TypeString, IsInput: true},
				{Name: "test_goal", Title: "Цель испытаний", Type: variable_domain.TypeString, IsInput: true},
				{Name: "test_location", Title: "Место проведения испытаний", Type: variable_domain.TypeString, IsInput: true},
				{Name: "test_date_start", Title: "Дата начала", Type: variable_domain.TypeString, IsInput: true},
				{Name: "test_date_end", Title: "Дата окончания", Type: variable_domain.TypeString, IsInput: true},
				{Name: "responsible", Title: "Ответственный", Type: variable_domain.TypeString, IsInput: true},
			},
			Data: readTemplate("gost_pmi.md"),
		},
	}
}

// gostOperatorManual — Руководство оператора (ГОСТ 19.505-79).
func gostOperatorManual() defaultTemplate {
	return defaultTemplate{
		Name:      "Руководство оператора (ГОСТ 19.505-79)",
		CreatedAt: dateGOSTOperatorCreated,
		Version: &versionInput{
			CreatedAt: dateGOSTOperatorCreated,
			Variables: []variableInput{
				{Name: "program_name", Title: "Наименование программы", Type: variable_domain.TypeString, IsInput: true},
				{Name: "program_version", Title: "Версия программы", Type: variable_domain.TypeString, IsInput: true},
				{Name: "min_memory_mb", Title: "Мин. объём ОЗУ (МБ)", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "min_disk_mb", Title: "Мин. место на диске (МБ)", Type: variable_domain.TypeInteger, IsInput: true},
				{Name: "executable", Title: "Имя исполняемого модуля", Type: variable_domain.TypeString, IsInput: true},
				{Name: "operator_role", Title: "Роль оператора", Type: variable_domain.TypeString, IsInput: true},
			},
			Data: readTemplate("gost_operator_manual.md"),
		},
	}
}
