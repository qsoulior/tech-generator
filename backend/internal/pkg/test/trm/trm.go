package test_trm

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type TrKey struct{}

type TrManager struct{}

func New() *TrManager {
	return &TrManager{}
}

func (m *TrManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(context.WithValue(ctx, TrKey{}, struct{}{}))
}

func (m *TrManager) DoWithSettings(ctx context.Context, _ trm.Settings, fn func(ctx context.Context) error) error {
	return fn(context.WithValue(ctx, TrKey{}, struct{}{}))
}
