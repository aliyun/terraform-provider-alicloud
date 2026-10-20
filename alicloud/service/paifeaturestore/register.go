package paifeaturestore

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PaiFeatureStore service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiFeatureStore",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_feature_store_feature_entity",
				Factory:  NewFeatureEntityResource,
			},
		},
	}
}
