package servicehandler

import (
	"testing"
)

func TestInvalidStatusCode(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Fail to panic on invalid status code")
		}
	}()
	RaiseHTTPException(1000, "testError")
}

func TestRaiseException(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("Fail to raise http exception")
		}
		if err.(HTTPException).StatusCode != 500 {
			t.Errorf("Invalid status code for internal server error")
		}
	}()
	RaiseHTTPException(500, "testError")
}
