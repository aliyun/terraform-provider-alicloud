package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlicloudEssNotifications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEssNotificationsRead,
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notifications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_types": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_encoding": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEssNotificationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeNotificationConfigurations"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	if scalingGroupId, ok := d.GetOk("scaling_group_id"); ok && scalingGroupId.(string) != "" {
		request["ScalingGroupId"] = scalingGroupId.(string)
	}
	var allNotifications []interface{}
	for {
		raw, err := client.RpcPost("Ess", "2014-08-28", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ess_notifications", action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.NotificationConfigurationModels.NotificationConfigurationModel", raw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, "$.NotificationConfigurationModels.NotificationConfigurationModel", raw)
		}
		addDebug(action, raw, request)

		if len(v.([]interface{})) < 1 {
			break
		}
		allNotifications = append(allNotifications, v.([]interface{})...)
		if len(v.([]interface{})) < PageSizeLarge {
			break
		} else {
			continue
		}
	}

	var filteredNotifications []interface{}
	idsMap := make(map[string]string)
	if ids, okIds := d.GetOk("ids"); okIds {
		for _, i := range ids.([]interface{}) {
			if i == nil {
				continue
			}
			idsMap[i.(string)] = i.(string)
		}
		for _, n := range allNotifications {
			var object map[string]interface{}
			object = n.(map[string]interface{})
			if _, ok := idsMap[object["NotificationArn"].(string)]; !ok {
				continue
			}
			filteredNotifications = append(filteredNotifications, n)
		}

	} else {
		filteredNotifications = allNotifications
	}

	return notificationsDescriptionAttribute(d, filteredNotifications, meta)
}

func notificationsDescriptionAttribute(d *schema.ResourceData, notifications []interface{}, meta interface{}) error {
	var ids []string
	var s = make([]map[string]interface{}, 0)
	for _, n := range notifications {
		var object map[string]interface{}
		object = n.(map[string]interface{})
		mapping := map[string]interface{}{
			"notification_arn": object["NotificationArn"].(string),
			"scaling_group_id": object["ScalingGroupId"].(string),
		}
		var notificationTypes []string
		if object["NotificationTypes"] != nil && len(object["NotificationTypes"].(map[string]interface{})["NotificationType"].([]interface{})) > 0 {
			for _, v := range object["NotificationTypes"].(map[string]interface{})["NotificationType"].([]interface{}) {
				notificationTypes = append(notificationTypes, v.(string))
			}
			mapping["notification_types"] = notificationTypes
		}
		if object["TimeZone"] != nil {
			mapping["time_zone"] = object["TimeZone"].(string)
		}
		if object["MessageEncoding"] != nil {
			mapping["message_encoding"] = object["MessageEncoding"].(string)
		}
		ids = append(ids, object["NotificationArn"].(string))
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("notifications", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
