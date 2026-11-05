package maxcompute_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/maxcompute"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test MaxCompute QuotaSchedule. >>> Resource test cases, automatically generated.
// Case QuotaSchedule_terraform测试 9281
func TestAccAliCloudMaxComputeQuotaSchedule_basic9281(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_quota_schedule.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMaxComputeQuotaScheduleMap9281)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return maxcompute.NewMaxComputeServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMaxComputeQuotaSchedule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmcqs%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMaxComputeQuotaScheduleBasicDependence9281)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
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
					"nickname": "${alicloud_max_compute_quota.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule_list.#": "3",
						"timezone":        "UTC+8",
						"nickname":        acctest.CHECKSET,
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

resource "alicloud_max_compute_quota" "default" {
  payment_type   = "Subscription"
  part_nick_name = var.name
  commodity_code = "odpsplus"
  commodity_data = "{\"CU\":50,\"ord_time\":\"1:Month\",\"autoRenew\":false}"
  sub_quota_info_list {
    nick_name = "os_${var.name}"
    type      = "FUXI_OFFLINE"
    parameter {
      min_cu = "30"
      max_cu = "50"
    }
  }
  sub_quota_info_list {
    nick_name = "sub_${var.name}"
    type      = "FUXI_OFFLINE"
    parameter {
      min_cu = "20"
      max_cu = "50"
    }
  }
}

resource "alicloud_max_compute_quota_plan" "default" {
  quota {
    parameter {
      elastic_reserved_cu = 0
    }
    sub_quota_info_list {
      nick_name = "sub_${var.name}"
      parameter {
        min_cu              = "0"
        max_cu              = "20"
        elastic_reserved_cu = "0"
      }
    }
    sub_quota_info_list {
      nick_name = "os_${var.name}"
      parameter {
        min_cu              = "50"
        max_cu              = "50"
        elastic_reserved_cu = "0"
      }

    }
  }

  plan_name = "quota_plan1"
  nickname  = alicloud_max_compute_quota.default.id
}

resource "alicloud_max_compute_quota_plan" "default2" {
  quota {
    parameter {
      elastic_reserved_cu = 0
    }
    sub_quota_info_list {
      nick_name = "sub_${var.name}"
      parameter {
        min_cu              = "0"
        max_cu              = "30"
        elastic_reserved_cu = "0"
      }
    }
    sub_quota_info_list {
      nick_name = "os_${var.name}"
      parameter {
        min_cu              = "50"
        max_cu              = "50"
        elastic_reserved_cu = "0"
      }

    }
  }

  plan_name = "quota_plan2"
  nickname  = alicloud_max_compute_quota.default.id
}

resource "alicloud_max_compute_quota_plan" "default3" {
  quota {
    parameter {
      elastic_reserved_cu = 0
    }
    sub_quota_info_list {
      nick_name = "sub_${var.name}"
      parameter {
        min_cu              = "40"
        max_cu              = "40"
        elastic_reserved_cu = "0"
      }
    }
    sub_quota_info_list {
      nick_name = "os_${var.name}"
      parameter {
        min_cu              = "10"
        max_cu              = "10"
        elastic_reserved_cu = "0"
      }

    }
  }

  plan_name = "quota_plan3"
  nickname  = alicloud_max_compute_quota.default.id
}

`, name)
}

// Test MaxCompute QuotaSchedule. <<< Resource test cases, automatically generated.
