package handlers

import (
	"net/http"

	"aether-console/internal/services/health"
)

type HealthHandler struct {
	svc *health.Service
}

func NewHealthHandler(svc *health.Service) *HealthHandler {
	return &HealthHandler{svc: svc}
}

func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	msg, err := h.svc.Check(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(msg))
}
