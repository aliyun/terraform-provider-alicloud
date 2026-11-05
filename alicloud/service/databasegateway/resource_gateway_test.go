package databasegateway_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/databasegateway"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var databasegatewayTestEndpoints sync.Map
var databasegatewayTestSignVersion sync.Map

// databasegatewayTestClientForRegion rebuilds the package-level sweeper client on the
// test account, mirroring the root suite's sharedClientForRegion. It returns a typed
// *connectivity.AliyunClient so the sweeper does not need a type assertion.
func databasegatewayTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	var accessKey, secretKey string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}
	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}
	securityToken := os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &databasegatewayTestEndpoints,
		SignVersion: &databasegatewayTestSignVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential

	return conf.Client()
}

func init() {
	resource.AddTestSweepers("alicloud_database_gateway_gateway", &resource.Sweeper{
		Name: "alicloud_database_gateway_gateway",
		F:    testSweepDatabaseGatewayGateway,
	})
}

func testSweepDatabaseGatewayGateway(region string) error {
	if sweep.SweepPreCheckWithRegions(region, true, connectivity.DatabaseGatewaySupportRegions) {
		log.Printf("[INFO] Skipping DatabaseGatewayGateway unsupported region: %s", region)
		return nil
	}
	client, err := databasegatewayTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "GetUserGateways"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
	}

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("dg", "2019-03-27", action, nil, request, true)
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

		if fmt.Sprint(response["Success"]) == "false" {
			log.Printf("[ERROR] Getting resource %s  failed!!!  Body: %v.", action, err)
			return nil
		}

		m, err := jsonpath.Get("$.Data", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Data", action, err)
			return nil
		}
		v, err := convert.ConvertJsonStringToList(m.(string))
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
			request := map[string]interface{}{
				"GatewayId": item["gatewayId"],
			}
			_, err = client.RpcPost("dg", "2019-03-27", action, nil, request, false)
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

func TestAccAliCloudDatabaseGatewayGateway_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_database_gateway_gateway.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDatabaseGatewayGatewayMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return databasegateway.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDatabaseGatewayGateway")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdatabasegatewaygateway%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDatabaseGatewayGatewayBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DatabaseGatewaySupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	"status": acctest.CHECKSET,
}

func AlicloudDatabaseGatewayGatewayBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
