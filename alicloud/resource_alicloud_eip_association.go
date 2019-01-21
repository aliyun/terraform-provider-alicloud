package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEipAssociationCreate,
		Read:   resourceAliyunEipAssociationRead,
		Delete: resourceAliyunEipAssociationDelete,

		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	args := vpc.CreateAssociateEipAddressRequest()
	args.AllocationId = Trim(d.Get("allocation_id").(string))
	args.InstanceId = Trim(d.Get("instance_id").(string))
	args.InstanceType = EcsInstance

	if strings.HasPrefix(args.InstanceId, "lb-") {
		args.InstanceType = SlbInstance
	}
	if strings.HasPrefix(args.InstanceId, "ngw-") {
		args.InstanceType = Nat
	}

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ar := args
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.AssociateEipAddress(ar)
		})
		if err != nil {
			if IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "new", args.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if err := vpcService.WaitForEip(args.AllocationId, InUse, 60); err != nil {
		return WrapError(err)
	}
	// There is at least 30 seconds delay for ecs instance
	if args.InstanceType == EcsInstance {
		time.Sleep(30 * time.Second)
	}

	d.SetId(args.AllocationId + ":" + args.InstanceId)

	return resourceAliyunEipAssociationRead(d, meta)
}

func resourceAliyunEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	eip, err := vpcService.DescribeEipAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", eip.InstanceId)
	d.Set("allocation_id", eip.AllocationId)
	return nil
}

func resourceAliyunEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	allocationId, instanceId, err := getAllocationIdAndInstanceId(d, meta)
	if err != nil {
		return WrapError(err)
	}

	request := vpc.CreateUnassociateEipAddressRequest()
	request.AllocationId = allocationId
	request.InstanceId = instanceId
	request.InstanceType = EcsInstance

	if strings.HasPrefix(instanceId, "lb-") {
		request.InstanceType = SlbInstance
	}
	if strings.HasPrefix(instanceId, "ngw-") {
		request.InstanceType = Nat
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateEipAddress(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InstanceIncorrectStatus, HaVipIncorrectStatus, TaskConflict}) {
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
		}

		if _, descErr := vpcService.DescribeEipAttachment(d.Id()); descErr != nil {
			if NotFoundError(descErr) {
				return nil
			}
			return resource.NonRetryableError(WrapError(descErr))
		}

		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))
	})
}

func getAllocationIdAndInstanceId(d *schema.ResourceData, meta interface{}) (string, string, error) {
	parts := strings.Split(d.Id(), ":")

	if len(parts) != 2 {
		return "", "", WrapError(Error("invalid resource id"))
	}
	return parts[0], parts[1], nil
}
