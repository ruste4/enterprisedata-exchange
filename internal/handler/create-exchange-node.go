package handler

import (
	"context"
	"encoding/json"
	"enterprisedata-exchange/internal/domain/entity"
	"enterprisedata-exchange/internal/lib/rest/response"
	"enterprisedata-exchange/internal/usecase"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

func CreateExchangeNode(ctx context.Context, log *slog.Logger, uc *usecase.ExchangeUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData entity.CreateExchangeNodeDto

		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			log.Error("failed to decode request body", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid json format"))
			return
		}
		defer r.Body.Close()

		_, err := uc.CreateExchangeNode(ctx, requestData)
		if err != nil {
			log.Error("integration not created", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("integration not found"))
			return
		}

		log.Info("Integration created")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.Ok())
	}
}
