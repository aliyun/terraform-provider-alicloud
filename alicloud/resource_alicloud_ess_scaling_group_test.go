package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_ess_scaling_group", &resource.Sweeper{
		Name: "alicloud_ess_scaling_group",
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

	for _, v := range groups {
		name := v.ScalingGroupName
		id := v.ScalingGroupId
		skip := true
		if !sweepAll() {
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
		}
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
	return nil
}

func TestAccAliCloudEssScalingGroup_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "1",
		"max_size":           "4",
		"desired_capacity":   "2",
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
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupDependence)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "4",
					"desired_capacity":   "2",
					"scaling_group_name": "${var.name}",
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
				}),
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "5",
					"desired_capacity":   "2",
					"scaling_group_name": "${var.name}",
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_size": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "5",
					"desired_capacity":   "3",
					"scaling_group_name": "${var.name}",
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_capacity": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "5",
					"desired_capacity":   "3",
					"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "5",
					"desired_capacity":   "3",
					"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"removal_policies.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "5",
					"desired_capacity":   "3",
					"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					"default_cooldown":   "200",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_cooldown": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "2",
					"max_size":           "5",
					"desired_capacity":   "3",
					"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					"default_cooldown":   "200",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_size": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "2",
					"max_size":           "5",
					"desired_capacity":   "3",
					"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					"default_cooldown":   "200",
					"vswitch_ids":        []string{"${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "4",
					"desired_capacity":   "2",
					"scaling_group_name": "${var.name}",
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_desiredCapacity(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "1",
		"max_size":           "4",
		"desired_capacity":   "2",
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
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupDependence)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "4",
					"desired_capacity":   "2",
					"scaling_group_name": "${var.name}",
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
				}),
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "5",
					"desired_capacity":   "3",
					"scaling_group_name": "${var.name}",
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_size":         "5",
						"desired_capacity": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "5",
					"desired_capacity":   "0",
					"scaling_group_name": "${var.name}",
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_size":         "0",
						"desired_capacity": "0",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingGroup_max_min_desiredCapacityRange(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "1",
		"max_size":           "4",
		"desired_capacity":   "2",
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
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupDependence)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "4",
					"desired_capacity":   "2",
					"scaling_group_name": "${var.name}",
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "2000",
					"max_size":           "2000",
					"desired_capacity":   "2000",
					"scaling_group_name": "${var.name}",
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_size":         "2000",
						"min_size":         "2000",
						"desired_capacity": "2000",
					}),
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

func TestAccAliClouddEssScalingGroup_withLaunchTemplateId(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"
	//checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)

	basicMap := map[string]string{
		"min_size":                "0",
		"max_size":                "4",
		"default_cooldown":        "20",
		"scaling_group_name":      fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand),
		"vswitch_id":              CHECKSET,
		"vswitch_ids.#":           "1",
		"removal_policies.#":      "2",
		"launch_template_version": "Default",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupTemplate)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "0",
					"max_size":                "4",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_id":              "${alicloud_vswitch.tmpVs.id}",
					"vswitch_ids":             []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"launch_template_id":      "${alicloud_ecs_launch_template.default3.id}",
					"launch_template_version": "Default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "0",
					"max_size":                "4",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"launch_template_id":      "${alicloud_ecs_launch_template.default3.id}",
					"launch_template_version": "Latest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"launch_template_id":      CHECKSET,
						"vswitch_id":              CHECKSET,
						"vswitch_ids.#":           "1",
						"launch_template_version": "Latest",
					}),
				),
			},
		},
	})

}

func TestAccAliClouddEssScalingGroup_withLaunchTemplateOverride(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"
	//checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)

	basicMap := map[string]string{
		"min_size":                "0",
		"max_size":                "4",
		"default_cooldown":        "20",
		"scaling_group_name":      fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand),
		"vswitch_id":              CHECKSET,
		"vswitch_ids.#":           "1",
		"removal_policies.#":      "2",
		"launch_template_version": "Default",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupTemplate)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "0",
					"max_size":                "4",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_id":              "${alicloud_vswitch.tmpVs.id}",
					"vswitch_ids":             []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"launch_template_id":      "${alicloud_ecs_launch_template.default3.id}",
					"launch_template_version": "Default",
					"launch_template_override": []map[string]string{{
						"instance_type": "${data.alicloud_instance_types.default.instance_types.0.id}",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "0",
					"max_size":                "4",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"launch_template_id":      "${alicloud_ecs_launch_template.default3.id}",
					"launch_template_version": "Latest",
					"launch_template_override": []map[string]string{{
						"instance_type":     "${data.alicloud_instance_types.default.instance_types.1.id}",
						"weighted_capacity": "4",
						"spot_price_limit":  "2.1",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"launch_template_id":         CHECKSET,
						"vswitch_id":                 CHECKSET,
						"vswitch_ids.#":              "1",
						"launch_template_override.#": "1",
						"launch_template_version":    "Latest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "0",
					"max_size":                "4",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"launch_template_id":      "${alicloud_ecs_launch_template.default4.id}",
					"launch_template_version": "Latest",
					"launch_template_override": []map[string]string{{
						"instance_type":     "${data.alicloud_instance_types.default.instance_types.0.id}",
						"weighted_capacity": "3",
						"spot_price_limit":  "1.2",
					},
						{
							"instance_type":     "${data.alicloud_instance_types.default.instance_types.1.id}",
							"weighted_capacity": "2",
							"spot_price_limit":  "1.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"launch_template_id":         CHECKSET,
						"vswitch_id":                 CHECKSET,
						"vswitch_ids.#":              "1",
						"launch_template_override.#": "2",
						"launch_template_version":    "Latest",
					}),
				),
			},
		},
	})

}

func TestAccAliClouddEssScalingGroup_withAlbServerGroup(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":                "0",
		"max_size":                "4",
		"default_cooldown":        "20",
		"scaling_group_name":      fmt.Sprintf("tf-testAccEssScalingGroupAlb-%d", rand),
		"vswitch_id":              CHECKSET,
		"vswitch_ids.#":           "1",
		"removal_policies.#":      "2",
		"launch_template_version": "Default",
		"alb_server_group.#":      "0",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupAlb-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupAlbServerGroup)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "0",
					"max_size":                "4",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_id":              "${alicloud_vswitch.tmpVs.id}",
					"vswitch_ids":             []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"launch_template_id":      "${alicloud_ecs_launch_template.default3.id}",
					"launch_template_version": "Default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":              []string{"alicloud_alb_server_group.default"},
					"min_size":                "0",
					"max_size":                "4",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"launch_template_id":      "${alicloud_ecs_launch_template.default4.id}",
					"launch_template_version": "Latest",
					"alb_server_group": []map[string]string{{
						"alb_server_group_id": "${alicloud_alb_server_group.default.0.id}",
						"weight":              "10",
						"port":                "20",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"launch_template_id":      CHECKSET,
						"vswitch_id":              CHECKSET,
						"vswitch_ids.#":           "1",
						"alb_server_group.#":      "1",
						"launch_template_version": "Latest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":              []string{"alicloud_alb_server_group.default"},
					"min_size":                "0",
					"max_size":                "4",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"launch_template_id":      "${alicloud_ecs_launch_template.default4.id}",
					"launch_template_version": "Latest",
					"alb_server_group": []map[string]string{{
						"alb_server_group_id": "${alicloud_alb_server_group.default.0.id}",
						"weight":              "10",
						"port":                "20",
					},
						{
							"alb_server_group_id": "${alicloud_alb_server_group.default.1.id}",
							"weight":              "20",
							"port":                "30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"launch_template_id":      CHECKSET,
						"vswitch_id":              CHECKSET,
						"vswitch_ids.#":           "1",
						"alb_server_group.#":      "2",
						"launch_template_version": "Latest",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":              []string{"alicloud_alb_server_group.default"},
					"min_size":                "0",
					"max_size":                "4",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"launch_template_id":      "${alicloud_ecs_launch_template.default4.id}",
					"launch_template_version": "Latest",
					"alb_server_group":        REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"launch_template_id":      CHECKSET,
						"vswitch_id":              CHECKSET,
						"vswitch_ids.#":           "1",
						"alb_server_group.#":      "0",
						"launch_template_version": "Latest",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":              []string{"alicloud_alb_server_group.default"},
					"min_size":                "0",
					"max_size":                "4",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"launch_template_id":      "${alicloud_ecs_launch_template.default4.id}",
					"launch_template_version": "Default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"launch_template_id":      CHECKSET,
						"vswitch_id":              CHECKSET,
						"vswitch_ids.#":           "1",
						"launch_template_version": "Default",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_costoptimized(t *testing.T) {
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
		"group_deletion_protection":                "false",
		"on_demand_percentage_above_base_capacity": "10",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupDependence)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COST_OPTIMIZED",
					"on_demand_base_capacity": "10",
					"on_demand_percentage_above_base_capacity": "10",
					"spot_instance_pools":                      "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COST_OPTIMIZED",
					"on_demand_base_capacity": "10",
					"on_demand_percentage_above_base_capacity": "10",
					"spot_instance_pools":                      "10",
					"spot_instance_remedy":                     "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_instance_remedy": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COST_OPTIMIZED",
					"on_demand_base_capacity": "10",
					"on_demand_percentage_above_base_capacity": "10",
					"spot_instance_pools":                      "10",
					"spot_instance_remedy":                     "true",
					"group_deletion_protection":                "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_deletion_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COST_OPTIMIZED",
					"on_demand_base_capacity": "0",
					"on_demand_percentage_above_base_capacity": "0",
					"spot_instance_pools":                      "10",
					"spot_instance_remedy":                     "true",
					"group_deletion_protection":                "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"on_demand_base_capacity":                  "0",
						"on_demand_percentage_above_base_capacity": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COST_OPTIMIZED",
					"on_demand_base_capacity": "8",
					"on_demand_percentage_above_base_capacity": "8",
					"spot_instance_pools":                      "10",
					"spot_instance_remedy":                     "true",
					"group_deletion_protection":                "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"on_demand_base_capacity":                  "8",
						"on_demand_percentage_above_base_capacity": "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COST_OPTIMIZED",
					"on_demand_base_capacity": "8",
					"on_demand_percentage_above_base_capacity": "8",
					"spot_instance_pools":                      "8",
					"spot_instance_remedy":                     "true",
					"group_deletion_protection":                "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_instance_pools":       "8",
						"group_deletion_protection": "false",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_composable(t *testing.T) {
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
		"group_deletion_protection":                "false",
		"on_demand_percentage_above_base_capacity": "10",
		"az_balance":                               "false",
		"allocation_strategy":                      "priority",
		"spot_allocation_strategy":                 "priority",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupDependence)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COMPOSABLE",
					"on_demand_base_capacity": "10",
					"on_demand_percentage_above_base_capacity": "10",
					"spot_instance_pools":                      "10",
					"az_balance":                               "false",
					"allocation_strategy":                      "priority",
					"spot_allocation_strategy":                 "priority",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COMPOSABLE",
					"on_demand_base_capacity": "10",
					"on_demand_percentage_above_base_capacity": "10",
					"spot_instance_pools":                      "10",
					"spot_instance_remedy":                     "true",
					"az_balance":                               "false",
					"allocation_strategy":                      "priority",
					"spot_allocation_strategy":                 "priority",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_instance_remedy":     "true",
						"az_balance":               "false",
						"allocation_strategy":      "priority",
						"spot_allocation_strategy": "priority",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COMPOSABLE",
					"on_demand_base_capacity": "10",
					"on_demand_percentage_above_base_capacity": "10",
					"spot_instance_pools":                      "10",
					"spot_instance_remedy":                     "true",
					"group_deletion_protection":                "true",
					"az_balance":                               "true",
					"allocation_strategy":                      "lowestPrice",
					"spot_allocation_strategy":                 "lowestPrice",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_deletion_protection": "true",
						"az_balance":                "true",
						"allocation_strategy":       "lowestPrice",
						"spot_allocation_strategy":  "lowestPrice",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COMPOSABLE",
					"on_demand_base_capacity": "0",
					"on_demand_percentage_above_base_capacity": "0",
					"spot_instance_pools":                      "10",
					"spot_instance_remedy":                     "true",
					"group_deletion_protection":                "true",
					"az_balance":                               "false",
					"allocation_strategy":                      "priority",
					"spot_allocation_strategy":                 "priority",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"on_demand_base_capacity":                  "0",
						"on_demand_percentage_above_base_capacity": "0",
						"az_balance":               "false",
						"allocation_strategy":      "priority",
						"spot_allocation_strategy": "priority",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COMPOSABLE",
					"on_demand_base_capacity": "8",
					"on_demand_percentage_above_base_capacity": "8",
					"spot_instance_pools":                      "10",
					"spot_instance_remedy":                     "true",
					"group_deletion_protection":                "true",
					"az_balance":                               "false",
					"allocation_strategy":                      "priority",
					"spot_allocation_strategy":                 "priority",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"on_demand_base_capacity":                  "8",
						"on_demand_percentage_above_base_capacity": "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COMPOSABLE",
					"on_demand_base_capacity": "8",
					"on_demand_percentage_above_base_capacity": "8",
					"spot_instance_pools":                      "8",
					"spot_instance_remedy":                     "true",
					"group_deletion_protection":                "false",
					"az_balance":                               "false",
					"allocation_strategy":                      "priority",
					"spot_allocation_strategy":                 "priority",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_instance_pools":       "8",
						"group_deletion_protection": "false",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_scalingPolicy_maxInstanceLifetime(t *testing.T) {
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
		"group_deletion_protection":                "false",
		"on_demand_percentage_above_base_capacity": "10",
		"az_balance":                               "false",
		"allocation_strategy":                      "priority",
		"spot_allocation_strategy":                 "priority",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupDependence)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COMPOSABLE",
					"on_demand_base_capacity": "10",
					"on_demand_percentage_above_base_capacity": "10",
					"spot_instance_pools":                      "10",
					"az_balance":                               "false",
					"allocation_strategy":                      "priority",
					"scaling_policy":                           "release",
					"spot_allocation_strategy":                 "priority",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_policy":          "forceRelease",
					"max_instance_lifetime":   "86400",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COMPOSABLE",
					"on_demand_base_capacity": "10",
					"on_demand_percentage_above_base_capacity": "10",
					"spot_instance_pools":                      "10",
					"spot_instance_remedy":                     "true",
					"az_balance":                               "false",
					"allocation_strategy":                      "priority",
					"spot_allocation_strategy":                 "priority",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_instance_remedy":     "true",
						"az_balance":               "false",
						"scaling_policy":           "forceRelease",
						"max_instance_lifetime":    "86400",
						"allocation_strategy":      "priority",
						"spot_allocation_strategy": "priority",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"scaling_policy":          "forceRelease",
					"max_instance_lifetime":   "86411",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COMPOSABLE",
					"on_demand_base_capacity": "10",
					"on_demand_percentage_above_base_capacity": "10",
					"spot_instance_pools":                      "10",
					"spot_instance_remedy":                     "true",
					"group_deletion_protection":                "false",
					"az_balance":                               "true",
					"allocation_strategy":                      "lowestPrice",
					"spot_allocation_strategy":                 "lowestPrice",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_deletion_protection": "false",
						"az_balance":                "true",
						"scaling_policy":            "forceRelease",
						"max_instance_lifetime":     "86411",
						"allocation_strategy":       "lowestPrice",
						"spot_allocation_strategy":  "lowestPrice",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"scaling_policy":          "forceRecycle",
					"max_instance_lifetime":   REMOVEKEY,
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COMPOSABLE",
					"on_demand_base_capacity": "10",
					"on_demand_percentage_above_base_capacity": "10",
					"spot_instance_pools":                      "10",
					"spot_instance_remedy":                     "true",
					"group_deletion_protection":                "false",
					"az_balance":                               "true",
					"allocation_strategy":                      "lowestPrice",
					"spot_allocation_strategy":                 "lowestPrice",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_deletion_protection": "false",
						"az_balance":                "true",
						"scaling_policy":            "forceRecycle",
						"max_instance_lifetime":     REMOVEKEY,
						"allocation_strategy":       "lowestPrice",
						"spot_allocation_strategy":  "lowestPrice",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_scalingPolicy_stopInstanceTimeout(t *testing.T) {
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
		"group_deletion_protection":                "false",
		"on_demand_percentage_above_base_capacity": "10",
		"az_balance":                               "false",
		"allocation_strategy":                      "priority",
		"spot_allocation_strategy":                 "priority",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupDependence)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COMPOSABLE",
					"on_demand_base_capacity": "10",
					"on_demand_percentage_above_base_capacity": "10",
					"spot_instance_pools":                      "10",
					"az_balance":                               "false",
					"allocation_strategy":                      "priority",
					"scaling_policy":                           "release",
					"spot_allocation_strategy":                 "priority",
					"stop_instance_timeout":                    "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":                "1",
					"max_size":                "1",
					"scaling_policy":          "forceRelease",
					"max_instance_lifetime":   "86400",
					"scaling_group_name":      "${var.name}",
					"default_cooldown":        "20",
					"vswitch_ids":             []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":        []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":         "COMPOSABLE",
					"on_demand_base_capacity": "10",
					"on_demand_percentage_above_base_capacity": "10",
					"spot_instance_pools":                      "10",
					"spot_instance_remedy":                     "true",
					"az_balance":                               "false",
					"allocation_strategy":                      "priority",
					"spot_allocation_strategy":                 "priority",
					"stop_instance_timeout":                    "240",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_instance_remedy":     "true",
						"az_balance":               "false",
						"scaling_policy":           "forceRelease",
						"max_instance_lifetime":    "86400",
						"allocation_strategy":      "priority",
						"spot_allocation_strategy": "priority",
						"stop_instance_timeout":    "240",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_vpc(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssScalingGroup_vpc-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupDependence)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "1",
					"scaling_group_name": "${var.name}",
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":    "BALANCE",
				}),
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "2",
					"scaling_group_name": "${var.name}",
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":    "BALANCE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_size": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "2",
					"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":    "BALANCE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "2",
					"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance"},
					"multi_az_policy":    "BALANCE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"removal_policies.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "2",
					"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					"default_cooldown":   "200",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance"},
					"multi_az_policy":    "BALANCE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_cooldown": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "2",
					"max_size":           "2",
					"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					"default_cooldown":   "200",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance"},
					"multi_az_policy":    "BALANCE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_size": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "1",
					"scaling_group_name": "${var.name}",
					"default_cooldown":   "20",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default2.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"multi_az_policy":    "BALANCE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_slb(t *testing.T) {
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
	rcSlb0 := resourceCheckInit("alicloud_slb_load_balancer.default.0", &slb, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rcSlb1 := resourceCheckInit("alicloud_slb_load_balancer.default.1", &slb, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup_slb-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupSlbDependence)
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
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_slb_load_balancer.default"},
					"min_size":           "1",
					"max_size":           "1",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"loadbalancer_ids":   []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_slb_load_balancer.default"},
					"min_size":           "1",
					"max_size":           "1",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"loadbalancer_ids":   []string{"${alicloud_slb_load_balancer.default.0.id}", "${alicloud_slb_load_balancer.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"loadbalancer_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_slb_load_balancer.default"},
					"min_size":           "1",
					"max_size":           "1",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"loadbalancer_ids":   []string{"${alicloud_slb_load_balancer.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					testAccCheck(map[string]string{
						"loadbalancer_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_slb_load_balancer.default"},
					"min_size":           "1",
					"max_size":           "2",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"loadbalancer_ids":   []string{"${alicloud_slb_load_balancer.default.0.id}", "${alicloud_slb_load_balancer.default.1.id}"},
				}),
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
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_slb_load_balancer.default"},
					"min_size":           "1",
					"max_size":           "2",
					"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"loadbalancer_ids":   []string{"${alicloud_slb_load_balancer.default.0.id}", "${alicloud_slb_load_balancer.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_slb_load_balancer.default"},
					"min_size":           "1",
					"max_size":           "2",
					"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance"},
					"loadbalancer_ids":   []string{"${alicloud_slb_load_balancer.default.0.id}", "${alicloud_slb_load_balancer.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"removal_policies.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_slb_load_balancer.default"},
					"min_size":           "1",
					"max_size":           "2",
					"default_cooldown":   "200",
					"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance"},
					"loadbalancer_ids":   []string{"${alicloud_slb_load_balancer.default.0.id}", "${alicloud_slb_load_balancer.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"default_cooldown": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_slb_load_balancer.default"},
					"min_size":           "2",
					"max_size":           "2",
					"default_cooldown":   "200",
					"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance"},
					"loadbalancer_ids":   []string{"${alicloud_slb_load_balancer.default.0.id}", "${alicloud_slb_load_balancer.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"min_size": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "1",
					"max_size":           "1",
					"default_cooldown":   "300",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"loadbalancer_ids":   []string{},
				}),
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

func TestAccAliCloudEssScalingGroup_tags(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "0",
		"max_size":           "4",
		"default_cooldown":   "20",
		"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand),
		"vswitch_ids.#":      "1",
		"removal_policies.#": "2",
		"tags.key1":          "value1",
		"tags.key2":          "value2",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupTemplate)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "4",
					"default_cooldown":   "20",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"tags":               map[string]string{"key1": "value1", "key2": "value2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "4",
					"default_cooldown":   "20",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"tags":               map[string]string{"key1": "value2", "key2": "value1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.key1": "value2",
						"tags.key2": "value1",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_healthCheckTypes(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "0",
		"max_size":           "4",
		"default_cooldown":   "20",
		"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand),
		"vswitch_ids.#":      "1",
		"removal_policies.#": "2",
		"health_check_type":  "NONE",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupTemplate)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "4",
					"default_cooldown":   "20",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"health_check_type":  "NONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "4",
					"default_cooldown":   "20",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"health_check_type":  REMOVEKEY,
					"health_check_types": []string{"ECS"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_type":    REMOVEKEY,
						"health_check_types.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "4",
					"default_cooldown":   "20",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"health_check_types": []string{"ECS", "LOAD_BALANCER"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_types.#": "2",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_resourceGroupId(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "0",
		"max_size":           "4",
		"default_cooldown":   "20",
		"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand),
		"vswitch_ids.#":      "1",
		"removal_policies.#": "2",
		"health_check_type":  "NONE",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupResourceGroupId)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "4",
					"default_cooldown":   "20",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"health_check_type":  "NONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "4",
					"default_cooldown":   "20",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"health_check_type":  "LOAD_BALANCER",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_type": "LOAD_BALANCER",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "4",
					"default_cooldown":   "20",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"health_check_type":  "ECS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_type": "ECS",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "4",
					"default_cooldown":   "20",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"resource_group_id":  "",
					"health_check_type":  "ECS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_type": "ECS",
						"resource_group_id": CHECKSET,
					}),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_eci(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "0",
		"max_size":           "4",
		"default_cooldown":   "20",
		"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand),
		"vswitch_ids.#":      "1",
		"removal_policies.#": "2",
		"group_type":         "ECI",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupTemplate)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "4",
					"default_cooldown":   "20",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"group_type":         "ECI",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "0",
					"default_cooldown":   "20",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"group_type":         "ECI",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_size": "0",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_eciInstance(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "0",
		"max_size":           "4",
		"default_cooldown":   "20",
		"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand),
		"vswitch_ids.#":      "1",
		"removal_policies.#": "2",
		"group_type":         "ECI",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupInstance)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "4",
					"default_cooldown":   "20",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"group_type":         "ECI",
					"container_group_id": "${alicloud_eci_container_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_ecsInstance(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "0",
		"max_size":           "4",
		"default_cooldown":   "20",
		"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand),
		"vswitch_ids.#":      "1",
		"removal_policies.#": "2",
		"group_type":         "ECS",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupInstance)
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
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "4",
					"default_cooldown":   "20",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.tmpVs.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"group_type":         "ECS",
					"instance_id":        "${alicloud_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_protected_instances(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alicloud_ess_scaling_group.default"
	//checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)

	basicMap := map[string]string{
		"min_size": "0",
		"max_size": "2",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssAttachmentConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingAttachmentConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           "0",
					"max_size":           "2",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":            "0",
					"max_size":            "2",
					"scaling_group_name":  "${var.name}",
					"vswitch_ids":         []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":    []string{"OldestInstance", "NewestInstance"},
					"protected_instances": []string{"${alicloud_instance.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protected_instances.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":            "0",
					"max_size":            "2",
					"scaling_group_name":  "${var.name}",
					"vswitch_ids":         []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":    []string{"OldestInstance", "NewestInstance"},
					"protected_instances": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protected_instances.#": "0",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudEssScalingGroup_rds(t *testing.T) {
	var v ess.ScalingGroup
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "alicloud_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "1",
		"max_size":           "1",
		"default_cooldown":   "300",
		"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup_rds-%d", rand),
		"vswitch_ids.#":      "1",
		"removal_policies.#": "2",
		"multi_az_policy":    "PRIORITY",
		"db_instance_ids.#":  "0",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})

	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroup_rds-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingGroupRdsDependence)
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
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_db_instance.default"},
					"min_size":           "1",
					"max_size":           "1",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"db_instance_ids":    []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_db_instance.default"},
					"min_size":           "1",
					"max_size":           "1",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"db_instance_ids":    []string{"${alicloud_db_instance.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":         []string{"alicloud_db_instance.default"},
					"min_size":           "1",
					"max_size":           "1",
					"default_cooldown":   "300",
					"scaling_group_name": "${var.name}",
					"vswitch_ids":        []string{"${alicloud_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"db_instance_ids":    []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ids.#":  "0",
						"min_size":           "1",
						"max_size":           "1",
						"default_cooldown":   "300",
						"removal_policies.#": "2",
						"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup_rds-%d", rand),
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

func resourceEssScalingGroupDependence(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}

	resource "alicloud_vswitch" "default2" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}`, EcsInstanceCommonTestCase, name)
}

func resourceEssScalingAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s" 
	}

    data "alicloud_images" "default1" {
		name_regex  = "^centos.*_64"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_instance_types" "c6" {
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}

	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		image_id = "${data.alicloud_images.default1.images.0.id}"
		instance_type = "${data.alicloud_instance_types.c6.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = true
		active = true
		enable = true
	}



	resource "alicloud_instance" "default" {
		image_id = "${data.alicloud_images.default1.images.0.id}"
		instance_type = "${data.alicloud_instance_types.c6.instance_types.0.id}"
		count = 1
		security_groups = ["${alicloud_security_group.default.id}"]
		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = "10"
		instance_charge_type = "PostPaid"
		system_disk_category = "cloud_efficiency"
		vswitch_id = "${alicloud_vswitch.default.id}"
		instance_name = "${var.name}"
	}

	resource "alicloud_ess_attachment" "default" {
		scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
		instance_ids = ["${alicloud_instance.default.0.id}"]
		force = true
		depends_on = [
        "alicloud_ess_scaling_group.default",
        "alicloud_instance.default"]
	}

	`, EcsInstanceCommonTestCase, name)
}

func resourceEssScalingGroupSlbDependence(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}

	resource "alicloud_slb_load_balancer" "default" {
	  count=2
	  load_balancer_name = "${var.name}"
	  vswitch_id = "${alicloud_vswitch.default.id}"
      load_balancer_spec = "slb.s1.small"
	}

	resource "alicloud_slb_listener" "default" {
	  count = 2
	  load_balancer_id = "${element(alicloud_slb_load_balancer.default.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "tcp"
	  bandwidth = "10"
	  health_check_type = "tcp"
	}`, EcsInstanceCommonTestCase, name)
}

func resourceEssScalingGroupRdsDependence(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}

	resource "alicloud_db_instance" "default" {
  		engine           = "MySQL"
  		engine_version   = "5.6"
  		instance_type    = "rds.mysql.s1.small"
  		instance_storage = "10"
  		vswitch_id       = "${alicloud_vswitch.default.id}"
  		instance_name    = var.name
	}`, EcsInstanceCommonTestCase, name)
}

func resourceEssScalingGroupInstance(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}

    resource "alicloud_vswitch" "tmpVs" {
  		vpc_id = "${alicloud_vpc.default.id}"
  		cidr_block = "172.16.1.0/24"
  		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  		name = "${var.name}-bar"
	}
	data "alicloud_images" "default1" {
		name_regex  = "^centos.*_64"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_instance_types" "c6" {
      //instance_type_family = "ecs.c6"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}

	resource "alicloud_instance" "default" {
		image_id = "${data.alicloud_images.default1.images.0.id}"
		instance_type = "${data.alicloud_instance_types.c6.instance_types.0.id}"
		security_groups = ["${alicloud_security_group.default1.id}"]
		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = "10"
		instance_charge_type = "PostPaid"
		system_disk_category = "cloud_efficiency"
		vswitch_id = "${alicloud_vswitch.tmpVs.id}"
		instance_name = "${var.name}"
	}
	
	resource "alicloud_eci_container_group" "default" {
	  container_group_name = "test"
	  cpu                  = 8.0
	  memory               = 16.0
	  restart_policy       = "OnFailure"
	  security_group_id    = alicloud_security_group.default1.id
	  vswitch_id           = alicloud_vswitch.tmpVs.id
	  auto_create_eip      = true
	  tags = {
		Created = "TF",
		For     = "example",
	  }
	  containers {
		image             = "registry.cn-beijing.aliyuncs.com/eci_open/nginx:alpine"
		name              = "nginx"
		working_dir       = "/tmp/nginx"
		image_pull_policy = "IfNotPresent"
		commands          = ["/bin/sh", "-c", "sleep 9999"]
		volume_mounts {
		  mount_path = "/tmp/example"
		  read_only  = false
		  name       = "empty1"
		}
		ports {
		  port     = 80
		  protocol = "TCP"
		}
		environment_vars {
		  key   = "name"
		  value = "nginx"
		}
		liveness_probe {
		  period_seconds        = "5"
		  initial_delay_seconds = "5"
		  success_threshold     = "1"
		  failure_threshold     = "3"
		  timeout_seconds       = "1"
		  exec {
			commands = ["cat /tmp/healthy"]
		  }
		}
		readiness_probe {
		  period_seconds        = "5"
		  initial_delay_seconds = "5"
		  success_threshold     = "1"
		  failure_threshold     = "3"
		  timeout_seconds       = "1"
		  exec {
			commands = ["cat /tmp/healthy"]
		  }
		}
	  }
	  init_containers {
		name              = "init-busybox"
		image             = "registry.cn-beijing.aliyuncs.com/eci_open/busybox:1.30"
		image_pull_policy = "IfNotPresent"
		commands          = ["echo"]
		args              = ["hello initcontainer"]
	  }
	  volumes {
		name = "empty1"
		type = "EmptyDirVolume"
	  }
	  volumes {
		name = "empty2"
		type = "EmptyDirVolume"
	  }
	}
	

	resource "alicloud_security_group" "default1" {
  		name   = var.name
  		vpc_id = "${alicloud_vpc.default.id}"
	}`, EcsInstanceCommonTestCase, name)
}

func resourceEssScalingGroupTemplate(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}

    resource "alicloud_vswitch" "tmpVs" {
  		vpc_id = "${alicloud_vpc.default.id}"
  		cidr_block = "172.16.1.0/24"
  		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  		name = "${var.name}-bar"
	}

	resource "alicloud_security_group" "default1" {
  		name   = var.name
  		vpc_id = "${alicloud_vpc.default.id}"
	}

	data "alicloud_images" "default3" {
		name_regex  = "^centos.*_64"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_instance_types" "c6" {
     availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}

	resource "alicloud_ecs_launch_template" "default3" {
  		launch_template_name = "tf-test3"
  		image_id             =  data.alicloud_images.default3.images.0.id
  		instance_charge_type = "PrePaid"
  		instance_type        =  data.alicloud_instance_types.c6.instance_types.0.id
  		internet_charge_type          = "PayByBandwidth"
  		internet_max_bandwidth_in     = "5"
  		internet_max_bandwidth_out    = "0"
  		io_optimized                  = "optimized"
  		network_type                  = "vpc"
  		security_enhancement_strategy = "Active"
  		spot_price_limit              = "5"
  		spot_strategy                 = "SpotWithPriceLimit"
  		security_group_id             = alicloud_security_group.default1.id
	}
	
	resource "alicloud_ecs_launch_template" "default4" {
  		launch_template_name = "tf-test4"
  		image_id             =  data.alicloud_images.default3.images.0.id
  		instance_charge_type = "PrePaid"
  		instance_type        =  data.alicloud_instance_types.c6.instance_types.0.id
  		internet_charge_type          = "PayByBandwidth"
  		internet_max_bandwidth_in     = "5"
  		internet_max_bandwidth_out    = "0"
  		io_optimized                  = "optimized"
  		network_type                  = "vpc"
  		security_enhancement_strategy = "Active"
  		spot_price_limit              = "5"
  		spot_strategy                 = "SpotWithPriceLimit"
  		security_group_id             = alicloud_security_group.default1.id
	}

	resource "alicloud_ecs_launch_template" "default1" {
  		launch_template_name = "tf-test"
  		image_id             =  data.alicloud_images.default.images.0.id
  		instance_charge_type = "PrePaid"
  		instance_type        =  data.alicloud_instance_types.default.instance_types.0.id
  		internet_charge_type          = "PayByBandwidth"
  		internet_max_bandwidth_in     = "5"
  		internet_max_bandwidth_out    = "0"
  		io_optimized                  = "optimized"
  		network_type                  = "vpc"
  		security_enhancement_strategy = "Active"
  		spot_price_limit              = "5"
  		spot_strategy                 = "SpotWithPriceLimit"
  		security_group_id             = alicloud_security_group.default1.id
	}

	resource "alicloud_ecs_launch_template" "default" {
  		launch_template_name = "tf-test1"
  		image_id             =  data.alicloud_images.default.images.0.id
  		instance_charge_type = "PrePaid"
  		instance_type        =  data.alicloud_instance_types.default.instance_types.0.id
  		internet_charge_type          = "PayByBandwidth"
  		internet_max_bandwidth_in     = "5"
  		internet_max_bandwidth_out    = "0"
  		io_optimized                  = "optimized"
  		network_type                  = "vpc"
  		security_enhancement_strategy = "Active"
  		spot_price_limit              = "5"
  		spot_strategy                 = "SpotWithPriceLimit"
  		security_group_id             = alicloud_security_group.default1.id
	}`, EcsInstanceCommonTestCase, name)
}

func resourceEssScalingGroupResourceGroupId(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}

    resource "alicloud_vswitch" "tmpVs" {
  		vpc_id = "${alicloud_vpc.default.id}"
  		cidr_block = "172.16.1.0/24"
  		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  		name = "${var.name}-bar"
	}

	resource "alicloud_security_group" "default1" {
  		name   = var.name
  		vpc_id = "${alicloud_vpc.default.id}"
	}

	resource "alicloud_ecs_launch_template" "default1" {
  		launch_template_name = "tf-test"
  		image_id             =  data.alicloud_images.default.images.0.id
  		instance_charge_type = "PrePaid"
  		instance_type        =  data.alicloud_instance_types.default.instance_types.0.id
  		internet_charge_type          = "PayByBandwidth"
  		internet_max_bandwidth_in     = "5"
  		internet_max_bandwidth_out    = "0"
  		io_optimized                  = "optimized"
  		network_type                  = "vpc"
  		security_enhancement_strategy = "Active"
  		spot_price_limit              = "5"
  		spot_strategy                 = "SpotWithPriceLimit"
  		security_group_id             = alicloud_security_group.default1.id
	}
	data "alicloud_resource_manager_resource_groups" "default" {
	}
	resource "alicloud_ecs_launch_template" "default" {
  		launch_template_name = "tf-test1"
  		image_id             =  data.alicloud_images.default.images.0.id
  		instance_charge_type = "PrePaid"
  		instance_type        =  data.alicloud_instance_types.default.instance_types.0.id
  		internet_charge_type          = "PayByBandwidth"
  		internet_max_bandwidth_in     = "5"
  		internet_max_bandwidth_out    = "0"
  		io_optimized                  = "optimized"
  		network_type                  = "vpc"
  		security_enhancement_strategy = "Active"
  		spot_price_limit              = "5"
  		spot_strategy                 = "SpotWithPriceLimit"
  		security_group_id             = alicloud_security_group.default1.id
	}`, EcsInstanceCommonTestCase, name)
}

func resourceEssScalingGroupAlbServerGroup(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}

    resource "alicloud_vswitch" "tmpVs" {
  		vpc_id = "${alicloud_vpc.default.id}"
  		cidr_block = "172.16.1.0/24"
  		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  		name = "${var.name}-bar"
	}

	resource "alicloud_security_group" "default1" {
  		name   = var.name
  		vpc_id = "${alicloud_vpc.default.id}"
	}

	data "alicloud_images" "default3" {
		name_regex  = "^centos.*_64"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_instance_types" "c6" {
       availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}

	resource "alicloud_ecs_launch_template" "default3" {
  		launch_template_name = "tf-test3"
  		image_id             =  data.alicloud_images.default3.images.0.id
  		instance_charge_type = "PrePaid"
  		instance_type        =  data.alicloud_instance_types.c6.instance_types.0.id
  		internet_charge_type          = "PayByBandwidth"
  		internet_max_bandwidth_in     = "5"
  		internet_max_bandwidth_out    = "0"
  		io_optimized                  = "optimized"
  		network_type                  = "vpc"
  		security_enhancement_strategy = "Active"
  		spot_price_limit              = "5"
  		spot_strategy                 = "SpotWithPriceLimit"
  		security_group_id             = alicloud_security_group.default1.id
	}

	resource "alicloud_ecs_launch_template" "default4" {
  		launch_template_name = "tf-test4"
  		image_id             =  data.alicloud_images.default3.images.0.id
  		instance_charge_type = "PrePaid"
  		instance_type        =  data.alicloud_instance_types.c6.instance_types.0.id
  		internet_charge_type          = "PayByBandwidth"
  		internet_max_bandwidth_in     = "5"
  		internet_max_bandwidth_out    = "0"
  		io_optimized                  = "optimized"
  		network_type                  = "vpc"
  		security_enhancement_strategy = "Active"
  		spot_price_limit              = "5"
  		spot_strategy                 = "SpotWithPriceLimit"
  		security_group_id             = alicloud_security_group.default1.id
	}

	resource "alicloud_ecs_launch_template" "default1" {
  		launch_template_name = "tf-test"
  		image_id             =  data.alicloud_images.default.images.0.id
  		instance_charge_type = "PrePaid"
  		instance_type        =  data.alicloud_instance_types.default.instance_types.0.id
  		internet_charge_type          = "PayByBandwidth"
  		internet_max_bandwidth_in     = "5"
  		internet_max_bandwidth_out    = "0"
  		io_optimized                  = "optimized"
  		network_type                  = "vpc"
  		security_enhancement_strategy = "Active"
  		spot_price_limit              = "5"
  		spot_strategy                 = "SpotWithPriceLimit"
  		security_group_id             = alicloud_security_group.default1.id
	}
	
	resource "alicloud_alb_server_group" "default" {
        count = 2
		server_group_name = "${var.name}"
		vpc_id = "${alicloud_vpc.default.id}"
		health_check_config {
		  health_check_enabled = "false"
		}
		sticky_session_config {
		  sticky_session_enabled = true
		  cookie                 = "tf-testAcc"
		  sticky_session_type    = "Server"
	  }
	}

	resource "alicloud_ecs_launch_template" "default" {
  		launch_template_name = "tf-test1"
  		image_id             =  data.alicloud_images.default.images.0.id
  		instance_charge_type = "PrePaid"
  		instance_type        =  data.alicloud_instance_types.default.instance_types.0.id
  		internet_charge_type          = "PayByBandwidth"
  		internet_max_bandwidth_in     = "5"
  		internet_max_bandwidth_out    = "0"
  		io_optimized                  = "optimized"
  		network_type                  = "vpc"
  		security_enhancement_strategy = "Active"
  		spot_price_limit              = "5"
  		spot_strategy                 = "SpotWithPriceLimit"
  		security_group_id             = alicloud_security_group.default1.id
	}`, EcsInstanceCommonTestCase, name)
}
