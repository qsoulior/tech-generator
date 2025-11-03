package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	test_trm "github.com/qsoulior/tech-generator/backend/internal/pkg/test/trm"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_version_create/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()
	trCtx := context.WithValue(ctx, test_trm.TrKey{}, struct{}{})

	tests := []struct {
		name  string
		in    domain.TemplateVersionCreateIn
		setup func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository)
	}{
		{
			name: "IsAuthor",
			in: domain.TemplateVersionCreateIn{
				AuthorID:   1,
				TemplateID: 10,
				Data:       []byte{1, 2, 3},
				Variables: []domain.Variable{
					{
						Name:       "var_1",
						Type:       variable_domain.TypeString,
						Expression: "expr_1",
						Constraints: []domain.VariableConstraint{
							{Name: "constraint_1_1", Expression: "expr_1_1", IsActive: true},
							{Name: "constraint_1_2", Expression: "expr_1_2", IsActive: false},
						},
					},
					{
						Name:       "var_2",
						Type:       variable_domain.TypeFloat,
						Expression: "expr_2",
						Constraints: []domain.VariableConstraint{
							{Name: "constraint_2_1", Expression: "expr_2_1", IsActive: true},
						},
					},
				},
			},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				templateVersion := domain.TemplateVersion{
					TemplateID: 10,
					AuthorID:   1,
					Data:       []byte{1, 2, 3},
				}
				templateVersionRepo.EXPECT().Create(trCtx, templateVersion).Return(int64(20), nil)

				variables := []domain.VariableToCreate{
					{VersionID: 20, Name: "var_1", Type: variable_domain.TypeString, Expression: "expr_1"},
					{VersionID: 20, Name: "var_2", Type: variable_domain.TypeFloat, Expression: "expr_2"},
				}
				variableRepo.EXPECT().Create(trCtx, variables).Return([]int64{31, 32}, nil)

				constraints := []domain.VariableConstraintToCreate{
					{VariableID: 31, Name: "constraint_1_1", Expression: "expr_1_1", IsActive: true},
					{VariableID: 31, Name: "constraint_1_2", Expression: "expr_1_2", IsActive: false},
					{VariableID: 32, Name: "constraint_2_1", Expression: "expr_2_1", IsActive: true},
				}
				variableConstraintRepo.EXPECT().Create(trCtx, constraints).Return(nil)

				templateToUpdate := domain.TemplateToUpdate{ID: 10, LastVersionID: 20}
				templateRepo.EXPECT().UpdateByID(trCtx, templateToUpdate).Return(nil)
			},
		},
		{
			name: "IsRootAuthor/NoConstraints",
			in: domain.TemplateVersionCreateIn{
				AuthorID:   1,
				TemplateID: 10,
				Data:       []byte{1, 2, 3},
				Variables: []domain.Variable{
					{
						Name:        "var_1",
						Type:        variable_domain.TypeString,
						Expression:  "expr_1",
						Constraints: []domain.VariableConstraint{},
					},
				},
			},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{AuthorID: 2, ProjectAuthorID: 1}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				templateVersion := domain.TemplateVersion{
					TemplateID: 10,
					AuthorID:   1,
					Data:       []byte{1, 2, 3},
				}
				templateVersionRepo.EXPECT().Create(trCtx, templateVersion).Return(int64(20), nil)

				variables := []domain.VariableToCreate{
					{VersionID: 20, Name: "var_1", Type: variable_domain.TypeString, Expression: "expr_1"},
				}
				variableRepo.EXPECT().Create(trCtx, variables).Return([]int64{31}, nil)

				templateToUpdate := domain.TemplateToUpdate{ID: 10, LastVersionID: 20}
				templateRepo.EXPECT().UpdateByID(trCtx, templateToUpdate).Return(nil)
			},
		},
		{
			name: "IsWriter/NoVariables",
			in: domain.TemplateVersionCreateIn{
				AuthorID:   1,
				TemplateID: 10,
				Data:       []byte{1, 2, 3},
				Variables:  []domain.Variable{},
			},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{
					AuthorID:        2,
					ProjectAuthorID: 3,
					Users:           []domain.TemplateUser{{ID: 1, Role: user_domain.RoleWrite}},
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				templateVersion := domain.TemplateVersion{
					TemplateID: 10,
					AuthorID:   1,
					Data:       []byte{1, 2, 3},
				}
				templateVersionRepo.EXPECT().Create(trCtx, templateVersion).Return(int64(20), nil)

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
			templateVersionRepo := NewMocktemplateVersionRepository(ctrl)
			variableRepo := NewMockvariableRepository(ctrl)
			variableConstraintRepo := NewMockvariableConstraintRepository(ctrl)
			trManager := test_trm.New()

			tt.setup(templateRepo, templateVersionRepo, variableRepo, variableConstraintRepo)

			usecase := New(templateRepo, templateVersionRepo, variableRepo, variableConstraintRepo, trManager)

			err := usecase.Handle(ctx, tt.in)
			require.NoError(t, err)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()
	trCtx := context.WithValue(ctx, test_trm.TrKey{}, struct{}{})

	validIn := domain.TemplateVersionCreateIn{
		AuthorID:   1,
		TemplateID: 10,
		Data:       []byte{1, 2, 3},
		Variables: []domain.Variable{
			{
				Name:       "var_1",
				Type:       variable_domain.TypeString,
				Expression: "expr_1",
				Constraints: []domain.VariableConstraint{
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
		in    domain.TemplateVersionCreateIn
		setup func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository)
		want  string
	}{
		{
			name: "in_Validate",
			in: domain.TemplateVersionCreateIn{
				AuthorID:   1,
				TemplateID: 10,
				Data:       []byte{1, 2, 3},
				Variables: []domain.Variable{
					{
						Name:        "var",
						Type:        "invalid",
						Expression:  "expr",
						Constraints: []domain.VariableConstraint{},
					},
				},
			},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
			},
			want: domain.ErrValueInvalid.Error(),
		},
		{
			name: "templateRepo_GetByID",
			in:   validIn,
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrTemplateNotFound",
			in:   validIn,
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, nil)
			},
			want: domain.ErrTemplateNotFound.Error(),
		},
		{
			name: "domain_ErrTemplateInvalid",
			in:   validIn,
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{AuthorID: 2, ProjectAuthorID: 3}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.ErrTemplateInvalid.Error(),
		},
		{
			name: "templateVersionRepo_Create",
			in:   validIn,
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateVersionRepo.EXPECT().Create(trCtx, gomock.Any()).Return(int64(0), errors.New("test2"))
			},
			want: "test2",
		},
		{
			name: "variableRepo_Create",
			in:   validIn,
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateVersionRepo.EXPECT().Create(trCtx, gomock.Any()).Return(int64(20), nil)
				variableRepo.EXPECT().Create(trCtx, gomock.Any()).Return(nil, errors.New("test3"))
			},
			want: "test3",
		},
		{
			name: "domain_ErrVariableIDsInvalid",
			in:   validIn,
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateVersionRepo.EXPECT().Create(trCtx, gomock.Any()).Return(int64(20), nil)
				variableRepo.EXPECT().Create(trCtx, gomock.Any()).Return([]int64{}, nil)
			},
			want: domain.ErrVariableIDsInvalid.Error(),
		},
		{
			name: "variableConstraintRepo_Create",
			in:   validIn,
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateVersionRepo.EXPECT().Create(trCtx, gomock.Any()).Return(int64(20), nil)
				variableRepo.EXPECT().Create(trCtx, gomock.Any()).Return([]int64{31}, nil)
				variableConstraintRepo.EXPECT().Create(trCtx, gomock.Any()).Return(errors.New("test4"))
			},
			want: "test4",
		},
		{
			name: "templateRepo_UpdateByID",
			in:   validIn,
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{AuthorID: 1, ProjectAuthorID: 2}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateVersionRepo.EXPECT().Create(trCtx, gomock.Any()).Return(int64(20), nil)
				variableRepo.EXPECT().Create(trCtx, gomock.Any()).Return([]int64{31}, nil)
				variableConstraintRepo.EXPECT().Create(trCtx, gomock.Any()).Return(nil)
				templateRepo.EXPECT().UpdateByID(trCtx, gomock.Any()).Return(errors.New("test5"))
			},
			want: "test5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			templateVersionRepo := NewMocktemplateVersionRepository(ctrl)
			variableRepo := NewMockvariableRepository(ctrl)
			variableConstraintRepo := NewMockvariableConstraintRepository(ctrl)
			trManager := test_trm.New()

			tt.setup(templateRepo, templateVersionRepo, variableRepo, variableConstraintRepo)

			usecase := New(templateRepo, templateVersionRepo, variableRepo, variableConstraintRepo, trManager)

			err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
