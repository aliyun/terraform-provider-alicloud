package alikafka

import (
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// unitTestProviderResources supplies the schema the migrated unit tests inspect.
// The root alicloud package cannot be imported here without an import cycle
// (alikafka -> acctest -> alicloud -> registry -> alikafka), so the resource
// schemas are resolved from the exported factories of this package instead.
// lintignore: R001
func unitTestProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"alicloud_alikafka_consumer_group":                 ResourceConsumerGroup(),
		"alicloud_alikafka_instance":                       ResourceInstance(),
		"alicloud_alikafka_instance_allowed_ip_attachment": ResourceInstanceAllowedIpAttachment(),
		"alicloud_alikafka_sasl_acl":                       ResourceSaslAcl(),
		"alicloud_alikafka_sasl_user":                      ResourceSaslUser(),
		"alicloud_alikafka_scheduled_scaling_rule":         ResourceScheduledScalingRule(),
		"alicloud_alikafka_topic":                          ResourceTopic(),
	}
}

// newInstanceDiff mirrors the root helper alicloud.newInstanceDiff
// (alicloud/common.go:457) for the migrated resources. It uses the local
// unitTestProviderResources so the unit test does not import acctest and so
// avoids the import cycle. Identical control flow to the root helper and to
// alicloud/service/ecs/disk_unit_test.go newInstanceDiff.
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
