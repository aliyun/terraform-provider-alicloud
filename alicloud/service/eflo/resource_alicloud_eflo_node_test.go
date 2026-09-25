package eflo_test

// TypeList Order Coverage: the multi-member TypeList attributes listed below
// are pre-existing stock debt surfaced by this directory-only migration (no
// schema or behavior change). They carry no convergent reorder acceptance
// test in this package; marked per scripts/collection-order/README.md
// stock-marker contract so the collection-order checker skips them.
// collection-order: bypass data_disk -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.bond_policy.bonds -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.machine_type_policy -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.machine_type_policy.bonds -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.node_policy -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.node_policy.bonds -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eflo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudEfloNode_basic10171(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloNodeMap10171)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloNode")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloNodeBasicDependence10171)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id": acctest.CHECKSET,
						"period":            "36",
						"discount_level":    acctest.CHECKSET,
						"billing_cycle":     "1month",
						"classify":          "gpuserver",
						"zone":              "cn-wulanchabu-b",
						"product_form":      "instance",
						"payment_ratio":     acctest.CHECKSET,
						"hpn_zone":          "B1",
						"server_arch":       "bmserver",
						"computing_server":  "efg2.C48cA3sen",
						"stage_num":         acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloNodeMap11276)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloNode")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloNodeBasicDependence11276)
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
						"resource_group_id": acctest.CHECKSET,
						"period":            "36",
						"discount_level":    acctest.CHECKSET,
						"billing_cycle":     "1month",
						"classify":          "gpuserver",
						"zone":              "cn-wulanchabu-a",
						"product_form":      "instance",
						"payment_ratio":     acctest.CHECKSET,
						"hpn_zone":          "A4",
						"server_arch":       "bmserver",
						"computing_server":  "efg2.C48eNH3ebn",
						"stage_num":         acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "discount_level", "login_password", "payment_ratio", "period", "product_form", "server_arch", "stage_num"},
			},
		},
	})
}

var AlicloudEfloNodeMap11276 = map[string]string{
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloNodeMap11280)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloNode")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloNodeBasicDependence11280)
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
						"resource_group_id": acctest.CHECKSET,
						"period":            "36",
						"discount_level":    acctest.CHECKSET,
						"billing_cycle":     "1month",
						"classify":          "gpuserver",
						"zone":              "cn-wulanchabu-a",
						"product_form":      "instance",
						"payment_ratio":     acctest.CHECKSET,
						"hpn_zone":          "A4",
						"server_arch":       "bmserver",
						"computing_server":  "efg2.C48eNH3ebn",
						"stage_num":         acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "discount_level", "login_password", "payment_ratio", "period", "product_form", "server_arch", "stage_num"},
			},
		},
	})
}

var AlicloudEfloNodeMap11280 = map[string]string{
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloNodeMap11275)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloNode")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloNodeBasicDependence11275)
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
						"resource_group_id": acctest.CHECKSET,
						"period":            "36",
						"discount_level":    acctest.CHECKSET,
						"billing_cycle":     "1m",
						"classify":          "gpuserver",
						"zone":              "cn-wulanchabu-a",
						"product_form":      "instance",
						"payment_ratio":     acctest.CHECKSET,
						"hpn_zone":          "A1",
						"server_arch":       "bmserver",
						"computing_server":  "efg1.nvga1n",
						"stage_num":         acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "discount_level", "login_password", "payment_ratio", "period", "product_form", "server_arch", "stage_num"},
			},
		},
	})
}

var AlicloudEfloNodeMap11275 = map[string]string{
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloNodeMap11278)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloNode")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloNodeBasicDependence11278)
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
						"resource_group_id": acctest.CHECKSET,
						"period":            "36",
						"discount_level":    acctest.CHECKSET,
						"billing_cycle":     "1month",
						"classify":          "gpuserver",
						"zone":              "cn-wulanchabu-a",
						"product_form":      "instance",
						"payment_ratio":     acctest.CHECKSET,
						"hpn_zone":          "A4",
						"server_arch":       "bmserver",
						"computing_server":  "efg2.C48eNH3ebn",
						"stage_num":         acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{"billing_cycle", "classify", "discount_level", "login_password", "payment_ratio", "period", "product_form", "server_arch", "stage_num"},
			},
		},
	})
}

var AlicloudEfloNodeMap11278 = map[string]string{
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
// RefundInstance refund branch. The boundary where payment_type is missing or unknown in state
// cannot be produced deterministically by this ACC case because the billing API is expected to
// return SubscriptionType for this prepaid node, so it is covered by
// TestUnitAlicloudEfloNodeDeleteRouting.
//
// zone, hpn_zone and computing_server are pinned because this provider version exposes no
// selection data source for eflo zones, HPN zones or machine types (only alicloud_eflo_vpds and
// alicloud_eflo_subnets exist) while the order request needs all three. The pinned values are
// also the ones the stock cases in this file already use. If the pinned specification stops being
// offered, Create fails at order time and the values must be re-pinned after checking
// DescribeZones, ListMachineTypes and ListMachineNetworkInfo; the destroy path under test is
// unaffected.
func TestAccAliCloudEfloNode_basicSubscriptionDestroy(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloNodeMapSubscriptionDestroy)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloNode")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloNodeBasicDependenceSubscriptionDestroy)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id": acctest.CHECKSET,
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
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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

// TestUnitAlicloudEfloNodeDeleteRouting drives the resource Delete entry point against an
// in-process RPC endpoint. It starts from every payment_type value that can appear in state and
// asserts the destroy request that is actually sent, so the Read->Delete routing chain and the
// billing-record identity check are covered without cloud access.
func TestUnitAlicloudEfloNodeDeleteRouting(t *testing.T) {
	const nodeID = "e01-target"
	billingRecord := func(subscriptionType string) map[string]interface{} {
		record := map[string]interface{}{
			"InstanceID":  nodeID,
			"ProductCode": "bccluster",
			"ProductType": "bccluster_eflocomputing_public_cn",
		}
		if subscriptionType != "" {
			record["SubscriptionType"] = subscriptionType
		}
		return record
	}
	otherInstanceRecord := billingRecord("PayAsYouGo")
	otherInstanceRecord["InstanceID"] = "e01-other"

	cases := []struct {
		name       string
		stateType  string
		record     map[string]interface{}
		gone       bool
		wantAction string
		wantErr    bool
	}{
		{
			name:       "empty state is resolved to the refund API",
			stateType:  "",
			record:     billingRecord("Subscription"),
			wantAction: "RefundInstance",
		},
		{
			name:       "known subscription state refunds the node",
			stateType:  "Subscription",
			record:     billingRecord("Subscription"),
			wantAction: "RefundInstance",
		},
		{
			name:      "already gone node is not deleted again",
			stateType: "",
			gone:      true,
		},
		{
			name:       "non-enumerated placeholder resolves from the billing record",
			stateType:  "<nil>",
			record:     billingRecord("PayAsYouGo"),
			wantAction: "DeleteNode",
		},
		{
			name:       "blank state is resolved from the billing record",
			stateType:  "  ",
			record:     billingRecord("PayAsYouGo"),
			wantAction: "DeleteNode",
		},
		{
			name:       "unsupported state value is resolved from the billing record",
			stateType:  "PrePaid",
			record:     billingRecord("Subscription"),
			wantAction: "RefundInstance",
		},
		{
			name:       "known pay as you go state releases the node without a billing lookup",
			stateType:  "PayAsYouGo",
			record:     nil,
			wantAction: "DeleteNode",
		},
		{
			name:      "billing record without a billing type fails without dropping state",
			stateType: "",
			record:    billingRecord(""),
			wantErr:   true,
		},
		{
			name:      "billing record of another instance fails without dropping state",
			stateType: "",
			record:    otherInstanceRecord,
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var mu sync.Mutex
			var actions []string
			deleted := false
			client := efloNodeRPCStubClient(t, func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				action := r.Form.Get("Action")
				actions = append(actions, action)
				w.Header().Set("Content-Type", "application/json")
				switch action {
				case "DescribeNode":
					if tc.gone {
						w.WriteHeader(http.StatusNotFound)
						fmt.Fprint(w, `{"Code":"RESOURCE_NOT_FOUND","Message":"node not found"}`)
						return
					}
					if deleted {
						fmt.Fprint(w, `{"RequestId":"test","NodeId":""}`)
						return
					}
					fmt.Fprintf(w, `{"RequestId":"test","NodeId":%q,"OperatingState":"Unused"}`, nodeID)
				case "QueryAvailableInstances":
					if tc.record == nil {
						t.Errorf("QueryAvailableInstances must not be called when state already knows the billing type")
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprint(w, `{"Code":"Unexpected","Message":"unexpected billing lookup"}`)
						return
					}
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

			node := eflo.ResourceNode()
			d := schema.TestResourceDataRaw(t, node.Schema, map[string]interface{}{
				"payment_type": tc.stateType,
			})
			d.SetId(nodeID)

			err := node.Delete(d, client)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("want an error when the billing type cannot be determined")
				}
				if d.Id() != nodeID {
					t.Fatalf("state id was dropped on failure: %q", d.Id())
				}
			} else if err != nil {
				t.Fatalf("unexpected destroy error: %v", err)
			}

			mu.Lock()
			actionsCopy := append([]string(nil), actions...)
			mu.Unlock()

			counts := map[string]int{}
			for _, action := range actionsCopy {
				counts[action]++
			}
			if tc.wantAction == "" {
				if counts["RefundInstance"] != 0 || counts["DeleteNode"] != 0 {
					t.Fatalf("destroy must not call a deletion API, got %v", actionsCopy)
				}
				return
			}
			if counts[tc.wantAction] != 1 {
				t.Fatalf("want exactly one %s call, got %v", tc.wantAction, actionsCopy)
			}
			other := "RefundInstance"
			if tc.wantAction == "RefundInstance" {
				other = "DeleteNode"
			}
			if counts[other] != 0 {
				t.Fatalf("unexpected %s call, got %v", other, actionsCopy)
			}
		})
	}
}

// TestUnitAlicloudEfloNodeBillingRecordSelection covers the Read-side billing lookup: the record
// stored as payment_type must belong to the node, so a response carrying only other records is
// treated as not found and Read keeps the billing type it already has.
func TestUnitAlicloudEfloNodeBillingRecordSelection(t *testing.T) {
	const nodeID = "e01-target"
	cases := []struct {
		name      string
		instances []map[string]interface{}
		wantType  string
		wantErr   bool
	}{
		{
			name: "matching record is selected from a multi-record response",
			instances: []map[string]interface{}{
				{"InstanceID": "e01-other", "SubscriptionType": "PayAsYouGo"},
				{"InstanceID": nodeID, "SubscriptionType": "Subscription"},
			},
			wantType: "Subscription",
		},
		{
			name:      "record of another instance is not used",
			instances: []map[string]interface{}{{"InstanceID": "e01-other", "SubscriptionType": "PayAsYouGo"}},
			wantErr:   true,
		},
		{
			name:      "record without an instance id is not used",
			instances: []map[string]interface{}{{"SubscriptionType": "Subscription"}},
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := efloNodeRPCStubClient(t, func(w http.ResponseWriter, r *http.Request) {
				if action := r.Form.Get("Action"); action != "QueryAvailableInstances" {
					t.Errorf("unexpected API call %s", action)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				list, err := json.Marshal(tc.instances)
				if err != nil {
					t.Errorf("marshal billing records: %v", err)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"Success":true,"Data":{"TotalCount":%d,"PageNum":1,"InstanceList":%s}}`, len(tc.instances), list)
			})

			d := schema.TestResourceDataRaw(t, eflo.ResourceNode().Schema, map[string]interface{}{"payment_type": ""})
			d.SetId(nodeID)

			record, err := eflo.NewEfloServiceV2(client).DescribeNodeQueryAvailableInstances(d)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("want an error when the response has no record for the node")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected billing lookup error: %v", err)
			}
			if got := fmt.Sprint(record["SubscriptionType"]); got != tc.wantType {
				t.Fatalf("got SubscriptionType %q, want %q", got, tc.wantType)
			}
		})
	}
}

// TestUnitAlicloudEfloNodeReadThenDelete drives the public resource lifecycle: Read persists
// payment_type into ResourceData from the billing response and Delete routes on that persisted
// value. A missing, null or foreign billing response must never be stored as a billing type, and
// Delete must fail without dropping state when the billing type cannot be determined.
func TestUnitAlicloudEfloNodeReadThenDelete(t *testing.T) {
	const nodeID = "e01-target"

	billingRecord := func(overrides map[string]interface{}) map[string]interface{} {
		record := map[string]interface{}{
			"InstanceID":  nodeID,
			"ProductCode": "bccluster",
			"ProductType": "bccluster_eflocomputing_public_cn",
		}
		for key, value := range overrides {
			record[key] = value
		}
		return record
	}
	foreignRecord := func() map[string]interface{} {
		return billingRecord(map[string]interface{}{"InstanceID": "e01-other", "SubscriptionType": "PayAsYouGo"})
	}

	cases := []struct {
		name          string
		stateType     string
		readRecords   []map[string]interface{}
		deleteRecords []map[string]interface{}
		wantStateType string
		wantAction    string
		wantErr       bool
	}{
		{
			name:          "missing billing type is not stored and Delete fails without dropping state",
			readRecords:   []map[string]interface{}{billingRecord(nil)},
			deleteRecords: []map[string]interface{}{billingRecord(nil)},
			wantStateType: "",
			wantErr:       true,
		},
		{
			name:          "null billing type is not stored and Delete fails without dropping state",
			readRecords:   []map[string]interface{}{billingRecord(map[string]interface{}{"SubscriptionType": nil})},
			deleteRecords: []map[string]interface{}{billingRecord(map[string]interface{}{"SubscriptionType": nil})},
			wantStateType: "",
			wantErr:       true,
		},
		{
			name:          "placeholder billing type is not stored and Delete fails without dropping state",
			readRecords:   []map[string]interface{}{billingRecord(map[string]interface{}{"SubscriptionType": "<nil>"})},
			deleteRecords: []map[string]interface{}{billingRecord(map[string]interface{}{"SubscriptionType": "<nil>"})},
			wantStateType: "",
			wantErr:       true,
		},
		{
			name:          "blank billing type is not stored and Delete fails without dropping state",
			readRecords:   []map[string]interface{}{billingRecord(map[string]interface{}{"SubscriptionType": "  "})},
			deleteRecords: []map[string]interface{}{billingRecord(map[string]interface{}{"SubscriptionType": "  "})},
			wantStateType: "",
			wantErr:       true,
		},
		{
			name:          "billing response overrides a different known state value",
			stateType:     "Subscription",
			readRecords:   []map[string]interface{}{billingRecord(map[string]interface{}{"SubscriptionType": "PayAsYouGo"})},
			wantStateType: "PayAsYouGo",
			wantAction:    "DeleteNode",
		},
		{
			name:          "Read keeps the known billing type when the response omits it",
			stateType:     "Subscription",
			readRecords:   []map[string]interface{}{billingRecord(nil)},
			deleteRecords: []map[string]interface{}{billingRecord(nil)},
			wantStateType: "Subscription",
			wantAction:    "RefundInstance",
		},
		{
			name:          "Read stores pay-as-you-go and Delete releases the node",
			readRecords:   []map[string]interface{}{billingRecord(map[string]interface{}{"SubscriptionType": "PayAsYouGo"})},
			wantStateType: "PayAsYouGo",
			wantAction:    "DeleteNode",
		},
		{
			name:          "Read stores subscription and Delete refunds the node",
			readRecords:   []map[string]interface{}{billingRecord(map[string]interface{}{"SubscriptionType": "Subscription"})},
			deleteRecords: []map[string]interface{}{billingRecord(nil)},
			wantStateType: "Subscription",
			wantAction:    "RefundInstance",
		},
		{
			name:          "Read ignores a foreign billing record and Delete fails without dropping state",
			stateType:     "Subscription",
			readRecords:   []map[string]interface{}{foreignRecord()},
			deleteRecords: []map[string]interface{}{foreignRecord()},
			wantStateType: "Subscription",
			wantErr:       true,
		},
		{
			name:          "empty state is resolved from the billing response returned during Delete",
			readRecords:   []map[string]interface{}{billingRecord(nil)},
			deleteRecords: []map[string]interface{}{billingRecord(map[string]interface{}{"SubscriptionType": "PayAsYouGo"})},
			wantStateType: "",
			wantAction:    "DeleteNode",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var mu sync.Mutex
			var actions []string
			var billingRequests []map[string]string
			readFinished := false
			deleteBillingIndex := 0
			deleted := false

			client := efloNodeRPCStubClient(t, func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
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
				case "DescribeDisks":
					fmt.Fprint(w, `{"RequestId":"test","Disks":{"Disk":[]}}`)
				case "ListTagResources":
					fmt.Fprint(w, `{"RequestId":"test","TagResources":{"TagResource":[]}}`)
				case "QueryAvailableInstances":
					billingRequests = append(billingRequests, map[string]string{
						"InstanceIDs": r.Form.Get("InstanceIDs"),
						"PageNum":     r.Form.Get("PageNum"),
						"PageSize":    r.Form.Get("PageSize"),
						"ProductCode": r.Form.Get("ProductCode"),
					})
					records := tc.readRecords
					if readFinished {
						if len(tc.deleteRecords) == 0 {
							t.Errorf("a billing lookup after Read requires deleteRecords in the case table")
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						end := deleteBillingIndex + 1
						if end > len(tc.deleteRecords) {
							end = len(tc.deleteRecords)
						}
						start := deleteBillingIndex
						if start >= len(tc.deleteRecords) {
							start = len(tc.deleteRecords) - 1
						}
						records = tc.deleteRecords[start:end]
						deleteBillingIndex++
					}
					recordList, err := json.Marshal(records)
					if err != nil {
						t.Errorf("marshal billing records: %v", err)
						return
					}
					fmt.Fprintf(w, `{"Success":true,"Data":{"TotalCount":%d,"PageNum":1,"InstanceList":%s}}`, len(records), recordList)
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

			node := eflo.ResourceNode()
			d := schema.TestResourceDataRaw(t, node.Schema, map[string]interface{}{
				"payment_type": tc.stateType,
			})
			d.SetId(nodeID)

			if err := node.Read(d, client); err != nil {
				t.Fatalf("unexpected read error: %v", err)
			}
			mu.Lock()
			readFinished = true
			mu.Unlock()
			if got := d.Get("payment_type"); got != tc.wantStateType {
				t.Fatalf("Read stored payment_type %q, want %q", got, tc.wantStateType)
			}

			err := node.Delete(d, client)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("want an error when the billing type cannot be determined")
				}
				if d.Id() != nodeID {
					t.Fatalf("state id was dropped on failure: %q", d.Id())
				}
			} else if err != nil {
				t.Fatalf("unexpected destroy error: %v", err)
			}

			mu.Lock()
			actionsCopy := append([]string(nil), actions...)
			billingRequestsCopy := append([]map[string]string(nil), billingRequests...)
			mu.Unlock()

			counts := map[string]int{}
			for _, action := range actionsCopy {
				counts[action]++
			}
			if tc.wantAction == "" {
				if counts["RefundInstance"] != 0 || counts["DeleteNode"] != 0 {
					t.Fatalf("destroy must not call a deletion API, got %v", actionsCopy)
				}
			} else {
				if counts[tc.wantAction] != 1 {
					t.Fatalf("want exactly one %s call, got %v", tc.wantAction, actionsCopy)
				}
				other := "RefundInstance"
				if tc.wantAction == "RefundInstance" {
					other = "DeleteNode"
				}
				if counts[other] != 0 {
					t.Fatalf("unexpected %s call, got %v", other, actionsCopy)
				}
			}

			if len(billingRequestsCopy) == 0 {
				t.Fatalf("QueryAvailableInstances was not called")
			}
			if got := billingRequestsCopy[0]; got["InstanceIDs"] != nodeID || got["ProductCode"] != "bccluster" {
				t.Fatalf("Read billing filter = %v, want InstanceIDs=%s ProductCode=bccluster", got, nodeID)
			}
			for _, request := range billingRequestsCopy[1:] {
				if request["InstanceIDs"] != nodeID || request["PageNum"] != "1" || request["PageSize"] != "100" {
					t.Fatalf("destroy billing filter = %v, want InstanceIDs=%s PageNum=1 PageSize=100", request, nodeID)
				}
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
	endpoints.Store("ecs", endpoint)
	return client
}
