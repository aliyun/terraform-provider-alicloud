package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunRouteTableCreate,
		Read:   resourceAliyunRouteTableRead,
		Update: resourceAliyunRouteTableUpdate,
		Delete: resourceAliyunRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceDescription,
			},
			"route_table_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	request := vpc.CreateCreateRouteTableRequest()
	request.RegionId = getRegionId(d, meta)

	request.VpcId = d.Get("vpc_id").(string)
	request.RouteTableName = d.Get("route_table_name").(string)
	request.Description = d.Get("description").(string)
	request.ClientToken = buildClientToken("TF-AllocateRouteTable")

	routeTable, err := client.vpcconn.CreateRouteTable(request)
	if err != nil {
		return err
	}
	d.SetId(routeTable.RouteTableId)
	time.Sleep(3 * time.Second)
	return resourceAliyunRouteTableUpdate(d, meta)
}

func resourceAliyunRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	resp, err := client.DescribeRouteTable(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe Route Table Attribute: %#v", err)
	}

	d.Set("vpc_id", resp.VpcId)
	d.Set("route_table_name", resp.RouteTableName)
	d.Set("description", resp.Description)

	return nil
}

func resourceAliyunRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	update := false
	request := vpc.CreateModifyRouteTableAttributesRequest()
	request.RouteTableId = d.Id()

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}

	if d.HasChange("route_table_name") {
		request.RouteTableName = d.Get("route_table_name").(string)
		update = true
	}

	if update {
		if _, err := meta.(*AliyunClient).vpcconn.ModifyRouteTableAttributes(request); err != nil {
			return err
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliyunRouteTableRead(d, meta)
	}

	d.Partial(false)

	return resourceAliyunRouteTableRead(d, meta)
}

func resourceAliyunRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	request := vpc.CreateDeleteRouteTableRequest()
	request.RouteTableId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.vpcconn.DeleteRouteTable(request); err != nil {
			return resource.NonRetryableError(err)
		}
		if _, err := client.DescribeRouteTable(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Error describing route table failed when deleting route table: %#v", err))
		}
		return nil
	})
}
