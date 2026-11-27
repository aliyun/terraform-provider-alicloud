package expressconnectrouter_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/expressconnectrouter"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testterraform "github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Test ExpressConnectRouter ExpressConnectRouterTrAssociation. >>> Resource test cases, automatically generated.
// Case TrAssociation 6355
func TestAccAliCloudExpressConnectRouterExpressConnectRouterTrAssociation_basic6355(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_tr_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterTrAssociationMap6355)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouterTrAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutertrassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterTrAssociationBasicDependence6355)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id":   acctest.DefaultRegionToTest,
					"ecr_id":                  "${alicloud_express_connect_router_express_connect_router.defaultpX0KlC.id}",
					"cen_id":                  "${alicloud_cen_instance.default418DC9.id}",
					"transit_router_owner_id": "${data.alicloud_account.current.id}",
					"transit_router_id":       "${alicloud_cen_transit_router.defaultRYcjsc.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_region_id": acctest.CHECKSET,
						"ecr_id":                acctest.CHECKSET,
						"description":           "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"${var.alowprefix1}", "${var.allowprefix3}", "${var.allowprefix2}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"${var.allowprefix3}", "${var.allowprefix2}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id":   acctest.DefaultRegionToTest,
					"ecr_id":                  "${alicloud_express_connect_router_express_connect_router.defaultpX0KlC.id}",
					"cen_id":                  "${alicloud_cen_instance.default418DC9.id}",
					"transit_router_owner_id": "${data.alicloud_account.current.id}",
					"allowed_prefixes": []string{
						"${var.alowprefix1}", "${var.allowprefix3}", "${var.allowprefix2}"},
					"transit_router_id": "${alicloud_cen_transit_router.defaultRYcjsc.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_region_id":   acctest.CHECKSET,
						"ecr_id":                  acctest.CHECKSET,
						"cen_id":                  acctest.CHECKSET,
						"transit_router_owner_id": acctest.CHECKSET,
						"allowed_prefixes.#":      "3",
						"transit_router_id":       acctest.CHECKSET,
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

var AlicloudExpressConnectRouterExpressConnectRouterTrAssociationMap6355 = map[string]string{
	"status":            acctest.CHECKSET,
	"create_time":       acctest.CHECKSET,
	"association_id":    acctest.CHECKSET,
	"transit_router_id": acctest.CHECKSET,
}

func AlicloudExpressConnectRouterExpressConnectRouterTrAssociationBasicDependence6355(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "alowprefix1" {
  default = "10.0.0.0/24"
}

variable "allowprefix2" {
  default = "10.0.1.0/24"
}

variable "allowprefix3" {
  default = "10.0.2.0/24"
}

variable "allowprefix4" {
  default = "10.0.3.0/24"
}

variable "asn" {
  default = "4200001003"
}

resource "alicloud_express_connect_router_express_connect_router" "defaultpX0KlC" {
  alibaba_side_asn = var.asn
}

resource "alicloud_cen_instance" "default418DC9" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "defaultRYcjsc" {
  cen_id = alicloud_cen_instance.default418DC9.id
}

data "alicloud_account" "current" {
}

`, name)
}

// Case TrAssociation 6355  twin
func TestAccAliCloudExpressConnectRouterExpressConnectRouterTrAssociation_basic6355_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_tr_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterTrAssociationMap6355)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouterTrAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutertrassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterTrAssociationBasicDependence6355)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id":   acctest.DefaultRegionToTest,
					"ecr_id":                  "${alicloud_express_connect_router_express_connect_router.defaultpX0KlC.id}",
					"cen_id":                  "${alicloud_cen_instance.default418DC9.id}",
					"transit_router_owner_id": "${data.alicloud_account.current.id}",
					"allowed_prefixes": []string{
						"${var.alowprefix1}", "${var.allowprefix3}", "${var.allowprefix2}"},
					"transit_router_id": "${alicloud_cen_transit_router.defaultRYcjsc.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_region_id":   acctest.CHECKSET,
						"ecr_id":                  acctest.CHECKSET,
						"cen_id":                  acctest.CHECKSET,
						"transit_router_owner_id": acctest.CHECKSET,
						"allowed_prefixes.#":      "3",
						"transit_router_id":       acctest.CHECKSET,
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

// Case TrAssociation 6355  raw
func TestAccAliCloudExpressConnectRouterExpressConnectRouterTrAssociation_basic6355_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_tr_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterTrAssociationMap6355)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouterTrAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutertrassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterTrAssociationBasicDependence6355)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id":   acctest.DefaultRegionToTest,
					"ecr_id":                  "${alicloud_express_connect_router_express_connect_router.defaultpX0KlC.id}",
					"cen_id":                  "${alicloud_cen_instance.default418DC9.id}",
					"transit_router_owner_id": "${data.alicloud_account.current.id}",
					"allowed_prefixes": []string{
						"${var.alowprefix1}", "${var.allowprefix3}", "${var.allowprefix2}"},
					"transit_router_id": "${alicloud_cen_transit_router.defaultRYcjsc.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_region_id":   acctest.CHECKSET,
						"ecr_id":                  acctest.CHECKSET,
						"cen_id":                  acctest.CHECKSET,
						"transit_router_owner_id": acctest.CHECKSET,
						"allowed_prefixes.#":      "3",
						"transit_router_id":       acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"${var.allowprefix3}", "${var.allowprefix2}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "0",
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

// Case TrAssociation allowed_prefixes collection order 6355
func TestAccAliCloudExpressConnectRouterExpressConnectRouterTrAssociation_allowedPrefixesOrder6355(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_tr_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterTrAssociationMap6355)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouterTrAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutertrassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterTrAssociationBasicDependence6355)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id":   acctest.DefaultRegionToTest,
					"ecr_id":                  "${alicloud_express_connect_router_express_connect_router.defaultpX0KlC.id}",
					"cen_id":                  "${alicloud_cen_instance.default418DC9.id}",
					"transit_router_owner_id": "${data.alicloud_account.current.id}",
					"allowed_prefixes": []string{
						"${var.alowprefix1}", "${var.allowprefix2}", "${var.allowprefix3}"},
					"transit_router_id": "${alicloud_cen_transit_router.defaultRYcjsc.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id":   acctest.DefaultRegionToTest,
					"ecr_id":                  "${alicloud_express_connect_router_express_connect_router.defaultpX0KlC.id}",
					"cen_id":                  "${alicloud_cen_instance.default418DC9.id}",
					"transit_router_owner_id": "${data.alicloud_account.current.id}",
					"allowed_prefixes": []string{
						"${var.allowprefix3}", "${var.allowprefix2}", "${var.alowprefix1}"},
					"transit_router_id": "${alicloud_cen_transit_router.defaultRYcjsc.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id":   acctest.DefaultRegionToTest,
					"ecr_id":                  "${alicloud_express_connect_router_express_connect_router.defaultpX0KlC.id}",
					"cen_id":                  "${alicloud_cen_instance.default418DC9.id}",
					"transit_router_owner_id": "${data.alicloud_account.current.id}",
					"allowed_prefixes": []string{
						"${var.alowprefix1}", "${var.allowprefix2}", "${var.allowprefix3}"},
					"transit_router_id": "${alicloud_cen_transit_router.defaultRYcjsc.transit_router_id}",
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// Case TrAssociation attr extend (create_attachment / allowed_prefixes_mode / description) 6355
func TestAccAliCloudExpressConnectRouterExpressConnectRouterTrAssociation_attrExtend6355(t *testing.T) {
	t.Setenv("DEBUG", "terraform")
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_tr_association.default"
	ra := acctest.ResourceAttrInit(resourceId, TRAssociationAttrExtendMap6355)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouterTrAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutertrassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterTrAssociationBasicDependence6355)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id":   acctest.DefaultRegionToTest,
					"ecr_id":                  "${alicloud_express_connect_router_express_connect_router.defaultpX0KlC.id}",
					"cen_id":                  "${alicloud_cen_instance.default418DC9.id}",
					"transit_router_owner_id": "${data.alicloud_account.current.id}",
					"transit_router_id":       "${alicloud_cen_transit_router.defaultRYcjsc.transit_router_id}",
					"allowed_prefixes": []string{
						"${var.alowprefix1}", "${var.allowprefix3}", "${var.allowprefix2}"},
					"create_attachment":     true,
					"allowed_prefixes_mode": "MatchMode",
					"description":           "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"create_attachment":     "true",
						"allowed_prefixes_mode": "MatchMode",
						"description":           "",
						"allowed_prefixes.#":    "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_attachment"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testacc-desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"create_attachment":     "true",
						"allowed_prefixes_mode": "MatchMode",
						"description":           "tf-testacc-desc",
						"allowed_prefixes.#":    "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_attachment"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes_mode": "IncrementalMode",
					"description":           "tf-testacc-desc-updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes_mode": "IncrementalMode",
						"description":           "tf-testacc-desc-updated",
						"allowed_prefixes.#":    "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_attachment"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        "",
						"allowed_prefixes.#": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_attachment"},
			},
		},
	})
}

var TRAssociationAttrExtendMap6355 = map[string]string{
	"status":                acctest.CHECKSET,
	"create_time":           acctest.CHECKSET,
	"association_id":        acctest.CHECKSET,
	"transit_router_id":     acctest.CHECKSET,
	"create_attachment":     acctest.CHECKSET,
	"allowed_prefixes_mode": acctest.CHECKSET,
}

func TestAccAliCloudExpressConnectRouterExpressConnectRouterTrAssociation_clearAllowedPrefixes6355(t *testing.T) {
	t.Setenv("DEBUG", "terraform")
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_tr_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterTrAssociationMap6355)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouterTrAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutertrassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterTrAssociationBasicDependence6355)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkTRAssociationDestroyWithParent(rac.CheckResourceDestroy()),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id":   acctest.DefaultRegionToTest,
					"ecr_id":                  "${alicloud_express_connect_router_express_connect_router.defaultpX0KlC.id}",
					"cen_id":                  "${alicloud_cen_instance.default418DC9.id}",
					"transit_router_owner_id": "${data.alicloud_account.current.id}",
					"transit_router_id":       "${alicloud_cen_transit_router.defaultRYcjsc.transit_router_id}",
					"allowed_prefixes_mode":   "MatchMode",
					"create_attachment":       false,
				}),
				ExpectError: regexp.MustCompile(`allowed_prefixes_mode.*non-empty.*allowed_prefixes`),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id":   acctest.DefaultRegionToTest,
					"ecr_id":                  "${alicloud_express_connect_router_express_connect_router.defaultpX0KlC.id}",
					"cen_id":                  "${alicloud_cen_instance.default418DC9.id}",
					"transit_router_owner_id": "${data.alicloud_account.current.id}",
					"transit_router_id":       "${alicloud_cen_transit_router.defaultRYcjsc.transit_router_id}",
					"allowed_prefixes": []string{
						"${var.alowprefix1}", "${var.allowprefix2}"},
					"allowed_prefixes_mode": "MatchMode",
					"create_attachment":     false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#":    "2",
						"allowed_prefixes_mode": "MatchMode",
						"create_attachment":     "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_attachment"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes":      acctest.REMOVEKEY,
					"allowed_prefixes_mode": "IncrementalMode",
				}),
				ExpectError: regexp.MustCompile(`allowed_prefixes_mode.*non-empty.*allowed_prefixes`),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"${var.alowprefix1}", "${var.allowprefix2}"},
					"allowed_prefixes_mode": "MatchMode",
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"${var.alowprefix1}", "${var.allowprefix2}"},
					"allowed_prefixes_mode": "MatchMode",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#":    "2",
						"allowed_prefixes_mode": "MatchMode",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_attachment"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes":      []string{},
					"allowed_prefixes_mode": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#":    "0",
						"allowed_prefixes_mode": "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_attachment"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes":      []string{},
					"allowed_prefixes_mode": acctest.REMOVEKEY,
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"${var.alowprefix1}", "${var.allowprefix2}"},
					"allowed_prefixes_mode": "MatchMode",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#":    "2",
						"allowed_prefixes_mode": "MatchMode",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_attachment"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes":      acctest.REMOVEKEY,
					"allowed_prefixes_mode": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#":    "0",
						"allowed_prefixes_mode": "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_attachment"},
			},
		},
	})
}

func checkTRAssociationDestroyWithParent(fallback resource.TestCheckFunc) resource.TestCheckFunc {
	return func(state *testterraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "alicloud_express_connect_router_tr_association" {
				continue
			}
			parts := strings.Split(rs.Primary.ID, ":")
			if len(parts) != 3 {
				return fmt.Errorf("invalid TR association ID %q", rs.Primary.ID)
			}
			service := expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
			_, err := service.DescribeExpressConnectRouterExpressConnectRouter(parts[0])
			if errs.NotFoundError(err) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("check ECR parent before TR association destroy: %w", err)
			}
			return fallback(state)
		}
		return nil
	}
}

// fastTRAssociationWaits removes the production status-poll delay so unit tests
// stay fast while still exercising the real wait path.
func fastTRAssociationWaits(t *testing.T) {
	t.Helper()
	t.Cleanup(expressconnectrouter.SetTRAssociationStateDelayForTest(time.Millisecond))
}

func TestUnitTRAssociationACheckerExhaustsDescribePages(t *testing.T) {
	const (
		resourceId = "alicloud_express_connect_router_tr_association.default"
		id         = "ecr-unit:assoc-unit:tr-unit"
	)
	target := `{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE"}`
	for _, tc := range []struct {
		name       string
		secondPage string
		exists     bool
		malformed  bool
	}{
		{name: "target_on_later_page", secondPage: `{"AssociationList":[` + target + `],"Success":true}`, exists: true},
		{name: "absent_after_all_pages", secondPage: `{"AssociationList":[],"Success":true}`},
		{name: "coded_parent_not_found", secondPage: `{"Code":"ResourceNotFound.EcrId","Message":"EcrId not found","Success":false}`},
		{name: "deleted_parent_empty_response", secondPage: `{}`},
		{name: "wrong_ecr", secondPage: `{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-other","TransitRouterId":"tr-unit"}],"Success":true}`},
		{name: "wrong_association", secondPage: `{"AssociationList":[{"AssociationId":"assoc-other","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit"}],"Success":true}`},
		{name: "foreign_association_needs_no_other_identity", secondPage: `{"AssociationList":[{"AssociationId":"assoc-other"}],"Success":true}`},
		{name: "wrong_transit_router", secondPage: `{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-other"}],"Success":true}`},
		{name: "wrong_node_type", secondPage: `{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"VPC","EcrId":"ecr-unit","TransitRouterId":"tr-unit"}],"Success":true}`},
		{name: "missing_list", secondPage: `{"Success":true}`, malformed: true},
		{name: "empty_code_not_found_message", secondPage: `{"Code":"","Message":"NotFound","Success":false}`, malformed: true},
		{name: "malformed_list", secondPage: `{"AssociationList":{},"Success":true}`, malformed: true},
		{name: "malformed_element", secondPage: `{"AssociationList":["unexpected"],"Success":true}`, malformed: true},
		{name: "empty_object", secondPage: `{"AssociationList":[{}],"Success":true}`, malformed: true},
		{name: "partial_match_missing_transit_router", secondPage: `{"AssociationList":[{"AssociationId":"assoc-unit","EcrId":"ecr-unit"}],"Success":true}`, malformed: true},
		{name: "partial_match_wrong_typed_ecr", secondPage: `{"AssociationList":[{"AssociationId":"assoc-unit","EcrId":42}],"Success":true}`, malformed: true},
		{name: "malformed_next_token", secondPage: `{"AssociationList":[],"NextToken":42,"Success":true}`, malformed: true},
		{name: "null_next_token", secondPage: `{"AssociationList":[],"NextToken":null,"Success":true}`, malformed: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var mu sync.Mutex
			var tokens []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %s", err)
				}
				w.Header().Set("Content-Type", "application/json")
				if action := r.Form.Get("Action"); action != "DescribeExpressConnectRouterAssociation" {
					t.Errorf("unexpected action: %q", action)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
					return
				}
				token := r.Form.Get("NextToken")
				mu.Lock()
				tokens = append(tokens, token)
				mu.Unlock()
				switch token {
				case "":
					_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-foreign","TransitRouterId":"tr-foreign"}],"NextToken":"page-2","Success":true}`))
				case "page-2":
					_, _ = w.Write([]byte(tc.secondPage))
				default:
					t.Errorf("unexpected NextToken: %q", token)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected NextToken"}`))
				}
			}))
			t.Cleanup(server.Close)

			service := expressconnectrouter.NewService(unitTestClient(t, server.URL))
			var object map[string]interface{}
			rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &object, func() interface{} {
				return service
			}, "DescribeExpressConnectRouterExpressConnectRouterTrAssociation")
			state := &testterraform.State{Modules: []*testterraform.ModuleState{{
				Path: []string{"root"},
				Resources: map[string]*testterraform.ResourceState{
					resourceId: {
						Type:    "alicloud_express_connect_router_tr_association",
						Primary: &testterraform.InstanceState{ID: id},
					},
				},
			}}}

			existsErr := rc.CheckResourceExists()(state)
			destroyErr := rc.CheckResourceDestroy()(state)
			refreshObject, refreshStatus, refreshErr := service.ExpressConnectRouterExpressConnectRouterTrAssociationStateRefreshFunc(id, "$.Status", nil)()
			if tc.malformed {
				if existsErr == nil {
					t.Fatal("a malformed Describe page must fail CheckResourceExists")
				}
				if destroyErr == nil {
					t.Fatal("a malformed Describe page must not pass CheckResourceDestroy")
				}
				if refreshErr == nil {
					t.Fatal("a malformed Describe page must fail the shared state refresh")
				}
			} else if tc.exists {
				if existsErr != nil {
					t.Fatalf("the later-page target must pass CheckResourceExists: %s", existsErr)
				}
				if destroyErr == nil {
					t.Fatal("a later-page target that still exists must fail CheckResourceDestroy")
				}
				if refreshErr != nil || refreshObject == nil || refreshStatus != "ACTIVE" {
					t.Fatalf("the later-page target must refresh as ACTIVE, object=%#v status=%q err=%v", refreshObject, refreshStatus, refreshErr)
				}
			} else {
				if existsErr == nil {
					t.Fatal("a target absent from every page must fail CheckResourceExists")
				}
				if destroyErr != nil {
					t.Fatalf("a target absent from every page must pass CheckResourceDestroy: %s", destroyErr)
				}
				if refreshErr != nil || refreshObject != nil || refreshStatus != "" {
					t.Fatalf("confirmed absence must refresh as empty, object=%#v status=%q err=%v", refreshObject, refreshStatus, refreshErr)
				}
			}

			mu.Lock()
			gotTokens := append([]string(nil), tokens...)
			mu.Unlock()
			expectedTokens := []string{"", "page-2", "", "page-2", "", "page-2"}
			if len(gotTokens) != len(expectedTokens) {
				t.Fatalf("both checkers and the shared refresh must exhaust two pages, got tokens %#v", gotTokens)
			}
			for i := range expectedTokens {
				if gotTokens[i] != expectedTokens[i] {
					t.Fatalf("unexpected page tokens: got %#v, want %#v", gotTokens, expectedTokens)
				}
			}
		})
	}
}

func TestUnitTRAssociationSharedDescribePropagatesTransportFailure(t *testing.T) {
	const (
		resourceId = "alicloud_express_connect_router_tr_association.default"
		id         = "ecr-unit:assoc-unit:tr-unit"
	)
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{name: "terminal_400", status: http.StatusBadRequest, body: `{"Code":"InvalidParameter","Message":"terminal","Success":false}`},
		{name: "empty_code_404", status: http.StatusNotFound, body: `{"Code":"","Message":"NotFound","Success":false}`},
		{name: "missing_code_404", status: http.StatusNotFound, body: `{"Message":"NotFound","Success":false}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var requests int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requests, 1)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			t.Cleanup(server.Close)
			service := expressconnectrouter.NewService(unitTestClient(t, server.URL))
			var object map[string]interface{}
			rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &object, func() interface{} {
				return service
			}, "DescribeExpressConnectRouterExpressConnectRouterTrAssociation")
			state := &testterraform.State{Modules: []*testterraform.ModuleState{{
				Path: []string{"root"},
				Resources: map[string]*testterraform.ResourceState{
					resourceId: {
						Type:    "alicloud_express_connect_router_tr_association",
						Primary: &testterraform.InstanceState{ID: id},
					},
				},
			}}}
			if _, err := service.DescribeExpressConnectRouterExpressConnectRouterTrAssociation(id); err == nil || errs.NotFoundError(err) {
				t.Fatalf("unsupported HTTP failure must retain its original error, got %v", err)
			}
			if err := rc.CheckResourceExists()(state); err == nil {
				t.Fatal("unsupported HTTP failure must fail CheckResourceExists")
			}
			if err := rc.CheckResourceDestroy()(state); err == nil {
				t.Fatal("unsupported HTTP failure must not pass CheckResourceDestroy")
			}
			refreshed, status, err := service.ExpressConnectRouterExpressConnectRouterTrAssociationStateRefreshFunc(id, "$.Status", nil)()
			if err == nil || refreshed != nil || status != "" {
				t.Fatalf("unsupported HTTP failure must fail refresh, object=%#v status=%q err=%v", refreshed, status, err)
			}
			if got := atomic.LoadInt32(&requests); got != 4 {
				t.Fatalf("four independent Describe calls = %d, want 4", got)
			}
		})
	}
}

// TestUnitTRAssociationCreateSerializesCreateAttachment covers the create request
// through the public resource entry: an omitted create_attachment must not reach
// the wire, while an explicit false/true must be serialized as configured.
func TestUnitTRAssociationCreateSerializesCreateAttachment(t *testing.T) {
	fastTRAssociationWaits(t)
	for _, tc := range []struct {
		name     string
		value    cty.Value
		expected string
	}{
		{name: "unset", value: cty.NullVal(cty.Bool), expected: ""},
		{name: "explicit_false", value: cty.False, expected: "false"},
		{name: "explicit_true", value: cty.True, expected: "true"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var mu sync.Mutex
			createForm := map[string][]string{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %s", err)
				}
				w.Header().Set("Content-Type", "application/json")
				switch action := r.Form.Get("Action"); action {
				case "CreateExpressConnectRouterAssociation":
					mu.Lock()
					createForm = r.PostForm
					mu.Unlock()
					_, _ = w.Write([]byte(`{"AssociationId":"assoc-unit","Success":true}`))
				case "DescribeExpressConnectRouterAssociation":
					_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE","AllowedPrefixes":[]}],"Success":true}`))
				default:
					t.Errorf("unexpected action: %q", action)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
				}
			}))
			t.Cleanup(server.Close)

			resourceSchema := expressconnectrouter.ResourceTRAssociation()
			attrs := map[string]cty.Value{}
			for name, attributeType := range resourceSchema.CoreConfigSchema().ImpliedType().AttributeTypes() {
				attrs[name] = cty.NullVal(attributeType)
			}
			attrs["ecr_id"] = cty.StringVal("ecr-unit")
			attrs["transit_router_id"] = cty.StringVal("tr-unit")
			attrs["association_region_id"] = cty.StringVal("cn-hangzhou")
			attrs["cen_id"] = cty.StringVal("cen-unit")
			attrs["transit_router_owner_id"] = cty.NumberIntVal(123456)
			attrs["allowed_prefixes"] = cty.ListVal([]cty.Value{cty.StringVal("10.0.0.0/24")})
			attrs["create_attachment"] = tc.value
			rawConfig := cty.ObjectVal(attrs)

			diff, err := resourceSchema.Diff(context.Background(), nil, sdktf.NewResourceConfigShimmed(rawConfig, resourceSchema.CoreConfigSchema()), nil)
			if err != nil {
				t.Fatalf("plan through the public resource entry failed: %s", err)
			}
			// PlanResourceChange carries the raw configuration on the diff; mirror
			// it so create reads the same source as production.
			diff.RawConfig = rawConfig
			diff.Meta = map[string]interface{}{
				schema.TimeoutKey: map[string]interface{}{
					schema.TimeoutCreate: 5 * time.Minute,
					schema.TimeoutUpdate: 5 * time.Minute,
					schema.TimeoutDelete: 5 * time.Minute,
				},
			}
			if _, diags := resourceSchema.Apply(context.Background(), nil, diff, unitTestClient(t, server.URL)); diags.HasError() {
				t.Fatalf("apply through the public resource entry failed: %v", diags)
			}

			mu.Lock()
			defer mu.Unlock()
			values, exists := createForm["CreateAttachment"]
			if tc.expected == "" {
				if exists {
					t.Fatalf("unset create_attachment must not be serialized, got %v", values)
				}
				return
			}
			if !exists || len(values) != 1 || values[0] != tc.expected {
				t.Fatalf("CreateAttachment serialized as %v, want %q", values, tc.expected)
			}
		})
	}
}

// TestUnitTRAssociationCreateRetriesTransientErrors covers the documented
// retryable create error: InternalErrorTemp must be retried and a later
// success must not surface an error to Terraform.
func TestUnitTRAssociationCreateRetriesTransientErrors(t *testing.T) {
	fastTRAssociationWaits(t)
	var mu sync.Mutex
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "CreateExpressConnectRouterAssociation":
			mu.Lock()
			attempts++
			current := attempts
			mu.Unlock()
			if current <= 2 {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"Code":"InternalErrorTemp","Message":"temporary failure","RequestId":"test"}`))
				return
			}
			_, _ = w.Write([]byte(`{"AssociationId":"assoc-unit","Success":true}`))
		case "DescribeExpressConnectRouterAssociation":
			_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE","AllowedPrefixes":[]}],"Success":true}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	if diags := applyTRAssociationCreate(t, server.URL); diags.HasError() {
		t.Fatalf("create should recover from InternalErrorTemp, got %v", diags)
	}
	mu.Lock()
	defer mu.Unlock()
	if attempts != 3 {
		t.Fatalf("create called %d times, want 3", attempts)
	}
}

// TestUnitTRAssociationCreateFastFailsTerminalErrors keeps terminal parameter
// errors non-retryable so they surface immediately to the user.
func TestUnitTRAssociationCreateFastFailsTerminalErrors(t *testing.T) {
	fastTRAssociationWaits(t)
	var mu sync.Mutex
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "CreateExpressConnectRouterAssociation":
			mu.Lock()
			attempts++
			mu.Unlock()
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"IllegalParamFormat.EcrId","Message":"terminal failure","RequestId":"test"}`))
		case "DescribeExpressConnectRouterAssociation":
			_, _ = w.Write([]byte(`{"AssociationList":[],"Success":true}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	if diags := applyTRAssociationCreate(t, server.URL); !diags.HasError() {
		t.Fatal("terminal create errors must fail the apply")
	}
	mu.Lock()
	defer mu.Unlock()
	if attempts != 1 {
		t.Fatalf("terminal create error must not be retried, called %d times", attempts)
	}
}

func applyTRAssociationCreate(t *testing.T, serverURL string) (diags diag.Diagnostics) {
	t.Helper()
	return applyTRAssociationCreateWithTimeout(t, serverURL, 5*time.Minute)
}

func TestUnitTRAssociationDeleteAttachmentPolicy(t *testing.T) {
	fastTRAssociationWaits(t)
	for _, tc := range []struct {
		name             string
		createAttachment string
		failWith         string
		succeedAfter     int
		wantCalls        int
		wantFlag         string
		wantError        string
	}{
		{name: "configured_true_may_delete_owned_attachment", createAttachment: "true", succeedAfter: 1, wantCalls: 1, wantFlag: "true"},
		{name: "explicit_false_preserves_attachment", createAttachment: "false", failWith: "OperationDenied.EcrAttachmentExisted", succeedAfter: 1, wantCalls: 1, wantError: "OperationDenied.EcrAttachmentExisted"},
		{name: "imported_unknown_preserves_attachment", failWith: "OperationDenied.EcrAttachmentExisted", succeedAfter: 1, wantCalls: 1, wantError: "OperationDenied.EcrAttachmentExisted"},
		{name: "transient_conflict_retries_without_widening", createAttachment: "false", failWith: "Conflict.Lock", succeedAfter: 2, wantCalls: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var deleteFlags []string
			var listCalls int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %s", err)
				}
				w.Header().Set("Content-Type", "application/json")
				switch action := r.Form.Get("Action"); action {
				case "ListTransitRouterEcrAttachments":
					atomic.AddInt32(&listCalls, 1)
					t.Errorf("attachment ownership must not be inferred from a pre-delete lookup")
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected lookup"}`))
				case "DeleteExpressConnectRouterAssociation":
					attempt := len(deleteFlags) + 1
					deleteFlags = append(deleteFlags, r.PostForm.Get("DeleteAttachment"))
					if tc.failWith != "" && (tc.wantError != "" || attempt < tc.succeedAfter) {
						w.WriteHeader(http.StatusBadRequest)
						_, _ = w.Write([]byte(`{"Code":"` + tc.failWith + `","Message":"delete fixture","RequestId":"test"}`))
						return
					}
					_, _ = w.Write([]byte(`{"Success":true}`))
				case "DescribeExpressConnectRouterAssociation":
					_, _ = w.Write([]byte(`{"AssociationList":[],"Success":true}`))
				default:
					t.Errorf("unexpected action: %q", action)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
				}
			}))
			t.Cleanup(server.Close)

			resourceSchema := expressconnectrouter.ResourceTRAssociation()
			attributes := map[string]string{
				"id":                    "ecr-unit:assoc-unit:tr-unit",
				"ecr_id":                "ecr-unit",
				"cen_id":                "cen-unit",
				"association_id":        "assoc-unit",
				"transit_router_id":     "tr-unit",
				"association_region_id": "cn-beijing",
			}
			if tc.createAttachment != "" {
				attributes["create_attachment"] = tc.createAttachment
			}
			state := &sdktf.InstanceState{ID: "ecr-unit:assoc-unit:tr-unit", Attributes: attributes}
			err := resourceSchema.Delete(resourceSchema.Data(state), unitTestClient(t, server.URL))
			if tc.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("delete error = %v, want %s", err, tc.wantError)
				}
			} else if err != nil {
				t.Fatalf("delete through the public resource entry failed: %s", err)
			}
			if len(deleteFlags) != tc.wantCalls {
				t.Fatalf("delete called %d times, want %d", len(deleteFlags), tc.wantCalls)
			}
			for i, flag := range deleteFlags {
				if flag != tc.wantFlag {
					t.Fatalf("request %d DeleteAttachment = %q, want %q", i+1, flag, tc.wantFlag)
				}
			}
			if got := atomic.LoadInt32(&listCalls); got != 0 {
				t.Fatalf("pre-delete attachment lookups = %d, want 0", got)
			}
		})
	}
}

// TestUnitTRAssociationUpdateSharesOneDeadline proves the description update
// cannot restart the update timeout after the prefix update consumed the
// budget. With a restarted budget the description action would be sent; with
// one shared deadline the update must stop instead.
func TestUnitTRAssociationUpdateSharesOneDeadline(t *testing.T) {
	fastTRAssociationWaits(t)
	var clockMu sync.Mutex
	// Start from the real wall clock so the operation's real deadline context is
	// not already expired before the first write.
	now := time.Now()
	t.Cleanup(expressconnectrouter.SetTRAssociationNowForTest(func() time.Time {
		clockMu.Lock()
		defer clockMu.Unlock()
		return now
	}))

	var mu sync.Mutex
	prefixUpdates, descriptionUpdates := 0, 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "ModifyExpressConnectRouterAssociationAllowedPrefix":
			mu.Lock()
			prefixUpdates++
			mu.Unlock()
			clockMu.Lock()
			// Consume more than the whole update timeout in this one stage.
			now = now.Add(10 * time.Minute)
			clockMu.Unlock()
			_, _ = w.Write([]byte(`{"Success":true}`))
		case "ModifyExpressConnectRouterAssociation":
			mu.Lock()
			descriptionUpdates++
			mu.Unlock()
			_, _ = w.Write([]byte(`{"Success":true}`))
		case "DescribeExpressConnectRouterAssociation":
			_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE","AllowedPrefixes":[]}],"Success":true}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	resourceSchema, before, diff := buildTRAssociationUpdateCase(t, 5*time.Minute)
	if _, diags := resourceSchema.Apply(context.Background(), before, diff, unitTestClient(t, server.URL)); !diags.HasError() {
		t.Fatal("an exhausted update deadline must fail the update instead of restarting the timeout")
	}

	mu.Lock()
	defer mu.Unlock()
	if prefixUpdates != 1 {
		t.Fatalf("prefix update called %d times, want 1", prefixUpdates)
	}
	if descriptionUpdates != 0 {
		t.Fatalf("description update must not restart the update timeout, called %d times", descriptionUpdates)
	}
}

func TestUnitTRAssociationReadValidatesDescribePages(t *testing.T) {
	const originalID = "ecr-unit:assoc-unit:tr-unit"
	for _, tc := range []struct {
		name      string
		body      string
		wantID    string
		wantError string
	}{
		{name: "normal_empty_page", body: `{"AssociationList":[],"Success":true}`},
		{name: "deleted_parent_empty_response", body: `{}`},
		{name: "coded_not_found", body: `{"Code":"ResourceNotFound.EcrId","Message":"EcrId not found","Success":false}`},
		{name: "empty_code_not_found_message", body: `{"Code":"","Message":"NotFound","Success":false}`, wantID: originalID, wantError: "NotFound"},
		{name: "foreign_association_needs_no_other_identity", body: `{"AssociationList":[{"AssociationId":"assoc-other"}],"Success":true}`},
		{name: "missing_list", body: `{"Success":true}`, wantID: originalID, wantError: "AssociationList"},
		{name: "empty_business_code", body: `{"Code":"","Message":"empty business code","Success":false}`, wantID: originalID, wantError: "empty business code"},
		{name: "malformed_list", body: `{"AssociationList":{},"Success":true}`, wantID: originalID, wantError: "AssociationList"},
		{name: "malformed_element", body: `{"AssociationList":["unexpected"],"Success":true}`, wantID: originalID, wantError: "AssociationList"},
		{name: "empty_object", body: `{"AssociationList":[{}],"Success":true}`, wantID: originalID, wantError: "AssociationId"},
		{name: "partial_match_missing_transit_router", body: `{"AssociationList":[{"AssociationId":"assoc-unit","EcrId":"ecr-unit"}],"Success":true}`, wantID: originalID, wantError: "TransitRouterId"},
		{name: "partial_match_wrong_typed_ecr", body: `{"AssociationList":[{"AssociationId":"assoc-unit","EcrId":42}],"Success":true}`, wantID: originalID, wantError: "EcrId"},
		{name: "malformed_next_token", body: `{"AssociationList":[],"NextToken":42,"Success":true}`, wantID: originalID, wantError: "NextToken"},
		{name: "null_next_token", body: `{"AssociationList":[],"NextToken":null,"Success":true}`, wantID: originalID, wantError: "NextToken"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %s", err)
				}
				w.Header().Set("Content-Type", "application/json")
				if action := r.Form.Get("Action"); action != "DescribeExpressConnectRouterAssociation" {
					t.Errorf("unexpected action: %q", action)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
					return
				}
				_, _ = w.Write([]byte(tc.body))
			}))
			t.Cleanup(server.Close)

			resourceSchema := expressconnectrouter.ResourceTRAssociation()
			data := resourceSchema.Data(importedTRAssociationState())
			err := resourceSchema.Read(data, unitTestClient(t, server.URL))
			if tc.wantError == "" {
				if err != nil {
					t.Fatalf("valid terminal page must be treated as confirmed absence: %s", err)
				}
			} else {
				if err == nil {
					t.Fatal("a malformed Describe page must fail Read")
				}
				if !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("Read error %q must contain %q", err, tc.wantError)
				}
			}
			if data.Id() != tc.wantID {
				t.Fatalf("Read left resource ID %q, want %q", data.Id(), tc.wantID)
			}
		})
	}
}

func TestUnitTRAssociationDeleteRejectsMalformedRefresh(t *testing.T) {
	fastTRAssociationWaits(t)
	for _, tc := range []struct {
		name      string
		body      string
		wantError string
	}{
		{name: "malformed_list", body: `{"AssociationList":{},"Success":true}`, wantError: "AssociationList"},
		{name: "empty_business_code", body: `{"Code":"","Message":"empty business code","Success":false}`, wantError: "empty business code"},
		{name: "empty_code_not_found_message", body: `{"Code":"","Message":"NotFound","Success":false}`, wantError: "NotFound"},
		{name: "malformed_identity", body: `{"AssociationList":[{"AssociationId":"assoc-unit","EcrId":"ecr-unit"}],"Success":true}`, wantError: "TransitRouterId"},
		{name: "malformed_next_token", body: `{"AssociationList":[],"NextToken":42,"Success":true}`, wantError: "NextToken"},
		{name: "null_next_token", body: `{"AssociationList":[],"NextToken":null,"Success":true}`, wantError: "NextToken"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var deleteCalls int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %s", err)
				}
				w.Header().Set("Content-Type", "application/json")
				switch action := r.Form.Get("Action"); action {
				case "ListTransitRouterEcrAttachments":
					_, _ = w.Write([]byte(`{"TransitRouterAttachments":[],"Success":true}`))
				case "DeleteExpressConnectRouterAssociation":
					atomic.AddInt32(&deleteCalls, 1)
					_, _ = w.Write([]byte(`{"Success":true}`))
				case "DescribeExpressConnectRouterAssociation":
					_, _ = w.Write([]byte(tc.body))
				default:
					t.Errorf("unexpected action: %q", action)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
				}
			}))
			t.Cleanup(server.Close)

			resourceSchema := expressconnectrouter.ResourceTRAssociation()
			err := resourceSchema.Delete(resourceSchema.Data(importedTRAssociationState()), unitTestClient(t, server.URL))
			if err == nil {
				t.Fatal("a malformed refresh response must not complete delete")
			}
			if !strings.Contains(err.Error(), tc.wantError) {
				t.Fatalf("delete error %q must contain %q", err, tc.wantError)
			}
			if got := atomic.LoadInt32(&deleteCalls); got != 1 {
				t.Fatalf("DeleteExpressConnectRouterAssociation called %d times, want 1", got)
			}
		})
	}
}

// TestUnitTRAssociationDeleteTreatsNotFoundAsSuccess covers both idempotent
// NotFound branches directly: one call, nil result and no follow-up wait.
func TestUnitTRAssociationDeleteTreatsNotFoundAsSuccess(t *testing.T) {
	fastTRAssociationWaits(t)
	for _, errorCode := range []string{"ResourceNotFound.AssociationId", "ResourceNotFound.EcrId"} {
		t.Run(errorCode, func(t *testing.T) {
			var mu sync.Mutex
			deleteCalls := 0
			describeCalls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %s", err)
				}
				w.Header().Set("Content-Type", "application/json")
				switch action := r.Form.Get("Action"); action {
				case "ListTransitRouterEcrAttachments":
					_, _ = w.Write([]byte(`{"TransitRouterAttachments":[],"Success":true}`))
				case "DeleteExpressConnectRouterAssociation":
					mu.Lock()
					deleteCalls++
					mu.Unlock()
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"` + errorCode + `","Message":"gone","RequestId":"test"}`))
				case "DescribeExpressConnectRouterAssociation":
					mu.Lock()
					describeCalls++
					mu.Unlock()
					_, _ = w.Write([]byte(`{"AssociationList":[],"Success":true}`))
				default:
					t.Errorf("unexpected action: %q", action)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
				}
			}))
			t.Cleanup(server.Close)

			resourceSchema := expressconnectrouter.ResourceTRAssociation()
			state := importedTRAssociationState()
			if err := resourceSchema.Delete(resourceSchema.Data(state), unitTestClient(t, server.URL)); err != nil {
				t.Fatalf("%s must be treated as an already deleted association, got %s", errorCode, err)
			}
			mu.Lock()
			defer mu.Unlock()
			if deleteCalls != 1 {
				t.Fatalf("%s must not be retried, delete called %d times", errorCode, deleteCalls)
			}
			if describeCalls != 0 {
				t.Fatalf("%s must return immediately without a state wait, describe called %d times", errorCode, describeCalls)
			}
		})
	}
}

func assertNoRequestAfter(t *testing.T, times []time.Time, deadline time.Time) {
	t.Helper()
	for _, ts := range times {
		if ts.After(deadline) {
			t.Fatalf("request started after the operation deadline: %s > %s", ts, deadline)
		}
	}
}

// assertDeadlineFailure requires the operation to fail for a deadline/timeout
// reason, not merely for the late fixture error.
func assertDeadlineFailure(t *testing.T, diags diag.Diagnostics) {
	t.Helper()
	if !diags.HasError() {
		t.Fatal("a callback blocked past the deadline must not turn into a late success")
	}
	text := ""
	for _, d := range diags {
		text += d.Summary + " " + d.Detail + " "
	}
	lower := strings.ToLower(text)
	if !strings.Contains(lower, "deadline") && !strings.Contains(lower, "timed out") && !strings.Contains(lower, "timeout") {
		t.Fatalf("expected a deadline/timeout error, got %s", text)
	}
}

// TestUnitTRAssociationCreateCallbackStopsAtDeadline holds the write callback
// until after the deadline and releases it afterwards: the operation must fail
// at the deadline, must not accept the late result and must not write again.
func TestUnitTRAssociationCreateCallbackStopsAtDeadline(t *testing.T) {
	var mu sync.Mutex
	createCalls := 0
	var createTimes []time.Time
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "CreateExpressConnectRouterAssociation":
			mu.Lock()
			createCalls++
			createTimes = append(createTimes, time.Now())
			mu.Unlock()
			time.Sleep(600 * time.Millisecond)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"Conflict.Lock","Message":"retryable","RequestId":"test"}`))
		case "DescribeExpressConnectRouterAssociation":
			_, _ = w.Write([]byte(`{"AssociationList":[],"Success":true}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	resourceSchema, diff, client := prepareTRAssociationCreate(t, server.URL, 200*time.Millisecond)
	start := time.Now()
	_, diags := resourceSchema.Apply(context.Background(), nil, diff, client)
	assertDeadlineFailure(t, diags)
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("create must stop at the configured deadline, took %s", elapsed)
	}
	mu.Lock()
	defer mu.Unlock()
	if createCalls != 1 {
		t.Fatalf("no write may happen after the deadline, create called %d times", createCalls)
	}
	assertNoRequestAfter(t, createTimes, start.Add(200*time.Millisecond))
}

// TestUnitTRAssociationUpdateCallbackStopsAtDeadline covers the update write
// path: a blocked prefix callback must not let the description update run or
// accept a late response after the deadline.
func TestUnitTRAssociationUpdateCallbackStopsAtDeadline(t *testing.T) {
	var mu sync.Mutex
	prefixUpdates, descriptionUpdates := 0, 0
	var prefixTimes []time.Time
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "ModifyExpressConnectRouterAssociationAllowedPrefix":
			mu.Lock()
			prefixUpdates++
			prefixTimes = append(prefixTimes, time.Now())
			mu.Unlock()
			time.Sleep(600 * time.Millisecond)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"Conflict.Lock","Message":"retryable","RequestId":"test"}`))
		case "ModifyExpressConnectRouterAssociation":
			mu.Lock()
			descriptionUpdates++
			mu.Unlock()
			_, _ = w.Write([]byte(`{"Success":true}`))
		case "DescribeExpressConnectRouterAssociation":
			_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE","AllowedPrefixes":[]}],"Success":true}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	resourceSchema, before, diff := buildTRAssociationUpdateCase(t, 200*time.Millisecond)
	client := unitTestClient(t, server.URL)
	start := time.Now()
	_, diags := resourceSchema.Apply(context.Background(), before, diff, client)
	assertDeadlineFailure(t, diags)
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("update must stop at the configured deadline, took %s", elapsed)
	}
	mu.Lock()
	defer mu.Unlock()
	if prefixUpdates != 1 {
		t.Fatalf("no write may happen after the deadline, prefix update called %d times", prefixUpdates)
	}
	if descriptionUpdates != 0 {
		t.Fatalf("the description update must not run after the deadline, called %d times", descriptionUpdates)
	}
	assertNoRequestAfter(t, prefixTimes, start.Add(200*time.Millisecond))
}

// TestUnitTRAssociationDeleteCallbackStopsAtDeadline covers the delete write
// path: the blocked delete callback must stop at the deadline and the late
// response must not be used.
func TestUnitTRAssociationDeleteCallbackStopsAtDeadline(t *testing.T) {
	var mu sync.Mutex
	deleteCalls := 0
	var deleteTimes []time.Time
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "ListTransitRouterEcrAttachments":
			_, _ = w.Write([]byte(`{"TransitRouterAttachments":[],"Success":true}`))
		case "DeleteExpressConnectRouterAssociation":
			mu.Lock()
			deleteCalls++
			deleteTimes = append(deleteTimes, time.Now())
			mu.Unlock()
			time.Sleep(600 * time.Millisecond)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"Conflict.Lock","Message":"retryable","RequestId":"test"}`))
		case "DescribeExpressConnectRouterAssociation":
			_, _ = w.Write([]byte(`{"AssociationList":[],"Success":true}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	resourceSchema := expressconnectrouter.ResourceTRAssociation()
	state := importedTRAssociationState()
	diff := &sdktf.InstanceDiff{
		Destroy: true,
		Meta: map[string]interface{}{
			schema.TimeoutKey: map[string]interface{}{schema.TimeoutDelete: 200 * time.Millisecond},
		},
	}
	client := unitTestClient(t, server.URL)
	start := time.Now()
	_, diags := resourceSchema.Apply(context.Background(), state, diff, client)
	assertDeadlineFailure(t, diags)
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("delete must stop at the configured deadline, took %s", elapsed)
	}
	mu.Lock()
	defer mu.Unlock()
	if deleteCalls != 1 {
		t.Fatalf("no write may happen after the deadline, delete called %d times", deleteCalls)
	}
	assertNoRequestAfter(t, deleteTimes, start.Add(200*time.Millisecond))
}

// TestUnitTRAssociationDeleteRefreshStopsAtDeadline covers the refresh path:
// the refresh callback blocks past the deadline and the wait must stop instead
// of accepting the late response or starting another refresh.
func TestUnitTRAssociationDeleteRefreshStopsAtDeadline(t *testing.T) {
	// Shorten only the initial poll delay so Refresh really executes before the
	// deadline, then require the wait to stop at the deadline.
	t.Cleanup(expressconnectrouter.SetTRAssociationStateDelayForTest(50 * time.Millisecond))
	var mu sync.Mutex
	deleteCalls := 0
	describeCalls := 0
	var describeTimes []time.Time
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "ListTransitRouterEcrAttachments":
			_, _ = w.Write([]byte(`{"TransitRouterAttachments":[],"Success":true}`))
		case "DeleteExpressConnectRouterAssociation":
			mu.Lock()
			deleteCalls++
			mu.Unlock()
			_, _ = w.Write([]byte(`{"Success":true}`))
		case "DescribeExpressConnectRouterAssociation":
			mu.Lock()
			describeCalls++
			describeTimes = append(describeTimes, time.Now())
			mu.Unlock()
			time.Sleep(600 * time.Millisecond)
			_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"DELETING","AllowedPrefixes":[]}],"Success":true}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	resourceSchema := expressconnectrouter.ResourceTRAssociation()
	state := importedTRAssociationState()
	diff := &sdktf.InstanceDiff{
		Destroy: true,
		Meta: map[string]interface{}{
			schema.TimeoutKey: map[string]interface{}{schema.TimeoutDelete: 200 * time.Millisecond},
		},
	}
	client := unitTestClient(t, server.URL)
	start := time.Now()
	_, diags := resourceSchema.Apply(context.Background(), state, diff, client)
	assertDeadlineFailure(t, diags)
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("the state wait must stop at the configured deadline, took %s", elapsed)
	}
	mu.Lock()
	defer mu.Unlock()
	if deleteCalls != 1 {
		t.Fatalf("exactly one delete write is expected, got %d", deleteCalls)
	}
	if describeCalls != 1 {
		t.Fatalf("no refresh may start after the deadline, describe called %d times", describeCalls)
	}
	assertNoRequestAfter(t, describeTimes, start.Add(200*time.Millisecond))
}

// TestUnitTRAssociationDeleteRefreshTreatsMissingParentAsNotFound preserves the
// shared describe classification: a post-delete Describe reporting an already
// missing parent must finish the destroy instead of failing it.
func TestUnitTRAssociationDeleteRefreshTreatsMissingParentAsNotFound(t *testing.T) {
	var mu sync.Mutex
	deleteCalls := 0
	describeCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "ListTransitRouterEcrAttachments":
			_, _ = w.Write([]byte(`{"TransitRouterAttachments":[],"Success":true}`))
		case "DeleteExpressConnectRouterAssociation":
			mu.Lock()
			deleteCalls++
			mu.Unlock()
			_, _ = w.Write([]byte(`{"Success":true}`))
		case "DescribeExpressConnectRouterAssociation":
			mu.Lock()
			describeCalls++
			mu.Unlock()
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"ResourceNotFound.EcrId","Message":"EcrId not found","RequestId":"test"}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	resourceSchema := expressconnectrouter.ResourceTRAssociation()
	state := importedTRAssociationState()
	if err := resourceSchema.Delete(resourceSchema.Data(state), unitTestClient(t, server.URL)); err != nil {
		t.Fatalf("a missing parent after delete must be treated as NotFound, got %s", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if deleteCalls != 1 {
		t.Fatalf("exactly one delete write is expected, got %d", deleteCalls)
	}
	if describeCalls < 1 {
		t.Fatal("refresh must execute so the missing parent is classified as NotFound")
	}
}

func prepareTRAssociationCreate(t *testing.T, serverURL string, timeout time.Duration) (*schema.Resource, *sdktf.InstanceDiff, *connectivity.AliyunClient) {
	t.Helper()
	resourceSchema := expressconnectrouter.ResourceTRAssociation()
	attrs := map[string]cty.Value{}
	for name, attributeType := range resourceSchema.CoreConfigSchema().ImpliedType().AttributeTypes() {
		attrs[name] = cty.NullVal(attributeType)
	}
	attrs["ecr_id"] = cty.StringVal("ecr-unit")
	attrs["transit_router_id"] = cty.StringVal("tr-unit")
	attrs["association_region_id"] = cty.StringVal("cn-hangzhou")
	attrs["cen_id"] = cty.StringVal("cen-unit")
	attrs["transit_router_owner_id"] = cty.NumberIntVal(123456)
	attrs["create_attachment"] = cty.True
	rawConfig := cty.ObjectVal(attrs)

	diff, err := resourceSchema.Diff(context.Background(), nil, sdktf.NewResourceConfigShimmed(rawConfig, resourceSchema.CoreConfigSchema()), nil)
	if err != nil {
		t.Fatalf("plan through the public resource entry failed: %s", err)
	}
	diff.RawConfig = rawConfig
	diff.Meta = map[string]interface{}{
		schema.TimeoutKey: map[string]interface{}{
			schema.TimeoutCreate: timeout,
			schema.TimeoutUpdate: 5 * time.Minute,
			schema.TimeoutDelete: 5 * time.Minute,
		},
	}
	return resourceSchema, diff, unitTestClient(t, serverURL)
}

func applyTRAssociationCreateWithTimeout(t *testing.T, serverURL string, timeout time.Duration) diag.Diagnostics {
	resourceSchema, diff, client := prepareTRAssociationCreate(t, serverURL, timeout)
	_, diags := resourceSchema.Apply(context.Background(), nil, diff, client)
	return diags
}

func TestUnitTRAssociationDescriptionValidation(t *testing.T) {
	for _, tc := range []struct {
		name       string
		value      string
		wantErrors int
	}{
		{name: "empty", value: ""},
		{name: "multiline", value: "line-1\nline-2"},
		{name: "max_length", value: strings.Repeat("a", 128)},
		{name: "too_long", value: strings.Repeat("a", 129), wantErrors: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := expressconnectrouter.ResourceTRAssociation().Schema["description"].ValidateFunc(tc.value, "description")
			if len(errors) != tc.wantErrors {
				t.Fatalf("description errors = %v, want %d", errors, tc.wantErrors)
			}
		})
	}
}

func TestUnitTRAssociationDeleteRejectsMalformedID(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusBadRequest)
	}))
	t.Cleanup(server.Close)

	resourceSchema := expressconnectrouter.ResourceTRAssociation()
	state := importedTRAssociationState()
	state.ID = "ecr-unit:assoc-unit"
	state.Attributes["id"] = state.ID
	err := resourceSchema.Delete(resourceSchema.Data(state), unitTestClient(t, server.URL))
	if err == nil {
		t.Fatal("a malformed resource ID must return an error instead of panicking")
	}
	if !strings.Contains(err.Error(), "Expected parts' length 3, got 2") {
		t.Fatalf("unexpected malformed-ID error: %s", err)
	}
	if got := atomic.LoadInt32(&calls); got != 0 {
		t.Fatalf("malformed ID issued %d cloud requests, want 0", got)
	}
}

func TestUnitTRAssociationUpdateClearingPrefixesOmitsMode(t *testing.T) {
	for _, tc := range []struct {
		name       string
		prefixes   cty.Value
		configMode cty.Value
	}{
		{name: "explicit_empty_with_computed_mode", prefixes: cty.ListValEmpty(cty.String), configMode: cty.NullVal(cty.String)},
		{name: "removed_with_computed_mode", prefixes: cty.NullVal(cty.List(cty.String)), configMode: cty.NullVal(cty.String)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fastTRAssociationWaits(t)
			var forms []url.Values
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %s", err)
				}
				w.Header().Set("Content-Type", "application/json")
				switch action := r.Form.Get("Action"); action {
				case "ModifyExpressConnectRouterAssociationAllowedPrefix":
					forms = append(forms, r.Form)
					_, _ = w.Write([]byte(`{"Success":true}`))
				case "DescribeExpressConnectRouterAssociation":
					_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","RegionId":"cn-beijing","Status":"ACTIVE","AllowedPrefixes":[],"AllowedPrefixesMode":"MatchMode","Description":"old-description"}],"Success":true}`))
				default:
					t.Errorf("unexpected action: %q", action)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
				}
			}))
			t.Cleanup(server.Close)

			resourceSchema, before, diff := buildTRAssociationUpdateCaseWithConfig(t, 5*time.Minute, tc.prefixes, tc.configMode, "old-description")
			updated, diags := resourceSchema.Apply(context.Background(), before, diff, unitTestClient(t, server.URL))
			if diags.HasError() {
				t.Fatalf("clearing prefixes failed: %v", diags)
			}
			if len(forms) != 1 {
				t.Fatalf("ModifyExpressConnectRouterAssociationAllowedPrefix calls = %d, want 1", len(forms))
			}
			if values, present := forms[0]["AllowedPrefixesMode"]; present {
				t.Fatalf("AllowedPrefixesMode must be omitted from a clear request, got %v", values)
			}
			assertFormArray(t, forms[0], "AllowedPrefixes")
			if got := updated.Attributes["allowed_prefixes.#"]; got != "0" {
				t.Fatalf("updated prefix count = %q, want 0", got)
			}
		})
	}
}

func TestUnitTRAssociationPrefixOrderSuppress(t *testing.T) {
	suppressTRAssociationPrefixReorder := expressconnectrouter.ResourceTRAssociation().Schema["allowed_prefixes"].DiffSuppressFunc
	t.Run("NilResourceData", func(t *testing.T) {
		if suppressTRAssociationPrefixReorder("allowed_prefixes.#", "1", "1", nil) {
			t.Fatal("a nil ResourceData carries no old/new pair and must not suppress")
		}
	})

	for _, tc := range []struct {
		name     string
		oldValue string
		newValue string
		expected bool
	}{
		{name: "ReorderOnly", oldValue: "10.0.1.0/24,10.0.2.0/24", newValue: "10.0.2.0/24,10.0.1.0/24", expected: true},
		{name: "LongerReorder", oldValue: "10.0.1.0/24,10.0.2.0/24,10.0.3.0/24", newValue: "10.0.3.0/24,10.0.1.0/24,10.0.2.0/24", expected: true},
		{name: "MemberAdded", oldValue: "10.0.1.0/24,10.0.2.0/24", newValue: "10.0.1.0/24,10.0.2.0/24,10.0.3.0/24", expected: false},
		{name: "MemberRemoved", oldValue: "10.0.1.0/24,10.0.2.0/24", newValue: "10.0.1.0/24", expected: false},
		{name: "MemberReplaced", oldValue: "10.0.1.0/24,10.0.2.0/24", newValue: "10.0.1.0/24,10.0.9.0/24", expected: false},
		{name: "MultiplicityChanged", oldValue: "10.0.1.0/24,10.0.1.0/24,10.0.2.0/24", newValue: "10.0.1.0/24,10.0.2.0/24,10.0.2.0/24", expected: false},
		{name: "OldEmpty", newValue: "10.0.1.0/24", expected: false},
		{name: "ExplicitClear", oldValue: "10.0.1.0/24,10.0.2.0/24", expected: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			oldPrefixes := []string{}
			if tc.oldValue != "" {
				oldPrefixes = strings.Split(tc.oldValue, ",")
			}
			newPrefixes := []string{}
			if tc.newValue != "" {
				newPrefixes = strings.Split(tc.newValue, ",")
			}
			attributes := map[string]string{
				"id":                      "ecr-unit:assoc-unit:tr-unit",
				"ecr_id":                  "ecr-unit",
				"cen_id":                  "cen-unit",
				"association_id":          "assoc-unit",
				"transit_router_id":       "tr-unit",
				"association_region_id":   "cn-beijing",
				"transit_router_owner_id": "123456",
				"create_attachment":       "true",
				"allowed_prefixes.#":      fmt.Sprint(len(oldPrefixes)),
				"allowed_prefixes_mode":   "MatchMode",
			}
			for i, prefix := range oldPrefixes {
				attributes[fmt.Sprintf("allowed_prefixes.%d", i)] = prefix
			}
			state := &sdktf.InstanceState{ID: "ecr-unit:assoc-unit:tr-unit", Attributes: attributes}
			newPrefixValues := make([]cty.Value, 0, len(newPrefixes))
			for _, prefix := range newPrefixes {
				newPrefixValues = append(newPrefixValues, cty.StringVal(prefix))
			}
			newPrefixList := cty.ListValEmpty(cty.String)
			if len(newPrefixValues) > 0 {
				newPrefixList = cty.ListVal(newPrefixValues)
			}
			resourceSchema := expressconnectrouter.ResourceTRAssociation()
			newMode := cty.StringVal("MatchMode")
			if tc.name == "ExplicitClear" {
				newMode = cty.NullVal(cty.String)
			}
			rawConfig := trAssociationPrefixModeRawConfig(resourceSchema, newPrefixList, newMode, true)
			state.RawConfig = rawConfig
			var data *schema.ResourceData
			resourceSchema.Schema["allowed_prefixes"].DiffSuppressFunc = func(key, old, new string, d *schema.ResourceData) bool {
				data = d
				return suppressTRAssociationPrefixReorder(key, old, new, d)
			}
			diff, err := resourceSchema.Diff(context.Background(), state, sdktf.NewResourceConfigShimmed(rawConfig, resourceSchema.CoreConfigSchema()), nil)
			if err != nil {
				t.Fatalf("build prefix diff: %v", err)
			}
			if data == nil {
				t.Fatal("prefix diff did not invoke the suppressor")
			}
			changed := diff != nil && !diff.Empty()
			if changed == tc.expected {
				t.Fatalf("changed = %t, want %t for old=%v new=%v", changed, !tc.expected, tc.oldValue, tc.newValue)
			}
			if got := suppressTRAssociationPrefixReorder("allowed_prefixes", tc.oldValue, tc.newValue, data); got != tc.expected {
				t.Fatalf("case %s: suppressor = %v, want %v", tc.name, got, tc.expected)
			}
		})
	}
}

func trAssociationPrefixModeRawConfig(resourceSchema *schema.Resource, prefixes, mode cty.Value, existing bool) cty.Value {
	attrs := map[string]cty.Value{}
	for name, attributeType := range resourceSchema.CoreConfigSchema().ImpliedType().AttributeTypes() {
		attrs[name] = cty.NullVal(attributeType)
	}
	attrs["ecr_id"] = cty.StringVal("ecr-unit")
	attrs["transit_router_id"] = cty.StringVal("tr-unit")
	attrs["association_region_id"] = cty.StringVal("cn-beijing")
	attrs["allowed_prefixes"] = prefixes
	attrs["allowed_prefixes_mode"] = mode
	if existing {
		attrs["cen_id"] = cty.StringVal("cen-unit")
		attrs["transit_router_owner_id"] = cty.NumberIntVal(123456)
		attrs["create_attachment"] = cty.True
	}
	return cty.ObjectVal(attrs)
}

func trAssociationPrefixModeState(existing bool) *sdktf.InstanceState {
	if !existing {
		return &sdktf.InstanceState{}
	}
	return &sdktf.InstanceState{
		ID: "ecr-unit:assoc-unit:tr-unit",
		Attributes: map[string]string{
			"id":                      "ecr-unit:assoc-unit:tr-unit",
			"ecr_id":                  "ecr-unit",
			"cen_id":                  "cen-unit",
			"association_id":          "assoc-unit",
			"transit_router_id":       "tr-unit",
			"association_region_id":   "cn-beijing",
			"transit_router_owner_id": "123456",
			"create_attachment":       "true",
			"allowed_prefixes.#":      "1",
			"allowed_prefixes.0":      "10.0.0.0/24",
			"allowed_prefixes_mode":   "MatchMode",
			"status":                  "ACTIVE",
		},
	}
}

func TestUnitTRAssociationPrefixModeConfigValidation(t *testing.T) {
	for _, tc := range []struct {
		name      string
		prefixes  cty.Value
		mode      cty.Value
		existing  bool
		wantError bool
	}{
		{name: "mode_with_prefix_is_valid", prefixes: cty.ListVal([]cty.Value{cty.StringVal("10.0.0.0/24")}), mode: cty.StringVal("MatchMode")},
		{name: "create_mode_without_prefixes_is_rejected", prefixes: cty.NullVal(cty.List(cty.String)), mode: cty.StringVal("MatchMode"), wantError: true},
		{name: "create_mode_with_empty_prefixes_is_rejected", prefixes: cty.ListValEmpty(cty.String), mode: cty.StringVal("MatchMode"), wantError: true},
		{name: "empty_prefixes_without_mode_is_valid", prefixes: cty.ListValEmpty(cty.String), mode: cty.NullVal(cty.String)},
		{name: "unknown_prefixes_defer_validation", prefixes: cty.UnknownVal(cty.List(cty.String)), mode: cty.StringVal("IncrementalMode")},
		{name: "unknown_mode_defers_validation", prefixes: cty.NullVal(cty.List(cty.String)), mode: cty.UnknownVal(cty.String)},
		{name: "update_mode_only_is_rejected", prefixes: cty.NullVal(cty.List(cty.String)), mode: cty.StringVal("IncrementalMode"), existing: true, wantError: true},
		{name: "update_explicit_empty_with_mode_is_rejected", prefixes: cty.ListValEmpty(cty.String), mode: cty.StringVal("MatchMode"), existing: true, wantError: true},
		{name: "update_clear_without_mode_is_valid", prefixes: cty.ListValEmpty(cty.String), mode: cty.NullVal(cty.String), existing: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resourceSchema := expressconnectrouter.ResourceTRAssociation()
			rawConfig := trAssociationPrefixModeRawConfig(resourceSchema, tc.prefixes, tc.mode, tc.existing)
			state := trAssociationPrefixModeState(tc.existing)
			state.RawConfig = rawConfig
			_, err := resourceSchema.Diff(context.Background(), state, sdktf.NewResourceConfigShimmed(rawConfig, resourceSchema.CoreConfigSchema()), nil)
			if tc.wantError {
				if err == nil || !strings.Contains(err.Error(), "allowed_prefixes_mode") || !strings.Contains(err.Error(), "non-empty") || !strings.Contains(err.Error(), "allowed_prefixes") {
					t.Fatalf("validation error = %v, want the prefix combination error", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestUnitTRAssociationResolvedUnknownPrefixModeValidation(t *testing.T) {
	for _, existing := range []bool{false, true} {
		for _, unknown := range []string{"prefixes", "mode"} {
			name := "create/" + unknown
			if existing {
				name = "update/" + unknown
			}
			t.Run(name, func(t *testing.T) {
				planPrefixes := cty.ListValEmpty(cty.String)
				planMode := cty.StringVal("MatchMode")
				if unknown == "prefixes" {
					planPrefixes = cty.UnknownVal(cty.List(cty.String))
				} else {
					planMode = cty.UnknownVal(cty.String)
				}
				resourceSchema := expressconnectrouter.ResourceTRAssociation()
				planRaw := trAssociationPrefixModeRawConfig(resourceSchema, planPrefixes, planMode, existing)
				resolvedRaw := trAssociationPrefixModeRawConfig(resourceSchema, cty.ListValEmpty(cty.String), cty.StringVal("MatchMode"), existing)
				state := trAssociationPrefixModeState(existing)
				state.RawConfig = planRaw
				if _, err := resourceSchema.Diff(context.Background(), state, sdktf.NewResourceConfigShimmed(planRaw, resourceSchema.CoreConfigSchema()), nil); err != nil {
					t.Fatalf("unknown plan values must defer validation: %v", err)
				}
				state.RawConfig = planRaw
				diff, err := resourceSchema.Diff(context.Background(), state, sdktf.NewResourceConfigShimmed(resolvedRaw, resourceSchema.CoreConfigSchema()), nil)
				if err != nil {
					t.Fatalf("build resolved apply diff with plan RawConfig: %v", err)
				}
				if diff == nil {
					t.Fatal("expected a resolved apply diff")
				}
				diff.RawConfig = planRaw
				diff.Meta = map[string]interface{}{
					schema.TimeoutKey: map[string]interface{}{
						schema.TimeoutCreate: 5 * time.Minute,
						schema.TimeoutUpdate: 5 * time.Minute,
						schema.TimeoutDelete: 5 * time.Minute,
					},
				}
				var requests int32
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					atomic.AddInt32(&requests, 1)
					t.Errorf("unexpected cloud request for resolved invalid prefix configuration")
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidRequest","Message":"unexpected request"}`))
				}))
				t.Cleanup(server.Close)
				_, diags := resourceSchema.Apply(context.Background(), state, diff, unitTestClient(t, server.URL))
				if !diags.HasError() || !strings.Contains(fmt.Sprintf("%v", diags), "allowed_prefixes_mode") || !strings.Contains(fmt.Sprintf("%v", diags), "non-empty") || !strings.Contains(fmt.Sprintf("%v", diags), "allowed_prefixes") {
					t.Fatalf("resolved apply error = %v, want the prefix combination error", diags)
				}
				if got := atomic.LoadInt32(&requests); got != 0 {
					t.Fatalf("resolved invalid apply issued %d cloud requests, want 0", got)
				}
			})
		}
	}
}

func TestUnitTRAssociationUpdatePairsPrefixesAndMode(t *testing.T) {
	for _, tc := range []struct {
		name     string
		prefixes []string
		mode     string
	}{
		{name: "mode_only_echoes_stored_prefixes", prefixes: []string{"10.0.0.0/24"}, mode: "IncrementalMode"},
		{name: "prefix_only_echoes_stored_mode", prefixes: []string{"10.9.0.0/24"}, mode: "MatchMode"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fastTRAssociationWaits(t)
			var forms []url.Values
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %s", err)
				}
				w.Header().Set("Content-Type", "application/json")
				switch action := r.Form.Get("Action"); action {
				case "ModifyExpressConnectRouterAssociationAllowedPrefix":
					forms = append(forms, r.Form)
					_, _ = w.Write([]byte(`{"Success":true}`))
				case "DescribeExpressConnectRouterAssociation":
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"Success": true,
						"AssociationList": []map[string]interface{}{{
							"AssociationId":       "assoc-unit",
							"AssociationNodeType": "TR",
							"EcrId":               "ecr-unit",
							"TransitRouterId":     "tr-unit",
							"RegionId":            "cn-beijing",
							"Status":              "ACTIVE",
							"AllowedPrefixes":     tc.prefixes,
							"AllowedPrefixesMode": tc.mode,
							"Description":         "old-description",
						}},
					})
				default:
					t.Errorf("unexpected action: %q", action)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
				}
			}))
			t.Cleanup(server.Close)

			resourceSchema, before, diff := buildTRAssociationUpdateCaseWithValues(t, 5*time.Minute, tc.prefixes, tc.mode, "old-description")
			_, diags := resourceSchema.Apply(context.Background(), before, diff, unitTestClient(t, server.URL))
			if diags.HasError() {
				t.Fatalf("paired prefix update failed: %v", diags)
			}
			if len(forms) != 1 {
				t.Fatalf("ModifyExpressConnectRouterAssociationAllowedPrefix calls = %d, want 1", len(forms))
			}
			if got := forms[0].Get("AllowedPrefixesMode"); got != tc.mode {
				t.Errorf("AllowedPrefixesMode = %q, want %q", got, tc.mode)
			}
			assertFormArray(t, forms[0], "AllowedPrefixes", tc.prefixes...)
		})
	}
}

func buildTRAssociationUpdateCase(t *testing.T, updateTimeout time.Duration) (*schema.Resource, *sdktf.InstanceState, *sdktf.InstanceDiff) {
	t.Helper()
	return buildTRAssociationUpdateCaseWithValues(t, updateTimeout, []string{"10.9.0.0/24"}, "MatchMode", "new-description")
}

func buildTRAssociationUpdateCaseWithValues(t *testing.T, updateTimeout time.Duration, prefixes []string, mode, description string) (*schema.Resource, *sdktf.InstanceState, *sdktf.InstanceDiff) {
	t.Helper()
	prefixValues := make([]cty.Value, len(prefixes))
	for i, prefix := range prefixes {
		prefixValues[i] = cty.StringVal(prefix)
	}
	return buildTRAssociationUpdateCaseWithConfig(t, updateTimeout, cty.ListVal(prefixValues), cty.StringVal(mode), description)
}

func buildTRAssociationUpdateCaseWithConfig(t *testing.T, updateTimeout time.Duration, prefixes, mode cty.Value, description string) (*schema.Resource, *sdktf.InstanceState, *sdktf.InstanceDiff) {
	t.Helper()
	resourceSchema := expressconnectrouter.ResourceTRAssociation()
	before := &sdktf.InstanceState{
		ID: "ecr-unit:assoc-unit:tr-unit",
		Attributes: map[string]string{
			"id":                      "ecr-unit:assoc-unit:tr-unit",
			"ecr_id":                  "ecr-unit",
			"cen_id":                  "cen-unit",
			"association_id":          "assoc-unit",
			"transit_router_id":       "tr-unit",
			"association_region_id":   "cn-beijing",
			"transit_router_owner_id": "123456",
			"create_attachment":       "true",
			"allowed_prefixes.#":      "1",
			"allowed_prefixes.0":      "10.0.0.0/24",
			"allowed_prefixes_mode":   "MatchMode",
			"description":             "old-description",
			"status":                  "ACTIVE",
		},
	}
	attrs := map[string]cty.Value{}
	for name, attributeType := range resourceSchema.CoreConfigSchema().ImpliedType().AttributeTypes() {
		attrs[name] = cty.NullVal(attributeType)
	}
	attrs["ecr_id"] = cty.StringVal("ecr-unit")
	attrs["transit_router_id"] = cty.StringVal("tr-unit")
	attrs["association_region_id"] = cty.StringVal("cn-beijing")
	attrs["cen_id"] = cty.StringVal("cen-unit")
	attrs["transit_router_owner_id"] = cty.NumberIntVal(123456)
	attrs["create_attachment"] = cty.True
	attrs["allowed_prefixes"] = prefixes
	attrs["allowed_prefixes_mode"] = mode
	attrs["description"] = cty.StringVal(description)
	rawConfig := cty.ObjectVal(attrs)

	diff, err := resourceSchema.Diff(context.Background(), before, sdktf.NewResourceConfigShimmed(rawConfig, resourceSchema.CoreConfigSchema()), nil)
	if err != nil {
		t.Fatalf("plan through the public resource entry failed: %s", err)
	}
	diff.RawConfig = rawConfig
	diff.Meta = map[string]interface{}{
		schema.TimeoutKey: map[string]interface{}{
			schema.TimeoutCreate: 5 * time.Minute,
			schema.TimeoutUpdate: updateTimeout,
			schema.TimeoutDelete: 5 * time.Minute,
		},
	}
	return resourceSchema, before, diff
}

func importedTRAssociationState() *sdktf.InstanceState {
	return &sdktf.InstanceState{
		ID: "ecr-unit:assoc-unit:tr-unit",
		Attributes: map[string]string{
			"id":                    "ecr-unit:assoc-unit:tr-unit",
			"ecr_id":                "ecr-unit",
			"cen_id":                "cen-unit",
			"association_id":        "assoc-unit",
			"transit_router_id":     "tr-unit",
			"association_region_id": "cn-beijing",
		},
	}
}

func unitTestClient(t *testing.T, serverURL string) *connectivity.AliyunClient {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoint := strings.TrimPrefix(serverURL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := connectivity.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("expressconnectrouter", endpoint)
	endpoints.Store("cbn", endpoint)
	return client
}

// TestUnitTRAssociationCreateFinalReadStopsAtDeadline covers the case where the create
// write and its state wait succeed inside the budget, but the final Describe that
// fills state blocks. The whole operation must still fail at the configured
// deadline instead of running past it on the shared context-free describe path.
func TestUnitTRAssociationCreateFinalReadStopsAtDeadline(t *testing.T) {
	t.Cleanup(expressconnectrouter.SetTRAssociationStateDelayForTest(10 * time.Millisecond))
	var mu sync.Mutex
	createCalls := 0
	var describeTimes []time.Time
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "CreateExpressConnectRouterAssociation":
			mu.Lock()
			createCalls++
			mu.Unlock()
			_, _ = w.Write([]byte(`{"AssociationId":"assoc-unit","Success":true}`))
		case "DescribeExpressConnectRouterAssociation":
			mu.Lock()
			describeTimes = append(describeTimes, time.Now())
			call := len(describeTimes)
			mu.Unlock()
			if call == 1 {
				// The state wait's first refresh returns ACTIVE immediately so the
				// write and wait finish inside the budget.
				_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE","AllowedPrefixes":[]}],"Success":true}`))
				return
			}
			// The final state-filling Describe blocks past the deadline.
			time.Sleep(600 * time.Millisecond)
			_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE","AllowedPrefixes":[]}],"Success":true}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	resourceSchema, diff, client := prepareTRAssociationCreate(t, server.URL, 300*time.Millisecond)
	start := time.Now()
	_, diags := resourceSchema.Apply(context.Background(), nil, diff, client)
	assertDeadlineFailure(t, diags)
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("the final read must stop at the configured deadline, took %s", elapsed)
	}
	mu.Lock()
	defer mu.Unlock()
	if createCalls != 1 {
		t.Fatalf("exactly one create write is expected, got %d", createCalls)
	}
	assertNoRequestAfter(t, describeTimes, start.Add(300*time.Millisecond))
}

// TestUnitTRAssociationUpdateFinalReadStopsAtDeadline covers the update path:
// the description write and its wait succeed, then the final Describe
// blocks. No further write may happen and the operation must stop at the deadline.
func TestUnitTRAssociationUpdateFinalReadStopsAtDeadline(t *testing.T) {
	t.Cleanup(expressconnectrouter.SetTRAssociationStateDelayForTest(10 * time.Millisecond))
	var mu sync.Mutex
	descriptionUpdates := 0
	var describeTimes []time.Time
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "ModifyExpressConnectRouterAssociationAllowedPrefix":
			_, _ = w.Write([]byte(`{"Success":true}`))
		case "ModifyExpressConnectRouterAssociation":
			mu.Lock()
			descriptionUpdates++
			mu.Unlock()
			_, _ = w.Write([]byte(`{"Success":true}`))
		case "DescribeExpressConnectRouterAssociation":
			mu.Lock()
			describeTimes = append(describeTimes, time.Now())
			call := len(describeTimes)
			mu.Unlock()
			if call <= 2 {
				// Each write's state wait returns ACTIVE immediately so both writes
				// and their waits finish inside the budget.
				_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE","AllowedPrefixes":[]}],"Success":true}`))
				return
			}
			time.Sleep(600 * time.Millisecond)
			_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE","AllowedPrefixes":[]}],"Success":true}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	resourceSchema, before, diff := buildTRAssociationUpdateCase(t, 300*time.Millisecond)
	client := unitTestClient(t, server.URL)
	start := time.Now()
	_, diags := resourceSchema.Apply(context.Background(), before, diff, client)
	assertDeadlineFailure(t, diags)
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("the final read must stop at the configured deadline, took %s", elapsed)
	}
	mu.Lock()
	defer mu.Unlock()
	if descriptionUpdates != 1 {
		t.Fatalf("exactly one description write is expected, got %d", descriptionUpdates)
	}
	assertNoRequestAfter(t, describeTimes, start.Add(300*time.Millisecond))
}

// TestUnitTRAssociationRPCPostPreCancelledContextStartsNoRequest verifies that
// the RPC wrapper itself rejects an already-cancelled context before starting
// a goroutine that can send a new request.
func TestUnitTRAssociationRPCPostPreCancelledContextStartsNoRequest(t *testing.T) {
	requestStarted := make(chan struct{}, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case requestStarted <- struct{}{}:
		default:
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"Success":true}`))
	}))
	t.Cleanup(server.Close)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	response, err := expressconnectrouter.TRAssociationRPCPostForTest(
		ctx, unitTestClient(t, server.URL), "ExpressConnectRouter", "2023-09-01",
		"CreateExpressConnectRouterAssociation", map[string]interface{}{}, map[string]interface{}{},
	)
	if err != context.Canceled {
		t.Fatalf("the wrapper must return context.Canceled, got %v", err)
	}
	if response != nil {
		t.Fatalf("a cancelled call must not return a response, got %#v", response)
	}
	select {
	case <-requestStarted:
		t.Fatal("no request may start with an already-cancelled context")
	case <-time.After(time.Second):
	}
}

func TestUnitTRAssociationRPCPostCancellationStopsInFlightRequest(t *testing.T) {
	arrived := make(chan struct{})
	transportCancelled := make(chan struct{})
	release := make(chan struct{})
	var arrivedOnce sync.Once
	var cancelledOnce sync.Once
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		arrivedOnce.Do(func() { close(arrived) })
		select {
		case <-r.Context().Done():
			cancelledOnce.Do(func() { close(transportCancelled) })
		case <-release:
		case <-time.After(3 * time.Second):
		}
	}))
	defer server.Close()
	defer close(release)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-arrived
		cancel()
	}()

	start := time.Now()
	response, err := expressconnectrouter.TRAssociationRPCPostForTest(
		ctx, unitTestClient(t, server.URL), "ExpressConnectRouter", "2023-09-01",
		"CreateExpressConnectRouterAssociation", map[string]interface{}{}, map[string]interface{}{},
	)
	if err == nil {
		t.Fatal("expected the in-flight request to fail after cancellation")
	}
	if response != nil {
		t.Fatalf("a cancelled call must not return a response, got %#v", response)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("the cancelled request did not stop promptly: %s", elapsed)
	}
	select {
	case <-transportCancelled:
	case <-time.After(time.Second):
		t.Fatal("the server did not observe cancellation of the request transport")
	}
	time.Sleep(100 * time.Millisecond)
	if got := requests.Load(); got != 1 {
		t.Fatalf("a cancelled request must not be retried, got %d requests", got)
	}
}

func TestUnitTRAssociationCreateExhaustedBudgetStartsNoRequest(t *testing.T) {
	var mu sync.Mutex
	totalRequests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		totalRequests++
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"Success":true}`))
	}))
	t.Cleanup(server.Close)

	now := time.Now()
	clockCalls := 0
	t.Cleanup(expressconnectrouter.SetTRAssociationNowForTest(func() time.Time {
		clockCalls++
		if clockCalls == 1 {
			return now
		}
		return now.Add(10 * time.Minute)
	}))
	resourceSchema, diff, client := prepareTRAssociationCreate(t, server.URL, time.Millisecond)
	_, diags := resourceSchema.Apply(context.Background(), nil, diff, client)
	if !diags.HasError() {
		t.Fatal("an exhausted create budget must fail")
	}
	mu.Lock()
	defer mu.Unlock()
	if totalRequests != 0 {
		t.Fatalf("no request may start with an exhausted create budget, got %d", totalRequests)
	}
}

// TestUnitTRAssociationStateRefreshDiscriminatesObjects verifies that the
// association refresh skips a wrong first entry and only accepts the object
// whose AssociationId/EcrId/TransitRouterId/AssociationNodeType all match.
func TestUnitTRAssociationStateRefreshDiscriminatesObjects(t *testing.T) {
	t.Cleanup(expressconnectrouter.SetTRAssociationStateDelayForTest(10 * time.Millisecond))
	var mu sync.Mutex
	deleteCalls := 0
	describeCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "ListTransitRouterEcrAttachments":
			_, _ = w.Write([]byte(`{"TransitRouterAttachments":[],"Success":true}`))
		case "DeleteExpressConnectRouterAssociation":
			mu.Lock()
			deleteCalls++
			mu.Unlock()
			_, _ = w.Write([]byte(`{"Success":true}`))
		case "DescribeExpressConnectRouterAssociation":
			mu.Lock()
			describeCalls++
			mu.Unlock()
			// The first entry is a different association that is still ACTIVE; the
			// target association is already gone. Taking the first entry would
			// wrongly report the target as alive and block the destroy.
			_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-other","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE","AllowedPrefixes":[]}],"Success":true}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	resourceSchema := expressconnectrouter.ResourceTRAssociation()
	state := importedTRAssociationState()
	if err := resourceSchema.Delete(resourceSchema.Data(state), unitTestClient(t, server.URL)); err != nil {
		t.Fatalf("a response containing only other objects must be treated as target-absent, got %s", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if deleteCalls != 1 {
		t.Fatalf("exactly one delete write is expected, got %d", deleteCalls)
	}
	if describeCalls < 1 {
		t.Fatal("refresh must execute so the wrong-first-entry response is evaluated")
	}
}

// TestUnitTRAssociationStateRefreshFindsTargetAfterWrongFirstEntry verifies that
// a wrong first object followed by the target returns the target's own status,
// not the first entry's.
func TestUnitTRAssociationStateRefreshFindsTargetAfterWrongFirstEntry(t *testing.T) {
	t.Cleanup(expressconnectrouter.SetTRAssociationStateDelayForTest(10 * time.Millisecond))
	var mu sync.Mutex
	createCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "CreateExpressConnectRouterAssociation":
			mu.Lock()
			createCalls++
			mu.Unlock()
			_, _ = w.Write([]byte(`{"AssociationId":"assoc-unit","Success":true}`))
		case "DescribeExpressConnectRouterAssociation":
			// Wrong object first with a non-target status; the target second is
			// ACTIVE so the create wait can finish and the final read fills state.
			_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-other","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"CREATING","AllowedPrefixes":[]},{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE","AllowedPrefixes":[],"Description":"after-fix"}],"Success":true}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	resourceSchema, diff, client := prepareTRAssociationCreate(t, server.URL, 5*time.Second)
	newState, diags := resourceSchema.Apply(context.Background(), nil, diff, client)
	if diags.HasError() {
		t.Fatalf("create must succeed when the target is present after a wrong first entry: %v", diags)
	}
	if newState == nil || newState.ID != "ecr-unit:assoc-unit:tr-unit" {
		t.Fatalf("the target association must be stored, got %+v", newState)
	}
	if got := newState.Attributes["description"]; got != "after-fix" {
		t.Fatalf("state must come from the target entry, description=%q", got)
	}
	if got := newState.Attributes["status"]; got != "ACTIVE" {
		t.Fatalf("the target status must be read, not the first entry's, status=%q", got)
	}
	mu.Lock()
	defer mu.Unlock()
	if createCalls != 1 {
		t.Fatalf("exactly one create write is expected, got %d", createCalls)
	}
}

func TestUnitTRAssociationCreateDescriptionSerialization(t *testing.T) {
	tests := []struct {
		name        string
		description cty.Value
		wantPresent bool
		wantValue   string
	}{
		{name: "omitted", description: cty.NullVal(cty.String)},
		{name: "explicit empty", description: cty.StringVal(""), wantPresent: true},
		{name: "nonempty", description: cty.StringVal("configured-description"), wantPresent: true, wantValue: "configured-description"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fastTRAssociationWaits(t)
			var mu sync.Mutex
			createCalls := 0
			var descriptionValues []string
			descriptionPresent := false
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %s", err)
				}
				w.Header().Set("Content-Type", "application/json")
				switch action := r.Form.Get("Action"); action {
				case "CreateExpressConnectRouterAssociation":
					mu.Lock()
					createCalls++
					descriptionValues, descriptionPresent = r.PostForm["Description"]
					mu.Unlock()
					_, _ = w.Write([]byte(`{"AssociationId":"assoc-unit","Success":true}`))
				case "DescribeExpressConnectRouterAssociation":
					_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE","AllowedPrefixes":[],"Description":""}],"Success":true}`))
				default:
					t.Errorf("unexpected action: %q", action)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
				}
			}))
			t.Cleanup(server.Close)

			resourceSchema := expressconnectrouter.ResourceTRAssociation()
			attrs := map[string]cty.Value{}
			for name, attributeType := range resourceSchema.CoreConfigSchema().ImpliedType().AttributeTypes() {
				attrs[name] = cty.NullVal(attributeType)
			}
			attrs["ecr_id"] = cty.StringVal("ecr-unit")
			attrs["transit_router_id"] = cty.StringVal("tr-unit")
			attrs["association_region_id"] = cty.StringVal("cn-hangzhou")
			attrs["cen_id"] = cty.StringVal("cen-unit")
			attrs["transit_router_owner_id"] = cty.NumberIntVal(123456)
			attrs["create_attachment"] = cty.False
			attrs["description"] = tc.description
			rawConfig := cty.ObjectVal(attrs)
			diff, err := resourceSchema.Diff(context.Background(), nil, sdktf.NewResourceConfigShimmed(rawConfig, resourceSchema.CoreConfigSchema()), nil)
			if err != nil {
				t.Fatalf("plan through the public resource entry failed: %s", err)
			}
			diff.RawConfig = rawConfig
			diff.Meta = map[string]interface{}{
				schema.TimeoutKey: map[string]interface{}{schema.TimeoutCreate: 5 * time.Second},
			}
			if _, diags := resourceSchema.Apply(context.Background(), nil, diff, unitTestClient(t, server.URL)); diags.HasError() {
				t.Fatalf("create through the public resource entry failed: %v", diags)
			}

			mu.Lock()
			defer mu.Unlock()
			if createCalls != 1 {
				t.Fatalf("exactly one create write is expected, got %d", createCalls)
			}
			if descriptionPresent != tc.wantPresent {
				t.Fatalf("Description presence=%t, want %t; values=%v", descriptionPresent, tc.wantPresent, descriptionValues)
			}
			if tc.wantPresent && (len(descriptionValues) != 1 || descriptionValues[0] != tc.wantValue) {
				t.Fatalf("Description serialized as %v, want exactly %q", descriptionValues, tc.wantValue)
			}
		})
	}
}

func TestUnitTRAssociationCreateFindsTargetOnLaterPage(t *testing.T) {
	fastTRAssociationWaits(t)
	var mu sync.Mutex
	describeCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "CreateExpressConnectRouterAssociation":
			_, _ = w.Write([]byte(`{"AssociationId":"assoc-unit","Success":true}`))
		case "DescribeExpressConnectRouterAssociation":
			mu.Lock()
			describeCalls++
			call := describeCalls
			mu.Unlock()
			for key, want := range map[string]string{
				"AssociationId": "assoc-unit", "EcrId": "ecr-unit",
				"TransitRouterId": "tr-unit", "AssociationNodeType": "TR",
			} {
				if got := r.PostForm.Get(key); got != want {
					t.Errorf("page request %s=%q, want %q", key, got, want)
				}
			}
			if call%2 == 1 {
				if got := r.PostForm.Get("NextToken"); got != "" {
					t.Errorf("first page request NextToken=%q, want empty", got)
				}
				_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-other","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"CREATING","AllowedPrefixes":[]}],"NextToken":"page-2","Success":true}`))
				return
			}
			if got := r.PostForm.Get("NextToken"); got != "page-2" {
				t.Errorf("second page request NextToken=%q, want page-2", got)
			}
			_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-unit","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE","AllowedPrefixes":[],"Description":"cross-page-target"}],"Success":true}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	resourceSchema, diff, client := prepareTRAssociationCreate(t, server.URL, 5*time.Second)
	newState, diags := resourceSchema.Apply(context.Background(), nil, diff, client)
	if diags.HasError() {
		t.Fatalf("create must follow NextToken in both the state wait and final read: %v", diags)
	}
	if newState == nil || newState.ID != "ecr-unit:assoc-unit:tr-unit" {
		t.Fatalf("the second-page target must be stored, got %+v", newState)
	}
	if got := newState.Attributes["description"]; got != "cross-page-target" {
		t.Fatalf("final read must fill state from the second-page target, description=%q", got)
	}
	if got := newState.Attributes["status"]; got != "ACTIVE" {
		t.Fatalf("state wait must use the second-page target status, got %q", got)
	}
	mu.Lock()
	defer mu.Unlock()
	if describeCalls != 4 {
		t.Fatalf("create wait and final read must each read two pages, got %d Describe calls", describeCalls)
	}
}

func TestUnitTRAssociationDeleteRefreshExhaustsPagesBeforeNotFound(t *testing.T) {
	fastTRAssociationWaits(t)
	var mu sync.Mutex
	deleteCalls := 0
	describeCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "ListTransitRouterEcrAttachments":
			_, _ = w.Write([]byte(`{"TransitRouterAttachments":[],"Success":true}`))
		case "DeleteExpressConnectRouterAssociation":
			mu.Lock()
			deleteCalls++
			mu.Unlock()
			_, _ = w.Write([]byte(`{"Success":true}`))
		case "DescribeExpressConnectRouterAssociation":
			mu.Lock()
			describeCalls++
			call := describeCalls
			mu.Unlock()
			for key, want := range map[string]string{
				"AssociationId": "assoc-unit", "EcrId": "ecr-unit",
				"TransitRouterId": "tr-unit", "AssociationNodeType": "TR",
			} {
				if got := r.PostForm.Get(key); got != want {
					t.Errorf("page request %s=%q, want %q", key, got, want)
				}
			}
			if call == 1 {
				if got := r.PostForm.Get("NextToken"); got != "" {
					t.Errorf("first page request NextToken=%q, want empty", got)
				}
				_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-other","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE"}],"NextToken":"page-2","Success":true}`))
				return
			}
			if got := r.PostForm.Get("NextToken"); got != "page-2" {
				t.Errorf("second page request NextToken=%q, want page-2", got)
			}
			_, _ = w.Write([]byte(`{"AssociationList":[{"AssociationId":"assoc-third","AssociationNodeType":"TR","EcrId":"ecr-unit","TransitRouterId":"tr-unit","Status":"ACTIVE"}],"Success":true}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	resourceSchema := expressconnectrouter.ResourceTRAssociation()
	if err := resourceSchema.Delete(resourceSchema.Data(importedTRAssociationState()), unitTestClient(t, server.URL)); err != nil {
		t.Fatalf("delete must report NotFound only after all pages exclude the target: %s", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if deleteCalls != 1 {
		t.Fatalf("exactly one delete write is expected, got %d", deleteCalls)
	}
	if describeCalls != 2 {
		t.Fatalf("delete refresh must exhaust both pages before NotFound, got %d Describe calls", describeCalls)
	}
}

// Test ExpressConnectRouter ExpressConnectRouterTrAssociation. <<< Resource test cases, automatically generated.
