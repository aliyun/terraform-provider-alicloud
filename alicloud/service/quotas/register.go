package quotas

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Quotas service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Quotas",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_quotas_template_service", Factory: ResourceTemplateService},
			{TypeName: "alicloud_quotas_template_applications", Factory: ResourceTemplateApplications},
			{TypeName: "alicloud_quotas_template_quota", Factory: ResourceTemplateQuota},
			{TypeName: "alicloud_quotas_application_info", Factory: ResourceQuotaApplication},
			{TypeName: "alicloud_quotas_quota_alarm", Factory: ResourceQuotaAlarm},
			{TypeName: "alicloud_quotas_quota_application", Factory: ResourceQuotaApplication},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_quotas_quotas", Factory: DataSourceQuotas},
			{TypeName: "alicloud_quotas_application_infos", Factory: DataSourceQuotaApplications},
			{TypeName: "alicloud_quotas_quota_alarms", Factory: DataSourceQuotaAlarms},
			{TypeName: "alicloud_quotas_quota_applications", Factory: DataSourceQuotaApplications},
			{TypeName: "alicloud_quotas_template_applications", Factory: DataSourceTemplateApplications},
		},
	}
}
