package rkvstore

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage is a placeholder for the rkvstore service; declarations move here when the product adopts the split layout.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "R-kvstore",
	}
}
