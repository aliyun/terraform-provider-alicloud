package alicloud

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sharedClientForRegion returns a common AlicloudClient setup needed for the sweeper
// functions for a given region
func sharedClientForRegionWithBackendRegions(region string, supported bool, regions []connectivity.Region) (interface{}, error) {
	find := false
	backupRegion := string(connectivity.APSouthEast1)
	backupRegionFind := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
		if string(r) == backupRegion {
			backupRegionFind = true
		}
	}

	if (find && !supported) || (!find && supported) {
		if supported {
			if backupRegionFind {
				log.Printf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, backupRegion)
				region = backupRegion
			}
		} else {
			if !backupRegionFind {
				log.Printf("Skipping unsupported region %s. Unsupported regions: %s. Using %s as this test region", region, regions, backupRegion)
				region = backupRegion
			}
		}
	}
	return sharedClientForRegion(region)
}

// sharedClientForRegion returns a common AlicloudClient setup needed for the sweeper
// functions for a give n region
func sharedClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}

	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}

	conf := connectivity.Config{
		Region:    connectivity.Region(region),
		RegionId:  region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Protocol:  "HTTPS",
		Endpoints: make(map[string]interface{}),
	}
	if accountId := os.Getenv("ALICLOUD_ACCOUNT_ID"); accountId != "" {
		conf.AccountId = accountId
	}

	// configures a default client for the region, using the above env vars
	client, err := conf.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}
