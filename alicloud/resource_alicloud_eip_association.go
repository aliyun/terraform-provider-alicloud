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
				Required: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceAliyunEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateAssociateEipAddressRequest()
	request.RegionId = client.RegionId
	request.AllocationId = Trim(d.Get("allocation_id").(string))
	request.InstanceId = Trim(d.Get("instance_id").(string))
	request.InstanceType = EcsInstance

	if strings.HasPrefix(request.InstanceId, "lb-") {
		request.InstanceType = SlbInstance
	}
	if strings.HasPrefix(request.InstanceId, "ngw-") {
		request.InstanceType = Nat
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = instanceType.(string)
	}
	if privateIPAddress, ok := d.GetOk("private_ip_address"); ok {
		request.PrivateIpAddress = privateIPAddress.(string)
	}
	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.AssociateEipAddress(request)
		})
		if err != nil {
			if IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eip_association", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if err := vpcService.WaitForEip(request.AllocationId, InUse, 60); err != nil {
		return WrapError(err)
	}
	// There is at least 30 seconds delay for ecs instance
	if request.InstanceType == EcsInstance {
		time.Sleep(30 * time.Second)
	}

	d.SetId(request.AllocationId + ":" + request.InstanceId)

	return resourceAliyunEipAssociationRead(d, meta)
}

func resourceAliyunEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeEipAssociation(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_id", object.InstanceId)
	d.Set("allocation_id", object.AllocationId)
	d.Set("instance_type", object.InstanceType)
	return nil
}

func resourceAliyunEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	allocationId, instanceId := parts[0], parts[1]
	if err != nil {
		return WrapError(err)
	}

	request := vpc.CreateUnassociateEipAddressRequest()
	request.RegionId = client.RegionId
	request.AllocationId = allocationId
	request.InstanceId = instanceId
	request.InstanceType = EcsInstance

	if strings.HasPrefix(instanceId, "lb-") {
		request.InstanceType = SlbInstance
	}
	if strings.HasPrefix(instanceId, "ngw-") {
		request.InstanceType = Nat
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = instanceType.(string)
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateEipAddress(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InstanceIncorrectStatus, HaVipIncorrectStatus, TaskConflict,
				HasBeenUsedBySnatTable, HasBeenUsedByForwardEntry}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpcService.WaitForEipAssociation(d.Id(), Available, DefaultTimeoutMedium))
}
