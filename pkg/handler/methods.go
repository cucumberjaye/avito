package handler

import (
	"encoding/json"
	"fmt"
	"github.com/cucumberjaye/balanceAPI"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) accrual(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/accrual" {
		if r.Method == http.MethodPost {
			userData, err := requestUnmarshal(r)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusBadRequest)
				return
			}
			err = h.repo.Add(userData)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, fmt.Sprintf("expect method POST at /api/accrual/, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		http.NotFound(w, r)
	}
}

func (h *Handler) writeOff(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/writeoff" {
		fmt.Println(r.Method)
		if r.Method == http.MethodPost {
			userData, err := requestUnmarshal(r)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusBadRequest)
				return
			}
			err = h.repo.Decrease(userData)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, fmt.Sprintf("expect method POST at /api/write-off/, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		http.NotFound(w, r)
		return
	}
}

func (h *Handler) transfer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/transfer" {
		if r.Method == http.MethodPost {
			usersData, err := twoUsersUnmarshal(r)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusBadRequest)
				return
			}
			err = h.repo.Transfer(usersData)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, fmt.Sprintf("expect method POST at /api/transfer/, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		http.NotFound(w, r)
		return
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
			http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusInternalServerError)
			return
		}
		jsonSum, err := json.Marshal(sum)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonSum)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, fmt.Sprintf("expect method GET at /api/balance/<id>, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

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
