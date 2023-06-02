package account

import (
	"fmt"
	"strings"

	"github.com/AlexandreMarcq/gozimbra/internal/utils"
)

type createAccountRequest struct {
	account  string
	password string
}

func NewCreateAccountRequest(account, password string) *createAccountRequest {
	return &createAccountRequest{account, password}
}

func (r *createAccountRequest) ToXML() string {
	return fmt.Sprintf(`<CreateAccountRequest xmlns="urn:zimbraAdmin" name="%s" password="%s"></CreateAccountRequest>`, r.account, r.password)
}

type getAccountRequest struct {
	account    string
	attributes []string
}

func NewGetAccountRequest(account string, attributes []string) *getAccountRequest {
	return &getAccountRequest{account, attributes}
}

func (r *getAccountRequest) ToXML() string {
	return fmt.Sprintf(`<GetAccountRequest xmlns="urn:zimbraAdmin" attrs="%s"><account by="name">%s</account></GetAccountRequest>`, strings.Join(r.attributes, ","), r.account)
}

type modifyAccountRequest struct {
	id         string
	attributes utils.AttrsMap
}

func NewModifyAccountRequest(id string, attributes utils.AttrsMap) *modifyAccountRequest {
	return &modifyAccountRequest{id, attributes}
}

func (r *modifyAccountRequest) ToXML() string {
	var sb strings.Builder

	for _, a := range r.attributes.Keys() {
		sb.WriteString(fmt.Sprintf("<a n=\"%s\">%s</a>", a, r.attributes[a]))
	}
	return fmt.Sprintf("<ModifyAccountRequest xmlns=\"urn:zimbraAdmin\" id=\"%s\">%s</ModifyAccountRequest>", r.id, sb.String())
}
