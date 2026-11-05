package kvcachestore

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Kvcachestore service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Kvcachestore",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_kvcachestore_kv_cache_store",
				Factory:  ResourceKvCacheStore,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_kvcachestore_kv_cache_stores",
				Factory:  DataSourceKvCacheStores,
			},
		},
	}
}
