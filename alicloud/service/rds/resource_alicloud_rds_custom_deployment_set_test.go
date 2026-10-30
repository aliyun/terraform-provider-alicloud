package rds_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Rds CustomDeploymentSet. >>> Resource test cases, automatically generated.
// Case DeploymentSetTest_1115 8966
func TestAccAliCloudRdsCustomDeploymentSet_basic8966(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom_deployment_set.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRdsCustomDeploymentSetMap8966)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rds.NewRdsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRdsCustomDeploymentSet")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdscustomdeploymentset%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRdsCustomDeploymentSetBasicDependence8966)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"description":                           acctest.CHECKSET,
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
	"status": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRdsCustomDeploymentSetMap8365)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rds.NewRdsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRdsCustomDeploymentSet")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdscustomdeploymentset%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRdsCustomDeploymentSetBasicDependence8365)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"description":                           acctest.CHECKSET,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudRdsCustomDeploymentSetBasicDependence8365(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Rds CustomDeploymentSet. <<< Resource test cases, automatically generated.
