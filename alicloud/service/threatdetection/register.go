package threatdetection

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ThreatDetection
// (Sas/2018-12-03) service. The custom check item resource and data source are
// implemented natively on terraform-plugin-framework and ride the mux v5 server
// alongside the SDK v2 half.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sas",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_threat_detection_custom_check_items",
				Factory:  NewCustomCheckItemsDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_threat_detection_custom_check_item",
				Factory:  NewCustomCheckItemResource,
			},
		},
	}
}
