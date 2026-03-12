package integrations


// Integrations holds provider implementations (Storage, Payment, Mail, Maps).
type Integrations struct {
	storage Storage
	payment Payment
	mail    Mail
	maps    Maps
}

// New builds Integrations from config (provider choice from settings).
func New(storage Storage, payment Payment, mail Mail, maps Maps) *Integrations {
	return &Integrations{storage: storage, payment: payment, mail: mail, maps: maps}
}

// Storage returns the configured storage provider.
func (i *Integrations) Storage() Storage {
	if i.storage != nil {
		return i.storage
	}
	return &noopStorage{}
}

// Payment returns the configured payment provider.
func (i *Integrations) Payment() Payment {
	if i.payment != nil {
		return i.payment
	}
	return &noopPayment{}
}

// Mail returns the configured mail provider.
func (i *Integrations) Mail() Mail {
	if i.mail != nil {
		return i.mail
	}
	return &noopMail{}
}

// Maps returns the configured maps provider.
func (i *Integrations) Maps() Maps {
	if i.maps != nil {
		return i.maps
	}
	return &noopMaps{}
}

type noopStorage struct{}

func (n *noopStorage) Presign(_, _ string, _ int) (string, error) { return "", nil }
func (n *noopStorage) Put(_, _ string, _ []byte) error             { return nil }
func (n *noopStorage) Get(_, _ string) ([]byte, error)             { return nil, nil }

type noopPayment struct{}

func (n *noopPayment) Pay(PayDetails) (PayResult, error) { return PayResult{}, nil }
func (n *noopPayment) Receive([]byte, string) (ReceiveResult, error) {
	return ReceiveResult{}, nil
}

type noopMail struct{}

func (n *noopMail) Send(_, _ string, _ map[string]string) error { return nil }

type noopMaps struct{}

func (n *noopMaps) Geocode(string) (float64, float64, error) { return 0, 0, nil }
