package cloudfirewall

import (
	"strings"

	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	sweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	terraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// unitTestProviderResources supplies the schema the migrated cloudfirewall
// unit tests inspect; the root alicloud package cannot be imported here
// because it would re-introduce the registry import cycle.
func unitTestProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"alicloud_cloud_firewall_address_book":                ResourceCloudFirewallAddressBook(),
		"alicloud_cloud_firewall_control_policy":              ResourceCloudFirewallControlPolicy(),
		"alicloud_cloud_firewall_control_policy_order":        ResourceCloudFirewallControlPolicyOrder(),
		"alicloud_cloud_firewall_instance":                    ResourceCloudFirewallInstance(),
		"alicloud_cloud_firewall_vpc_firewall_control_policy": ResourceCloudFirewallVpcFirewallControlPolicy(),
	}
}

// unitTestSharedClientForRegion delegates to the public sweep helper so the
// migrated cloudfirewall unit test shares the same region-isolated shared
// client and credential path as the rest of the provider; no per-package
// credential cache or STS/Config.Client construction is reimplemented here.
func unitTestSharedClientForRegion(region string) (interface{}, error) {
	return sweep.SharedClientForRegion(region)
}

// newInstanceDiff mirrors the root helper for the migrated resources.
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
