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

// Test Ens Network. >>> Resource test cases, automatically generated.
// Case 5077
func TestAccAliCloudEnsNetwork_basic5077(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_network.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEnsNetworkMap5077)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ens.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEnsNetwork")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensnetwork%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEnsNetworkBasicDependence5077)
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
					"network_name":  name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block":    "192.168.2.0/24",
						"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
						"network_name":  name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "NetworkDescription_autotest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "NetworkDescription_autotest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "NetworkDescription_UPDATE_autotest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "NetworkDescription_UPDATE_autotest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_name":  name + "_update",
					"description":   "NetworkDescription_autotest",
					"cidr_block":    "192.168.2.0/24",
					"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_name":  name + "_update",
						"description":   "NetworkDescription_autotest",
						"cidr_block":    "192.168.2.0/24",
						"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
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

var AlicloudEnsNetworkMap5077 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudEnsNetworkBasicDependence5077(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 5077  twin
func TestAccAliCloudEnsNetwork_basic5077_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_network.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEnsNetworkMap5077)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ens.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEnsNetwork")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensnetwork%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEnsNetworkBasicDependence5077)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"network_name":  name,
					"description":   "NetworkDescription_UPDATE_autotest",
					"cidr_block":    "192.168.2.0/24",
					"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_name":  name,
						"description":   "NetworkDescription_UPDATE_autotest",
						"cidr_block":    "192.168.2.0/24",
						"ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
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

// Test Ens Network. <<< Resource test cases, automatically generated.
