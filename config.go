package main

import (
	"os"

	ini "gopkg.in/ini.v1"
)

type ctfT struct {
	title string
}

type configT struct {
	ip          string
	port        int
	useTLS      bool
	certificate *os.File
	keyFile     *os.File
	db          *ctfDB
	ctfPrefs    *ctfT
}

func loadConfig() (*configT, error) {
	config := &configT{}
	cfg, err := ini.Load("config.ini")
	if err != nil {
		return nil, err
	}

	config.ip = cfg.Section("server").Key("ip").String()
	config.port, _ = cfg.Section("server").Key("port").Int()
	config.useTLS, _ = cfg.Section("server").Key("use_tls").Bool()
	if config.useTLS {
		config.certificate, err = os.Open(cfg.Section("server").Key("certificate").String())
		if err != nil {
			return nil, err
		}
		config.keyFile, err = os.Open(cfg.Section("server").Key("key_file").String())
		if err != nil {
			return nil, err
		}
	}
	config.ctfPrefs = &ctfT{}
	config.ctfPrefs.title = cfg.Section("ctf").Key("title").String()

	return config, nil
}
