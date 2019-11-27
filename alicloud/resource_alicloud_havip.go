package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunHaVip() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunHaVipCreate,
		Read:   resourceAliyunHaVipRead,
		Update: resourceAliyunHaVipUpdate,
		Delete: resourceAliyunHaVipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"ip_address": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunHaVipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	haVipService := HaVipService{client}

	request := vpc.CreateCreateHaVipRequest()
	request.RegionId = client.RegionId

	request.VSwitchId = d.Get("vswitch_id").(string)
	request.IpAddress = d.Get("ip_address").(string)
	request.Description = d.Get("description").(string)
	request.ClientToken = buildClientToken(request.GetActionName())

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateHaVip(request)
	})
	if err != nil {
		return err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	havip, _ := raw.(*vpc.CreateHaVipResponse)
	d.SetId(havip.HaVipId)
	if err := haVipService.WaitForHaVip(havip.HaVipId, Available, 2*DefaultTimeout); err != nil {
		return fmt.Errorf("WaitHaVip %s got error: %#v, %s", Available, err, havip.HaVipId)
	}
	return resourceAliyunHaVipRead(d, meta)
}

func resourceAliyunHaVipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	haVipService := HaVipService{client}

	resp, err := haVipService.DescribeHaVip(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe HaVip Attribute: %#v", err)
	}

	d.Set("vswitch_id", resp.VSwitchId)
	d.Set("ip_address", resp.IpAddress)
	d.Set("description", resp.Description)

	return nil
}

func resourceAliyunHaVipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("description") {
		request := vpc.CreateModifyHaVipAttributeRequest()
		request.RegionId = client.RegionId
		request.HaVipId = d.Id()
		request.Description = d.Get("description").(string)
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyHaVipAttribute(request)
		})
		if err != nil {
			return err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAliyunHaVipRead(d, meta)
}

func resourceAliyunHaVipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	haVipService := HaVipService{client}

	if err := haVipService.WaitForHaVip(d.Id(), Available, 2*DefaultTimeout); err != nil {
		return fmt.Errorf("WaitHaVip %s got error: %#v, %s", Available, err, d.Id())
	}
	request := vpc.CreateDeleteHaVipRequest()
	request.RegionId = client.RegionId
	request.HaVipId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteHaVip(request)
		})
		if err != nil {
			if IsExceptedError(err, InvalidHaVipIdNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if _, err := haVipService.DescribeHaVip(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Error describing havip failed when deleting HaVip: %#v", err))
		}
		return nil
	})
}
