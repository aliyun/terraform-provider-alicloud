package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunRouteTableAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunRouteTableAttachmentCreate,
		Read:   resourceAliyunRouteTableAttachmentRead,
		Delete: resourceAliyunRouteTableAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"route_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunRouteTableAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	routeTableService := RouteTableService{client}

	args := vpc.CreateAssociateRouteTableRequest()
	args.RouteTableId = Trim(d.Get("route_table_id").(string))
	args.VSwitchId = Trim(d.Get("vswitch_id").(string))
	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ar := args
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.AssociateRouteTable(ar)
		})
		if err != nil {
			if IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(fmt.Errorf("AssociateRouteTable got an error: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("AssociateRouteTable got an error: %#v", err))
		}
		return nil
	}); err != nil {
		return err
	}

	err := routeTableService.WaitForRouteTableAttachment(args.RouteTableId, args.VSwitchId, DefaultTimeout)
	//check the route table attachment
	if err != nil {
		return fmt.Errorf("Wait for route table attachment got error: %#v", err)
	}

	d.SetId(args.RouteTableId + COLON_SEPARATED + args.VSwitchId)

	return resourceAliyunRouteTableAttachmentRead(d, meta)
}

func resourceAliyunRouteTableAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	routeTableService := RouteTableService{client}

	routeTableId, vSwitchId, err := routeTableService.GetRouteTableIdAndVSwitchId(d, meta)
	routeTable, err := routeTableService.DescribeRouteTable(routeTableId)

	if len(routeTable.VSwitchIds.VSwitchId) <= 0 {
		d.SetId("")
		return nil
	}
	//Finding the vSwitchId
	err = routeTableService.DescribeRouteTableAttachment(routeTableId, vSwitchId)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe Route Table Attribute: %#v", err)
	}

	d.Set("route_table_id", routeTableId)
	d.Set("vswitch_id", vSwitchId)
	return nil
}

func resourceAliyunRouteTableAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	routeTableService := RouteTableService{client}

	routeTableId, vSwitchId, err := routeTableService.GetRouteTableIdAndVSwitchId(d, meta)
	if err != nil {
		return err
	}

	request := vpc.CreateUnassociateRouteTableRequest()
	request.RouteTableId = routeTableId
	request.VSwitchId = vSwitchId

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateRouteTable(request)
		})
		//Waiting for unassociate the route table
		if err != nil {
			if IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(fmt.Errorf("Unassociate Route Table timeout and got an error:%#v.", err))
			}
		}
		//Eusure the vswitch has been unassociated truly.
		err = routeTableService.DescribeRouteTableAttachment(routeTableId, vSwitchId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("Unassociate Route Table timeout."))
	})
}
