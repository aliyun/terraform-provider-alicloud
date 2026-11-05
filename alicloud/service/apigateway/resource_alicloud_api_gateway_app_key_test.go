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

// TestAccAliCloudApiGatewayAppKey_update verifies that changing
// secret_reset_token rotates the AppSecret in place via ResetAppSecret,
// while the AppKey identity (and the resource id) stays the same.
func TestAccAliCloudApiGatewayAppKey_update(t *testing.T) {
	resourceId := "alicloud_api_gateway_app_key.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-apigw-appkey-%d", rand)

	var appSecretBefore, appKeyBefore string
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckApiGatewayAppKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayAppKeyUpdate(name, "initial"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceId, "app_id", "alicloud_api_gateway_app.default", "id"),
					resource.TestCheckResourceAttr(resourceId, "secret_reset_token", "initial"),
					resource.TestCheckResourceAttrSet(resourceId, "app_key"),
					resource.TestCheckResourceAttrSet(resourceId, "app_secret"),
					testAccExtractApiGatewayAppKeyAttrs(resourceId, &appKeyBefore, &appSecretBefore),
				),
			},
			{
				// Changing secret_reset_token drives an in-place ResetAppSecret.
				Config: testAccApiGatewayAppKeyUpdate(name, "rotated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceId, "app_id", "alicloud_api_gateway_app.default", "id"),
					resource.TestCheckResourceAttr(resourceId, "secret_reset_token", "rotated"),
					resource.TestCheckResourceAttrSet(resourceId, "app_secret"),
					// The AppKey identity must not change: rotating it would be a ForceNew.
					testAccCheckApiGatewayAppKeyAttr(resourceId, "app_key", appKeyBefore),
					// The AppSecret must differ after rotation.
					testAccCheckApiGatewayAppKeyAttrChanged(resourceId, "app_secret", appSecretBefore),
				),
			},
		},
	})
}

func testAccApiGatewayAppKeyUpdate(name, token string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_api_gateway_app" "default" {
  name        = var.name
  description = "tf_example"
}

resource "alicloud_api_gateway_app_key" "default" {
  app_id            = alicloud_api_gateway_app.default.id
  secret_reset_token = "%s"
}
`, name, token)
}

// testAccExtractApiGatewayAppKeyAttrs captures the current app_key and
// app_secret from state into the provided pointers so a later step can assert
// they changed (or stayed the same) across rotation.
func testAccExtractApiGatewayAppKeyAttrs(id string, appKey, appSecret *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("resource %s not found in state", id)
		}
		*appKey = rs.Primary.Attributes["app_key"]
		*appSecret = rs.Primary.Attributes["app_secret"]
		return nil
	}
}

func testAccCheckApiGatewayAppKeyAttr(id, attr, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("resource %s not found in state", id)
		}
		if got := rs.Primary.Attributes[attr]; got != expected {
			return fmt.Errorf("%s: expected %s = %q, got %q", id, attr, expected, got)
		}
		return nil
	}
}

func testAccCheckApiGatewayAppKeyAttrChanged(id, attr, before string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("resource %s not found in state", id)
		}
		got := rs.Primary.Attributes[attr]
		if got == before {
			return fmt.Errorf("%s: expected %s to change after rotation, but it stayed %q", id, attr, got)
		}
		return nil
	}
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
