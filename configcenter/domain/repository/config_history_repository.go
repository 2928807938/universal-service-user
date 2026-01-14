package repository

import (
	"context"
	"github.com/2928807938/universal-service-user/configcenter/domain/entity"
)

// ConfigHistoryRepository records config changes.
type ConfigHistoryRepository interface {
	Create(ctx context.Context, history *entity.ConfigHistory) error
}
