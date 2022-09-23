package handler

import (
	"encoding/json"
	"fmt"
	"github.com/cucumberjaye/balanceAPI"
	"io"
	"net/http"
	"strconv"
	"strings"
	"log"
)

func requestUnmarshal(r *http.Request) (balanceAPI.UserData, error) {
	var result balanceAPI.UserData
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf(err.Error())
		return result, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Fatalf(err.Error())
		return result, err
	}

	return result, nil
}

func twoUsersUnmarshal(r *http.Request) (balanceAPI.TwoUsers, error) {
	var result balanceAPI.TwoUsers
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf(err.Error())
		return result, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Fatalf(err.Error())
		return result, err
	}

	return result, nil
}

func (h *Handler) accrual(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/accrual/" {
		if r.Method == http.MethodPost {
			userData, err := requestUnmarshal(r)
			if err != nil {
				http.Error(w, fmt.Sprintf("error on server"), http.StatusInternalServerError)
			}
			err = h.repo.Add(userData)
			if err != nil {
				http.Error(w, fmt.Sprintf("error on server"), http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, fmt.Sprintf("expect method POST at /api/accrual/, got %v", r.Method), http.StatusMethodNotAllowed)
		}
	}
}

func (h *Handler) writeOff(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/writeoff/" {
		if r.Method == http.MethodPost {
			userData, err := requestUnmarshal(r)
			if err != nil {
				http.Error(w, fmt.Sprintf("error on server"), http.StatusInternalServerError)
			}
			err = h.repo.Decrease(userData)
			if err != nil {
				http.Error(w, fmt.Sprintf("error on server"), http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, fmt.Sprintf("expect method POST at /api/write-off/, got %v", r.Method), http.StatusMethodNotAllowed)
		}
	}
}

func (h *Handler) transfer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/transfer/" {
		if r.Method == http.MethodPost {
			usersData, err := twoUsersUnmarshal(r)
			if err != nil {
				http.Error(w, fmt.Sprintf("error on server"), http.StatusInternalServerError)
			}
			err = h.repo.Transfer(usersData)
			if err != nil {
				http.Error(w, fmt.Sprintf("error on server"), http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
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
		sum, err := h.repo.GetBalance(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("error on server"), http.StatusInternalServerError)
		}
		jsonSum, err := json.Marshal(sum)
		if err != nil {
			http.Error(w, fmt.Sprintf("error on server"), http.StatusInternalServerError)
		}
		_, err = w.Write(jsonSum)
		if err != nil {
			http.Error(w, fmt.Sprintf("error on server"), http.StatusInternalServerError)
		}

	} else {
		http.Error(w, fmt.Sprintf("expect method GET at /api/balance/<id>, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}
