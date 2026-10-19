package cms

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CMS service. The
// Framework-native transformer resource and transformers data source replace
// the SDKv2 prototypes that lived under the provider root; under the
// split-layout (#10858) they now belong to the cms service package.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cms",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_cms_transformer",
				Factory:  NewTransformerResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_cms_transformers",
				Factory:  NewTransformersDataSource,
			},
		},
	}
}
