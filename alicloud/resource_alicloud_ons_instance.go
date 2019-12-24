package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudOnsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOnsInstanceCreate,
		Read:   resourceAlicloudOnsInstanceRead,
		Update: resourceAlicloudOnsInstanceUpdate,
		Delete: resourceAlicloudOnsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},

			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},

			// Computed Values
			"instance_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"instance_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"release_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudOnsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	request := ons.CreateOnsInstanceCreateRequest()
	request.RegionId = client.RegionId
	request.InstanceName = d.Get("name").(string)
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = v.(string)
	}

	var response *ons.OnsInstanceCreateResponse
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsInstanceCreate(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{OnsThrottlingUser}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ons.OnsInstanceCreateResponse)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ons_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(response.Data.InstanceId)

	return resourceAlicloudOnsInstanceRead(d, meta)
}

func resourceAlicloudOnsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	response, err := onsService.DescribeOnsInstance(d.Id())

	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", response.InstanceBaseInfo.InstanceName)
	d.Set("instance_type", response.InstanceBaseInfo.InstanceType)
	d.Set("instance_status", response.InstanceBaseInfo.InstanceStatus)
	d.Set("release_time", time.Unix(int64(response.InstanceBaseInfo.ReleaseTime)/1000, 0).Format("2006-01-02 03:04:05"))

	return nil
}

func resourceAlicloudOnsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	attributeUpdate := false

	request := ons.CreateOnsInstanceUpdateRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()

	if d.HasChange("name") {
		var name string
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
		request.InstanceName = name
		attributeUpdate = true
	}

	if d.HasChange("remark") {
		var remark string
		if v, ok := d.GetOk("remark"); ok {
			remark = v.(string)
		}
		request.Remark = remark
		attributeUpdate = true
	}

	if attributeUpdate {
		raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsInstanceUpdate(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlicloudOnsInstanceRead(d, meta)
}

func resourceAlicloudOnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	request := ons.CreateOnsInstanceDeleteRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsInstanceDelete(request)
		})
		if err != nil {
			if IsExceptedError(err, OnsInstanceNotEmpty) {
				return resource.RetryableError(err)
			}
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
		if IsExceptedError(err, OnsInstanceNotExist) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(onsService.WaitForOnsInstance(d.Id(), Deleted, DefaultTimeoutMedium))
}
