package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
	"time"
)

func resourceAliyunVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpcCreate,
		Read:   resourceAliyunVpcRead,
		Update: resourceAliyunVpcUpdate,
		Delete: resourceAliyunVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cidr_block": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 2 || len(value) > 128 {
						errors = append(errors, fmt.Errorf("%s cannot be longer than 128 characters", k))
					}

					if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
						errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
					}

					return
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 2 || len(value) > 256 {
						errors = append(errors, fmt.Errorf("%s cannot be longer than 256 characters", k))

					}
					return
				},
			},
			"router_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"router_table_id": &schema.Schema{
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute router_table_id has been deprecated and replaced with route_table_id.",
			},
			"route_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunVpcCreate(d *schema.ResourceData, meta interface{}) error {

	args, err := buildAliyunVpcArgs(d, meta)
	if err != nil {
		return err
	}

	ecsconn := meta.(*AliyunClient).ecsconn

	var vpc *ecs.CreateVpcResponse
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		resp, err := ecsconn.CreateVpc(args)
		if err != nil {
			if IsExceptedError(err, VpcQuotaExceeded) {
				return resource.NonRetryableError(fmt.Errorf("The number of VPC has quota has reached the quota limit in your account, and please use existing VPCs or remove some of them."))
			}
			if IsExceptedError(err, UnknownError) {
				return resource.RetryableError(fmt.Errorf("Vpc is still creating result from some unknown error -- try again"))
			}
			return resource.NonRetryableError(err)
		}
		vpc = resp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Create vpc got an error :%#v", err)
	}

	d.SetId(vpc.VpcId)

	err = ecsconn.WaitForVpcAvailable(args.RegionId, vpc.VpcId, 60)
	if err != nil {
		return fmt.Errorf("Timeout when WaitForVpcAvailable")
	}

	return resourceAliyunVpcUpdate(d, meta)
}

func resourceAliyunVpcRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	vpc, err := client.DescribeVpc(d.Id())
	if err != nil {
		return err
	}

	if vpc == nil {
		d.SetId("")
		return nil
	}

	d.Set("cidr_block", vpc.CidrBlock)
	d.Set("name", vpc.VpcName)
	d.Set("description", vpc.Description)
	d.Set("router_id", vpc.VRouterId)
	vrouters, _, err := client.vpcconn.DescribeVRouters(&ecs.DescribeVRoutersArgs{
		VRouterId: vpc.VRouterId,
		RegionId:  getRegion(d, meta),
	})
	if len(vrouters) > 0 && len(vrouters[0].RouteTableIds.RouteTableId) > 0 {
		d.Set("router_table_id", vrouters[0].RouteTableIds.RouteTableId[0])
		d.Set("route_table_id", vrouters[0].RouteTableIds.RouteTableId[0])
	} else {
		d.Set("router_table_id", "")
		d.Set("route_table_id", "")
	}

	return nil
}

func resourceAliyunVpcUpdate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*AliyunClient).ecsconn

	d.Partial(true)

	attributeUpdate := false
	args := &ecs.ModifyVpcAttributeArgs{
		VpcId: d.Id(),
	}

	if d.HasChange("name") {
		d.SetPartial("name")
		args.VpcName = d.Get("name").(string)

		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		args.Description = d.Get("description").(string)

		attributeUpdate = true
	}

	if attributeUpdate {
		if err := conn.ModifyVpcAttribute(args); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAliyunVpcRead(d, meta)
}

func resourceAliyunVpcDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.DeleteVpc(d.Id())

		if err != nil {
			return resource.RetryableError(fmt.Errorf("Vpc in use - trying again while it is deleted."))
		}

		args := &ecs.DescribeVpcsArgs{
			RegionId: getRegion(d, meta),
			VpcId:    d.Id(),
		}
		vpc, _, descErr := conn.DescribeVpcs(args)
		if descErr != nil {
			return resource.NonRetryableError(err)
		} else if vpc == nil || len(vpc) < 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Vpc in use - trying again while it is deleted."))
	})
}

func buildAliyunVpcArgs(d *schema.ResourceData, meta interface{}) (*ecs.CreateVpcArgs, error) {
	args := &ecs.CreateVpcArgs{
		RegionId:  getRegion(d, meta),
		CidrBlock: d.Get("cidr_block").(string),
	}

	if v := d.Get("name").(string); v != "" {
		args.VpcName = v
	}

	if v := d.Get("description").(string); v != "" {
		args.Description = v
	}

	return args, nil
}
