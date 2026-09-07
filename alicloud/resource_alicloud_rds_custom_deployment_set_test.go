package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Rds CustomDeploymentSet. >>> Resource test cases, automatically generated.
// Case DeploymentSetTest_1115 8966
func TestAccAliCloudRdsCustomDeploymentSet_basic8966(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom_deployment_set.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsCustomDeploymentSetMap8966)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsCustomDeploymentSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdscustomdeploymentset%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsCustomDeploymentSetBasicDependence8966)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"on_unable_to_redeploy_failed_instance": "CancelMembershipAndStart",
					"custom_deployment_set_name":            name,
					"description":                           "2024:11:19 1010:1111:0707",
					"group_count":                           "3",
					"strategy":                              "Availability",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"on_unable_to_redeploy_failed_instance": "CancelMembershipAndStart",
						"custom_deployment_set_name":            name,
						"description":                           CHECKSET,
						"group_count":                           "3",
						"strategy":                              "Availability",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"group_count", "on_unable_to_redeploy_failed_instance"},
			},
		},
	})
}

var AlicloudRdsCustomDeploymentSetMap8966 = map[string]string{
	"status": CHECKSET,
}

func AlicloudRdsCustomDeploymentSetBasicDependence8966(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case DeploymentSetTest 8365
func TestAccAliCloudRdsCustomDeploymentSet_basic8365(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom_deployment_set.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsCustomDeploymentSetMap8365)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsCustomDeploymentSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdscustomdeploymentset%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsCustomDeploymentSetBasicDependence8365)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"on_unable_to_redeploy_failed_instance": "CancelMembershipAndStart",
					"custom_deployment_set_name":            name,
					"description":                           "2024:11:19 1010:1111:0707",
					"group_count":                           "3",
					"strategy":                              "Availability",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"on_unable_to_redeploy_failed_instance": "CancelMembershipAndStart",
						"custom_deployment_set_name":            name,
						"description":                           CHECKSET,
						"group_count":                           "3",
						"strategy":                              "Availability",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"group_count", "on_unable_to_redeploy_failed_instance"},
			},
		},
	})
}

var AlicloudRdsCustomDeploymentSetMap8365 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRdsCustomDeploymentSetBasicDependence8365(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Rds CustomDeploymentSet. <<< Resource test cases, automatically generated.
