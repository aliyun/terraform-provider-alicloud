package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudOnsTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOnsTopicCreate,
		Read:   resourceAlicloudOnsTopicRead,
		Update: resourceAlicloudOnsTopicUpdate,
		Delete: resourceAlicloudOnsTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"message_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 4, 5}),
			},
			"perm": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{2, 4, 6}),
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"tags": tagsSchema(),
			"topic_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"topic"},
				ValidateFunc:  validation.StringLenBetween(1, 64),
			},
			"topic": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'topic' has been deprecated from version 1.97.0. Use 'topic_name' instead.",
				ConflictsWith: []string{"topic_name"},
				ValidateFunc:  validation.StringLenBetween(1, 64),
			},
		},
	}
}

func resourceAlicloudOnsTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ons.CreateOnsTopicCreateRequest()
	request.InstanceId = d.Get("instance_id").(string)
	request.MessageType = requests.NewInteger(d.Get("message_type").(int))
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = v.(string)
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request.Topic = v.(string)
	} else if v, ok := d.GetOk("topic"); ok {
		request.Topic = v.(string)
	} else {
		return WrapError(Error(`[ERROR] Argument "topic" or "topic_name" must be set one!`))
	}

	wait := incrementalWait(3*time.Second, 10*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsTopicCreate(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.User"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		d.SetId(fmt.Sprintf("%v:%v", request.InstanceId, request.Topic))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ons_topic", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudOnsTopicUpdate(d, meta)
}
func resourceAlicloudOnsTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}
	object, err := onsService.DescribeOnsTopic(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ons_topic onsService.DescribeOnsTopic Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("instance_id", parts[0])
	d.Set("topic_name", parts[1])
	d.Set("topic", parts[1])
	d.Set("message_type", object.MessageType)
	d.Set("remark", object.Remark)

	tags := make(map[string]string)
	for _, t := range object.Tags.Tag {
		tags[t.Key] = t.Value
	}
	d.Set("tags", tags)

	onsTopicStatusObject, err := onsService.OnsTopicStatus(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("perm", onsTopicStatusObject.Perm)
	return nil
}
func resourceAlicloudOnsTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := onsService.SetResourceTagsForTopic(d, "TOPIC"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	if d.HasChange("perm") {
		request := ons.CreateOnsTopicUpdateRequest()
		request.InstanceId = parts[0]
		request.Topic = parts[1]
		request.Perm = requests.NewInteger(d.Get("perm").(int))
		raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsTopicUpdate(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("perm")
	}
	d.Partial(false)
	return resourceAlicloudOnsTopicRead(d, meta)
}
func resourceAlicloudOnsTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := ons.CreateOnsTopicDeleteRequest()
	request.InstanceId = parts[0]
	request.Topic = parts[1]
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsTopicDelete(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.User"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
