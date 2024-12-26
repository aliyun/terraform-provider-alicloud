// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudConfigDelivery() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudConfigDeliveryCreate,
		Read:   resourceAliCloudConfigDeliveryRead,
		Update: resourceAliCloudConfigDeliveryUpdate,
		Delete: resourceAliCloudConfigDeliveryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"configuration_item_change_notification": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"configuration_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delivery_channel_condition": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delivery_channel_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delivery_channel_target_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delivery_channel_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"non_compliant_notification": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"oversized_data_oss_target_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntInSlice([]int{1, 0}),
			},
		},
	}
}

func resourceAliCloudConfigDeliveryCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateConfigDeliveryChannel"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]interface{})

	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("configuration_item_change_notification"); ok {
		request["ConfigurationItemChangeNotification"] = v
	}
	if v, ok := d.GetOk("configuration_snapshot"); ok {
		request["ConfigurationSnapshot"] = v
	}
	if v, ok := d.GetOk("delivery_channel_condition"); ok {
		request["DeliveryChannelCondition"] = v
	}
	if v, ok := d.GetOk("delivery_channel_name"); ok {
		request["DeliveryChannelName"] = v
	}
	request["DeliveryChannelTargetArn"] = d.Get("delivery_channel_target_arn")
	request["DeliveryChannelType"] = d.Get("delivery_channel_type")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("non_compliant_notification"); ok {
		request["NonCompliantNotification"] = v
	}
	if v, ok := d.GetOk("oversized_data_oss_target_arn"); ok {
		request["OversizedDataOSSTargetArn"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_delivery", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DeliveryChannelId"]))

	return resourceAliCloudConfigDeliveryUpdate(d, meta)
}

func resourceAliCloudConfigDeliveryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configServiceV2 := ConfigServiceV2{client}

	objectRaw, err := configServiceV2.DescribeConfigDelivery(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_delivery DescribeConfigDelivery Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["ConfigurationItemChangeNotification"] != nil {
		d.Set("configuration_item_change_notification", objectRaw["ConfigurationItemChangeNotification"])
	}
	if objectRaw["ConfigurationSnapshot"] != nil {
		d.Set("configuration_snapshot", objectRaw["ConfigurationSnapshot"])
	}
	if objectRaw["DeliveryChannelCondition"] != nil {
		d.Set("delivery_channel_condition", objectRaw["DeliveryChannelCondition"])
	}
	if objectRaw["DeliveryChannelName"] != nil {
		d.Set("delivery_channel_name", objectRaw["DeliveryChannelName"])
	}
	if objectRaw["DeliveryChannelTargetArn"] != nil {
		d.Set("delivery_channel_target_arn", objectRaw["DeliveryChannelTargetArn"])
	}
	if objectRaw["DeliveryChannelType"] != nil {
		d.Set("delivery_channel_type", objectRaw["DeliveryChannelType"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["NonCompliantNotification"] != nil {
		d.Set("non_compliant_notification", objectRaw["NonCompliantNotification"])
	}
	if objectRaw["OversizedDataOSSTargetArn"] != nil {
		d.Set("oversized_data_oss_target_arn", objectRaw["OversizedDataOSSTargetArn"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", formatInt(objectRaw["Status"]))
	}

	return nil
}

func resourceAliCloudConfigDeliveryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configServiceV2 := ConfigServiceV2{client}
	objectRaw, err := configServiceV2.DescribeConfigDelivery(d.Id())
	if err != nil {
		return WrapError(err)
	}
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateConfigDeliveryChannel"

	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DeliveryChannelId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("configuration_item_change_notification") {
		update = true
		request["ConfigurationItemChangeNotification"] = d.Get("configuration_item_change_notification")
	}

	if !d.IsNewResource() && d.HasChange("configuration_snapshot") {
		update = true
		request["ConfigurationSnapshot"] = d.Get("configuration_snapshot")
	}

	if !d.IsNewResource() && d.HasChange("delivery_channel_condition") {
		update = true
		request["DeliveryChannelCondition"] = d.Get("delivery_channel_condition")
	}

	if !d.IsNewResource() && d.HasChange("delivery_channel_name") {
		update = true
		request["DeliveryChannelName"] = d.Get("delivery_channel_name")
	}

	if !d.IsNewResource() && d.HasChange("delivery_channel_target_arn") {
		update = true
	}
	request["DeliveryChannelTargetArn"] = d.Get("delivery_channel_target_arn")
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("non_compliant_notification") {
		update = true
		request["NonCompliantNotification"] = d.Get("non_compliant_notification")
	}

	if !d.IsNewResource() && d.HasChange("oversized_data_oss_target_arn") {
		update = true
		request["OversizedDataOSSTargetArn"] = d.Get("oversized_data_oss_target_arn")
	}

	// status default to 1 after creating
	if v, ok := d.GetOkExists("status"); ok && fmt.Sprint(objectRaw["Status"]) != fmt.Sprint(v) {
		update = true
		request["Status"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudConfigDeliveryRead(d, meta)
}

func resourceAliCloudConfigDeliveryDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteConfigDeliveryChannel"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]interface{})

	request = make(map[string]interface{})
	query["DeliveryChannelId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, false)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
