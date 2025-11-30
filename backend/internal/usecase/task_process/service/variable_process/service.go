package variable_process_service

import (
	"context"
	"strconv"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/samber/lo"

	task_domain "github.com/qsoulior/tech-generator/backend/internal/domain/task"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Handle(ctx context.Context, in domain.VariableProcessIn) (map[string]any, error) {
	variablesByName := lo.KeyBy(in.Variables, func(v domain.Variable) string { return v.Name })

	dependencies := buildDependencies(variablesByName)

	variableNames, err := sortDependencies(dependencies)
	if err != nil {
		return nil, err
	}

	var variableErrors []task_domain.VariableError
	variableValues := make(map[string]any)

	for name, value := range in.Payload {
		variableValue, variableError := parseVariable(name, value, variablesByName)
		if variableError != nil {
			variableErrors = append(variableErrors, *variableError)
		}

		variableValues[name] = variableValue
	}

	for _, name := range variableNames {
		variable := variablesByName[name]

		variableValue, variableError := processVariable(variable, variableValues)
		if variableError != nil {
			variableErrors = append(variableErrors, *variableError)
		}

		variableValues[name] = variableValue
	}

	if len(variableErrors) != 0 {
		return nil, &task_domain.ProcessError{VariableErrors: variableErrors}
	}

	return variableValues, nil
}

func buildDependencies(variablesByName map[string]domain.Variable) map[string][]string {
	names := lo.Keys(variablesByName)
	dependents := make(map[string][]string)

	for _, v := range variablesByName {
		dependencies := extractDependencies(lo.FromPtr(v.Expression), names)
		dependents[v.Name] = dependencies
	}

	return dependents
}

func extractDependencies(expr string, names []string) []string {
	return lo.Filter(names, func(name string, _ int) bool {
		return strings.Contains(expr, name)
	})
}

func sortDependencies(dependencies map[string][]string) ([]string, error) {
	var dfs func(string) error

	sorted := make([]string, 0, len(dependencies))

	gray := make(map[string]struct{})
	black := make(map[string]struct{})

	dfs = func(dependent string) error {
		if _, ok := gray[dependent]; ok {
			return &task_domain.ProcessError{Message: task_domain.MessageCycle}
		}

		gray[dependent] = struct{}{}
		for _, dependency := range dependencies[dependent] {
			if _, ok := black[dependency]; ok {
				continue
			}

			if err := dfs(dependency); err != nil {
				return err
			}
		}

		delete(gray, dependent)
		black[dependent] = struct{}{}
		sorted = append(sorted, dependent)
		return nil
	}

	for dependency := range dependencies {
		if _, ok := black[dependency]; ok {
			continue
		}

		if err := dfs(dependency); err != nil {
			return []string{}, err
		}
	}

	return sorted, nil
}

var parsers = map[variable_domain.Type]func(s string) (any, error){
	variable_domain.TypeInteger: func(s string) (any, error) { return strconv.ParseInt(s, 10, 64) },
	variable_domain.TypeFloat:   func(s string) (any, error) { return strconv.ParseFloat(s, 64) },
	variable_domain.TypeString:  func(s string) (any, error) { return s, nil },
}

func parseVariable(name, value string, variablesByName map[string]domain.Variable) (any, *task_domain.VariableError) {
	variable, ok := variablesByName[name]
	if !ok {
		return nil, &task_domain.VariableError{ID: variable.ID, Name: variable.Name, Message: task_domain.MessageVariableNotFound}
	}

	parser, ok := parsers[variable.Type]
	if !ok {
		return nil, &task_domain.VariableError{ID: variable.ID, Name: variable.Name, Message: task_domain.MessageVariableTypeUnknown}
	}

	parsedValue, err := parser(value)
	if err != nil {
		return nil, &task_domain.VariableError{ID: variable.ID, Name: variable.Name, Message: task_domain.MessageVariableParse}
	}

	return parsedValue, nil
}

func processVariable(variable domain.Variable, values map[string]any) (any, *task_domain.VariableError) {
	value, variableError := processVariableValue(variable, values)
	if variableError != nil {
		return nil, variableError
	}

	constraintErrors := processConstraints(variable.Name, value, variable.Constraints)
	if len(constraintErrors) != 0 {
		return nil, &task_domain.VariableError{ID: variable.ID, Name: variable.Name, ConstraintErrors: constraintErrors}
	}

	return value, nil
}

func processVariableValue(variable domain.Variable, values map[string]any) (any, *task_domain.VariableError) {
	if variable.IsInput {
		return values[variable.Name], nil
	}

	program, err := expr.Compile(lo.FromPtr(variable.Expression), expr.Env(values))
	if err != nil {
		return nil, &task_domain.VariableError{ID: variable.ID, Name: variable.Name, Message: task_domain.MessageVariableCompile}
	}

	value, err := expr.Run(program, values)
	if err != nil {
		return nil, &task_domain.VariableError{ID: variable.ID, Name: variable.Name, Message: task_domain.MessageVariableExec}
	}

	return value, nil
}

func processConstraints(name string, value any, constraints []domain.Constraint) []task_domain.ConstraintError {
	var constraintErrors []task_domain.ConstraintError

	for _, constraint := range constraints {
		if !constraint.IsActive {
			continue
		}

		constraintError := processConstraint(name, value, constraint)
		if constraintError != nil {
			constraintErrors = append(constraintErrors, *constraintError)
		}
	}

	return constraintErrors
}

func processConstraint(name string, value any, constraint domain.Constraint) *task_domain.ConstraintError {
	env := map[string]any{name: value}

	program, err := expr.Compile(constraint.Expression, expr.Env(env), expr.AsBool())
	if err != nil {
		return &task_domain.ConstraintError{ID: constraint.ID, Name: constraint.Name, Message: task_domain.MessageVariableCompile}
	}

	check, err := expr.Run(program, env)
	if err != nil {
		return &task_domain.ConstraintError{ID: constraint.ID, Name: constraint.Name, Message: task_domain.MessageVariableExec}
	}

	if !check.(bool) {
		return &task_domain.ConstraintError{ID: constraint.ID, Name: constraint.Name, Message: task_domain.MessageConstraintCheck}
	}

	return nil
}
