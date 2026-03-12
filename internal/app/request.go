package app

import (
	"context"
	"net/http"
)

// Request wraps the HTTP request and auth/role config.
type Request struct {
	*http.Request
	AllowedRoles    []string
	AllowedUserIDs  []string
	resolvedAuth    any
}

// Context returns the request context.
func (r *Request) Context() context.Context {
	if r.Request == nil {
		return context.Background()
	}
	return r.Request.Context()
}

// Auth returns the current authenticated user (nil if unauthenticated).
func (r *Request) Auth() any {
	return r.resolvedAuth
}

// SetAuth sets the resolved auth user (used by middleware).
func (r *Request) SetAuth(user any) {
	r.resolvedAuth = user
}
