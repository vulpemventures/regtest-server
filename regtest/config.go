package main

import (
	"errors"
	"math"
	"os"
	"regexp"
	"strconv"
)

var defaultEnv = map[string]string{
	"PORT":    "8000",
	"ADDRESS": "localhost",
}

type Config struct {
	Port    string
	Address string
}

func generateConfigFromEnv() (Config, error) {
	address := getStrEnv("ADDRESS")
	port, err := getIntEnv("PORT")
	if err != nil {
		return Config{}, err
	}

	// checks
	if !validPort(port) {
		return Config{}, errors.New("Invalid port: port must be a number betweem 1 and 65535")
	}
	if !validAddress(address) {
		return Config{}, errors.New("Invalid address: address must be 'localhost' or a valid IP address")
	}

	return Config{
		Port:    strconv.Itoa(port),
		Address: address,
	}, nil
}

// get env vars
func getStrEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultEnv[key]
}

func getIntEnv(key string) (int, error) {
	return strconv.Atoi(getStrEnv(key))
}

// check env vars
func validPort(val int) bool {
	return val > 0 && val < math.MaxUint16
}
func validAddress(val string) bool {
	return val == "localhost" || regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]`).MatchString(val)
}
