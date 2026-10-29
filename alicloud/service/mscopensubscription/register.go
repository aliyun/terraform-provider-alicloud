package mscopensubscription

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the MscOpenSubscription service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "MscOpenSubscription",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_msc_sub_contact",
				Factory:  ResourceContact,
			},
			{
				TypeName: "alicloud_msc_sub_subscription",
				Factory:  ResourceSubscription,
			},
			{
				TypeName: "alicloud_msc_sub_webhook",
				Factory:  ResourceWebhook,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_msc_sub_contacts",
				Factory:  DataSourceContacts,
			},
			{
				TypeName: "alicloud_msc_sub_subscriptions",
				Factory:  DataSourceSubscriptions,
			},
			{
				TypeName: "alicloud_msc_sub_webhooks",
				Factory:  DataSourceWebhooks,
			},
			{
				TypeName: "alicloud_msc_sub_contact_verification_message",
				Factory:  DataSourceContactVerificationMessage,
			},
		},
	}
}
