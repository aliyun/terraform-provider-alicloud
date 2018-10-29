package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunHaVipAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunHaVipAttachmentCreate,
		Read:   resourceAliyunHaVipAttachmentRead,
		Delete: resourceAliyunHaVipAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"havip_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunHaVipAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	haVipService := HaVipService{client}

	args := vpc.CreateAssociateHaVipRequest()
	args.HaVipId = Trim(d.Get("havip_id").(string))
	args.InstanceId = Trim(d.Get("instance_id").(string))
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		ar := args
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.AssociateHaVip(ar)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{TaskConflict, IncorrectHaVipStatus, InvalidVipStatus}) {
				return resource.RetryableError(fmt.Errorf("AssociateHaVip got an error: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("AssociateHaVip got an error: %#v", err))
		}
		return nil
	}); err != nil {
		return err
	}
	//check the havip attachment
	if err := haVipService.WaitForHaVipAttachment(args.HaVipId, args.InstanceId, 5*DefaultTimeout); err != nil {
		return fmt.Errorf("Wait for havip attachment got error: %#v", err)
	}

	d.SetId(args.HaVipId + COLON_SEPARATED + args.InstanceId)

	return resourceAliyunHaVipAttachmentRead(d, meta)
}

func resourceAliyunHaVipAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	haVipService := HaVipService{client}

	haVipId, instanceId, err := getHaVipIdAndInstanceId(d, meta)
	err = haVipService.DescribeHaVipAttachment(haVipId, instanceId)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe HaVip Attribute: %#v", err)
	}

	d.Set("havip_id", haVipId)
	d.Set("instance_id", instanceId)
	return nil
}

func resourceAliyunHaVipAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	haVipService := HaVipService{client}

	haVipId, instanceId, err := getHaVipIdAndInstanceId(d, meta)
	if err != nil {
		return err
	}

	request := vpc.CreateUnassociateHaVipRequest()
	request.HaVipId = haVipId
	request.InstanceId = instanceId

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateHaVip(request)
		})
		//Waiting for unassociate the havip
		if err != nil {
			if IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(fmt.Errorf("Unassociate HaVip timeout and got an error:%#v.", err))
			}
		}
		//Eusure the instance has been unassociated truly.
		err = haVipService.DescribeHaVipAttachment(haVipId, instanceId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("Unassociate HaVip timeout."))
	})
}
