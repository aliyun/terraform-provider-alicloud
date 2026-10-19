package cms

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CloudMonitor
// (Cms) alert-webhook resource and data source. Both are native
// terraform-plugin-framework declarations, registered through the framework
// provider assembly loop.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cms",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_cms_alert_webhook",
				Factory:  NewAlertWebhookResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_cms_alert_webhooks",
				Factory:  NewAlertWebhooksDataSource,
			},
		},
	}
}
