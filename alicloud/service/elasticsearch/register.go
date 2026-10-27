package elasticsearch

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Elasticsearch service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Elasticsearch",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_elasticsearch_instance",
				Factory:  ResourceInstance,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_elasticsearch_instances",
				Factory:  DataSourceInstances,
			},
			{
				TypeName: "alicloud_elasticsearch_zones",
				Factory:  DataSourceZones,
			},
		},
	}
}
