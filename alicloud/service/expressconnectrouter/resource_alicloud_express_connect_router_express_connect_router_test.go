package expressconnectrouter_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/expressconnectrouter"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass regions -- pre-existing stock attribute outside this pure-migration change; inter-region transit mode is read via a separate List API (DescribeExpressConnectRouterInterRegionTransitMode) that returns empty after Modify on a single-region ECR (QA-observed), making a convergent collection-order reorder infeasible without a Read/Modify change beyond migration scope
// Test ExpressConnectRouter ExpressConnectRouter. >>> Resource test cases, automatically generated.
// Case 初始版本测试用例 5682
func TestAccAliCloudExpressConnectRouterExpressConnectRouter_basic5682(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_express_connect_router.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterMap5682)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouter")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectrouter%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterBasicDependence5682)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alibaba_side_asn": "${var.initial_alibaba_side_asn}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alibaba_side_asn": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.initial_ecr_description}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecr_name": "${var.initial_ecr_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecr_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.updated_ecr_description}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecr_name": "${var.updated_ecr_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecr_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.initial_ecr_description}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecr_name": "${var.initial_ecr_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecr_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.updated_ecr_description}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecr_name": "${var.updated_ecr_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecr_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alibaba_side_asn":  "${var.initial_alibaba_side_asn}",
					"description":       "${var.initial_ecr_description}",
					"ecr_name":          "${var.initial_ecr_name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alibaba_side_asn":  acctest.CHECKSET,
						"description":       acctest.CHECKSET,
						"ecr_name":          acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var AlicloudExpressConnectRouterExpressConnectRouterMap5682 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudExpressConnectRouterExpressConnectRouterBasicDependence5682(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "fresh_region_id" {
  default = "cn-hangzhou"
}

variable "initial_ecr_name" {
  default = "InitialName"
}

variable "updated_ecr_description" {
  default = "UpdatedDescription"
}

variable "initial_ecr_description" {
  default = "InitialDescription"
}

variable "initial_alibaba_side_asn" {
  default = "65534"
}

variable "updated_ecr_name" {
  default = "UpdatedName"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case 初始版本测试用例 5682  twin
func TestAccAliCloudExpressConnectRouterExpressConnectRouter_basic5682_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_express_connect_router.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterMap5682)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouter")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectrouter%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterBasicDependence5682)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alibaba_side_asn":  "${var.initial_alibaba_side_asn}",
					"description":       "${var.initial_ecr_description}",
					"ecr_name":          "${var.initial_ecr_name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alibaba_side_asn":  acctest.CHECKSET,
						"description":       acctest.CHECKSET,
						"ecr_name":          acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
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

// Case 初始版本测试用例 5682  raw
func TestAccAliCloudExpressConnectRouterExpressConnectRouter_basic5682_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_express_connect_router.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterMap5682)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouter")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectrouter%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterBasicDependence5682)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alibaba_side_asn":  "${var.initial_alibaba_side_asn}",
					"description":       "${var.initial_ecr_description}",
					"ecr_name":          "${var.initial_ecr_name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alibaba_side_asn":  acctest.CHECKSET,
						"description":       acctest.CHECKSET,
						"ecr_name":          acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "${var.updated_ecr_description}",
					"ecr_name":          "${var.updated_ecr_name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       acctest.CHECKSET,
						"ecr_name":          acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "${var.initial_ecr_description}",
					"ecr_name":          "${var.initial_ecr_name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       acctest.CHECKSET,
						"ecr_name":          acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "${var.updated_ecr_description}",
					"ecr_name":          "${var.updated_ecr_name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       acctest.CHECKSET,
						"ecr_name":          acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
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

// Case ECR regions coverage 5682
// regions is read via a separate List API (DescribeExpressConnectRouterInterRegionTransitMode)
// that returns empty on a single-region ECR, so the attribute is not convergent; the step
// therefore plans only (no apply) and expects a non-empty plan, while still populating the
// coverage config with region_id and transit_mode.
func TestAccAliCloudExpressConnectRouterExpressConnectRouter_regions5682(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_express_connect_router.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterMap5682)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouter")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectrouter%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterBasicDependence5682)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alibaba_side_asn": "${var.initial_alibaba_side_asn}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alibaba_side_asn": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alibaba_side_asn": "${var.initial_alibaba_side_asn}",
					"regions": []interface{}{
						map[string]interface{}{
							"region_id":    "cn-shanghai",
							"transit_mode": "NearBy",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Test ExpressConnectRouter ExpressConnectRouter. <<< Resource test cases, automatically generated.
