package paidsw

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAI DSW API
// service.
//
// PAI DSW (Data Science Workshop) exposes a RESTful/ROA API (paths like
// /api/v2/instances/{InstanceId}/snapshots) and has no dedicated Go SDK wired
// into this provider. Every CRUD call on its resources therefore goes through
// the generic ROA entry (RoaPost/RoaGet/RoaDelete) on connectivity.AliyunClient.
// The location service does not resolve "pai-dsw", so the product endpoint is
// pinned in connectivity/endpoint.go (irregularProductEndpoint) to
// pai-dsw.<region>.aliyuncs.com.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiDswApi",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_dsw_api_instance_snapshot",
				Factory:  NewInstanceSnapshotResource,
			},
		},
	}
}
