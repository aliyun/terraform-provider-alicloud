package alicloud

import (
	"fmt"
	"github.com/aliyun/credentials-go/credentials"
	"log"
	"os"
	"strings"
	"sync"
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

var endpoints sync.Map

// sharedClientForRegion returns a common AlicloudClient setup needed for the sweeper
// functions for a give n region
func sharedClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey, securityToken string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}

	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}

	securityToken = os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:    connectivity.Region(region),
		RegionId:  region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Protocol:  "HTTPS",
		Endpoints: &endpoints,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	if accountId := os.Getenv("ALICLOUD_ACCOUNT_ID"); accountId != "" {
		conf.AccountId = accountId
	}
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential

	// configures a default client for the region, using the above env vars
	client, err := conf.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func sweepAll() bool {
	return os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") == "true"
}
