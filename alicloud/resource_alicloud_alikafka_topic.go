// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudAlikafkaTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlikafkaTopicCreate,
		Read:   resourceAliCloudAlikafkaTopicRead,
		Update: resourceAliCloudAlikafkaTopicUpdate,
		Delete: resourceAliCloudAlikafkaTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(16 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"compact_topic": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"configs": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"local_topic": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"partition_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"topic": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudAlikafkaTopicCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTopic"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Get("instance_id")
	request["Topic"] = d.Get("topic")
	request["RegionId"] = client.RegionId

	request["Remark"] = d.Get("remark")
	if v, ok := d.GetOkExists("compact_topic"); ok {
		request["CompactTopic"] = v
	}
	if v, ok := d.GetOkExists("partition_num"); ok {
		request["PartitionNum"] = v
	}
	if v, ok := d.GetOkExists("local_topic"); ok {
		request["LocalTopic"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("configs"); ok {
		request["Config"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_topic", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], request["Topic"]))

	alikafkaServiceV2 := AlikafkaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"0"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, alikafkaServiceV2.AlikafkaTopicStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAlikafkaTopicRead(d, meta)
}

func resourceAliCloudAlikafkaTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaServiceV2 := AlikafkaServiceV2{client}

	objectRaw, err := alikafkaServiceV2.DescribeAlikafkaTopic(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alikafka_topic DescribeAlikafkaTopic Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("compact_topic", objectRaw["CompactTopic"])
	d.Set("configs", objectRaw["TopicConfig"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("local_topic", objectRaw["LocalTopic"])
	d.Set("partition_num", objectRaw["PartitionNum"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("remark", objectRaw["Remark"])
	d.Set("status", objectRaw["Status"])
	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("topic", objectRaw["Topic"])

	tagsMaps, _ := jsonpath.Get("$.Tags.TagVO", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudAlikafkaTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyTopicRemark"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["Topic"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("remark") {
		update = true
	}
	request["Remark"] = d.Get("remark")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("alikafka", "2019-09-16", action, query, request, true)
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "UpdateTopicConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["Topic"] = parts[1]
	request["RegionId"] = client.RegionId

	// In the UpdateTopicConfig, the Config and Value are required; For Terraform, if Config, Value and Configs are set simultaneously, only Configs takes effect.
	request["Config"] = "skipConfig"
	request["Value"] = "skipValue"

	if d.HasChange("configs") {
		update = true
	}
	if v, ok := d.GetOk("configs"); ok {
		request["Configs"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("alikafka", "2019-09-16", action, query, request, true)
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

	if d.HasChange("partition_num") {
		o, n := d.GetChange("partition_num")
		oldPartitionNum := o.(int)
		newPartitionNum := n.(int)

		if newPartitionNum < oldPartitionNum {
			return WrapError(errors.New("partition_num only support adjust to a greater value"))
		} else {
			action = "ModifyPartitionNum"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["InstanceId"] = parts[0]
			request["Topic"] = parts[1]
			request["RegionId"] = client.RegionId
			request["AddPartitionNum"] = newPartitionNum - oldPartitionNum

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("alikafka", "2019-09-16", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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
	}

	if d.HasChange("tags") {
		alikafkaServiceV2 := AlikafkaServiceV2{client}
		if err := alikafkaServiceV2.SetResourceTags(d, "TOPIC"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudAlikafkaTopicRead(d, meta)
}

func resourceAliCloudAlikafkaTopicDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteTopic"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["Topic"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	alikafkaServiceV2 := AlikafkaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, alikafkaServiceV2.AlikafkaTopicStateRefreshFunc(d.Id(), "Topic", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
