package sas

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the SAS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sas",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_threat_detection_property_sca_details",
				Factory:  NewThreatDetectionPropertyScaDetailsDataSource,
			},
		},
	}
}
