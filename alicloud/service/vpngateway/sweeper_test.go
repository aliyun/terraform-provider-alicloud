// Package vpngateway_test.
package vpngateway_test

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// vpngatewayTestEndpoints and vpngatewayTestSignVersion stand in for the flat
// layout's shared endpoint cache and sign-version map; the split-package tests
// cannot reach the root sharedClientForRegion, so this product keeps its own
// (empty) maps, mirroring the actiontrail/nlb split-package convention.
var (
	vpngatewayTestEndpoints   sync.Map
	vpngatewayTestSignVersion sync.Map
)

// vpngatewayTestClientForRegion rebuilds the package-level sweeper client on
// the exported connectivity API: the flat-suite sharedClientForRegion is not
// reachable from this package. It mirrors the root helper's (interface{}, error)
// signature so existing sweeper bodies that assert rawClient.(*connectivity.AliyunClient)
// keep working unchanged.
func vpngatewayTestClientForRegion(region string) (interface{}, error) {
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
		Endpoints:   &vpngatewayTestEndpoints,
		SignVersion: &vpngatewayTestSignVersion,
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
	credentialConfig := new(credentials.Config).SetType("access_key").
		SetAccessKeyId(accessKey).
		SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential

	client, err := conf.Client()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// vpngatewaySweepAll mirrors the root-package sweepAll gate; it is not
// reachable from a service test package.
func vpngatewaySweepAll() bool {
	return os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") == "true"
}
