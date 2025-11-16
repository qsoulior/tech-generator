package service

import (
	"context"
	"errors"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	test_trm "github.com/qsoulior/tech-generator/backend/internal/pkg/test/trm"
	"github.com/qsoulior/tech-generator/backend/internal/service/version_create/domain"
)

func TestService_Handle_Success(t *testing.T) {
	ctx := context.Background()
	trCtx := context.WithValue(ctx, test_trm.TrKey{}, struct{}{})

	tests := []struct {
		name  string
		in    domain.VersionCreateIn
		setup func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository)
	}{
		{
			name: "VariablesConstraints",
			in: domain.VersionCreateIn{
				AuthorID:   1,
				TemplateID: 10,
				Data:       []byte{1, 2, 3},
				Variables: []domain.Variable{
					{
						Name:       "var_1",
						Type:       variable_domain.TypeString,
						Expression: lo.ToPtr("expr_1"),
						IsInput:    false,
						Constraints: []domain.Constraint{
							{Name: "constraint_1_1", Expression: "expr_1_1", IsActive: true},
							{Name: "constraint_1_2", Expression: "expr_1_2", IsActive: false},
						},
					},
					{
						Name:       "var_2",
						Type:       variable_domain.TypeFloat,
						Expression: lo.ToPtr("expr_2"),
						IsInput:    false,
						Constraints: []domain.Constraint{
							{Name: "constraint_2_1", Expression: "expr_2_1", IsActive: true},
						},
					},
				},
			},
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				templateVersion := domain.Version{
					TemplateID: 10,
					AuthorID:   1,
					Data:       []byte{1, 2, 3},
				}
				versionRepo.EXPECT().Create(trCtx, templateVersion).Return(int64(20), nil)

				variables := []domain.VariableToCreate{
					{VersionID: 20, Name: "var_1", Type: variable_domain.TypeString, Expression: lo.ToPtr("expr_1")},
					{VersionID: 20, Name: "var_2", Type: variable_domain.TypeFloat, Expression: lo.ToPtr("expr_2")},
				}
				variableRepo.EXPECT().Create(trCtx, variables).Return([]int64{31, 32}, nil)

				constraints := []domain.ConstraintToCreate{
					{VariableID: 31, Name: "constraint_1_1", Expression: "expr_1_1", IsActive: true},
					{VariableID: 31, Name: "constraint_1_2", Expression: "expr_1_2", IsActive: false},
					{VariableID: 32, Name: "constraint_2_1", Expression: "expr_2_1", IsActive: true},
				}
				constraintRepo.EXPECT().Create(trCtx, constraints).Return(nil)

				templateToUpdate := domain.TemplateToUpdate{ID: 10, LastVersionID: 20}
				templateRepo.EXPECT().UpdateByID(trCtx, templateToUpdate).Return(nil)
			},
		},
		{
			name: "NoConstraints",
			in: domain.VersionCreateIn{
				AuthorID:   1,
				TemplateID: 10,
				Data:       []byte{1, 2, 3},
				Variables: []domain.Variable{
					{
						Name:        "var_1",
						Type:        variable_domain.TypeString,
						Expression:  lo.ToPtr("expr_1"),
						Constraints: []domain.Constraint{},
					},
				},
			},
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				templateVersion := domain.Version{
					TemplateID: 10,
					AuthorID:   1,
					Data:       []byte{1, 2, 3},
				}
				versionRepo.EXPECT().Create(trCtx, templateVersion).Return(int64(20), nil)

				variables := []domain.VariableToCreate{
					{VersionID: 20, Name: "var_1", Type: variable_domain.TypeString, Expression: lo.ToPtr("expr_1")},
				}
				variableRepo.EXPECT().Create(trCtx, variables).Return([]int64{31}, nil)

				templateToUpdate := domain.TemplateToUpdate{ID: 10, LastVersionID: 20}
				templateRepo.EXPECT().UpdateByID(trCtx, templateToUpdate).Return(nil)
			},
		},
		{
			name: "NoVariables",
			in: domain.VersionCreateIn{
				AuthorID:   1,
				TemplateID: 10,
				Data:       []byte{1, 2, 3},
				Variables:  []domain.Variable{},
			},
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				templateVersion := domain.Version{
					TemplateID: 10,
					AuthorID:   1,
					Data:       []byte{1, 2, 3},
				}
				versionRepo.EXPECT().Create(trCtx, templateVersion).Return(int64(20), nil)

				templateToUpdate := domain.TemplateToUpdate{ID: 10, LastVersionID: 20}
				templateRepo.EXPECT().UpdateByID(trCtx, templateToUpdate).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			versionRepo := NewMockversionRepository(ctrl)
			variableRepo := NewMockvariableRepository(ctrl)
			constraintRepo := NewMockconstraintRepository(ctrl)
			trManager := test_trm.New()

			tt.setup(templateRepo, versionRepo, variableRepo, constraintRepo)

			usecase := New(templateRepo, versionRepo, variableRepo, constraintRepo, trManager)

			err := usecase.Handle(ctx, tt.in)
			require.NoError(t, err)
		})
	}
}

func TestService_Handle_Error(t *testing.T) {
	ctx := context.Background()
	trCtx := context.WithValue(ctx, test_trm.TrKey{}, struct{}{})

	validIn := domain.VersionCreateIn{
		AuthorID:   1,
		TemplateID: 10,
		Data:       []byte{1, 2, 3},
		Variables: []domain.Variable{
			{
				Name:       "var_1",
				Type:       variable_domain.TypeString,
				Expression: lo.ToPtr("expr_1"),
				Constraints: []domain.Constraint{
					{
						Name:       "constraint_1_1",
						Expression: "expr_1_1",
						IsActive:   true,
					},
				},
			},
		},
	}

	tests := []struct {
		name  string
		in    domain.VersionCreateIn
		setup func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository)
		want  string
	}{
		{
			name: "in_Validate",
			in: domain.VersionCreateIn{
				AuthorID:   1,
				TemplateID: 10,
				Data:       []byte{1, 2, 3},
				Variables: []domain.Variable{
					{
						Name:        "var",
						Type:        "invalid",
						Expression:  lo.ToPtr("expr"),
						Constraints: []domain.Constraint{},
					},
				},
			},
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
			},
			want: domain.ErrValueInvalid.Error(),
		},
		{
			name: "versionRepo_Create",
			in:   validIn,
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				versionRepo.EXPECT().Create(trCtx, gomock.Any()).Return(int64(0), errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "variableRepo_Create",
			in:   validIn,
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				versionRepo.EXPECT().Create(trCtx, gomock.Any()).Return(int64(20), nil)
				variableRepo.EXPECT().Create(trCtx, gomock.Any()).Return(nil, errors.New("test2"))
			},
			want: "test2",
		},
		{
			name: "domain_ErrVariableIDsInvalid",
			in:   validIn,
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				versionRepo.EXPECT().Create(trCtx, gomock.Any()).Return(int64(20), nil)
				variableRepo.EXPECT().Create(trCtx, gomock.Any()).Return([]int64{}, nil)
			},
			want: domain.ErrVariableIDsInvalid.Error(),
		},
		{
			name: "constraintRepo_Create",
			in:   validIn,
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				versionRepo.EXPECT().Create(trCtx, gomock.Any()).Return(int64(20), nil)
				variableRepo.EXPECT().Create(trCtx, gomock.Any()).Return([]int64{31}, nil)
				constraintRepo.EXPECT().Create(trCtx, gomock.Any()).Return(errors.New("test3"))
			},
			want: "test3",
		},
		{
			name: "templateRepo_UpdateByID",
			in:   validIn,
			setup: func(templateRepo *MocktemplateRepository, versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				versionRepo.EXPECT().Create(trCtx, gomock.Any()).Return(int64(20), nil)
				variableRepo.EXPECT().Create(trCtx, gomock.Any()).Return([]int64{31}, nil)
				constraintRepo.EXPECT().Create(trCtx, gomock.Any()).Return(nil)
				templateRepo.EXPECT().UpdateByID(trCtx, gomock.Any()).Return(errors.New("test4"))
			},
			want: "test4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			versionRepo := NewMockversionRepository(ctrl)
			variableRepo := NewMockvariableRepository(ctrl)
			constraintRepo := NewMockconstraintRepository(ctrl)
			trManager := test_trm.New()

			tt.setup(templateRepo, versionRepo, variableRepo, constraintRepo)

			usecase := New(templateRepo, versionRepo, variableRepo, constraintRepo, trManager)

			err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
