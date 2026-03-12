package result

import "testing"

func TestOk(t *testing.T) {
	r := Ok("ok")
	if r.Status() != 200 || r.IsError() {
		t.Fatal("Ok should have status 200 and not be error")
	}
	if r.Body() != "ok" {
		t.Fatal("body should be ok")
	}
}

func TestErr(t *testing.T) {
	r := Err(401, "Unauthorized")
	if r.Status() != 401 || !r.IsError() {
		t.Fatal("Err should have status 401 and be error")
	}
	if r.Body().(string) != "Unauthorized" {
		t.Fatal("body should be Unauthorized")
	}
}
