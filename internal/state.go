package internal

import (
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var (
	stateFolder = "/etc/traeflare/"
	stateFile   = stateFolder + "state.yaml"
	mu          sync.Mutex
)

type state struct {
	WanIP   string
	Routers map[string]Router
}

func generateState() error {
	log.Info().Msg("Generate the state file")

	// First check if directory exists and if not then create it
	if _, err := os.Stat(stateFolder); os.IsNotExist(err) {
		err = os.Mkdir(stateFolder, 0744)
		if err != nil {
			return err
		}
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
	log.Info().Msg("Reading state file data")

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

	mu.Lock()
	buffer, err := os.ReadFile(stateFile)
	mu.Unlock()
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
	log.Info().Msg("Writing to the state file")

	data, err := yaml.Marshal(newState)
	if err != nil {
		return err
	}

	mu.Lock()
	err = os.WriteFile(stateFile, data, 0644)
	mu.Unlock()
	if err != nil {
		return err
	}

	return nil
}

func CompareStateToConfig(config TraefikConfig) error {
	log.Info().Msg("Comparing state file to config")

	s, err := getState()
	if err != nil {
		return err
	}

	changed := false

	// Check if any new subdomains were added to the config
	for k, v := range config.HTTP.Routers {
		_, ok := s.Routers[k]
		if !ok {
			s.Routers[k] = v
			changed = true
			// TODO: Perform Cloudflare DNS add
			log.Info().Msg("Performing Cloudflare DNS add")
		}
	}

	// Check if any subdomains were removed from the config
	for k := range s.Routers {
		_, ok := config.HTTP.Routers[k]
		if !ok {
			delete(s.Routers, k)
			changed = true
			// TODO: Perform Cloudflare DNS remove
			log.Info().Msg("Performing Cloudflare DNS remove")
		}
	}

	if changed {
		if err = writeState(s); err != nil {
			return err
		}
	}

	return nil
}

func CompareStateToWanIP(wanIP string) error {
	log.Info().Msg("Comparing state file WAN IP to actual WAN IP")

	s, err := getState()
	if err != nil {
		return err
	}

	// Check if the WAN IP has changed, update it if it has
	if s.WanIP != wanIP {
		s.WanIP = wanIP

		// TODO: Update Cloudflare DNS records with new WAN IP
		if err = writeState(s); err != nil {
			return err
		}
	}

	return nil
}
