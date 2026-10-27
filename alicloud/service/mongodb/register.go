package mongodb

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{Name: "mongodb", SDKResources: []conns.SDKResource{{TypeName: "alicloud_mongodb_account", Factory: ResourceAccount}}}
}
