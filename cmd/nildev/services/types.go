package services

type (
	// Environment type
	Environment string

	// Platform is component for creating required artifacts and interacting with specific platform
	Platform interface {
		// Setup services
		Setup(...Name) error
		Run(...Name) error
	}

	// Name of service
	Name string

	// EnvVar type
	EnvVar struct {
		Name   Name
		Values map[Environment]string
	}

	// EnvVars type
	EnvVars []*EnvVar

	// Service represents service configuration
	Service struct {
		Name      Name     `json:"name"`
		Project   string   `json:"project"`
		Image     string   `json:"image"`
		IsPublic  bool     `json:"isPublic"`
		Subscribe []Name   `json:"subscribe"`
		Require   Services `json:"require"`
		Secrets   Secrets  `json:"secrets"`
		EnvVars   EnvVars  `json:"env"`
	}

	// Services type
	Services []*Service

	// Secret type
	Secret struct {
		Name Name
		Data map[string]string
	}

	// Secrets type
	Secrets []*Secret
)

func (s *Services) LoadSecrets(secrets map[string]interface{}) {
	for serviceName, secret := range secrets {
		for secretName, data := range secret.(map[string]interface{}) {
			scr := &Secret{
				Data: map[string]string{},
			}
			scr.Name = Name(secretName)
			for key, val := range data.(map[string]interface{}) {
				scr.Data[key] = val.(string)
			}
			srvc := s.find(serviceName)
			if srvc != nil {
				srvc.Secrets = append(srvc.Secrets, scr)
			}
		}
	}
}

func (s *Services) LoadEnvVars(envVars map[string]interface{}) {
	for serviceName, envs := range envVars {
		for varName, dataPerEnv := range envs.(map[string]interface{}) {
			ev := &EnvVar{
				Values: map[Environment]string{},
			}
			ev.Name = Name(varName)
			for env, val := range dataPerEnv.(map[string]interface{}) {
				ev.Values[Environment(env)] = val.(string)
			}
			srvc := s.find(serviceName)
			if srvc != nil {
				srvc.EnvVars = append(srvc.EnvVars, ev)
			}

		}
	}
}

func (s *Services) find(name string) *Service {
	for _, service := range *s {
		if service.Name == Name(name) {
			return service
		}
	}

	return nil
}

func (scrts *Secrets) find(name string) *Secret {
	for _, secret := range *scrts {
		if secret.Name == Name(name) {
			return secret
		}
	}

	return nil
}

func (ev *EnvVars) find(name string) *EnvVar {
	for _, env := range *ev {
		if env.Name == Name(name) {
			return env
		}
	}

	return nil
}
