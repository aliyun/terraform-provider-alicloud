package alicloud

// collection-order: bypass data_disk -- pre-existing TypeList attribute, not added or changed by this destroy-routing fix
// collection-order: bypass ip_allocation_policy -- pre-existing TypeList attribute, not added or changed by this destroy-routing fix
// collection-order: bypass ip_allocation_policy.bond_policy.bonds -- pre-existing nested TypeList, not added or changed by this destroy-routing fix
// collection-order: bypass ip_allocation_policy.machine_type_policy -- pre-existing nested TypeList, not added or changed by this destroy-routing fix
// collection-order: bypass ip_allocation_policy.machine_type_policy.bonds -- pre-existing nested TypeList, not added or changed by this destroy-routing fix
// collection-order: bypass ip_allocation_policy.node_policy -- pre-existing nested TypeList, not added or changed by this destroy-routing fix
// collection-order: bypass ip_allocation_policy.node_policy.bonds -- pre-existing nested TypeList, not added or changed by this destroy-routing fix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				Config: testAccConfig(map[string]interface{}{
					"hostname":       "tfacc-host",
					"cluster_id":     "ehc-tfacc-placeholder",
					"node_group_id":  "ng-tfacc-placeholder",
					"login_password": "YourPassword_123",
					"machine_type":   "efg1.large",
					"node_type":      "ei",
					"user_data":      "I3RmYWNjCg==",
					"vswitch_id":     "vsw-tfacc-placeholder",
					"vpc_id":         "vpc-tfacc-placeholder",
					"install_pai":    "true",
					"data_disk": []map[string]interface{}{
						{
							"category":          "cloud_essd",
							"performance_level": "PL1",
							"size":              "2000",
						},
					},
					"ip_allocation_policy": []map[string]interface{}{
						{
							"machine_type_policy": []map[string]interface{}{
								{
									"bonds": []map[string]interface{}{
										{
											"name":   "tfacc-bond",
											"subnet": "vsw-tfacc-placeholder",
										},
									},
									"machine_type": "efg1.large",
								},
							},
							"bond_policy": []map[string]interface{}{
								{
									"bonds": []map[string]interface{}{
										{
											"name":   "tfacc-bond-bp",
											"subnet": "vsw-tfacc-placeholder-bp",
										},
									},
									"bond_default_subnet": "vsw-tfacc-placeholder-default",
								},
							},
							"node_policy": []map[string]interface{}{
								{
									"bonds": []map[string]interface{}{
										{
											"name":   "tfacc-bond-np",
											"subnet": "vsw-tfacc-placeholder-np",
										},
									},
									"node_id":  "ein-tfacc-placeholder",
									"hostname": "tfacc-node",
								},
							},
						},
					},
					"status":       "Running",
					"payment_type": "Bogus",
				}),
				ExpectError: regexp.MustCompile("expected payment_type to be one of"),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "computing_server", "discount_level", "payment_ratio", "period", "product_form", "renew_period", "renewal_status", "server_arch", "stage_num"},
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "discount_level", "login_password", "payment_ratio", "period", "product_form", "server_arch", "stage_num"},
			},
		},
	})
}

var AlicloudEfloNodeMap11276 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "discount_level", "login_password", "payment_ratio", "period", "product_form", "server_arch", "stage_num"},
			},
		},
	})
}

var AlicloudEfloNodeMap11280 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"period":            "36",
					"discount_level":    "36",
					"billing_cycle":     "1m",
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"period":            "36",
						"discount_level":    CHECKSET,
						"billing_cycle":     "1m",
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
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "discount_level", "login_password", "payment_ratio", "period", "product_form", "server_arch", "stage_num"},
			},
		},
	})
}

var AlicloudEfloNodeMap11275 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "discount_level", "login_password", "payment_ratio", "period", "product_form", "server_arch", "stage_num"},
			},
		},
	})
}

var AlicloudEfloNodeMap11278 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudEfloNodeBasicDependence11278(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// TestAccAliCloudEfloNode_basicSubscriptionDestroy covers the destroy path for a Subscription eflo
// node. payment_type is set explicitly to "Subscription" so that Delete routes to the
// RefundInstance refund branch. The boundary where payment_type is missing, empty or "<nil>"
// in state cannot be produced by this ACC case because Read backfills it, so it is covered by
// TestUnitAlicloudEfloNodeResolvePaymentType.
//
// zone, hpn_zone and computing_server are pinned because the provider has no data source for
// eflo zones, HPN zones or machine types (only alicloud_eflo_vpds and alicloud_eflo_subnets
// exist) while the order request needs all three. The pinned values are resolvable through the
// eflo selection APIs: DescribeZones lists cn-wulanchabu-a, ListMachineTypes lists
// efg2.C48eNH3ebn and ListMachineNetworkInfo resolves the A4 + efg2.C48eNH3ebn topology. They are
// also the values the stock cases in this file already pin. This is the documented exception:
// there is no supported selection alternative at this provider version. If the pinned
// specification stops being offered, Create fails at order time and the values must be re-pinned
// from the same selection APIs; the destroy path under test is unaffected.
func TestAccAliCloudEfloNode_basicSubscriptionDestroy(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloNodeMapSubscriptionDestroy)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloNodeBasicDependenceSubscriptionDestroy)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":      "Subscription",
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
						"payment_type":      "Subscription",
						"resource_group_id": CHECKSET,
						"period":            "36",
						"billing_cycle":     "1month",
						"classify":          "gpuserver",
						"zone":              "cn-wulanchabu-a",
						"product_form":      "instance",
						"hpn_zone":          "A4",
						"server_arch":       "bmserver",
						"computing_server":  "efg2.C48eNH3ebn",
						"renewal_status":    "AutoRenewal",
						"renew_period":      "36",
						"status":            "Unused",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "discount_level", "login_password", "payment_ratio", "period", "product_form", "server_arch", "stage_num"},
			},
		},
	})
}

var AlicloudEfloNodeMapSubscriptionDestroy = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudEfloNodeBasicDependenceSubscriptionDestroy(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Test Eflo Node. <<< Resource test cases, automatically generated.

func TestUnitAlicloudEfloNodeNormalizePaymentType(t *testing.T) {
	cases := []struct {
		name string
		raw  interface{}
		want string
		ok   bool
	}{
		{"missing key is unknown", map[string]interface{}{}["SubscriptionType"], "", false},
		{"nil value is unknown", nil, "", false},
		{"legacy nil placeholder is unknown", "<nil>", "", false},
		{"empty string is unknown", "", "", false},
		{"blank string is unknown", "  ", "", false},
		{"unsupported value is unknown", "PrePaid", "", false},
		{"subscription", "Subscription", "Subscription", true},
		{"pay as you go", "PayAsYouGo", "PayAsYouGo", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := normalizeEfloNodePaymentType(c.raw)
			if got != c.want || ok != c.ok {
				t.Fatalf("normalizeEfloNodePaymentType(%v) = (%q, %v), want (%q, %v)", c.raw, got, ok, c.want, c.ok)
			}
		})
	}
}

// TestUnitAlicloudEfloNodeResolvePaymentType covers the destroy boundary: a payment_type that is
// missing, empty or "<nil>" must not fall back to "Subscription" (which refunds a prepaid
// instance) and must not report success while deleting nothing.
func TestUnitAlicloudEfloNodeResolvePaymentType(t *testing.T) {
	t.Run("state billing type is used without a lookup", func(t *testing.T) {
		calls := 0
		got, err := resolveEfloNodePaymentType("Subscription", func() (map[string]interface{}, error) {
			calls++
			return map[string]interface{}{"SubscriptionType": "PayAsYouGo"}, nil
		})
		if err != nil || got != "Subscription" {
			t.Fatalf("got (%q, %v), want (Subscription, nil)", got, err)
		}
		if calls != 0 {
			t.Fatalf("billing lookup must not run when state already knows the billing type")
		}
	})

	t.Run("empty state resolves to pay as you go instead of refunding", func(t *testing.T) {
		got, err := resolveEfloNodePaymentType("", func() (map[string]interface{}, error) {
			return map[string]interface{}{"SubscriptionType": "PayAsYouGo"}, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "PayAsYouGo" {
			t.Fatalf("got %q, want PayAsYouGo so the destroy is routed to DeleteNode and not RefundInstance", got)
		}
	})

	t.Run("legacy nil marker resolves from the billing api", func(t *testing.T) {
		got, err := resolveEfloNodePaymentType("<nil>", func() (map[string]interface{}, error) {
			return map[string]interface{}{"SubscriptionType": "Subscription"}, nil
		})
		if err != nil || got != "Subscription" {
			t.Fatalf("got (%q, %v), want (Subscription, nil)", got, err)
		}
	})

	t.Run("unknown billing type fails instead of skipping the destroy", func(t *testing.T) {
		got, err := resolveEfloNodePaymentType("", func() (map[string]interface{}, error) {
			return map[string]interface{}{"ProductCode": "bccluster"}, nil
		})
		if err == nil {
			t.Fatalf("got (%q, nil), want an error so destroy is not reported as successful", got)
		}
	})

	t.Run("billing lookup error is propagated", func(t *testing.T) {
		_, err := resolveEfloNodePaymentType(nil, func() (map[string]interface{}, error) {
			return nil, fmt.Errorf("billing api unavailable")
		})
		if err == nil {
			t.Fatalf("want the billing lookup error to be propagated")
		}
	})
}

// TestUnitAlicloudEfloNodeDestroyAction pins the destroy routing: a pay-as-you-go node must never
// be sent to RefundInstance, and an unresolved billing type must not resolve to any API so the
// caller fails instead of removing state without deleting the node.
func TestUnitAlicloudEfloNodeDestroyAction(t *testing.T) {
	cases := map[string]string{
		"Subscription": "RefundInstance",
		"PayAsYouGo":   "DeleteNode",
		"":             "",
		"<nil>":        "",
	}

	for billingType, want := range cases {
		if got := efloNodeDestroyAction(billingType); got != want {
			t.Fatalf("efloNodeDestroyAction(%q) = %q, want %q", billingType, got, want)
		}
	}
}

// TestUnitAlicloudEfloNodeBillingRecordIdentity covers the identity check applied to the billing
// response: the destroy API and refund product code come from this record, so a record of another
// instance must never be used even when QueryAvailableInstances returns it.
func TestUnitAlicloudEfloNodeBillingRecordIdentity(t *testing.T) {
	requested := map[string]interface{}{"InstanceID": "e01-target", "SubscriptionType": "Subscription"}
	other := map[string]interface{}{"InstanceID": "e01-other", "SubscriptionType": "PayAsYouGo"}

	t.Run("matching record is selected from a multi-record response", func(t *testing.T) {
		got, ok := efloNodeBillingRecord([]interface{}{other, requested}, "e01-target")
		if !ok || got["InstanceID"] != "e01-target" {
			t.Fatalf("got (%v, %v), want the record of e01-target", got, ok)
		}
	})

	t.Run("record of another instance is rejected", func(t *testing.T) {
		if _, ok := efloNodeBillingRecord([]interface{}{other}, "e01-target"); ok {
			t.Fatalf("the billing record of e01-other must not decide the destroy of e01-target")
		}
	})

	t.Run("record without an instance id is rejected", func(t *testing.T) {
		if _, ok := efloNodeBillingRecord([]interface{}{map[string]interface{}{"SubscriptionType": "Subscription"}}, "e01-target"); ok {
			t.Fatalf("a record whose identity cannot be confirmed must not be selected")
		}
	})
}

// TestUnitAlicloudEfloNodeDeleteRouting drives resourceAliCloudEfloNodeDelete against an
// in-process RPC endpoint: it starts from the payment_type values that can appear in state
// (empty and the legacy "<nil>" placeholder) and asserts the destroy request that is actually
// sent, so the Read->Delete routing chain is covered without cloud access.
func TestUnitAlicloudEfloNodeDeleteRouting(t *testing.T) {
	const nodeID = "e01-target"
	billingRecord := func(subscriptionType, instanceID string) map[string]interface{} {
		return map[string]interface{}{
			"InstanceID":       instanceID,
			"ProductCode":      "bccluster",
			"ProductType":      "bccluster_eflocomputing_public_cn",
			"SubscriptionType": subscriptionType,
		}
	}

	cases := []struct {
		name       string
		stateType  string
		record     map[string]interface{}
		wantAction string
		wantErr    bool
	}{
		{
			name:       "unknown state is resolved to the refund API",
			stateType:  "",
			record:     billingRecord("Subscription", nodeID),
			wantAction: "RefundInstance",
		},
		{
			name:       "unknown state on a pay as you go record releases the node",
			stateType:  "<nil>",
			record:     billingRecord("PayAsYouGo", nodeID),
			wantAction: "DeleteNode",
		},
		{
			name:      "billing record of another instance fails without dropping state",
			stateType: "",
			record:    billingRecord("PayAsYouGo", "e01-other"),
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var actions []string
			deleted := false
			client := efloNodeRPCStubClient(t, func(w http.ResponseWriter, r *http.Request) {
				action := r.Form.Get("Action")
				actions = append(actions, action)
				w.Header().Set("Content-Type", "application/json")
				switch action {
				case "DescribeNode":
					if deleted {
						fmt.Fprint(w, `{"RequestId":"test","NodeId":""}`)
						return
					}
					fmt.Fprintf(w, `{"RequestId":"test","NodeId":%q,"OperatingState":"Unused"}`, nodeID)
				case "QueryAvailableInstances":
					record, err := json.Marshal(tc.record)
					if err != nil {
						t.Errorf("marshal billing record: %v", err)
						return
					}
					fmt.Fprintf(w, `{"Success":true,"Data":{"TotalCount":1,"PageNum":1,"InstanceList":[%s]}}`, record)
				case "RefundInstance":
					deleted = true
					fmt.Fprint(w, `{"Success":true}`)
				case "DeleteNode":
					deleted = true
					fmt.Fprint(w, `{"RequestId":"test"}`)
				default:
					t.Errorf("unexpected API call %s", action)
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprint(w, `{"Code":"Unexpected","Message":"unexpected API call"}`)
				}
			})

			d := schema.TestResourceDataRaw(t, resourceAliCloudEfloNode().Schema, map[string]interface{}{
				"payment_type": tc.stateType,
			})
			d.SetId(nodeID)

			err := resourceAliCloudEfloNodeDelete(d, client)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("want an error when the billing record belongs to another instance")
				}
				if d.Id() != nodeID {
					t.Fatalf("state id was dropped on failure: %q", d.Id())
				}
			} else if err != nil {
				t.Fatalf("unexpected destroy error: %v", err)
			}

			counts := map[string]int{}
			for _, action := range actions {
				counts[action]++
			}
			if tc.wantAction == "" {
				if counts["RefundInstance"] != 0 || counts["DeleteNode"] != 0 {
					t.Fatalf("destroy must not call a deletion API, got %v", actions)
				}
				return
			}
			if counts[tc.wantAction] != 1 {
				t.Fatalf("want exactly one %s call, got %v", tc.wantAction, actions)
			}
			other := "RefundInstance"
			if tc.wantAction == "RefundInstance" {
				other = "DeleteNode"
			}
			if counts[other] != 0 {
				t.Fatalf("unexpected %s call, got %v", other, actions)
			}
		})
	}
}

// efloNodeRPCStubClient returns a client whose RPC requests are served by handler, so resource
// entry points can be exercised without credentials or network access.
func efloNodeRPCStubClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse RPC form: %v", err)
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build test credential: %v", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	cfg := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-wulanchabu", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := cfg.Client()
	if err != nil {
		t.Fatalf("build test client: %v", err)
	}
	endpoints.Store("bssopenapi", endpoint)
	endpoints.Store("eflo_controller", endpoint)
	return client
}
