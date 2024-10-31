package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fintech-backend/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    
    // Crear mock de la base de datos
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error creating mock db: %v", err)
    }
    defer db.Close()

    // Crear handler con la base de datos mock
    handler := NewBalanceHandler(db)

    tests := []struct {
        name           string
        userID         string
        setupMock      func()
        expectedStatus int
        expectedBody   models.BalanceResponse
    }{
        {
            name:   "successful balance retrieval",
            userID: "1",
            setupMock: func() {
                rows := sqlmock.NewRows([]string{"balance", "debits", "credits"}).
                    AddRow(100.50, 2, 3)
                mock.ExpectQuery("^SELECT").
                    WithArgs(1).
                    WillReturnRows(rows)
            },
            expectedStatus: http.StatusOK,
            expectedBody: models.BalanceResponse{
                Balance:      100.50,
                TotalDebits:  2,
                TotalCredits: 3,
            },
        },
        {
            name:   "invalid user ID",
            userID: "invalid",
            setupMock: func() {
                // No need to setup mock as it should fail before DB query
            },
            expectedStatus: http.StatusBadRequest,
        },
        {
            name:   "database error",
            userID: "1",
            setupMock: func() {
                mock.ExpectQuery("^SELECT").
                    WithArgs(1).
                    WillReturnError(sql.ErrConnDone)
            },
            expectedStatus: http.StatusInternalServerError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mock
            tt.setupMock()

            // Create request
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            c.Params = []gin.Param{{Key: "user_id", Value: tt.userID}}

            // Execute request
            handler.GetBalance(c)

            // Assert results
            assert.Equal(t, tt.expectedStatus, w.Code)
            
            if tt.expectedStatus == http.StatusOK {
                var response models.BalanceResponse
                err := json.NewDecoder(w.Body).Decode(&response)
                assert.NoError(t, err)
                assert.Equal(t, tt.expectedBody, response)
            }
        })
    }
}