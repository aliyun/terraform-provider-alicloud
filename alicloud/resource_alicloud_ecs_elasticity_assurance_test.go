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
func TestAccAliCloudEcsElasticityAssurance_basic1716(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_elasticity_assurance.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsElasticityAssuranceMap1716)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsElasticityAssurance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sEcsElasticityAssurance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsElasticityAssuranceBasicDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_amount":                     "1",
					"description":                         "before",
					"zone_ids":                            []string{"${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"},
					"private_pool_options_name":           name,
					"period":                              "1",
					"private_pool_options_match_criteria": "Open",
					"instance_type":                       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"period_unit":                         "Month",
					"auto_renew":                          "true",
					"auto_renew_period":                   "2",
					"auto_renew_period_unit":              "Year",
					"assurance_times":                     "Unlimited",
					"start_time":                          time.Now().Add(1 * time.Hour).Format("2006-01-02T15:04:05Z"),
					"resource_group_id":                   "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_amount":                     "1",
						"description":                         "before",
						"zone_ids.#":                          "1",
						"private_pool_options_name":           name,
						"period":                              "1",
						"private_pool_options_match_criteria": "Open",
						"instance_type.#":                     "1",
						"period_unit":                         "Month",
						"auto_renew":                          "true",
						"auto_renew_period":                   "2",
						"auto_renew_period_unit":              "Year",
						"assurance_times":                     "Unlimited",
						"start_time":                          CHECKSET,
						"resource_group_id":                   CHECKSET,
						"tags.%":                              "2",
						"tags.Created":                        "TF",
						"tags.For":                            "Test",
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

func TestAccAliCloudEcsElasticityAssurance_basic1717(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_elasticity_assurance.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsElasticityAssuranceMap1716)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsElasticityAssurance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sEcsElasticityAssurance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsElasticityAssuranceBasicDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_ids":          []string{"${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"},
					"instance_type":     []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"instance_amount":   "1",
					"auto_renew":        "true",
					"auto_renew_period": "1",
					"period_unit":       "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_ids.#":        "1",
						"instance_type.#":   "1",
						"instance_amount":   "1",
						"auto_renew":        "true",
						"auto_renew_period": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_pool_options_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_pool_options_name": name,
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
					"instance_amount": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_amount": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":             "true",
					"auto_renew_period":      "2",
					"auto_renew_period_unit": "Year",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":             "true",
						"auto_renew_period":      "2",
						"auto_renew_period_unit": "Year",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AliCloudEcsElasticityAssuranceMap1716 = map[string]string{}

func AliCloudEcsElasticityAssuranceBasicDependence1716(name string) string {
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
`, name)
}
