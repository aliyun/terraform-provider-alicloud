package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudAlikafkaTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlikafkaTopicCreate,
		Read:   resourceAlicloudAlikafkaTopicRead,
		Delete: resourceAlicloudAlikafkaTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAlikafkaStringLen,
			},
			"local_topic": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"compact_topic": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"partition_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAlikafkaPartitionNum,
			},
			"remark": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAlikafkaStringLen,
			},
		},
	}
}

func resourceAlicloudAlikafkaTopicCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	instanceId := d.Get("instance_id").(string)
	regionId := client.RegionId
	topic := d.Get("topic").(string)

	request := alikafka.CreateCreateTopicRequest()
	request.InstanceId = instanceId
	request.RegionId = regionId
	request.Topic = topic
	if v, ok := d.GetOk("local_topic"); ok {
		request.LocalTopic = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("compact_topic"); ok {
		request.CompactTopic = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("partition_num"); ok {
		request.PartitionNum = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = v.(string)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.CreateTopic(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AlikafkaThrottlingUser}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_topic", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(instanceId + ":" + topic)
	return resourceAlicloudAlikafkaTopicRead(d, meta)
}

func resourceAlicloudAlikafkaTopicRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	object, err := alikafkaService.DescribeAlikafkaTopic(d.Id())
	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("topic", object.Topic)
	d.Set("local_topic", object.LocalTopic)
	d.Set("compact_topic", object.CompactTopic)
	d.Set("partition_num", object.PartitionNum)
	d.Set("remark", object.Remark)

	return nil
}

func resourceAlicloudAlikafkaTopicDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := alikafka.CreateDeleteTopicRequest()
	request.Topic = topic
	request.InstanceId = instanceId
	request.RegionId = client.RegionId

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DeleteTopic(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AlikafkaThrottlingUser}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(alikafkaService.WaitForAlikafkaTopic(d.Id(), Deleted, DefaultTimeoutMedium))
}
