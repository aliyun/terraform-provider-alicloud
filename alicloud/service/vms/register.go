package vms

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Voice Messaging
// Service (Dyvmsapi) service. Name is the API product code ("Dyvmsapi") used by
// RpcPost and the conns ServicePackage identifier; the package directory follows
// the canonical namespace (vms), mirroring the sibling sms (Dysmsapi) package.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Dyvmsapi",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_vms_skill",
				Factory:  NewSkillResource,
			},
		},
	}
}
