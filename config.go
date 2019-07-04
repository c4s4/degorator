package main

import (
	"fmt"
	"regexp"
)

var config Config

type Config struct {
	Port       int         `yaml:"port"`
	Users      []User      `yaml:"users"`
	Operations []Operation `yaml:"operations"`
	Target     *Target     `yaml:"target"`
}

type User struct {
	Name  string `yaml:"name"`
	Token string `yaml:"token"`
}

type Operation struct {
	Path       string                `yaml:"path"`
	Method     string                `yaml:"method"`
	Parameters map[string]*Parameter `yaml:"parameters"`
	Target     *Target               `yaml:"target"`
}

type Parameter struct {
	Optional bool   `yaml:"optional"`
	Regexp   string `yaml:"regexp`
	Compiled *regexp.Regexp
}

type Target struct {
	Host string `yaml:"host"`
	Path string `yaml:"path"`
}

func (config *Config) Compile() error {
	for i := range config.Operations {
		for name := range config.Operations[i].Parameters {
			compiled, err := regexp.Compile(config.Operations[i].Parameters[name].Regexp)
			if err != nil {
				return fmt.Errorf("Regexp '%s' doesn't compile: %v",
					config.Operations[i].Parameters[name].Regexp, err)
			}
			config.Operations[i].Parameters[name].Compiled = compiled
		}
	}
	return nil
}
