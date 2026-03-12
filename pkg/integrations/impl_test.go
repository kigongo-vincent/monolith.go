package integrations

import "testing"

func TestNewNilReturnsNoops(t *testing.T) {
	i := New(nil, nil, nil, nil)
	if i.Storage() == nil {
		t.Fatal("Storage() should not be nil")
	}
	if i.Payment() == nil {
		t.Fatal("Payment() should not be nil")
	}
	if i.Mail() == nil {
		t.Fatal("Mail() should not be nil")
	}
	if i.Maps() == nil {
		t.Fatal("Maps() should not be nil")
	}
}

func TestNoopStorage(t *testing.T) {
	s := &noopStorage{}
	url, err := s.Presign("b", "k", 60)
	if err != nil || url != "" {
		t.Errorf("Presign want \"\",nil got %q,%v", url, err)
	}
	if err := s.Put("b", "k", []byte("x")); err != nil {
		t.Error(err)
	}
	data, err := s.Get("b", "k")
	if err != nil || data != nil {
		t.Errorf("Get want nil,nil got %v,%v", data, err)
	}
}
