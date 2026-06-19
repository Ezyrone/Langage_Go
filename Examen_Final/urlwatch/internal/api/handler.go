package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jory/urlwatch/internal/domain"
	"github.com/jory/urlwatch/internal/pool"
)

type Handler struct {
	checker domain.Checker
	store   domain.Store
	logger  *slog.Logger
}

func NewHandler(checker domain.Checker, store domain.Store, logger *slog.Logger) *Handler {
	return &Handler{checker: checker, store: store, logger: logger}
}

func (h *Handler) PostChecks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only POST is accepted")
		return
	}

	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body: "+err.Error())
		return
	}

	urls, opts, err := validateAndNormalize(req)
	if err != nil {
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			writeError(w, http.StatusBadRequest, "invalid_request", ve.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}

	batch := pool.Run(r.Context(), h.checker, urls, opts)
	batch.ID = generateID()

	if err := h.store.Save(r.Context(), batch); err != nil {
		h.logger.Error("failed to save batch", "error", err)
		writeError(w, http.StatusInternalServerError, "internal", "failed to persist batch")
		return
	}

	h.logger.Info("batch created", "batch_id", batch.ID, "total", batch.Summary.Total, "up", batch.Summary.Up, "down", batch.Summary.Down)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(batch)
}

func (h *Handler) GetCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET is accepted")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/v1/checks/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid_request", "missing batch id")
		return
	}

	batch, err := h.store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrBatchNotFound) {
			writeError(w, http.StatusNotFound, "batch_not_found", "aucun lot avec l'id "+id)
			return
		}
		h.logger.Error("failed to get batch", "error", err)
		writeError(w, http.StatusInternalServerError, "internal", "failed to retrieve batch")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batch)
}

func Healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorBody{
		Error: ErrorDetail{Code: code, Message: message},
	})
}

func generateID() string {
	b := make([]byte, 3)
	rand.Read(b)
	return "b_" + hex.EncodeToString(b)
}
