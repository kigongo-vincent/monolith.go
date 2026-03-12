package app

import "net/http"

// App provides request/response and optional WS/cache APIs per request.
type App struct {
	Request  *Request
	Response *Response
}

// NewApp builds an App from the HTTP request and response writer.
func NewApp(r *http.Request, w http.ResponseWriter) *App {
	return &App{
		Request:  &Request{Request: r},
		Response: &Response{ResponseWriter: w},
	}
}
