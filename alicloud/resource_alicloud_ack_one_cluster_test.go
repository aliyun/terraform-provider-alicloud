package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test AckOne Cluster. >>> Resource test cases, automatically generated.
// Case 4593
func TestAccAliCloudAckOneCluster_basic4593(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ack_one_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudAckOneClusterMap4593)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckOneServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckOneCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sackonecluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckOneClusterBasicDependence4593)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.AckOneSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": name,
					"network": []map[string]interface{}{
						{
							"vpc_id": "${alicloud_vpc.defaultVpc.id}",
							"vswitches": []string{
								"${alicloud_vswitch.defaultyVSwitch.id}"},
						},
					},
					"profile": "Default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name": CHECKSET,
						"profile":      "Default",
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

func TestAccAliCloudAckOneCluster_basic4593_XFlow(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ack_one_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudAckOneClusterMap4593)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckOneServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckOneCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sackonecluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckOneClusterBasicDependence4593)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.AckOneSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": name,
					"network": []map[string]interface{}{
						{
							"vpc_id": "${alicloud_vpc.defaultVpc.id}",
							"vswitches": []string{
								"${alicloud_vswitch.defaultyVSwitch.id}"},
						},
					},
					"profile": "XFlow",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name": CHECKSET,
						"profile":      "XFlow",
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

func TestAccAliCloudAckOneCluster_basic4593_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ack_one_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudAckOneClusterMap4593)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckOneServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckOneCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sackonecluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckOneClusterBasicDependence4593)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.AckOneSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"network": []map[string]interface{}{
						{
							"vpc_id": "${alicloud_vpc.defaultVpc.id}",
							"vswitches": []string{
								"${alicloud_vswitch.defaultyVSwitch.id}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudAckOneClusterMap4593 = map[string]string{
	"status":       CHECKSET,
	"create_time":  CHECKSET,
	"cluster_name": CHECKSET,
	"profile":      CHECKSET,
}

func AlicloudAckOneClusterBasicDependence4593(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

}

resource "alicloud_vswitch" "defaultyVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  cidr_block   = "172.16.2.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name

}


`, name)
}

// Test AckOne Cluster. <<< Resource test cases, automatically generated.
