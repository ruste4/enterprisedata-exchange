//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"enterprisedata-exchange/internal/config"
	"enterprisedata-exchange/internal/domain/repository"
	"enterprisedata-exchange/internal/domain/service"
	sqliteRep "enterprisedata-exchange/internal/repository/sqlite"
	"enterprisedata-exchange/internal/usecase"
	"enterprisedata-exchange/pkg/database"
	"enterprisedata-exchange/pkg/logger"
	"log/slog"

	"github.com/google/wire"
)

var BaseSet = wire.NewSet(
	providerConfig,
	providerLogger,
	providerDBConnect,

	wire.Bind(new(repository.ExchangeNodeRepository), new(*sqliteRep.ExchangeNodeSqliteRepository)),
	providerExchangeNodeRepository,
)

func InitExchangeService() (*service.ExchangeNodeService, func(), error) {
	wire.Build(
		BaseSet,
		service.NewExchangeNodeService,
	)

	return nil, nil, nil
}

func InitUseCase() (*usecase.ExchangeUseCase, func(), error) {
	wire.Build(
		BaseSet,
		service.NewExchangeNodeService,
		providerFileService,
		usecase.NewExchangeNodeUseCase,
	)

	return nil, nil, nil
}

func providerFileService(cfg *config.Config, log *slog.Logger) (*service.FileService, func(), error) {
	return service.NewFileService(cfg, log), nil, nil
}

func providerConfig() (*config.Config, func(), error) {
	return config.MustLoad(), nil, nil
}

func providerLogger(cfg *config.Config) (*slog.Logger, func(), error) {
	return logger.SetupLogger(cfg.Env), nil, nil
}

func providerDBConnect(cfg *config.Config) (*sql.DB, func(), error) {
	conn, err := database.Connect(cfg)
	return conn, nil, err
}

func providerExchangeNodeRepository(log *slog.Logger, db *sql.DB) (*sqliteRep.ExchangeNodeSqliteRepository, func(), error) {
	return sqliteRep.NewExchangeNodeSqliteRepository(log, db), nil, nil
}
