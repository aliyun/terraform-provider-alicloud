package sas

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage is a placeholder for the sas service; declarations move here when the product adopts the split layout.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sas",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_threat_detection_sca_details",
				Factory:  NewScaDetailsDataSource,
			},
		},
	}
}
