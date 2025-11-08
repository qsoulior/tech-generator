package service

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	variable_domain "github.com/qsoulior/tech-generator/backend/internal/domain/variable"
	"github.com/qsoulior/tech-generator/backend/internal/service/version_get/domain"
)

func TestService_Handle_Success(t *testing.T) {
	ctx := context.Background()
	createdAt := gofakeit.Date()
	versionID := gofakeit.Int64()

	tests := []struct {
		name  string
		setup func(versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository)
		want  domain.Version
	}{
		{
			name: "Variables",
			setup: func(versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				version := domain.Version{
					ID:         versionID,
					TemplateID: 1,
					Number:     2,
					CreatedAt:  createdAt,
					Data:       []byte{1, 2, 3},
				}
				versionRepo.EXPECT().GetByID(ctx, versionID).Return(&version, nil)

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
				variableRepo.EXPECT().ListByVersionID(ctx, versionID).Return(variables, nil)

				constraints := []domain.Constraint{
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
				constraintRepo.EXPECT().ListByVariableIDs(ctx, []int64{31, 32}).Return(constraints, nil)
			},
			want: domain.Version{
				ID:         versionID,
				TemplateID: 1,
				Number:     2,
				CreatedAt:  createdAt,
				Data:       []byte{1, 2, 3},
				Variables: []domain.Variable{
					{
						ID:         31,
						Name:       "var1",
						Type:       variable_domain.TypeString,
						Expression: "expr31",
						Constraints: []domain.Constraint{
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
						Constraints: []domain.Constraint{
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
			name: "NoVariables",
			setup: func(versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				version := domain.Version{
					ID:         versionID,
					TemplateID: 1,
					Number:     2,
					CreatedAt:  createdAt,
					Data:       []byte{1, 2, 3},
				}
				versionRepo.EXPECT().GetByID(ctx, versionID).Return(&version, nil)

				variables := []domain.Variable{}
				variableRepo.EXPECT().ListByVersionID(ctx, versionID).Return(variables, nil)
			},
			want: domain.Version{
				ID:         versionID,
				TemplateID: 1,
				Number:     2,
				CreatedAt:  createdAt,
				Data:       []byte{1, 2, 3},
				Variables:  []domain.Variable{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			versionRepo := NewMockversionRepository(ctrl)
			variableRepo := NewMockvariableRepository(ctrl)
			constraintRepo := NewMockconstraintRepository(ctrl)

			tt.setup(versionRepo, variableRepo, constraintRepo)

			usecase := New(versionRepo, variableRepo, constraintRepo)

			got, err := usecase.Handle(ctx, versionID)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.want, *got)
		})
	}
}

func TestService_Handle_Error(t *testing.T) {
	ctx := context.Background()
	versionID := gofakeit.Int64()

	tests := []struct {
		name  string
		setup func(versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository)
		want  string
	}{
		{
			name: "versionRepo_GetByID",
			setup: func(versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				versionRepo.EXPECT().GetByID(ctx, versionID).Return(nil, errors.New("test3"))
			},
			want: "test3",
		},
		{
			name: "domain_ErrTemplateVersionNotFound",
			setup: func(versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				versionRepo.EXPECT().GetByID(ctx, versionID).Return(nil, nil)
			},
			want: domain.ErrVersionNotFound.Error(),
		},
		{
			name: "variableRepo_ListByVersionID",
			setup: func(versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				version := domain.Version{ID: versionID, Data: []byte{1, 2, 3}}
				versionRepo.EXPECT().GetByID(ctx, versionID).Return(&version, nil)
				variableRepo.EXPECT().ListByVersionID(ctx, versionID).Return(nil, errors.New("test4"))
			},
			want: "test4",
		},
		{
			name: "constraintRepo_ListByVariableIDs",
			setup: func(versionRepo *MockversionRepository, variableRepo *MockvariableRepository, constraintRepo *MockconstraintRepository) {
				version := domain.Version{ID: versionID, Data: []byte{1, 2, 3}}
				versionRepo.EXPECT().GetByID(ctx, versionID).Return(&version, nil)

				variables := []domain.Variable{
					{
						ID:         31,
						Name:       "var1",
						Type:       variable_domain.TypeString,
						Expression: "expr31",
					},
				}
				variableRepo.EXPECT().ListByVersionID(ctx, versionID).Return(variables, nil)
				constraintRepo.EXPECT().ListByVariableIDs(ctx, []int64{31}).Return(nil, errors.New("test5"))
			},
			want: "test5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			versionRepo := NewMockversionRepository(ctrl)
			variableRepo := NewMockvariableRepository(ctrl)
			constraintRepo := NewMockconstraintRepository(ctrl)

			tt.setup(versionRepo, variableRepo, constraintRepo)

			usecase := New(versionRepo, variableRepo, constraintRepo)

			_, err := usecase.Handle(ctx, versionID)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
