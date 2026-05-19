package data_process_service

import (
	"bytes"
	"context"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

// templateFuncs is the sprig text/template helper set with process-environment
// accessors removed so a template cannot exfiltrate the worker's secrets.
var templateFuncs = func() template.FuncMap {
	funcs := sprig.TxtFuncMap()
	delete(funcs, "env")
	delete(funcs, "expandenv")
	delete(funcs, "getHostByName")
	return funcs
}()

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Handle(ctx context.Context, in domain.DataProcessIn) ([]byte, error) {
	tmpl, err := template.New("").Funcs(templateFuncs).Parse(string(in.Data))
	if err != nil {
		return nil, &task_domain.ProcessError{
			Message:  task_domain.MessageTemplateParse,
			Template: buildTemplateError(in.Data, err),
		}
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, in.Values)
	if err != nil {
		return nil, &task_domain.ProcessError{
			Message:  task_domain.MessageTemplateExec,
			Template: buildTemplateError(in.Data, err),
		}
	}

	return buf.Bytes(), nil
}

// templateErrRe matches the canonical Go text/template diagnostic prefix:
// "template: <name>:<line>[:<col>]: <message>". The name segment is optional
// content up to the first colon, line/col are decimal digits.
var templateErrRe = regexp.MustCompile(`^template:\s*[^:]*:(\d+)(?::(\d+))?:\s*(.+)$`)

// buildTemplateError extracts line/column/snippet from a Go template parse or
// execution error. When the error message does not match the expected format,
// returns nil so the caller falls back to the high-level message only.
func buildTemplateError(data []byte, err error) *task_domain.TemplateError {
	msg := err.Error()
	m := templateErrRe.FindStringSubmatch(msg)
	if m == nil {
		return nil
	}

	line, convErr := strconv.Atoi(m[1])
	if convErr != nil || line < 1 {
		return nil
	}

	col := 0
	if m[2] != "" {
		if c, e := strconv.Atoi(m[2]); e == nil && c > 0 {
			col = c
		}
	}

	return &task_domain.TemplateError{
		Line:    line,
		Column:  col,
		Snippet: extractLine(data, line),
		Detail:  strings.TrimSpace(m[3]),
	}
}

func extractLine(data []byte, line int) string {
	if line < 1 {
		return ""
	}
	lines := strings.SplitN(string(data), "\n", line+1)
	if len(lines) < line {
		return ""
	}
	return strings.TrimRight(lines[line-1], "\r")
}
