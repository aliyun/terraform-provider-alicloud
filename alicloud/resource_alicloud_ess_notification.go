package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		},
	}
}

func resourceAlicloudEssNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateNotificationConfiguration"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	conn, err := client.NewEssClient()
	request["ScalingGroupId"] = d.Get("scaling_group_id").(string)
	request["NotificationArn"] = d.Get("notification_arn").(string)
	if v, ok := d.GetOk("time_zone"); ok {
		request["TimeZone"] = v
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
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	d.SetId(fmt.Sprintf("%s:%s", request["ScalingGroupId"], request["NotificationArn"]))
	return resourceAlicloudEssNotificationRead(d, meta)
}

func resourceAlicloudEssNotificationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	object, err := essService.DescribeEssNotification(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("scaling_group_id", object["ScalingGroupId"])
	d.Set("notification_arn", object["NotificationArn"])
	d.Set("time_zone", object["TimeZone"])
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
	conn, err := client.NewEssClient()
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
	runtime := util.RuntimeOptions{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	runtime.SetAutoretry(true)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
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
		if IsExpectedErrors(err, []string{"NotificationConfigurationNotExist", "InvalidScalingGroupId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(essService.WaitForEssNotification(d.Id(), Deleted, DefaultTimeout))
}
