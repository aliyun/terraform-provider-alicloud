package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudConfigDelivery() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudConfigDeliveryCreate,
		Read:   resourceAlicloudConfigDeliveryRead,
		Update: resourceAlicloudConfigDeliveryUpdate,
		Delete: resourceAlicloudConfigDeliveryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"configuration_item_change_notification": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"configuration_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("delivery_channel_type"); ok && v.(string) == "OSS" {
						return false
					}
					return true
				},
			},
			"delivery_channel_condition": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("delivery_channel_type"); ok && v.(string) == "MNS" {
						return false
					}
					return true
				},
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MNS", "OSS", "SLS"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"non_compliant_notification": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("delivery_channel_type"); ok && (v.(string) == "MNS" || v.(string) == "SLS") {
						return false
					}
					return true
				},
			},
			"oversized_data_oss_target_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
		},
	}
}

func resourceAlicloudConfigDeliveryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateConfigDeliveryChannel"
	request := make(map[string]interface{})
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("configuration_item_change_notification"); ok {
		request["ConfigurationItemChangeNotification"] = v
	}
	if v, ok := d.GetOkExists("configuration_snapshot"); ok {
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
	if v, ok := d.GetOkExists("non_compliant_notification"); ok {
		request["NonCompliantNotification"] = v
	}
	if v, ok := d.GetOk("oversized_data_oss_target_arn"); ok {
		request["OversizedDataOSSTargetArn"] = v
	}
	request["ClientToken"] = buildClientToken("CreateConfigDeliveryChannel")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &runtime)
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
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_delivery", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DeliveryChannelId"]))

	return resourceAlicloudConfigDeliveryUpdate(d, meta)
}
func resourceAlicloudConfigDeliveryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	object, err := configService.DescribeConfigDelivery(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_delivery configService.DescribeConfigDelivery Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("configuration_item_change_notification", object["ConfigurationItemChangeNotification"])
	d.Set("configuration_snapshot", object["ConfigurationSnapshot"])
	d.Set("delivery_channel_condition", object["DeliveryChannelCondition"])
	d.Set("delivery_channel_name", object["DeliveryChannelName"])
	d.Set("delivery_channel_target_arn", object["DeliveryChannelTargetArn"])
	d.Set("delivery_channel_type", object["DeliveryChannelType"])
	d.Set("description", object["Description"])
	d.Set("non_compliant_notification", object["NonCompliantNotification"])
	d.Set("oversized_data_oss_target_arn", object["OversizedDataOSSTargetArn"])
	d.Set("status", formatInt(object["Status"]))
	return nil
}
func resourceAlicloudConfigDeliveryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"DeliveryChannelId": d.Id(),
	}
	if v, ok := d.GetOkExists("configuration_item_change_notification"); ok {
		request["ConfigurationItemChangeNotification"] = v
	}
	if !d.IsNewResource() && d.HasChange("configuration_item_change_notification") {
		update = true
	}
	if v, ok := d.GetOkExists("configuration_snapshot"); ok {
		request["ConfigurationSnapshot"] = v
	}
	if !d.IsNewResource() && d.HasChange("configuration_snapshot") {
		update = true
	}
	if v, ok := d.GetOk("delivery_channel_condition"); ok {
		request["DeliveryChannelCondition"] = v
	}
	if !d.IsNewResource() && d.HasChange("delivery_channel_condition") {
		update = true
	}
	if v, ok := d.GetOk("delivery_channel_name"); ok {
		request["DeliveryChannelName"] = v
	}
	if !d.IsNewResource() && d.HasChange("delivery_channel_name") {
		update = true
	}
	if !d.IsNewResource() && d.HasChange("delivery_channel_target_arn") {
		update = true
	}
	request["DeliveryChannelTargetArn"] = d.Get("delivery_channel_target_arn")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOkExists("non_compliant_notification"); ok {
		request["NonCompliantNotification"] = v
	}
	if !d.IsNewResource() && d.HasChange("non_compliant_notification") {
		update = true
	}
	if v, ok := d.GetOk("oversized_data_oss_target_arn"); ok {
		request["OversizedDataOSSTargetArn"] = v
	}
	if !d.IsNewResource() && d.HasChange("oversized_data_oss_target_arn") {
		update = true
	}
	if v, ok := d.GetOkExists("status"); ok {
		request["Status"] = v
	}
	if d.HasChange("status") {
		update = true
	}
	if update {
		action := "UpdateConfigDeliveryChannel"
		request["ClientToken"] = buildClientToken("UpdateConfigDeliveryChannel")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &runtime)
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
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudConfigDeliveryRead(d, meta)
}
func resourceAlicloudConfigDeliveryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteConfigDeliveryChannel"
	var response map[string]interface{}
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DeliveryChannelId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
