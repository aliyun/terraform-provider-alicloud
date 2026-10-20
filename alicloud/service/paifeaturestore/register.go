package paifeaturestore

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PaiFeatureStore
// service. FeatureView is the first resource onboarded here; it is implemented
// natively on the Terraform Plugin Framework and registered through the
// service-package aggregator.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiFeatureStore",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_paifeaturestore_feature_view",
				Factory:  NewFeatureViewResource,
			},
		},
	}
}
