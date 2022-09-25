package handler

import (
	"github.com/cucumberjaye/balanceAPI/pkg/repository"
	"net/http"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/accrual", h.accrual)
	mux.HandleFunc("/api/writeoff", h.writeOff)
	mux.HandleFunc("/api/transfer", h.transfer)
	mux.HandleFunc("/api/balance/", h.userBalance)

	return mux
}
