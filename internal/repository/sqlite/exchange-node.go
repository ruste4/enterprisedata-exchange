package sqliteRep

import (
	"context"
	"database/sql"
	"enterprisedata-exchange/internal/domain/entity"
	"fmt"
	"log/slog"
)

func NewExchangeNodeSqliteRepository(log *slog.Logger, db *sql.DB) *ExchangeNodeSqliteRepository {
	return &ExchangeNodeSqliteRepository{
		log: log,
		db:  db,
	}
}

type ExchangeNodeSqliteRepository struct {
	log *slog.Logger
	db  *sql.DB
}

func (r *ExchangeNodeSqliteRepository) Create(ctx context.Context, node *entity.ExchangeNode) (*entity.ExchangeNode, error) {
	const op = "repository.Create"

	query := `
		INSERT INTO exchange_nodes (this_node_code, c_description, node_code, prefix, this_prefix, c_state)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id, this_node_code, c_description, node_code, prefix, this_prefix, c_state, created_at, updated_at
	`

	var createdNode entity.ExchangeNode
	err := r.db.QueryRowContext(ctx, query,
		node.ThisNodeCode,
		node.Description,
		node.NodeCode,
		node.Prefix,
		node.ThisPrefix,
		node.State,
	).Scan(
		&createdNode.ID,
		&createdNode.ThisNodeCode,
		&createdNode.Description,
		&createdNode.NodeCode,
		&createdNode.Prefix,
		&createdNode.ThisPrefix,
		&createdNode.State,
		&createdNode.CreatedAt,
		&createdNode.UpdatedAt,
	)

	if err != nil {
		r.log.Error("failed to create exchange node",
			slog.String("operation", op),
			slog.Any("error", err),
			slog.Any("node", node),
		)
		return nil, fmt.Errorf("%s: %d", op, err)
	}

	return &createdNode, nil
}
