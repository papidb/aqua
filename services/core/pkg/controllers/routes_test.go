package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/config"
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
	t.Log("hi")

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
