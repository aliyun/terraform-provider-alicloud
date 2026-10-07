// Package alicloud. This file is hand-written from the Cms 2024-03-30 OpenAPI definition.
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsOncallSchedule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_oncall_schedule.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsOncallScheduleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsOncallSchedule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsOncallScheduleBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"oncall_schedule_name": name,
					"rotations": []interface{}{
						map[string]interface{}{
							"active_days":                []interface{}{1, 2, 3, 4, 5},
							"contacts":                   []interface{}{name},
							"rotation_end_time":          "2026-01-01 23:59:59",
							"rotation_name":              name,
							"rotation_start_time":        "2026-01-01 00:00:00",
							"shift_length":               1,
							"shift_recurrence_frequency": "Daily",
							"start_date":                 "2026-01-01",
							"time_zone":                  "Asia/Shanghai",
						},
					},
					"shift_robot_id": name,
					"source":         name,
					"substitudes":    []interface{}{name},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oncall_schedule_name": name,
						"shift_robot_id":       name,
						"source":               name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oncall_schedule_name": name,
					"rotations": []interface{}{
						map[string]interface{}{
							"active_days":                []interface{}{6, 7},
							"contacts":                   []interface{}{fmt.Sprintf("%s-sub", name)},
							"rotation_end_time":          "2026-01-02 23:59:59",
							"rotation_name":              fmt.Sprintf("%s-updated", name),
							"rotation_start_time":        "2026-01-02 00:00:00",
							"shift_length":               2,
							"shift_recurrence_frequency": "Weekly",
							"start_date":                 "2026-01-02",
							"time_zone":                  "UTC",
						},
					},
					"shift_robot_id": fmt.Sprintf("%s-updated", name),
					"source":         fmt.Sprintf("%s-updated", name),
					"substitudes":    []interface{}{fmt.Sprintf("%s-sub", name)},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oncall_schedule_name": name,
						"shift_robot_id":       fmt.Sprintf("%s-updated", name),
						"source":               fmt.Sprintf("%s-updated", name),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oncall_schedule_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oncall_schedule_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rotations", "substitudes", "shift_robot_id", "source"},
			},
		},
	})
}

func TestAccAliCloudCmsOncallSchedule_basic_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_oncall_schedule.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsOncallScheduleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsOncallSchedule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsOncallScheduleBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"oncall_schedule_name": name,
					"rotations": []interface{}{
						map[string]interface{}{
							"active_days":                []interface{}{1, 2, 3},
							"contacts":                   []interface{}{name},
							"rotation_end_time":          "2026-01-01 23:59:59",
							"rotation_name":              name,
							"rotation_start_time":        "2026-01-01 00:00:00",
							"shift_length":               1,
							"shift_recurrence_frequency": "Daily",
							"start_date":                 "2026-01-01",
							"time_zone":                  "Asia/Shanghai",
						},
					},
					"shift_robot_id": name,
					"source":         name,
					"substitudes":    []interface{}{name},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oncall_schedule_name": name,
						"shift_robot_id":       name,
						"source":               name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rotations", "substitudes", "shift_robot_id", "source"},
			},
		},
	})
}

var AliCloudCmsOncallScheduleMap = map[string]string{
	"oncall_schedule_id": CHECKSET,
}

func AliCloudCmsOncallScheduleBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}
