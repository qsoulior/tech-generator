package data_process_service

import (
	"bytes"
	"context"
	"text/template"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Handle(ctx context.Context, in domain.DataProcessIn) ([]byte, error) {
	tmpl, err := template.New("").Parse(string(in.Data))
	if err != nil {
		return nil, &task_domain.ProcessError{Message: task_domain.MessageTemplateParse}
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, in.Values)
	if err != nil {
		return nil, &task_domain.ProcessError{Message: task_domain.MessageTemplateExec}
	}

	return buf.Bytes(), nil
}
