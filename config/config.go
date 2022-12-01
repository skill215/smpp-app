package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/creasty/defaults"
	yaml "gopkg.in/yaml.v3"
)

type SmppConfig struct {
	Server struct {
		Addr     string `default:"localhost" yaml:"addr"`
		Port     uint16 `default:"5588" yaml:"port"`
		User     string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"server"`
	Client struct {
		Type  string `default:"transmitter" yaml:"bind-type"`
		Count uint16 `default:"10" yaml:"conn-num"`
	}
	Message struct {
		Send struct {
			Src struct {
				Npi   uint16 `default:"1" yaml:"npi"`
				Ton   uint16 `default:"1" yaml:"ton"`
				Oaddr string `yaml:"oaddr"`
			} `yaml:"src"`
			Dst struct {
				Npi uint16 `default:"1" yaml:"npi"`
				Ton uint16 `default:"1" yaml:"ton"`
			} `yaml:"dst"`
			Traffic   int    `yaml:"traffic"`
			RequireSR bool   `default:"false" yaml:"require-sr"`
			Content   string `yaml:"content"`
		} `yaml:"send"`
	} `yaml:"message"`
}

type AppConfig struct {
	App struct {
		SmppConn []SmppConfig `yaml:"smpp"`
		Rest     struct {
			Addr string `default:"0.0.0.0" yaml:"addr"`
			Port uint16 `default:"5000" yaml:"port"`
		} `yaml:"rest"`
		Log struct {
			Level string `default:"info" yaml:"level"`
		} `yaml:"log"`
	} `yaml:"serivce"`
}

func (c *AppConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(c)

	type plain AppConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	return nil
}

func GetSmppConf() (*AppConfig, error) {
	c := &AppConfig{}
	yamlFile, err := os.ReadFile("smpp-app.yaml")
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		return nil, err
	}
	log.Printf("%+v", c)
	return c, nil
}

func (ac *AppConfig) GetRestAddr() string {
	return fmt.Sprintf("%s:%d", ac.App.Rest.Addr, ac.App.Rest.Port)
}

func (s *SmppConfig) IsTransmitter() bool {
	return !strings.EqualFold("receiver", s.Client.Type)
}
