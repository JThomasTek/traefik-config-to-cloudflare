package internal

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
	"gopkg.in/yaml.v3"
)

type cf struct {
	cloudflareAPI *cloudflare.API
	zoneID        string
}

var (
	cloudflareData cf
	commentMessage string = "Managed by traeflare: "
	proxied        *bool  = new(bool)
)

func InitializeCloudflareAPIToken(token string, zoneID string) error {
	var err error

	cloudflareData.cloudflareAPI, err = cloudflare.NewWithAPIToken(token)
	cloudflareData.zoneID = zoneID
	*proxied = true

	return err
}

func InitializeCloudflareAPIKey(key string, email string, zoneID string) error {
	var err error

	cloudflareData.cloudflareAPI, err = cloudflare.New(key, email)
	cloudflareData.zoneID = zoneID
	*proxied = true

	return err
}

func AddSubdomain(routerIdentifier string, rule string, wanIP string) error {
	ctx := context.Background()

	_, err := cloudflareData.cloudflareAPI.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(cloudflareData.zoneID), cloudflare.CreateDNSRecordParams{
		Type:    "A",
		Name:    rule,
		Content: wanIP,
		Comment: fmt.Sprintf("%s%s", commentMessage, routerIdentifier),
		TTL:     1,
		Proxied: proxied,
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateWanIP(s state) error {
	ctx := context.Background()

	records, _, err := cloudflareData.cloudflareAPI.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(cloudflareData.zoneID), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record.Comment) >= 22 {
			substr := record.Comment[len(commentMessage):]
			_, ok := s.Routers[substr]
			if ok {
				cloudflareData.cloudflareAPI.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(cloudflareData.zoneID), cloudflare.UpdateDNSRecordParams{
					ID:      record.ID,
					Content: s.WanIP,
				})
			}
		}
	}

	return nil
}

func ListDNSRecords() error {
	ctx := context.Background()

	records, _, err := cloudflareData.cloudflareAPI.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(cloudflareData.zoneID), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		return err
	}

	recordData, err := yaml.Marshal(records)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n\n", string(recordData))

	return nil
}

func DeleteSubdomain(routerIdentifier string) error {
	ctx := context.Background()

	records, _, err := cloudflareData.cloudflareAPI.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(cloudflareData.zoneID), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record.Comment) >= 22 {
			substr := record.Comment[len(commentMessage):]
			if substr == routerIdentifier {
				err = cloudflareData.cloudflareAPI.DeleteDNSRecord(ctx, cloudflare.ZoneIdentifier(cloudflareData.zoneID), record.ID)
				if err != nil {
					return err
				}

				break
			}
		}
	}

	return nil
}
