package dockerlink_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/syb-devs/dockerlink"
)

type envVars map[string]string

var getLinkTests = []struct {
	Env          envVars
	Name         string
	Port         int
	Protocol     string
	ExpectedLink *dockerlink.Link
	ExpectedErr  error
}{
	{
		envVars{},
		"MYSQL", 3306, "TCP",
		nil,
		dockerlink.ErrLinkNotDefined,
	},
	{
		envVars{
			"MYSQL_PORT_3306_TCP":       "tcp://172.17.0.5:5432",
			"MYSQL_PORT_3306_TCP_ADDR":  "172.17.0.5",
			"MYSQL_PORT_3306_TCP_PORT":  "5432",
			"MYSQL_PORT_3306_TCP_PROTO": "tcp",
		},
		"MYSQL", 3306, "",
		&dockerlink.Link{Name: "MYSQL", Protocol: "tcp", ExposedPort: 3306, Port: 5432, Address: "172.17.0.5"},
		nil,
	},
	{
		envVars{
			"MYSQL_PORT_3306_TCP": "tcp://172.17.0.5:5432",
		},
		"MYSQL", 3306, "",
		nil,
		dockerlink.ErrPortNotDefined,
	},
	{
		envVars{
			"MYSQL_PORT_3306_TCP":      "tcp://172.17.0.5:5432",
			"MYSQL_PORT_3306_TCP_PORT": "weird",
		},
		"MYSQL", 3306, "",
		nil,
		dockerlink.ErrPortNotDefined,
	},
	{
		envVars{
			"MYSQL_PORT_3306_TCP":       "tcp://172.17.0.5:5432",
			"MYSQL_PORT_3306_TCP_PORT":  "5432",
			"MYSQL_PORT_3306_TCP_PROTO": "tcp",
		},
		"MYSQL", 3306, "",
		nil,
		dockerlink.ErrAddressNotDefined,
	},
}

func TestGetLink(t *testing.T) {
	for _, test := range getLinkTests {
		setupEnv(test.Env)

		link, err := dockerlink.GetLink(test.Name, test.Port, test.Protocol)
		if !reflect.DeepEqual(link, test.ExpectedLink) {
			t.Errorf("expecting link to be %+v, but got %+v", test.ExpectedLink, link)
		}
		if !reflect.DeepEqual(err, test.ExpectedErr) {
			t.Errorf("expecting error to be %+v, but got %+v", test.ExpectedErr, err)
		}
	}
}

func setupEnv(vars map[string]string) {
	os.Clearenv()
	for key, value := range vars {
		os.Setenv(key, value)
	}
}
