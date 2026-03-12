package integrations

// Storage interface: same API for local, S3, Cloudinary.
type Storage interface {
	Presign(bucket, key string, ttlSeconds int) (string, error)
	Put(bucket, key string, data []byte) error
	Get(bucket, key string) ([]byte, error)
}

// Payment interface: Pay and webhook Receive.
type Payment interface {
	Pay(details PayDetails) (PayResult, error)
	Receive(payload []byte, signature string) (ReceiveResult, error)
}

// PayDetails for payment requests.
type PayDetails struct {
	Amount   int64
	Currency string
	Email    string
	Ref      string
}

// PayResult from payment provider.
type PayResult struct {
	ID     string
	URL    string
	Status string
}

// ReceiveResult from webhook.
type ReceiveResult struct {
	ID     string
	Status string
}

// Mail interface for transactional email.
type Mail interface {
	Send(to, template string, data map[string]string) error
}

// Maps interface for geocoding etc.
type Maps interface {
	Geocode(address string) (lat, lng float64, err error)
}
