package main

import (
	"os"
	"path/filepath"

	"github.com/JThomasTek/traefik-config-to-cloudflare/internal"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/* TODO: 1. Create main infinite loop that checks for WAN IP changes or subdomain changes and updates accordingly
1.a Create a state file that stores the current WAN IP and subdomains
2. Add logging
3. Add command line flags for config file location, cloudflare credentials, etc.
4. Add support for multiple domains
5. Add ability to disable WAN IP updates
*/

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var err error

	configFile := "./test.yaml"

	if os.Getenv("CLOUDFLARE_API_TOKEN") != "" {
		err = internal.InitializeCloudflareAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
	}

	traefikConfigWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		// fmt.Printf("creating a new watcher: %s\n", err)
		log.Fatal().Err(err).Msg("")
	}

	defer traefikConfigWatcher.Close()

	go internal.TraefikConfigWatcher(traefikConfigWatcher, configFile)
	go internal.WanIPCheck(60)

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
