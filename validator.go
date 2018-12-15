package main

import (
	"net"
	"path/filepath"
	"regexp"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

// isLoopbackIP - validates to loopback IPv4
func isLoopbackIP(fl validator.FieldLevel) bool {

	ip := net.ParseIP(fl.Field().String())
	if ip == nil {
		return false
	}

	return ip.IsLoopback()
}

// isFilepathClean - checks if path to fiel is clean
func isFilepathClean(fl validator.FieldLevel) bool {
	path := fl.Field().String()

	return strings.EqualFold(path, filepath.Clean(path))
}

// isValidDBTableName - checks for DB.Table name validity
func isValidDBTableName(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]+\.[a-zA-Z0-9]+$`).MatchString(fl.Field().String())
}
