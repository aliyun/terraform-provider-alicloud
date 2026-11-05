package schedulerx

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the schedulerx service.
// The package directory is the normalized namespace (schedulerx); the API product
// code stays "schedulerx2" because that is the string used by the RPC calls and the
// root provider endpoint configuration, and it does not change with the package name.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "schedulerx2",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_schedulerx_app_group",
				Factory:  ResourceAppGroup,
			},
			{
				TypeName: "alicloud_schedulerx_job",
				Factory:  ResourceJob,
			},
			{
				TypeName: "alicloud_schedulerx_namespace",
				Factory:  ResourceNamespace,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_schedulerx_namespaces",
				Factory:  DataSourceNamespaces,
			},
		},
	}
}
