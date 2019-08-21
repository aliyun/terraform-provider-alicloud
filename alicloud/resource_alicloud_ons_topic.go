package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
				ValidateFunc: validateOnsTopic,
			},
			"message_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOnsTopicMessageType,
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOnsTopicRemark,
			},
			"perm": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateOnsTopicPerm,
			},
		},
	}
}

func resourceAlicloudOnsTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	instanceId := d.Get("instance_id").(string)
	topic := d.Get("topic").(string)

	request := ons.CreateOnsTopicCreateRequest()
	request.RegionId = client.RegionId
	request.Topic = topic
	request.InstanceId = instanceId
	request.MessageType = requests.NewInteger(d.Get("message_type").(int))
	request.PreventCache = onsService.GetPreventCache()

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = v.(string)
	}

	raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsTopicCreate(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ons_topic", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(instanceId + ":" + topic)

	if err = onsService.WaitForOnsTopic(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAlicloudOnsTopicRead(d, meta)
}

func resourceAlicloudOnsTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	object, err := onsService.DescribeOnsTopic(d.Id())
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
	d.Set("message_type", object.MessageType)
	d.Set("remark", object.Remark)

	return nil
}

func resourceAlicloudOnsTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := ons.CreateOnsTopicUpdateRequest()
	request.RegionId = client.RegionId
	request.InstanceId = instanceId
	request.Topic = topic
	request.PreventCache = onsService.GetPreventCache()

	var perm int
	if d.HasChange("perm") {
		perm = d.Get("perm").(int)
		request.Perm = requests.NewInteger(perm)
		raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsTopicUpdate(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlicloudOnsTopicRead(d, meta)
}

func resourceAlicloudOnsTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := ons.CreateOnsTopicDeleteRequest()
	request.RegionId = client.RegionId
	request.Topic = topic
	request.InstanceId = instanceId
	request.PreventCache = onsService.GetPreventCache()

	raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsTopicDelete(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{AuthResourceOwnerError}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(onsService.WaitForOnsTopic(d.Id(), Deleted, DefaultTimeoutMedium))
}
