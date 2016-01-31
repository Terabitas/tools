package main

import (
	"os"
	"strings"

	"path/filepath"

	"github.com/nildev/tools/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/nildev/tools/Godeps/_workspace/src/github.com/codeskyblue/go-sh"
	"github.com/nildev/tools/Godeps/_workspace/src/github.com/nildev/lib/log"
	"github.com/nildev/tools/Godeps/_workspace/src/github.com/nildev/project/setup"
	"github.com/nildev/tools/cmd/nildev/auth"
	"github.com/nildev/tools/cmd/nildev/inout"
	"github.com/nildev/tools/cmd/nildev/routes"
	"github.com/nildev/tools/cmd/nildev/state"
)

func main() {
	var ndState state.State

	app := cli.NewApp()
	app.Name = "nildev"
	app.Usage = "Tool for code generation"
	app.EnableBashCompletion = true
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "verbosity",
			Value: 1,
			Usage: "logging level",
		},
	}
	app.Before = func(c *cli.Context) error {
		ndState = state.Load()
		// setup logging here
		return nil
	}

	app.Commands = []cli.Command{
		{
			// go run nildev.go auth --provider=bitbucket.org
			Name:  "auth",
			Usage: "authenticate yourself by using one of supported providers: bitbucket|github",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "provider",
					Value: "bitbucket.org",
					Usage: "provider: bitbucket.org",
				},
			},
			Action: func(c *cli.Context) {
				ndState.Token = string(auth.Auth(c.String("provider")))
				ndState.Provider = c.String("provider")

				state.Persist(ndState)
			},
		},
		{
			// go run nildev.go create --config config.json github.com/nildev/echo
			Name:  "create",
			Usage: "create new project in $GOPATH/src/ based on selected template",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "template",
					Value: "git@github.com:nildev/prj-tpl-basic-api.git",
					Usage: "path to project template",
				},
				cli.StringFlag{
					Name:  "config",
					Value: "config.json",
					Usage: "path to json file with values to be used while rendering project template",
				},
			},
			Action: func(c *cli.Context) {
				pathToGoSrc := os.Getenv("GOPATH")
				if pathToGoSrc == "" {
					log.Fatalf("$GOPATH env is empty ? Do you have golang development environment setup ?")
				}

				pathToGoSrc += "/src"
				if len(c.Args()) == 0 {
					log.Fatalf("Please provide path to project, for example: github.com/your_username/repo_name|bitbucket.org/your_username/repo_name")
				}

				if c.String("config") == "" {
					log.Fatalf("Please provide path to config")
				}

				if c.String("template") == "" {
					log.Fatalf("Please provide path to project template git repository, for example git@github.com:nildev/prj-tpl-basic-api.git")
				}

				init := setup.NewInitializer()
				cfgLoader := setup.NewConfigLoader(c.String("config"))
				cfg := cfgLoader.Read(c.String("config"))

				destDir := pathToGoSrc + string(filepath.Separator) + c.Args()[0]

				init.Setup(cfg, destDir, c.String("template"))
			},
		},
		{
			// go run nildev.go build bitbucket.org/nildev/account

			// 1. build container and push it so it would have to be only deployed, also used for local testing
			// 2. push request to be deployed on nildev.io remote infrastructure
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "build project",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "registry",
					Value: "",
					Usage: "registry where to push image, for example quay.io. By default docker hub registry is used, which equals to empty string",
				},
			},
			Action: func(c *cli.Context) {
				pathToGoSrc := os.Getenv("GOPATH")
				if pathToGoSrc == "" {
					log.Fatalf("$GOPATH env is empty ? Do you have golang development environment setup ?")
				}

				pathToGoSrc += "/src"
				if len(c.Args()) == 0 {
					log.Fatalf("Please provide repository to build: your_username/repo_name")
				}

				regsitry := c.String("registry")
				if regsitry != "" {
					regsitry += "/"
				}

				tag := "latest"
				path := c.Args()[0]
				rez := strings.Split(c.Args()[0], ":")
				if len(rez) == 2 {
					tag = rez[1]
					path = rez[0]
				}

				providerParts := strings.SplitN(path, "/", 2)
				if len(providerParts) != 2 {
					log.Fatalf("Given path [%s] is invalid, should be : provider.com/your_username/repo_name", c.Args()[0])
				}
				provider := providerParts[0]
				justRepoName := providerParts[1]
				repoName := regsitry + providerParts[1]
				fullRepoName := provider + string(filepath.Separator) + justRepoName

				log.Infof("src path [%s]", pathToGoSrc)
				log.Infof("Tag [%s]", tag)
				log.Infof("Provider [%s]", provider)
				log.Infof("Org/Project [%s]", justRepoName)
				log.Infof("Registry/Org/Project [%s]", repoName)
				log.Infof("Provider/Org/Project [%s]", fullRepoName)

				s := sh.NewSession()
				//s.ShowCMD = true

				params := []string{
					"run",
					"--rm",
					"--net=host",
					"-v", "/var/run/docker.sock:/var/run/docker.sock",
					"-v", pathToGoSrc + ":/src",
					"nildev/api-builder:latest",
					fullRepoName,
					"github.com/nildev/api-host",
					repoName + ":" + tag,
				}

				log.Infof("Building ... ")
				log.Infof("First time it can take more time because images needs to be downloaded. \n\n")
				s = s.Command("docker", params)
				out, err := s.Output()

				if err != nil {
					log.Fatalf("----\n %s\n\n", err)
				}

				if len(out) > 0 {
					log.Infof("----\n %s\n\n", out)
				}
			},
		},
		{
			// go run nildev.go run --env dev github.com/nildev/echo
			Name:  "run",
			Usage: "run service localy",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "variables, v",
					Value: "",
					Usage: "file with variables to be passed as env variables to your service container",
				},
				cli.StringFlag{
					Name:  "env",
					Value: "",
					Usage: "env indicates which docker-compose from project should be launched, for example docker-compose-${env}.yml",
				},
				cli.StringFlag{
					Name:  "registry",
					Value: "",
					Usage: "registry used to fetch image from, for example quay.io. By default docker hub registry is used, which equals to empty string",
				},
			},
			Action: func(c *cli.Context) {
				log.Infof("%s", c.String("v"))
				log.Infof("%s", c.String("dcf"))
				log.Infof("%s", c.String("registry"))

				// if docker-compose found, run it
				// docker run -d -p "8080:8080" -e "ND_BITBUCKET_CLIENT_ID=xxxx" -e "ND_BITBUCKET_SECRETE=xxx" -e "ND_DATABASE_NAME=nildev" -e "ND_MONGODB_URL=mongodb://192.168.99.100:27017/nildev" blackhole/account:latest
			},
		},
		{
			// go run nildev.go deploy --v file.json github.com/nildev/echo
			Name:  "deploy",
			Usage: "deploy service on nildev.io",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "variables, v",
					Value: "",
					Usage: "file with variables to be passed as env variables to your service container",
				},
				cli.StringFlag{
					Name:  "registry",
					Value: "",
					Usage: "registry used to fetch image from, for example quay.io. By default docker hub registry is used, which equals to empty string",
				},
			},
			Action: func(c *cli.Context) {
				log.Infof("%s", c.String("v"))
				log.Infof("%s", c.String("registry"))

				// use kubernetes Go client to create Pod
			},
		},
		{
			// go run nildev.go io --sourceDir=$GOPATH/src/bitbucket.org/nildev/ping
			Name:  "inout",
			Aliases: []string{"io"},
			Usage: "generate *inout service* integration code required for nildev service container",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "sourceDir",
					Value: "",
					Usage: "full path to service source directory",
				},
				cli.StringFlag{
					Name:  "pathToTpl",
					Value: "",
					Usage: "full path to integration code template",
				},
				cli.StringFlag{
					Name:  "basePattern",
					Value: "/api/v1",
					Usage: "base pattern for api endpoints",
				},
			},
			Action: func(c *cli.Context) {
				tplPath := c.String("pathToTpl")
				basePattern := c.String("basePattern")
				if c.String("sourceDir") == "" {
					log.Fatalf("Please provide path to service source directory")
				}

				inout.Generate(c.String("sourceDir"), tplPath, basePattern)
			},
		},
		{
			// go run nildev.go r --services=bitbucket.org/nildev/ping --containerDir=$GOPATH/src/bitbucket.org/nildev/blackhole
			Name:  "routes",
			Aliases: []string{"r"},
			Usage: "generate routes for given services",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "services",
					Value: "",
					Usage: "full path to service source directory",
				},
				cli.StringFlag{
					Name:  "containerDir",
					Value: ".",
					Usage: "full path to service container source directory",
				},
				cli.StringFlag{
					Name:  "pathToTpl",
					Value: "",
					Usage: "full path to integration code template",
				},
			},
			Action: func(c *cli.Context) {
				tplPath := c.String("pathToTpl")
				if c.String("containerDir") == "" {
					log.Fatalf("Please provide path to service container source directory")
				}

				if c.String("services") == "" {
					log.Fatalf("Please provide comma seperated pathes to services source directory")
				}

				pathesToServices := strings.Split(c.String("services"), ",")

				routes.Generate(c.String("containerDir"), pathesToServices, tplPath)
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Debugf("Could not run [codegen] cmd, %s", err)
	}
}
