package internal

import "github.com/cloudflare/cloudflare-go"

var (
	cloudflareAPI *cloudflare.API
)

func InitializeCloudflareAPIToken(token string) error {
	var err error

	cloudflareAPI, err = cloudflare.NewWithAPIToken(token)

	return err
}
