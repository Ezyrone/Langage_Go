package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jory/urlwatch/internal/checker"
	"github.com/jory/urlwatch/internal/domain"
	"github.com/jory/urlwatch/internal/store"
)

func setupTestRouter() http.Handler {
	mock := checker.NewMockChecker(map[string]domain.CheckResult{
		"https://go.dev": {URL: "https://go.dev", StatusCode: 200, OK: true, LatencyMs: 10},
	}, 0)
	st := store.NewMemoryStore()
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	h := NewHandler(mock, st, logger)
	return NewRouter(h, logger)
}

func TestPostChecks_Success(t *testing.T) {
	router := setupTestRouter()

	body := `{"urls":["https://go.dev"],"options":{"concurrency":2,"timeout_ms":3000}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/checks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var batch domain.Batch
	if err := json.NewDecoder(rec.Body).Decode(&batch); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if batch.ID == "" {
		t.Error("batch_id should not be empty")
	}
	if batch.Summary.Total != 1 {
		t.Errorf("total = %d, want 1", batch.Summary.Total)
	}
	if batch.Summary.Up != 1 {
		t.Errorf("up = %d, want 1", batch.Summary.Up)
	}
}

func TestPostChecks_InvalidBody(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/v1/checks", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestPostChecks_InvalidURL(t *testing.T) {
	router := setupTestRouter()

	body := `{"urls":["not-a-url"]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/checks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestGetCheck_NotFound(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/v1/checks/nonexistent", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	var errBody ErrorBody
	if err := json.NewDecoder(rec.Body).Decode(&errBody); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if errBody.Error.Code != "batch_not_found" {
		t.Errorf("error code = %q, want %q", errBody.Error.Code, "batch_not_found")
	}
}

func TestGetCheck_AfterPost(t *testing.T) {
	router := setupTestRouter()

	body := `{"urls":["https://go.dev"]}`
	postReq := httptest.NewRequest(http.MethodPost, "/v1/checks", bytes.NewBufferString(body))
	postReq.Header.Set("Content-Type", "application/json")
	postRec := httptest.NewRecorder()
	router.ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, want %d", postRec.Code, http.StatusCreated)
	}

	var created domain.Batch
	json.NewDecoder(postRec.Body).Decode(&created)

	getReq := httptest.NewRequest(http.MethodGet, "/v1/checks/"+created.ID, nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want %d", getRec.Code, http.StatusOK)
	}

	var fetched domain.Batch
	json.NewDecoder(getRec.Body).Decode(&fetched)
	if fetched.ID != created.ID {
		t.Errorf("batch_id = %q, want %q", fetched.ID, created.ID)
	}
}

func TestHealthz(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestPostChecks_ValidationConcurrency(t *testing.T) {
	router := setupTestRouter()

	body := `{"urls":["https://go.dev"],"options":{"concurrency":100}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/checks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
