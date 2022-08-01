package alicloud

import (
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func resourceAlicloudMNSTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMNSTopicCreate,
		Read:   resourceAlicloudMNSTopicRead,
		Update: resourceAlicloudMNSTopicUpdate,
		Delete: resourceAlicloudMNSTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 256),
			},
			"maximum_message_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      65536,
				ValidateFunc: validation.IntBetween(1024, 65536),
			},
			"logging_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAlicloudMNSTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("name"); ok {
		request["TopicName"] = v
	}
	if v, ok := d.GetOkExists("maximum_message_size"); ok {
		request["MaxMessageSize"] = v
	}
	if v, ok := d.GetOkExists("logging_enabled"); ok {
		request["EnableLogging"] = v
	}

	var response map[string]interface{}
	action := "CreateTopic"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mns_topic", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(d.Get("name").(string))

	return resourceAlicloudMNSTopicRead(d, meta)
}

func resourceAlicloudMNSTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mnsService := MnsService{client}

	object, err := mnsService.DescribeMessageServiceTopic(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object["TopicName"])
	d.Set("maximum_message_size", object["MaxMessageSize"])
	d.Set("logging_enabled", object["LoggingEnabled"])
	return nil
}

func resourceAlicloudMNSTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TopicName": d.Id(),
	}

	attributeUpdate := false

	if !d.IsNewResource() && d.HasChange("logging_enabled") {
		attributeUpdate = true
		if v, ok := d.GetOkExists("logging_enabled"); ok {
			request["EnableLogging"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("maximum_message_size") {
		attributeUpdate = true
		if v, ok := d.GetOkExists("maximum_message_size"); ok {
			request["MaxMessageSize"] = v
		}
	}

	if attributeUpdate {
		action := "SetTopicAttributes"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudMNSTopicRead(d, meta)
}

func resourceAlicloudMNSTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"TopicName": d.Id(),
	}

	action := "DeleteTopic"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
