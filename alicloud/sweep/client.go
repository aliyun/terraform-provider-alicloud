package sweep

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// SharedClientForRegionWithBackendRegions selects a backend region before
// creating the shared sweeper client.
func SharedClientForRegionWithBackendRegions(region string, supported bool, regions []connectivity.Region) (interface{}, error) {
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
			// Try to use backup region first
			if backupRegionFind {
				log.Printf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, backupRegion)
				region = backupRegion
			} else if len(regions) > 0 {
				// If backup region is not available, use the first supported region
				targetRegion := string(regions[0])
				log.Printf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, targetRegion)
				region = targetRegion
			}
		} else {
			if !backupRegionFind {
				log.Printf("Skipping unsupported region %s. Unsupported regions: %s. Using %s as this test region", region, regions, backupRegion)
				region = backupRegion
			}
		}
	}
	return SharedClientForRegion(region)
}

var signVersion sync.Map
var regionalEndpoints sync.Map

// SharedEndpointCacheForRegion returns the product endpoint cache for a region.
// SDK endpoint caches are keyed by product. Sharing one across regions can route
// a Hangzhou request to a previously cached Beijing endpoint.
func SharedEndpointCacheForRegion(region string) *sync.Map {
	cache, _ := regionalEndpoints.LoadOrStore(region, new(sync.Map))
	return cache.(*sync.Map)
}

// SharedClientForRegion returns a common client setup needed for sweeper
// functions for a given region.
func SharedClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey, securityToken string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}

	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}

	securityToken = os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   SharedEndpointCacheForRegion(region),
		SignVersion: &signVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	accountId := os.Getenv("ALICLOUD_ACCOUNT_ID")
	if accountId == "" {
		accountId = os.Getenv("ALIBABA_CLOUD_ACCOUNT_ID")
	}
	if accountId != "" {
		conf.AccountId = accountId
	}
	accountType := os.Getenv("ALICLOUD_ACCOUNT_TYPE")
	if accountType == "" {
		accountType = os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE")
	}
	if accountType != "" {
		conf.AccountType = accountType
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
