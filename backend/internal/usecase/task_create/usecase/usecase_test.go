package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	user_domain "github.com/qsoulior/tech-generator/backend/internal/domain/user"
	test_trm "github.com/qsoulior/tech-generator/backend/internal/pkg/test/trm"
	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_create/domain"
)

func TestUsecase_Handle_Success(t *testing.T) {
	ctx := context.Background()
	trCtx := context.WithValue(ctx, test_trm.TrKey{}, struct{}{})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	versionRepo := NewMockversionRepository(ctrl)
	taskRepo := NewMocktaskRepository(ctrl)
	publisher := NewMockpublisher(ctrl)
	trManager := test_trm.New()

	in := domain.TaskCreateIn{
		VersionID: 100,
		CreatorID: 1,
		Payload:   map[string]string{"k": "v"},
	}

	version := &domain.Version{
		ProjectAuthorID:  1,
		TemplateAuthorID: 2,
		TemplateUsers:    nil,
	}

	versionRepo.EXPECT().GetByID(ctx, in.VersionID).Return(version, nil)
	taskRepo.EXPECT().Insert(trCtx, in).Return(int64(50), nil)
	publisher.EXPECT().PublishTaskCreated(trCtx, int64(50)).Return(nil)

	usecase := New(versionRepo, taskRepo, publisher, trManager)
	err := usecase.Handle(ctx, in)
	require.NoError(t, err)
}

func TestUsecase_Handle_Error(t *testing.T) {
	ctx := context.Background()
	trCtx := context.WithValue(ctx, test_trm.TrKey{}, struct{}{})

	testErr := errors.New("test error")

	validIn := domain.TaskCreateIn{
		VersionID: 100,
		CreatorID: 1,
		Payload:   map[string]string{"k": "v"},
	}

	validVersion := &domain.Version{
		ProjectAuthorID:  1,
		TemplateAuthorID: 2,
		TemplateUsers:    nil,
	}

	tests := []struct {
		name  string
		setup func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository, publisher *Mockpublisher)
		in    domain.TaskCreateIn
		want  error
	}{
		{
			name: "versionRepo_GetByID",
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository, publisher *Mockpublisher) {
				versionRepo.EXPECT().GetByID(ctx, validIn.VersionID).Return(nil, testErr)
			},
			in:   validIn,
			want: testErr,
		},
		{
			name: "versionRepo_GetByID_NotFound",
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository, publisher *Mockpublisher) {
				versionRepo.EXPECT().GetByID(ctx, validIn.VersionID).Return(nil, nil)
			},
			in:   validIn,
			want: domain.ErrVersionNotFound,
		},
		{
			name: "version_Invalid_NoPermission",
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository, publisher *Mockpublisher) {
				version := &domain.Version{
					ProjectAuthorID:  999,
					TemplateAuthorID: 998,
					TemplateUsers:    nil,
				}
				versionRepo.EXPECT().GetByID(ctx, validIn.VersionID).Return(version, nil)
			},
			in:   validIn,
			want: domain.ErrVersionInvalid,
		},
		{
			name: "version_Invalid_WrongRole",
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository, publisher *Mockpublisher) {
				version := &domain.Version{
					ProjectAuthorID:  999,
					TemplateAuthorID: 998,
					TemplateUsers: []domain.TemplateUser{
						{ID: validIn.CreatorID, Role: user_domain.Role("invalid")},
					},
				}
				versionRepo.EXPECT().GetByID(ctx, validIn.VersionID).Return(version, nil)
			},
			in:   validIn,
			want: domain.ErrVersionInvalid,
		},
		{
			name: "taskRepo_Insert",
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository, publisher *Mockpublisher) {
				versionRepo.EXPECT().GetByID(ctx, validIn.VersionID).Return(validVersion, nil)
				taskRepo.EXPECT().Insert(trCtx, validIn).Return(int64(0), testErr)
			},
			in:   validIn,
			want: testErr,
		},
		{
			name: "publisher_PublishTaskCreated",
			setup: func(versionRepo *MockversionRepository, taskRepo *MocktaskRepository, publisher *Mockpublisher) {
				versionRepo.EXPECT().GetByID(ctx, validIn.VersionID).Return(validVersion, nil)
				taskRepo.EXPECT().Insert(trCtx, validIn).Return(int64(50), nil)
				publisher.EXPECT().PublishTaskCreated(trCtx, int64(50)).Return(testErr)
			},
			in:   validIn,
			want: testErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			versionRepo := NewMockversionRepository(ctrl)
			taskRepo := NewMocktaskRepository(ctrl)
			publisher := NewMockpublisher(ctrl)
			trManager := test_trm.New()
			tt.setup(versionRepo, taskRepo, publisher)

			usecase := New(versionRepo, taskRepo, publisher, trManager)
			err := usecase.Handle(ctx, tt.in)
			require.ErrorIs(t, err, tt.want)
		})
	}
}
