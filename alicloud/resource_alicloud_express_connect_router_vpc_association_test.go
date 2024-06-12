package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ExpressConnectRouter ExpressConnectRouterVpcAssociation. >>> Resource test cases, automatically generated.
// Case 初始版本测试用例 6348
func TestAccAliCloudExpressConnectRouterExpressConnectRouterVpcAssociation_basic6348(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_vpc_association.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationMap6348)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectRouterServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectRouterExpressConnectRouterVpcAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervpcassociation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationBasicDependence6348)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": defaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_region_id": CHECKSET,
						"vpc_id":                CHECKSET,
						"ecr_id":                CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"172.16.4.0/24", "172.16.3.0/24", "172.16.2.0/24", "172.16.1.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"172.16.5.0/25", "172.16.6.0/24"},
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
					"allowed_prefixes": []string{
						"172.16.4.0/24", "172.16.3.0/24", "172.16.2.0/24", "172.16.1.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": defaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
					"allowed_prefixes": []string{
						"172.16.4.0/24", "172.16.3.0/24", "172.16.2.0/24", "172.16.1.0/24"},
					"vpc_owner_id": "1511928242963727",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_region_id": CHECKSET,
						"vpc_id":                CHECKSET,
						"ecr_id":                CHECKSET,
						"allowed_prefixes.#":    "4",
						"vpc_owner_id":          CHECKSET,
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

var AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationMap6348 = map[string]string{
	"status":         CHECKSET,
	"create_time":    CHECKSET,
	"association_id": CHECKSET,
}

func AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationBasicDependence6348(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default8qAtD6" {
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_express_connect_router_express_connect_router" "defaultM9YxGW" {
  alibaba_side_asn = "65533"
}


`, name)
}

// Case 初始版本测试用例 6348  twin
func TestAccAliCloudExpressConnectRouterExpressConnectRouterVpcAssociation_basic6348_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_vpc_association.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationMap6348)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectRouterServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectRouterExpressConnectRouterVpcAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervpcassociation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationBasicDependence6348)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": defaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
					"allowed_prefixes": []string{
						"172.16.4.0/24", "172.16.3.0/24", "172.16.2.0/24", "172.16.1.0/24"},
					"vpc_owner_id": "1511928242963727",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_region_id": CHECKSET,
						"vpc_id":                CHECKSET,
						"ecr_id":                CHECKSET,
						"allowed_prefixes.#":    "4",
						"vpc_owner_id":          CHECKSET,
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

// Case 初始版本测试用例 6348  raw
func TestAccAliCloudExpressConnectRouterExpressConnectRouterVpcAssociation_basic6348_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_vpc_association.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationMap6348)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectRouterServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectRouterExpressConnectRouterVpcAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervpcassociation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationBasicDependence6348)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": defaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
					"allowed_prefixes": []string{
						"172.16.4.0/24", "172.16.3.0/24", "172.16.2.0/24", "172.16.1.0/24"},
					"vpc_owner_id": "1511928242963727",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_region_id": CHECKSET,
						"vpc_id":                CHECKSET,
						"ecr_id":                CHECKSET,
						"allowed_prefixes.#":    "4",
						"vpc_owner_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"172.16.5.0/25", "172.16.6.0/24"},
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
					"allowed_prefixes": []string{
						"172.16.4.0/24", "172.16.3.0/24", "172.16.2.0/24", "172.16.1.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "4",
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

// Test ExpressConnectRouter ExpressConnectRouterVpcAssociation. <<< Resource test cases, automatically generated.
