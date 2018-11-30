package idserver

import (
	"fmt"
	"github.com/eyedeekay/sam-forwarder/config"
	"strconv"
)

type Option func(*IDServer) error

// Config sets the host of the SAM Bridge to use
func Config(s *i2ptunconf.Conf) func(*IDServer) error {
	return func(c *IDServer) error {
		c.Conf = s
		return nil
	}
}

// SAMHost sets the host of the SAM Bridge to use
func SAMHost(s string) func(*IDServer) error {
	return func(c *IDServer) error {
		c.host = s
		return nil
	}
}

// SAMPort sets the port of the SAM bridge to use
func SAMPort(s string) func(*IDServer) error {
	return func(c *IDServer) error {
		val, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if val > 0 && val < 65536 {
			c.port = s
			return nil
		}
		return fmt.Errorf("port is %s invalid")
	}
}

// ConfigPath sets the server config file path
func ConfigPath(s string) func(*IDServer) error {
	return func(c *IDServer) error {
		c.path = s
		return nil
	}
}

// GenerateConfig sets the path to the keys, if no keys are present, they will be generated.
func GenerateConfig(s bool) func(*IDServer) error {
	return func(c *IDServer) error {
		c.generatekeys = s
		return nil
	}
}
