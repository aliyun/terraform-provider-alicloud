package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSwitchCreate,
		Read:   resourceAliyunSwitchRead,
		Update: resourceAliyunSwitchUpdate,
		Delete: resourceAliyunSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSwitchCIDRNetworkAddress,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliyunSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	var vswitchID string
	request, err := buildAliyunSwitchArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateVSwitch(&args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{TaskConflict, UnknownError, InvalidStatusRouteEntry,
				InvalidCidrBlockOverlapped, Throttling, TokenProcessing}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*vpc.CreateVSwitchResponse)
		vswitchID = response.VSwitchId
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vswitch", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	d.SetId(vswitchID)
	if err := vpcService.WaitForVSwitch(vswitchID, Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return resourceAliyunSwitchRead(d, meta)
}

func resourceAliyunSwitchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	vswitch, err := vpcService.DescribeVSwitch(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("availability_zone", vswitch.ZoneId)
	d.Set("vpc_id", vswitch.VpcId)
	d.Set("cidr_block", vswitch.CidrBlock)
	d.Set("name", vswitch.VSwitchName)
	d.Set("description", vswitch.Description)

	return nil
}

func resourceAliyunSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	update := false
	request := vpc.CreateModifyVSwitchAttributeRequest()
	request.VSwitchId = d.Id()

	if d.HasChange("name") {
		request.VSwitchName = d.Get("name").(string)
		update = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}
	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVSwitchAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}
	return resourceAliyunSwitchRead(d, meta)
}

func resourceAliyunSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := vpc.CreateDeleteVSwitchRequest()
	request.VSwitchId = d.Id()
	err := resource.Retry(6*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVSwitch(request)
		})
		if err != nil {
			if IsExceptedError(err, VswitcInvalidRegionId) {
				return resource.NonRetryableError(err)
			}
			if IsExceptedError(err, InvalidVswitchIDNotFound) {
				return nil
			}

			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpcService.WaitForVSwitch(d.Id(), Deleted, DefaultTimeout))
}

func buildAliyunSwitchArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateVSwitchRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	zoneID := Trim(d.Get("availability_zone").(string))
	zone, err := ecsService.DescribeZone(zoneID)
	if err != nil {
		return nil, WrapError(err)
	}
	err = ecsService.ResourceAvailable(zone, ResourceTypeVSwitch)
	if err != nil {
		return nil, WrapError(err)
	}
	request := vpc.CreateCreateVSwitchRequest()
	request.VpcId = Trim(d.Get("vpc_id").(string))
	request.ZoneId = zoneID
	request.CidrBlock = Trim(d.Get("cidr_block").(string))

	if v, ok := d.GetOk("name"); ok && v != "" {
		request.VSwitchName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v != "" {
		request.Description = v.(string)
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	return request, nil
}
