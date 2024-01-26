package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

var stateFile = "/var/traeflare/state.yaml"

type state struct {
	WanIP      string
	Subdomains []string
}

func generateState() error {
	// First check if directory exists and if not then create it
	if _, err := os.Stat("/var/traeflare"); os.IsNotExist(err) {
		os.Mkdir("/var/traeflare", 0755)
	}

	// Create the state file
	file, err := os.Create(stateFile)
	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}

func getState() (state, error) {
	// First check that state file exists
	_, err := os.Stat(stateFile)
	if err != nil {
		// If it doesn't exist, generate it
		err = generateState()

		if err != nil {
			return state{}, err
		}
	}

	var s state

	file, err := os.Open(stateFile)
	if err != nil {
		return s, err
	}

	defer file.Close()

	buffer, err := os.ReadFile(stateFile)
	if err != nil {
		return s, err
	}

	err = yaml.Unmarshal(buffer, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func writeState(newState state) error {
	var oldState state

	file, err := os.Open(stateFile)
	if err != nil {
		return err
	}

	defer file.Close()

	buffer, err := os.ReadFile(stateFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(buffer, &oldState)
	if err != nil {
		return err
	}

	return nil
}

func updateState(config TraefikConfig) error {
	return nil
}

func CompareStateToConfig(config TraefikConfig) error {
	return nil
}

func CompareStateToWanIP(wanIP string) error {
	return nil
}
