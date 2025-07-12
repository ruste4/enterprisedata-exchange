package service

import (
	"enterprisedata-exchange/internal/config"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileService struct {
	exchangeFilesPath string
	tempFilesPath     string
	logger            *slog.Logger
}

func NewFileService(cfg *config.Config, log *slog.Logger) *FileService {
	return &FileService{
		exchangeFilesPath: cfg.ExchangeFiles,
		tempFilesPath:     cfg.TempFiles,
		logger:            log,
	}
}

func (fs *FileService) WriteFile(dir string, filename string, file multipart.File) (string, error) {
	const op = "FileService.WriteFile"

	filePath := filepath.Join(dir, "/"+filename)
	dst, err := os.Create(filePath)
	if err != nil {
		fs.logger.Error("failed to create file", "op", op, "filePath", filePath, "err", err)
		return "", err

	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		fs.logger.Error("failed to copy file data", "op", op, "filePath", filePath, "err", err)
		return "", err
	}

	return filePath, nil
}
