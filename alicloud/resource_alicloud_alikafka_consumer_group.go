// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAliKafkaConsumerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAliKafkaConsumerGroupCreate,
		Read:   resourceAliCloudAliKafkaConsumerGroupRead,
		Update: resourceAliCloudAliKafkaConsumerGroupUpdate,
		Delete: resourceAliCloudAliKafkaConsumerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"consumer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"description": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "Field `description` has been deprecated from provider version 1.268.0. New field `remark` instead.",
			},
		},
	}
}

func resourceAliCloudAliKafkaConsumerGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateConsumerGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("consumer_id"); ok {
		request["ConsumerId"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("remark"); ok {
		request["Remark"] = v
	} else if v, ok := d.GetOk("description"); ok {
		request["Remark"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_consumer_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], request["ConsumerId"]))

	return resourceAliCloudAliKafkaConsumerGroupRead(d, meta)
}

func resourceAliCloudAliKafkaConsumerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaServiceV2 := AlikafkaServiceV2{client}

	objectRaw, err := alikafkaServiceV2.DescribeAliKafkaConsumerGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alikafka_consumer_group DescribeAliKafkaConsumerGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("remark", objectRaw["Remark"])
	d.Set("consumer_id", objectRaw["ConsumerId"])
	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("description", objectRaw["Remark"])

	tagsMaps, _ := jsonpath.Get("$.Tags.TagVO", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudAliKafkaConsumerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("tags") {
		alikafkaServiceV2 := AlikafkaServiceV2{client}
		if err := alikafkaServiceV2.SetResourceTags(d, "CONSUMERGROUP"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudAliKafkaConsumerGroupRead(d, meta)
}

func resourceAliCloudAliKafkaConsumerGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteConsumerGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["ConsumerId"] = parts[1]
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

	aliKafkaServiceV2 := AlikafkaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, aliKafkaServiceV2.AliKafkaConsumerGroupStateRefreshFunc(d.Id(), "ConsumerId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
