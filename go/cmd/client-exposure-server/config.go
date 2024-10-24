package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	ServerPort string `yaml:"port"`
	Db         Db     `yaml:"db"`
	Mq         Mq     `yaml:"mq"`
}

type Db struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
}

type Mq struct {
	Url   string `yaml:"url"`
	Queue Queue  `yaml:"queue"`
}

type Queue struct {
	Deals string `yaml:"deals"`
	Rates string `yaml:"rates"`
}

func GetDefaultConfig() *ConfigStruct {
	configpath := os.Getenv("CONFIGPATH")
	if configpath == "" {
		configpath = "config.yaml"
	}
	return getConfig(configpath)
}

func getConfig(filename string) *ConfigStruct {
	filename, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var config ConfigStruct
	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		panic(err)
	}

	return &config
}
