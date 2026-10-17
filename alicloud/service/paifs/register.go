package paifs

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PaiFeatureStore
// service. FeatureStoreTask is read-only, so it is registered as a data source.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiFeatureStore",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_pai_feature_store_tasks",
				Factory:  NewPaiFeatureStoreTasksDataSource,
			},
		},
	}
}
