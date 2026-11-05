package aligreen

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Aligreen service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Aligreen",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_aligreen_audit_callback",
				Factory:  resourceAliCloudAligreenAuditCallback,
			},
			{
				TypeName: "alicloud_aligreen_biz_type",
				Factory:  resourceAliCloudAligreenBizType,
			},
			{
				TypeName: "alicloud_aligreen_callback",
				Factory:  resourceAliCloudAligreenCallback,
			},
			{
				TypeName: "alicloud_aligreen_image_lib",
				Factory:  resourceAliCloudAligreenImageLib,
			},
			{
				TypeName: "alicloud_aligreen_keyword_lib",
				Factory:  resourceAliCloudAligreenKeywordLib,
			},
			{
				TypeName: "alicloud_aligreen_oss_stock_task",
				Factory:  resourceAliCloudAligreenOssStockTask,
			},
		},
	}
}
