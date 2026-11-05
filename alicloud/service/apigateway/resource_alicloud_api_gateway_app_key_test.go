package apigateway_test

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/apigateway"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_app_key", &resource.Sweeper{
		Name: "alicloud_api_gateway_app_key",
		F:    testSweepApiGatewayAppKey,
	})
}

// testSweepApiGatewayAppKey removes lingering App Keys from tf-testacc Apps so
// the parent App sweeper is not blocked by remaining credentials. It reuses the
// CloudApiService DescribeAppSecurities/DeleteAppKey paths used by the resource.
func testSweepApiGatewayAppKey(region string) error {
	if sweep.SweepPreCheckWithRegions(region, false, connectivity.ApiGatewayNoSupportedRegions) {
		log.Printf("[INFO] Skipping API Gateway unsupported region: %s", region)
		return nil
	}
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	cloudApiService := apigateway.NewCloudApiService(client)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	req := cloudapi.CreateDescribeAppAttributesRequest()
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeAppAttributes(req)
	})
	if err != nil {
		return fmt.Errorf("Error Describe Apps: %s", err)
	}
	apps, _ := raw.(*cloudapi.DescribeAppAttributesResponse)

	for _, v := range apps.Apps.AppAttribute {
		name := v.AppName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping app: %s", name)
			continue
		}
		appId := v.AppId
		securities, err := cloudApiService.DescribeAppSecurities(appId)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			log.Printf("[ERROR] DescribeAppSecurities for App %d: %s", appId, err)
			continue
		}
		for _, m := range securities {
			appKey := fmt.Sprint(m["AppKey"])
			if err := cloudApiService.DeleteAppKey(strconv.FormatInt(appId, 10), appKey); err != nil {
				log.Printf("[ERROR] Failed to delete AppKey %d:%s: %s", appId, appKey, err)
			}
		}
	}
	return nil
}

func TestAccAliCloudApiGatewayAppKey_basic(t *testing.T) {
	resourceId := "alicloud_api_gateway_app_key.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-apigw-appkey-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckApiGatewayAppKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayAppKeyBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceId, "app_id", "alicloud_api_gateway_app.default", "id"),
					resource.TestCheckResourceAttrSet(resourceId, "app_key"),
					resource.TestCheckResourceAttrSet(resourceId, "app_secret"),
					resource.TestCheckResourceAttrSet(resourceId, "app_code"),
					resource.TestCheckResourceAttrSet(resourceId, "created_time"),
					resource.TestCheckResourceAttrSet(resourceId, "modified_time"),
					resource.TestMatchResourceAttr(resourceId, "id", regexp.MustCompile(`^\d+:\S+$`)),
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

func testAccApiGatewayAppKeyBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_api_gateway_app" "default" {
  name        = var.name
  description = "tf_example"
}

resource "alicloud_api_gateway_app_key" "default" {
  app_id = alicloud_api_gateway_app.default.id
}
`, name)
}

func testAccCheckApiGatewayAppKeyDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	cloudApiService := apigateway.NewCloudApiService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_api_gateway_app_key" {
			continue
		}
		parts := strings.Split(rs.Primary.ID, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid alicloud_api_gateway_app_key id %q", rs.Primary.ID)
		}
		appId, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid app_id in alicloud_api_gateway_app_key id %q: %s", rs.Primary.ID, err)
		}
		if _, err := cloudApiService.FindAppSecurity(appId, parts[1]); err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing app key %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("app key %s still exists", rs.Primary.ID)
	}
	return nil
}
