package app

import "net/http"

// Response wraps the HTTP response writer and cache config.
type Response struct {
	http.ResponseWriter
	cacheTTL  int
	cacheType string
}

// EnableCache enables response caching with TTL in seconds and cache type from settings.
func (r *Response) EnableCache(ttlSeconds int, cacheType string) {
	r.cacheTTL = ttlSeconds
	r.cacheType = cacheType
}

// CacheTTL returns the cache TTL in seconds (0 if disabled).
func (r *Response) CacheTTL() int { return r.cacheTTL }

// CacheType returns the cache backend type.
func (r *Response) CacheType() string { return r.cacheType }
