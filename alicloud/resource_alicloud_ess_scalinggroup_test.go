package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
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

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
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
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "1",
		"max_size":           "1",
		"default_cooldown":   "20",
		"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand),
		"vswitch_ids.#":      "2",
		"removal_policies.#": "2",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingGroup(EcsInstanceCommonTestCase, rand),
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
				Config: testAccEssScalingGroupUpdateMaxSize(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_size": "2",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupUpdateScalingGroupName(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					}),
				),
			},
			{
				Config: testAccEssScalingGroupUpdateRemovalPolicies(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"removal_policies.#": "1",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupUpdateDefaultCooldown(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_cooldown": "200",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupUpdateMinSize(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_size": "2",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupModifyVSwitchIds(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccEssScalingGroup(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})

}

func TestAccAlicloudEssScalingGroup_costoptimized(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":                                 "1",
		"max_size":                                 "1",
		"default_cooldown":                         "20",
		"scaling_group_name":                       fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand),
		"vswitch_ids.#":                            "2",
		"removal_policies.#":                       "2",
		"on_demand_base_capacity":                  "10",
		"spot_instance_pools":                      "10",
		"spot_instance_remedy":                     "false",
		"on_demand_percentage_above_base_capacity": "10",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingGroupCostOptimized(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccEssScalingGroupSpotInstanceRemedy(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_instance_remedy": "true",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupOnDemandBaseCapacity(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"on_demand_base_capacity": "8",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupOnDemandPercentageAboveBaseCapacity(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"on_demand_percentage_above_base_capacity": "8",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupSpotInstancePools(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_instance_pools": "8",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudEssScalingGroup_vpc(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "1",
		"max_size":           "1",
		"default_cooldown":   "20",
		"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup_vpc-%d", rand),
		"vswitch_ids.#":      "2",
		"removal_policies.#": "2",
		"multi_az_policy":    "BALANCE",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingGroupVpc(EcsInstanceCommonTestCase, rand),
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
				Config: testAccEssScalingGroupVpcUpdateMaxSize(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_size": "2",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupVpcUpdateScalingGroupName(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					}),
				),
			},
			{
				Config: testAccEssScalingGroupVpcUpdateRemovalPolicies(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"removal_policies.#": "1",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupVpcUpdateDefaultCooldown(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_cooldown": "200",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupVpcUpdateMinSize(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_size": "2",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupVpc(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})

}

func TestAccAlicloudEssScalingGroup_slb(t *testing.T) {
	var v ess.ScalingGroup
	var slb *slb.DescribeLoadBalancerAttributeResponse
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "1",
		"max_size":           "1",
		"default_cooldown":   "300",
		"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup_slb-%d", rand),
		"vswitch_ids.#":      "1",
		"removal_policies.#": "2",
		"multi_az_policy":    "PRIORITY",
		"loadbalancer_ids.#": "0",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rcSlb0 := resourceCheckInit("alicloud_slb.default.0", &slb, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rcSlb1 := resourceCheckInit("alicloud_slb.default.1", &slb, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingGroupSlbempty(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccEssScalingGroupSlb(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"loadbalancer_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupSlbDetach(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					testAccCheck(map[string]string{
						"loadbalancer_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupSlbUpdateMaxSize(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"max_size":           "2",
						"loadbalancer_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupSlbUpdateScalingGroupName(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					}),
				),
			},
			{
				Config: testAccEssScalingGroupSlbUpdateRemovalPolicies(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"removal_policies.#": "1",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupSlbUpdateDefaultCooldown(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"default_cooldown": "200",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupSlbUpdateMinSize(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"min_size": "2",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupSlbempty(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"loadbalancer_ids.#": "0",
						"min_size":           "1",
						"max_size":           "1",
						"default_cooldown":   "300",
						"removal_policies.#": "2",
						"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup_slb-%d", rand),
					}),
				),
			},
		},
	})

}

func testAccCheckEssScalingGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_group" {
			continue
		}

		if _, err := essService.DescribeEssScalingGroup(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(fmt.Errorf("Scaling group %s still exists.", rs.Primary.ID))
	}

	return nil
}

func testAccEssScalingGroup(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupUpdateMaxSize(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupUpdateScalingGroupName(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupUpdateRemovalPolicies(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupUpdateDefaultCooldown(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 200
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupUpdateMinSize(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 2
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 200
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
	}`, common, rand)
}
func testAccEssScalingGroupVpc(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_vpc-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
		multi_az_policy = "BALANCE"
	}`, common, rand)
}

func testAccEssScalingGroupVpcUpdateMaxSize(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_vpc-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
		multi_az_policy = "BALANCE"
	}`, common, rand)
}

func testAccEssScalingGroupVpcUpdateScalingGroupName(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
		multi_az_policy = "BALANCE"
	}`, common, rand)
}

func testAccEssScalingGroupVpcUpdateRemovalPolicies(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
		multi_az_policy = "BALANCE"
	}`, common, rand)
}

func testAccEssScalingGroupVpcUpdateDefaultCooldown(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 200
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
		multi_az_policy = "BALANCE"
	}`, common, rand)
}

func testAccEssScalingGroupVpcUpdateMinSize(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 2
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 200
		vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
		multi_az_policy = "BALANCE"
	}`, common, rand)
}

func testAccEssScalingGroupSlb(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_slb-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "1"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance", "NewestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	  loadbalancer_ids = ["${alicloud_slb.default.0.id}","${alicloud_slb.default.1.id}"]
	  depends_on = ["alicloud_slb_listener.default"]
	}

	resource "alicloud_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_slb_listener" "default" {
	  count = 2
	  load_balancer_id = "${element(alicloud_slb.default.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "tcp"
	  bandwidth = "10"
	  health_check_type = "tcp"
	}
	`, common, rand)
}

func testAccEssScalingGroupSlbDetach(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_slb-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "1"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance", "NewestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	  loadbalancer_ids = ["${alicloud_slb.default.0.id}"]
	  depends_on = ["alicloud_slb_listener.default"]
	}

	resource "alicloud_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_slb_listener" "default" {
	  count = 2
	  load_balancer_id = "${element(alicloud_slb.default.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "tcp"
	  bandwidth = "10"
	  health_check_type = "tcp"
	}
	`, common, rand)
}

func testAccEssScalingGroupSlbempty(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_slb-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "1"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance", "NewestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	  loadbalancer_ids = []
	}`, common, rand)
}

func testAccEssScalingGroupSlbUpdateMaxSize(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_slb-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "2"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance", "NewestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	  loadbalancer_ids = ["${alicloud_slb.default.0.id}","${alicloud_slb.default.1.id}"]
	  depends_on = ["alicloud_slb_listener.default"]
	}

	resource "alicloud_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_slb_listener" "default" {
	  count = 2
	  load_balancer_id = "${element(alicloud_slb.default.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "tcp"
	  bandwidth = "10"
	  health_check_type = "tcp"
	}
	`, common, rand)
}

func testAccEssScalingGroupSlbUpdateScalingGroupName(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "2"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance", "NewestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	  loadbalancer_ids = ["${alicloud_slb.default.0.id}","${alicloud_slb.default.1.id}"]
	  depends_on = ["alicloud_slb_listener.default"]
	}

	resource "alicloud_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_slb_listener" "default" {
	  count = 2
	  load_balancer_id = "${element(alicloud_slb.default.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "tcp"
	  bandwidth = "10"
	  health_check_type = "tcp"
	}
	`, common, rand)
}

func testAccEssScalingGroupSlbUpdateRemovalPolicies(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "2"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	  loadbalancer_ids = ["${alicloud_slb.default.0.id}","${alicloud_slb.default.1.id}"]
	  depends_on = ["alicloud_slb_listener.default"]
	}

	resource "alicloud_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_slb_listener" "default" {
	  count = 2
	  load_balancer_id = "${element(alicloud_slb.default.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "tcp"
	  bandwidth = "10"
	  health_check_type = "tcp"
	}
	`, common, rand)
}

func testAccEssScalingGroupSlbUpdateDefaultCooldown(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	  loadbalancer_ids = ["${alicloud_slb.default.0.id}","${alicloud_slb.default.1.id}"]
	  depends_on = ["alicloud_slb_listener.default"]
	}

	resource "alicloud_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_slb_listener" "default" {
	  count = 2
	  load_balancer_id = "${element(alicloud_slb.default.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "tcp"
	  bandwidth = "10"
	  health_check_type = "tcp"
	}
	`, common, rand)
}

func testAccEssScalingGroupSlbUpdateMinSize(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "2"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	  loadbalancer_ids = ["${alicloud_slb.default.0.id}","${alicloud_slb.default.1.id}"]
	  depends_on = ["alicloud_slb_listener.default"]
	}

	resource "alicloud_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_slb_listener" "default" {
	  count = 2
	  load_balancer_id = "${element(alicloud_slb.default.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "tcp"
	  bandwidth = "10"
	  health_check_type = "tcp"
	}
	`, common, rand)
}

func testAccEssScalingGroupModifyVSwitchIds(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 2
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 200
		vswitch_ids = ["${alicloud_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupCostOptimized(common string, rand int) string {
	return fmt.Sprintf(`
    %s
    variable "name" {
        default = "tf-testAccEssScalingGroup-%d"
    }
    
    resource "alicloud_vswitch" "default2" {
          vpc_id = "${alicloud_vpc.default.id}"
          cidr_block = "172.16.1.0/24"
          availability_zone = "${data.alicloud_zones.default.zones.0.id}"
          name = "${var.name}-bar"
    }
    
    resource "alicloud_ess_scaling_group" "default" {
        min_size = 1
        max_size = 1
        scaling_group_name = "${var.name}"
        default_cooldown = 20
        vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
        removal_policies = ["OldestInstance", "NewestInstance"]
        multi_az_policy = "COST_OPTIMIZED"
        on_demand_base_capacity = "10"
        on_demand_percentage_above_base_capacity = "10"
        spot_instance_pools = "10"
    }`, common, rand)
}

func testAccEssScalingGroupSpotInstanceRemedy(common string, rand int) string {
	return fmt.Sprintf(`
    %s
    variable "name" {
        default = "tf-testAccEssScalingGroup-%d"
    }
    
    resource "alicloud_vswitch" "default2" {
          vpc_id = "${alicloud_vpc.default.id}"
          cidr_block = "172.16.1.0/24"
          availability_zone = "${data.alicloud_zones.default.zones.0.id}"
          name = "${var.name}-bar"
    }
    
    resource "alicloud_ess_scaling_group" "default" {
        min_size = 1
        max_size = 1
        scaling_group_name = "${var.name}"
        default_cooldown = 20
        vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
        removal_policies = ["OldestInstance", "NewestInstance"]
        multi_az_policy = "COST_OPTIMIZED"
        on_demand_base_capacity = "10"
        on_demand_percentage_above_base_capacity = "10"
        spot_instance_pools = "10"
		spot_instance_remedy = true
    }`, common, rand)
}

func testAccEssScalingGroupOnDemandBaseCapacity(common string, rand int) string {
	return fmt.Sprintf(`
    %s
    variable "name" {
        default = "tf-testAccEssScalingGroup-%d"
    }
    
    resource "alicloud_vswitch" "default2" {
          vpc_id = "${alicloud_vpc.default.id}"
          cidr_block = "172.16.1.0/24"
          availability_zone = "${data.alicloud_zones.default.zones.0.id}"
          name = "${var.name}-bar"
    }
    
    resource "alicloud_ess_scaling_group" "default" {
        min_size = 1
        max_size = 1
        scaling_group_name = "${var.name}"
        default_cooldown = 20
        vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
        removal_policies = ["OldestInstance", "NewestInstance"]
        multi_az_policy = "COST_OPTIMIZED"
        on_demand_base_capacity = "8"
        on_demand_percentage_above_base_capacity = "10"
        spot_instance_pools = "10"
		spot_instance_remedy = true
    }`, common, rand)
}

func testAccEssScalingGroupOnDemandPercentageAboveBaseCapacity(common string, rand int) string {
	return fmt.Sprintf(`
    %s
    variable "name" {
        default = "tf-testAccEssScalingGroup-%d"
    }
    
    resource "alicloud_vswitch" "default2" {
          vpc_id = "${alicloud_vpc.default.id}"
          cidr_block = "172.16.1.0/24"
          availability_zone = "${data.alicloud_zones.default.zones.0.id}"
          name = "${var.name}-bar"
    }
    
    resource "alicloud_ess_scaling_group" "default" {
        min_size = 1
        max_size = 1
        scaling_group_name = "${var.name}"
        default_cooldown = 20
        vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
        removal_policies = ["OldestInstance", "NewestInstance"]
        multi_az_policy = "COST_OPTIMIZED"
        on_demand_base_capacity = "8"
        on_demand_percentage_above_base_capacity = "8"
        spot_instance_pools = "10"
		spot_instance_remedy = true
    }`, common, rand)
}

func testAccEssScalingGroupSpotInstancePools(common string, rand int) string {
	return fmt.Sprintf(`
    %s
    variable "name" {
        default = "tf-testAccEssScalingGroup-%d"
    }
    
    resource "alicloud_vswitch" "default2" {
          vpc_id = "${alicloud_vpc.default.id}"
          cidr_block = "172.16.1.0/24"
          availability_zone = "${data.alicloud_zones.default.zones.0.id}"
          name = "${var.name}-bar"
    }
    
    resource "alicloud_ess_scaling_group" "default" {
        min_size = 1
        max_size = 1
        scaling_group_name = "${var.name}"
        default_cooldown = 20
        vswitch_ids = ["${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"]
        removal_policies = ["OldestInstance", "NewestInstance"]
        multi_az_policy = "COST_OPTIMIZED"
        on_demand_base_capacity = "8"
        on_demand_percentage_above_base_capacity = "8"
        spot_instance_pools = "8"
		spot_instance_remedy = true
    }`, common, rand)
}
