package data_process_service

import (
	"bytes"
	"context"
	"text/template"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Handle(ctx context.Context, in domain.DataProcessIn) ([]byte, error) {
	tmpl, err := template.New("").Parse(string(in.Data))
	if err != nil {
		return nil, &domain.ProcessError{Message: domain.MessageTemplateParse}
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, in.Values)
	if err != nil {
		return nil, &domain.ProcessError{Message: domain.MessageTemplateExec}
	}

	return buf.Bytes(), nil
}
