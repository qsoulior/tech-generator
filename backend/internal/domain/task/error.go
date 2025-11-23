package task_domain

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

var (
	MessageCycle         = "Обнаружена цикличная зависимость"
	MessageCompile       = "Ошибка компиляции"
	MessageRun           = "Ошибка выполнения"
	MessageCheck         = "Нарушение ограничения"
	MessageTemplateParse = "Ошибка парсинга шаблона"
	MessageTemplateExec  = "Ошибка выполнения шаблона"
)

type ProcessError struct {
	Message        string          `json:"message,omitempty"`
	VariableErrors []VariableError `json:"variable_errors,omitempty"`
}

func (e *ProcessError) Error() string {
	variableMessages := lo.Map(e.VariableErrors, func(ve VariableError, _ int) string { return ve.Error() })
	return fmt.Sprintf("%s\n%s", e.Message, strings.Join(variableMessages, "\n"))
}

type VariableError struct {
	ID               int64             `json:"id"`
	Name             string            `json:"name"`
	Message          string            `json:"message,omitempty"`
	ConstraintErrors []ConstraintError `json:"constraint_errors,omitempty"`
}

func (e *VariableError) Error() string {
	constraintMessages := lo.Map(e.ConstraintErrors, func(ce ConstraintError, _ int) string { return ce.Error() })
	return fmt.Sprintf("[%d] %s: %s\n%s", e.ID, e.Name, e.Message, strings.Join(constraintMessages, "\n"))
}

type ConstraintError struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message,omitempty"`
}

func (e *ConstraintError) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.ID, e.Name, e.Message)
}
