package template

import (
	"bufio"
	"github.com/juju/errors"
	"github.com/nildev/lib/log"
	"github.com/nildev/lib/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type (
	goPathLoader struct {
		tpl  []byte
		org  string
		name string
		ver  string
	}
)

// NewGoPathLoader constructor
func NewGoPathLoader() Loader {
	return &goPathLoader{}
}

// Load goes through the all $GOPATH and searches for template
// checks for all *.tpl files
func (gpl *goPathLoader) Load(org, name, version string) ([]byte, error) {
	gpl.org = org
	gpl.name = name
	gpl.ver = version

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return []byte{}, errors.New("$GOPATH is not set!")
	}

	rootDir := gopath + string(filepath.Separator) + "src"

	filepath.Walk(rootDir, gpl.visit)
	return gpl.tpl, nil
}

func (gpl *goPathLoader) visit(path string, f os.FileInfo, err error) error {
	//fmt.Printf(" -- checking [%s/%s]", path, f.Name())
	if !f.IsDir() {
		if strings.Contains(f.Name(), ".tpl") {
			rez := gpl.analyseFile(path)
			if rez {
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				remaining, err := utils.PopLine(data)
				if err != nil {
					return err
				}

				gpl.tpl = remaining
			}
		}
	}

	return nil
}

func (gpl *goPathLoader) analyseFile(pathToFile string) bool {
	f, err := os.Open(pathToFile)
	if err != nil {
		log.Errorf("error opening file: %v\n", err)
		return false
	}
	r := bufio.NewReader(f)
	ln, prefix, err := r.ReadLine()
	if err != nil {
		log.Errorf("error reading first line at file: %v\n", err)
		return false
	}

	if prefix {
		log.Errorf("Line too long: %v\n", err)
		return false
	}

	//fmt.Printf("%s ---> %s\n\n", pathToFile, ln)
	rgx, _ := regexp.Compile("^//nildev:template ([a-z]+):([a-z-_]+) (v[0-9]+.[0-9]+.[0-9]+)$")
	matches := rgx.FindAllStringSubmatch(string(ln), -1)

	if len(matches) > 0 {
		if len(matches[0]) != 4 {
			return false
		}
		if matches[0][1] != gpl.org || matches[0][2] != gpl.name {
			return false
		}

		if gpl.ver != "" && matches[0][3] != gpl.ver {
			return false
		}
		return true
	}

	return false
}
