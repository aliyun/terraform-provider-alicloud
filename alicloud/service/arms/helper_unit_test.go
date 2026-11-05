package arms_test

import (
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/arms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// unitTestProviderResources mirrors the former root test helper so the arms
// unit tests can resolve resource schemas without depending on the root
// Provider(). It returns every public SDKv2 resource factory registered by
// the arms ServicePackage, so any p["alicloud_arms_*"] lookup in the test
// files resolves to the migrated schema.
func unitTestProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"alicloud_arms_synthetic_task":        arms.ResourceSyntheticTask(),
		"alicloud_arms_grafana_workspace":     arms.ResourceGrafanaWorkspace(),
		"alicloud_arms_env_custom_job":        arms.ResourceEnvCustomJob(),
		"alicloud_arms_env_service_monitor":   arms.ResourceEnvServiceMonitor(),
		"alicloud_arms_env_pod_monitor":       arms.ResourceEnvPodMonitor(),
		"alicloud_arms_addon_release":         arms.ResourceAddonRelease(),
		"alicloud_arms_env_feature":           arms.ResourceEnvFeature(),
		"alicloud_arms_environment":           arms.ResourceEnvironment(),
		"alicloud_arms_prometheus_monitoring": arms.ResourcePrometheusMonitoring(),
		"alicloud_arms_alert_contact":         arms.ResourceAlertContact(),
		"alicloud_arms_alert_robot":           arms.ResourceAlertRobot(),
		"alicloud_arms_alert_contact_group":   arms.ResourceAlertContactGroup(),
		"alicloud_arms_dispatch_rule":         arms.ResourceDispatchRule(),
		"alicloud_arms_prometheus_alert_rule": arms.ResourcePrometheusAlertRule(),
		"alicloud_arms_prometheus":            arms.ResourcePrometheus(),
		"alicloud_arms_integration_exporter":  arms.ResourceIntegrationExporter(),
		"alicloud_arms_remote_write":          arms.ResourceRemoteWrite(),
	}
}

// newInstanceDiff mirrors the former root helper (alicloud/common.go) and the
// product-subpackage precedent in service/ecs (disk_unit_test.go): it resolves
// schemas through unitTestProviderResources instead of the root Provider(), so
// the migrated arms_test package stays free of root-package imports.
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
