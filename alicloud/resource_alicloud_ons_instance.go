package alicloud

import (
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/hashicorp/terraform/helper/schema"
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
				ValidateFunc: validateOnsInstanceName,
			},

			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOnsInstanceRemark,
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
	request.InstanceName = d.Get("name").(string)
	request.PreventCache = onsService.GetPreventCache()
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = v.(string)
	}

	raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsInstanceCreate(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ons_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw)

	var response *ons.OnsInstanceCreateResponse
	response, _ = raw.(*ons.OnsInstanceCreateResponse)
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
	request.InstanceId = d.Id()
	request.PreventCache = onsService.GetPreventCache()

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
		addDebug(request.GetActionName(), raw)
	}

	return resourceAlicloudOnsInstanceRead(d, meta)
}

func resourceAlicloudOnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	request := ons.CreateOnsInstanceDeleteRequest()
	request.InstanceId = d.Id()
	request.PreventCache = onsService.GetPreventCache()

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsInstanceDelete(request)
		})
		if err != nil {
			if IsExceptedError(err, OnsInstanceNotEmpty) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
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
