package repository

import (
	"context"
	"enterprisedata-exchange/internal/domain/entity"
)

type ExchangeNodeRepository interface {
	Create(ctx context.Context, node *entity.ExchangeNode) (*entity.ExchangeNode, error)
}
