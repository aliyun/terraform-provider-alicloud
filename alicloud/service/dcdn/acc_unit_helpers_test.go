package dcdn

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

// Internal-test counterparts of the flat alicloud test helpers for the
// TestUnit* cases that exercise unexported CRUD funcs. These files live in
// package dcdn (internal) so they may call resourceAlicloudDcdn* directly; the
// external package dcdn_test ACC files keep using acctest.Provider. Package
// dcdn must not import alicloud/acctest — that would close the cycle
// acctest -> alicloud -> provider/registry -> service/dcdn — so the helpers
// here mirror alicloud/common.go and alicloud_sweeper_test.go without crossing
// that edge.

var unitTestRegionalEndpoints sync.Map
var unitTestSignVersion sync.Map

// unitTestSharedEndpointCacheForRegion returns the per-region endpoint cache so
// parallel TestUnit cases targeting different regions do not share a single
// global sync.Map and overwrite each other's endpoints. Mirrors the root helper
// and the external dcdnTestSharedEndpointCacheForRegion used by the ACC files.
func unitTestSharedEndpointCacheForRegion(region string) *sync.Map {
	cache, _ := unitTestRegionalEndpoints.LoadOrStore(region, new(sync.Map))
	return cache.(*sync.Map)
}

// unitTestProviderResources supplies the schema the migrated unit tests inspect;
// the root alicloud package cannot be imported here.
func unitTestProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"alicloud_dcdn_ipa_domain":                   resourceAlicloudDcdnIpaDomain(),
		"alicloud_dcdn_waf_domain":                   resourceAlicloudDcdnWafDomain(),
		"alicloud_dcdn_waf_policy":                   resourceAlicloudDcdnWafPolicy(),
		"alicloud_dcdn_waf_policy_domain_attachment": resourceAlicloudDcdnWafPolicyDomainAttachment(),
	}
}

// unitTestSharedClientForRegion mirrors the root test helper so the unit test
// can be skipped when no cloud credentials are configured.
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
		Endpoints:   unitTestSharedEndpointCacheForRegion(region),
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

// newInstanceDiff mirrors alicloud/common.go's helper for the migrated
// resources. Provider() is replaced with unitTestProviderResources, WrapErrorf
// with errs.WrapErrorf and mapSort with helper.MapSort — the root aliases these
// respectively.
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
						// If the attribute has complex elements, there should remove the key, not setting it to empty
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
