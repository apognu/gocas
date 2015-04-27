package util

import (
	"encoding/xml"

	"github.com/apognu/gocas/authenticator"
)

type LoginRequestorData struct {
	Ticket   string
	Service  string
	Type     string
	Message  string
	Username string
	ShowForm bool
	Url      string
	Logout   string
}

type CASServiceResponse struct {
	XMLName xml.Name `xml:"cas:serviceResponse"`
	Xmlns   string   `xml:"xmlns:cas,attr"`
	Success *CASAuthenticationSuccess
	Failure *CASAuthenticationFailure
}

type CASAuthenticationSuccess struct {
	XMLName xml.Name `xml:"cas:authenticationSuccess"`
	User    CASUser
}

type CASAuthenticationFailure struct {
	XMLName xml.Name `xml:"cas:authenticationFailure"`
	Code    string   `xml:"code,attr"`
	Message string   `xml:",chardata"`
}

type CASUser struct {
	XMLName xml.Name `xml:"cas:user"`
	User    string   `xml:",chardata"`
}

func NewCASResponse() CASServiceResponse {
	return CASServiceResponse{
		Xmlns: "http://www.yale.edu/tp/cas",
	}
}

func NewCASSuccessResponse(u string) []byte {
	s := NewCASResponse()
	s.Success = &CASAuthenticationSuccess{
		User: CASUser{User: u},
	}
	x, _ := xml.Marshal(s)
	return x
}

func NewCASFailureResponse(c string, msg string) []byte {
	f := NewCASResponse()
	f.Failure = &CASAuthenticationFailure{
		Code:    c,
		Message: msg,
	}
	x, _ := xml.Marshal(f)
	return x
}

type Authenticator interface {
	Auth(u string, p string) bool
}

var AvailableAuthenticators = map[string]Authenticator{
	"dummy": authenticator.Dummy{},
	"ldap":  authenticator.Ldap{},
}
