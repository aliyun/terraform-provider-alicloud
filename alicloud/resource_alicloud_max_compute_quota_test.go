package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test MaxCompute Quota. >>> Resource test cases, automatically generated.
// Case Quota发布terraform-包年包月标准计算 10404
func TestAccAliCloudMaxComputeQuota_basic10404(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_quota.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeQuotaMap10404)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeQuota")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmaxcompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeQuotaBasicDependence10404)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sub_quota_info_list": []map[string]interface{}{
						{
							"parameter": []map[string]interface{}{
								{
									"min_cu":              "60",
									"max_cu":              "60",
									"enable_priority":     "false",
									"force_reserved_min":  "false",
									"scheduler_type":      "Fifo",
									"single_job_cu_limit": "10",
								},
							},
							"nick_name": "os_${var.part_nick_name}",
							"type":      "FUXI_OFFLINE",
						},
						{
							"parameter": []map[string]interface{}{
								{
									"min_cu":             "10",
									"max_cu":             "10",
									"scheduler_type":     "Fifo",
									"enable_priority":    "true",
									"force_reserved_min": "true",
								},
							},
							"nick_name": "${var.sub_quota_nickname_1}",
							"type":      "FUXI_OFFLINE",
						},
						{
							"parameter": []map[string]interface{}{
								{
									"min_cu":             "10",
									"max_cu":             "10",
									"force_reserved_min": "false",
									"scheduler_type":     "Fifo",
									"enable_priority":    "false",
								},
							},
							"nick_name": "${var.sub_quota_nickname_2}",
							"type":      "FUXI_OFFLINE",
						},
					},
					"payment_type":   "Subscription",
					"part_nick_name": "${var.part_nick_name}",
					"commodity_data": "{\\\"CU\\\":80,\\\"ord_time\\\":\\\"1:Month\\\",\\\"autoRenew\\\":false} ",
					"commodity_code": "odpsplus",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sub_quota_info_list.#": "3",
						"payment_type":          "Subscription",
						"part_nick_name":        CHECKSET,
						"commodity_data":        CHECKSET,
						"commodity_code":        "odpsplus",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sub_quota_info_list": []map[string]interface{}{
						{
							"parameter": []map[string]interface{}{
								{
									"min_cu":              "10",
									"max_cu":              "10",
									"enable_priority":     "true",
									"force_reserved_min":  "true",
									"single_job_cu_limit": "16",
									"scheduler_type":      "Fair",
								},
							},
							"nick_name": "os_${var.part_nick_name}",
							"type":      "FUXI_OFFLINE",
						},
						{
							"parameter": []map[string]interface{}{
								{
									"min_cu":             "10",
									"max_cu":             "10",
									"scheduler_type":     "Fair",
									"enable_priority":    "false",
									"force_reserved_min": "false",
								},
							},
							"nick_name": "${var.sub_quota_nickname_1}",
							"type":      "FUXI_OFFLINE",
						},
						{
							"parameter": []map[string]interface{}{
								{
									"min_cu":             "10",
									"max_cu":             "10",
									"scheduler_type":     "Fair",
									"enable_priority":    "true",
									"force_reserved_min": "true",
								},
							},
							"nick_name": "${var.sub_quota_nickname_2}",
							"type":      "FUXI_OFFLINE",
						},
						{
							"parameter": []map[string]interface{}{
								{
									"min_cu":             "50",
									"max_cu":             "50",
									"force_reserved_min": "true",
								},
							},
							"nick_name": "${var.sub_quota_nickname_3}",
							"type":      "FUXI_ONLINE",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sub_quota_info_list.#": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sub_quota_info_list": []map[string]interface{}{
						{
							"parameter": []map[string]interface{}{
								{
									"min_cu":             "80",
									"max_cu":             "80",
									"enable_priority":    "true",
									"force_reserved_min": "true",
								},
							},
							"nick_name": "os_${var.part_nick_name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sub_quota_info_list.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"commodity_code", "commodity_data", "part_nick_name"},
			},
		},
	})
}

var AlicloudMaxComputeQuotaMap10404 = map[string]string{}

func AlicloudMaxComputeQuotaBasicDependence10404(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "part_nick_name" {
  default = "TFTest17216"
}

variable "sub_quota_nickname_3" {
  default = "sub398816"
}

variable "sub_quota_nickname_1" {
  default = "sub169716"
}

variable "sub_quota_nickname_2" {
  default = "sub223116"
}


`, name)
}

// Case 后付费测试用例 6979
func TestAccAliCloudMaxComputeQuota_basic6979(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_quota.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeQuotaMap6979)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeQuota")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmaxcompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeQuotaBasicDependence6979)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":   "PayAsYouGo",
					"commodity_code": "odps",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"commodity_code", "commodity_data", "part_nick_name"},
			},
		},
	})
}

var AlicloudMaxComputeQuotaMap6979 = map[string]string{}

func AlicloudMaxComputeQuotaBasicDependence6979(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test MaxCompute Quota. <<< Resource test cases, automatically generated.
