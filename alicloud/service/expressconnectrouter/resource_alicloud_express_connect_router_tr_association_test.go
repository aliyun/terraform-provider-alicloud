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
		},
	})
}

// Test ExpressConnectRouter ExpressConnectRouterTrAssociation. <<< Resource test cases, automatically generated.
