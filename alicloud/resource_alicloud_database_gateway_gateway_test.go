package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_database_gateway_gateway", &resource.Sweeper{
		Name: "alicloud_database_gateway_gateway",
		F:    testSweepDatabaseGatewayGateway,
	})
}

func testSweepDatabaseGatewayGateway(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.DatabaseGatewaySupportRegions) {
		log.Printf("[INFO] Skipping DatabaseGatewayGateway unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "GetUserGateways"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}

	var response map[string]interface{}
	conn, err := client.NewDgClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-27"), StringPointer("AK"), nil, request, &runtime)
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

		if fmt.Sprint(response["Success"]) == "false" {
			log.Printf("[ERROR] Getting resource %s  failed!!!  Body: %v.", action, err)
			return nil
		}

		m, err := jsonpath.Get("$.Data", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Data", action, err)
			return nil
		}
		v, err := convertJsonStringToList(m.(string))
		for _, v := range v {
			item := v.(map[string]interface{})
			if _, ok := item["gatewayName"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["gatewayName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DatabaseGateway Gateway: %s", item["gatewayName"].(string))
				continue
			}

			action := "DeleteGateway"
			conn, err := client.NewDgClient()
			if err != nil {
				return WrapError(err)
			}
			request := map[string]interface{}{
				"GatewayId": item["gatewayId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-27"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete DatabaseGateway Gateway (%s): %s", item["gatewayId"].(string), err)
			}
			log.Printf("[INFO] Delete DatabaseGateway Gateway success: %s ", item["gatewayId"].(string))
		}

		if len(v) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudDatabaseGatewayGateway_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_database_gateway_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudDatabaseGatewayGatewayMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DgService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDatabaseGatewayGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdatabasegatewaygateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDatabaseGatewayGatewayBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DatabaseGatewaySupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": name,
					"gateway_desc": name + "Desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name": name,
						"gateway_desc": name + "Desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_desc": name + "DescUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_desc": name + "DescUpdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": name,
					"gateway_desc": name + "Desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name": name,
						"gateway_desc": name + "Desc",
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

var AlicloudDatabaseGatewayGatewayMap0 = map[string]string{
	"status": CHECKSET,
}

func AlicloudDatabaseGatewayGatewayBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
