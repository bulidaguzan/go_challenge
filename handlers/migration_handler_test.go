package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMigrateCSV(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    
    // Create mock db
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error creating mock db: %v", err)
    }
    defer db.Close()

    // Create handler
    handler := NewMigrationHandler(db)

    // Test cases
    tests := []struct {
        name           string
        csvContent     string
        setupMock      func()
        expectedStatus int
    }{
        {
            name: "successful migration",
            csvContent: `id,user_id,amount,datetime
1,1,100.00,2024-01-01T00:00:00Z
2,1,-50.00,2024-01-02T00:00:00Z`,
            setupMock: func() {
                mock.ExpectPrepare("WITH del AS")
                mock.ExpectExec("WITH del AS").
                    WithArgs(1, 1, 100.00, "2024-01-01T00:00:00Z").
                    WillReturnResult(sqlmock.NewResult(1, 1))
                mock.ExpectExec("WITH del AS").
                    WithArgs(2, 1, -50.00, "2024-01-02T00:00:00Z").
                    WillReturnResult(sqlmock.NewResult(1, 1))
            },
            expectedStatus: http.StatusOK,
        },
        // Add more test cases as needed
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mock
            tt.setupMock()

            // Create multipart form
            body := new(bytes.Buffer)
            writer := multipart.NewWriter(body)
            part, err := writer.CreateFormFile("file", "test.csv")
            assert.NoError(t, err)
            part.Write([]byte(tt.csvContent))
            writer.Close()

            // Create request
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            c.Request, _ = http.NewRequest("POST", "/migrate", body)
            c.Request.Header.Set("Content-Type", writer.FormDataContentType())

            // Execute request
            handler.MigrateCSV(c)

            // Assert results
            assert.Equal(t, tt.expectedStatus, w.Code)
        })
    }
}