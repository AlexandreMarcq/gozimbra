package auth_test

import (
	"fmt"
	"testing"

	"github.com/AlexandreMarcq/gozimbra/internal/auth"
	"github.com/AlexandreMarcq/gozimbra/test"
)

func TestAuthRequest(t *testing.T) {
	testUser := "USER"
	testPass := "PASS"
	want := fmt.Sprintf(`<AuthRequest xmlns="urn:zimbraAdmin" name="%s" password="%s"/>`,
		testUser, testPass)
	got := auth.NewAuthRequest(testUser, testPass).ToXML()

	test.AssertXML(t, got, want)
}
