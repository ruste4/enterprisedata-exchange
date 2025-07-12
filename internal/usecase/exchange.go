package usecase

import (
	"context"
	"enterprisedata-exchange/internal/config"
	"enterprisedata-exchange/internal/domain/entity"
	"enterprisedata-exchange/internal/domain/service"
	"fmt"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
)

const (
	fileName = "data.zip"
)

type ExchangeUseCase struct {
	cfg       *config.Config
	log       *slog.Logger
	enService *service.ExchangeNodeService
	fService  *service.FileService
}

func NewExchangeNodeUseCase(
	cfg *config.Config,
	log *slog.Logger,
	enService *service.ExchangeNodeService,
	fService *service.FileService,
) *ExchangeUseCase {
	return &ExchangeUseCase{
		cfg:       cfg,
		log:       log,
		enService: enService,
		fService:  fService,
	}
}

func (uc *ExchangeUseCase) CreateExchangeNode(ctx context.Context, nodeDto entity.CreateExchangeNodeDto) (*entity.ExchangeNode, error) {
	const op = "ExchangeUseCase.CreateExchangeNode"

	createdNode, err := uc.enService.CreateExchangeNode(ctx, nodeDto)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return createdNode, nil
}

func (uc *ExchangeUseCase) PutFile(ctx context.Context, login string, sessionId string, partNumber string, file multipart.File) error {
	const op = "ExchangeUseCase.Putfile"
	sessionPath := filepath.Dir(uc.cfg.TempFiles + "/" + login + "/" + sessionId)
	if err := os.MkdirAll(sessionPath, 0755); err != nil {
		uc.log.Error("failed to create session directory", "op", op, "path", sessionPath, "err", err)
		return err
	}

	_, err := uc.fService.WriteFile(sessionPath, fileName+"."+partNumber, file)
	if err != nil {
		uc.log.Error("failed to write file part", "op", op, "sessionPath", sessionPath, "partNumber", partNumber, "err", err)
		return fmt.Errorf("%s: failed to write file part %s: %w", op, partNumber, err)
	}

	return nil
}
