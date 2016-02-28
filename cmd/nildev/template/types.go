package template

type (
	// Loader is responsible for loading template and return it
	// how templates are found depends on implementation
	Loader interface {
		Load(org, name, version string) ([]byte, error)
	}
)
