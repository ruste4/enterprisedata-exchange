package handler

import (
	"context"
	"encoding/json"
	"enterprisedata-exchange/internal/lib"
	"enterprisedata-exchange/internal/lib/rest/response"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	paramFile = "./ibparam.json" //todo переместить в конфиги
)

func GetIbParams(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "app.queries.handlers.version.GetIbParams"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		) // Todo check

		file, err := os.Open(paramFile)
		if err != nil {
			log.Error("file error", "file", paramFile)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("file error"))
			return
		}
		defer file.Close()

		var params map[string]interface{}
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&params); err != nil {
			log.Error("json decode error", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("json decode error"))
			return
		}

		params["InfobasePrefix"] = "BR"
		params["InfobaseDescription"] = "Интеграция с Бизнес.ру"
		params["ThisNodeCode"] = lib.GenerateUUID()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(params); err != nil {
			log.Error("json encode error", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("json encode error"))
		}
	}
}
