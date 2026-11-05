package cloudmonitorservice_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// unitTestProviderResources builds the resource schema map for the migrated
// cloudmonitorservice package's TestUnit tests, mirroring the root
// Provider().ResourcesMap reference (which is unavailable from the external
// cloudmonitorservice_test package). Factory is the exported conns.SDKResource
// field, so the schemas are built without naming the package's unexported
// factory functions.
//
// Only SDKResources are included, faithfully matching Provider().ResourcesMap:
// the SDKv2 provider keeps ResourcesMap and DataSourcesMap as separate maps. The
// TestUnit tests reference resource schemas (resource_*_test.go). Merging
// SDKDataSources into the same map would be wrong — alicloud_cms_monitor_group_instances
// is registered as both a resource (Factory resourceAlicloudCmsMonitorGroupInstances)
// and a data source (Factory dataSourceAlicloudCmsMonitorGroupInstances) under the
// same TypeName, so the data-source factory would overwrite the resource factory
// and dInit.Set("group_id", ...) would fail with "Invalid address to set" because
// the data-source schema has no group_id.
func unitTestProviderResources() map[string]*schema.Resource {
	m := map[string]*schema.Resource{}
	sp := cloudmonitorservice.ServicePackage()
	for _, r := range sp.SDKResources {
		m[r.TypeName] = r.Factory()
	}
	return m
}
