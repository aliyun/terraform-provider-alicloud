package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

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

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
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

func testAccCheckRouteTableDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	routeTableService := VpcService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_route_table" {
			continue
		}
		instance, err := routeTableService.DescribeRouteTable(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if instance.RouteTableId != "" {
			return WrapError(Error("Route Table %s still exist", instance.RouteTableId))
		}
	}
	return nil
}

func TestAccAlicloudRouteTableBasic(t *testing.T) {
	var v vpc.RouterTableListType
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "alicloud_route_table.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"vpc_id":      CHECKSET,
		"name":        fmt.Sprintf("tf-testAccRouteTable%d", rand),
		"description": "",
	})
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RouteTableNoSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRouteTableConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccRouteTable%d_change", rand),
					}),
				),
			},
			{
				Config: testAccRouteTableConfig_description(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccRouteTable%d_description", rand),
					}),
				),
			},
			{
				Config: testAccRouteTableConfig_tags(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccRouteTableConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         fmt.Sprintf("tf-testAccRouteTable%d_all", rand),
						"description":  fmt.Sprintf("tf-testAccRouteTable%d_description_all", rand),
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudRouteTableMulti(t *testing.T) {
	var v vpc.RouterTableListType
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "alicloud_route_table.default.4"
	ra := resourceAttrInit(resourceId, map[string]string{
		"vpc_id":      CHECKSET,
		"name":        fmt.Sprintf("tf-testAccRouteTable%d", rand),
		"description": "",
	})
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RouteTableNoSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccRouteTableConfigBasic(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccRouteTable%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}	

resource "alicloud_route_table" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  name = "${var.name}"
}
`, rand)
}

func testAccRouteTableConfig_name(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccRouteTable%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}	

resource "alicloud_route_table" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  name = "${var.name}_change"
}
`, rand)
}

func testAccRouteTableConfig_description(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccRouteTable%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}	

resource "alicloud_route_table" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  name = "${var.name}_change"
  description = "${var.name}_description"
}
`, rand)
}

func testAccRouteTableConfig_tags(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccRouteTable%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}	

resource "alicloud_route_table" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  name = "${var.name}_change"
  description = "${var.name}_description"
  tags 		= {
		Created = "TF"
		For 	= "acceptance test"
  }
}
`, rand)
}

func testAccRouteTableConfig_all(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccRouteTable%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}	

resource "alicloud_route_table" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  name = "${var.name}_all"
  description = "${var.name}_description_all"
}
`, rand)
}

func testAccRouteTableConfigMulti(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccRouteTable%d"
}
resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}	

resource "alicloud_route_table" "default" {
  count = 5
  vpc_id = "${alicloud_vpc.default.id}"
  name = "${var.name}"
}
`, rand)
}
