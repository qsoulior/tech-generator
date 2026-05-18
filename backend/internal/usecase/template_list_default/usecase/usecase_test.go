package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	sorting_domain "github.com/qsoulior/tech-generator/backend/internal/domain/sorting"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/template_list_default/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	in := domain.TemplateListDefaultIn{
		Page:   2,
		Size:   5,
		Filter: domain.TemplateListDefaultFilter{},
		Sorting: &sorting_domain.Sorting{
			Attribute: "name",
			Direction: "asc",
		},
	}

	want := domain.TemplateListDefaultOut{
		Templates:      make([]domain.Template, 5),
		TotalTemplates: 11,
		TotalPages:     3,
	}
	gofakeit.Slice(&want.Templates)

	templateRepo := NewMocktemplateRepository(ctrl)
	templateRepo.EXPECT().ListDefault(ctx, in).Return(want.Templates, nil)
	templateRepo.EXPECT().GetTotalDefault(ctx, in).Return(want.TotalTemplates, nil)

	usecase := New(templateRepo)
	got, err := usecase.Handle(ctx, in)
	require.NoError(t, err)
	require.Equal(t, want, *got)
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		setup func(templateRepo *MocktemplateRepository)
		in    domain.TemplateListDefaultIn
		want  string
	}{
		{
			name:  "in_Validate",
			setup: func(templateRepo *MocktemplateRepository) {},
			in:    domain.TemplateListDefaultIn{Page: 0, Size: 0},
			want:  domain.ErrValueInvalid.Error(),
		},
		{
			name: "templateRepo_ListDefault",
			setup: func(templateRepo *MocktemplateRepository) {
				templateRepo.EXPECT().ListDefault(ctx, gomock.Any()).Return(nil, errors.New("test1"))
			},
			in:   domain.TemplateListDefaultIn{Page: 1, Size: 1},
			want: "test1",
		},
		{
			name: "templateRepo_GetTotalDefault",
			setup: func(templateRepo *MocktemplateRepository) {
				templateRepo.EXPECT().ListDefault(ctx, gomock.Any()).Return([]domain.Template{}, nil)
				templateRepo.EXPECT().GetTotalDefault(ctx, gomock.Any()).Return(int64(0), errors.New("test2"))
			},
			in:   domain.TemplateListDefaultIn{Page: 1, Size: 1},
			want: "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			templateRepo := NewMocktemplateRepository(ctrl)
			tt.setup(templateRepo)

			usecase := New(templateRepo)
			_, err := usecase.Handle(ctx, tt.in)
			require.ErrorContains(t, err, tt.want)
		})
	}
}
