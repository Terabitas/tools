package services

import . "gopkg.in/check.v1"

type TypeSuite struct{}

var _ = Suite(&TypeSuite{})

func (s *TypeSuite) TestIfSecretsAreLoadedCorrectly(c *C) {

	secrets := map[string]interface{}{
		"empty": map[string]interface{}{
			"prod-db-cred": map[string]interface{}{
				"username": "xxx",
				"password": "ccc",
			},
			"prod-token-paypal": map[string]interface{}{
				"token": "xxx",
			},
		},
		"does-not-exists": map[string]interface{}{
			"prod-db-cred": map[string]interface{}{
				"username": "xxx",
				"password": "ccc",
			},
			"prod-token-paypal": map[string]interface{}{
				"token": "xxx",
			},
		},
	}

	srvcs := Services{
		&Service{
			Name:    "empty",
			Secrets: Secrets{},
			EnvVars: EnvVars{},
		},
		&Service{
			Name:    "none",
			Secrets: Secrets{},
			EnvVars: EnvVars{},
		},
	}

	srvcs.LoadSecrets(secrets)
	c.Assert(len(srvcs.find("empty").Secrets), Equals, 2)
	c.Assert(
		srvcs.find("empty").Secrets.find("prod-db-cred").Data,
		DeepEquals,
		map[string]string{"username": "xxx", "password": "ccc"},
	)
	c.Assert(
		srvcs.find("empty").Secrets.find("prod-token-paypal").Data,
		DeepEquals,
		map[string]string{"token": "xxx"},
	)
	c.Assert(len(srvcs.find("none").Secrets), Equals, 0)
	c.Assert(srvcs.find("not-exists"), IsNil)
}

func (s *TypeSuite) TestIfEnvVarsAreLoadedCorrectly(c *C) {

	envs := map[string]interface{}{
		"empty": map[string]interface{}{
			"CONF_STRATEGY": map[string]interface{}{
				"live":    "xxx",
				"staging": "ccc",
				"testing": "vvv",
			},
			"CONF_STRATEGY_X": map[string]interface{}{
				"live":    "xxx",
				"staging": "ccc",
				"testing": "vvv",
			},
		},
		"does-not-exists": map[string]interface{}{
			"CONF_STRATEGY": map[string]interface{}{
				"live":    "xxx",
				"staging": "ccc",
				"testing": "vvv",
			},
			"CONF_STRATEGY_X": map[string]interface{}{
				"live":    "xxx",
				"staging": "ccc",
				"testing": "vvv",
			},
		},
	}

	srvcs := Services{
		&Service{
			Name:    "empty",
			Secrets: Secrets{},
			EnvVars: EnvVars{},
		},
		&Service{
			Name:    "none",
			Secrets: Secrets{},
			EnvVars: EnvVars{},
		},
	}

	srvcs.LoadEnvVars(envs)
	c.Assert(len(srvcs.find("empty").EnvVars), Equals, 2)
	c.Assert(
		srvcs.find("empty").EnvVars.find("CONF_STRATEGY").Values,
		DeepEquals,
		map[Environment]string{
			"live":    "xxx",
			"staging": "ccc",
			"testing": "vvv",
		},
	)
	c.Assert(
		srvcs.find("empty").EnvVars.find("CONF_STRATEGY_X").Values,
		DeepEquals,
		map[Environment]string{
			"live":    "xxx",
			"staging": "ccc",
			"testing": "vvv",
		},
	)
	c.Assert(len(srvcs.find("none").EnvVars), Equals, 0)
}
