package internal

import (
	"math"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Router struct {
	Rule    string `yaml:"rule,omitempty"`
	Service string `yaml:"service,omitempty"`
}

type TraefikConfig struct {
	HTTP struct {
		Routers map[string]Router `yaml:"routers,omitempty"`
	} `yaml:"http,omitempty"`
}

func readTraefikConfig(filename string) (TraefikConfig, error) {
	var config TraefikConfig

	// Read the config file
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	// Unmarshal the config into a simplified TraefikConfig struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func handleConfigChange(filename string) {
	traefikConfig, err := readTraefikConfig(filename)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	log.Info().Str("Rule", traefikConfig.HTTP.Routers["gitlab"].Rule).Str("WAN_IP", GetWANIP()).Msg("Config updated")
}

func TraefikConfigWatcher(w *fsnotify.Watcher, filename string) {
	var (
		// Wait 100ms for new events; each new event resets the timer.
		waitTime = 100 * time.Millisecond

		// Keep track of the timers, as path -> timer.
		mu     sync.Mutex
		timers = make(map[string]*time.Timer)

		// Callback we run.
		eventHandler = func(e fsnotify.Event) {
			handleConfigChange(filename)

			mu.Lock()
			delete(timers, e.Name)
			mu.Unlock()
		}
	)

	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}

			if event.Name == filename && event.Has(fsnotify.Write) {
				mu.Lock()
				t, ok := timers[event.Name]
				mu.Unlock()

				if !ok {
					t = time.AfterFunc(math.MaxInt64, func() { eventHandler(event) })
					t.Stop()

					mu.Lock()
					timers[event.Name] = t
					mu.Unlock()
				}

				t.Reset(waitTime)
			}
		case err, ok := <-w.Errors:
			if !ok {
				return
			}

			log.Error().Err(err)
		}
	}
}
