package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	//repo repo
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/accrual/", h.accrual)
	mux.HandleFunc("/api/write-off/", h.writeOff)
	mux.HandleFunc("/api/transfer/", h.transfer)
	mux.HandleFunc("/api/balance/", h.userBalance)

	return mux
}

func (h *Handler) accrual(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/accrual/" {
		if r.Method == http.MethodPost {
			//
		} else {
			http.Error(w, fmt.Sprintf("expect method POST at /api/accrual/, got %v", r.Method), http.StatusMethodNotAllowed)
		}
	}
}

func (h *Handler) writeOff(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/write-off/" {
		if r.Method == http.MethodPost {
			//
		} else {
			http.Error(w, fmt.Sprintf("expect method POST at /api/write-off/, got %v", r.Method), http.StatusMethodNotAllowed)
		}
	}
}

func (h *Handler) transfer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/transfer/" {
		if r.Method == http.MethodPost {
			//
		} else {
			http.Error(w, fmt.Sprintf("expect method POST at /api/transfer/, got %v", r.Method), http.StatusMethodNotAllowed)
		}
	}
}

func (h *Handler) userBalance(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "expect /api/balance/<id> in balance handler", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodGet {
		id = id * 2
	} else {
		http.Error(w, fmt.Sprintf("expect method GET at /api/balance/<id>, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}
