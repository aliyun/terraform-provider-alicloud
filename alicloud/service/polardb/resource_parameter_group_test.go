package polardb_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_polardb_parameter_group", &resource.Sweeper{
		Name: "alicloud_polardb_parameter_group",
		F:    testSweepPolarDBParameterGroup,
	})
}

func testSweepPolarDBParameterGroup(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeParameterGroups"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	var response map[string]interface{}
	polarDBParameterGroupIds := make([]string, 0)
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_polardb_parameter_group", action, errs.AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.ParameterGroups", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.ParameterGroups", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		skip := true
		item := v.(map[string]interface{})
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["ParameterGroupName"])), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping PolarDBParameterGroup Instance: %v", item["ParameterGroupName"])
			continue
		}
		polarDBParameterGroupIds = append(polarDBParameterGroupIds, fmt.Sprint(item["ParameterGroupId"]))
	}

	for _, id := range polarDBParameterGroupIds {
		log.Printf("[INFO] Deleting PolarDBParameterGroup Instance: %s", id)
		deleteAction := "DeleteParameterGroup"
		if err != nil {
			return errs.WrapError(err)
		}
		request = map[string]interface{}{
			"RegionId":         client.RegionId,
			"ParameterGroupId": id,
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(3*time.Minute, func() *retry.RetryError {
			_, err = client.RpcPost("polardb", "2017-08-01", deleteAction, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete PolarDBParameterGroup Instance (%s): %s", polarDBParameterGroupIds, err)
		}
	}
	return nil
}

func TestAccAliCloudPolarDBParameterGroup_basic00(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_parameter_group.default"
	ra := acctest.ResourceAttrInit(resourceId, resourceAliCloudPolarDbParameterGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDBParameterGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sPolarDBParameterGroup-name%d", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAliCloudPolarDbParameterGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":       "tf_testAcc",
					"db_type":    "MySQL",
					"db_version": "8.0",
					"parameters": []map[string]interface{}{
						{
							"param_name":  "wait_timeout",
							"param_value": "86400",
						},
						{
							"param_name":  "innodb_old_blocks_time",
							"param_value": "1000",
						},
					},
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         "tf_testAcc",
						"db_type":      "MySQL",
						"db_version":   "8.0",
						"parameters.#": "2",
						"description":  name,
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

var resourceAliCloudPolarDbParameterGroupMap = map[string]string{}

func resourceAliCloudPolarDbParameterGroupBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

`, name)
}

// Test PolarDb ParameterGroup. >>> Resource test cases, automatically generated.
// Case 属性全覆盖_多参数 11763
func TestAccAliCloudPolarDbParameterGroup_basic11763(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_parameter_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPolarDbParameterGroupMap11763)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbParameterGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPolarDbParameterGroupBasicDependence11763)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_version":           "8.0",
					"parameter_group_name": name,
					"db_type":              "MySQL",
					"parameters": []map[string]interface{}{
						{
							"param_name":  "wait_timeout",
							"param_value": "86400",
						},
						{
							"param_name":  "innodb_old_blocks_time",
							"param_value": "1000",
						},
						{
							"param_name":  "default_time_zone",
							"param_value": "SYSTEM",
						},
					},
					"description": "tf_testAcc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_version":           CHECKSET,
						"parameter_group_name": name,
						"db_type":              "MySQL",
						"parameters.#":         "3",
						"description":          "tf_testAcc",
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

var AlicloudPolarDbParameterGroupMap11763 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPolarDbParameterGroupBasicDependence11763(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test PolarDb ParameterGroup. <<< Resource test cases, automatically generated.
