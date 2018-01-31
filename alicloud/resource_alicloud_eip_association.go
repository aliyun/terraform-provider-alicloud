package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEipAssociationCreate,
		Read:   resourceAliyunEipAssociationRead,
		Delete: resourceAliyunEipAssociationDelete,

		Schema: map[string]*schema.Schema{
			"allocation_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*AliyunClient).ecsconn

	args := ecs.AssociateEipAddressArgs{
		AllocationId: d.Get("allocation_id").(string),
		InstanceId:   d.Get("instance_id").(string),
		InstanceType: ecs.EcsInstance,
	}
	if strings.HasPrefix(args.InstanceId, "lb-") {
		args.InstanceType = ecs.SlbInstance
	}
	if strings.HasPrefix(args.InstanceId, "ngw-") {
		args.InstanceType = ecs.Nat
	}

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ar := args
		if err := conn.NewAssociateEipAddress(&ar); err != nil {
			if IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(fmt.Errorf("AssociateEip got an error: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("AssociateEip got an error: %#v", err))
		}
		return nil
	}); err != nil {
		return err
	}

	if err := conn.WaitForEip(getRegion(d, meta), args.AllocationId, ecs.EipStatusInUse, 60); err != nil {
		return fmt.Errorf("Error Waitting for EIP allocated: %#v", err)
	}

	d.SetId(args.AllocationId + ":" + args.InstanceId)

	return resourceAliyunEipAssociationRead(d, meta)
}

func resourceAliyunEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	allocationId, instanceId, err := getAllocationIdAndInstanceId(d, meta)
	if err != nil {
		return err
	}

	eip, err := client.DescribeEipAddress(allocationId)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe Eip Attribute: %#v", err)
	}

	if eip.InstanceId != instanceId {
		d.SetId("")
		return nil
	}

	d.Set("instance_id", eip.InstanceId)
	d.Set("allocation_id", allocationId)
	return nil
}

func resourceAliyunEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*AliyunClient).ecsconn

	allocationId, instanceId, err := getAllocationIdAndInstanceId(d, meta)
	if err != nil {
		return err
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.UnassociateEipAddress(allocationId, instanceId)

		if err != nil {
			if IsExceptedError(err, InstanceIncorrectStatus) ||
				IsExceptedError(err, HaVipIncorrectStatus) ||
				IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(fmt.Errorf("Unassociat EIP timeout and got an error:%#v.", err))
			}
		}

		args := &ecs.DescribeEipAddressesArgs{
			RegionId:     getRegion(d, meta),
			AllocationId: allocationId,
		}

		eips, _, descErr := conn.DescribeEipAddresses(args)

		if descErr != nil {
			return resource.NonRetryableError(descErr)
		} else if eips == nil || len(eips) < 1 {
			return nil
		}
		for _, eip := range eips {
			if eip.Status != ecs.EipStatusAvailable {
				return resource.RetryableError(fmt.Errorf("Unassociat EIP timeout and got an error:%#v.", err))
			}
		}

		return nil
	})
}

func getAllocationIdAndInstanceId(d *schema.ResourceData, meta interface{}) (string, string, error) {
	parts := strings.Split(d.Id(), ":")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}
