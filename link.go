package dockerlink

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var ErrLinkNotDefined = errors.New("link not defined")
var ErrPortNotDefined = errors.New("port not defined")

// Link represents a Docker container link
type Link struct {
	Name        string
	Protocol    string
	ExposedPort int
	Port        int
	Address     string
}

// GetLink returns a Link configured with Docker defined environment vars
func GetLink(name string, port int, proto string) (*Link, error) {
	if proto == "" {
		proto = "TCP"
	}
	prefix := fmt.Sprintf(
		"%s_PORT_%d_%s",
		strings.ToUpper(name),
		port,
		strings.ToUpper(proto))

	if os.Getenv(prefix) == "" {
		return nil, ErrLinkNotDefined
	}

	envPort := os.Getenv(fmt.Sprintf("%s_PORT", prefix))
	if envPort == "" {
		return nil, ErrPortNotDefined
	}
	portInt, err := strconv.Atoi(envPort)
	if err != nil {
		return nil, ErrPortNotDefined
	}

	l := &Link{
		Name:        name,
		Protocol:    strings.ToLower(proto),
		ExposedPort: port,
		Port:        portInt,
		Address:     os.Getenv(fmt.Sprintf("%s_ADDR", prefix)),
	}
	return l, nil
}
