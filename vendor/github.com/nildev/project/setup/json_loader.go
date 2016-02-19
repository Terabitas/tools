package setup

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	jsonLoader struct{}
)

// NewJSONLoader returns instance of JSON config Loader
func NewJSONLoader() ConfigReader {
	return &jsonLoader{}
}

// Read json config
func (yl *jsonLoader) Read(filepath string) Config {

	file, _ := os.Open(filepath)
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	return configuration
}
