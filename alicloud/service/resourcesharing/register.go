package resourcesharing

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage is a placeholder for the resourcesharing service; declarations move here when the product adopts the split layout.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ResourceSharing",
	}
}
