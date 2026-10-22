package drds

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DRDS (PolarDB-X)
// service. The legacy root-package registration of alicloud_drds_polardbx_instance
// is removed; the public type name, factory and TestAcc names are preserved here.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Drds",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_drds_polardbx_instance",
				Factory:  ResourcePolardbxInstance,
			},
		},
	}
}
