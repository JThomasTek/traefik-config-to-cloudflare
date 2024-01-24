package main

import (
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

/* TODO: 1. Create main infinite loop that checks for WAN IP changes or subdomain changes and updates accordingly
1.a Create a state file that stores the current WAN IP and subdomains
2. Add logging
3. Add command line flags for config file location, cloudflare credentials, etc.
4. Add support for multiple domains
5. Add ability to disable WAN IP updates
*/

type Router struct {
	Rule    string `yaml:"rule,omitempty"`
	Service string `yaml:"service,omitempty"`
}

type TraefikConfig struct {
	HTTP struct {
		Routers map[string]Router `yaml:"routers,omitempty"`
	} `yaml:"http,omitempty"`
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// var cloudflareAPI *cloudflare.API
	var err error

	configFile := "./test.yaml"

	traefikConfigWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		// fmt.Printf("creating a new watcher: %s\n", err)
		log.Fatal().Err(err).Msg("")
	}

	defer traefikConfigWatcher.Close()

	go configWatcher(traefikConfigWatcher, configFile)

	st, err := os.Lstat(configFile)
	if err != nil {
		// fmt.Printf("getting file info: %s\n", err)
		log.Fatal().Err(err).Msg("")
	}

	if st.IsDir() {
		// fmt.Printf("%s is a directory\n", configFile)
		log.Fatal().Msgf("%s is a directory\n", configFile)
	}

	err = traefikConfigWatcher.Add(filepath.Dir(configFile))
	if err != nil {
		// fmt.Printf("adding a new watcher: %s\n", err)
		log.Fatal().Err(err).Msg("")
	}

	// cloudflareAPI, err = cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))

	// if err != nil {
	// 	cloudflareAPI, err = cloudflare.New(os.Getenv("CLOUDFLARE_API_KEY"), os.Getenv("CLOUDFLARE_EMAIL"))

	// 	if err != nil {
	// 		log.Fatal("No valid API credentials provided")
	// 	}
	// }

	// ctx := context.Background()

	// var proxied *bool = new(bool)
	// *proxied = false

	// record, err := cloudflareAPI.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(os.Getenv("CLOUDFLARE_ZONE_ID")), cloudflare.CreateDNSRecordParams{
	// 	Type:    "A",
	// 	Content: string(resBody),
	// 	Proxied: proxied,
	// 	Name:    "ex.teknand.io",
	// 	TTL:     1,
	// })

	// record, err := cloudflareAPI.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(os.Getenv("CLOUDFLARE_ZONE_ID")), cloudflare.UpdateDNSRecordParams{
	// 	TTL:     1,
	// 	Name:    "ex.teknand.io",
	// 	Content: string(resBody),
	// 	ID:      "7f28579be91ec960b2cee0e9463db31d",
	// })

	// err = cloudflareAPI.DeleteDNSRecord(ctx, cloudflare.ZoneIdentifier(os.Getenv("CLOUDFLARE_ZONE_ID")), "7f28579be91ec960b2cee0e9463db31d")

	// records, resultInfo, err := cloudflareAPI.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(os.Getenv("CLOUDFLARE_ZONE_ID")), cloudflare.ListDNSRecordsParams{})

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// tconfig, err := readTraefikConfig("test.yaml")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(getWANIP())
	// fmt.Println(tconfig)
	// fmt.Println(resultInfo.Count, records)

	log.Info().Msg("Watching for config changes")
	<-make(chan struct{})
}

func getWANIP() string {
	res, err := http.Get("https://ipv4.icanhazip.com")

	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	resBody, _ := io.ReadAll(res.Body)

	return string(resBody)
}

func readTraefikConfig(filename string) (TraefikConfig, error) {
	var config TraefikConfig

	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func handleConfigChange(filename string) {
	traefikConfig, err := readTraefikConfig(filename)
	if err != nil {
		log.Error().Err(err)
		// fmt.Println(fmt.Errorf("error: %s", err))
	}

	log.Info().Str("Rule", traefikConfig.HTTP.Routers["gitlab"].Rule).Msg("Config updated")
	// fmt.Println(traefikConfig)
}

func configWatcher(w *fsnotify.Watcher, filename string) {
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
