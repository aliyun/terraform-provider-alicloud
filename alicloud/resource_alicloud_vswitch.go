package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSwitchCIDRNetworkAddress,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliyunSwitchCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	var vswitchID string
	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args, err := buildAliyunSwitchArgs(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Building CreateVSwitchArgs got an error: %#v", err))
		}
		resp, err := client.vpcconn.CreateVSwitch(args)
		if err != nil {
			if IsExceptedError(err, TaskConflict) ||
				IsExceptedError(err, UnknownError) ||
				IsExceptedError(err, InvalidStatusRouteEntry) ||
				IsExceptedError(err, InvalidCidrBlockOverlapped) {
				return resource.RetryableError(fmt.Errorf("Creating Vswitch got an error: %#v", err))
			}
			return resource.NonRetryableError(err)
		}
		vswitchID = resp.VSwitchId
		return nil
	}); err != nil {
		return err
	}

	d.SetId(vswitchID)

	if err := client.WaitForVSwitch(vswitchID, Available, 300); err != nil {
		return fmt.Errorf("WaitForVSwitchAvailable got a error: %s", err)
	}

	return resourceAliyunSwitchUpdate(d, meta)
}

func resourceAliyunSwitchRead(d *schema.ResourceData, meta interface{}) error {

	vswitch, err := meta.(*AliyunClient).DescribeVswitch(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("availability_zone", vswitch.ZoneId)
	d.Set("vpc_id", vswitch.VpcId)
	d.Set("cidr_block", vswitch.CidrBlock)
	d.Set("name", vswitch.VSwitchName)
	d.Set("description", vswitch.Description)

	return nil
}

func resourceAliyunSwitchUpdate(d *schema.ResourceData, meta interface{}) error {

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
		if _, err := meta.(*AliyunClient).vpcconn.ModifyVSwitchAttribute(request); err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceAliyunSwitchRead(d, meta)
}

func resourceAliyunSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	request := vpc.CreateDeleteVSwitchRequest()
	request.VSwitchId = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.vpcconn.DeleteVSwitch(request); err != nil {
			if IsExceptedError(err, VswitcInvalidRegionId) {
				return resource.NonRetryableError(err)
			}
			if IsExceptedError(err, InvalidVswitchIDNotFound) {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("Delete vswitch timeout and got an error: %#v.", err))
		}

		if _, err := client.DescribeVswitch(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func buildAliyunSwitchArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateVSwitchRequest, error) {

	client := meta.(*AliyunClient)

	zoneID := Trim(d.Get("availability_zone").(string))

	zone, err := client.DescribeZone(zoneID)
	if err != nil {
		return nil, err
	}

	err = client.ResourceAvailable(zone, ResourceTypeVSwitch)
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

	return request, nil
}
