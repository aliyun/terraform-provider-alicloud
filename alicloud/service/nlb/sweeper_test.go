package nlb

import (
	"fmt"
	"os"
	"sync"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// nlbTestEndpoints and nlbTestSignVersion stand in for the flat layout's
// shared endpoint cache; the split-package tests cannot reach the root
// sharedClientForRegion, so this product keeps its own (empty) maps, mirroring
// the dts split-package convention.
var (
	nlbTestEndpoints   sync.Map
	nlbTestSignVersion sync.Map
)

// nlbTestClientForRegion rebuilds the package-level sweeper client on the
// split-package test boundary. It mirrors the dts pattern: read AK/SK/token
// from the environment, build a connectivity.Config with local (empty) endpoint
// maps, attach an access_key credential, and return the realised client.
func nlbTestClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}
	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}
	securityToken := os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &nlbTestEndpoints,
		SignVersion: &nlbTestSignVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	credentialConfig := new(credentials.Config).SetType("access_key").
		SetAccessKeyId(accessKey).
		SetAccessKeySecret(secretKey)
	if v := securityToken; v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential
	return conf.Client()
}

// sweepAll mirrors the root helper of the same name so the migrated sweeper
// funcs keep their original gate without referencing the flat layout.
func sweepAll() bool {
	return os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") == "true"
}
