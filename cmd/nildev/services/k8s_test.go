package services

import (
	"os"
	"testing"

	"io/ioutil"

	. "gopkg.in/check.v1"
)

type K8sSuite struct{}

var _ = Suite(&K8sSuite{})
var (
	servicesContent = `[
  {
    "image":"nildev/empty",
    "name":"empty",
    "isPublic": true,
    "subscribe":["xxx"]
  },
  {
    "image":"nildev/account",
    "name":"account",
    "isPublic": true,
    "subscribe":[]
  }
]
`
	secretsContent = `{
  "empty":{
    "prod-db-cred" :{
      "username":"aaa",
      "password":"bbb"
    },
    "prod-paypal-token":{
      "token":"aaa"
    },
    "service-cfg":{
      "cfg":"{\"key\":\"value\"} | base64"
    }
  }
}`
	envVarsContent = `{
  "empty": {
    "CONF_STRATEGY": {
      "live": "binary_search",
      "staging": "random",
      "testing": "fake"
    }
  }
}`
)

const (
	pathToServices = "./services.json"
	pathToSecrets  = "./secrets.json"
	pathToEnv      = "./env.json"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	destroy()
	os.Exit(code)
}

func setup() {
	ioutil.WriteFile(pathToServices, []byte(servicesContent), 0644)
	ioutil.WriteFile(pathToSecrets, []byte(secretsContent), 0644)
	ioutil.WriteFile(pathToEnv, []byte(envVarsContent), 0644)
}

func destroy() {
	os.Remove(pathToServices)
	os.Remove(pathToSecrets)
	os.Remove(pathToEnv)
}

func (s *K8sSuite) TestIfK8sPlatformIsCreated(c *C) {
	var x Platform
	plt, err := MakeKubernetesPlatform("./build", "live", pathToServices, pathToSecrets, pathToEnv)

	rez := &KubernetesPlatform{
		Services: Services{
			&Service{
				Image:     "nildev/empty",
				Name:      "empty",
				IsPublic:  true,
				Subscribe: []Name{"xxx"},
				Secrets: Secrets{
					&Secret{},
					&Secret{},
					&Secret{},
				},
				EnvVars: EnvVars{
					&EnvVar{},
				},
			},
			&Service{
				Image:     "nildev/account",
				Name:      "account",
				IsPublic:  true,
				Subscribe: []Name{},
				Secrets:   Secrets{},
				EnvVars:   EnvVars{},
			},
		},
	}

	c.Assert(plt, Implements, &x)
	c.Assert(err, IsNil)
	c.Assert(plt.Services.find("empty").Image, Equals, rez.Services.find("empty").Image)
	c.Assert(plt.Services.find("empty").IsPublic, Equals, rez.Services.find("empty").IsPublic)
	c.Assert(plt.Services.find("empty").Subscribe, DeepEquals, rez.Services.find("empty").Subscribe)
	c.Assert(len(plt.Services.find("empty").Secrets), Equals, len(rez.Services.find("empty").Secrets))
	c.Assert(len(plt.Services.find("empty").EnvVars), Equals, len(rez.Services.find("empty").EnvVars))

	c.Assert(plt.Services.find("account").Image, Equals, rez.Services.find("account").Image)
	c.Assert(plt.Services.find("account").IsPublic, Equals, rez.Services.find("account").IsPublic)
	c.Assert(plt.Services.find("account").Subscribe, DeepEquals, rez.Services.find("account").Subscribe)
	c.Assert(len(plt.Services.find("account").Secrets), Equals, len(rez.Services.find("account").Secrets))
	c.Assert(len(plt.Services.find("account").EnvVars), Equals, len(rez.Services.find("account").EnvVars))
}

func (s *K8sSuite) TestIfPlatformIsCreated(c *C) {
	plt, err := MakePlatform("kubernetes", "./build", "live", pathToServices, pathToSecrets, pathToEnv)

	c.Assert(plt, NotNil)
	c.Assert(err, IsNil)
	c.Assert(len(plt.(*KubernetesPlatform).Services), Equals, 2)
}
