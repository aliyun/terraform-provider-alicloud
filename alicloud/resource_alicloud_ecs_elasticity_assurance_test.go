package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudEcsElasticityAssurance_basic1716(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_elasticity_assurance.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsElasticityAssuranceMap1716)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsElasticityAssurance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sEcsElasticityAssurance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsElasticityAssuranceBasicDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_amount":                     "1",
					"description":                         "before",
					"zone_ids":                            []string{"${data.alicloud_zones.default.zones[0].id}"},
					"private_pool_options_name":           "test_before",
					"period":                              "1",
					"private_pool_options_match_criteria": "Open",
					"instance_type":                       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"period_unit":                         "Month",
					"assurance_times":                     "Unlimited",
					"start_time":                          time.Now().Add(1 * time.Hour).Format("2006-01-02T15:04:05Z"),
					"resource_group_id":                   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "tfTestAcc0",
						"For":     "Tftestacc 0",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_amount":                     "1",
						"description":                         "before",
						"zone_ids.#":                          "1",
						"private_pool_options_name":           "test_before",
						"period":                              "1",
						"private_pool_options_match_criteria": "Open",
						"instance_type.#":                     "1",
						"period_unit":                         "Month",
						"assurance_times":                     "Unlimited",
						"start_time":                          CHECKSET,
						"resource_group_id":                   CHECKSET,
						"tags.%":                              "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

func TestAccAlicloudEcsElasticityAssurance_basic1717(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_elasticity_assurance.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsElasticityAssuranceMap1716)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsElasticityAssurance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sEcsElasticityAssurance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsElasticityAssuranceBasicDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_ids":        []string{"${data.alicloud_zones.default.zones[0].id}"},
					"instance_type":   []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"instance_amount": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_ids.#":      "1",
						"instance_type.#": "1",
						"instance_amount": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc1",
						"For":     "Tftestacc 1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_pool_options_name": "test_after_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_pool_options_name": "test_after_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "after_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "after_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_amount":                     "1",
					"description":                         "after",
					"private_pool_options_name":           "test_after",
					"period":                              "2",
					"private_pool_options_match_criteria": "Open",
					"instance_type":                       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_amount":                     "1",
						"description":                         "after",
						"private_pool_options_name":           "test_after",
						"period":                              "2",
						"private_pool_options_match_criteria": "Open",
						"instance_type.#":                     "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

var AlicloudEcsElasticityAssuranceMap1716 = map[string]string{}

func AlicloudEcsElasticityAssuranceBasicDependence1716(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}

data "alicloud_instance_types" "default" {
	instance_type_family = "ecs.c6"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}

`, name)
}
