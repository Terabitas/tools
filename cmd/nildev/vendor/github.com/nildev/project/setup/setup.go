package setup

type (
	initializer struct {
		repoLoader RepositoryLoader
		renderer   Renderer
	}
)

// NewInitializer returns instance of initializer
func NewInitializer() Initializer {
	repoLoader := NewBashGitRepoLoader("master")
	renderer := NewRenderer()
	return &initializer{
		repoLoader: repoLoader,
		renderer:   renderer,
	}
}

// Setup implementation
func (i *initializer) Setup(cfg Config, destDir, templateRepo string) {
	// setup project dir by getting template
	i.repoLoader.Load(templateRepo, destDir)

	// iterate over dir and replace variables with config values
	i.renderer.Render(cfg, destDir)
}
