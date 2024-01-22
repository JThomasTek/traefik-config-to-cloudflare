package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
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
	// var cloudflareAPI *cloudflare.API
	var err error

	configFile := "./test.yaml"

	traefikConfigWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("creating a new watcher: %s\n", err)
		os.Exit(1)
	}

	defer traefikConfigWatcher.Close()

	go configWatcher(traefikConfigWatcher, configFile)

	st, err := os.Lstat(configFile)
	if err != nil {
		fmt.Printf("getting file info: %s\n", err)
		os.Exit(2)
	}

	if st.IsDir() {
		fmt.Printf("%s is a directory\n", configFile)
		os.Exit(3)
	}

	err = traefikConfigWatcher.Add(filepath.Dir(configFile))
	if err != nil {
		fmt.Printf("adding a new watcher: %s\n", err)
		os.Exit(4)
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

	fmt.Printf("%v running; press ^C to exit\n", time.Now())
	<-make(chan struct{})
}

func getWANIP() string {
	res, err := http.Get("https://ipv4.icanhazip.com")

	if err != nil {
		log.Fatal(err)
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
		fmt.Println(fmt.Errorf("error: %s", err))
	}

	fmt.Println(traefikConfig)
}

func configWatcher(w *fsnotify.Watcher, filename string) {
	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}

			if event.Name == filename && event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println(getWANIP())
				handleConfigChange(filename)
			}
		case err, ok := <-w.Errors:
			if !ok {
				return
			}

			log.Println("error:", err)
		}
	}
}
