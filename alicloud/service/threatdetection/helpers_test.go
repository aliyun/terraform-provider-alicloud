// Package threatdetection_test.
package threatdetection_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

var threatdetectionTestEndpoints sync.Map
var threatdetectionTestSignVersion sync.Map

// threatdetectionTestClientForRegion rebuilds the package-level sweeper client
// on the exported connectivity API: the root-package sharedClientForRegion is
// not reachable from this external test package. It returns interface{} to
// mirror the root sharedClientForRegion signature so existing sweeper bodies
// (rawClient.(*connectivity.AliyunClient)) compile unchanged.
func threatdetectionTestClientForRegion(region string) (interface{}, error) {
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
		Endpoints:   &threatdetectionTestEndpoints,
		SignVersion: &threatdetectionTestSignVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
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

	return conf.Client()
}

// threatdetectionSweepAll mirrors the root-package sweepAll gate; it is not
// reachable from a service test package.
func threatdetectionSweepAll() bool {
	return os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") == "true"
}

// threatdetectionTestSweepPreCheckWithRegions mirrors the root-package
// testSweepPreCheckWithRegions gate; it is not reachable from a service test
// package.
func threatdetectionTestSweepPreCheckWithRegions(region string, supported bool, regions []connectivity.Region) bool {
	find := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
	}
	return (find && !supported) || (!find && supported)
}

// queryAvailableInstanceList mirrors BssOpenApiService.QueryAvailableInstanceList,
// which is not constructible from this external test package (its client field
// is unexported). It reproduces the BssOpenApi QueryAvailableInstances RPC.
func queryAvailableInstanceList(client *connectivity.AliyunClient, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl string) ([]interface{}, error) {
	var response map[string]interface{}
	var endpoint string
	action := "QueryAvailableInstances"
	request := map[string]interface{}{
		"ProductCode":      productCode,
		"ProductType":      productType,
		"SubscriptionType": "Subscription",
	}
	if client.IsInternationalAccount() {
		request = map[string]interface{}{
			"ProductCode": productCodeIntl,
			"ProductType": productTypeIntl,
		}
	}
	if instanceRegion != "" {
		request["Region"] = instanceRegion
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		response = resp
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = productTypeIntl
				request["ProductCode"] = productCodeIntl
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, errs.WrapError(err)
	}
	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	return v.([]interface{}), nil
}

// testAccPreCheckForCleanUpInstances mirrors the root-package helper of the
// same name; the original is not reachable from a service test package. The
// root sendMessage notification is reduced to a log line (disclosed gap for
// this paid, AT001-excluded instance test).
func testAccPreCheckForCleanUpInstances(t *testing.T, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl string) {
	rawClient, err := threatdetectionTestClientForRegion(acctest.DefaultRegionToTest)
	if err != nil {
		t.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	accountId, err := client.AccountId()
	if err != nil {
		t.Errorf("error getting AliCloud client: %s", err)
	}
	deadline := time.Now().Add(30 * time.Minute)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if time.Now().After(deadline) {
				log.Println("Deadline reached, stopping waiting.")
				return
			}
			instances, err := queryAvailableInstanceList(client, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl)
			if err != nil {
				t.Errorf("error querying available instances: %s", err)
				return
			}
			instanceId := ""
			for _, instance := range instances {
				v := instance.(map[string]interface{})
				if v["Status"].(string) != "Normal" {
					continue
				}
				instanceId = v["InstanceID"].(string)
			}
			if instanceId == "" {
				return
			}
			t.Logf("waiting for cleaning up existed instances %s in region %s...", instanceId, instanceRegion)
			log.Printf("[Critical] Please cleaning up instance before running integration test. AccountId: %s, ProductCode: %s, InstanceId: %s, Time: %s",
				accountId, productCode, instanceId, time.Now().String())
			time.Sleep(3 * time.Minute)
		}
	}
}
