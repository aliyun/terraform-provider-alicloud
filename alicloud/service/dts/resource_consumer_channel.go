package dts

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceConsumerChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsumerChannelCreate,
		Read:   resourceConsumerChannelRead,
		Update: resourceConsumerChannelUpdate,
		Delete: resourceConsumerChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"consumer_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"consumer_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,

				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"consumer_group_password": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(8, 32),
			},
			"consumer_group_user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,

				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-zA-Z0-9_]{1,16}$"), "The length of the name is limited to `1` to `16` characters. It can contain uppercase and lowercase letters, numbers, underscores (_)"),
			},
			"dts_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConsumerChannelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateConsumerChannel"
	request := make(map[string]interface{})
	var err error
	request["ConsumerGroupName"] = d.Get("consumer_group_name")
	request["ConsumerGroupPassword"] = d.Get("consumer_group_password")
	request["ConsumerGroupUserName"] = d.Get("consumer_group_user_name")
	request["DtsInstanceId"] = d.Get("dts_instance_id")

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"OperationDenied.JobStatus"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_dts_consumer_channel", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DtsInstanceId"], ":", response["ConsumerGroupID"]))

	return resourceConsumerChannelRead(d, meta)
}
func resourceConsumerChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := NewService(client)
	object, err := dtsService.DescribeDtsConsumerChannel(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dts_consumer_channel dtsService.DescribeDtsConsumerChannel Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("consumer_group_id", parts[1])
	d.Set("dts_instance_id", parts[0])
	d.Set("consumer_group_name", object["ConsumerGroupName"])
	d.Set("consumer_group_user_name", object["ConsumerGroupUserName"])
	return nil
}
func resourceConsumerChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"ConsumerGroupId": parts[1],
		"DtsInstanceId":   parts[0],
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("consumer_group_password") {
		update = true
		request["ConsumerGroupPassword"] = d.Get("consumer_group_password")
	}
	if update {
		action := "ModifyConsumerChannel"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"OperationDenied.JobStatus"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	return resourceConsumerChannelRead(d, meta)
}
func resourceConsumerChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeleteConsumerChannel"
	var response map[string]interface{}
	request := map[string]interface{}{
		"ConsumerGroupId": parts[1],
		"DtsInstanceId":   parts[0],
	}

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"OperationDenied.JobStatus"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	time.Sleep(10 * time.Second)
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
