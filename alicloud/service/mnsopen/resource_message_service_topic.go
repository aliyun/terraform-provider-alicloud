package mnsopen

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	tftags "github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	tfvalidators "github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudMessageServiceTopic() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAliCloudMessageServiceTopicCreate,
		Read:   ResourceAliCloudMessageServiceTopicRead,
		Update: ResourceAliCloudMessageServiceTopicUpdate,
		Delete: ResourceAliCloudMessageServiceTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_logging": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"logging_enabled"},
			},
			"max_message_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tfvalidators.IntBetween(1024, 65536),
			},
			"enable_sse": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"sse_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tfvalidators.StringInSlice([]string{"SMQ", "KMS"}, false),
			},
			"sse_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tfvalidators.StringInSlice([]string{"AES-256-GCM"}, false),
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tftags.TagsSchemaWithElements(),
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: tfvalidators.StringInSlice([]string{"normal", "fifo"}, false),
			},
			"logging_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field `logging_enabled` has been deprecated from provider version 1.241.0. New field `enable_logging` instead.",
			},
		},
	}
}

func ResourceAliCloudMessageServiceTopicCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTopic"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TopicName"] = d.Get("topic_name")

	if v, ok := d.GetOkExists("enable_logging"); ok {
		request["EnableLogging"] = v
	} else if v, ok := d.GetOkExists("logging_enabled"); ok {
		request["EnableLogging"] = v
	}
	if v, ok := d.GetOkExists("max_message_size"); ok {
		request["MaxMessageSize"] = v
	}
	// A new topic has SSE off by default, so EnableSSE is only sent when
	// true. d.Get is used instead of GetOkExists, which may treat an
	// existing value as non-existent.
	if v := d.Get("enable_sse"); v.(bool) {
		request["EnableSSE"] = v
	}
	if v, ok := d.GetOk("sse_type"); ok {
		request["SseType"] = v
	}
	if v, ok := d.GetOk("sse_algorithm"); ok {
		request["SseAlgorithm"] = v
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KmsKeyId"] = v
	}
	if v, ok := d.GetOk("topic_type"); ok {
		request["TopicType"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tftags.ConvertTags(v.(map[string]interface{}))
		request = tftags.ExpandTagsToMap(request, tagsMap)
	}

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)
		if err != nil {
			if tferrs.IsExpectedErrors(err, []string{"KMSServiceUnavailable"}) || tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)

	if err != nil {
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, "alicloud_message_service_topic", action, tferrs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["TopicName"]))

	return ResourceAliCloudMessageServiceTopicRead(d, meta)
}

func ResourceAliCloudMessageServiceTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	messageServiceServiceV2 := MessageServiceServiceV2{client}

	objectRaw, err := messageServiceServiceV2.DescribeMessageServiceTopic(d.Id())
	if err != nil {
		if tferrs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_message_service_topic DescribeMessageServiceTopic Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return tferrs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("enable_logging", objectRaw["LoggingEnabled"])
	d.Set("max_message_size", objectRaw["MaxMessageSize"])
	d.Set("enable_sse", objectRaw["EnableSSE"])
	d.Set("sse_type", objectRaw["SseType"])
	d.Set("sse_algorithm", objectRaw["SseAlgorithm"])
	d.Set("kms_key_id", objectRaw["KmsKeyId"])
	d.Set("topic_name", objectRaw["TopicName"])
	d.Set("topic_type", objectRaw["TopicType"])
	d.Set("logging_enabled", objectRaw["LoggingEnabled"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tftags.TagsToMap(tagsMaps))

	return nil
}

func ResourceAliCloudMessageServiceTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	action := "SetTopicAttributes"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TopicName"] = d.Id()

	if d.HasChange("enable_logging") {
		update = true

		if v, ok := d.GetOkExists("enable_logging"); ok {
			request["EnableLogging"] = v
		}
	}

	if d.HasChange("logging_enabled") {
		update = true

		if v, ok := d.GetOkExists("logging_enabled"); ok {
			request["EnableLogging"] = v
		}
	}

	if d.HasChange("max_message_size") {
		update = true

		if v, ok := d.GetOkExists("max_message_size"); ok {
			request["MaxMessageSize"] = v
		}
	}

	if d.HasChange("enable_sse") || d.HasChange("sse_type") || d.HasChange("sse_algorithm") || d.HasChange("kms_key_id") {
		update = true

		// The whole SSE family is sent together on any SSE change. EnableSSE
		// is set unconditionally via d.Get (not GetOkExists, which may treat
		// an existing value as non-existent), so the current switch value is
		// always sent, including the explicit false that disables SSE.
		request["EnableSSE"] = d.Get("enable_sse")

		if v, ok := d.GetOk("sse_type"); ok {
			request["SseType"] = v
		}

		if v, ok := d.GetOk("sse_algorithm"); ok {
			request["SseAlgorithm"] = v
		}

		if v, ok := d.GetOk("kms_key_id"); ok {
			request["KmsKeyId"] = v
		}
	}

	if update {
		wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)
			if err != nil {
				if tferrs.IsExpectedErrors(err, []string{"KMSServiceUnavailable"}) || tfhelper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		tflogs.AddDebug(action, response, request)
		if err != nil {
			return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		messageServiceServiceV2 := MessageServiceServiceV2{client}
		if err := messageServiceServiceV2.SetResourceTags(d, "topic"); err != nil {
			return tferrs.WrapError(err)
		}
	}

	return ResourceAliCloudMessageServiceTopicRead(d, meta)
}

func ResourceAliCloudMessageServiceTopicDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTopic"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TopicName"] = d.Id()

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)

		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)

	if err != nil {
		if tferrs.NotFoundError(err) {
			return nil
		}
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
