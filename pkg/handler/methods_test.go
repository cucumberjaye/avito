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
