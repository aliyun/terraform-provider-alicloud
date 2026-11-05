package dlfnext

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ServicePackage returns the registration declaration of the DlfNext service.
// Both public types are native Framework implementations.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DlfNext",
		Resources: []conns.Resource{
			{TypeName: "alicloud_dlf_next_catalog", Factory: func() resource.Resource { return NewCatalogResource() }},
		},
		DataSources: []conns.DataSource{
			{TypeName: "alicloud_dlf_next_catalogs", Factory: func() datasource.DataSource { return NewCatalogsDataSource() }},
		},
	}
}
