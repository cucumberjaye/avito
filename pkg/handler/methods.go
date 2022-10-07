package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumberjaye/balanceAPI"
	"github.com/cucumberjaye/balanceAPI/pkg/apilayer"
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
	var res float64
	if r.Method == http.MethodGet {
		sum, err := h.repo.GetBalance(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusInternalServerError)
			return
		}
		res = float64(sum)

		query := r.URL.Query()
		if query["currency"] != nil {
			newSum, err := apilayer.Convert(query["currency"][0], res)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			res = newSum
		}
		jsonSum, err := json.Marshal(res)
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

func (h *Handler) userTransactions(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "expect /api/transactions/<id> in transactions handler", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var history []balanceAPI.Transactions
	if r.Method == http.MethodGet {
		query := r.URL.Query()
		if query["sort"] != nil {
			history, err = h.repo.GetTransactions(id, query["sort"][0])
		} else {
			history, err = h.repo.GetTransactions(id, "")
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusInternalServerError)
			return
		}

		jsonHistory, err := json.Marshal(history)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(jsonHistory)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusInternalServerError)
			return
		}
	}
}

func requestUnmarshal(r *http.Request) (balanceAPI.UserData, error) {
	var result balanceAPI.UserData
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf(err.Error())
		return result, err
	}

	count := strings.Split(string(b), ",")
	if len(count) != 5 {
		return result, errors.New("invalid input body")
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Printf(err.Error())
		return result, err
	}

	if result.Sum <= 0 {
		return result, errors.New("invalid input body")
	}

	return result, nil
}

func twoUsersUnmarshal(r *http.Request) (balanceAPI.TwoUsers, error) {
	var result balanceAPI.TwoUsers
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf(err.Error())
		return result, err
	}

	count := strings.Split(string(b), ",")
	if len(count) != 8 {
		return result, errors.New("invalid input body")
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Printf(err.Error())
		return result, err
	}

	if result.Sum <= 0 {
		return result, errors.New("invalid input body")
	}

	return result, nil
}
