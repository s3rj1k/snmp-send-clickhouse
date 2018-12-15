package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// MergedTableMapKey is used for mapping SNMP data to network box
type MergedTableMapKey struct {
	Community string `json:"community" validate:"required,printascii"`
	Address   string `json:"address" validate:"required,printascii"`
}

// Config main application configuration
type Config struct {
	DB struct {
		Local struct {
			Path string `json:"path" validate:"required,printascii,min=4,clean_filepath"`
		} `json:"local" validate:"required"`

		Clickhouse struct {
			Address string `json:"address" validate:"required,ipv4,loopback"`
			Port    int    `json:"port" validate:"required,min=1025,max=65535"`
			Table   string `json:"table" validate:"required,printascii,table_name"`
		} `json:"clickhouse" validate:"required"`
	} `json:"db" validate:"required"`

	SNMP []MergedTableMapKey `json:"snmp" validate:"required,unique"`
}

// GetConfig - get application configuration
func GetConfig(path string) (*Config, error) {

	var err error

	// create config object
	c = new(Config)

	// read config file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config error: %s", err)
	}

	// convert config file to object
	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, fmt.Errorf("config error: %s", err)
	}

	// additional structe validation
	err = Validate.Struct(c)
	if err != nil {
		return nil, fmt.Errorf("config error: %s", err)
	}

	return c, nil
}
