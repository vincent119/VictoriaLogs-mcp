package victorialogs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vincent119/victorialogs-mcp/internal/util"
)

func TestClient_Health(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, util.AuthConfig{}, 10*time.Second)
	defer client.Close()

	result, err := client.Health(context.Background())
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}

	if result.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", result.Status)
	}
}

func TestClient_Health_Unhealthy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(server.URL, util.AuthConfig{}, 10*time.Second)
	defer client.Close()

	result, err := client.Health(context.Background())
	if err != nil {
		t.Fatalf("Health check should not error for unhealthy: %v", err)
	}

	if result.Status != "unhealthy" {
		t.Errorf("Expected status 'unhealthy', got '%s'", result.Status)
	}
}

func TestClient_Query_InvalidQuery(t *testing.T) {
	client := NewClient("http://localhost:9428", util.AuthConfig{}, 10*time.Second)
	defer client.Close()

	_, err := client.Query(context.Background(), QueryParams{
		Query: "",
	})

	if err != ErrInvalidQuery {
		t.Errorf("Expected ErrInvalidQuery, got %v", err)
	}
}

func TestClient_Query_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid query syntax"))
	}))
	defer server.Close()

	client := NewClient(server.URL, util.AuthConfig{}, 10*time.Second)
	defer client.Close()

	_, err := client.Query(context.Background(), QueryParams{
		Query: "invalid[",
	})

	if err == nil {
		t.Error("Expected error for invalid query")
	}

	if _, ok := err.(*APIError); !ok {
		t.Errorf("Expected APIError, got %T", err)
	}
}

func TestClient_GetMaxResults(t *testing.T) {
	client := NewClient("http://localhost:9428", util.AuthConfig{}, 10*time.Second)
	defer client.Close()

	if client.GetMaxResults() != 5000 {
		t.Errorf("Expected default 5000, got %d", client.GetMaxResults())
	}
}

func TestClient_WithMaxResults(t *testing.T) {
	client := NewClient("http://localhost:9428", util.AuthConfig{}, 10*time.Second, WithMaxResults(1000))
	defer client.Close()

	if client.GetMaxResults() != 1000 {
		t.Errorf("Expected 1000, got %d", client.GetMaxResults())
	}
}

func TestClient_Close(t *testing.T) {
	client := NewClient("http://localhost:9428", util.AuthConfig{}, 10*time.Second)
	client.Close()
	client.Close() // Double close should not panic
}
