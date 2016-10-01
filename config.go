package main

import (
	"github.com/koding/multiconfig"
)

type config struct {
	WebappPath  string `default:"."`
	Port        string `default:":80"`
	Mp3Location string `default:"."`
}

func getConfig() *config {
	c := &config{}
	multiconfig.MustLoadWithPath("config.toml", c)

	return c
}
