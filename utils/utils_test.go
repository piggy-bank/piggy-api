package utils

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func DummyTest(t *testing.T) {
	scope := "dev"
	want := CalculateProfile()
	if scope != want {
		t.Fatalf(`This is a dummy test`)
	}
	println("Hello from GitHub Actions")
}
