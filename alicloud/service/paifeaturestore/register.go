package paifeaturestore

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PaiFeatureStore
// service. PaiFeatureStore exposes a ROA-style OpenAPI (version 2023-06-21)
// whose endpoint is registered in connectivity/endpoint.go under the key
// "paifeaturestore".
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiFeatureStore",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_feature_store_llm_config",
				Factory:  NewLLMConfigResource,
			},
		},
	}
}
