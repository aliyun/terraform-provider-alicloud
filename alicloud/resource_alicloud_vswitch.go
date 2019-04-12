package alicloud

import (
	"fmt"
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
				return resource.RetryableError(fmt.Errorf("Creating Vswitch is timeout and got an error: %#v", err))
			}
			return resource.NonRetryableError(err)
		}
		resp, _ := raw.(*vpc.CreateVSwitchResponse)
		vswitchID = resp.VSwitchId
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "vswitch", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(vswitchID)

	if err := vpcService.WaitForVSwitch(vswitchID, Available, 300); err != nil {
		return WrapError(err)
	}

	return resourceAliyunSwitchUpdate(d, meta)
}

func resourceAliyunSwitchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	vswitch, err := vpcService.DescribeVswitch(d.Id())

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

	d.Partial(true)

	attributeUpdate := false
	request := vpc.CreateModifyVSwitchAttributeRequest()
	request.VSwitchId = d.Id()

	if d.HasChange("name") {
		d.SetPartial("name")
		request.VSwitchName = d.Get("name").(string)

		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		request.Description = d.Get("description").(string)

		attributeUpdate = true
	}
	if attributeUpdate {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVSwitchAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

	}

	d.Partial(false)

	return resourceAliyunSwitchRead(d, meta)
}

func resourceAliyunSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateDeleteVSwitchRequest()
	request.VSwitchId = d.Id()
	return resource.Retry(6*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVSwitch(request)
		})
		if err != nil {
			if IsExceptedError(err, VswitcInvalidRegionId) {
				return resource.NonRetryableError(err)
			}
			if IsExceptedError(err, InvalidVswitchIDNotFound) {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("Delete vswitch timeout and got an error: %#v.", err))
		}

		if _, err := vpcService.DescribeVswitch(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func buildAliyunSwitchArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateVSwitchRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	zoneID := Trim(d.Get("availability_zone").(string))

	zone, err := ecsService.DescribeZone(zoneID)
	if err != nil {
		return nil, err
	}

	err = ecsService.ResourceAvailable(zone, ResourceTypeVSwitch)
	if err != nil {
		return nil, err
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
