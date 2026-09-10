package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test MaxCompute QuotaSchedule. >>> Resource test cases, automatically generated.
// Case QuotaSchedule_terraform测试 9281
func TestAccAliCloudMaxComputeQuotaSchedule_basic9281(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_quota_schedule.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeQuotaScheduleMap9281)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeQuotaSchedule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smaxcomputequotaschedule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeQuotaScheduleBasicDependence9281)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_list": []map[string]interface{}{
						{
							"plan": "${alicloud_max_compute_quota_plan.default.plan_name}",
							"condition": []map[string]interface{}{
								{
									"at": "00:00",
								},
							},
							"type": "daily",
						},
						{
							"plan": "${alicloud_max_compute_quota_plan.default2.plan_name}",
							"condition": []map[string]interface{}{
								{
									"at": "01:00",
								},
							},
							"type": "daily",
						},
						{
							"plan": "${alicloud_max_compute_quota_plan.default3.plan_name}",
							"condition": []map[string]interface{}{
								{
									"at": "02:00",
								},
							},
							"type": "daily",
						},
					},
					"timezone": "UTC+8",
					"nickname": "${var.quota_nick_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule_list.#": "3",
						"timezone":        "UTC+8",
						"nickname":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_list": []map[string]interface{}{
						{
							"plan": "${alicloud_max_compute_quota_plan.default3.plan_name}",
							"condition": []map[string]interface{}{
								{
									"at": "00:00",
								},
							},
							"type": "daily",
						},
						{
							"plan": "${alicloud_max_compute_quota_plan.default2.plan_name}",
							"condition": []map[string]interface{}{
								{
									"at": "08:00",
								},
							},
							"type": "daily",
						},
						{
							"condition": []map[string]interface{}{
								{
									"at": "10:00",
								},
							},
							"plan": "${alicloud_max_compute_quota_plan.default.plan_name}",
							"type": "daily",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_list": []map[string]interface{}{
						{
							"plan": "${alicloud_max_compute_quota_plan.default2.plan_name}",
							"condition": []map[string]interface{}{
								{
									"at": "00:00",
								},
							},
							"type": "daily",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_list": []map[string]interface{}{
						{
							"plan": "Default",
							"condition": []map[string]interface{}{
								{
									"at": "00:00",
								},
							},
							"type": "daily",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule_list.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudMaxComputeQuotaScheduleMap9281 = map[string]string{}

func AlicloudMaxComputeQuotaScheduleBasicDependence9281(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


variable "elastic_reserved_cu" {
  default = "0"
}

variable "quota_nick_name" {
  default = "os_terrform_p"
}

resource "alicloud_max_compute_quota_plan" "default" {
  quota {
    parameter {
      elastic_reserved_cu = 50
    }
    sub_quota_info_list {
      nick_name = "sub_quota"
      parameter {
        min_cu              = "0"
        max_cu              = "20"
        elastic_reserved_cu = "30"
      }
    }
    sub_quota_info_list {
      nick_name = "os_terrform"
      parameter {
        min_cu              = "50"
        max_cu              = "50"
        elastic_reserved_cu = "20"
      }

    }
  }

  plan_name = "quota_plan1"
  nickname  = "os_terrform_p"
}

resource "alicloud_max_compute_quota_plan" "default2" {
  quota {
    parameter {
      elastic_reserved_cu = 50
    }
    sub_quota_info_list {
      nick_name = "sub_quota"
      parameter {
        min_cu              = "0"
        max_cu              = "20"
        elastic_reserved_cu = "20"
      }
    }
    sub_quota_info_list {
      nick_name = "os_terrform"
      parameter {
        min_cu              = "50"
        max_cu              = "50"
        elastic_reserved_cu = "30"
      }

    }
  }

  plan_name = "quota_plan2"
  nickname  = "os_terrform_p"
}

resource "alicloud_max_compute_quota_plan" "default3" {
  quota {
    parameter {
      elastic_reserved_cu = 50
    }
    sub_quota_info_list {
      nick_name = "sub_quota"
      parameter {
        min_cu              = "40"
        max_cu              = "40"
        elastic_reserved_cu = "40"
      }
    }
    sub_quota_info_list {
      nick_name = "os_terrform"
      parameter {
        min_cu              = "10"
        max_cu              = "10"
        elastic_reserved_cu = "10"
      }

    }
  }

  plan_name = "quota_plan3"
  nickname  = "os_terrform_p"
}

`, name)
}

// Test MaxCompute QuotaSchedule. <<< Resource test cases, automatically generated.
