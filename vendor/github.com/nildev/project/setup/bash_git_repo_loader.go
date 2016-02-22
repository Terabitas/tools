package setup

import (
	"github.com/codeskyblue/go-sh"
	"log"
	"os"
	"path/filepath"
)

type (
	bashGitRepoLoader struct {
		branch string
	}
)

// NewBashGitRepoLoader returns instance of new bash repo loader
func NewBashGitRepoLoader(branch string) RepositoryLoader {
	return &bashGitRepoLoader{
		branch: branch,
	}
}

// Load repo through the bash which uses installed git
func (grl *bashGitRepoLoader) Load(repoPath string, destDir string) {
	shSession := sh.NewSession()
	shSession.Command("git", "clone", repoPath, destDir)
	rez, err := shSession.Output()
	log.Printf("%s", rez)
	if err != nil {
		log.Fatal(err)
	}

	err = os.RemoveAll(destDir + string(filepath.Separator) + ".git")
	if err != nil {
		log.Fatal(err)
	}
}
