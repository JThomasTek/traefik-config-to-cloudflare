package main

import (
	"os"
	"path/filepath"

	"github.com/JThomasTek/traefik-config-to-cloudflare/internal"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/* TODO: 1. Create main infinite loop that checks for WAN IP changes or subdomain changes and updates accordingly -DONE
1.a Create a state file that stores the current WAN IP and subdomains -DONE
2. Add logging
3. Add command line flags for config file location, cloudflare credentials, etc.
4. Add support for multiple domains
5. Add ability to disable WAN IP updates
*/

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var err error

	traefikConfigFile := "/etc/traefik/config.yaml"

	if os.Getenv("CLOUDFLARE_API_TOKEN") != "" {
		err = internal.InitializeCloudflareAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"), os.Getenv("CLOUDFLARE_ZONE_ID"))
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
	}

	err = internal.InitialWanIPCheck()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	internal.InitialConfigCheck(traefikConfigFile)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	traefikConfigWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	defer traefikConfigWatcher.Close()

	go internal.TraefikConfigWatcher(traefikConfigWatcher, traefikConfigFile)
	go internal.WanIPCheck(60)

	st, err := os.Lstat(traefikConfigFile)
	if err != nil {
		// fmt.Printf("getting file info: %s\n", err)
		log.Fatal().Err(err).Msg("")
	}

	if st.IsDir() {
		// fmt.Printf("%s is a directory\n", configFile)
		log.Fatal().Msgf("%s is a directory\n", traefikConfigFile)
	}

	err = traefikConfigWatcher.Add(filepath.Dir(traefikConfigFile))
	if err != nil {
		// fmt.Printf("adding a new watcher: %s\n", err)
		log.Fatal().Err(err).Msg("")
	}

	log.Info().Msg("Watching for config changes")
	<-make(chan struct{})
}
