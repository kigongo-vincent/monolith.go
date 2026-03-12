package result

// Result is the return type for handlers. Use Ok or Err to build it.
type Result struct {
	status  int
	body    any
	isError bool
}

// Ok returns a successful result. String → text/plain, struct/slice → JSON.
func Ok(value any) Result {
	return Result{status: 200, body: value, isError: false}
}

// Err returns an error result with HTTP status and message.
func Err(statusCode int, message string) Result {
	return Result{status: statusCode, body: message, isError: true}
}

// Status returns the HTTP status code.
func (r Result) Status() int { return r.status }

// Body returns the response body (any for JSON/text).
func (r Result) Body() any { return r.body }

// IsError reports whether this result is an error.
func (r Result) IsError() bool { return r.isError }
