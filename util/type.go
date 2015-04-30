package util

import (
	"encoding/xml"
	"time"

	"github.com/apognu/gocas/config"
)

type LoginRequestorData struct {
	Config   *config.Config
	Session  LoginRequestorSession
	Message  LoginRequestorMessage
	ShowForm bool
}

type LoginRequestorSession struct {
	Ticket   string
	Service  string
	Url      string
	Username string
}

type LoginRequestorMessage struct {
	Type    string
	Message string
}

type FailedLogin struct {
	Id        string    `bson:"_id"`
	Ip        string    `bson:"ip,omitempty"`
	Username  string    `bson:"username,omitempty"`
	Count     uint      `bson:"count"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type CASServiceResponse struct {
	XMLName      xml.Name `xml:"cas:serviceResponse"`
	Xmlns        string   `xml:"xmlns:cas,attr"`
	Success      *CASAuthenticationSuccess
	Failure      *CASAuthenticationFailure
	ProxySuccess *CASProxySuccess
	ProxyFailure *CASProxyFailure
}

type CASAuthenticationSuccess struct {
	XMLName xml.Name `xml:"cas:authenticationSuccess"`
	User    CASUser
	PgtIou  CASPgtIou `xml:",omitempty"`
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

type CASPgtIou struct {
	XMLName xml.Name `xml:"cas:proxyGrantingTicket"`
	Ticket  string   `xml:",chardata"`
}

type CASProxySuccess struct {
	XMLName xml.Name `xml:"cas:proxyTicket"`
	Ticket  string   `xml:",chardata"`
}

type CASProxyFailure struct {
	XMLName xml.Name `xml:"cas:proxyFailure"`
	Code    string   `xml:"string"`
	Message string   `xml:",chardata"`
}

func NewCASResponse() CASServiceResponse {
	return CASServiceResponse{
		Xmlns: "http://www.yale.edu/tp/cas",
	}
}

func NewCASSuccessResponse(u string, pgtiou string) []byte {
	s := NewCASResponse()
	s.Success = &CASAuthenticationSuccess{
		User:   CASUser{User: u},
		PgtIou: CASPgtIou{Ticket: pgtiou},
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

func NewCASProxySuccessResponse(pt string) []byte {
	s := NewCASResponse()
	s.ProxySuccess = &CASProxySuccess{
		Ticket: pt,
	}
	x, _ := xml.Marshal(s)
	return x
}

func NewCASProxyFailureResponse(c string, msg string) []byte {
	f := NewCASResponse()
	f.ProxyFailure = &CASProxyFailure{
		Code:    c,
		Message: msg,
	}
	x, _ := xml.Marshal(f)
	return x
}
