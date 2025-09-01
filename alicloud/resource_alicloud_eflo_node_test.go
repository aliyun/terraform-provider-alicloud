package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEfloNode_basic10171(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloNodeMap10171)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloNodeBasicDependence10171)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"period":            "36",
					"discount_level":    "36",
					"billing_cycle":     "1month",
					"classify":          "gpuserver",
					"zone":              "cn-wulanchabu-b",
					"product_form":      "instance",
					"payment_ratio":     "0",
					"hpn_zone":          "B1",
					"server_arch":       "bmserver",
					"computing_server":  "efg2.C48cA3sen",
					"stage_num":         "36",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"period":            "36",
						"discount_level":    CHECKSET,
						"billing_cycle":     "1month",
						"classify":          "gpuserver",
						"zone":              "cn-wulanchabu-b",
						"product_form":      "instance",
						"payment_ratio":     CHECKSET,
						"hpn_zone":          "B1",
						"server_arch":       "bmserver",
						"computing_server":  "efg2.C48cA3sen",
						"stage_num":         CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
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
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "computing_server", "discount_level", "hpn_zone", "payment_ratio", "period", "product_form", "renew_period", "renewal_status", "server_arch", "stage_num", "zone"},
			},
		},
	})
}

var AlicloudEfloNodeMap10171 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEfloNodeBasicDependence10171(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Test Eflo Node. >>> Resource test cases, automatically generated.
// Case learn_eflocomputing_public_intl购买_install_true 11276
func TestAccAliCloudEfloNode_basic11276(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloNodeMap11276)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloNodeBasicDependence11276)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"period":            "36",
					"discount_level":    "36",
					"billing_cycle":     "1month",
					"classify":          "gpuserver",
					"zone":              "cn-wulanchabu-a",
					"product_form":      "instance",
					"payment_ratio":     "0",
					"hpn_zone":          "A4",
					"server_arch":       "bmserver",
					"computing_server":  "efg2.C48eNH3ebn",
					"stage_num":         "36",
					"renewal_status":    "AutoRenewal",
					"renew_period":      "36",
					"status":            "Unused",
					"install_pai":       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"period":            "36",
						"product_type":      "learn_eflocomputing_public_intl",
						"discount_level":    CHECKSET,
						"billing_cycle":     "1month",
						"classify":          "gpuserver",
						"zone":              "cn-wulanchabu-a",
						"product_form":      "instance",
						"payment_ratio":     CHECKSET,
						"hpn_zone":          "A4",
						"server_arch":       "bmserver",
						"computing_server":  "efg2.C48eNH3ebn",
						"stage_num":         CHECKSET,
						"renewal_status":    "AutoRenewal",
						"renew_period":      "36",
						"status":            "Unused",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
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
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "computing_server", "discount_level", "hpn_zone", "payment_ratio", "period", "product_form", "renew_period", "renewal_status", "server_arch", "stage_num", "zone", "install_pai"},
			},
		},
	})
}

var AlicloudEfloNodeMap11276 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEfloNodeBasicDependence11276(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case bccluster_eflocomputing_public_intl购买_install_false 11280
func TestAccAliCloudEfloNode_basic11280(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloNodeMap11280)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloNodeBasicDependence11280)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"period":            "36",
					"discount_level":    "36",
					"billing_cycle":     "1month",
					"classify":          "gpuserver",
					"zone":              "cn-wulanchabu-a",
					"product_form":      "instance",
					"payment_ratio":     "0",
					"hpn_zone":          "A4",
					"server_arch":       "bmserver",
					"computing_server":  "efg2.C48eNH3ebn",
					"stage_num":         "36",
					"renewal_status":    "AutoRenewal",
					"renew_period":      "36",
					"status":            "Unused",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"period":            "36",
						"discount_level":    CHECKSET,
						"billing_cycle":     "1month",
						"classify":          "gpuserver",
						"zone":              "cn-wulanchabu-a",
						"product_form":      "instance",
						"payment_ratio":     CHECKSET,
						"hpn_zone":          "A4",
						"server_arch":       "bmserver",
						"computing_server":  "efg2.C48eNH3ebn",
						"stage_num":         CHECKSET,
						"renewal_status":    "AutoRenewal",
						"renew_period":      "36",
						"status":            "Unused",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
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
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "computing_server", "discount_level", "hpn_zone", "payment_ratio", "period", "product_form", "renew_period", "renewal_status", "server_arch", "stage_num", "zone", "install_pai"},
			},
		},
	})
}

var AlicloudEfloNodeMap11280 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEfloNodeBasicDependence11280(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case learn_eflocomputing_public_cn购买_install_true 11275
func TestAccAliCloudEfloNode_basic11275(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloNodeMap11275)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloNodeBasicDependence11275)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"period":            "36",
					"discount_level":    "36",
					"billing_cycle":     "1month",
					"classify":          "gpuserver",
					"zone":              "cn-wulanchabu-a",
					"product_form":      "instance",
					"payment_ratio":     "0",
					"hpn_zone":          "A1",
					"server_arch":       "bmserver",
					"computing_server":  "efg1.nvga1n",
					"stage_num":         "36",
					"renewal_status":    "AutoRenewal",
					"renew_period":      "36",
					"status":            "Unused",
					"install_pai":       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"period":            "36",
						"discount_level":    CHECKSET,
						"billing_cycle":     "1month",
						"classify":          "gpuserver",
						"zone":              "cn-wulanchabu-a",
						"product_form":      "instance",
						"payment_ratio":     CHECKSET,
						"hpn_zone":          "A1",
						"server_arch":       "bmserver",
						"computing_server":  "efg1.nvga1n",
						"stage_num":         CHECKSET,
						"renewal_status":    "AutoRenewal",
						"renew_period":      "36",
						"status":            "Unused",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
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
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "computing_server", "discount_level", "hpn_zone", "payment_ratio", "period", "product_form", "renew_period", "renewal_status", "server_arch", "stage_num", "zone", "install_pai"},
			},
		},
	})
}

var AlicloudEfloNodeMap11275 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEfloNodeBasicDependence11275(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case bccluster_eflocomputing_public_cn购买_install_false 11278
func TestAccAliCloudEfloNode_basic11278(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloNodeMap11278)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloNodeBasicDependence11278)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"period":            "36",
					"discount_level":    "36",
					"billing_cycle":     "1month",
					"classify":          "gpuserver",
					"zone":              "cn-wulanchabu-a",
					"product_form":      "instance",
					"payment_ratio":     "0",
					"hpn_zone":          "A4",
					"server_arch":       "bmserver",
					"computing_server":  "efg2.C48eNH3ebn",
					"stage_num":         "36",
					"renewal_status":    "AutoRenewal",
					"renew_period":      "36",
					"status":            "Unused",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"period":            "36",
						"discount_level":    CHECKSET,
						"billing_cycle":     "1month",
						"classify":          "gpuserver",
						"zone":              "cn-wulanchabu-a",
						"product_form":      "instance",
						"payment_ratio":     CHECKSET,
						"hpn_zone":          "A4",
						"server_arch":       "bmserver",
						"computing_server":  "efg2.C48eNH3ebn",
						"stage_num":         CHECKSET,
						"renewal_status":    "AutoRenewal",
						"renew_period":      "36",
						"status":            "Unused",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
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
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "computing_server", "discount_level", "hpn_zone", "payment_ratio", "period", "product_form", "renew_period", "renewal_status", "server_arch", "stage_num", "zone", "install_pai"},
			},
		},
	})
}

var AlicloudEfloNodeMap11278 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEfloNodeBasicDependence11278(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Test Eflo Node. <<< Resource test cases, automatically generated.
