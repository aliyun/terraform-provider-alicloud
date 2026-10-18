package cms

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CMS (Cloud
// Monitor Service) service. New Resource/Data Source types use the
// terraform-plugin-framework native interfaces and are served by the
// framework half of the muxed provider; the legacy SDKv2 CMS resources
// remain registered in provider.go and coexist unchanged.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cms",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_cms_contact",
				Factory:  NewContactResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_cms_contacts",
				Factory:  NewContactsDataSource,
			},
		},
	}
}
