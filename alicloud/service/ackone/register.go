package ackone

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the AckOne service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "AckOne",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_ack_one_cluster",
				Factory:  ResourceCluster,
			},
			{
				TypeName: "alicloud_ack_one_membership_attachment",
				Factory:  ResourceMembershipAttachment,
			},
		},
	}
}
