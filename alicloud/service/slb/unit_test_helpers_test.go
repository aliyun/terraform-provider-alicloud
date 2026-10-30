package slb

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// These test helpers mirror the ones in alicloud/service/ecs/disk_unit_test.go.
// The migrated slb unit tests patch the unexported CRUD funcs and the
// connectivity.AliyunClient via gomonkey, so they must live in the internal
// package slb test scope. They cannot import alicloud/acctest (which imports
// the root alicloud package that in turn wires slb.ServicePackage, forming a
// cycle), nor alicloud/sweep for the same reason; the local helpers below
// supply the equivalent provider resource map and shared client instead.

var unitTestEndpoints sync.Map
var unitTestSignVersion sync.Map

// unitTestProviderResources supplies the schema the migrated slb unit tests
// inspect; the root alicloud package cannot be imported here. Only the types
// exercised by the TestUnitAlicloudSlb* cases are registered.
func unitTestProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"alicloud_slb_acl":               ResourceAcl(),
		"alicloud_slb_ca_certificate":    ResourceCaCertificate(),
		"alicloud_slb_tls_cipher_policy": ResourceTlsCipherPolicy(),
	}
}

// unitTestSharedClientForRegion mirrors alicloud/sweep.SharedClientForRegion
// so the unit test can be skipped when no cloud credentials are configured.
func unitTestSharedClientForRegion(region string) (interface{}, error) {
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
		Endpoints:   &unitTestEndpoints,
		SignVersion: &unitTestSignVersion,
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

	client, err := conf.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// newInstanceDiff mirrors the root helper for the migrated resource.
// lintignore: R001
func newInstanceDiff(resourceName string, attributes, attributesDiff map[string]interface{}, state *terraform.InstanceState) (*terraform.InstanceDiff, error) {
	p := unitTestProviderResources()
	dOld, _ := schema.InternalMap(p[resourceName].Schema).Data(state, nil)
	dNew, _ := schema.InternalMap(p[resourceName].Schema).Data(state, nil)
	for key, value := range attributes {
		err := dOld.Set(key, value)
		if err != nil {
			return nil, errs.WrapErrorf(err, "[ERROR] the field %s setting error.", key)
		}
	}
	for key, value := range attributesDiff {
		attributes[key] = value
	}

	for key, value := range attributes {
		err := dNew.Set(key, value)
		if err != nil {
			return nil, errs.WrapErrorf(err, "[ERROR] the field %s setting error.", key)
		}
	}

	diff := terraform.NewInstanceDiff()
	objectKey := ""
	for _, key := range helper.MapSort(dNew.State().Attributes) {
		newValue := dNew.State().Attributes[key]
		if objectKey != "" && !strings.HasPrefix(key, objectKey) {
			objectKey = ""
		}
		if objectKey == "" {
			for _, suffix := range []string{"#", "%"} {
				if strings.HasSuffix(key, suffix) {
					objectKey = strings.TrimSuffix(key, suffix)
					break
				}
			}
		}
		oldValue, ok := dOld.State().Attributes[key]
		if ok && oldValue == newValue {
			continue
		}
		if oldValue == "" {
			for _, suffix := range []string{"#", "%"} {
				if strings.HasSuffix(key, suffix) {
					oldValue = "0"
				}
			}
		}
		diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: oldValue, New: newValue}

		if objectKey != "" {
			for removeKey, removeValue := range dOld.State().Attributes {
				if strings.HasPrefix(removeKey, objectKey) {
					if _, ok := dNew.State().Attributes[removeKey]; !ok {
						if len(strings.Split(removeKey, ".")) > 2 {
							delete(diff.Attributes, removeKey)
						} else {
							diff.Attributes[removeKey] = &terraform.ResourceAttrDiff{Old: removeValue, New: ""}
						}
					}
				}
			}
			objectKey = ""
		}
	}
	return diff, nil
}

// unitTestMockClient builds a mock-only connectivity.AliyunClient that does
// not depend on real cloud credentials or network access, so offline gomonkey
// data-source unit tests no longer SKIP when ALICLOUD_ACCESS_KEY/REGION are
// unset. It injects dummy access keys, skips region validation and account-type
// probing (the latter would otherwise call BssOpenApi/STS to detect the caller
// identity), and pre-populates the SLB endpoint so dataSourceTlsCipherPoliciesRead's
// RpcPost -> rpcRequest -> loadApiEndpoint resolves "slb" locally instead of
// probing the location service. gomonkey then intercepts (*client.Client).DoRequest,
// so no real RPC ever leaves the process; production Read logic is unchanged.
func unitTestMockClient() (*connectivity.AliyunClient, error) {
	const region = "cn-hangzhou"
	accessKey, secretKey := "MockAccessKey", "MockSecretKey"
	credentialConfig := new(credentials.Config).SetType("access_key").
		SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	endpoints := &sync.Map{}
	endpoints.Store("slb", "slb."+region+".aliyuncs.com")
	signVersion := &sync.Map{}
	conf := connectivity.Config{
		Region:               connectivity.Region(region),
		RegionId:             region,
		AccessKey:            accessKey,
		SecretKey:            secretKey,
		Protocol:             "HTTPS",
		Endpoints:            endpoints,
		SignVersion:          signVersion,
		SkipRegionValidation: true,
		AccountId:            "1000000000",
		AccountType:          "Domestic",
		Credential:           credential,
	}
	return conf.Client()
}
