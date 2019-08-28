package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudOnsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOnsGroupCreate,
		Read:   resourceAlicloudOnsGroupRead,
		Update: resourceAlicloudOnsGroupUpdate,
		Delete: resourceAlicloudOnsGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOnsGroupId,
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOnsGroupRemark,
			},
			"read_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudOnsGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	instanceId := d.Get("instance_id").(string)
	groupId := d.Get("group_id").(string)

	request := ons.CreateOnsGroupCreateRequest()
	request.RegionId = client.RegionId
	request.GroupId = groupId
	request.InstanceId = instanceId
	request.PreventCache = onsService.GetPreventCache()

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = v.(string)
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsGroupCreate(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{OnsThrottlingUser}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ons_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(instanceId + ":" + groupId)

	if err = onsService.WaitForOnsGroup(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAlicloudOnsGroupRead(d, meta)
}

func resourceAlicloudOnsGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	object, err := onsService.DescribeOnsGroup(d.Id())

	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("group_id", object.GroupId)
	d.Set("remark", object.Remark)

	return nil
}

func resourceAlicloudOnsGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	groupId := parts[1]

	request := ons.CreateOnsGroupConsumerUpdateRequest()
	request.RegionId = client.RegionId
	request.InstanceId = instanceId
	request.GroupId = groupId
	request.PreventCache = onsService.GetPreventCache()

	if d.HasChange("read_enable") {
		readEnable := d.Get("read_enable").(bool)
		request.ReadEnable = requests.NewBoolean(readEnable)
		raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsGroupConsumerUpdate(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlicloudOnsGroupRead(d, meta)
}

func resourceAlicloudOnsGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	groupId := parts[1]

	request := ons.CreateOnsGroupDeleteRequest()
	request.RegionId = client.RegionId
	request.InstanceId = instanceId
	request.GroupId = groupId
	request.PreventCache = onsService.GetPreventCache()

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsGroupDelete(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{OnsThrottlingUser}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, []string{AuthResourceOwnerError}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(onsService.WaitForOnsGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}
