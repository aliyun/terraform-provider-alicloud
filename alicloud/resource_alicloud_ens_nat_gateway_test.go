package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens NatGateway. >>> Resource test cases, automatically generated.
// Case Nat网关测试_20240507 6657
func TestAccAliCloudEnsNatGateway_basic6657(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_nat_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsNatGatewayMap6657)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsNatGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensnatgateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsNatGatewayBasicDependence6657)
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
					"vswitch_id":    "${alicloud_ens_vswitch.defaulteFw783.id}",
					"ens_region_id": "${alicloud_ens_vswitch.defaulteFw783.ens_region_id}",
					"network_id":    "${alicloud_ens_vswitch.defaulteFw783.network_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":    CHECKSET,
						"ens_region_id": CHECKSET,
						"network_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nat_name": "test1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nat_name": "test1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nat_name": "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nat_name": "test2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":    "${alicloud_ens_vswitch.defaulteFw783.id}",
					"ens_region_id": "${alicloud_ens_vswitch.defaulteFw783.ens_region_id}",
					"network_id":    "${alicloud_ens_vswitch.defaulteFw783.network_id}",
					"instance_type": "enat.default",
					"nat_name":      "test1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":    CHECKSET,
						"ens_region_id": CHECKSET,
						"network_id":    CHECKSET,
						"instance_type": "enat.default",
						"nat_name":      "test1",
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

var AlicloudEnsNatGatewayMap6657 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEnsNatGatewayBasicDependence6657(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "defaultObbrL7" {
  network_name  = var.name
  description   = "测试用例-测试NAT使用"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "defaulteFw783" {
  cidr_block    = "10.0.8.0/24"
  vswitch_name  = var.name
  ens_region_id = alicloud_ens_network.defaultObbrL7.ens_region_id
  network_id    = alicloud_ens_network.defaultObbrL7.id
}


`, name)
}

// Case Nat网关测试_20240507 6657  twin
func TestAccAliCloudEnsNatGateway_basic6657_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_nat_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsNatGatewayMap6657)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsNatGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensnatgateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsNatGatewayBasicDependence6657)
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
					"vswitch_id":    "${alicloud_ens_vswitch.defaulteFw783.id}",
					"ens_region_id": "${alicloud_ens_vswitch.defaulteFw783.ens_region_id}",
					"network_id":    "${alicloud_ens_vswitch.defaulteFw783.network_id}",
					"instance_type": "enat.default",
					"nat_name":      "test1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":    CHECKSET,
						"ens_region_id": CHECKSET,
						"network_id":    CHECKSET,
						"instance_type": "enat.default",
						"nat_name":      "test1",
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

// Case Nat网关测试_20240507 6657  raw
func TestAccAliCloudEnsNatGateway_basic6657_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_nat_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsNatGatewayMap6657)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsNatGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensnatgateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsNatGatewayBasicDependence6657)
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
					"vswitch_id":    "${alicloud_ens_vswitch.defaulteFw783.id}",
					"ens_region_id": "${alicloud_ens_vswitch.defaulteFw783.ens_region_id}",
					"network_id":    "${alicloud_ens_vswitch.defaulteFw783.network_id}",
					"instance_type": "enat.default",
					"nat_name":      "test1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":    CHECKSET,
						"ens_region_id": CHECKSET,
						"network_id":    CHECKSET,
						"instance_type": "enat.default",
						"nat_name":      "test1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nat_name": "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nat_name": "test2",
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

// Test Ens NatGateway. <<< Resource test cases, automatically generated.
