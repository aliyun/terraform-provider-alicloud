package realtimecompute

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{Name: "realtimecompute", Resources: []conns.Resource{{TypeName: "alicloud_realtime_compute_vvp_namespace", Factory: NewVvpNamespaceResource}}}
}
