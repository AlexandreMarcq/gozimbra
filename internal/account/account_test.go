package account_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/AlexandreMarcq/gozimbra/internal/account"
	"github.com/AlexandreMarcq/gozimbra/internal/utils"
	"github.com/AlexandreMarcq/gozimbra/test"
)

func TestCreateAccountRequest(t *testing.T) {
	testAccount := "ACCOUNT"
	testPassword := "PASSWORD"
	want := fmt.Sprintf(`<CreateAccountRequest xmlns="urn:zimbraAdmin" name="%s" password="%s"></CreateAccountRequest>`, testAccount, testPassword)
	got := account.NewCreateAccountRequest(testAccount, testPassword).ToXML()

	test.AssertXML(t, got, want)
}

func TestGetAccountRequest(t *testing.T) {
	testAccount := "ACCOUNT"
	testAttributes := []string{"ATTR1", "ATTR2"}
	want := fmt.Sprintf(`<GetAccountRequest xmlns="urn:zimbraAdmin" attrs="%s"><account by="name">%s</account></GetAccountRequest>`, strings.Join(testAttributes, ","), testAccount)
	got := account.NewGetAccountRequest(testAccount, testAttributes).ToXML()

	test.AssertXML(t, got, want)
}

func TestModifyAccountRequest(t *testing.T) {
	testId := "ID"
	testAttributes := make(utils.AttrsMap)
	testAttributes["ATTR1"] = "value1"
	testAttributes["ATTR2"] = "value2"
	want := fmt.Sprintf(`<ModifyAccountRequest xmlns="urn:zimbraAdmin" id="%s"><a n="%s">%s</a><a n="%s">%s</a></ModifyAccountRequest>`, testId, "ATTR1", testAttributes["ATTR1"], "ATTR2", testAttributes["ATTR2"])
	got := account.NewModifyAccountRequest(testId, testAttributes).ToXML()

	test.AssertXML(t, got, want)
}

func TestSetPasswordRequest(t *testing.T) {
	testId := "ID"
	testPassword := "TEST_PASSWORD"
	want := fmt.Sprintf(`<SetPasswordRequest xmlns="urn:zimbraAdmin" id="%s" newPassword="%s"/>`, testId, testPassword)
	got := account.NewSetPasswordRequest(testId, testPassword).ToXML()

	test.AssertXML(t, got, want)
}
