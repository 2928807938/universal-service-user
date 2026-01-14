package repository

import (
	"context"
	"universal-service-user/configcenter/domain/entity"
)

// ConfigHistoryRepository records config changes.
type ConfigHistoryRepository interface {
	Create(ctx context.Context, history *entity.ConfigHistory) error
}
