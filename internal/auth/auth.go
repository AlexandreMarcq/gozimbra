package auth

import "fmt"

type authRequest struct {
	username string
	password string
}

func NewAuthRequest(username, password string) *authRequest {
	return &authRequest{username, password}
}

func (r *authRequest) ToXML() string {
	return fmt.Sprintf(`<AuthRequest xmlns="urn:zimbraAdmin" name="%s" password="%s"/>`, r.username, r.password)
}
