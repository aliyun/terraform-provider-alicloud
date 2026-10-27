package scdn

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

// unitTestProviderResources supplies the schema the migrated unit test
// inspects; the root alicloud package cannot be imported here.
func unitTestProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"alicloud_scdn_domain":        ResourceScdnDomain(),
		"alicloud_scdn_domain_config": ResourceScdnDomainConfig(),
	}
}

var unitTestEndpoints sync.Map
var unitTestSignVersion sync.Map

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
