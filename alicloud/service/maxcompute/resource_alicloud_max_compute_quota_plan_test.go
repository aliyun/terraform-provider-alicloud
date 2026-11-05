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

// Test MaxCompute QuotaPlan. >>> Resource test cases, automatically generated.
// Case QuotaPlan_terraform测试 9566
func TestAccAliCloudMaxComputeQuotaPlan_basic9566(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_quota_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMaxComputeQuotaPlanMap9566)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return maxcompute.NewMaxComputeServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMaxComputeQuotaPlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmcqp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMaxComputeQuotaPlanBasicDependence9566)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota": []map[string]interface{}{
						{
							"parameter": []map[string]interface{}{
								{
									"elastic_reserved_cu": "${var.elastic_reserved_cu}",
								},
							},
							"sub_quota_info_list": []map[string]interface{}{
								{
									"nick_name": "sub_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "0",
											"max_cu":              "20",
											"elastic_reserved_cu": "${var.elastic_reserved_cu}",
										},
									},
								},
								{
									"nick_name": "os_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "50",
											"max_cu":              "50",
											"elastic_reserved_cu": "0",
										},
									},
								},
							},
						},
					},
					"plan_name": "quota_plan",
					"nickname":  "${alicloud_max_compute_quota.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name": acctest.CHECKSET,
						"nickname":  acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota": []map[string]interface{}{
						{
							"parameter": []map[string]interface{}{
								{
									"elastic_reserved_cu": "${var.update_elastic_reserved_cu}",
								},
							},
							"sub_quota_info_list": []map[string]interface{}{
								{
									"nick_name": "sub_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "30",
											"max_cu":              "40",
											"elastic_reserved_cu": "0",
										},
									},
								},
								{
									"nick_name": "os_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "20",
											"max_cu":              "40",
											"elastic_reserved_cu": "${var.update_elastic_reserved_cu}",
										},
									},
								},
							},
						},
					},
					"is_effective": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_effective": "false",
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

var AlicloudMaxComputeQuotaPlanMap9566 = map[string]string{}

func AlicloudMaxComputeQuotaPlanBasicDependence9566(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "elastic_reserved_cu" {
  default = "0"
}

variable "update_elastic_reserved_cu" {
  default = "0"
}

variable "plan_name" {
  default = "TFPlan1737081504"
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


`, name)
}

// Case QuotaPlanApply_terraform测试(update) 9854
func TestAccAliCloudMaxComputeQuotaPlan_basic9854(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_quota_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMaxComputeQuotaPlanMap9854)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return maxcompute.NewMaxComputeServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMaxComputeQuotaPlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmcqp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMaxComputeQuotaPlanBasicDependence9854)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota": []map[string]interface{}{
						{
							"parameter": []map[string]interface{}{
								{
									"elastic_reserved_cu": "${var.elastic_reserved_cu}",
								},
							},
							"sub_quota_info_list": []map[string]interface{}{
								{
									"nick_name": "sub_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "0",
											"max_cu":              "20",
											"elastic_reserved_cu": "${var.elastic_reserved_cu}",
										},
									},
								},
								{
									"nick_name": "os_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "50",
											"max_cu":              "50",
											"elastic_reserved_cu": "0",
										},
									},
								},
							},
						},
					},
					"plan_name":    "${var.plan_name}",
					"is_effective": "false",
					"nickname":     "${alicloud_max_compute_quota.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name":    acctest.CHECKSET,
						"is_effective": "false",
						"nickname":     acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota": []map[string]interface{}{
						{
							"parameter": []map[string]interface{}{
								{
									"elastic_reserved_cu": "${var.update_elastic_reserved_cu}",
								},
							},
							"sub_quota_info_list": []map[string]interface{}{
								{
									"nick_name": "os_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "30",
											"max_cu":              "50",
											"elastic_reserved_cu": "0",
										},
									},
								},
								{
									"nick_name": "sub_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "20",
											"max_cu":              "50",
											"elastic_reserved_cu": "${var.update_elastic_reserved_cu}",
										},
									},
								},
							},
						},
					},
					"is_effective": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_effective": "false",
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

var AlicloudMaxComputeQuotaPlanMap9854 = map[string]string{}

func AlicloudMaxComputeQuotaPlanBasicDependence9854(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "elastic_reserved_cu" {
  default = 0
}

variable "update_elastic_reserved_cu" {
  # 50 (not 0) so the applied plan differs from the built-in Default plan
  # (QUOTA_PLAN_DUPLICATE) without drifting quota attributes terraform tracks;
  # the service requires ElasticReservedCU to be a multiple of 50
  default = 50
}

variable "plan_name" {
  default = "TFPlan1737081505"
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


`, name)
}

// Case QuotaPlanApply_terraform测试(create) 9853
func TestAccAliCloudMaxComputeQuotaPlan_basic9853(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_quota_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMaxComputeQuotaPlanMap9853)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return maxcompute.NewMaxComputeServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMaxComputeQuotaPlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmcqp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMaxComputeQuotaPlanBasicDependence9853)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota": []map[string]interface{}{
						{
							"parameter": []map[string]interface{}{
								{
									"elastic_reserved_cu": "${var.elastic_reserved_cu}",
								},
							},
							"sub_quota_info_list": []map[string]interface{}{
								{
									"nick_name": "sub_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "20",
											"max_cu":              "50",
											"elastic_reserved_cu": "${var.elastic_reserved_cu}",
										},
									},
								},
								{
									"nick_name": "os_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "30",
											"max_cu":              "50",
											"elastic_reserved_cu": "0",
										},
									},
								},
							},
						},
					},
					"plan_name":    "${var.plan_name}",
					"is_effective": "false",
					"nickname":     "${alicloud_max_compute_quota.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name":    acctest.CHECKSET,
						"is_effective": "false",
						"nickname":     acctest.CHECKSET,
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

var AlicloudMaxComputeQuotaPlanMap9853 = map[string]string{}

func AlicloudMaxComputeQuotaPlanBasicDependence9853(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "elastic_reserved_cu" {
  # 50 (not 0) so the applied plan differs from the built-in Default plan
  # (QUOTA_PLAN_DUPLICATE) without drifting quota attributes terraform tracks;
  # the service requires ElasticReservedCU to be a multiple of 50
  default = "50"
}

variable "update_elastic_reserved_cu" {
  default = "0"
}

variable "plan_name" {
  default = "TFPlan1737081505"
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


`, name)
}

// Case QuotaPlan_terraform无序数组BUG复现 9869
func TestAccAliCloudMaxComputeQuotaPlan_basic9869(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_quota_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMaxComputeQuotaPlanMap9869)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return maxcompute.NewMaxComputeServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMaxComputeQuotaPlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmcqp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMaxComputeQuotaPlanBasicDependence9869)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota": []map[string]interface{}{
						{
							"parameter": []map[string]interface{}{
								{
									"elastic_reserved_cu": "${var.elastic_reserved_cu}",
								},
							},
							"sub_quota_info_list": []map[string]interface{}{
								{
									"nick_name": "sub_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "0",
											"max_cu":              "20",
											"elastic_reserved_cu": "${var.elastic_reserved_cu}",
										},
									},
								},
								{
									"nick_name": "os_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "50",
											"max_cu":              "50",
											"elastic_reserved_cu": "0",
										},
									},
								},
							},
						},
					},
					"plan_name": "${var.plan_name}",
					"nickname":  "${alicloud_max_compute_quota.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name": acctest.CHECKSET,
						"nickname":  acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota": []map[string]interface{}{
						{
							"parameter": []map[string]interface{}{
								{
									"elastic_reserved_cu": "${var.update_elastic_reserved_cu}",
								},
							},
							"sub_quota_info_list": []map[string]interface{}{
								{
									"nick_name": "os_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "30",
											"max_cu":              "40",
											"elastic_reserved_cu": "0",
										},
									},
								},
								{
									"nick_name": "sub_${var.name}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "20",
											"max_cu":              "40",
											"elastic_reserved_cu": "${var.update_elastic_reserved_cu}",
										},
									},
								},
							},
						},
					},
					"is_effective": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_effective": "false",
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

var AlicloudMaxComputeQuotaPlanMap9869 = map[string]string{}

func AlicloudMaxComputeQuotaPlanBasicDependence9869(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "elastic_reserved_cu" {
  default = "0"
}

variable "update_elastic_reserved_cu" {
  default = "0"
}

variable "plan_name" {
  default = "TFPlan1737081505"
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


`, name)
}

// Test MaxCompute QuotaPlan. <<< Resource test cases, automatically generated.
