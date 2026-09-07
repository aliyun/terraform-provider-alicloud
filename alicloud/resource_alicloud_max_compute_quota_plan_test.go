package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test MaxCompute QuotaPlan. >>> Resource test cases, automatically generated.
// Case QuotaPlan_terraform测试 9566
func TestAccAliCloudMaxComputeQuotaPlan_basic9566(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_quota_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeQuotaPlanMap9566)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeQuotaPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smaxcomputequotaplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeQuotaPlanBasicDependence9566)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
									"nick_name": "sub_quota",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "0",
											"max_cu":              "20",
											"elastic_reserved_cu": "${var.elastic_reserved_cu}",
										},
									},
								},
								{
									"nick_name": "os_terrform",
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
					"nickname":  "os_terrform_p",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name": CHECKSET,
						"nickname":  CHECKSET,
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
									"nick_name": "sub_quota",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "30",
											"max_cu":              "40",
											"elastic_reserved_cu": "0",
										},
									},
								},
								{
									"nick_name": "os_terrform",
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
  default = "50"
}

variable "sub1" {
  default = "sub1"
}

variable "sub_max_cu" {
  default = "35"
}

variable "update_elastic_reserved_cu" {
  default = "0"
}

variable "plan_name" {
  default = "TFPlan1737081504"
}

variable "part_name" {
  default = "TFTest"
}


`, name)
}

// Case QuotaPlanApply_terraform测试(update) 9854
func TestAccAliCloudMaxComputeQuotaPlan_basic9854(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_quota_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeQuotaPlanMap9854)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeQuotaPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smaxcomputequotaplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeQuotaPlanBasicDependence9854)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
									"nick_name": "sub_quota",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "0",
											"max_cu":              "37",
											"elastic_reserved_cu": "${var.elastic_reserved_cu}",
										},
									},
								},
								{
									"nick_name": "os_terrform",
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
					"nickname":     "os_terrform_p",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name":    CHECKSET,
						"is_effective": "false",
						"nickname":     CHECKSET,
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
									"nick_name": "os_terrform",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "30",
											"max_cu":              "50",
											"elastic_reserved_cu": "0",
										},
									},
								},
								{
									"nick_name": "sub_quota",
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
					"is_effective": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_effective": "true",
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
  default = 50
}

variable "sub1" {
  default = "sub1"
}

variable "sub_max_cu" {
  default = "920"
}

variable "update_elastic_reserved_cu" {
  default = 0
}

variable "plan_name" {
  default = "TFPlan1737081505"
}

variable "part_name" {
  default = "TFTest"
}


`, name)
}

// Case QuotaPlanApply_terraform测试(create) 9853
func TestAccAliCloudMaxComputeQuotaPlan_basic9853(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_quota_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeQuotaPlanMap9853)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeQuotaPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smaxcomputequotaplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeQuotaPlanBasicDependence9853)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
									"nick_name": "${var.sub1}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "0",
											"max_cu":              "370",
											"elastic_reserved_cu": "${var.elastic_reserved_cu}",
										},
									},
								},
								{
									"nick_name": "os_${{ref(variable, partName)}}",
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
					"is_effective": "true",
					"nickname":     "os_${{ref(variable, partName)}}_p",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name":    CHECKSET,
						"is_effective": "true",
						"nickname":     "os_${{ref(variable, partName)}}_p",
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
  default = <<EOF
50
EOF
}

variable "sub1" {
  default = "sub1"
}

variable "sub_max_cu" {
  default = "214"
}

variable "update_elastic_reserved_cu" {
  default = <<EOF
0
EOF
}

variable "plan_name" {
  default = "TFPlan1737081505"
}

variable "part_name" {
  default = "TFTest"
}


`, name)
}

// Case QuotaPlan_terraform无序数组BUG复现 9869
func TestAccAliCloudMaxComputeQuotaPlan_basic9869(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_quota_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeQuotaPlanMap9869)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeQuotaPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smaxcomputequotaplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeQuotaPlanBasicDependence9869)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
									"nick_name": "${var.sub1}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "0",
											"max_cu":              "982",
											"elastic_reserved_cu": "${var.elastic_reserved_cu}",
										},
									},
								},
								{
									"nick_name": "os_${{ref(variable, partName)}}",
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name": CHECKSET,
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
									"nick_name": "os_${{ref(variable, partName)}}",
									"parameter": []map[string]interface{}{
										{
											"min_cu":              "30",
											"max_cu":              "757",
											"elastic_reserved_cu": "0",
										},
									},
								},
								{
									"nick_name": "${var.sub1}",
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
  default = <<EOF
50
EOF
}

variable "sub1" {
  default = "sub1"
}

variable "default_sub" {
  default = "os_TFTest"
}

variable "sub_max_cu" {
  default = "151"
}

variable "update_elastic_reserved_cu" {
  default = <<EOF
0
EOF
}

variable "plan_name" {
  default = "TFPlan1737081505"
}

variable "part_name" {
  default = "TFTest"
}


`, name)
}

// Test MaxCompute QuotaPlan. <<< Resource test cases, automatically generated.
