package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/creasty/defaults"
	yaml "gopkg.in/yaml.v3"
)

type AddrConfig struct {
	Prefix       string `yaml:"prefix"`
	Suffix       string `yaml:"suffix"`
	GenerateLen  int    `default:"6" yaml:"generate-length"`
	GenerateType string `default:"sequence" yaml:"generate-type"`
	Start        int    `default:"0" yaml:"start"`
	Stop         int    `yaml:"stop"`
}

type MessageConfig struct {
	Send struct {
		TextFile               string  `yaml:"text-file"`
		UrlFile                string  `yaml:"url-file"`
		ContentMode            string  `default:"random" yaml:"content-mode"`
		PreDefinedContentRatio float64 `default:"0.5" yaml:"pre-defined-content-ratio"`
		Src                    struct {
			Npi   uint16 `default:"1" yaml:"npi"`
			Ton   uint16 `default:"1" yaml:"ton"`
			Oaddr string `yaml:"oaddr"`
		} `yaml:"src"`
		Dst struct {
			Npi   uint16     `default:"1" yaml:"npi"`
			Ton   uint16     `default:"1" yaml:"ton"`
			Daddr AddrConfig `yaml:"daddr"`
		} `yaml:"dst"`
		RequireSR bool   `default:"false" yaml:"require-sr"`
		Content   string `yaml:"content"`
		Dcs       int    `default:"0" yaml:"dcs"`
	} `yaml:"send"`
}

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
	Message MessageConfig `yaml:"message"`
}

type AppConfig struct {
	App struct {
		SmppConn []SmppConfig `yaml:"smpp"`
		Rest     struct {
			Addr string `default:"0.0.0.0" yaml:"addr"`
			Port uint16 `default:"8080" yaml:"port"`
		} `yaml:"rest"`
		Log struct {
			Level string `default:"info" yaml:"level"`
		} `yaml:"log"`
	} `yaml:"service"`
}

func (c *AppConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(c)

	type plain AppConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	return nil
}

func GetSmppConf(path string) (*AppConfig, error) {
	c := &AppConfig{}
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("read conf err %v\n", err)
		return nil, err
	}
	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		fmt.Printf("parse conf err %v\n", err)
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
