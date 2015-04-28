package config

import (
	"io/ioutil"
	"net/url"

	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Title               string `yaml:"title"`
	Url                 string `yaml:"url"`
	UrlPrefix           string `yaml:"url_prefix"`
	RestApi             bool   `yaml:"rest_api"`
	TrustAuthentication string `yaml:"trust_authentication"`
	Listen              string `yaml:"listen"`
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

var c *Config

func Get() *Config {
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

	u, err := url.Parse(Get().Url)
	if err != nil {
		logrus.Fatalf("cannot parse base URL: %s", Get().Url)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		logrus.Fatalf("only schemes 'http' and 'https' are supported: %s", Get().Url)
	}
	u.Path, u.RawQuery = "", ""
	logrus.Infof("normalizing base URL to %s", u)
	Get().Url = u.String()

	if Get().TrustAuthentication != "" && Get().TrustAuthentication != "on-gateway" && Get().TrustAuthentication != "always" && Get().TrustAuthentication != "never" {
		logrus.Fatalf("setting 'trust_authentication' should be 'never', 'on-gateway' or 'always'")
	}
}
