package redis

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Redis (R-kvstore
// TairInstance) service for the flat layout. Only the types migrated into this
// package are declared here; sibling Redis resources/data sources that remain in
// the root package keep their existing provider.go map entries.
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
