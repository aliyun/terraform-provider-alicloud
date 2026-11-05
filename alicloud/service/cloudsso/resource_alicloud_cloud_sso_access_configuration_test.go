package cloudsso_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	convert "github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	cloudsso "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudsso"
	sweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testSweepCloudSsoDirectoryAccessConfiguration(region, directoryId string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"",
	}
	action := "ListAccessConfigurations"
	request := map[string]interface{}{}
	request["DirectoryId"] = directoryId

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	resp, err := jsonpath.Get("$.AccessConfigurations", response)
	totalCount, formatErr := convert.FormatInt(response["TotalCounts"])
	if formatErr != nil {
		return formatErr
	}
	if totalCount != 0 && err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.AccessAssignments", action, err)
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
		action := "DeleteAccessConfiguration"
		req := map[string]interface{}{
			"DirectoryId":                   directoryId,
			"AccessConfigurationId":         item["AccessConfigurationId"],
			"ForceRemovePermissionPolicies": true,
		}
		_, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, req, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Cloud Sso AccessAssignment (%s): %s", item["AccessConfigurationName"].(string), err)
		}
		log.Printf("[INFO] Delete Cloud Sso AccessAssignment success: %s ", item["AccessConfigurationName"].(string))
	}
	return nil
}

func TestAccAliCloudCloudSSOAccessConfiguration_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	resourceId := "alicloud_cloud_sso_access_configuration.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudSSOAccessConfigurationMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudsso.NewCloudssoService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudSsoAccessConfiguration")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccconfig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudSSOAccessConfigurationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_id":              "${alicloud_cloud_sso_directory.default.id}",
					"access_configuration_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":              acctest.CHECKSET,
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	resourceId := "alicloud_cloud_sso_access_configuration.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudSSOAccessConfigurationMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudsso.NewCloudssoService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudSsoAccessConfiguration")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccconfig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudSSOAccessConfigurationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_id":                     "${alicloud_cloud_sso_directory.default.id}",
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
						"directory_id":              acctest.CHECKSET,
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
	"session_duration":        acctest.CHECKSET,
	"access_configuration_id": acctest.CHECKSET,
}

func AliCloudCloudSSOAccessConfigurationBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_cloud_sso_directory" "default" {
	  directory_name    = var.name
	}
	locals{
	  directory_id = alicloud_cloud_sso_directory.default.id
	}
`, name)
}
