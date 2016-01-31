package setup

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"text/template"
)

type (
	renderer struct {
		cfg Config
	}
)

// NewRenderer returns new renderer
func NewRenderer() Renderer {
	return &renderer{}
}

// Render impl
func (r *renderer) Render(cfg Config, destDir string) {
	r.cfg = cfg
	walkDirs(destDir, cfg)

	var err error

	// render files
	err = filepath.Walk(destDir, r.visit)
	if err != nil {
		log.Fatalf("Error while iterating over directory: %s", err)
	}
}

func renderFile(path string, cfg Config) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("%s", err)
	}

	output, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Could not open output file: %s", err)
	}
	defer func() {
		err := output.Close()
		if err != nil {
			log.Fatal("Could not close file!", err)
		}
	}()

	tpl, err := template.New("render").Parse(string(b))
	if err != nil {
		log.Fatalf("Could not parse template: %s", err)
	}
	tpl.Execute(output, cfg)
}

func render(name string, tpl string, data interface{}) []byte {
	var cnt []byte
	var rendered *bytes.Buffer
	rendered = bytes.NewBuffer(cnt)
	tmpl, err := template.New(name).Parse(tpl)
	if err != nil {
		log.Fatalf("Could not parse [%s], [%v]", name, err)
	}
	err = tmpl.Execute(rendered, data)
	if err != nil {
		log.Fatalf("Could not render, [%v]", err)
	}

	return rendered.Bytes()
}

func (r *renderer) visit(path string, f os.FileInfo, err error) error {
	// rename folder
	newFileName := string(render(f.Name(), path, r.cfg))

	if path != newFileName {
		os.Rename(path, newFileName)
	}

	if !f.IsDir() {
		renderFile(newFileName, r.cfg)
	}

	return nil
}

func readDirs(root string) []string {
	f, _ := os.Open(root)
	names, _ := f.Readdirnames(-1)
	f.Close()

	sort.Strings(names)

	return names
}

func hasRenderable(names []string, pattern string) bool {
	for _, n := range names {
		x, _ := filepath.Match(pattern, n)
		if x {
			return true
		}
	}

	return false
}

func renderDirs(path string, dirs []string, cfg Config) {
	for _, name := range dirs {
		oldPath := path + string(filepath.Separator) + name
		newFileName := string(render(name, path+string(filepath.Separator)+name, cfg))
		os.Rename(oldPath, newFileName)
	}
}

func walkDirs(root string, cfg Config) {
	dirs := readDirs(root)

	if hasRenderable(dirs, "{{.*}}") {
		renderDirs(root, dirs, cfg)
		dirs = readDirs(root)
	}

	for _, name := range dirs {
		walkDirs(root+string(filepath.Separator)+name, cfg)
	}
}
