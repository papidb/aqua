package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/config"
	"github.com/papidb/aqua/pkg/entities/customers"
	middlewares "github.com/papidb/aqua/services/core/pkg/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

var env = &config.Env{
	PostgresDatabase: "database",
	PostgresPassword: "password",
	PostgresUser:     "user",
}

func TestMain(m *testing.M) {
	teardown, err := config.MustStartPostgresContainer(env)
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	fmt.Println("started postgres container")
	m.Run()

	ctx := context.Background()
	if teardown != nil && teardown(ctx,
		testcontainers.RemoveVolumes(),
		testcontainers.StopContext(ctx),
	) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestHelloWorldHandler(t *testing.T) {
	r := gin.New()
	r.GET("/", helloWorldHandler)

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "status code should be 200")
	// Check the response body
	expected := "{\"message\":\"Hello World\"}"
	assert.Equal(t, expected, rr.Body.String(), "response body should be Hello World")
}

func TestHealthHandler(t *testing.T) {
	// Create a new instance of the server

	app, err := config.New(*env)

	if err != nil {
		t.Fatal(err)
		return
	}

	r := gin.New()

	r.GET("/health", func(ctx *gin.Context) {
		healthHandler(ctx, app)
	})

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Convert *bytes.Buffer to map
	var result map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	fmt.Println(result)
	// Check the response body

	assert.Equal(t, result["status"], "up", "status should be up")
	assert.Equal(t, result["message"], "It's healthy", "message should be It's healthy")
}

func TestCreateCustomerHandler(t *testing.T) {
	r := gin.New()
	r.POST(
		"/customers",
		middlewares.ValidationMiddleware(&customers.CreateCustomerDTO{}),
		createCustomerHandler,
	)

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid payload",
			payload:        `{"name":"Jane Doe","email":"jane.doe@example.com"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Validation passed"}`,
		},
		{
			name:           "Invalid payload - missing required fields",
			payload:        `{"name":"","email":""}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"We could not validate your request.","data":{"validation_error":{"email":"cannot be blank","name":"cannot be blank"}}}`,
		},
		{
			name:           "Invalid email payload",
			payload:        `{"name":"Jane Doe","email":"invalid.com"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"We could not validate your request.","data":{"validation_error":{"email":"must be a valid email address"}}}`,
		},
		{
			name:           "Malformed JSON",
			payload:        `{"name":"Jane Doe","email":"jane.doe@example.com",`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"details":"unexpected EOF","error":"Invalid request format"}`,
		},
		{
			name:           "Unsupported media type",
			payload:        `<xml><name>Jane Doe</name></xml>`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"details":"invalid character '<' looking for beginning of value", "error":"Invalid request format"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test HTTP request
			req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer([]byte(tt.payload)))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Perform the test HTTP request
			r.ServeHTTP(w, req)

			// Assert the status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert the response body
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}

}
