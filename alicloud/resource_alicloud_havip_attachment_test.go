package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// At present, only white list users can operate HaVip Resource. So close havip sweeper.
func init() {
	resource.AddTestSweepers("alicloud_havip_attachment", &resource.Sweeper{
		Name: "alicloud_havip_attachment",
		F:    testSweepHaVipAttachment,
	})
}

func testSweepHaVipAttachment(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	conn := client.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var haVips []vpc.HaVip
	req := vpc.CreateDescribeHaVipsRequest()
	req.RegionId = conn.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := conn.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeHaVips(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving HaVips: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeHaVipsResponse)
		if resp == nil || len(resp.HaVips.HaVip) < 1 {
			break
		}
		haVips = append(haVips, resp.HaVips.HaVip...)

		if len(resp.HaVips.HaVip) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, haVip := range haVips {
		id := haVip.HaVipId
		desc := haVip.Description
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(desc), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping HaVip: (%s)", id)
			continue
		}
		for _, instance := range haVip.AssociatedInstances.AssociatedInstance {
			log.Printf("[INFO] Unassociating HaVip: (%s)", id)
			req := vpc.CreateUnassociateHaVipRequest()
			req.HaVipId = id
			req.InstanceId = instance
			_, err := conn.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.UnassociateHaVip(req)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to unassociate HaVip %s): %s", id, err)
			}
		}
	}
	return nil
}

// At present, only white list users can operate HaVip Resource.
func SkipTestAccAlicloudVPCHavipAttachmentBasic(t *testing.T) {
	resourceId := "alicloud_havip_attachment.foo"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckHaVipAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHaVipAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHaVipAttachmentExists("alicloud_havip_attachment.foo"),
					resource.TestCheckResourceAttrSet(
						"alicloud_havip_attachment.foo", "havip_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_havip_attachment.foo", "instance_id"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckHaVipAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No HaVip ID is set")
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		haVipService := HaVipService{client}
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if len(parts) != 2 {
			return fmt.Errorf("invalid resource id")
		}
		err := haVipService.DescribeHaVipAttachment(parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("Describe HaVip attachment error %#v", err)
		}
		return nil
	}
}

func testAccCheckHaVipAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	haVipService := HaVipService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_havip_attachment" {
			continue
		}
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if len(parts) != 2 {
			return fmt.Errorf("invalid resource id")
		}
		err := haVipService.DescribeHaVipAttachment(parts[0], parts[1])
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describe HaVip attachment error %#v", err)
		}
	}
	return nil
}

const testAccHaVipAttachmentConfig = `

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_images" "default" {
	name_regex = "^ubuntu"
	most_recent = true
	owners = "system"
}

variable "name" {
	default = "tf-testAccHaVipAttachment"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_havip" "foo" {
	vswitch_id = "${alicloud_vswitch.foo.id}"
	description = "${var.name}"
}

resource "alicloud_havip_attachment" "foo" {
	havip_id = "${alicloud_havip.foo.id}"
	instance_id = "${alicloud_instance.foo.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "foo" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_id = "${alicloud_vswitch.foo.id}"
	image_id = "${data.alicloud_images.default.images.0.id}"
	# series III
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	instance_name = "${var.name}"
	user_data = "echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"
}

`
