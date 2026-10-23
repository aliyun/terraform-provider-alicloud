package cms

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Cms service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cms",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_cms_endpoint_connector",
				Factory:  NewEndpointConnectorResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_cms_endpoint_connectors",
				Factory:  NewEndpointConnectorsDataSource,
			},
		},
	}
}
