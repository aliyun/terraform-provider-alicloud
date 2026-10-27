package redis

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Redis service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Redis",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_redis_tair_instance",
				Factory:  ResourceTairInstance,
			},
		},
	}
}
