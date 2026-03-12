package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewApp(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	a := NewApp(req, rec)
	if a.Request == nil || a.Response == nil {
		t.Fatal("Request and Response should be set")
	}
	if a.Request.Context() == nil {
		t.Fatal("Context should be non-nil")
	}
}

func TestRequestAuthSetAuth(t *testing.T) {
	req := &Request{}
	if req.Auth() != nil {
		t.Fatal("Auth() should be nil initially")
	}
	req.SetAuth("user1")
	if req.Auth() != "user1" {
		t.Errorf("Auth want user1 got %v", req.Auth())
	}
}

func TestResponseEnableCache(t *testing.T) {
	rec := httptest.NewRecorder()
	resp := &Response{ResponseWriter: rec}
	resp.EnableCache(300, "redis")
	if resp.CacheTTL() != 300 {
		t.Errorf("CacheTTL want 300 got %d", resp.CacheTTL())
	}
	if resp.CacheType() != "redis" {
		t.Errorf("CacheType want redis got %s", resp.CacheType())
	}
}
