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
	"github.com/papidb/aqua/pkg/http/server"
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
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Check the response body
	expected := "{\"message\":\"Hello World\"}"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHealthHandler(t *testing.T) {
	// Create a new instance of the server

	app, err := config.New(*env)

	if err != nil {
		t.Fatal(err)
		return
	}
	s := &server.Server{
		App: app,
	}

	r := gin.New()
	t.Log("hi")

	r.GET("/health", func(ctx *gin.Context) {
		healthHandler(ctx, s)
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

	if result["status"] != "up" {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), "up")
	}

	if result["message"] != "It's healthy" {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), "It's healthy")
	}

}
