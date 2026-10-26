package sas

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage is the registration declaration of the Sas service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sas",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_security_center_sca_details",
				Factory:  NewScaDetailsDataSource,
			},
		},
	}
}
