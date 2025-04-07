package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testSweepCloudSsoDirectoryAccessConfigurationProvisioning(region, directoryId string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"",
	}
	action := "ListAccessConfigurationProvisionings"
	request := map[string]interface{}{}
	request["DirectoryId"] = directoryId

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	resp, err := jsonpath.Get("$.AccessConfigurationProvisionings", response)
	if formatInt(response["TotalCounts"]) != 0 && err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.AccessConfigurationProvisionings", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["AccessConfigurationName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Cloud Sso AccessConfigurationName: %s", item["AccessConfigurationName"].(string))
			continue
		}
		action := "DeprovisionAccessConfiguration"
		req := map[string]interface{}{
			"DirectoryId":           directoryId,
			"AccessConfigurationId": item["AccessConfigurationId"],
			"TargetType":            item["TargetType"],
			"TargetId":              item["TargetId"],
		}
		_, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, req, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Cloud Sso AccessAssignment (%s): %s", item["AccessConfigurationName"].(string), err)
		}
		log.Printf("[INFO] Delete Cloud Sso AccessAssignment success: %s ", item["AccessConfigurationName"].(string))
	}
	return nil
}

func TestAccAlicloudCloudSSOAccessConfigurationProvisioning_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_access_configuration_provisioning.default"
	checkoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCloudSSOAccessConfigurationProvisioningMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudssoService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSsoAccessConfigurationProvisioning")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudSSOAccessConfigurationProvisioningBasicDependence0)
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
					"directory_id":            "${local.directory_id}",
					"access_configuration_id": "${alicloud_cloud_sso_access_configuration.default.access_configuration_id}",
					"target_type":             "RD-Account",
					"target_id":               "${data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":            CHECKSET,
						"access_configuration_id": CHECKSET,
						"target_type":             "RD-Account",
						"target_id":               CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudCloudSSOAccessConfigurationProvisioningMap0 = map[string]string{}

func AlicloudCloudSSOAccessConfigurationProvisioningBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_cloud_sso_directories" "default" {}
data "alicloud_resource_manager_resource_directories" "default" {}
resource "alicloud_cloud_sso_directory" "default" {
  count             = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name    = var.name
}
locals{
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}
resource "alicloud_cloud_sso_access_configuration" "default" {
  access_configuration_name = var.name
  directory_id = local.directory_id
}
`, name)
}
