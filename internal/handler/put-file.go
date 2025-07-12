package handler

import (
	"context"
	resp "enterprisedata-exchange/internal/lib/rest/response"
	"enterprisedata-exchange/internal/usecase"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func PutFile(ctx context.Context, log *slog.Logger, uc *usecase.ExchangeUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.PutFile"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		login, _, _ := r.BasicAuth()

		sessionId := chi.URLParam(r, "SessionId")
		if sessionId == "" {
			log.Error("Query param SessionId empty")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Query parameter SessionId empty"))
			return
		}

		partNumber := chi.URLParam(r, "PartNumber")
		if partNumber == "" {
			log.Error("Query param PartNumber empty")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Query parameter 'PartNumber empty"))
			return
		}

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Error("failed to parse multipart form", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to parse form"))
			return
		}

		file, header, err := r.FormFile("file") // "file" - имя поля в форме
		if err != nil {
			log.Error("failed to get file from form", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("file is required"))
			return
		}
		defer file.Close()

		log.Info("File received",
			slog.String("filename", header.Filename),
			slog.Int64("size", header.Size),
			slog.String("content_type", header.Header.Get("Content-Type")),
		)

		err = uc.PutFile(ctx, login, sessionId, partNumber, file)
		if err != nil {
			log.Error("failed to зге file", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to put file"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp.Ok())
	}
}
