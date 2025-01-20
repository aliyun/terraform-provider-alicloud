package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test MaxCompute TunnelQuotaTimer. >>> Resource test cases, automatically generated.
// Case TunnelQuotaTimer_terraform测试 9976
func TestAccAliCloudMaxComputeTunnelQuotaTimer_basic9976(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_tunnel_quota_timer.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeTunnelQuotaTimerMap9976)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeTunnelQuotaTimer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smaxcomputetunnelquotatimer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeTunnelQuotaTimerBasicDependence9976)
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
					"quota_timer": []map[string]interface{}{
						{
							"begin_time": "00:00",
							"end_time":   "01:00",
							"tunnel_quota_parameter": []map[string]interface{}{
								{
									"slot_num":                  "50",
									"elastic_reserved_slot_num": "50",
								},
							},
						},
						{
							"begin_time": "01:00",
							"end_time":   "02:00",
							"tunnel_quota_parameter": []map[string]interface{}{
								{
									"slot_num":                  "50",
									"elastic_reserved_slot_num": "50",
								},
							},
						},
						{
							"begin_time": "02:00",
							"end_time":   "24:00",
							"tunnel_quota_parameter": []map[string]interface{}{
								{
									"slot_num":                  "50",
									"elastic_reserved_slot_num": "50",
								},
							},
						},
					},
					"nickname":  "ot_terraform_p",
					"time_zone": "Asia/Shanghai",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_timer.#": "3",
						"nickname":      "ot_terraform_p",
						"time_zone":     "Asia/Shanghai",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_timer": []map[string]interface{}{
						{
							"begin_time": "00:00",
							"end_time":   "20:00",
							"tunnel_quota_parameter": []map[string]interface{}{
								{
									"slot_num":                  "50",
									"elastic_reserved_slot_num": "40",
								},
							},
						},
						{
							"begin_time": "20:00",
							"end_time":   "22:00",
							"tunnel_quota_parameter": []map[string]interface{}{
								{
									"slot_num":                  "50",
									"elastic_reserved_slot_num": "40",
								},
							},
						},
						{
							"begin_time": "22:00",
							"end_time":   "23:00",
							"tunnel_quota_parameter": []map[string]interface{}{
								{
									"slot_num":                  "50",
									"elastic_reserved_slot_num": "40",
								},
							},
						},
						{
							"begin_time": "23:00",
							"end_time":   "24:00",
							"tunnel_quota_parameter": []map[string]interface{}{
								{
									"slot_num":                  "50",
									"elastic_reserved_slot_num": "40",
								},
							},
						},
					},
					"time_zone": "Asia/Tokyo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_timer.#": "4",
						"time_zone":     "Asia/Tokyo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_timer": []map[string]interface{}{
						{
							"begin_time": "00:00",
							"end_time":   "20:00",
							"tunnel_quota_parameter": []map[string]interface{}{
								{
									"slot_num":                  "50",
									"elastic_reserved_slot_num": "50",
								},
							},
						},
						{
							"begin_time": "20:00",
							"end_time":   "23:00",
							"tunnel_quota_parameter": []map[string]interface{}{
								{
									"slot_num":                  "50",
									"elastic_reserved_slot_num": "50",
								},
							},
						},
						{
							"begin_time": "23:00",
							"end_time":   "24:00",
							"tunnel_quota_parameter": []map[string]interface{}{
								{
									"slot_num":                  "50",
									"elastic_reserved_slot_num": "50",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_timer.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_timer": []map[string]interface{}{
						{
							"begin_time": "00:00",
							"end_time":   "24:00",
							"tunnel_quota_parameter": []map[string]interface{}{
								{
									"slot_num":                  "50",
									"elastic_reserved_slot_num": "50",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_timer.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_timer": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_timer.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"time_zone"},
			},
		},
	})
}

var AlicloudMaxComputeTunnelQuotaTimerMap9976 = map[string]string{}

func AlicloudMaxComputeTunnelQuotaTimerBasicDependence9976(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test MaxCompute TunnelQuotaTimer. <<< Resource test cases, automatically generated.
