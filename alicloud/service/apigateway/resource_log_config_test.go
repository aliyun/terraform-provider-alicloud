package apigateway_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/apigateway"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_log_config", &resource.Sweeper{
		Name: "alicloud_api_gateway_log_config",
		F:    testSweepApiGatewayLogConfig,
	})
}

func testSweepApiGatewayLogConfig(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_api_gateway_log_config", action, errs.AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.LogInfos.LogInfo", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.LogInfos.LogInfo", response)
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
			return errs.WrapError(err)
		}
		request = map[string]interface{}{
			"LogType": id,
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(3*time.Minute, func() *retry.RetryError {
			_, err = client.RpcPost("CloudAPI", "2016-07-14", deleteAction, nil, request, false)
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
			log.Printf("[ERROR] Failed to delete ApiGatewayLogConfig Instance (%s): %s", ApiGatewayLogConfigIds, err)
		}
	}
	return nil
}

func TestAccAliCloudApiGatewayLogConfig_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_log_config.default"
	ra := acctest.ResourceAttrInit(resourceId, ResourceLogConfigMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewCloudApiService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayLogConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s-name%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, ResourceLogConfigBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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

var ResourceLogConfigMap = map[string]string{}

func ResourceLogConfigBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

`, name)
}

// lintignore: R001
