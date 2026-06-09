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
)

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_log_config", &resource.Sweeper{
		Name: "alicloud_api_gateway_log_config",
		F:    testSweepApiGatewayLogConfig,
	})
}

func testSweepApiGatewayLogConfig(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}
	action := "DescribeLogConfig"
	request := make(map[string]interface{})
	var response map[string]interface{}
	ApiGatewayLogConfigIds := make([]string, 0)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_api_gateway_log_config", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.LogInfos.LogInfo", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.LogInfos.LogInfo", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		skip := true
		item := v.(map[string]interface{})
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["SlsProject"])), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping ApiGatewayLogConfig Instance: %v", item["SlsProject"])
			continue
		}
		ApiGatewayLogConfigIds = append(ApiGatewayLogConfigIds, fmt.Sprint(item["LogType"]))
	}

	for _, id := range ApiGatewayLogConfigIds {
		log.Printf("[INFO] Deleting ApiGatewayLogConfig Instance: %s", id)
		deleteAction := "DeleteLogConfig"
		if err != nil {
			return WrapError(err)
		}
		request = map[string]interface{}{
			"LogType": id,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			_, err = client.RpcPost("CloudAPI", "2016-07-14", deleteAction, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete ApiGatewayLogConfig Instance (%s): %s", ApiGatewayLogConfigIds, err)
		}
	}
	return nil
}

func TestAccAliCloudApiGatewayLogConfig_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_log_config.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudApiGatewayLogConfigMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayLogConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s-name%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudApiGatewayLogConfigBasicDependence)
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
					"sls_project":   name,
					"sls_log_store": name,
					"log_type":      "PROVIDER",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_project":   name,
						"sls_log_store": name,
						"log_type":      "PROVIDER",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_project":   name + "-update",
					"sls_log_store": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_project":   name + "-update",
						"sls_log_store": name + "-update",
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

var resourceAlicloudApiGatewayLogConfigMap = map[string]string{}

func resourceAlicloudApiGatewayLogConfigBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

`, name)
}

// lintignore: R001
