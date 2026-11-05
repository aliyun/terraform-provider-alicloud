package eip

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Eip",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_eip_address",
				Factory:  resourceAliCloudEipAddress,
			},
			{
				TypeName: "alicloud_eip",
				Factory:  resourceAliCloudEipAddress,
			},
			{
				TypeName: "alicloud_eip_association",
				Factory:  resourceAliCloudEipAssociation,
			},
			{
				TypeName: "alicloud_eip_segment_address",
				Factory:  resourceAliCloudEipSegmentAddress,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_eip_addresses",
				Factory:  dataSourceAlicloudEipAddresses,
			},
			{
				TypeName: "alicloud_eips",
				Factory:  dataSourceAlicloudEipAddresses,
			},
		},
	}
}
