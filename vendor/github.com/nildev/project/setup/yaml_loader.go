package setup

type (
	yamlLoader struct{}
)

// NewYamlLoader returns instance of YAML config Loader
func NewYamlLoader() ConfigReader {

	return &yamlLoader{}
}

// Read yaml config
func (yl *yamlLoader) Read(filepath string) Config {

	return Config{}
}
