package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_route_table", &resource.Sweeper{
		Name: "alicloud_route_table",
		F:    testSweepRouteTable,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_route_table_attachment",
		},
	})
}

func testSweepRouteTable(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.RouteTableNoSupportedRegions) {
		log.Printf("[INFO] Skipping Route Table unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var routeTables []vpc.RouterTableListType
	req := vpc.CreateDescribeRouteTableListRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouteTableList(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving RouteTables: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeRouteTableListResponse)
		if resp == nil || len(resp.RouterTableList.RouterTableListType) < 1 {
			break
		}
		routeTables = append(routeTables, resp.RouterTableList.RouterTableListType...)

		if len(resp.RouterTableList.RouterTableListType) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, vtb := range routeTables {
		name := vtb.RouteTableName
		id := vtb.RouteTableId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Route Table: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Route Table: %s (%s)", name, id)
		req := vpc.CreateDeleteRouteTableRequest()
		req.RouteTableId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteRouteTable(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Route Table (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudRouteTable_basic(t *testing.T) {
	var routeTable vpc.DescribeRouteTableListResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RouteTableNoSupportedRegions)
		},
		// module name
		IDRefreshName: "alicloud_route_table.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableListExists("alicloud_route_table.foo", &routeTable),
					resource.TestCheckResourceAttrSet(
						"alicloud_route_table.foo", "vpc_id"),
					resource.TestCheckResourceAttr(
						"alicloud_route_table.foo", "name", "tf-testAcc_route_table"),
					resource.TestCheckResourceAttr(
						"alicloud_route_table.foo", "description", "tf-testAcc_route_table"),
				),
			},
		},
	})
}

func testAccCheckRouteTableListExists(n string, routeTable *vpc.DescribeRouteTableListResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route Table ID is set")
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		routeTableService := RouteTableService{client}
		instance, err := routeTableService.DescribeRouteTable(rs.Primary.ID)
		if err != nil {
			return err
		}
		if routeTable == nil || len((*routeTable).RouterTableList.RouterTableListType) <= 0 {
			return err
		}
		(*routeTable).RouterTableList.RouterTableListType[0].RouteTableId = instance.RouteTableId
		return nil
	}
}

func testAccCheckRouteTableDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	routeTableService := RouteTableService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_route_table" {
			continue
		}
		instance, err := routeTableService.DescribeRouteTable(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describe Route Table error %#v", err)
		}
		if instance.RouteTableId != "" {
			return fmt.Errorf("Route Table %s still exist", instance.RouteTableId)
		}
	}
	return nil
}

const testAccRouteTableConfig = `

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "tf-testAccVpcConfig"
}	

resource "alicloud_route_table" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  name = "tf-testAcc_route_table"
  description = "tf-testAcc_route_table"
}

`
