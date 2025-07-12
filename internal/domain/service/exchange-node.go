package service

import (
	"context"
	"enterprisedata-exchange/internal/domain/entity"
	"enterprisedata-exchange/internal/domain/repository"
	"fmt"
	"log/slog"
)

type ExchangeNodeService struct {
	log slog.Logger
	rep repository.ExchangeNodeRepository
}

func NewExchangeNodeService(log *slog.Logger, rep repository.ExchangeNodeRepository) *ExchangeNodeService {
	return &ExchangeNodeService{
		log: *log,
		rep: rep,
	}
}

func (s *ExchangeNodeService) CreateExchangeNode(ctx context.Context, nodeDto entity.CreateExchangeNodeDto) (*entity.ExchangeNode, error) {
	const op = "service.CreateExchangeNode"

	newNode := &entity.ExchangeNode{
		ThisNodeCode: nodeDto.MainExchangeParameters.CorrespondentNodeCode,
		Description:  nodeDto.MainExchangeParameters.ThisInfobaseDescription,
		NodeCode:     nodeDto.MainExchangeParameters.NodeCode,
		Prefix:       nodeDto.MainExchangeParameters.SourceInfobasePrefix,
		ThisPrefix:   nodeDto.MainExchangeParameters.DestinationInfobasePrefix,
		State:        "active",
	}

	createdNode, err := s.rep.Create(ctx, newNode)
	if err != nil {
		s.log.Error("exchange node not created",
			slog.Any("error", err),
			slog.Any("nodeDto", nodeDto),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info(
		fmt.Sprintf("exchange node created id:%d", createdNode.ID),
		slog.Any("nodeDto", nodeDto),
		slog.Any("createdNode", createdNode),
	)

	return createdNode, nil
}
