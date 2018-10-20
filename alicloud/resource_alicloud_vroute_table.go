package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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

			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	routeTableService := RouteTableService{client}

	request := vpc.CreateCreateRouteTableRequest()
	request.RegionId = client.RegionId

	request.VpcId = d.Get("vpc_id").(string)
	request.RouteTableName = d.Get("name").(string)
	request.Description = d.Get("description").(string)
	request.ClientToken = buildClientToken("TF-AllocateRouteTable")

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateRouteTable(request)
	})
	if err != nil {
		return err
	}
	routeTable, _ := raw.(*vpc.CreateRouteTableResponse)
	d.SetId(routeTable.RouteTableId)

	if err := routeTableService.WaitForRouteTable(routeTable.RouteTableId, DefaultTimeout); err != nil {
		return fmt.Errorf("Wait for route table got error: %#v", err)
	}

	return resourceAliyunRouteTableRead(d, meta)
}

func resourceAliyunRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	routeTableService := RouteTableService{client}

	resp, err := routeTableService.DescribeRouteTable(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe Route Table Attribute: %#v", err)
	}

	d.Set("vpc_id", resp.VpcId)
	d.Set("name", resp.RouteTableName)
	d.Set("description", resp.Description)

	return nil
}

func resourceAliyunRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	update := false
	request := vpc.CreateModifyRouteTableAttributesRequest()
	request.RouteTableId = d.Id()

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}

	if d.HasChange("name") {
		request.RouteTableName = d.Get("name").(string)
		update = true
	}

	if update {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyRouteTableAttributes(request)
		})
		if err != nil {
			return err
		}
	}

	return resourceAliyunRouteTableRead(d, meta)
}

func resourceAliyunRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	routeTableService := RouteTableService{client}

	request := vpc.CreateDeleteRouteTableRequest()
	request.RouteTableId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteRouteTable(request)
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if _, err := routeTableService.DescribeRouteTable(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Error describing route table failed when deleting route table: %#v", err))
		}
		return resource.RetryableError(fmt.Errorf("Delete Route Table timeout."))
	})
}
