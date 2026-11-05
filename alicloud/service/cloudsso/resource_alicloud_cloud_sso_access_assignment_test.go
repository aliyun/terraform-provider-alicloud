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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudsso"
	sweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testSweepCloudSSODirectoryAccessAssignment(region, directoryId string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"",
	}
	action := "ListAccessAssignments"
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

	resp, err := jsonpath.Get("$.AccessAssignments", response)
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
			log.Printf("[INFO] Skipping Cloud SSO AccessConfigurationName: %s", item["AccessConfigurationName"].(string))
			continue
		}
		action := "DeleteAccessAssignment"
		req := map[string]interface{}{
			"DirectoryId":           directoryId,
			"AccessConfigurationId": item["AccessConfigurationId"],
			"TargetType":            item["TargetType"],
			"TargetId":              item["TargetId"],
			"PrincipalType":         item["PrincipalType"],
			"PrincipalId":           item["PrincipalId"],
		}
		_, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, req, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Cloud SSO AccessAssignment (%s): %s", item["AccessConfigurationName"].(string), err)
		}
		log.Printf("[INFO] Delete Cloud SSO AccessAssignment success: %s ", item["AccessConfigurationName"].(string))
	}
	return nil
}

func TestAccAliCloudCloudSSOAccessAssignment_basic10051(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_access_assignment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudSSOAccessAssignmentMap10051)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudsso.NewCloudSSOServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudSSOAccessAssignment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudSSOAccessAssignmentBasicDependence10051)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.CloudSsoSupportRegions)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_id":            "${local.directory_id}",
					"principal_id":            "${alicloud_cloud_sso_user.default.user_id}",
					"target_type":             "RD-Account",
					"access_configuration_id": "${alicloud_cloud_sso_access_configuration.default.access_configuration_id}",
					"principal_type":          "User",
					"target_id":               "${data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":            acctest.CHECKSET,
						"access_configuration_id": acctest.CHECKSET,
						"target_type":             "RD-Account",
						"target_id":               acctest.CHECKSET,
						"principal_type":          "User",
						"principal_id":            acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deprovision_strategy"},
			},
		},
	})
}

func TestAccAliCloudCloudSSOAccessAssignment_basic10051_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_access_assignment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudSSOAccessAssignmentMap10051)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudsso.NewCloudSSOServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudSSOAccessAssignment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudSSOAccessAssignmentBasicDependence10051)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.CloudSsoSupportRegions)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_id":            "${local.directory_id}",
					"principal_id":            "${alicloud_cloud_sso_user.default.user_id}",
					"target_type":             "RD-Account",
					"access_configuration_id": "${alicloud_cloud_sso_access_configuration.default.access_configuration_id}",
					"principal_type":          "User",
					"target_id":               "${data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id}",
					"deprovision_strategy":    "DeprovisionForLastAccessAssignmentOnAccount",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":            acctest.CHECKSET,
						"access_configuration_id": acctest.CHECKSET,
						"target_type":             "RD-Account",
						"target_id":               acctest.CHECKSET,
						"principal_type":          "User",
						"principal_id":            acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deprovision_strategy"},
			},
		},
	})
}

func AliCloudCloudSSOAccessAssignmentBasicDependence10051(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_directories" "default" {
	}


	resource "alicloud_cloud_sso_directory" "default" {
  		directory_name = var.name
	}

	resource "alicloud_cloud_sso_group" "default" {
  		directory_id = local.directory_id
  		group_name   = var.name
  		description  = var.name
	}

	resource "alicloud_cloud_sso_user" "default" {
  		directory_id = local.directory_id
  		user_name    = var.name
	}

	resource "alicloud_cloud_sso_access_configuration" "default" {
  		directory_id              = local.directory_id
  		access_configuration_name = var.name
	}

	locals {
  		directory_id = alicloud_cloud_sso_directory.default.id
	}
`, name)
}

func TestAccAliCloudCloudSSOAccessAssignment_basic10355(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_access_assignment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudSSOAccessAssignmentMap10051)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudsso.NewCloudSSOServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudSSOAccessAssignment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudSSOAccessAssignmentBasicDependence10051)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.CloudSsoSupportRegions)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_id":            "${local.directory_id}",
					"principal_id":            "${alicloud_cloud_sso_group.default.group_id}",
					"target_type":             "RD-Account",
					"access_configuration_id": "${alicloud_cloud_sso_access_configuration.default.access_configuration_id}",
					"principal_type":          "Group",
					"target_id":               "${data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":            acctest.CHECKSET,
						"access_configuration_id": acctest.CHECKSET,
						"target_type":             "RD-Account",
						"target_id":               acctest.CHECKSET,
						"principal_type":          "Group",
						"principal_id":            acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deprovision_strategy"},
			},
		},
	})
}

func TestAccAliCloudCloudSSOAccessAssignment_basic10355_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_access_assignment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudSSOAccessAssignmentMap10051)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudsso.NewCloudSSOServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudSSOAccessAssignment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudSSOAccessAssignmentBasicDependence10051)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.CloudSsoSupportRegions)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_id":            "${local.directory_id}",
					"principal_id":            "${alicloud_cloud_sso_group.default.group_id}",
					"target_type":             "RD-Account",
					"access_configuration_id": "${alicloud_cloud_sso_access_configuration.default.access_configuration_id}",
					"principal_type":          "Group",
					"target_id":               "${data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id}",
					"deprovision_strategy":    "DeprovisionForLastAccessAssignmentOnAccount",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":            acctest.CHECKSET,
						"access_configuration_id": acctest.CHECKSET,
						"target_type":             "RD-Account",
						"target_id":               acctest.CHECKSET,
						"principal_type":          "Group",
						"principal_id":            acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deprovision_strategy"},
			},
		},
	})
}

var AliCloudCloudSSOAccessAssignmentMap10051 = map[string]string{
	"create_time": acctest.CHECKSET,
}
