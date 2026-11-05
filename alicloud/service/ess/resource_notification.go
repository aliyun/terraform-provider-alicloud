package ess

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
)

func resourceAlicloudEssNotification() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEssNotificationCreate,
		Read:   resourceAlicloudEssNotificationRead,
		Update: resourceAlicloudEssNotificationUpdate,
		Delete: resourceAlicloudEssNotificationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"notification_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"notification_types": {
				Required: true,
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"message_encoding": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"PlainText", "Base64"}, false),
			},
		},
	}
}

func resourceAlicloudEssNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateNotificationConfiguration"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	var err error
	request["ScalingGroupId"] = d.Get("scaling_group_id").(string)
	request["NotificationArn"] = d.Get("notification_arn").(string)
	if v, ok := d.GetOk("time_zone"); ok {
		request["TimeZone"] = v
	}
	if v, ok := d.GetOk("message_encoding"); ok {
		request["MessageEncoding"] = v
	}
	if v, ok := d.GetOk("notification_types"); ok {
		notificationTypes := make([]string, 0)
		notificationTypeList := v.(*schema.Set).List()
		if len(notificationTypeList) > 0 {
			for _, n := range notificationTypeList {
				notificationTypes = append(notificationTypes, n.(string))
			}
		}
		if len(notificationTypes) > 0 {
			request["NotificationType"] = notificationTypes
		}
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ess", "2014-08-28", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%s:%s", request["ScalingGroupId"], request["NotificationArn"]))
	return resourceAlicloudEssNotificationRead(d, meta)
}

func resourceAlicloudEssNotificationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	object, err := essService.DescribeEssNotification(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("scaling_group_id", object["ScalingGroupId"])
	d.Set("notification_arn", object["NotificationArn"])
	d.Set("time_zone", object["TimeZone"])
	d.Set("message_encoding", object["MessageEncoding"])
	notificationTypes, _ := jsonpath.Get("$.NotificationTypes", object)
	notificationType, _ := jsonpath.Get("$.NotificationType", notificationTypes)
	d.Set("notification_types", notificationType)
	return nil
}

func resourceAlicloudEssNotificationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyNotificationConfiguration"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	var err error
	request["RegionId"] = client.RegionId
	parts := strings.SplitN(d.Id(), ":", 2)
	request["ScalingGroupId"] = parts[0]
	request["NotificationArn"] = parts[1]
	v := d.Get("notification_types")
	notificationTypes := make([]string, 0)
	notificationTypeList := v.(*schema.Set).List()
	if len(notificationTypeList) > 0 {
		for _, n := range notificationTypeList {
			notificationTypes = append(notificationTypes, n.(string))
		}
	}
	if len(notificationTypes) > 0 {
		request["NotificationType"] = notificationTypes
	}
	if d.HasChange("time_zone") {
		if v, ok := d.GetOk("time_zone"); ok {
			request["TimeZone"] = v
		}
	}
	if d.HasChange("message_encoding") {
		if v, ok := d.GetOk("message_encoding"); ok {
			request["MessageEncoding"] = v
		}
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("Ess", "2014-08-28", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	return resourceAlicloudEssNotificationRead(d, meta)
}

func resourceAlicloudEssNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	request := ess.CreateDeleteNotificationConfigurationRequest()
	request.RegionId = client.RegionId
	parts := strings.SplitN(d.Id(), ":", 2)

	request.ScalingGroupId = parts[0]
	request.NotificationArn = parts[1]

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteNotificationConfiguration(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotificationConfigurationNotExist", "InvalidScalingGroupId.NotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errs.WrapError(essService.WaitForEssNotification(d.Id(), Deleted, DefaultTimeout))
}
