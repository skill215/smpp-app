package main

import (
	"os"

	"github.com/creasty/defaults"
	yaml "gopkg.in/yaml.v3"
)

type SmppConfig struct {
	SmppApp struct {
		SmppServer struct {
			Addr     string `default:"localhost" yaml:"addr"`
			Port     uint16 `default:"8443" yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
		} `yaml:"smpp-server"`
		SmppClient struct {
			Type string `default:"receiver" yaml:"bind-type"`
		} `yaml:"smpp-client"`
		SmppMessage struct {
			Src struct {
				Npi uint8 `default:"1" yaml:"npi"`
				Ton uint8 `default:"1" yaml:"ton"`
			} `yaml:"src"`
			Dst struct {
				Npi uint8 `default:"1" yaml:"npi"`
				Ton uint8 `default:"1" yaml:"ton"`
			} `yaml:"dst"`
			Content string `yaml:"content"`
		} `yaml:"smpp-ao"`
		Service struct {
			Addr string `default:"0.0.0.0" yaml:"addr"`
			Port uint16 `default:"8080" yaml:"port"`
		} `yaml:"service"`
		Log struct {
			Level string `default:"info" yaml:"level"`
		} `yaml:"log"`
	} `yaml:"smpp-app"`
}

func (c *SmppConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(c)

	type plain SmppConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	return nil
}

func GetSmppConf() (*SmppConfig, error) {
	c := &SmppConfig{}
	yamlFile, err := os.ReadFile("smpp-app.yaml")
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		return nil, err
	}
	return c, nil
}
