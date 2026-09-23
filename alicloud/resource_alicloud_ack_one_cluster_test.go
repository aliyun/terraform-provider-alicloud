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
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
							"vswitches": []string{
								"${data.alicloud_vswitches.default.ids.0}",
							},
						},
					},
					"profile":        "Default",
					"argocd_enabled": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":   CHECKSET,
						"profile":        "Default",
						"argocd_enabled": "false",
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
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
							"vswitches": []string{
								"${data.alicloud_vswitches.default.ids.0}"},
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
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
							"vswitches": []string{
								"${data.alicloud_vswitches.default.ids.0}"},
						},
					},
					"argocd_enabled": false,
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

func TestAccAliCloudAckOneCluster_updateArgocdEnabled(t *testing.T) {
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
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
							"vswitches": []string{
								"${data.alicloud_vswitches.default.ids.0}"},
						},
					},
					"profile":        "Default",
					"argocd_enabled": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":   CHECKSET,
						"profile":        "Default",
						"argocd_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": name,
					"network": []map[string]interface{}{
						{
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
							"vswitches": []string{
								"${data.alicloud_vswitches.default.ids.0}"},
						},
					},
					"profile":        "Default",
					"argocd_enabled": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":   CHECKSET,
						"profile":        "Default",
						"argocd_enabled": "false",
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

var AlicloudAckOneClusterMap4593 = map[string]string{
	"status":       CHECKSET,
	"create_time":  CHECKSET,
	"cluster_name": CHECKSET,
	"profile":      CHECKSET,
}

// use existing vpc and vswitch because https://project.aone.alibaba-inc.com/v2/project/206563/bug/63062057
func AlicloudAckOneClusterBasicDependence4593(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "default" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.default.ids.0
  security_group_name = "${var.name}"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}
`, name)
}

// Test AckOne Cluster. <<< Resource test cases, automatically generated.
