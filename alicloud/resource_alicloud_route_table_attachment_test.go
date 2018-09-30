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
)

func init() {
	resource.AddTestSweepers("alicloud_route_table_attachment", &resource.Sweeper{
		Name: "alicloud_route_table_attachment",
		F:    testSweepRouteTableAttachment,
	})
}

func testSweepRouteTableAttachment(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	conn := client.(*AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var routeTables []vpc.RouterTableListType
	req := vpc.CreateDescribeRouteTableListRequest()
	req.RegionId = conn.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		resp, err := conn.vpcconn.DescribeRouteTableList(req)
		if err != nil {
			return fmt.Errorf("Error retrieving RouteTables: %s", err)
		}
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
		for _, vswitch := range vtb.VSwitchIds.VSwitchId {
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
			log.Printf("[INFO] Unassociating Route Table: %s (%s)", name, id)
			req := vpc.CreateUnassociateRouteTableRequest()
			req.RouteTableId = id
			req.VSwitchId = vswitch
			if _, err := conn.vpcconn.UnassociateRouteTable(req); err != nil {
				log.Printf("[ERROR] Failed to unassociate Route Table (%s (%s)): %s", name, id, err)
			}
		}
	}
	return nil
}

func TestAccAlicloudRouteTableAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "alicloud_route_table_attachment.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteTableAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteTableAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableAttachmentExists("alicloud_route_table_attachment.foo"),
					resource.TestCheckResourceAttrSet(
						"alicloud_route_table_attachment.foo", "vswitch_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_route_table_attachment.foo", "route_table_id"),
				),
			},
		},
	})
}

func testAccCheckRouteTableAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route Table ID is set")
		}
		client := testAccProvider.Meta().(*AliyunClient)
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if len(parts) != 2 {
			return fmt.Errorf("invalid resource id")
		}
		err := client.DescribeRouteTableAttachment(parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("Describe Route Table attachment error %#v", err)
		}
		return nil
	}
}

func testAccCheckRouteTableAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_route_table_attachment" {
			continue
		}
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if len(parts) != 2 {
			return fmt.Errorf("invalid resource id")
		}
		err := client.DescribeRouteTableAttachment(parts[0], parts[1])
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describe Route Table attachment error %#v", err)
		}
	}
	return nil
}

const testAccRouteTableAttachmentConfig = `

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "tf-testAcc_route_table_attachment"
}
 data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}
 resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "tf-testAcc_route_table_attachment"
}

resource "alicloud_route_table" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
    name = "tf-testAcc_route_table_attachment"
    description = "tf-testAcc_route_table_attachment"
}

resource "alicloud_route_table_attachment" "foo" {
	vswitch_id = "${alicloud_vswitch.foo.id}"
	route_table_id = "${alicloud_route_table.foo.id}"
}
`
