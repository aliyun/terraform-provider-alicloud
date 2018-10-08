package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceDescription,
			},
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},

			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunHaVipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	request := vpc.CreateCreateHaVipRequest()
	request.RegionId = getRegionId(d, meta)

	request.VSwitchId = d.Get("vswitch_id").(string)
	request.IpAddress = d.Get("ip_address").(string)
	request.Description = d.Get("description").(string)
	request.ClientToken = buildClientToken("TF-AllocateHaVip")

	havip, err := client.vpcconn.CreateHaVip(request)
	if err != nil {
		return err
	}
	d.SetId(havip.HaVipId)
	if err := client.WaitForHaVip(havip.HaVipId, Available, 2*DefaultTimeout); err != nil {
		return fmt.Errorf("WaitHaVip %s got error: %#v, %s", Available, err, havip.HaVipId)
	}
	return resourceAliyunHaVipRead(d, meta)
}

func resourceAliyunHaVipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	resp, err := client.DescribeHaVip(d.Id())
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
	client := meta.(*AliyunClient)
	if d.HasChange("description") {
		request := vpc.CreateModifyHaVipAttributeRequest()
		request.HaVipId = d.Id()
		request.Description = d.Get("description").(string)
		if _, err := client.vpcconn.ModifyHaVipAttribute(request); err != nil {
			return err
		}
	}

	return resourceAliyunHaVipRead(d, meta)
}

func resourceAliyunHaVipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	request := vpc.CreateDeleteHaVipRequest()
	request.HaVipId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.vpcconn.DeleteHaVip(request); err != nil {
			if IsExceptedError(err, InvalidHaVipIdNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		if _, err := client.DescribeHaVip(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Error describing havip failed when deleting HaVip: %#v", err))
		}
		return nil
	})
}
