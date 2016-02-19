package setup

type (
	// Config holds arbitrary data that is loaded
	// from config file
	Config map[string]interface{}

	// Initializer is a component which iterates over repo template
	// and replaces variables with provided config values
	Initializer interface {
		Setup(cfg Config, destDir, templateRepo string)
	}

	// ConfigReader reads config file and returns Config
	ConfigReader interface {
		Read(filepath string) Config
	}

	// RepositoryLoader creates template repo in dest dir
	RepositoryLoader interface {
		Load(repoPath string, destDir string)
	}

	// Renderer iterates over destDir and replaces variables with cfg values
	Renderer interface {
		Render(cfg Config, destDir string)
	}
)
