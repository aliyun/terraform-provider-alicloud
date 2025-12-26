package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudThreatDetectionCheckConfig_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_check_config.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionCheckConfigMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionCheckConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionBaselineStrategy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionCheckConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"config_standard_ids": []map[string]interface{}{
						{
							"remove_ids": []int{3},
							"add_ids":    []int{1, 2},
						},
					},
					"config_requirement_ids": []map[string]interface{}{
						{
							"remove_ids": []int{4},
							"add_ids":    []int{1, 2, 3},
						},
					},
					"system_config":     false,
					"end_time":          18,
					"enable_auto_check": true,
					"enable_add_check":  true,
					"start_time":        12,
					"cycle_days":        []int{7, 1, 2},
					"removed_check": []map[string]interface{}{
						{
							"check_id":   370,
							"section_id": 515,
						},
					},
					"added_check": []map[string]interface{}{
						{
							"check_id":   5,
							"section_id": 384,
						},
					},
					"vendors":   []string{"ALIYUN"},
					"configure": "not",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_standard_ids.#":              "1",
						"config_standard_ids.0.remove_ids.#": "1",
						"config_standard_ids.0.remove_ids.0": "3",
						"config_standard_ids.0.add_ids.#":    "2",
						"config_standard_ids.0.add_ids.0":    "1",
						"config_standard_ids.0.add_ids.1":    "2",

						"config_requirement_ids.#":              "1",
						"config_requirement_ids.0.remove_ids.#": "1",
						"config_requirement_ids.0.remove_ids.0": "4",
						"config_requirement_ids.0.add_ids.#":    "3",
						"config_requirement_ids.0.add_ids.0":    "1",
						"config_requirement_ids.0.add_ids.1":    "2",
						"config_requirement_ids.0.add_ids.2":    "3",
						"system_config":                         "false",
						"end_time":                              "18",
						"enable_auto_check":                     "true",
						"enable_add_check":                      "true",
						"start_time":                            "12",
						"cycle_days.#":                          "3",
						"removed_check.#":                       "1",
						"removed_check.0.check_id":              "370",
						"removed_check.0.section_id":            "515",
						"added_check.#":                         "1",
						"added_check.0.check_id":                "5",
						"added_check.0.section_id":              "384",
						"vendors.#":                             "1",
						"vendors.0":                             "ALIYUN",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_standard_ids": []map[string]interface{}{
						{
							"remove_ids": []int{4},
							"add_ids":    []int{1, 3},
						},
					},
					"config_requirement_ids": []map[string]interface{}{
						{
							"remove_ids": []int{4},
							"add_ids":    []int{1, 2, 3},
						},
					},
					"system_config":     false,
					"end_time":          12,
					"enable_auto_check": false,
					"enable_add_check":  false,
					"start_time":        6,
					"cycle_days":        []int{4},
					"removed_check": []map[string]interface{}{
						{
							"check_id":   1607,
							"section_id": 1000000000026,
						},
					},
					"added_check": []map[string]interface{}{
						{
							"check_id":   1607,
							"section_id": 1000000000026,
						},
					},
					"vendors":   []string{"ALIYUN"},
					"configure": "all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_standard_ids.#":              "1",
						"config_standard_ids.0.remove_ids.#": "1",
						"config_standard_ids.0.remove_ids.0": "4",
						"config_standard_ids.0.add_ids.#":    "2",
						"config_standard_ids.0.add_ids.0":    "1",
						"config_standard_ids.0.add_ids.1":    "3",

						"config_requirement_ids.#":              "1",
						"config_requirement_ids.0.remove_ids.#": "1",
						"config_requirement_ids.0.remove_ids.0": "4",
						"config_requirement_ids.0.add_ids.#":    "3",
						"config_requirement_ids.0.add_ids.0":    "1",
						"config_requirement_ids.0.add_ids.1":    "2",
						"config_requirement_ids.0.add_ids.2":    "3",
						"system_config":                         "false",
						"end_time":                              "12",
						"enable_auto_check":                     "false",
						"enable_add_check":                      "false",
						"start_time":                            "6",
						"cycle_days.#":                          "1",
						"removed_check.#":                       "1",
						"removed_check.0.check_id":              "1607",
						"removed_check.0.section_id":            "1000000000026",
						"added_check.#":                         "1",
						"added_check.0.check_id":                "1607",
						"added_check.0.section_id":              "1000000000026",
						"vendors.#":                             "1",
						"vendors.0":                             "ALIYUN",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vendors", "system_config", "removed_check", "added_check", "config_standard_ids", "config_requirement_ids", "configure"},
			},
		},
	})
}

var AlicloudThreatDetectionCheckConfigMap = map[string]string{}

func AlicloudThreatDetectionCheckConfigDependence(name string) string {
	return fmt.Sprintf(` 
variable "name" {
default = "%s"
}
`, name)
}
