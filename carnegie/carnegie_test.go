package carnegie

import (
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	config := viper.New()
	_, err := New(config)
	if err == nil {
		t.Fatalf("error should bubble up from NewCache")
	}
	config.Set("backend", "consul")
	carnegie, err := New(config)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if carnegie.CacheInterval != 60*time.Second {
		t.Fatalf("default time should be set to 60s")
	}
	if carnegie.Server.Addr != ":8181" {
		t.Fatalf("default address should be set to :8181")
	}
	config.Set("interval", "3m0s")
	carnegie, err = New(config)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if carnegie.CacheInterval != 3*time.Minute {
		t.Fatalf("default time should be set to 60s")
	}
}

func TestStart(t *testing.T) {
	config := viper.New()
	config.Set("backend", "consul")
	carnegie, err := New(config)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	go carnegie.Start()
}

func TestHandler(t *testing.T) {
	config := viper.New()
	config.Set("backend", "consul")
	carnegie, err := New(config)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// empty host
	req, err := http.NewRequest("GET", "http://localhost:8181", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	w := httptest.NewRecorder()
	carnegie.Handler(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("should return 404 on no backend")
	}
	req.Host = "test.com"
	w = httptest.NewRecorder()
	carnegie.Handler(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("should return 200")
	}
}
