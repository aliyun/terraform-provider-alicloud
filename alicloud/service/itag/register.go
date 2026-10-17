package itag

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage registers the iTAG (openitag 2022-06-16) service resources and
// data sources with the provider assembly loop. The API product code "Openitag"
// is what the ROA helpers receive; the endpoint lookup lowercases it to "openitag".
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Itag",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_itag_label_user",
				Factory:  NewItagLabelUserResource,
			},
		},
	}
}
