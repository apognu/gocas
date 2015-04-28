package config

import (
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Title               string `yaml:"title"`
	UrlPrefix           string `yaml:"url_prefix"`
	TrustAuthentication bool   `yaml:"trust_authentication"`
	Mongo               struct {
		Host string `yaml:"host"`
	} `yaml:"mongo"`
	Services       []string `yaml:"services"`
	TicketValidity struct {
		LoginTicket          int `yaml:"login_ticket"`
		TicketGrantingTicket int `yaml:"ticket_granting_ticket"`
		ServiceTicket        int `yaml:"service_ticket"`
	} `yaml:"ticket_validity"`
	Protocol      string `yaml:"protocol"`
	Authenticator string `yaml:"authenticator"`
	Ldap          struct {
		Host string `yaml:"host"`
		Base string `yaml:"base"`
		Dn   string `yaml:"dn"`
	} `yaml:"ldap"`
	Oauth struct {
		ClientID    string   `yaml:"client_id"`
		Secret      string   `yaml:"secret"`
		AuthURL     string   `yaml:"auth_url"`
		TokenURL    string   `yaml:"token_url"`
		RedirectURL string   `yaml:"redirect_url"`
		Scopes      []string `yaml:"scopes"`
		UserinfoURL string   `yaml:"userinfo_url"`
	} `yaml:"oauth"`
}

var c Config

func Get() Config {
	return c
}

func Set(p string) {
	f, err := ioutil.ReadFile(p)
	if err != nil {
		logrus.Fatalf("error parsing configuration file: %s", err)
	}
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		logrus.Fatalf("error parsing configuration file: %s", err)
	}
}
