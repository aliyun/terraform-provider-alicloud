package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCenInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenInstanceAttachmentCreate,
		Read:   resourceAlicloudCenInstanceAttachmentRead,
		Delete: resourceAlicloudCenInstanceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"child_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"child_instance_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"child_instance_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCenInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	cenId := d.Get("instance_id").(string)
	instanceId := d.Get("child_instance_id").(string)
	instanceRegionId := d.Get("child_instance_region_id").(string)

	instanceType, err := GetCenChildInstanceType(instanceId)
	if err != nil {
		return WrapError(err)
	}

	request := cbn.CreateAttachCenChildInstanceRequest()
	request.RegionId = client.RegionId
	request.CenId = cenId
	request.ChildInstanceId = instanceId
	request.ChildInstanceType = instanceType
	request.ChildInstanceRegionId = instanceRegionId
	if v := d.Get("child_instance_owner_id").(string); v != "" {
		request.ChildInstanceOwnerId = requests.Integer(v)
	}
	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.AttachCenChildInstance(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidCenInstanceStatus, InvalidChildInstanceStatus}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_instance_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(cenId + COLON_SEPARATED + instanceId)

	if err := cenService.WaitForCenInstanceAttachment(d.Id(), Status("Attached"), DefaultCenTimeoutLong); err != nil {
		return WrapError(err)
	}
	return resourceAlicloudCenInstanceAttachmentRead(d, meta)
}

func resourceAlicloudCenInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}

	object, err := cenService.DescribeCenInstanceAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.CenId)
	d.Set("child_instance_id", object.ChildInstanceId)
	d.Set("child_instance_region_id", object.ChildInstanceRegionId)
	d.Set("child_instance_owner_id", strconv.Itoa(object.ChildInstanceOwnerId))

	return nil
}

func resourceAlicloudCenInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	instanceRegionId := d.Get("child_instance_region_id").(string)
	cenId, instanceId, err := cenService.GetCenIdAndAnotherId(d.Id())
	if err != nil {
		return WrapError(err)
	}
	instanceType, err := GetCenChildInstanceType(instanceId)
	if err != nil {
		return WrapError(err)
	}

	request := cbn.CreateDetachCenChildInstanceRequest()
	request.RegionId = client.RegionId
	request.CenId = cenId
	request.ChildInstanceId = instanceId
	request.ChildInstanceType = instanceType
	request.ChildInstanceRegionId = instanceRegionId
	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		raw, err = client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DetachCenChildInstance(request)
		})
		if err != nil {
			if IsExceptedError(err, InvalidCenInstanceStatus) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExceptedError(err, ParameterInstanceIdNotExist) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(cenService.WaitForCenInstanceAttachment(d.Id(), Deleted, DefaultCenTimeoutLong))
}
