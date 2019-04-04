package main

import (
	"io/ioutil"
	"os"
	"time"

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
	pages       []string
	startTime   int64
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

	config.startTime = time.Now().Unix()
	pages, err := ioutil.ReadDir("./pages")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		pName := page.Name()[:len(page.Name())-5]
		config.pages = append(config.pages, pName)
	}

	return config, nil
}
