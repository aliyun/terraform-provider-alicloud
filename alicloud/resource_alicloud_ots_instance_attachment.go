package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudOtsInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsInstanceAttachmentCreate,
		Read:   resourceAliyunOtsInstanceAttachmentRead,
		Delete: resourceAliyunOtsInstanceAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpc_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunOtsInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	req := ots.CreateBindInstance2VpcRequest()
	req.InstanceName = d.Get("instance_name").(string)
	req.InstanceVpcName = d.Get("vpc_name").(string)
	req.VirtualSwitchId = d.Get("vswitch_id").(string)

	if vsw, err := vpcService.DescribeVSwitch(d.Get("vswitch_id").(string)); err != nil {
		return err
	} else {
		req.VpcId = vsw.VpcId
	}

	_, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.BindInstance2Vpc(req)
	})
	if err != nil {
		return fmt.Errorf("Failed to bind instance with error: %s", err)
	}

	d.SetId(req.InstanceName)
	return resourceAliyunOtsInstanceAttachmentRead(d, meta)
}

func resourceAliyunOtsInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	inst, err := otsService.DescribeOtsInstanceVpc(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to describe instance vpc with error: %s", err)
	}
	// There is a bug that inst does not contain instance name and vswitch ID, so this resource does not support import function.
	//d.Set("instance_name", inst.InstanceName)
	d.Set("vpc_name", inst.InstanceVpcName)
	d.Set("vpc_id", inst.VpcId)
	return nil
}

func resourceAliyunOtsInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	inst, err := otsService.DescribeOtsInstanceVpc(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return fmt.Errorf("When unbinding vpc, failed to describe instance vpc with error: %s", err)
	}
	req := ots.CreateUnbindInstance2VpcRequest()
	req.InstanceName = d.Id()
	req.InstanceVpcName = inst.InstanceVpcName

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.UnbindInstance2Vpc(req)
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if _, err := otsService.DescribeOtsInstanceVpc(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("failed to describe instance with error: %s", err))
		}
		return resource.RetryableError(fmt.Errorf("delete instance timeout."))
	})
}
