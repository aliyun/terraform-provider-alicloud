package ens_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ens"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Ens Vswitch. >>> Resource test cases, automatically generated.
// Case 5060
func TestAccAliCloudEnsVswitch_basic5060(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_vswitch.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEnsVswitchMap5060)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ens.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEnsVswitch")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensvswitch%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEnsVswitchBasicDependence5060)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"network_id":    acctest.CHECKSET,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"network_id":  acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEnsVswitchMap5060)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ens.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEnsVswitch")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensvswitch%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEnsVswitchBasicDependence5060)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"network_id":    acctest.CHECKSET,
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
