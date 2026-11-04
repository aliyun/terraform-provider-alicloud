package ga

import (
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Shared white-box test helpers for the internal package ga TestUnit cases.
//
// The black-box acceptance tests live in the external ga_test package, which is
// free to import alicloud/acctest. The mock-based TestUnit cases, however, must
// compile in the INTERNAL package ga so they can call the unexported resource
// CRUD funcs (resourceAlicloudGa...{Create,Read,Update,Delete}) directly. An
// internal package ga test file CANNOT import alicloud/acctest: acctest
// transitively imports the root alicloud package, which imports
// alicloud/provider/registry, which imports alicloud/service/ga, forming an
// import cycle in the test binary ("import cycle not allowed in test").
//
// Following the databasegateway gateway_unit_test.go precedent, the provider
// resource map and instance-diff helper are reproduced locally here instead of
// importing acctest. The shared test client itself delegates to
// alicloud/sweep.SharedClientForRegion (region-isolated endpoint cache) rather
// than being reproduced, so this package does not carry a second copy of the
// sweeper client. This single file holds the shared definitions so each
// resource's TestUnit file can reference them without redeclaring them.

// unitTestProviderResources supplies the schema every TestUnit inspects. It
// maps the full Terraform resource type to its exported accessor; the root
// alicloud package cannot be imported here.
func unitTestProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"alicloud_ga_accelerator":                               ResourceAliCloudGaAccelerator(),
		"alicloud_ga_listener":                                  ResourceAliCloudGaListener(),
		"alicloud_ga_bandwidth_package":                         ResourceAliCloudGaBandwidthPackage(),
		"alicloud_ga_bandwidth_package_attachment":              ResourceAliCloudGaBandwidthPackageAttachment(),
		"alicloud_ga_endpoint_group":                            ResourceAliCloudGaEndpointGroup(),
		"alicloud_ga_ip_set":                                    ResourceAliCloudGaIpSet(),
		"alicloud_ga_forwarding_rule":                           ResourceAliCloudGaForwardingRule(),
		"alicloud_ga_acl":                                       ResourceAliCloudGaAcl(),
		"alicloud_ga_acl_attachment":                            ResourceAliCloudGaAclAttachment(),
		"alicloud_ga_acl_entry_attachment":                      ResourceAliCloudGaAclEntryAttachment(),
		"alicloud_ga_additional_certificate":                    ResourceAliCloudGaAdditionalCertificate(),
		"alicloud_ga_accelerator_spare_ip_attachment":           ResourceAlicloudGaAcceleratorSpareIpAttachment(),
		"alicloud_ga_access_log":                                ResourceAlicloudGaAccessLog(),
		"alicloud_ga_basic_accelerator":                         ResourceAliCloudGaBasicAccelerator(),
		"alicloud_ga_basic_endpoint_group":                      ResourceAliCloudGaBasicEndpointGroup(),
		"alicloud_ga_basic_ip_set":                              ResourceAliCloudGaBasicIpSet(),
		"alicloud_ga_basic_accelerate_ip":                       ResourceAliCloudGaBasicAccelerateIp(),
		"alicloud_ga_basic_endpoint":                            ResourceAliCloudGaBasicEndpoint(),
		"alicloud_ga_basic_accelerate_ip_endpoint_relation":     ResourceAlicloudGaBasicAccelerateIpEndpointRelation(),
		"alicloud_ga_domain":                                    ResourceAlicloudGaDomain(),
		"alicloud_ga_custom_routing_endpoint_group":             ResourceAliCloudGaCustomRoutingEndpointGroup(),
		"alicloud_ga_custom_routing_endpoint_group_destination": ResourceAliCloudGaCustomRoutingEndpointGroupDestination(),
		"alicloud_ga_custom_routing_endpoint":                   ResourceAliCloudGaCustomRoutingEndpoint(),
		"alicloud_ga_custom_routing_endpoint_traffic_policy":    ResourceAliCloudGaCustomRoutingEndpointTrafficPolicy(),
	}
}

// unitTestSharedClientForRegion delegates to sweep.SharedClientForRegion so the
// white-box TestUnit shares the canonical, region-isolated sweeper client
// (per-region endpoint cache via sweep.SharedEndpointCacheForRegion) instead of a
// package-global cache that could route one region's request to another region's
// cached endpoint. Returns the concrete *connectivity.AliyunClient boxed as
// interface{} so call sites that assert rawClient.(*connectivity.AliyunClient)
// keep working unchanged; the signature is preserved.
func unitTestSharedClientForRegion(region string) (interface{}, error) {
	return sweep.SharedClientForRegion(region)
}

// newInstanceDiff mirrors the root acctest.NewInstanceDiff helper for the
// migrated resource. Reproduced here verbatim from the databasegateway
// precedent because importing acctest would form a test import cycle.
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
