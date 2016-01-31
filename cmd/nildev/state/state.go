package state

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	DefaultStateFile = ".nildev"
)

// State type
type State struct {
	Token    string `json:"token"`
	Provider string `json:"provider"`
}

// Persist state
func Persist(s State) {
	d, err := json.Marshal(s)
	if err != nil {
		panic(fmt.Sprintf("Could not save state: [%s]", err))
	}
	ioutil.WriteFile(DefaultStateFile, d, 0777)
}

// Load state
func Load() State {
	var s State
	b, err := ioutil.ReadFile(DefaultStateFile)
	if err != nil {
		return s
	}

	err = json.Unmarshal(b, &s)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal state: [%s]", err))
	}

	return s
}
