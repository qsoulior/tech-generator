package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_get_by_id/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()
	createdAt := gofakeit.Date()

	tests := []struct {
		name  string
		in    domain.TemplateGetByIDIn
		setup func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository)
		want  domain.TemplateGetByIDOut
	}{
		{
			name: "IsProjectAuthor",
			in:   domain.TemplateGetByIDIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{
					LastVersionID:   lo.ToPtr[int64](20),
					AuthorID:        2,
					ProjectAuthorID: 1,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				version := domain.TemplateVersion{
					ID:        20,
					Number:    2,
					CreatedAt: createdAt,
					Data:      []byte{1, 2, 3},
				}
				templateVersionRepo.EXPECT().GetByID(ctx, int64(20)).Return(&version, nil)

				variables := []domain.Variable{
					{
						ID:         31,
						Name:       "var1",
						Type:       variable_domain.TypeString,
						Expression: "expr31",
					},
					{
						ID:         32,
						Name:       "var2",
						Type:       variable_domain.TypeInteger,
						Expression: "expr32",
					},
				}
				variableRepo.EXPECT().ListByVersionID(ctx, int64(20)).Return(variables, nil)

				constraints := []domain.VariableConstraint{
					{
						ID:         41,
						VariableID: 31,
						Name:       "constraint1",
						Expression: "expr41",
						IsActive:   false,
					},
					{
						ID:         42,
						VariableID: 31,
						Name:       "constraint2",
						Expression: "expr42",
						IsActive:   true,
					},
					{
						ID:         43,
						VariableID: 32,
						Name:       "constraint3",
						Expression: "expr43",
						IsActive:   true,
					},
				}
				variableConstraintRepo.EXPECT().ListByVariableIDs(ctx, []int64{31, 32}).Return(constraints, nil)
			},
			want: domain.TemplateGetByIDOut{
				VersionID:     20,
				VersionNumber: 2,
				CreatedAt:     createdAt,
				Data:          []byte{1, 2, 3},
				Variables: []domain.Variable{
					{
						ID:         31,
						Name:       "var1",
						Type:       variable_domain.TypeString,
						Expression: "expr31",
						Constraints: []domain.VariableConstraint{
							{
								ID:         41,
								VariableID: 31,
								Name:       "constraint1",
								Expression: "expr41",
								IsActive:   false,
							},
							{
								ID:         42,
								VariableID: 31,
								Name:       "constraint2",
								Expression: "expr42",
								IsActive:   true,
							},
						},
					},
					{
						ID:         32,
						Name:       "var2",
						Type:       variable_domain.TypeInteger,
						Expression: "expr32",
						Constraints: []domain.VariableConstraint{
							{
								ID:         43,
								VariableID: 32,
								Name:       "constraint3",
								Expression: "expr43",
								IsActive:   true,
							},
						},
					},
				},
			},
		},
		{
			name: "IsAuthor/NoVariables",
			in:   domain.TemplateGetByIDIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{
					LastVersionID:   lo.ToPtr[int64](20),
					AuthorID:        1,
					ProjectAuthorID: 2,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				version := domain.TemplateVersion{
					ID:        20,
					Number:    2,
					CreatedAt: createdAt,
					Data:      []byte{1, 2, 3},
				}
				templateVersionRepo.EXPECT().GetByID(ctx, int64(20)).Return(&version, nil)

				variables := []domain.Variable{}
				variableRepo.EXPECT().ListByVersionID(ctx, int64(20)).Return(variables, nil)
			},
			want: domain.TemplateGetByIDOut{
				VersionID:     20,
				VersionNumber: 2,
				CreatedAt:     createdAt,
				Data:          []byte{1, 2, 3},
				Variables:     []domain.Variable{},
			},
		},
		{
			name: "IsReader/NoLastVersion",
			in:   domain.TemplateGetByIDIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{
					LastVersionID:   nil,
					AuthorID:        2,
					ProjectAuthorID: 3,
					Users:           []domain.TemplateUser{{ID: 1, Role: user_domain.RoleRead}},
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.TemplateGetByIDOut{VersionID: 0, VersionNumber: 0},
		},
		{
			name: "IsWriter/NoLastVersion",
			in:   domain.TemplateGetByIDIn{TemplateID: 10, UserID: 1},
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{
					LastVersionID:   nil,
					AuthorID:        2,
					ProjectAuthorID: 3,
					Users:           []domain.TemplateUser{{ID: 1, Role: user_domain.RoleWrite}},
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.TemplateGetByIDOut{VersionID: 0, VersionNumber: 0},
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

			tt.setup(templateRepo, templateVersionRepo, variableRepo, variableConstraintRepo)

			usecase := New(templateRepo, templateVersionRepo, variableRepo, variableConstraintRepo)

			got, err := usecase.Handle(ctx, tt.in)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.want, *got)
		})
	}
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()
	in := domain.TemplateGetByIDIn{TemplateID: 10, UserID: 1}

	tests := []struct {
		name  string
		setup func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository)
		want  string
	}{
		{
			name: "templateRepo_GetByID",
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, errors.New("test1"))
			},
			want: "test1",
		},
		{
			name: "domain_ErrTemplateNotFound",
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(nil, nil)
			},
			want: domain.ErrTemplateNotFound.Error(),
		},
		{
			name: "domain_ErrTemplateInvalid",
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{
					LastVersionID:   lo.ToPtr[int64](20),
					AuthorID:        2,
					ProjectAuthorID: 3,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
			},
			want: domain.ErrTemplateInvalid.Error(),
		},
		{
			name: "templateVersionRepo_GetByID",
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{
					LastVersionID:   lo.ToPtr[int64](20),
					AuthorID:        2,
					ProjectAuthorID: 1,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateVersionRepo.EXPECT().GetByID(ctx, int64(20)).Return(nil, errors.New("test3"))
			},
			want: "test3",
		},
		{
			name: "domain_ErrTemplateVersionNotFound",
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{
					LastVersionID:   lo.ToPtr[int64](20),
					AuthorID:        2,
					ProjectAuthorID: 1,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)
				templateVersionRepo.EXPECT().GetByID(ctx, int64(20)).Return(nil, nil)
			},
			want: domain.ErrTemplateVersionNotFound.Error(),
		},
		{
			name: "variableRepo_ListByVersionID",
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{
					LastVersionID:   lo.ToPtr[int64](20),
					AuthorID:        2,
					ProjectAuthorID: 1,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				version := domain.TemplateVersion{ID: 20, Data: []byte{1, 2, 3}}
				templateVersionRepo.EXPECT().GetByID(ctx, int64(20)).Return(&version, nil)
				variableRepo.EXPECT().ListByVersionID(ctx, int64(20)).Return(nil, errors.New("test4"))
			},
			want: "test4",
		},
		{
			name: "variableConstraintRepo_ListByVariableIDs",
			setup: func(templateRepo *MocktemplateRepository, templateVersionRepo *MocktemplateVersionRepository, variableRepo *MockvariableRepository, variableConstraintRepo *MockvariableConstraintRepository) {
				template := domain.Template{
					LastVersionID:   lo.ToPtr[int64](20),
					AuthorID:        2,
					ProjectAuthorID: 1,
				}
				templateRepo.EXPECT().GetByID(ctx, int64(10)).Return(&template, nil)

				version := domain.TemplateVersion{ID: 20, Data: []byte{1, 2, 3}}
				templateVersionRepo.EXPECT().GetByID(ctx, int64(20)).Return(&version, nil)

				variables := []domain.Variable{
					{
						ID:         31,
						Name:       "var1",
						Type:       variable_domain.TypeString,
						Expression: "expr31",
					},
				}
				variableRepo.EXPECT().ListByVersionID(ctx, int64(20)).Return(variables, nil)
				variableConstraintRepo.EXPECT().ListByVariableIDs(ctx, []int64{31}).Return(nil, errors.New("test5"))
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

			tt.setup(templateRepo, templateVersionRepo, variableRepo, variableConstraintRepo)

			usecase := New(templateRepo, templateVersionRepo, variableRepo, variableConstraintRepo)

			_, err := usecase.Handle(ctx, in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
