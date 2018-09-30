package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
	client := meta.(*AliyunClient)

	args := vpc.CreateAssociateRouteTableRequest()
	args.RouteTableId = Trim(d.Get("route_table_id").(string))
	args.VSwitchId = Trim(d.Get("vswitch_id").(string))
	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ar := args
		if _, err := client.vpcconn.AssociateRouteTable(ar); err != nil {
			if IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(fmt.Errorf("AssociateRouteTable got an error: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("AssociateRouteTable got an error: %#v", err))
		}
		return nil
	}); err != nil {
		return err
	}

	err := client.WaitForRouteTableAttachment(args.RouteTableId, args.VSwitchId, DefaultTimeout)
	//check the route table attachment
	if err != nil {
		return fmt.Errorf("Wait for route table attachment got error: %#v", err)
	}

	d.SetId(args.RouteTableId + COLON_SEPARATED + args.VSwitchId)

	return resourceAliyunRouteTableAttachmentRead(d, meta)
}

func resourceAliyunRouteTableAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	routeTableId, vSwitchId, err := getRouteTableIdAndVSwitchId(d, meta)
	routeTable, err := client.DescribeRouteTable(routeTableId)

	if len(routeTable.VSwitchIds.VSwitchId) <= 0 {
		d.SetId("")
		return nil
	}
	//Finding the vSwitchId
	err = client.DescribeRouteTableAttachment(routeTableId, vSwitchId)

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

	client := meta.(*AliyunClient)

	routeTableId, vSwitchId, err := getRouteTableIdAndVSwitchId(d, meta)
	if err != nil {
		return err
	}

	request := vpc.CreateUnassociateRouteTableRequest()
	request.RouteTableId = routeTableId
	request.VSwitchId = vSwitchId

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.vpcconn.UnassociateRouteTable(request)
		//Waiting for unassociate the route table
		if err != nil {
			if IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(fmt.Errorf("Unassociate Route Table timeout and got an error:%#v.", err))
			}
		}
		//Eusure the vswitch has been unassociated truly.
		err = client.DescribeRouteTableAttachment(routeTableId, vSwitchId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("Unassociate Route Table timeout."))
	})
}
