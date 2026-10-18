package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
	}
	action := "ListAccessConfigurationProvisionings"
	request := map[string]interface{}{}
	request["DirectoryId"] = directoryId

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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

		accessConfigurationName, _ := item["AccessConfigurationName"].(string)
		if protectedSweepVpcName(accessConfigurationName) {
			log.Printf("[INFO] Skipping protected Cloud Sso AccessConfigurationProvisioning: %s", accessConfigurationName)
			continue
		}
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(accessConfigurationName), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Cloud Sso AccessConfigurationName: %s", accessConfigurationName)
				continue
			}
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
			log.Printf("[ERROR] Failed to delete Cloud Sso AccessConfigurationProvisioning (%s): %s", accessConfigurationName, err)
		} else {
			log.Printf("[INFO] Delete Cloud Sso AccessConfigurationProvisioning success: %s ", accessConfigurationName)
		}
	}
	return nil
}

func init() {
	resource.AddTestSweepers("alicloud_cloud_sso_access_configuration_provisioning", &resource.Sweeper{
		Name: "alicloud_cloud_sso_access_configuration_provisioning",
		F:    testSweepCloudSsoAccessConfigurationProvisioning,
	})
}

// testSweepCloudSsoAccessConfigurationProvisioning is the registered sweeper
// entry point for alicloud_cloud_sso_access_configuration_provisioning.
// Provisionings live inside a directory, so this wrapper enumerates every
// directory in the region and delegates to
// testSweepCloudSsoDirectoryAccessConfigurationProvisioning for each one.
func testSweepCloudSsoAccessConfigurationProvisioning(region string) error {
	directoryIds, err := testSweepCloudSsoDirectoryIds(region)
	if err != nil {
		return err
	}
	for _, directoryId := range directoryIds {
		if err := testSweepCloudSsoDirectoryAccessConfigurationProvisioning(region, directoryId); err != nil {
			log.Printf("[ERROR] Failed to sweep Cloud Sso AccessConfigurationProvisioning in directory %s: %s", directoryId, err)
		}
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
