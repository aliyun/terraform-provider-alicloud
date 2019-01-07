package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_ess_scalinggroup", &resource.Sweeper{
		Name: "alicloud_ess_scalinggroup",
		F:    testSweepEssGroups,
	})
}

func testSweepEssGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var groups []ess.ScalingGroup
	req := ess.CreateDescribeScalingGroupsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingGroups(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Scaling groups: %s", err)
		}
		resp, _ := raw.(*ess.DescribeScalingGroupsResponse)
		if resp == nil || len(resp.ScalingGroups.ScalingGroup) < 1 {
			break
		}
		groups = append(groups, resp.ScalingGroups.ScalingGroup...)

		if len(resp.ScalingGroups.ScalingGroup) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range groups {
		name := v.ScalingGroupName
		id := v.ScalingGroupId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Scaling Group: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting Scaling Group: %s (%s)", name, id)
		req := ess.CreateDeleteScalingGroupRequest()
		req.ScalingGroupId = id
		req.ForceDelete = requests.NewBoolean(true)
		_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DeleteScalingGroup(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Scaling Group (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		time.Sleep(2 * time.Minute)
	}
	return nil
}

func TestAccAlicloudEssScalingGroup_basic(t *testing.T) {
	var sg ess.ScalingGroup
	rand1 := acctest.RandIntRange(10000, 999999)
	rand2 := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_group.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingGroup(EcsInstanceCommonTestCase, rand1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingGroupExists(
						"alicloud_ess_scaling_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "min_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "max_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "scaling_group_name", fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand1)),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "removal_policies.#", "2"),
				),
			},

			{
				Config: testAccEssScalingGroup_update(EcsInstanceCommonTestCase, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingGroupExists(
						"alicloud_ess_scaling_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "min_size", "2"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "max_size", "2"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "scaling_group_name", fmt.Sprintf("tf-testAccEssScalingGroup_update-%d", rand2)),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "removal_policies.#", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudEssScalingGroup_vpc(t *testing.T) {
	var sg ess.ScalingGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_group.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingGroup_vpc(EcsInstanceCommonTestCase, acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingGroupExists(
						"alicloud_ess_scaling_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "min_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "max_size", "2"),
					resource.TestMatchResourceAttr(
						"alicloud_ess_scaling_group.foo", "scaling_group_name", regexp.MustCompile("^tf-testAccEssScalingGroup_vpc-*")),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "removal_policies.#", "2"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "vswitch_ids.#", "2"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.foo", "multi_az_policy", "BALANCE"),
				),
			},
		},
	})

}

func TestAccAlicloudEssScalingGroup_slb(t *testing.T) {
	var sg ess.ScalingGroup
	var slb slb.DescribeLoadBalancerAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_group.scaling",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingGroup_slb(EcsInstanceCommonTestCase, acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingGroupExists(
						"alicloud_ess_scaling_group.scaling", &sg),
					testAccCheckSlbExists(
						"alicloud_slb.instance.0", &slb),
					testAccCheckSlbExists(
						"alicloud_slb.instance.1", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.scaling", "min_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.scaling", "max_size", "1"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.scaling", "loadbalancer_ids.#", "2"),
				),
			},
		},
	})

}

func TestAccAlicloudEssScalingGroup_slbempty(t *testing.T) {
	var sg ess.ScalingGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_group.scaling",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingGroup_slbempty(EcsInstanceCommonTestCase, acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingGroupExists(
						"alicloud_ess_scaling_group.scaling", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_ess_scaling_group.scaling", "loadbalancer_ids.#", "0"),
				),
			},
		},
	})

}

func testAccCheckEssScalingGroupExists(n string, d *ess.ScalingGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ESS Scaling Group ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		essService := EssService{client}
		attr, err := essService.DescribeScalingGroupById(rs.Primary.ID)
		log.Printf("[DEBUG] check scaling group %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = attr
		return nil
	}
}

func testAccCheckEssScalingGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_group" {
			continue
		}

		if _, err := essService.DescribeScalingGroupById(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Scaling group %s still exists.", rs.Primary.ID)
	}

	return nil
}

func testAccEssScalingGroup(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup-%d"
	}
	
	resource "alicloud_vswitch" "bar" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.bar.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
		enable = true
		active = true
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		system_disk_category = "cloud_efficiency"
		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = 10
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	
	`, common, rand)
}

func testAccEssScalingGroup_update(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_update-%d"
	}
	
	resource "alicloud_vswitch" "bar" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 2
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.bar.id}"]
		removal_policies = ["OldestInstance"]
	}
	
	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
		enable = true
		active = true
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		system_disk_category = "cloud_efficiency"
		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = 10
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	`, common, rand)
}
func testAccEssScalingGroup_vpc(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_vpc-%d"
	}
	
	resource "alicloud_vswitch" "bar" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.bar.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
		multi_az_policy = "BALANCE"
	}
	
	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
		enable = true
		active = true
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		system_disk_category = "cloud_efficiency"
		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = 10
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = "true"
	}
	`, common, rand)
}

func testAccEssScalingGroup_slb(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_slb-%d"
	}
	
	resource "alicloud_ess_scaling_group" "scaling" {
	  min_size = "1"
	  max_size = "1"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance", "NewestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	  loadbalancer_ids = ["${alicloud_slb.instance.0.id}","${alicloud_slb.instance.1.id}"]
	  depends_on = ["alicloud_slb_listener.tcp"]
	}
	
	resource "alicloud_ess_scaling_configuration" "config" {
	  scaling_group_id = "${alicloud_ess_scaling_group.scaling.id}"
	  active = true
	  enable = true
	  image_id = "${data.alicloud_images.default.images.0.id}"
	  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	  security_group_id = "${alicloud_security_group.default.id}"
	  force_delete = "true"
	  internet_charge_type = "PayByTraffic"
	}
	
	resource "alicloud_slb" "instance" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alicloud_vswitch.default.id}"
	}
	resource "alicloud_slb_listener" "tcp" {
	  count = 2
	  load_balancer_id = "${element(alicloud_slb.instance.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "tcp"
	  bandwidth = "10"
	  health_check_type = "tcp"
	}
	`, common, rand)
}

func testAccEssScalingGroup_slbempty(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_slbempty-%d"
	}
	
	resource "alicloud_ess_scaling_group" "scaling" {
	  min_size = "1"
	  max_size = "1"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance", "NewestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	  loadbalancer_ids = []
	}
	
	resource "alicloud_ess_scaling_configuration" "config" {
	  scaling_group_id = "${alicloud_ess_scaling_group.scaling.id}"
	  active = true
	  enable = true
	  image_id = "${data.alicloud_images.default.images.0.id}"
	  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	  security_group_id = "${alicloud_security_group.default.id}"
	  force_delete = "true"
	  internet_charge_type = "PayByTraffic"
	}
	`, common, rand)
}
