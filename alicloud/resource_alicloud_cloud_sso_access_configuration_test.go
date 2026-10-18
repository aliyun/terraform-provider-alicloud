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

func testSweepCloudSsoDirectoryAccessConfiguration(region, directoryId string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListAccessConfigurations"
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

	resp, err := jsonpath.Get("$.AccessConfigurations", response)
	if formatInt(response["TotalCounts"]) != 0 && err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.AccessAssignments", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})

		accessConfigurationName, _ := item["AccessConfigurationName"].(string)
		if protectedSweepVpcName(accessConfigurationName) {
			log.Printf("[INFO] Skipping protected Cloud Sso AccessConfiguration: %s", accessConfigurationName)
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
		action := "DeleteAccessConfiguration"
		req := map[string]interface{}{
			"DirectoryId":                   directoryId,
			"AccessConfigurationId":         item["AccessConfigurationId"],
			"ForceRemovePermissionPolicies": true,
		}
		_, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, req, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Cloud Sso AccessConfiguration (%s): %s", accessConfigurationName, err)
		} else {
			log.Printf("[INFO] Delete Cloud Sso AccessConfiguration success: %s ", accessConfigurationName)
		}
	}
	return nil
}

func init() {
	resource.AddTestSweepers("alicloud_cloud_sso_access_configuration", &resource.Sweeper{
		Name: "alicloud_cloud_sso_access_configuration",
		F:    testSweepCloudSsoAccessConfiguration,
		// DeleteAccessConfiguration fails with
		// DeletionConflict.AccessConfiguration.Provisioning while a config is
		// still provisioned to an account, so deprovision first.
		Dependencies: []string{
			"alicloud_cloud_sso_access_configuration_provisioning",
			"alicloud_cloud_sso_access_assignment",
		},
	})
}

// testSweepCloudSsoAccessConfiguration is the registered sweeper entry point
// for alicloud_cloud_sso_access_configuration. Access configurations live
// inside a directory, so this wrapper enumerates every directory in the region
// and delegates to testSweepCloudSsoDirectoryAccessConfiguration for each one.
func testSweepCloudSsoAccessConfiguration(region string) error {
	directoryIds, err := testSweepCloudSsoDirectoryIds(region)
	if err != nil {
		return err
	}
	for _, directoryId := range directoryIds {
		if err := testSweepCloudSsoDirectoryAccessConfiguration(region, directoryId); err != nil {
			log.Printf("[ERROR] Failed to sweep Cloud Sso AccessConfiguration in directory %s: %s", directoryId, err)
		}
	}
	return nil
}

func TestAccAliCloudCloudSSOAccessConfiguration_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	resourceId := "alicloud_cloud_sso_access_configuration.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudSSOAccessConfigurationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudssoService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSsoAccessConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudSSOAccessConfigurationBasicDependence0)
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
					"directory_id":              "${data.alicloud_cloud_sso_directories.default.directories.0.id}",
					"access_configuration_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":              CHECKSET,
						"access_configuration_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"session_duration": "1800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"session_duration": "1800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"relay_state": "https://cloudsso.console.aliyun.com/test1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"relay_state": "https://cloudsso.console.aliyun.com/test1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permission_policies": []map[string]interface{}{
						{
							"permission_policy_type": "System",
							"permission_policy_name": "ReadOnlyAccess",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permission_policies.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permission_policies": []map[string]interface{}{
						{
							"permission_policy_type": "System",
							"permission_policy_name": "ReadOnlyAccess",
						},
						{
							"permission_policy_type": "System",
							"permission_policy_name": "AliyunECSFullAccess",
						},
						{
							"permission_policy_type":     "Inline",
							"permission_policy_name":     "oos",
							"permission_policy_document": `{\"Statement\": [{\"Action\": \"ecs:Get\",\"Effect\": \"Allow\",\"Resource\": \"*\"}],\"Version\": \"1\"}`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permission_policies.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permission_policies": []map[string]interface{}{
						{
							"permission_policy_type":     "Inline",
							"permission_policy_name":     "oos",
							"permission_policy_document": `{\"Statement\": [{\"Action\": \"oss:*\",\"Effect\": \"Allow\",\"Resource\": \"*\"}],\"Version\": \"1\"}`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permission_policies.#": "1",
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

func TestAccAliCloudCloudSSOAccessConfiguration_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	resourceId := "alicloud_cloud_sso_access_configuration.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudSSOAccessConfigurationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudssoService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSsoAccessConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudSSOAccessConfigurationBasicDependence0)
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
					"directory_id":                     "${data.alicloud_cloud_sso_directories.default.directories.0.id}",
					"access_configuration_name":        name,
					"session_duration":                 "1800",
					"relay_state":                      "https://cloudsso.console.aliyun.com/test1",
					"description":                      name,
					"force_remove_permission_policies": "true",
					"permission_policies": []map[string]interface{}{
						{
							"permission_policy_type": "System",
							"permission_policy_name": "ReadOnlyAccess",
						},
						{
							"permission_policy_type": "System",
							"permission_policy_name": "AliyunECSFullAccess",
						},
						{
							"permission_policy_type":     "Inline",
							"permission_policy_name":     "oos",
							"permission_policy_document": `{\"Statement\": [{\"Action\": \"ecs:Get\",\"Effect\": \"Allow\",\"Resource\": \"*\"}],\"Version\": \"1\"}`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":              CHECKSET,
						"access_configuration_name": name,
						"session_duration":          "1800",
						"relay_state":               "https://cloudsso.console.aliyun.com/test1",
						"description":               name,
						"permission_policies.#":     "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_remove_permission_policies"},
			},
		},
	})
}

var AliCloudCloudSSOAccessConfigurationMap0 = map[string]string{
	"session_duration":        CHECKSET,
	"access_configuration_id": CHECKSET,
}

func AliCloudCloudSSOAccessConfigurationBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_cloud_sso_directories" "default" {
	}
`, name)
}
