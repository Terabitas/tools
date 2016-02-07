package services

import "github.com/juju/errors"

const (
	PlatformKubernetes = "kubernetes"
)

// MakePlatform returns specific platform
func MakePlatform(platform, buildDir, env, pathToServicesFile, pathToSecretsFile, pathToEnvFile string) (Platform, error) {

	switch platform {
	case PlatformKubernetes:
		plt, err := MakeKubernetesPlatform(buildDir, env, pathToServicesFile, pathToSecretsFile, pathToEnvFile)
		if err != nil {
			return nil, err
		}
		return plt, nil
	}

	return nil, errors.Trace(errors.Errorf("Platform [%s] is not supported", platform))
}
