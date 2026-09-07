package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens Vswitch. >>> Resource test cases, automatically generated.
// Case 5060
func TestAccAliCloudEnsVswitch_basic5060(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_vswitch.default"
	ra := resourceAttrInit(resourceId, AliCloudEnsVswitchMap5060)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsVswitch")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensvswitch%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEnsVswitchBasicDependence5060)
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
					"cidr_block":    "192.168.2.0/24",
					"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
					"network_id":    "${alicloud_ens_network.default23T2cD.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block":    "192.168.2.0/24",
						"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
						"network_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "VSwitchDescription_autotest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "VSwitchDescription_autotest",
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

var AliCloudEnsVswitchMap5060 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"network_id":  CHECKSET,
}

func AliCloudEnsVswitchBasicDependence5060(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ens_network" "default23T2cD" {
  network_name = var.name

  description   = "VSwitchDescription_autotest"
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}


`, name)
}

// Case 5060  twin
func TestAccAliCloudEnsVswitch_basic5060_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_vswitch.default"
	ra := resourceAttrInit(resourceId, AliCloudEnsVswitchMap5060)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsVswitch")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensvswitch%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEnsVswitchBasicDependence5060)
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
					"description":   "VSwitchDescription_UPDATE_autotest",
					"cidr_block":    "192.168.2.0/24",
					"vswitch_name":  name,
					"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
					"network_id":    "${alicloud_ens_network.default23T2cD.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   "VSwitchDescription_UPDATE_autotest",
						"cidr_block":    "192.168.2.0/24",
						"vswitch_name":  name,
						"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
						"network_id":    CHECKSET,
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

// Test Ens Vswitch. <<< Resource test cases, automatically generated.
