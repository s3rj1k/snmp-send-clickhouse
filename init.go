package main

import (
	"flag"
	"log"
	"os"

	validator "gopkg.in/go-playground/validator.v9"
)

// nolint: gochecknoglobals
var (
	// Config object
	c *Config

	// InfoLogger defines info level logger object
	InfoLogger *log.Logger
	// ErrorLogger defines error level logger object
	ErrorLogger *log.Logger

	// Validate defines validator object
	Validate *validator.Validate
)

// nolint: gochecknoinits
func init() {

	var err error

	// CMD Config Path
	var cmdConfigPath string

	// configure loggers
	InfoLogger = log.New(os.Stdout, "INF: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERR: ", log.Ldate|log.Ltime|log.Lshortfile)

	flag.StringVar(&cmdConfigPath, "config", "config.json", "path to config file")
	flag.Parse()

	// initialize validator object
	Validate = validator.New()

	// register custom validation functions
	err = Validate.RegisterValidation("loopback", isLoopbackIP)
	if err != nil {
		ErrorLogger.Fatalf("validator error: %v", err)
	}
	err = Validate.RegisterValidation("clean_filepath", isFilepathClean)
	if err != nil {
		ErrorLogger.Fatalf("validator error: %v", err)
	}
	err = Validate.RegisterValidation("table_name", isValidDBTableName)
	if err != nil {
		ErrorLogger.Fatalf("validator error: %v", err)
	}

	// get config data
	c, err = GetConfig(cmdConfigPath)
	if err != nil {
		ErrorLogger.Fatal(err)
	}
}
