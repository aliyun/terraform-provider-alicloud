package kms

import (
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// unitTestProviderResources mirrors the root Provider().ResourcesMap for the
// migrated kms types that participate in newInstanceDiff. The split-package
// internal test cannot reach the root Provider (which would re-introduce the
// registry import cycle), so the unit tests look the factories up through this
// local map instead. Only the resources used by the TestUnit funcs are wired
// here, matching the ecs/nlb split-package convention.
func unitTestProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"alicloud_kms_alias":       resourceAlicloudKmsAlias(),
		"alicloud_kms_key":         resourceAliCloudKmsKey(),
		"alicloud_kms_key_version": resourceAlicloudKmsKeyVersion(),
		"alicloud_kms_secret":      resourceAliCloudKmsSecret(),
	}
}

// newInstanceDiff mirrors the root helper of the same name so the migrated
// unit tests can build an InstanceDiff without depending on the unexported
// root symbol. It follows the ecs/nlb split-package convention verbatim.
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
