package ens

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ENS service.
//
// alicloud_ens_bucket_acl and alicloud_ens_bucket_acls are new objects, so per
// the v2 implementation strategy they are written as Plugin Framework native
// (resource.Resource / datasource.DataSource + tfsdk.Schema) and registered
// here, not in the SDK v2 provider.go map literal. The muxed provider serves
// them through the framework half.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ens",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_ens_bucket_acl",
				Factory:  NewBucketAclResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_ens_bucket_acls",
				Factory:  NewBucketAclsDataSource,
			},
		},
	}
}
