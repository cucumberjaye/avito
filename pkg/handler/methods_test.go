package handler

import (
	"bytes"
	"errors"
	"github.com/cucumberjaye/balanceAPI"
	"github.com/cucumberjaye/balanceAPI/pkg/repository"
	mock_repository "github.com/cucumberjaye/balanceAPI/pkg/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Accrual(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockBalance, userData balanceAPI.UserData)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUserData        balanceAPI.UserData
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"id":1,"name":"test","surname":"test","sum":100,"comment":"test"}`,
			inputUserData: balanceAPI.UserData{
				User: balanceAPI.User{
					Id:      1,
					Name:    "test",
					Surname: "test",
				},
				Sum:     100,
				Comment: "test",
			},
			mockBehavior: func(s *mock_repository.MockBalance, userData balanceAPI.UserData) {
				s.EXPECT().Add(userData).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"id":1}`,
			inputUserData:        balanceAPI.UserData{},
			mockBehavior:         func(s *mock_repository.MockBalance, userData balanceAPI.UserData) {},
			expectedStatusCode:   400,
			expectedResponseBody: "invalid input body\n",
		},
		{
			name:      "Repository Failure",
			inputBody: `{"id":1,"name":"test","surname":"test","sum":100,"comment":"test"}`,
			inputUserData: balanceAPI.UserData{
				User: balanceAPI.User{
					Id:      1,
					Name:    "test",
					Surname: "test",
				},
				Sum:     100,
				Comment: "test",
			},
			mockBehavior: func(s *mock_repository.MockBalance, userData balanceAPI.UserData) {
				s.EXPECT().Add(userData).Return(errors.New("repository failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "repository failure\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			add := mock_repository.NewMockBalance(c)
			testCase.mockBehavior(add, testCase.inputUserData)

			repos := &repository.Repository{Balance: add}
			handler := NewHandler(repos)

			//Init Endpoint
			r := http.NewServeMux()
			r.HandleFunc("/api/accrual", handler.accrual)

			//Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/accrual", bytes.NewBufferString(testCase.inputBody))

			//Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_WriteOff(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockBalance, userData balanceAPI.UserData)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUserData        balanceAPI.UserData
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"id":1,"name":"test","surname":"test","sum":100,"comment":"test"}`,
			inputUserData: balanceAPI.UserData{
				User: balanceAPI.User{
					Id:      1,
					Name:    "test",
					Surname: "test",
				},
				Sum:     100,
				Comment: "test",
			},
			mockBehavior: func(s *mock_repository.MockBalance, userData balanceAPI.UserData) {
				s.EXPECT().Decrease(userData).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"id":1}`,
			inputUserData:        balanceAPI.UserData{},
			mockBehavior:         func(s *mock_repository.MockBalance, userData balanceAPI.UserData) {},
			expectedStatusCode:   400,
			expectedResponseBody: "invalid input body\n",
		},
		{
			name:      "Repository Failure",
			inputBody: `{"id":1,"name":"test","surname":"test","sum":100,"comment":"test"}`,
			inputUserData: balanceAPI.UserData{
				User: balanceAPI.User{
					Id:      1,
					Name:    "test",
					Surname: "test",
				},
				Sum:     100,
				Comment: "test",
			},
			mockBehavior: func(s *mock_repository.MockBalance, userData balanceAPI.UserData) {
				s.EXPECT().Decrease(userData).Return(errors.New("repository failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "repository failure\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			decrease := mock_repository.NewMockBalance(c)
			testCase.mockBehavior(decrease, testCase.inputUserData)

			repos := &repository.Repository{Balance: decrease}
			handler := NewHandler(repos)

			//Init Endpoint
			r := http.NewServeMux()
			r.HandleFunc("/api/writeoff", handler.writeOff)

			//Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/writeoff", bytes.NewBufferString(testCase.inputBody))

			//Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_Transfer(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockBalance, userData balanceAPI.TwoUsers)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUserData        balanceAPI.TwoUsers
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"add_user":{"id":1,"name":"test","surname":"test"},"decrease_user":{"id":2,"name":"test","surname":"test"},"sum":100,"comment":"test"}`,
			inputUserData: balanceAPI.TwoUsers{
				AddMoneyUser: balanceAPI.User{
					Id:      1,
					Name:    "test",
					Surname: "test",
				},
				DecreaseMoneyUser: balanceAPI.User{
					Id:      2,
					Name:    "test",
					Surname: "test",
				},
				Sum:     100,
				Comment: "test",
			},
			mockBehavior: func(s *mock_repository.MockBalance, userData balanceAPI.TwoUsers) {
				s.EXPECT().Transfer(userData).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"id":1}`,
			inputUserData:        balanceAPI.TwoUsers{},
			mockBehavior:         func(s *mock_repository.MockBalance, userData balanceAPI.TwoUsers) {},
			expectedStatusCode:   400,
			expectedResponseBody: "invalid input body\n",
		},
		{
			name:      "Repository Failure",
			inputBody: `{"add_user":{"id":1,"name":"test","surname":"test"},"decrease_user":{"id":2,"name":"test","surname":"test"},"sum":100,"comment":"test"}`,
			inputUserData: balanceAPI.TwoUsers{
				AddMoneyUser: balanceAPI.User{
					Id:      1,
					Name:    "test",
					Surname: "test",
				},
				DecreaseMoneyUser: balanceAPI.User{
					Id:      2,
					Name:    "test",
					Surname: "test",
				},
				Sum:     100,
				Comment: "test",
			},
			mockBehavior: func(s *mock_repository.MockBalance, userData balanceAPI.TwoUsers) {
				s.EXPECT().Transfer(userData).Return(errors.New("repository failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "repository failure\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			transfer := mock_repository.NewMockBalance(c)
			testCase.mockBehavior(transfer, testCase.inputUserData)

			repos := &repository.Repository{Balance: transfer}
			handler := NewHandler(repos)

			//Init Endpoint
			r := http.NewServeMux()
			r.HandleFunc("/api/transfer", handler.transfer)

			//Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/transfer", bytes.NewBufferString(testCase.inputBody))

			//Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_GetBalance(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockBalance, userId int)

	testTable := []struct {
		name                 string
		inputUserData        int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "OK",
			inputUserData: 1,
			mockBehavior: func(s *mock_repository.MockBalance, userId int) {
				s.EXPECT().GetBalance(userId).Return(100, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `100`,
		},
		{
			name:          "Repository Failure",
			inputUserData: 1,
			mockBehavior: func(s *mock_repository.MockBalance, userId int) {
				s.EXPECT().GetBalance(userId).Return(0, errors.New("repository failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "repository failure\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			getBalance := mock_repository.NewMockBalance(c)
			testCase.mockBehavior(getBalance, testCase.inputUserData)

			repos := &repository.Repository{Balance: getBalance}
			handler := NewHandler(repos)

			//Init Endpoint
			r := http.NewServeMux()
			r.HandleFunc("/api/balance/", handler.userBalance)

			//Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/balance/1", nil)

			//Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_GetTransactions(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockBalance, userId int, sortBy string)

	testTable := []struct {
		name                 string
		inputUserData        int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "OK",
			inputUserData: 1,
			mockBehavior: func(s *mock_repository.MockBalance, userId int, sortBy string) {
				s.EXPECT().GetTransactions(1, "").Return([]balanceAPI.Transactions{balanceAPI.Transactions{Sum: 100, Comment: "test", Date: "test", No: 1}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"No":1,"Sum":100,"Comment":"test","Date":"test"}]`,
		},
		{
			name:          "Repository Failure",
			inputUserData: 1,
			mockBehavior: func(s *mock_repository.MockBalance, userId int, sortBy string) {
				s.EXPECT().GetTransactions(userId, "").Return(nil, errors.New("repository failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "repository failure\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			getTransactions := mock_repository.NewMockBalance(c)
			testCase.mockBehavior(getTransactions, testCase.inputUserData, "")

			repos := &repository.Repository{Balance: getTransactions}
			handler := NewHandler(repos)

			//Init Endpoint
			r := http.NewServeMux()
			r.HandleFunc("/api/transactions/", handler.userTransactions)

			//Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/transactions/1", nil)

			//Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
