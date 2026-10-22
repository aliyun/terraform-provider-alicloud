package amqp

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Amqp service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Amqp",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_amqp_instance",
				Factory:  ResourceAmqpInstance,
			},
		},
	}
}
