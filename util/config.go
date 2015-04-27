package util

import (
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Mongo struct {
		Host string `yaml:"host"`
	} `yaml:"mongo"`
	Services       []string `yaml:"services"`
	TicketValidity struct {
		LoginTicket          int `yaml:"login_ticket"`
		TicketGrantingTicket int `yaml:"ticket_granting_ticket"`
		ServiceTicket        int `yaml:"service_ticket"`
	} `yaml:"ticket_validity"`
	Authenticator string `yaml:"authenticator"`
}

var c Config

func GetConfig() Config {
	return c
}

func SetConfig(p string) {
	f, err := ioutil.ReadFile(p)
	if err != nil {
		logrus.Fatalf("error parsing configuration file: %s", err)
	}
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		logrus.Fatalf("error parsing configuration file: %s", err)
	}
}
